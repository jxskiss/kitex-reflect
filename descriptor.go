package kitexreflect

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/remote"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
	svc "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl/reflectionservice"
)

// NewDescriptorProvider creates a ProviderImpl which implements [generic.DescriptorProvider].
func NewDescriptorProvider(
	ctx context.Context,
	serviceName string,
	kclientOptions []client.Option,
	providerOpts ...Option,
) (*ProviderImpl, error) {
	cli, err := svc.NewClient(serviceName, kclientOptions...)
	if err != nil {
		return nil, err
	}

	impl := &ProviderImpl{
		debounceInterval: defaultDebounceInterval,
		errorHandler:     defaultErrorHandler,
		cli:              cli,
		serviceName:      serviceName,
		updating:         1,
		updates:          make(chan *descriptor.ServiceDescriptor, 1),
	}
	for _, opt := range providerOpts {
		if opt.applyDescProv != nil {
			opt.applyDescProv(impl)
		}
	}

	if !impl.initAsync {
		desc, version, err := doRequestDescriptor(ctx, impl.cli, "")
		if err != nil {
			return nil, err
		}
		impl.setNewDescriptor(desc, version)
		impl.updating = 0
	} else {
		notify := make(chan descriptor.Router, 1)
		placeholderDesc := newPlaceholderDescriptor(serviceName, notify)
		impl.updates <- placeholderDesc
		go impl.asyncInitialize(ctx, notify)
	}

	return impl, nil
}

func (p *ProviderImpl) asyncInitialize(ctx context.Context, ready chan descriptor.Router) {
	defer atomic.StoreInt64(&p.updating, 0)
	const (
		errSleep           = 300 * time.Millisecond
		unknownMethodSleep = time.Second
	)
	for {
		if p.IsClosed() {
			close(ready)
			return
		}
		desc, version, err := doRequestDescriptor(ctx, p.cli, "")
		if err == nil {
			p.setNewDescriptor(desc, version)
			ready <- desc.Router
			close(ready)
			return
		}
		errMsg := fmt.Sprintf("failed to init descriptor for service %s", p.serviceName)
		p.errorHandler(err, errMsg)
		sleep := errSleep
		if isUnknownMethodError(err) {
			sleep = unknownMethodSleep
		}
		time.Sleep(sleep)
	}
}

var _ generic.DescriptorProvider = &ProviderImpl{}

// ProviderImpl connects to a Kitex service which has reflection support,
// it calls the service's method ReflectService to build service descriptor.
// It implements [generic.DescriptorProvider].
type ProviderImpl struct {
	debounceInterval time.Duration
	errorHandler     func(err error, msg string)

	cli         svc.Client
	serviceName string

	initAsync bool

	preUpdateSec int64
	updating     int64

	closeMu sync.Mutex
	closed  int64
	updates chan *descriptor.ServiceDescriptor
	version string
}

// Close closes the provider and the channel returned by Provide.
func (p *ProviderImpl) Close() error {
	p.closeMu.Lock()
	if atomic.CompareAndSwapInt64(&p.closed, 0, 1) {
		close(p.updates)
	}
	p.closeMu.Unlock()
	return nil
}

// IsClosed checks whether the provider has been closed.
func (p *ProviderImpl) IsClosed() bool {
	return atomic.LoadInt64(&p.closed) > 0
}

// Provide returns a channel for service descriptors.
func (p *ProviderImpl) Provide() <-chan *descriptor.ServiceDescriptor {
	return p.updates
}

// Update triggers updating service descriptor from backend service.
// The actual rpc requests will be debounced according to debounceInterval.
// This method is safe to call concurrently.
func (p *ProviderImpl) Update(ctx context.Context) {
	if p.IsClosed() {
		return
	}
	preUpdateTime := time.Unix(atomic.LoadInt64(&p.preUpdateSec), 0)
	if time.Since(preUpdateTime) < p.debounceInterval {
		return
	}

	if atomic.CompareAndSwapInt64(&p.updating, 0, 1) {
		go p.doUpdate(ctx)
	}
}

func (p *ProviderImpl) doUpdate(ctx context.Context) {
	defer atomic.StoreInt64(&p.updating, 0)

	const (
		retryCnt = 2
		sleep    = 100 * time.Millisecond
	)
	var (
		desc    *descriptor.ServiceDescriptor
		payload *idl.ReflectServiceRespPayload
	)
	for i := 0; i <= retryCnt; i++ {
		if i > 0 {
			time.Sleep(sleep)
		}
		req := newReflectServiceRequest(p.version)
		resp, err := p.cli.ReflectService(ctx, req)
		if err != nil {
			msg := fmt.Sprintf("failed to call %s.ReflectService", p.serviceName)
			p.errorHandler(err, msg)
			continue
		}
		payload, err = idl.UnmarshalReflectServiceRespPayload(resp.Payload)
		if err != nil {
			msg := fmt.Sprintf("failed to unmarshal response payload: %v", err)
			p.errorHandler(err, msg)
			return
		}
		if payload.Version == "" || payload.Version == p.version {
			// The IDL is not changed, nothing to update.
			return
		}

		desc, err = BuildServiceDescriptor(payload)
		if err != nil {
			msg := fmt.Sprintf("failed to build descriptor for service %s", p.serviceName)
			p.errorHandler(err, msg)
			return
		}
	}

	if desc != nil {
		p.setNewDescriptor(desc, payload.Version)
	}
	return
}

func (p *ProviderImpl) setNewDescriptor(desc *descriptor.ServiceDescriptor, version string) {
	p.closeMu.Lock()
	defer p.closeMu.Unlock()
	if p.IsClosed() {
		return
	}
	select {
	case <-p.updates:
	default:
	}
	p.updates <- desc
	p.version = version
	atomic.StoreInt64(&p.preUpdateSec, time.Now().Unix())
}

func newPlaceholderDescriptor(serviceName string, notify <-chan descriptor.Router) *descriptor.ServiceDescriptor {
	desc := &descriptor.ServiceDescriptor{
		Name:      "KitexReflectPlaceholderService",
		Functions: nil,
		Router: &placeholderRouter{
			serviceName: serviceName,
			initTime:    time.Now(),
			notify:      notify,
		},
	}
	return desc
}

type placeholderRouter struct {
	serviceName string
	initTime    time.Time

	notify <-chan descriptor.Router
	real   atomic.Value // descriptor.Router
}

func (r *placeholderRouter) Handle(_ descriptor.Route) {
}

func (r *placeholderRouter) Lookup(req *descriptor.HTTPRequest) (*descriptor.FunctionDescriptor, error) {
	const timeout = time.Second
	var realr descriptor.Router
	select {
	case realr = <-r.notify:
		if realr != nil {
			r.real.Store(realr)
		} else {
			realr, _ = r.real.Load().(descriptor.Router)
			for realr == nil {
				runtime.Gosched()
				realr, _ = r.real.Load().(descriptor.Router)
			}
		}
	case <-time.After(timeout):
		return nil, fmt.Errorf("descriptor for service %s is not ready after %v",
			r.serviceName, time.Since(r.initTime))
	}
	return realr.Lookup(req)
}

func isUnknownMethodError(err error) bool {
	var transErr *remote.TransError
	if errors.As(err, &transErr) {
		return transErr.TypeID() == remote.UnknownMethod
	}
	return false
}
