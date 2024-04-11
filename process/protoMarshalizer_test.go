package process

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

func recovedMarshal(obj interface{}) (buf []byte, err error) {
	protoMarshaller := ProtoMarshalizer{}

	defer func() {
		if p := recover(); p != nil {
			if panicError, ok := p.(error); ok {
				err = panicError
			} else {
				err = fmt.Errorf("%#v", p)
			}
			buf = nil
		}
	}()
	buf, err = protoMarshaller.Marshal(obj)
	return
}

func recovedUnmarshal(obj interface{}, buf []byte) (err error) {
	protoMarshaller := ProtoMarshalizer{}

	defer func() {
		if p := recover(); p != nil {
			if panicError, ok := p.(error); ok {
				err = panicError
			} else {
				err = fmt.Errorf("%#v", p)
			}
		}
	}()
	err = protoMarshaller.Unmarshal(obj, buf)
	return
}

func TestProtoMarshalizer_Marshal(t *testing.T) {
	t.Parallel()

	outportBlock := data.ShardOutportBlock{}
	encNode, err := recovedMarshal(&outportBlock)
	assert.Nil(t, err)
	assert.NotNil(t, encNode)
}

func TestProtoMarshalizer_MarshalWrongObj(t *testing.T) {
	t.Parallel()

	obj := "multiversx"
	encNode, err := recovedMarshal(obj)
	assert.Nil(t, encNode)
	assert.NotNil(t, err)
}

func TestProtoMarshalizer_Unmarshal(t *testing.T) {
	t.Parallel()

	protoMarshaller := ProtoMarshalizer{}
	outportBlock := data.ShardOutportBlock{}

	encNode, _ := protoMarshaller.Marshal(outportBlock)
	newNode := &data.ShardOutportBlock{}

	err := recovedUnmarshal(newNode, encNode)
	assert.Nil(t, err)
	assert.Equal(t, outportBlock.ProtoReflect(), newNode.ProtoReflect())
}

func TestProtoMarshalizer_UnmarshalWrongObj(t *testing.T) {
	t.Parallel()

	protoMarshaller := ProtoMarshalizer{}
	outportBlock := data.ShardOutportBlock{}

	encNode, _ := protoMarshaller.Marshal(outportBlock)
	err := recovedUnmarshal([]byte{}, encNode)
	assert.NotNil(t, err)
}
