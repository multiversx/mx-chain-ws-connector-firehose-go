package process

import (
	"google.golang.org/protobuf/proto"
)

type ProtoMarshaller struct{}

func (p ProtoMarshaller) Marshal(obj interface{}) ([]byte, error) {

	return proto.Marshal(obj)
}

func (p ProtoMarshaller) Unmarshal(obj interface{}, buff []byte) error {
	//TODO implement me
	panic("implement me")
}

func (p ProtoMarshaller) IsInterfaceNil() bool {
	//TODO implement me
	panic("implement me")
}
