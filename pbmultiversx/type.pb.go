// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: type.proto

package pbmultiversx

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	outport "github.com/multiversx/mx-chain-core-go/data/outport"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type BlockHeader struct {
	Height       uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Hash         string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	PreviousNum  uint64 `protobuf:"varint,3,opt,name=previous_num,json=previousNum,proto3" json:"previous_num,omitempty"`
	PreviousHash string `protobuf:"bytes,4,opt,name=previous_hash,json=previousHash,proto3" json:"previous_hash,omitempty"`
	FinalNum     uint64 `protobuf:"varint,5,opt,name=final_num,json=finalNum,proto3" json:"final_num,omitempty"`
	FinalHash    string `protobuf:"bytes,6,opt,name=final_hash,json=finalHash,proto3" json:"final_hash,omitempty"`
	Timestamp    uint64 `protobuf:"varint,7,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (m *BlockHeader) Reset()      { *m = BlockHeader{} }
func (*BlockHeader) ProtoMessage() {}
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eaed4801c3a9059, []int{0}
}
func (m *BlockHeader) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockHeader.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHeader.Merge(m, src)
}
func (m *BlockHeader) XXX_Size() int {
	return m.Size()
}
func (m *BlockHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHeader proto.InternalMessageInfo

func (m *BlockHeader) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockHeader) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *BlockHeader) GetPreviousNum() uint64 {
	if m != nil {
		return m.PreviousNum
	}
	return 0
}

func (m *BlockHeader) GetPreviousHash() string {
	if m != nil {
		return m.PreviousHash
	}
	return ""
}

func (m *BlockHeader) GetFinalNum() uint64 {
	if m != nil {
		return m.FinalNum
	}
	return 0
}

func (m *BlockHeader) GetFinalHash() string {
	if m != nil {
		return m.FinalHash
	}
	return ""
}

func (m *BlockHeader) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type Block struct {
	Header          *BlockHeader          `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	MultiversxBlock *outport.OutportBlock `protobuf:"bytes,2,opt,name=MultiversxBlock,proto3" json:"MultiversxBlock,omitempty"`
}

func (m *Block) Reset()      { *m = Block{} }
func (*Block) ProtoMessage() {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eaed4801c3a9059, []int{1}
}
func (m *Block) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Block.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(m, src)
}
func (m *Block) XXX_Size() int {
	return m.Size()
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetHeader() *BlockHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Block) GetMultiversxBlock() *outport.OutportBlock {
	if m != nil {
		return m.MultiversxBlock
	}
	return nil
}

func init() {
	proto.RegisterType((*BlockHeader)(nil), "pbmultiversx.BlockHeader")
	proto.RegisterType((*Block)(nil), "pbmultiversx.Block")
}

func init() { proto.RegisterFile("type.proto", fileDescriptor_8eaed4801c3a9059) }

var fileDescriptor_8eaed4801c3a9059 = []byte{
	// 348 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcd, 0x4e, 0xf2, 0x40,
	0x18, 0x85, 0x3b, 0xdf, 0x07, 0x28, 0x53, 0x8c, 0xc9, 0x98, 0x98, 0xfa, 0x37, 0x41, 0xdc, 0xb0,
	0xa1, 0x8d, 0xe8, 0xd6, 0x0d, 0x0b, 0xc3, 0x46, 0x4d, 0x7a, 0x03, 0x66, 0x28, 0x23, 0x9d, 0xc8,
	0x30, 0x4d, 0x3b, 0x25, 0xb0, 0xf3, 0x12, 0xbc, 0x0c, 0x2f, 0xc5, 0x25, 0x4b, 0x12, 0x37, 0x32,
	0x6c, 0x5c, 0x72, 0x09, 0xa6, 0x6f, 0xf9, 0x69, 0x5c, 0xb5, 0xef, 0x99, 0xf3, 0x9c, 0xe4, 0x1c,
	0x8c, 0xf5, 0x34, 0xe2, 0x6e, 0x14, 0x2b, 0xad, 0x48, 0x2d, 0xea, 0xc9, 0x74, 0xa8, 0xc5, 0x98,
	0xc7, 0xc9, 0xe4, 0xf4, 0x7e, 0x20, 0x74, 0x98, 0xf6, 0xdc, 0x40, 0x49, 0x6f, 0x27, 0x7b, 0x72,
	0xd2, 0x0a, 0x42, 0x26, 0x46, 0xad, 0x40, 0xc5, 0xbc, 0x35, 0x50, 0x5e, 0x9f, 0x69, 0xe6, 0xa9,
	0x54, 0x47, 0x2a, 0xd6, 0x9b, 0x6f, 0x67, 0xa8, 0x82, 0xd7, 0x3c, 0xb5, 0xf1, 0x85, 0xb0, 0x0d,
	0x77, 0x97, 0xb3, 0x3e, 0x8f, 0xc9, 0x31, 0xae, 0x84, 0x5c, 0x0c, 0x42, 0xed, 0xa0, 0x3a, 0x6a,
	0x96, 0xfc, 0xf5, 0x45, 0x08, 0x2e, 0x85, 0x2c, 0x09, 0x9d, 0x7f, 0x75, 0xd4, 0xac, 0xfa, 0xf0,
	0x4f, 0x2e, 0x71, 0x2d, 0x8a, 0xf9, 0x58, 0xa8, 0x34, 0x79, 0x1e, 0xa5, 0xd2, 0xf9, 0x0f, 0x84,
	0xbd, 0xd1, 0x1e, 0x53, 0x49, 0xae, 0xf0, 0xc1, 0xd6, 0x02, 0x7c, 0x09, 0xf8, 0x2d, 0xd7, 0xcd,
	0x72, 0xce, 0x70, 0xf5, 0x45, 0x8c, 0xd8, 0x10, 0x42, 0xca, 0x10, 0xb2, 0x0f, 0x42, 0x96, 0x70,
	0x81, 0x71, 0xfe, 0x08, 0x78, 0x05, 0xf0, 0xdc, 0x0e, 0xec, 0x39, 0xae, 0x6a, 0x21, 0x79, 0xa2,
	0x99, 0x8c, 0x9c, 0x3d, 0x60, 0x77, 0x42, 0x63, 0x8a, 0xcb, 0x50, 0x8e, 0x5c, 0x67, 0xb5, 0xb2,
	0x82, 0x50, 0xcb, 0x6e, 0x9f, 0xb8, 0xc5, 0x35, 0xdd, 0xc2, 0x02, 0xfe, 0xda, 0x48, 0xee, 0xf0,
	0xe1, 0xc3, 0xd6, 0x01, 0x06, 0x28, 0x6f, 0xb7, 0x8f, 0xf2, 0xe9, 0xdc, 0xa7, 0xc2, 0x9a, 0xfe,
	0x5f, 0x6f, 0xe7, 0x76, 0xb6, 0xa0, 0xd6, 0x7c, 0x41, 0xad, 0xd5, 0x82, 0xa2, 0x37, 0x43, 0xd1,
	0x87, 0xa1, 0xe8, 0xd3, 0x50, 0x34, 0x33, 0x14, 0x7d, 0x1b, 0x8a, 0x7e, 0x0c, 0xb5, 0x56, 0x86,
	0xa2, 0xf7, 0x25, 0xb5, 0x66, 0x4b, 0x6a, 0xcd, 0x97, 0xd4, 0xea, 0x55, 0x20, 0xfa, 0xe6, 0x37,
	0x00, 0x00, 0xff, 0xff, 0x7e, 0x76, 0x43, 0x19, 0xf9, 0x01, 0x00, 0x00,
}

func (this *BlockHeader) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*BlockHeader)
	if !ok {
		that2, ok := that.(BlockHeader)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Height != that1.Height {
		return false
	}
	if this.Hash != that1.Hash {
		return false
	}
	if this.PreviousNum != that1.PreviousNum {
		return false
	}
	if this.PreviousHash != that1.PreviousHash {
		return false
	}
	if this.FinalNum != that1.FinalNum {
		return false
	}
	if this.FinalHash != that1.FinalHash {
		return false
	}
	if this.Timestamp != that1.Timestamp {
		return false
	}
	return true
}
func (this *Block) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Block)
	if !ok {
		that2, ok := that.(Block)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Header.Equal(that1.Header) {
		return false
	}
	if !this.MultiversxBlock.Equal(that1.MultiversxBlock) {
		return false
	}
	return true
}
func (this *BlockHeader) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 11)
	s = append(s, "&pbmultiversx.BlockHeader{")
	s = append(s, "Height: "+fmt.Sprintf("%#v", this.Height)+",\n")
	s = append(s, "Hash: "+fmt.Sprintf("%#v", this.Hash)+",\n")
	s = append(s, "PreviousNum: "+fmt.Sprintf("%#v", this.PreviousNum)+",\n")
	s = append(s, "PreviousHash: "+fmt.Sprintf("%#v", this.PreviousHash)+",\n")
	s = append(s, "FinalNum: "+fmt.Sprintf("%#v", this.FinalNum)+",\n")
	s = append(s, "FinalHash: "+fmt.Sprintf("%#v", this.FinalHash)+",\n")
	s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Block) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&pbmultiversx.Block{")
	if this.Header != nil {
		s = append(s, "Header: "+fmt.Sprintf("%#v", this.Header)+",\n")
	}
	if this.MultiversxBlock != nil {
		s = append(s, "MultiversxBlock: "+fmt.Sprintf("%#v", this.MultiversxBlock)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringType(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *BlockHeader) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockHeader) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockHeader) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Timestamp != 0 {
		i = encodeVarintType(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x38
	}
	if len(m.FinalHash) > 0 {
		i -= len(m.FinalHash)
		copy(dAtA[i:], m.FinalHash)
		i = encodeVarintType(dAtA, i, uint64(len(m.FinalHash)))
		i--
		dAtA[i] = 0x32
	}
	if m.FinalNum != 0 {
		i = encodeVarintType(dAtA, i, uint64(m.FinalNum))
		i--
		dAtA[i] = 0x28
	}
	if len(m.PreviousHash) > 0 {
		i -= len(m.PreviousHash)
		copy(dAtA[i:], m.PreviousHash)
		i = encodeVarintType(dAtA, i, uint64(len(m.PreviousHash)))
		i--
		dAtA[i] = 0x22
	}
	if m.PreviousNum != 0 {
		i = encodeVarintType(dAtA, i, uint64(m.PreviousNum))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintType(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x12
	}
	if m.Height != 0 {
		i = encodeVarintType(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Block) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Block) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Block) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MultiversxBlock != nil {
		{
			size, err := m.MultiversxBlock.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintType(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Header != nil {
		{
			size, err := m.Header.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintType(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintType(dAtA []byte, offset int, v uint64) int {
	offset -= sovType(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BlockHeader) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovType(uint64(m.Height))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovType(uint64(l))
	}
	if m.PreviousNum != 0 {
		n += 1 + sovType(uint64(m.PreviousNum))
	}
	l = len(m.PreviousHash)
	if l > 0 {
		n += 1 + l + sovType(uint64(l))
	}
	if m.FinalNum != 0 {
		n += 1 + sovType(uint64(m.FinalNum))
	}
	l = len(m.FinalHash)
	if l > 0 {
		n += 1 + l + sovType(uint64(l))
	}
	if m.Timestamp != 0 {
		n += 1 + sovType(uint64(m.Timestamp))
	}
	return n
}

func (m *Block) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Header != nil {
		l = m.Header.Size()
		n += 1 + l + sovType(uint64(l))
	}
	if m.MultiversxBlock != nil {
		l = m.MultiversxBlock.Size()
		n += 1 + l + sovType(uint64(l))
	}
	return n
}

func sovType(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozType(x uint64) (n int) {
	return sovType(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *BlockHeader) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&BlockHeader{`,
		`Height:` + fmt.Sprintf("%v", this.Height) + `,`,
		`Hash:` + fmt.Sprintf("%v", this.Hash) + `,`,
		`PreviousNum:` + fmt.Sprintf("%v", this.PreviousNum) + `,`,
		`PreviousHash:` + fmt.Sprintf("%v", this.PreviousHash) + `,`,
		`FinalNum:` + fmt.Sprintf("%v", this.FinalNum) + `,`,
		`FinalHash:` + fmt.Sprintf("%v", this.FinalHash) + `,`,
		`Timestamp:` + fmt.Sprintf("%v", this.Timestamp) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Block) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Block{`,
		`Header:` + strings.Replace(this.Header.String(), "BlockHeader", "BlockHeader", 1) + `,`,
		`MultiversxBlock:` + strings.Replace(fmt.Sprintf("%v", this.MultiversxBlock), "OutportBlock", "outport.OutportBlock", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringType(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *BlockHeader) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowType
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BlockHeader: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockHeader: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthType
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthType
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousNum", wireType)
			}
			m.PreviousNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PreviousNum |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthType
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthType
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PreviousHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FinalNum", wireType)
			}
			m.FinalNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FinalNum |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FinalHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthType
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthType
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FinalHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipType(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthType
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthType
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Block) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowType
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Block: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Block: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Header", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthType
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthType
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Header == nil {
				m.Header = &BlockHeader{}
			}
			if err := m.Header.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiversxBlock", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowType
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthType
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthType
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MultiversxBlock == nil {
				m.MultiversxBlock = &outport.OutportBlock{}
			}
			if err := m.MultiversxBlock.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipType(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthType
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthType
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipType(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowType
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowType
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowType
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthType
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupType
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthType
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthType        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowType          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupType = fmt.Errorf("proto: unexpected end of group")
)
