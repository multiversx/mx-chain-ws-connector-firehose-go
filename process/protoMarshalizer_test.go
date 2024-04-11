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

func TestGogoProtoMarshalizer_Marshal(t *testing.T) {
	outportBlock := data.ShardOutportBlock{}
	encNode, err := recovedMarshal(&outportBlock)
	assert.Nil(t, err)
	assert.NotNil(t, encNode)
}

func TestGogoProtoMarshalizer_MarshalWrongObj(t *testing.T) {
	obj := "multiversx"
	encNode, err := recovedMarshal(obj)
	assert.Nil(t, encNode)
	assert.NotNil(t, err)
}

func TestGogoProtoMarshalizer_Unmarshal(t *testing.T) {
	protoMarshaller := ProtoMarshalizer{}
	outportBlock := data.ShardOutportBlock{}

	encNode, _ := protoMarshaller.Marshal(outportBlock)
	newNode := &data.ShardOutportBlock{}

	err := recovedUnmarshal(newNode, encNode)
	assert.Nil(t, err)
	assert.Equal(t, outportBlock.ProtoReflect(), newNode.ProtoReflect())
}

func TestGogoProtoMarshalizer_UnmarshalWrongObj(t *testing.T) {
	protoMarshaller := ProtoMarshalizer{}
	outportBlock := data.ShardOutportBlock{}

	encNode, _ := protoMarshaller.Marshal(outportBlock)
	err := recovedUnmarshal([]byte{}, encNode)
	assert.NotNil(t, err)
}
