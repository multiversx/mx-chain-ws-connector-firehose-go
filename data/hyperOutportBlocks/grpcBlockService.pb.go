// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: grpcBlockService.proto

package hyperOutportBlocks

import (
	duration "github.com/golang/protobuf/ptypes/duration"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BlockHashRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *BlockHashRequest) Reset() {
	*x = BlockHashRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpcBlockService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHashRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHashRequest) ProtoMessage() {}

func (x *BlockHashRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcBlockService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHashRequest.ProtoReflect.Descriptor instead.
func (*BlockHashRequest) Descriptor() ([]byte, []int) {
	return file_grpcBlockService_proto_rawDescGZIP(), []int{0}
}

func (x *BlockHashRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type BlockHashStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash            string             `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	PollingInterval *duration.Duration `protobuf:"bytes,2,opt,name=pollingInterval,proto3" json:"pollingInterval,omitempty"`
}

func (x *BlockHashStreamRequest) Reset() {
	*x = BlockHashStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpcBlockService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHashStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHashStreamRequest) ProtoMessage() {}

func (x *BlockHashStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcBlockService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHashStreamRequest.ProtoReflect.Descriptor instead.
func (*BlockHashStreamRequest) Descriptor() ([]byte, []int) {
	return file_grpcBlockService_proto_rawDescGZIP(), []int{1}
}

func (x *BlockHashStreamRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *BlockHashStreamRequest) GetPollingInterval() *duration.Duration {
	if x != nil {
		return x.PollingInterval
	}
	return nil
}

type BlockNonceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (x *BlockNonceRequest) Reset() {
	*x = BlockNonceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpcBlockService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockNonceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockNonceRequest) ProtoMessage() {}

func (x *BlockNonceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcBlockService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockNonceRequest.ProtoReflect.Descriptor instead.
func (*BlockNonceRequest) Descriptor() ([]byte, []int) {
	return file_grpcBlockService_proto_rawDescGZIP(), []int{2}
}

func (x *BlockNonceRequest) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

type BlockNonceStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce           uint64             `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	PollingInterval *duration.Duration `protobuf:"bytes,2,opt,name=pollingInterval,proto3" json:"pollingInterval,omitempty"`
}

func (x *BlockNonceStreamRequest) Reset() {
	*x = BlockNonceStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpcBlockService_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockNonceStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockNonceStreamRequest) ProtoMessage() {}

func (x *BlockNonceStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpcBlockService_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockNonceStreamRequest.ProtoReflect.Descriptor instead.
func (*BlockNonceStreamRequest) Descriptor() ([]byte, []int) {
	return file_grpcBlockService_proto_rawDescGZIP(), []int{3}
}

func (x *BlockNonceStreamRequest) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *BlockNonceStreamRequest) GetPollingInterval() *duration.Duration {
	if x != nil {
		return x.PollingInterval
	}
	return nil
}

var File_grpcBlockService_proto protoreflect.FileDescriptor

var file_grpcBlockService_proto_rawDesc = []byte{
	0x0a, 0x16, 0x67, 0x72, 0x70, 0x63, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x68, 0x79, 0x70, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x6f,
	0x72, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x2f, 0x68, 0x79, 0x70, 0x65, 0x72, 0x4f, 0x75,
	0x74, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x26, 0x0a, 0x10, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0x71, 0x0a, 0x16, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x61, 0x73, 0x68, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x43, 0x0a, 0x0f, 0x70, 0x6f, 0x6c, 0x6c, 0x69, 0x6e,
	0x67, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x70, 0x6f, 0x6c, 0x6c,
	0x69, 0x6e, 0x67, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x22, 0x29, 0x0a, 0x11, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x22, 0x74, 0x0a, 0x17, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e,
	0x6f, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0f, 0x70, 0x6f, 0x6c, 0x6c, 0x69,
	0x6e, 0x67, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x70, 0x6f, 0x6c,
	0x6c, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x32, 0xa9, 0x01, 0x0a,
	0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x4b, 0x0a, 0x0c,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x42, 0x79, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x79, 0x70, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x6f, 0x72, 0x74,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x22, 0x00, 0x30, 0x01, 0x12, 0x4d, 0x0a, 0x0d, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x73, 0x42, 0x79, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x48, 0x79, 0x70, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x6f, 0x72, 0x74, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x22, 0x00, 0x30, 0x01, 0x32, 0x9c, 0x01, 0x0a, 0x0a, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x46, 0x65, 0x74, 0x63, 0x68, 0x12, 0x45, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x42, 0x79, 0x48, 0x61, 0x73, 0x68, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x79, 0x70, 0x65, 0x72,
	0x4f, 0x75, 0x74, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x47,
	0x0a, 0x0f, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x42, 0x79, 0x4e, 0x6f, 0x6e, 0x63,
	0x65, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e,
	0x6f, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x79, 0x70, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x6f, 0x72, 0x74,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x22, 0x00, 0x42, 0x51, 0x5a, 0x4f, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x76, 0x65, 0x72, 0x73, 0x78,
	0x2f, 0x6d, 0x78, 0x2d, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2d, 0x77, 0x73, 0x2d, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2d, 0x66, 0x69, 0x72, 0x65, 0x68, 0x6f, 0x73, 0x65, 0x2d,
	0x67, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x68, 0x79, 0x70, 0x65, 0x72, 0x4f, 0x75, 0x74,
	0x70, 0x6f, 0x72, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_grpcBlockService_proto_rawDescOnce sync.Once
	file_grpcBlockService_proto_rawDescData = file_grpcBlockService_proto_rawDesc
)

func file_grpcBlockService_proto_rawDescGZIP() []byte {
	file_grpcBlockService_proto_rawDescOnce.Do(func() {
		file_grpcBlockService_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpcBlockService_proto_rawDescData)
	})
	return file_grpcBlockService_proto_rawDescData
}

var file_grpcBlockService_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_grpcBlockService_proto_goTypes = []interface{}{
	(*BlockHashRequest)(nil),        // 0: proto.BlockHashRequest
	(*BlockHashStreamRequest)(nil),  // 1: proto.BlockHashStreamRequest
	(*BlockNonceRequest)(nil),       // 2: proto.BlockNonceRequest
	(*BlockNonceStreamRequest)(nil), // 3: proto.BlockNonceStreamRequest
	(*duration.Duration)(nil),       // 4: google.protobuf.Duration
	(*HyperOutportBlock)(nil),       // 5: proto.HyperOutportBlock
}
var file_grpcBlockService_proto_depIdxs = []int32{
	4, // 0: proto.BlockHashStreamRequest.pollingInterval:type_name -> google.protobuf.Duration
	4, // 1: proto.BlockNonceStreamRequest.pollingInterval:type_name -> google.protobuf.Duration
	1, // 2: proto.BlockStream.BlocksByHash:input_type -> proto.BlockHashStreamRequest
	3, // 3: proto.BlockStream.BlocksByNonce:input_type -> proto.BlockNonceStreamRequest
	0, // 4: proto.BlockFetch.GetBlockByHash:input_type -> proto.BlockHashRequest
	2, // 5: proto.BlockFetch.GetBlockByNonce:input_type -> proto.BlockNonceRequest
	5, // 6: proto.BlockStream.BlocksByHash:output_type -> proto.HyperOutportBlock
	5, // 7: proto.BlockStream.BlocksByNonce:output_type -> proto.HyperOutportBlock
	5, // 8: proto.BlockFetch.GetBlockByHash:output_type -> proto.HyperOutportBlock
	5, // 9: proto.BlockFetch.GetBlockByNonce:output_type -> proto.HyperOutportBlock
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpcBlockService_proto_init() }
func file_grpcBlockService_proto_init() {
	if File_grpcBlockService_proto != nil {
		return
	}
	file_data_hyperOutportBlocks_hyperOutportBlock_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_grpcBlockService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHashRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpcBlockService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHashStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpcBlockService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockNonceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpcBlockService_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockNonceStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpcBlockService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_grpcBlockService_proto_goTypes,
		DependencyIndexes: file_grpcBlockService_proto_depIdxs,
		MessageInfos:      file_grpcBlockService_proto_msgTypes,
	}.Build()
	File_grpcBlockService_proto = out.File
	file_grpcBlockService_proto_rawDesc = nil
	file_grpcBlockService_proto_goTypes = nil
	file_grpcBlockService_proto_depIdxs = nil
}
