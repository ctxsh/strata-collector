package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	typesv1 "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceOpts struct {
	Cache   cache.Cache
	Client  client.Client
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Service struct {
	cache    cache.Cache
	client   client.Client
	enabled  bool
	interval time.Duration
	logger   logr.Logger
	metrics  *strata.Metrics
	obj      *v1beta1.Discovery
	prefix   string
	selector metav1.LabelSelector

	sendChans map[string]chan<- resource.Resource
	stopChan  chan struct{}
	stopOnce  sync.Once
	sync.Mutex
}

func NewService(obj *v1beta1.Discovery, opts *ServiceOpts) *Service {
	return &Service{
		cache:     opts.Cache,
		client:    opts.Client,
		enabled:   *obj.Spec.Enabled,
		interval:  time.Duration(*obj.Spec.IntervalSeconds) * time.Second,
		logger:    opts.Logger,
		metrics:   opts.Metrics,
		obj:       obj,
		prefix:    *obj.Spec.Prefix,
		selector:  obj.Spec.Selector,
		sendChans: make(map[string]chan<- resource.Resource),
		stopChan:  make(chan struct{}),
	}
}

func (s *Service) AddChan(key string, ch chan<- resource.Resource) {
	s.Lock()
	defer s.Unlock()

	s.sendChans[key] = ch
}

func (s *Service) Start() {
	// TODO: manage context better
	go s.start(context.Background())
}

func (s *Service) Stop() {
	s.Lock()
	defer s.Unlock()

	s.stopOnce.Do(func() {
		close(s.stopChan)
	})
}

func (s *Service) start(ctx context.Context) {
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

func (s *Service) intervalRun(ctx context.Context) {
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
	s.send(resources)
	s.updateStatus(ctx)
}

func (s *Service) discover(ctx context.Context) []resource.Resource {
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

	return resources
}

func (s *Service) send(resources []resource.Resource) {
	s.Lock()
	defer s.Unlock()

	for n, sendChan := range s.sendChans {
		if sendChan == nil {
			s.logger.V(8).Info("send channel is nil", "collector", n)
			continue
		}
		s.logger.V(8).Info("sending resources", "collector", n)
		// TODO: make me async
		for i := 0; i < len(resources); i++ {
			sendChan <- resources[i]
		}
	}
}

// discoverPods lists all pods that match the selector and if the scrape annotation
// is set to true, will create the collection resource and send it to the collector.
func (s *Service) discoverPods(ctx context.Context, res *[]resource.Resource) error {
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
func (s *Service) discoverServices(ctx context.Context, res *[]resource.Resource) error {
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
	err := s.cache.Get(ctx, typesv1.NamespacedName{
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

func (s *Service) updateStatus(ctx context.Context) error {
	s.logger.V(8).Info("updating discovery status")

	obj := s.obj.DeepCopy()
	obj.Status = v1beta1.DiscoveryStatus{
		Active:         s.enabled,
		LastDiscovered: metav1.Now(),
		Ready:          fmt.Sprintf("%d/%d", len(s.sendChans), len(s.obj.Spec.Collectors)),
	}

	return s.client.Update(ctx, obj)
}
