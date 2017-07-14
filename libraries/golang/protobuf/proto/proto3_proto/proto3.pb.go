// Code generated by protoc-gen-go.
// source: proto3_proto/proto3.proto
// DO NOT EDIT!

/*
Package proto3_proto is a generated protocol buffer package.

It is generated from these files:
	proto3_proto/proto3.proto

It has these top-level messages:
	Message
	Nested
	MessageWithMap
*/
package proto3_proto

import proto "github.com/insionng/yougam/libraries/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/insionng/yougam/libraries/golang/protobuf/ptypes/any"
import testdata "github.com/insionng/yougam/libraries/golang/protobuf/proto/testdata"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type Message_Humour int32

const (
	Message_UNKNOWN     Message_Humour = 0
	Message_PUNS        Message_Humour = 1
	Message_SLAPSTICK   Message_Humour = 2
	Message_BILL_BAILEY Message_Humour = 3
)

var Message_Humour_name = map[int32]string{
	0: "UNKNOWN",
	1: "PUNS",
	2: "SLAPSTICK",
	3: "BILL_BAILEY",
}
var Message_Humour_value = map[string]int32{
	"UNKNOWN":     0,
	"PUNS":        1,
	"SLAPSTICK":   2,
	"BILL_BAILEY": 3,
}

func (x Message_Humour) String() string {
	return proto.EnumName(Message_Humour_name, int32(x))
}
func (Message_Humour) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Message struct {
	Name         string                           `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Hilarity     Message_Humour                   `protobuf:"varint,2,opt,name=hilarity,enum=proto3_proto.Message_Humour" json:"hilarity,omitempty"`
	HeightInCm   uint32                           `protobuf:"varint,3,opt,name=height_in_cm,json=heightInCm" json:"height_in_cm,omitempty"`
	Data         []byte                           `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	ResultCount  int64                            `protobuf:"varint,7,opt,name=result_count,json=resultCount" json:"result_count,omitempty"`
	TrueScotsman bool                             `protobuf:"varint,8,opt,name=true_scotsman,json=trueScotsman" json:"true_scotsman,omitempty"`
	Score        float32                          `protobuf:"fixed32,9,opt,name=score" json:"score,omitempty"`
	Key          []uint64                         `protobuf:"varint,5,rep,name=key" json:"key,omitempty"`
	Nested       *Nested                          `protobuf:"bytes,6,opt,name=nested" json:"nested,omitempty"`
	RFunny       []Message_Humour                 `protobuf:"varint,16,rep,name=r_funny,json=rFunny,enum=proto3_proto.Message_Humour" json:"r_funny,omitempty"`
	Terrain      map[string]*Nested               `protobuf:"bytes,10,rep,name=terrain" json:"terrain,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Proto2Field  *testdata.SubDefaults            `protobuf:"bytes,11,opt,name=proto2_field,json=proto2Field" json:"proto2_field,omitempty"`
	Proto2Value  map[string]*testdata.SubDefaults `protobuf:"bytes,13,rep,name=proto2_value,json=proto2Value" json:"proto2_value,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Anything     *google_protobuf.Any             `protobuf:"bytes,14,opt,name=anything" json:"anything,omitempty"`
	ManyThings   []*google_protobuf.Any           `protobuf:"bytes,15,rep,name=many_things,json=manyThings" json:"many_things,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetNested() *Nested {
	if m != nil {
		return m.Nested
	}
	return nil
}

func (m *Message) GetTerrain() map[string]*Nested {
	if m != nil {
		return m.Terrain
	}
	return nil
}

func (m *Message) GetProto2Field() *testdata.SubDefaults {
	if m != nil {
		return m.Proto2Field
	}
	return nil
}

func (m *Message) GetProto2Value() map[string]*testdata.SubDefaults {
	if m != nil {
		return m.Proto2Value
	}
	return nil
}

func (m *Message) GetAnything() *google_protobuf.Any {
	if m != nil {
		return m.Anything
	}
	return nil
}

func (m *Message) GetManyThings() []*google_protobuf.Any {
	if m != nil {
		return m.ManyThings
	}
	return nil
}

type Nested struct {
	Bunny string `protobuf:"bytes,1,opt,name=bunny" json:"bunny,omitempty"`
}

func (m *Nested) Reset()                    { *m = Nested{} }
func (m *Nested) String() string            { return proto.CompactTextString(m) }
func (*Nested) ProtoMessage()               {}
func (*Nested) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type MessageWithMap struct {
	ByteMapping map[bool][]byte `protobuf:"bytes,1,rep,name=byte_mapping,json=byteMapping" json:"byte_mapping,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *MessageWithMap) Reset()                    { *m = MessageWithMap{} }
func (m *MessageWithMap) String() string            { return proto.CompactTextString(m) }
func (*MessageWithMap) ProtoMessage()               {}
func (*MessageWithMap) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MessageWithMap) GetByteMapping() map[bool][]byte {
	if m != nil {
		return m.ByteMapping
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "proto3_proto.Message")
	proto.RegisterType((*Nested)(nil), "proto3_proto.Nested")
	proto.RegisterType((*MessageWithMap)(nil), "proto3_proto.MessageWithMap")
	proto.RegisterEnum("proto3_proto.Message_Humour", Message_Humour_name, Message_Humour_value)
}

var fileDescriptor0 = []byte{
	// 617 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x92, 0x5d, 0x6b, 0xdb, 0x3c,
	0x14, 0xc7, 0x1f, 0xc5, 0xa9, 0x93, 0x1e, 0x3b, 0xad, 0xd1, 0xd3, 0x81, 0x1a, 0xc6, 0xf0, 0x32,
	0x18, 0x66, 0x2f, 0xee, 0xc8, 0x28, 0x94, 0x31, 0x36, 0xda, 0xae, 0x65, 0xa1, 0x69, 0x16, 0x9c,
	0x76, 0x65, 0x57, 0x46, 0x49, 0x95, 0xc4, 0x2c, 0x96, 0x83, 0x2d, 0x0f, 0xfc, 0x75, 0xf6, 0x29,
	0x77, 0x39, 0x24, 0x39, 0xa9, 0x5b, 0xb2, 0xed, 0xca, 0xd2, 0xf1, 0xef, 0xbc, 0xe8, 0xff, 0x3f,
	0xb0, 0xbf, 0x4c, 0x13, 0x91, 0xbc, 0x0d, 0xd5, 0xe7, 0x40, 0x5f, 0x7c, 0xf5, 0xc1, 0x76, 0xf5,
	0x57, 0x7b, 0x7f, 0x96, 0x24, 0xb3, 0x05, 0xd3, 0xc8, 0x38, 0x9f, 0x1e, 0x50, 0x5e, 0x68, 0xb0,
	0xfd, 0xbf, 0x60, 0x99, 0xb8, 0xa5, 0x82, 0x1e, 0xc8, 0x83, 0x0e, 0x76, 0x7e, 0x99, 0xd0, 0xb8,
	0x64, 0x59, 0x46, 0x67, 0x0c, 0x63, 0xa8, 0x73, 0x1a, 0x33, 0x82, 0x5c, 0xe4, 0x6d, 0x07, 0xea,
	0x8c, 0x8f, 0xa0, 0x39, 0x8f, 0x16, 0x34, 0x8d, 0x44, 0x41, 0x6a, 0x2e, 0xf2, 0x76, 0xba, 0x8f,
	0xfd, 0x6a, 0x43, 0xbf, 0x4c, 0xf6, 0x3f, 0xe7, 0x71, 0x92, 0xa7, 0xc1, 0x9a, 0xc6, 0x2e, 0xd8,
	0x73, 0x16, 0xcd, 0xe6, 0x22, 0x8c, 0x78, 0x38, 0x89, 0x89, 0xe1, 0x22, 0xaf, 0x15, 0x80, 0x8e,
	0xf5, 0xf8, 0x69, 0x2c, 0xfb, 0xc9, 0x71, 0x48, 0xdd, 0x45, 0x9e, 0x1d, 0xa8, 0x33, 0x7e, 0x0a,
	0x76, 0xca, 0xb2, 0x7c, 0x21, 0xc2, 0x49, 0x92, 0x73, 0x41, 0x1a, 0x2e, 0xf2, 0x8c, 0xc0, 0xd2,
	0xb1, 0x53, 0x19, 0xc2, 0xcf, 0xa0, 0x25, 0xd2, 0x9c, 0x85, 0xd9, 0x24, 0x11, 0x59, 0x4c, 0x39,
	0x69, 0xba, 0xc8, 0x6b, 0x06, 0xb6, 0x0c, 0x8e, 0xca, 0x18, 0xde, 0x83, 0xad, 0x6c, 0x92, 0xa4,
	0x8c, 0x6c, 0xbb, 0xc8, 0xab, 0x05, 0xfa, 0x82, 0x1d, 0x30, 0xbe, 0xb3, 0x82, 0x6c, 0xb9, 0x86,
	0x57, 0x0f, 0xe4, 0x11, 0xbf, 0x02, 0x93, 0xb3, 0x4c, 0xb0, 0x5b, 0x62, 0xba, 0xc8, 0xb3, 0xba,
	0x7b, 0xf7, 0x5f, 0x37, 0x50, 0xff, 0x82, 0x92, 0xc1, 0x87, 0xd0, 0x48, 0xc3, 0x69, 0xce, 0x79,
	0x41, 0x1c, 0xd7, 0xf8, 0xa7, 0x18, 0x66, 0x7a, 0x2e, 0x59, 0xfc, 0x1e, 0x1a, 0x82, 0xa5, 0x29,
	0x8d, 0x38, 0x01, 0xd7, 0xf0, 0xac, 0x6e, 0x67, 0x73, 0xda, 0x95, 0x86, 0xce, 0xb8, 0x48, 0x8b,
	0x60, 0x95, 0x82, 0x8f, 0x40, 0x5b, 0xdc, 0x0d, 0xa7, 0x11, 0x5b, 0xdc, 0x12, 0x4b, 0x0d, 0xfa,
	0xc8, 0x5f, 0xd9, 0xe9, 0x8f, 0xf2, 0xf1, 0x27, 0x36, 0xa5, 0xf9, 0x42, 0x64, 0x81, 0xa5, 0xd1,
	0x73, 0x49, 0xe2, 0xde, 0x3a, 0xf3, 0x07, 0x5d, 0xe4, 0x8c, 0xb4, 0x54, 0xf3, 0xe7, 0x9b, 0x9b,
	0x0f, 0x15, 0xf9, 0x55, 0x82, 0x7a, 0x80, 0xb2, 0x94, 0x8a, 0xe0, 0x37, 0xd0, 0xa4, 0xbc, 0x10,
	0xf3, 0x88, 0xcf, 0xc8, 0x4e, 0xa9, 0x94, 0x5e, 0x35, 0x7f, 0xb5, 0x6a, 0xfe, 0x31, 0x2f, 0x82,
	0x35, 0x85, 0x0f, 0xc1, 0x8a, 0x29, 0x2f, 0x42, 0x75, 0xcb, 0xc8, 0xae, 0xea, 0xbd, 0x39, 0x09,
	0x24, 0x78, 0xa5, 0xb8, 0xf6, 0x10, 0xec, 0xaa, 0x0c, 0x2b, 0xcb, 0xf4, 0x4e, 0x2a, 0xcb, 0x5e,
	0xc0, 0x96, 0x7e, 0x4e, 0xed, 0x2f, 0x8e, 0x69, 0xe4, 0x5d, 0xed, 0x08, 0xb5, 0xaf, 0xc1, 0x79,
	0xf8, 0xb6, 0x0d, 0x55, 0x5f, 0xde, 0xaf, 0xfa, 0x07, 0x79, 0xef, 0xca, 0x76, 0x3e, 0x82, 0xa9,
	0x6d, 0xc6, 0x16, 0x34, 0xae, 0x07, 0x17, 0x83, 0x2f, 0x37, 0x03, 0xe7, 0x3f, 0xdc, 0x84, 0xfa,
	0xf0, 0x7a, 0x30, 0x72, 0x10, 0x6e, 0xc1, 0xf6, 0xa8, 0x7f, 0x3c, 0x1c, 0x5d, 0xf5, 0x4e, 0x2f,
	0x9c, 0x1a, 0xde, 0x05, 0xeb, 0xa4, 0xd7, 0xef, 0x87, 0x27, 0xc7, 0xbd, 0xfe, 0xd9, 0x37, 0xc7,
	0xe8, 0x3c, 0x01, 0x53, 0x0f, 0x2b, 0x97, 0x75, 0xac, 0x96, 0x4a, 0xcf, 0xa3, 0x2f, 0x9d, 0x9f,
	0x08, 0x76, 0x4a, 0x73, 0x6e, 0x22, 0x31, 0xbf, 0xa4, 0x4b, 0x3c, 0x04, 0x7b, 0x5c, 0x08, 0x16,
	0xc6, 0x74, 0xb9, 0x94, 0x4e, 0x20, 0x25, 0xea, 0xeb, 0x8d, 0x86, 0x96, 0x39, 0xfe, 0x49, 0x21,
	0xd8, 0xa5, 0xe6, 0x4b, 0x5f, 0xc7, 0x77, 0x91, 0xf6, 0x07, 0x70, 0x1e, 0x02, 0x55, 0x71, 0x9a,
	0x5a, 0x9c, 0xbd, 0xaa, 0x38, 0x76, 0x45, 0x85, 0xb1, 0xa9, 0x5b, 0xff, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0x54, 0x4a, 0xfa, 0x41, 0xa1, 0x04, 0x00, 0x00,
}
