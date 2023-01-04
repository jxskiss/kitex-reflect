package kitexreflect

import (
	"testing"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/stretchr/testify/assert"

	"github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl/reflectionservice"
)

func TestDiffComparator(t *testing.T) {
	kitexProv, err := generic.NewThriftFileProvider("./idl.thrift")
	assert.Nil(t, err)
	kitexDesc := <-kitexProv.Provide()

	payload := reflectionservice.GetReflectServiceRespPayload()
	pluginDesc, err := BuildServiceDescriptor(payload)
	assert.Nil(t, err)

	diff := NewDiffComparator(kitexDesc, pluginDesc)
	err = diff.Check()
	assert.Nil(t, err)

	allMsgs := diff.Messages()
	errMsgs := diff.ErrorMessages()
	infoMsgs := diff.InfoMessages()
	for _, msg := range errMsgs {
		t.Log(msg)
	}
	assert.True(t, len(allMsgs) > 0)
	assert.True(t, len(errMsgs) == 0)
	assert.True(t, len(infoMsgs) == len(allMsgs))
}
