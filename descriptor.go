package kitexreflect

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/generic/thrift"
	"github.com/cloudwego/thriftgo/generator/golang/extension/meta"
	"github.com/cloudwego/thriftgo/parser"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
	svc "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl/reflectionservice"
)

const defaultDebounceInterval = 10 * time.Second

func defaultErrorHandler(err error, msg string) {
	log.Printf("ERROR: %s: %v", msg, err)
}

// NewDescriptorProvider creates a ProviderImpl which implements [generic.DescriptorProvider].
func NewDescriptorProvider(ctx context.Context, serviceName string, opts ...client.Option) (*ProviderImpl, error) {
	cli, err := svc.NewClient(serviceName, opts...)
	if err != nil {
		return nil, err
	}
	firstReq := newReflectServiceRequest("")
	firstResp, err := cli.ReflectService(ctx, firstReq)
	if err != nil {
		return nil, err
	}
	respPayload, err := idl.UnmarshalReflectServiceRespPayload(firstResp.Payload)
	if err != nil {
		return nil, err
	}
	desc, err := BuildServiceDescriptor(respPayload)
	if err != nil {
		return nil, err
	}

	impl := &ProviderImpl{
		DebounceInterval: defaultDebounceInterval,
		ErrorHandler:     defaultErrorHandler,
		cli:              cli,
		serviceName:      serviceName,
		updates:          make(chan *descriptor.ServiceDescriptor, 1),
		version:          respPayload.Version,
	}
	impl.updates <- desc
	return impl, nil
}

var _ generic.DescriptorProvider = &ProviderImpl{}

// ProviderImpl connects to a Kitex service which has reflection support,
// it calls the service's method ReflectService to build service descriptor.
// It implements [generic.DescriptorProvider].
type ProviderImpl struct {

	// DebounceInterval sets max interval to debounce requests
	// sent to backend service. The default is 10 seconds.
	DebounceInterval time.Duration

	// ErrorHandler optionally specifies an error handler function
	// for errors happened when updating service descriptor from
	// backend service.
	// By default, it logs error using the [log] package.
	ErrorHandler func(err error, msg string)

	cli         svc.Client
	serviceName string

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
// The actual rpc requests will be debounced according to DebounceInterval.
// This method is safe to call concurrently.
func (p *ProviderImpl) Update(ctx context.Context) {
	if p.IsClosed() {
		return
	}
	preUpdateTime := time.Unix(atomic.LoadInt64(&p.preUpdateSec), 0)
	if time.Since(preUpdateTime) < p.DebounceInterval {
		return
	}

	if atomic.CompareAndSwapInt64(&p.updating, 0, 1) {
		atomic.StoreInt64(&p.preUpdateSec, time.Now().Unix())
		go p.doUpdate(ctx)
	}
}

func (p *ProviderImpl) doUpdate(ctx context.Context) {
	defer atomic.StoreInt64(&p.updating, 0)

	req := newReflectServiceRequest(p.version)
	resp, err := p.cli.ReflectService(ctx, req)
	if err != nil {
		msg := fmt.Sprintf("failed to call %s.ReflectService", p.serviceName)
		p.ErrorHandler(err, msg)
		return
	}
	payload, err := idl.UnmarshalReflectServiceRespPayload(resp.Payload)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal response payload: %v", err)
		p.ErrorHandler(err, msg)
		return
	}
	if payload.Version == "" || payload.Version == p.version {
		// The IDL is not changed.
		return
	}

	desc, err := BuildServiceDescriptor(payload)
	if err != nil {
		msg := fmt.Sprintf("failed to build descriptor for service %s", p.serviceName)
		p.ErrorHandler(err, msg)
		return
	}

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
	p.version = payload.Version
}

func newReflectServiceRequest(existingIDLVersion string) *ReflectServiceRequest {
	payload := &idl.ReflectServiceReqPayload{
		ExistingIDLVersion: existingIDLVersion,
	}
	payloadBuf, _ := idl.MarshalReflectServiceReqPayload(payload)
	return &idl.ReflectServiceRequest{
		Payload: payloadBuf,
	}
}

// BuildServiceDescriptor builds a [descriptor.ServiceDescriptor] from a ReflectServiceResponse.
func BuildServiceDescriptor(payload *idl.ReflectServiceRespPayload) (*descriptor.ServiceDescriptor, error) {
	if len(payload.IDL) == 0 {
		return nil, fmt.Errorf("IDL bytes is empty")
	}
	builder := &descriptorBuilder{
		payload: payload,
	}
	return builder.Build()
}

type descriptorBuilder struct {
	payload *ReflectServiceRespPayload
}

func (p *descriptorBuilder) Build() (*descriptor.ServiceDescriptor, error) {
	gzr, _ := gzip.NewReader(bytes.NewBuffer(p.payload.IDL))
	rawBuf, err := io.ReadAll(gzr)
	if err != nil {
		return nil, fmt.Errorf("cannot decompress IDL bytes: %w", err)
	}

	ast := &parser.Thrift{}
	err = meta.Unmarshal(rawBuf, ast)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal parser.Thrift: %w", err)
	}

	parseMode := thrift.DefaultParseMode()
	if p.payload.IsCombineService {
		parseMode = thrift.CombineServices
	}

	desc, err := thrift.Parse(ast, parseMode)
	if err != nil {
		return nil, fmt.Errorf("cannot parse thirft IDL: %w", err)
	}
	return desc, nil
}
