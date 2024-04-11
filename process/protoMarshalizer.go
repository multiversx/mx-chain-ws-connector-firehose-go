package process

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// ProtoMarshalizer marshals/unmarshals proto.Message(s)
type ProtoMarshalizer struct{}

// Marshal any object into a slice of bytes.
func (p *ProtoMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	message, ok := obj.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot marshal %T into a proto.Message", obj)
	}
	return proto.Marshal(message)
}

// Unmarshal a slice of bytes into a proto message.
func (p *ProtoMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	message, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("cannot marshal %T into a proto.Message", obj)
	}
	return proto.Unmarshal(buff, message)
}

func (p *ProtoMarshalizer) IsInterfaceNil() bool {
	return p == nil
}
