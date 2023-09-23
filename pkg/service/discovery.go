package service

import (
	"context"
	"sync"
	"time"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DiscoveryOpts struct {
	Cache    cache.Cache
	Client   client.Client
	Logger   logr.Logger
	Metrics  *strata.Metrics
	Registry *Registry
}

type Discovery struct {
	cache    cache.Cache
	client   client.Client
	registry *Registry
	enabled  bool
	interval time.Duration
	logger   logr.Logger
	metrics  *strata.Metrics
	obj      *v1beta1.Discovery
	prefix   string
	selector metav1.LabelSelector
	stopChan chan struct{}
	stopOnce sync.Once
	sync.Mutex
}

func NewDiscovery(obj *v1beta1.Discovery, opts *DiscoveryOpts) *Discovery {
	return &Discovery{
		cache:    opts.Cache,
		client:   opts.Client,
		registry: opts.Registry,
		enabled:  *obj.Spec.Enabled,
		interval: time.Duration(*obj.Spec.IntervalSeconds) * time.Second,
		logger:   opts.Logger,
		metrics:  opts.Metrics,
		obj:      obj,
		prefix:   *obj.Spec.Prefix,
		selector: obj.Spec.Selector,
		stopChan: make(chan struct{}),
	}
}

func (s *Discovery) Start() {
	// TODO: manage context better
	go s.start(context.Background())
}

func (s *Discovery) Stop() {
	s.Lock()
	defer s.Unlock()

	s.stopOnce.Do(func() {
		close(s.stopChan)
	})
}

func (s *Discovery) start(ctx context.Context) {
	s.logger.Info("starting discovery service")
	s.intervalRun(ctx)

	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-s.stopChan:
			s.logger.V(8).Info("shutting down discovery service")
			return
		case <-ticker.C:
			s.logger.V(8).Info("running discovery")
			s.intervalRun(ctx)
		}
	}
}

func (s *Discovery) intervalRun(ctx context.Context) {
	// TODO: make this into a goroutine and add a mutex that represents
	// that discovery is running.  If it is running, then don't start another
	// discovery, but instead emit a metric to let the operator know that the
	// run is being skipped.  This will help inform the operator that there are
	// too many resources to monitor and that they should consider reducing.
	// The other option would be to come back and allow multiple discovery workers
	// but there's some complexity in dealing with the controller cache that I don't
	// think that we want to tackle right now.  That being said, we probably won't
	// actually hit this case since we are only interacting with the cache/api or the
	// service is blocked on the channel meaning that the collector workers are not
	// ablel to keep up.  None of which are solved by adding more discovery workers.
	resources := s.discover(ctx)
	ready, inFlight := s.send(ctx, resources)
	err := s.updateStatus(ctx, len(resources), ready, inFlight)
	if err != nil {
		s.logger.Error(err, "unable to update status")
	}
}

func (s *Discovery) discover(ctx context.Context) []resource.Resource {
	s.Lock()
	defer s.Unlock()

	resources := make([]resource.Resource, 0)

	s.logger.V(8).Info("discovering resources", "spec", s.obj.Spec)

	if *s.obj.Spec.Resources.Pods {
		s.logger.V(8).Info("discovering pods")
		if err := s.discoverPods(ctx, &resources); err != nil {
			s.logger.Error(err, "unable to discover pods")
		}
	}

	if *s.obj.Spec.Resources.Services {
		s.logger.V(8).Info("discovering services")
		if err := s.discoverServices(ctx, &resources); err != nil {
			s.logger.Error(err, "unable to discover services")
		}
	}

	return resources
}

func (s *Discovery) send(ctx context.Context, resources []resource.Resource) (int, int) {
	// TODO: look at some of the race conditions between getting the send channels
	// and sending the resources.  There's a chance that the send channel may not
	// exist because of a collector deletion, so we should probably make sure that
	// we are checking in the finalizers for the collector and block until the the
	// send has completed.
	s.Lock()
	defer s.Unlock()

	var ready int
	var inFlight int

	for _, objRef := range s.obj.Spec.Collectors {
		nn := types.NamespacedName{
			Namespace: objRef.Namespace,
			Name:      objRef.Name,
		}

		// Grab the collector from the cache - don't know if need this anymore - we can
		// send this over to the collector and only worry about whether or not the collector
		// is in the registry.
		var collector v1beta1.Collector
		err := s.cache.Get(ctx, nn, &collector)
		if err != nil {
			if client.IgnoreNotFound(err) == nil {
				s.logger.V(8).Info("collector not found", "collector", nn)
			} else {
				s.logger.Error(err, "error retrieving collector", "collector", nn)
			}
			continue
		}

		// If the collector is not enabled, then skip it.
		if !*collector.Spec.Enabled {
			continue
		}

		s.logger.V(8).Info("collector found, sending resources", "collector", nn)

		// Start the fun stuff
		ready++
		err = s.registry.SendResources(nn, resources)
		if err != nil {
			continue
		}

		s.logger.V(8).Info("resources sent", "collector", nn)

		i, err := s.registry.GetInFlightResources(nn)
		if err != nil {
			continue
		}

		inFlight += i
	}

	return ready, inFlight
}

// discoverPods lists all pods that match the selector and if the scrape annotation
// is set to true, will create the collection resource and send it to the collector.
func (s *Discovery) discoverPods(ctx context.Context, res *[]resource.Resource) error {
	var list corev1.PodList

	// TODO: how do we handle MatchExpressions?
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(s.selector.MatchLabels),
	}

	err := s.cache.List(ctx, &list, opts)
	if err != nil {
		return err
	}

	for _, pod := range list.Items {
		// TODO: configurable prefix for annotations
		cr := resource.New(pod.GetAnnotations(), s.prefix)
		if !cr.Scrape {
			continue
		}

		s.logger.V(8).Info("pod found", "obj", pod.ObjectMeta)

		cr = cr.WithMetadata(pod.DeepCopy()).
			WithIP(pod.Status.PodIP).
			WithAnnotations(pod.Annotations).
			WithLabels(pod.Labels)
		*res = append(*res, *cr)
	}

	return nil
}

// discoverServices lists all services that match the selector and if the scrape
// annotation is set to true, will create the collection resource and send it to
// the collector.  If the service is a headless service, then it will discover
// the endpoints and create the collection resource for each endpoint.
func (s *Discovery) discoverServices(ctx context.Context, res *[]resource.Resource) error {
	var list corev1.ServiceList
	err := s.cache.List(ctx, &list, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(s.selector.MatchLabels),
	})
	if err != nil {
		return err
	}

	for _, svc := range list.Items {
		// TODO: configurable prefix for annotations
		cr := resource.New(svc.Annotations, s.prefix)
		if !cr.Scrape {
			continue
		}

		if svc.Spec.ClusterIP == "None" {
			if !*s.obj.Spec.Resources.Pods {
				return nil
			}

			s.logger.V(8).Info("headless service encountered, discovering endpoints")
			return s.discoverEndpoints(ctx, svc, res)
		}

		s.logger.V(8).Info("service found", "obj", svc.ObjectMeta)
		cr = cr.WithMetadata(svc.DeepCopy()).
			WithIP(svc.Spec.ClusterIP).
			WithAnnotations(svc.Annotations).
			WithLabels(svc.Labels)
		*res = append(*res, *cr)
	}

	return nil
}

// discoverEndpoints lists all endpoints that match the selector and sends the
// collection resource to the collector.  We don't need to check the scrape annotation
// since that is handled by the service discovery.
// TODO: instead of passing the service object along, create a metadata struct which
// will have all the info and can be attached to the resource.
func (s *Discovery) discoverEndpoints(ctx context.Context, svc corev1.Service, res *[]resource.Resource) error {
	var endpoints corev1.Endpoints
	err := s.cache.Get(ctx, types.NamespacedName{
		Namespace: svc.GetNamespace(),
		Name:      svc.GetName(),
	}, &endpoints, &client.GetOptions{})
	if err != nil {
		return err
	}

	for _, sset := range endpoints.Subsets {
		for _, addr := range sset.Addresses {
			cr := resource.New(svc.Annotations, s.prefix)
			// We're not checking for the scrape condition here as that we are using
			// the parent service as the authority for this and it's already been checked.
			s.logger.V(8).Info("pod found", "obj", svc.ObjectMeta, "ip", addr.IP)
			cr = cr.WithMetadataRef(addr.TargetRef).
				WithIP(addr.IP).
				WithAnnotations(svc.Annotations).
				WithLabels(svc.Labels)
			*res = append(*res, *cr)
		}
	}

	return nil
}

func (s *Discovery) updateStatus(ctx context.Context, count int, ready int, inFlight int) error {
	s.Lock()
	defer s.Unlock()

	var obj v1beta1.Discovery
	err := s.cache.Get(ctx, types.NamespacedName{
		Namespace: s.obj.GetNamespace(),
		Name:      s.obj.GetName(),
	}, &obj)
	if err != nil {
		return err
	}

	obj.Status = v1beta1.DiscoveryStatus{
		Active:                   s.enabled,
		LastDiscovered:           metav1.Now(),
		ReadyCollectors:          ready,
		TotalCollectors:          len(s.obj.Spec.Collectors),
		DiscoveredResourcesCount: count,
		InFlightResources:        inFlight,
	}

	s.logger.V(8).Info("updating discovery status", "status", obj.Status)

	return s.client.Status().Update(ctx, &obj)
}
