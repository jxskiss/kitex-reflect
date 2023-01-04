package kitexreflectidl

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
)

const logPrefix = "KitexReflect: "

// MarshalReflectServiceReqPayload encodes a ReflectServiceReqPayload with binary protocol.
func MarshalReflectServiceReqPayload(payload *ReflectServiceReqPayload) ([]byte, error) {
	s := thrift.NewTSerializer()
	return s.Write(context.Background(), payload)
}

// UnmarshalReflectServiceReqPayload decodes a ReflectServiceReqPayload with binary protocol.
func UnmarshalReflectServiceReqPayload(bs []byte) (*ReflectServiceReqPayload, error) {
	payload := &ReflectServiceReqPayload{}
	d := thrift.NewTDeserializer()
	err := d.Read(payload, bs)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// MarshalReflectServiceRespPayload encodes a ReflectServiceRespPayload with binary protocol.
func MarshalReflectServiceRespPayload(payload *ReflectServiceRespPayload) ([]byte, error) {
	s := thrift.NewTSerializer()
	return s.Write(context.Background(), payload)
}

// UnmarshalReflectServiceRespPayload decodes a ReflectServiceRespPayload with binary protocol.
func UnmarshalReflectServiceRespPayload(bs []byte) (*ReflectServiceRespPayload, error) {
	payload := &ReflectServiceRespPayload{}
	d := thrift.NewTDeserializer()
	err := d.Read(payload, bs)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

type RequestHandler func(ctx context.Context, req *ReflectServiceRequest, resp *ReflectServiceResponse) error

type PluginImpl struct {
	Version string

	GetIDLBytes func() []byte

	mw func(next RequestHandler) RequestHandler
}

func (p *PluginImpl) SetMiddleware(mw func(next RequestHandler) RequestHandler) {
	p.mw = mw
}

func (p *PluginImpl) ServeRequest(ctx context.Context, req *ReflectServiceRequest, resp *ReflectServiceResponse) error {
	reqPayload, err := UnmarshalReflectServiceReqPayload(req.GetPayload())
	if err != nil {
		return fmt.Errorf("cannot unmarshal ReflectServiceReqPayload: %w", err)
	}

	respPayload := &ReflectServiceRespPayload{
		Version: p.Version,
		IDL:     nil,
	}
	if reqPayload.ExistingIDLVersion != p.Version {
		respPayload.IDL = p.GetIDLBytes()
	}
	payloadBuf, err := MarshalReflectServiceRespPayload(respPayload)
	if err != nil {
		return fmt.Errorf("cannot marshal ReflectServiceRespPayload: %w", err)
	}
	resp.SetPayload(payloadBuf)
	return nil
}

func (p *PluginImpl) NewRespPayload() *ReflectServiceRespPayload {
	return &ReflectServiceRespPayload{
		Version: p.Version,
		IDL:     p.GetIDLBytes(),
	}
}

func (p *PluginImpl) NewMethodInfo() serviceinfo.MethodInfo {
	return serviceinfo.NewMethodInfo(
		p.reflectServiceHandler,
		newReflectionServiceReflectServiceArgs,
		newReflectionServiceReflectServiceResult,
		false,
	)
}

func (p *PluginImpl) reflectServiceHandler(ctx context.Context, _ interface{}, arg, result interface{}) error {
	methodHandler := p.ServeRequest
	if p.mw != nil {
		methodHandler = p.mw(methodHandler)
	}
	realArg := arg.(*ReflectionServiceReflectServiceArgs)
	realResult := result.(*ReflectionServiceReflectServiceResult)
	response := NewReflectServiceResponse()
	err := methodHandler(ctx, realArg.Request, response)
	if err != nil {
		const format = logPrefix + "failed to serve ReflectService request: %v"
		klog.CtxErrorf(ctx, format, err)
		return err
	}
	realResult.Success = response
	return nil
}

func newReflectionServiceReflectServiceArgs() interface{} {
	return NewReflectionServiceReflectServiceArgs()
}

func newReflectionServiceReflectServiceResult() interface{} {
	return NewReflectionServiceReflectServiceResult()
}
