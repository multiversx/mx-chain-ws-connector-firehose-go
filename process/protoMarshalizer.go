package process

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type ProtoMarshaller struct{}

func (p *ProtoMarshaller) Marshal(obj interface{}) ([]byte, error) {
	message, ok := obj.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot marshal %T into a proto.Message", obj)
	}
	return proto.Marshal(message)
}

func (p *ProtoMarshaller) Unmarshal(obj interface{}, buff []byte) error {
	message, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("cannot marshal %T into a proto.Message", obj)
	}
	return proto.Unmarshal(buff, message)
}

func (p *ProtoMarshaller) IsInterfaceNil() bool {
	return p == nil
}
