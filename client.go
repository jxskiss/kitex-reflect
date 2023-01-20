package kitexreflect

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

// NewGenericClient creates a genericclient.Client for a service.
// gFactory should be a factory function from the [generic] package.
func NewGenericClient(
	ctx context.Context,
	serviceName string,
	gFactory func(p generic.DescriptorProvider) (generic.Generic, error),
	kcOptions []client.Option,
	providerOpts ...Option,
) (genericclient.Client, error) {
	prov, err := NewDescriptorProvider(ctx, serviceName, kcOptions, providerOpts...)
	if err != nil {
		return nil, err
	}
	g, err := gFactory(prov)
	if err != nil {
		return nil, err
	}
	cli, err := genericclient.NewClient(serviceName, g, kcOptions...)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
