package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"google.golang.org/protobuf/proto"
)

var protoMarshaller = &marshal.GogoProtoMarshalizer{}

type outportBlockCaster struct {
	outportBlock *outport.OutportBlock
}

// NewOutportBlockCaster returns a new instance of outportBlockCaster.
func NewOutportBlockCaster(outportBlock *outport.OutportBlock) *outportBlockCaster {
	return &outportBlockCaster{outportBlock: outportBlock}
}

// MarshalTo will marshal the underlying outport.OutportBlock to the given proto.Message
func (o *outportBlockCaster) MarshalTo(m proto.Message) error {
	obBytes, err := protoMarshaller.Marshal(o.outportBlock)
	if err != nil {
		return fmt.Errorf("failed to marshal gogo outport block: %w", err)
	}

	err = proto.Unmarshal(obBytes, m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to proto message [%s]: %w",
			m.ProtoReflect().Descriptor().Name(), err)
	}

	return nil
}
