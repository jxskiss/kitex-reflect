package kitexreflect

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"unicode"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
	svc "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl/reflectionservice"
)

func NewDescriptorProvider(ctx context.Context, serviceName string, reloadInterval time.Duration, opts ...client.Option) (generic.DescriptorProvider, error) {
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

	impl := &providerImpl{
		cli:            cli,
		reloadInterval: reloadInterval,
		close:          make(chan struct{}),
		updates:        make(chan *descriptor.ServiceDescriptor, 1),
	}
	impl.updates <- desc

	// TODO: we may don't need to do reloading at regular interval,
	//   it may be better to trigger reloading by caller.
	go impl.startReloading()

	return impl, nil
}

type providerImpl struct {
	cli            svc.Client
	reloadInterval time.Duration

	close   chan struct{}
	updates chan *descriptor.ServiceDescriptor
}

func (p *providerImpl) Close() error {
	close(p.close)
	return nil
}

func (p *providerImpl) Provide() <-chan *descriptor.ServiceDescriptor {
	return p.updates
}

func (p *providerImpl) startReloading() {

	// TODO

	select {
	case <-p.close:
		close(p.updates)
	}
}

func newReflectServiceRequest() *ReflectServiceRequest {
	return &ReflectServiceRequest{}
}

func BuildServiceDescriptor(ctx context.Context, resp *ReflectServiceResponse) (*descriptor.ServiceDescriptor, error) {
	builder := &descriptorBuilder{
		resp: resp,
		desc: &descriptor.ServiceDescriptor{},
	}
	return builder.Build()
}

type descriptorBuilder struct {
	resp *ReflectServiceResponse
	desc *descriptor.ServiceDescriptor
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
	p.desc.Name = tDesc.Name
	p.desc.Functions = make(map[string]*descriptor.FunctionDescriptor, len(tDesc.Functions))
	for _, tFuncDesc := range tDesc.Functions {
		fDesc, err := p.convertFunctionDesc(tFuncDesc)
		if err != nil {
			return nil, err
		}
		p.desc.Functions[fDesc.Name] = fDesc
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
		Name:           tDesc.Name,
		Oneway:         tDesc.Oneway,
		Request:        reqDesc,
		Response:       rspDesc,
		HasRequestBase: tDesc.HasRequestBase,
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
	structDesc, err := p.convertStructDesc(tDesc.Struct)
	if err != nil {
		return nil, err
	}
	typEnum, err := p.convertTypeEnum(tDesc.Type)
	if err != nil {
		return nil, err
	}
	desc := &descriptor.TypeDescriptor{
		Name:          tDesc.Name,
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
		Name:           tDesc.Name,
		FieldsByID:     make(map[int32]*descriptor.FieldDescriptor, len(tDesc.Fields)),
		FieldsByName:   make(map[string]*descriptor.FieldDescriptor, len(tDesc.Fields)),
		RequiredFields: make(map[int32]*descriptor.FieldDescriptor),
		DefaultFields:  make(map[string]*descriptor.FieldDescriptor),
	}
	for _, fDesc := range tDesc.Fields {
		tmp, err := p.convertFieldDesc(fDesc)
		if err != nil {
			return nil, err
		}
		desc.FieldsByID[tmp.ID] = tmp
		desc.FieldsByName[tmp.Name] = tmp
		if fDesc.Required {
			desc.RequiredFields[tmp.ID] = tmp
		}
		if fDesc.DefaultValue != nil {
			desc.DefaultFields[tmp.Name] = tmp
		}
	}
	return desc, nil
}

func (p *descriptorBuilder) convertFieldDesc(tDesc *idl.FieldDesc) (*descriptor.FieldDescriptor, error) {
	typDesc, err := p.convertTypeDesc(tDesc.Type)
	if err != nil {
		return nil, err
	}
	desc := &descriptor.FieldDescriptor{
		Name:         tDesc.Name,
		Alias:        tDesc.Alias,
		ID:           tDesc.ID,
		Required:     tDesc.Required,
		Optional:     tDesc.Optional,
		DefaultValue: nil,
		IsException:  tDesc.IsException,
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
	if desc.HTTPMapping == nil {
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
	str := *field.DefaultValue
	switch field.Type.Type {
	case idl.Type_BOOL:
		return strconv.ParseBool(str)
	case idl.Type_BYTE:
		if isDigit(str) {
			val, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return nil, err
			}
			return byte(val), nil
		}
		if len(str) == 1 {
			return byte(str[0]), nil
		}
		return 0, fmt.Errorf("invalid default value %s for type byte", str)
	case idl.Type_DOUBLE:
		return strconv.ParseFloat(str, 64)
	case idl.Type_I16:
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		return int16(val), nil
	case idl.Type_I32:
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		return int32(val), nil
	case idl.Type_I64:
		return strconv.ParseInt(str, 10, 64)
	case idl.Type_STRING:
		return str, nil
	case idl.Type_STRUCT,
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
}

func isDigit(str string) bool {
	for _, x := range str {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return str != ""
}
