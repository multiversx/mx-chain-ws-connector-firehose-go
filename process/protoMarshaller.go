package process

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// ProtoMarshaller marshals/unmarshals proto.Message(s)
type ProtoMarshaller struct{}

// Marshal any object into a slice of bytes.
func (p *ProtoMarshaller) Marshal(obj interface{}) ([]byte, error) {
	message, ok := obj.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot marshal %T into a proto.Message", obj)
	}
	return proto.Marshal(message)
}

// Unmarshal a slice of bytes into a proto message.
func (p *ProtoMarshaller) Unmarshal(obj interface{}, buff []byte) error {
	message, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("cannot marshal %T into a proto.Message", obj)
	}
	return proto.Unmarshal(buff, message)
}

// IsInterfaceNil returns nil if there is no value under the interface.
func (p *ProtoMarshaller) IsInterfaceNil() bool {
	return p == nil
}
