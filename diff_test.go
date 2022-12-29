package kitexreflect

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/stretchr/testify/assert"

	"github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl/reflectionservice"
)

func TestDiffComparator(t *testing.T) {
	kitexProv, err := generic.NewThriftFileProvider("./idl.thrift")
	assert.Nil(t, err)
	kitexDesc := <-kitexProv.Provide()

	ctx := context.Background()
	payload, err := reflectionservice.GetReflectServiceRespPayload(ctx)
	assert.Nil(t, err)

	reflectResp := &ReflectServiceResponse{}
	reflectResp.SetPayload(payload)
	pluginDesc, err := BuildServiceDescriptor(ctx, reflectResp)
	assert.Nil(t, err)

	diff := NewDiffComparator(kitexDesc, pluginDesc)
	err = diff.Check()
	assert.Nil(t, err)

	allMsgs := diff.Messages()
	errMsgs := diff.ErrorMessages()
	infoMsgs := diff.InfoMessages()
	assert.True(t, len(allMsgs) > 0)
	assert.True(t, len(errMsgs) == 0)
	assert.True(t, len(infoMsgs) == len(allMsgs))
}
