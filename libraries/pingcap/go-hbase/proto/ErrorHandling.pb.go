// Code generated by protoc-gen-go.
// source: ErrorHandling.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/insionng/yougam/libraries/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

// *
// Protobuf version of a java.lang.StackTraceElement
// so we can serialize exceptions.
type StackTraceElementMessage struct {
	DeclaringClass   *string `protobuf:"bytes,1,opt,name=declaring_class" json:"declaring_class,omitempty"`
	MethodName       *string `protobuf:"bytes,2,opt,name=method_name" json:"method_name,omitempty"`
	FileName         *string `protobuf:"bytes,3,opt,name=file_name" json:"file_name,omitempty"`
	LineNumber       *int32  `protobuf:"varint,4,opt,name=line_number" json:"line_number,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *StackTraceElementMessage) Reset()         { *m = StackTraceElementMessage{} }
func (m *StackTraceElementMessage) String() string { return proto1.CompactTextString(m) }
func (*StackTraceElementMessage) ProtoMessage()    {}

func (m *StackTraceElementMessage) GetDeclaringClass() string {
	if m != nil && m.DeclaringClass != nil {
		return *m.DeclaringClass
	}
	return ""
}

func (m *StackTraceElementMessage) GetMethodName() string {
	if m != nil && m.MethodName != nil {
		return *m.MethodName
	}
	return ""
}

func (m *StackTraceElementMessage) GetFileName() string {
	if m != nil && m.FileName != nil {
		return *m.FileName
	}
	return ""
}

func (m *StackTraceElementMessage) GetLineNumber() int32 {
	if m != nil && m.LineNumber != nil {
		return *m.LineNumber
	}
	return 0
}

// *
// Cause of a remote failure for a generic exception. Contains
// all the information for a generic exception as well as
// optional info about the error for generic info passing
// (which should be another protobuffed class).
type GenericExceptionMessage struct {
	ClassName        *string                     `protobuf:"bytes,1,opt,name=class_name" json:"class_name,omitempty"`
	Message          *string                     `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	ErrorInfo        []byte                      `protobuf:"bytes,3,opt,name=error_info" json:"error_info,omitempty"`
	Trace            []*StackTraceElementMessage `protobuf:"bytes,4,rep,name=trace" json:"trace,omitempty"`
	XXX_unrecognized []byte                      `json:"-"`
}

func (m *GenericExceptionMessage) Reset()         { *m = GenericExceptionMessage{} }
func (m *GenericExceptionMessage) String() string { return proto1.CompactTextString(m) }
func (*GenericExceptionMessage) ProtoMessage()    {}

func (m *GenericExceptionMessage) GetClassName() string {
	if m != nil && m.ClassName != nil {
		return *m.ClassName
	}
	return ""
}

func (m *GenericExceptionMessage) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

func (m *GenericExceptionMessage) GetErrorInfo() []byte {
	if m != nil {
		return m.ErrorInfo
	}
	return nil
}

func (m *GenericExceptionMessage) GetTrace() []*StackTraceElementMessage {
	if m != nil {
		return m.Trace
	}
	return nil
}

// *
// Exception sent across the wire when a remote task needs
// to notify other tasks that it failed and why
type ForeignExceptionMessage struct {
	Source           *string                  `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	GenericException *GenericExceptionMessage `protobuf:"bytes,2,opt,name=generic_exception" json:"generic_exception,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *ForeignExceptionMessage) Reset()         { *m = ForeignExceptionMessage{} }
func (m *ForeignExceptionMessage) String() string { return proto1.CompactTextString(m) }
func (*ForeignExceptionMessage) ProtoMessage()    {}

func (m *ForeignExceptionMessage) GetSource() string {
	if m != nil && m.Source != nil {
		return *m.Source
	}
	return ""
}

func (m *ForeignExceptionMessage) GetGenericException() *GenericExceptionMessage {
	if m != nil {
		return m.GenericException
	}
	return nil
}

func init() {
}
