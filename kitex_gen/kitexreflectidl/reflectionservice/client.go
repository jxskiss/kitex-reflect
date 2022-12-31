// Code generated by Kitex v0.4.4. DO NOT EDIT.

package reflectionservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	kitexreflectidl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	ReflectService(ctx context.Context, request *kitexreflectidl.ReflectServiceRequest, callOptions ...callopt.Option) (r *kitexreflectidl.ReflectServiceResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kReflectionServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kReflectionServiceClient struct {
	*kClient
}

func (p *kReflectionServiceClient) ReflectService(ctx context.Context, request *kitexreflectidl.ReflectServiceRequest, callOptions ...callopt.Option) (r *kitexreflectidl.ReflectServiceResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ReflectService(ctx, request)
}
