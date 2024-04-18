package process

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

func recoveredMarshal(obj interface{}) (buf []byte, err error) {
	protoMarshaller := ProtoMarshaller{}

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

func recoveredUnmarshal(obj interface{}, buf []byte) (err error) {
	protoMarshaller := ProtoMarshaller{}

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

func TestProtoMarshaller_Marshal(t *testing.T) {
	t.Parallel()

	outportBlock := data.ShardOutportBlock{}
	encNode, err := recoveredMarshal(&outportBlock)
	assert.Nil(t, err)
	assert.NotNil(t, encNode)
}

func TestProtoMarshaller_MarshalWrongObj(t *testing.T) {
	t.Parallel()

	obj := "multiversx"
	encNode, err := recoveredMarshal(obj)
	assert.Nil(t, encNode)
	assert.Equal(t, errors.New("cannot marshal string into a proto.Message"), err)
}

func TestProtoMarshaller_Unmarshal(t *testing.T) {
	t.Parallel()

	protoMarshaller := ProtoMarshaller{}
	outportBlock := data.ShardOutportBlock{}

	encNode, _ := protoMarshaller.Marshal(&outportBlock)
	newNode := &data.ShardOutportBlock{}

	err := recoveredUnmarshal(newNode, encNode)
	assert.Nil(t, err)
	assert.Equal(t, outportBlock.ProtoReflect(), newNode.ProtoReflect())
}

func TestProtoMarshaller_UnmarshalWrongObj(t *testing.T) {
	t.Parallel()

	protoMarshaller := ProtoMarshaller{}
	outportBlock := data.ShardOutportBlock{}

	encNode, _ := protoMarshaller.Marshal(&outportBlock)
	err := recoveredUnmarshal([]byte{}, encNode)
	assert.NotNil(t, err)
}
