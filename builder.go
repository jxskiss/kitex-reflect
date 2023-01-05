package kitexreflect

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/generic/thrift"
	"github.com/cloudwego/thriftgo/generator/golang/extension/meta"
	"github.com/cloudwego/thriftgo/parser"

	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
	svc "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl/reflectionservice"
)

func doRequestDescriptor(ctx context.Context, cli svc.Client, existingIDLVersion string) (
	desc *descriptor.ServiceDescriptor, version string, err error) {
	request := newReflectServiceRequest(existingIDLVersion)
	response, err := cli.ReflectService(ctx, request)
	if err != nil {
		return nil, "", err
	}
	respPayload, err := idl.UnmarshalReflectServiceRespPayload(response.Payload)
	if err != nil {
		return nil, "", err
	}
	desc, err = BuildServiceDescriptor(respPayload)
	if err != nil {
		return nil, "", err
	}
	return desc, respPayload.Version, nil
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
