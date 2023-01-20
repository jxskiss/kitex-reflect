package kitexreflect

import (
	"context"
	"errors"
	"fmt"
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
	kcOptions []client.Option,
	providerOpts ...Option,
) (*ProviderImpl, error) {
	cli, err := svc.NewClient(serviceName, kcOptions...)
	if err != nil {
		return nil, err
	}

	impl := &ProviderImpl{
		debounceInterval: defaultDebounceInterval,
		errorHandler:     defaultErrorHandler,
		cli:              cli,
		serviceName:      serviceName,
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
	} else {
		notify := make(chan *descriptor.ServiceDescriptor, 1)
		placeholderDesc := newPlaceholderDescriptor(serviceName, notify)
		impl.updates <- placeholderDesc
		impl.updating = 1
		go impl.asyncInitialize(ctx, notify)
	}

	if impl.autoUpdateInterval > 0 {
		go impl.autoUpdate()
	}

	return impl, nil
}

func (p *ProviderImpl) asyncInitialize(ctx context.Context, ready chan *descriptor.ServiceDescriptor) {
	defer atomic.StoreInt64(&p.updating, 0)

	const (
		errSleep = 200 * time.Millisecond
		maxSleep = time.Minute
	)

	sleep := errSleep
	for {
		if p.IsClosed() {
			close(ready)
			return
		}
		desc, version, err := doRequestDescriptor(ctx, p.cli, "")
		if err == nil {
			p.setNewDescriptor(desc, version)
			ready <- desc
			close(ready)
			return
		}
		errMsg := fmt.Sprintf("failed to init descriptor for service %s", p.serviceName)
		p.errorHandler(err, errMsg)
		sleep = min(2*sleep, maxSleep)
		time.Sleep(sleep)
	}
}

func (p *ProviderImpl) autoUpdate() {
	ticker := time.NewTicker(p.autoUpdateInterval)
	defer ticker.Stop()

	ctx := context.Background()
	for range ticker.C {
		if p.IsClosed() {
			return
		}
		p.Update(ctx)
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

	initAsync          bool
	autoUpdateInterval time.Duration

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
		sleep    = 300 * time.Millisecond
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

func newPlaceholderDescriptor(serviceName string, notify <-chan *descriptor.ServiceDescriptor) *descriptor.ServiceDescriptor {
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

	notify   <-chan *descriptor.ServiceDescriptor
	realDesc atomic.Value // *descriptor.ServiceDescriptor
}

func (p *placeholderRouter) Handle(_ descriptor.Route) {
}

func (p *placeholderRouter) Lookup(req *descriptor.HTTPRequest) (*descriptor.FunctionDescriptor, error) {
	realDesc, _ := p.realDesc.Load().(*descriptor.ServiceDescriptor)
	if realDesc != nil {
		return realDesc.Router.Lookup(req)
	}

	const timeout = 500 * time.Millisecond
	select {
	case realDesc = <-p.notify:
		if realDesc != nil {
			p.realDesc.Store(realDesc)
		} else {
			time.Sleep(time.Millisecond)
		}
	case <-time.After(timeout):
	}

	if realDesc == nil {
		realDesc, _ = p.realDesc.Load().(*descriptor.ServiceDescriptor)
	}
	if realDesc == nil {
		return nil, fmt.Errorf("descriptor for service %s is not ready after %v",
			p.serviceName, time.Since(p.initTime))
	}
	return realDesc.Router.Lookup(req)
}

//nolint:unused
func isUnknownMethodError(err error) bool {
	var transErr *remote.TransError
	if errors.As(err, &transErr) {
		return transErr.TypeID() == remote.UnknownMethod
	}
	return false
}

func min[T ~int | ~int32 | ~int64](a, b T) T {
	if a < b {
		return a
	}
	return b
}
