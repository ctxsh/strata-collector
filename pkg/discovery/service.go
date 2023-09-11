package discovery

import (
	"context"
	"sync"
	"time"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	typesv1 "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceOpts struct {
	Client          client.Client
	Enabled         bool
	IntervalSeconds int64
	Prefix          string
	Selector        metav1.LabelSelector
	Logger          logr.Logger
	Metrics         *strata.Metrics
}

type Service struct {
	namespacedName typesv1.NamespacedName
	client         client.Client
	enabled        bool
	interval       time.Duration
	prefix         string
	selector       metav1.LabelSelector
	logger         logr.Logger
	metrics        *strata.Metrics

	stopChan chan struct{}
	stopOnce sync.Once
	sync.Mutex
}

func NewService(namespace, name string, opts *ServiceOpts) *Service {
	return &Service{
		namespacedName: typesv1.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
		client:   opts.Client,
		enabled:  opts.Enabled,
		interval: time.Duration(opts.IntervalSeconds) * time.Second,
		prefix:   opts.Prefix,
		selector: opts.Selector,
		logger:   opts.Logger,
		metrics:  opts.Metrics,
		stopChan: make(chan struct{}),
	}
}

func (s *Service) NamespacedName() typesv1.NamespacedName {
	return s.namespacedName
}

func (s *Service) Start(sendChan chan<- resource.Resource) {
	s.logger.Info("starting discovery service")
	// TODO: manage context better
	go s.start(context.Background(), sendChan)
}

func (s *Service) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopChan)
	})
}

func (s *Service) start(ctx context.Context, sendChan chan<- resource.Resource) {
	// Initial discovery run
	s.discover(ctx, sendChan)

	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-s.stopChan:
			s.logger.V(8).Info("worker received stop")
			return
		case <-ticker.C:
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
			s.discover(ctx, sendChan)
		}
	}
}

func (s *Service) discover(ctx context.Context, sendChan chan<- resource.Resource) {
	s.Lock()
	defer s.Unlock()

	resources := make([]resource.Resource, 0)

	s.logger.V(8).Info("discovering pods")
	err := s.discoverPods(ctx, &resources)
	if err != nil {
		s.logger.Error(err, "unable to discover pods")
	}

	s.logger.V(8).Info("discovering services")
	err = s.discoverServices(ctx, &resources)
	if err != nil {
		s.logger.Error(err, "unable to discover services")
	}

	// TODO: I think the discover functions should probably return a list of
	// resources that were discovered and then we can send them to the collector
	// at the end.  This will most likely help us manage metrics and state better.
	// the only thing I'm not a fan of is the memory allocation every iteration
	// which could impact GC, but I think that's a tradeoff that we can make initially.

	// TODO: send metrics to the collector
	if len(resources) > 0 {
		s.logger.V(8).Info("sending resources to collector", "count", len(resources))
		for _, resource := range resources {
			sendChan <- resource
		}
	}

	// TODO: update the status of the discovery service
	s.logger.V(8).Info("finished discovery run", "discovered", len(resources))
}

// discoverPods lists all pods that match the selector and if the scrape annotation
// is set to true, will create the collection resource and send it to the collector.
func (s *Service) discoverPods(ctx context.Context, res *[]resource.Resource) error {
	var list corev1.PodList

	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(s.selector.MatchLabels),
	}

	err := s.client.List(ctx, &list, opts)
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
func (s *Service) discoverServices(ctx context.Context, res *[]resource.Resource) error {
	var list corev1.ServiceList
	err := s.client.List(ctx, &list, &client.ListOptions{
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
func (s *Service) discoverEndpoints(ctx context.Context, svc corev1.Service, res *[]resource.Resource) error {
	var endpoints corev1.Endpoints
	err := s.client.Get(ctx, typesv1.NamespacedName{
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
