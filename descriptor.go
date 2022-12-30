package kitexreflect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"

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
	firstResp, err := cli.ReflectService(ctx, newReflectServiceRequest())
	if err != nil {
		return nil, err
	}
	desc, err := BuildServiceDescriptor(ctx, firstResp)
	if err != nil {
		return nil, err
	}

	impl := &ProviderImpl{
		DebounceInterval: defaultDebounceInterval,
		ErrorHandler:     defaultErrorHandler,
		cli:              cli,
		serviceName:      serviceName,
		updates:          make(chan *descriptor.ServiceDescriptor, 1),
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
	closed  bool
	updates chan *descriptor.ServiceDescriptor
}

// Close closes the provider and the channel returned by Provide.
func (p *ProviderImpl) Close() error {
	p.closeMu.Lock()
	if !p.closed {
		p.closed = true
		close(p.updates)
	}
	p.closeMu.Unlock()
	return nil
}

func (p *ProviderImpl) isClosed() (ret bool) {
	p.closeMu.Lock()
	ret = p.closed
	p.closeMu.Unlock()
	return
}

// Provide returns a channel for service descriptors.
func (p *ProviderImpl) Provide() <-chan *descriptor.ServiceDescriptor {
	return p.updates
}

// Update triggers updating service descriptor from backend service.
// The actual rpc requests will be debounced according to DebounceInterval.
// This method is safe to call concurrently.
func (p *ProviderImpl) Update(ctx context.Context) {
	if p.isClosed() {
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

	resp, err := p.cli.ReflectService(ctx, newReflectServiceRequest())
	if err != nil {
		msg := fmt.Sprintf("failed to call %s.ReflectService", p.serviceName)
		p.ErrorHandler(err, msg)
		return
	}
	desc, err := BuildServiceDescriptor(ctx, resp)
	if err != nil {
		msg := fmt.Sprintf("failed to build descriptor for service %s", p.serviceName)
		p.ErrorHandler(err, msg)
		return
	}

	p.closeMu.Lock()
	defer p.closeMu.Unlock()
	if p.closed {
		return
	}
	select {
	case <-p.updates:
	default:
	}
	p.updates <- desc
}

func newReflectServiceRequest() *ReflectServiceRequest {
	return &ReflectServiceRequest{}
}

// BuildServiceDescriptor builds a [descriptor.ServiceDescriptor] from a ReflectServiceResponse.
func BuildServiceDescriptor(_ context.Context, resp *ReflectServiceResponse) (*descriptor.ServiceDescriptor, error) {
	builder := &descriptorBuilder{
		resp: resp,
		desc: &descriptor.ServiceDescriptor{},
	}
	return builder.Build()
}

type descriptorBuilder struct {
	resp *ReflectServiceResponse
	desc *descriptor.ServiceDescriptor

	idlDesc *idl.ServiceDesc
}

func (p *descriptorBuilder) Build() (*descriptor.ServiceDescriptor, error) {
	payload, err := idl.UnmarshalReflectServiceRespPayload(p.resp.Payload)
	if err != nil {
		return nil, err
	}
	tDesc := payload.GetServiceDesc()
	if tDesc == nil {
		return nil, errors.New("service descriptor not available")
	}
	p.idlDesc = tDesc
	p.desc.Name = tDesc.Name
	p.desc.Functions = make(map[string]*descriptor.FunctionDescriptor, len(tDesc.Functions))
	p.desc.Router = descriptor.NewRouter()
	for _, tFuncDesc := range tDesc.Functions {
		fDesc, err := p.convertFunctionDesc(tFuncDesc)
		if err != nil {
			return nil, err
		}
		p.desc.Functions[fDesc.Name] = fDesc
		for _, ann := range tFuncDesc.Annotations {
			for _, v := range ann.GetValues() {
				if handle, ok := descriptor.FindAnnotation(ann.GetKey(), v); ok {
					if nr, ok := handle.(descriptor.NewRoute); ok {
						p.desc.Router.Handle(nr(v, fDesc))
						break
					}
				}
			}
		}
	}
	return p.desc, nil
}

func (p *descriptorBuilder) convertFunctionDesc(tDesc *idl.FunctionDesc) (*descriptor.FunctionDescriptor, error) {
	reqDesc, err := p.convertTypeDesc(tDesc.Request)
	if err != nil {
		return nil, err
	}
	rspDesc, err := p.convertTypeDesc(tDesc.Response)
	if err != nil {
		return nil, err
	}
	desc := &descriptor.FunctionDescriptor{
		Name:           tDesc.GetName(),
		Oneway:         tDesc.GetOneway(),
		Request:        reqDesc,
		Response:       rspDesc,
		HasRequestBase: tDesc.GetHasRequestBase(),
	}
	return desc, nil
}

func (p *descriptorBuilder) convertTypeDesc(tDesc *idl.TypeDesc) (*descriptor.TypeDescriptor, error) {
	if tDesc == nil {
		return nil, nil
	}
	keyDesc, err := p.convertTypeDesc(tDesc.Key)
	if err != nil {
		return nil, err
	}
	elemDesc, err := p.convertTypeDesc(tDesc.Elem)
	if err != nil {
		return nil, err
	}
	idlStructDesc := tDesc.Struct
	if idlStructDesc == nil && tDesc.StructIdx != nil && len(p.idlDesc.StructList) > int(*tDesc.StructIdx) {
		idlStructDesc = p.idlDesc.StructList[*tDesc.StructIdx]
	}
	structDesc, err := p.convertStructDesc(idlStructDesc)
	if err != nil {
		return nil, err
	}
	typEnum, err := p.convertTypeEnum(tDesc.GetType())
	if err != nil {
		return nil, err
	}
	desc := &descriptor.TypeDescriptor{
		Name:          tDesc.GetName(),
		Type:          typEnum,
		Key:           keyDesc,
		Elem:          elemDesc,
		Struct:        structDesc,
		IsRequestBase: tDesc.GetIsRequestBase(),
	}
	return desc, nil
}

func (p *descriptorBuilder) convertStructDesc(tDesc *idl.StructDesc) (*descriptor.StructDescriptor, error) {
	if tDesc == nil {
		return nil, nil
	}
	desc := &descriptor.StructDescriptor{
		Name:           tDesc.GetName(),
		FieldsByID:     make(map[int32]*descriptor.FieldDescriptor, len(tDesc.Fields)),
		FieldsByName:   make(map[string]*descriptor.FieldDescriptor, len(tDesc.Fields)),
		RequiredFields: make(map[int32]*descriptor.FieldDescriptor),
		DefaultFields:  make(map[string]*descriptor.FieldDescriptor),
	}
	for _, fDesc := range tDesc.Fields {
		_f, err := p.convertFieldDesc(tDesc.GetName(), fDesc)
		if err != nil {
			return nil, err
		}
		desc.FieldsByID[_f.ID] = _f
		desc.FieldsByName[_f.FieldName()] = _f
		if _f.Required {
			desc.RequiredFields[_f.ID] = _f
		}
		if _f.DefaultValue != nil {
			desc.DefaultFields[_f.Name] = _f
		}
	}
	return desc, nil
}

func (p *descriptorBuilder) convertFieldDesc(structName string, tDesc *idl.FieldDesc) (*descriptor.FieldDescriptor, error) {
	typDesc, err := p.convertTypeDesc(tDesc.Type)
	if err != nil {
		return nil, err
	}
	desc := &descriptor.FieldDescriptor{
		Name:         tDesc.GetName(),
		Alias:        tDesc.GetAlias(),
		ID:           tDesc.GetID(),
		Required:     tDesc.GetRequired(),
		Optional:     tDesc.GetOptional(),
		DefaultValue: nil,
		IsException:  tDesc.GetIsException(),
		Type:         typDesc,
		HTTPMapping:  nil,
		ValueMapping: nil,
	}

	// Default value.
	if tDesc.DefaultValue != nil {
		desc.DefaultValue, err = p.parseDefaultValue(tDesc)
		if err != nil {
			return nil, err
		}
	}

	// Mapping.
	for _, ann := range tDesc.Annotations {
		for _, v := range ann.Values {
			if handle, ok := descriptor.FindAnnotation(ann.GetKey(), v); ok {
				switch h := handle.(type) {
				case descriptor.NewHTTPMapping:
					desc.HTTPMapping = h(v)
				case descriptor.NewValueMapping:
					desc.ValueMapping = h(v)
				case descriptor.NewFieldMapping:
					// execute at compile time
					h(v).Handle(desc)
				case nil:
					// none annotation
				default:
					// not supported annotation type
					return nil, fmt.Errorf("not supported handle type: %T", handle)
				}
			}
		}
	}
	if desc.HTTPMapping == nil && structName != "" && desc.FieldName() != "" {
		desc.HTTPMapping = descriptor.DefaultNewMapping(desc.FieldName())
	}

	return desc, nil
}

var typeEnumTable = [...]descriptor.Type{
	idl.Type_BOOL:   descriptor.BOOL,
	idl.Type_BYTE:   descriptor.BYTE,
	idl.Type_DOUBLE: descriptor.DOUBLE,
	idl.Type_I16:    descriptor.I16,
	idl.Type_I32:    descriptor.I32,
	idl.Type_I64:    descriptor.I64,
	idl.Type_STRING: descriptor.STRING,
	idl.Type_STRUCT: descriptor.STRUCT,
	idl.Type_MAP:    descriptor.MAP,
	idl.Type_SET:    descriptor.SET,
	idl.Type_LIST:   descriptor.LIST,
	idl.Type_UTF8:   descriptor.UTF8,
	idl.Type_UTF16:  descriptor.UTF16,
	idl.Type_JSON:   descriptor.JSON,
}

func (p *descriptorBuilder) convertTypeEnum(tTyp idl.Type) (descriptor.Type, error) {
	var ret descriptor.Type
	if int(tTyp) < len(typeEnumTable) {
		ret = typeEnumTable[tTyp]
	}
	if ret > 0 {
		return ret, nil
	}
	return 0, fmt.Errorf("not supported type enum: %d", tTyp)
}

func (p *descriptorBuilder) parseDefaultValue(field *idl.FieldDesc) (interface{}, error) {
	var out interface{}
	str := *field.DefaultValue
	switch *field.Type.Type {
	case idl.Type_BOOL:
		out = new(bool)
	case idl.Type_DOUBLE:
		out = new(float64)
	case idl.Type_I16:
		out = new(int16)
	case idl.Type_I32:
		out = new(int32)
	case idl.Type_I64:
		out = new(int64)
	case idl.Type_STRING:
		out = new(string)
	case idl.Type_BYTE,
		idl.Type_STRUCT,
		idl.Type_MAP,
		idl.Type_SET,
		idl.Type_LIST,
		idl.Type_UTF8,
		idl.Type_UTF16,
		idl.Type_JSON:
		return nil, fmt.Errorf("not supported default value for type %v", field.Type.Type)
	default:
		return nil, fmt.Errorf("not supported field type: %v", field.Type.Type)
	}
	err := json.Unmarshal([]byte(str), out)
	if err != nil {
		return nil, fmt.Errorf("cannot parse default value: %w", err)
	}
	out = reflect.Indirect(reflect.ValueOf(out)).Interface()
	return out, nil
}
