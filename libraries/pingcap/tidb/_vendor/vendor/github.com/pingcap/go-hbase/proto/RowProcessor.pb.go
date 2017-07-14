// Code generated by protoc-gen-go.
// source: RowProcessor.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/insionng/yougam/libraries/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

type ProcessRequest struct {
	RowProcessorClassName              *string `protobuf:"bytes,1,req,name=row_processor_class_name" json:"row_processor_class_name,omitempty"`
	RowProcessorInitializerMessageName *string `protobuf:"bytes,2,opt,name=row_processor_initializer_message_name" json:"row_processor_initializer_message_name,omitempty"`
	RowProcessorInitializerMessage     []byte  `protobuf:"bytes,3,opt,name=row_processor_initializer_message" json:"row_processor_initializer_message,omitempty"`
	NonceGroup                         *uint64 `protobuf:"varint,4,opt,name=nonce_group" json:"nonce_group,omitempty"`
	Nonce                              *uint64 `protobuf:"varint,5,opt,name=nonce" json:"nonce,omitempty"`
	XXX_unrecognized                   []byte  `json:"-"`
}

func (m *ProcessRequest) Reset()         { *m = ProcessRequest{} }
func (m *ProcessRequest) String() string { return proto1.CompactTextString(m) }
func (*ProcessRequest) ProtoMessage()    {}

func (m *ProcessRequest) GetRowProcessorClassName() string {
	if m != nil && m.RowProcessorClassName != nil {
		return *m.RowProcessorClassName
	}
	return ""
}

func (m *ProcessRequest) GetRowProcessorInitializerMessageName() string {
	if m != nil && m.RowProcessorInitializerMessageName != nil {
		return *m.RowProcessorInitializerMessageName
	}
	return ""
}

func (m *ProcessRequest) GetRowProcessorInitializerMessage() []byte {
	if m != nil {
		return m.RowProcessorInitializerMessage
	}
	return nil
}

func (m *ProcessRequest) GetNonceGroup() uint64 {
	if m != nil && m.NonceGroup != nil {
		return *m.NonceGroup
	}
	return 0
}

func (m *ProcessRequest) GetNonce() uint64 {
	if m != nil && m.Nonce != nil {
		return *m.Nonce
	}
	return 0
}

type ProcessResponse struct {
	RowProcessorResult []byte `protobuf:"bytes,1,req,name=row_processor_result" json:"row_processor_result,omitempty"`
	XXX_unrecognized   []byte `json:"-"`
}

func (m *ProcessResponse) Reset()         { *m = ProcessResponse{} }
func (m *ProcessResponse) String() string { return proto1.CompactTextString(m) }
func (*ProcessResponse) ProtoMessage()    {}

func (m *ProcessResponse) GetRowProcessorResult() []byte {
	if m != nil {
		return m.RowProcessorResult
	}
	return nil
}

func init() {
}
