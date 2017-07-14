// Code generated by protoc-gen-go.
// source: raftpb.proto
// DO NOT EDIT!

/*
Package raftpb is a generated protocol buffer package.

It is generated from these files:
	raftpb.proto

It has these top-level messages:
	Entry
	SnapshotMetadata
	Snapshot
	Message
	HardState
	ConfState
	ConfChange
*/
package raftpb

import proto "github.com/insionng/yougam/libraries/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type EntryType int32

const (
	EntryType_EntryNormal     EntryType = 0
	EntryType_EntryConfChange EntryType = 1
)

var EntryType_name = map[int32]string{
	0: "EntryNormal",
	1: "EntryConfChange",
}
var EntryType_value = map[string]int32{
	"EntryNormal":     0,
	"EntryConfChange": 1,
}

func (x EntryType) Enum() *EntryType {
	p := new(EntryType)
	*p = x
	return p
}
func (x EntryType) String() string {
	return proto.EnumName(EntryType_name, int32(x))
}
func (x *EntryType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(EntryType_value, data, "EntryType")
	if err != nil {
		return err
	}
	*x = EntryType(value)
	return nil
}
func (EntryType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type MessageType int32

const (
	MessageType_MsgHup                 MessageType = 0
	MessageType_MsgBeat                MessageType = 1
	MessageType_MsgPropose             MessageType = 2
	MessageType_MsgAppend              MessageType = 3
	MessageType_MsgAppendResponse      MessageType = 4
	MessageType_MsgRequestVote         MessageType = 5
	MessageType_MsgRequestVoteResponse MessageType = 6
	MessageType_MsgSnapshot            MessageType = 7
	MessageType_MsgHeartbeat           MessageType = 8
	MessageType_MsgHeartbeatResponse   MessageType = 9
	MessageType_MsgUnreachable         MessageType = 10
	MessageType_MsgSnapStatus          MessageType = 11
	MessageType_MsgCheckQuorum         MessageType = 12
	MessageType_MsgTransferLeader      MessageType = 13
	MessageType_MsgTimeoutNow          MessageType = 14
)

var MessageType_name = map[int32]string{
	0:  "MsgHup",
	1:  "MsgBeat",
	2:  "MsgPropose",
	3:  "MsgAppend",
	4:  "MsgAppendResponse",
	5:  "MsgRequestVote",
	6:  "MsgRequestVoteResponse",
	7:  "MsgSnapshot",
	8:  "MsgHeartbeat",
	9:  "MsgHeartbeatResponse",
	10: "MsgUnreachable",
	11: "MsgSnapStatus",
	12: "MsgCheckQuorum",
	13: "MsgTransferLeader",
	14: "MsgTimeoutNow",
}
var MessageType_value = map[string]int32{
	"MsgHup":                 0,
	"MsgBeat":                1,
	"MsgPropose":             2,
	"MsgAppend":              3,
	"MsgAppendResponse":      4,
	"MsgRequestVote":         5,
	"MsgRequestVoteResponse": 6,
	"MsgSnapshot":            7,
	"MsgHeartbeat":           8,
	"MsgHeartbeatResponse":   9,
	"MsgUnreachable":         10,
	"MsgSnapStatus":          11,
	"MsgCheckQuorum":         12,
	"MsgTransferLeader":      13,
	"MsgTimeoutNow":          14,
}

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}
func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}
func (x *MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MessageType_value, data, "MessageType")
	if err != nil {
		return err
	}
	*x = MessageType(value)
	return nil
}
func (MessageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ConfChangeType int32

const (
	ConfChangeType_AddNode    ConfChangeType = 0
	ConfChangeType_RemoveNode ConfChangeType = 1
)

var ConfChangeType_name = map[int32]string{
	0: "AddNode",
	1: "RemoveNode",
}
var ConfChangeType_value = map[string]int32{
	"AddNode":    0,
	"RemoveNode": 1,
}

func (x ConfChangeType) Enum() *ConfChangeType {
	p := new(ConfChangeType)
	*p = x
	return p
}
func (x ConfChangeType) String() string {
	return proto.EnumName(ConfChangeType_name, int32(x))
}
func (x *ConfChangeType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ConfChangeType_value, data, "ConfChangeType")
	if err != nil {
		return err
	}
	*x = ConfChangeType(value)
	return nil
}
func (ConfChangeType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Entry struct {
	EntryType        *EntryType `protobuf:"varint,1,opt,name=entry_type,enum=raftpb.EntryType" json:"entry_type,omitempty"`
	Term             *uint64    `protobuf:"varint,2,opt,name=term" json:"term,omitempty"`
	Index            *uint64    `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
	Data             []byte     `protobuf:"bytes,4,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *Entry) Reset()                    { *m = Entry{} }
func (m *Entry) String() string            { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()               {}
func (*Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Entry) GetEntryType() EntryType {
	if m != nil && m.EntryType != nil {
		return *m.EntryType
	}
	return EntryType_EntryNormal
}

func (m *Entry) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *Entry) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *Entry) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type SnapshotMetadata struct {
	ConfState        *ConfState `protobuf:"bytes,1,opt,name=conf_state" json:"conf_state,omitempty"`
	Index            *uint64    `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	Term             *uint64    `protobuf:"varint,3,opt,name=term" json:"term,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *SnapshotMetadata) Reset()                    { *m = SnapshotMetadata{} }
func (m *SnapshotMetadata) String() string            { return proto.CompactTextString(m) }
func (*SnapshotMetadata) ProtoMessage()               {}
func (*SnapshotMetadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SnapshotMetadata) GetConfState() *ConfState {
	if m != nil {
		return m.ConfState
	}
	return nil
}

func (m *SnapshotMetadata) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *SnapshotMetadata) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

type Snapshot struct {
	Data             []byte            `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
	Metadata         *SnapshotMetadata `protobuf:"bytes,2,opt,name=metadata" json:"metadata,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Snapshot) Reset()                    { *m = Snapshot{} }
func (m *Snapshot) String() string            { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()               {}
func (*Snapshot) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Snapshot) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Snapshot) GetMetadata() *SnapshotMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Message struct {
	MsgType          *MessageType `protobuf:"varint,1,opt,name=msg_type,enum=raftpb.MessageType" json:"msg_type,omitempty"`
	To               *uint64      `protobuf:"varint,2,opt,name=to" json:"to,omitempty"`
	From             *uint64      `protobuf:"varint,3,opt,name=from" json:"from,omitempty"`
	Term             *uint64      `protobuf:"varint,4,opt,name=term" json:"term,omitempty"`
	LogTerm          *uint64      `protobuf:"varint,5,opt,name=log_term" json:"log_term,omitempty"`
	Index            *uint64      `protobuf:"varint,6,opt,name=index" json:"index,omitempty"`
	Entries          []*Entry     `protobuf:"bytes,7,rep,name=entries" json:"entries,omitempty"`
	Commit           *uint64      `protobuf:"varint,8,opt,name=commit" json:"commit,omitempty"`
	Snapshot         *Snapshot    `protobuf:"bytes,9,opt,name=snapshot" json:"snapshot,omitempty"`
	Reject           *bool        `protobuf:"varint,10,opt,name=reject" json:"reject,omitempty"`
	RejectHint       *uint64      `protobuf:"varint,11,opt,name=reject_hint" json:"reject_hint,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Message) GetMsgType() MessageType {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return MessageType_MsgHup
}

func (m *Message) GetTo() uint64 {
	if m != nil && m.To != nil {
		return *m.To
	}
	return 0
}

func (m *Message) GetFrom() uint64 {
	if m != nil && m.From != nil {
		return *m.From
	}
	return 0
}

func (m *Message) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *Message) GetLogTerm() uint64 {
	if m != nil && m.LogTerm != nil {
		return *m.LogTerm
	}
	return 0
}

func (m *Message) GetIndex() uint64 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

func (m *Message) GetEntries() []*Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *Message) GetCommit() uint64 {
	if m != nil && m.Commit != nil {
		return *m.Commit
	}
	return 0
}

func (m *Message) GetSnapshot() *Snapshot {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

func (m *Message) GetReject() bool {
	if m != nil && m.Reject != nil {
		return *m.Reject
	}
	return false
}

func (m *Message) GetRejectHint() uint64 {
	if m != nil && m.RejectHint != nil {
		return *m.RejectHint
	}
	return 0
}

type HardState struct {
	Term             *uint64 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	Vote             *uint64 `protobuf:"varint,2,opt,name=vote" json:"vote,omitempty"`
	Commit           *uint64 `protobuf:"varint,3,opt,name=commit" json:"commit,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *HardState) Reset()                    { *m = HardState{} }
func (m *HardState) String() string            { return proto.CompactTextString(m) }
func (*HardState) ProtoMessage()               {}
func (*HardState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *HardState) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *HardState) GetVote() uint64 {
	if m != nil && m.Vote != nil {
		return *m.Vote
	}
	return 0
}

func (m *HardState) GetCommit() uint64 {
	if m != nil && m.Commit != nil {
		return *m.Commit
	}
	return 0
}

type ConfState struct {
	Nodes            []uint64 `protobuf:"varint,1,rep,name=nodes" json:"nodes,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ConfState) Reset()                    { *m = ConfState{} }
func (m *ConfState) String() string            { return proto.CompactTextString(m) }
func (*ConfState) ProtoMessage()               {}
func (*ConfState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ConfState) GetNodes() []uint64 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type ConfChange struct {
	Id               *uint64         `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	ChangeType       *ConfChangeType `protobuf:"varint,2,opt,name=change_type,enum=raftpb.ConfChangeType" json:"change_type,omitempty"`
	NodeId           *uint64         `protobuf:"varint,3,opt,name=node_id" json:"node_id,omitempty"`
	Context          []byte          `protobuf:"bytes,4,opt,name=context" json:"context,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *ConfChange) Reset()                    { *m = ConfChange{} }
func (m *ConfChange) String() string            { return proto.CompactTextString(m) }
func (*ConfChange) ProtoMessage()               {}
func (*ConfChange) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ConfChange) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *ConfChange) GetChangeType() ConfChangeType {
	if m != nil && m.ChangeType != nil {
		return *m.ChangeType
	}
	return ConfChangeType_AddNode
}

func (m *ConfChange) GetNodeId() uint64 {
	if m != nil && m.NodeId != nil {
		return *m.NodeId
	}
	return 0
}

func (m *ConfChange) GetContext() []byte {
	if m != nil {
		return m.Context
	}
	return nil
}

func init() {
	proto.RegisterType((*Entry)(nil), "raftpb.Entry")
	proto.RegisterType((*SnapshotMetadata)(nil), "raftpb.SnapshotMetadata")
	proto.RegisterType((*Snapshot)(nil), "raftpb.Snapshot")
	proto.RegisterType((*Message)(nil), "raftpb.Message")
	proto.RegisterType((*HardState)(nil), "raftpb.HardState")
	proto.RegisterType((*ConfState)(nil), "raftpb.ConfState")
	proto.RegisterType((*ConfChange)(nil), "raftpb.ConfChange")
	proto.RegisterEnum("raftpb.EntryType", EntryType_name, EntryType_value)
	proto.RegisterEnum("raftpb.MessageType", MessageType_name, MessageType_value)
	proto.RegisterEnum("raftpb.ConfChangeType", ConfChangeType_name, ConfChangeType_value)
}

var fileDescriptor0 = []byte{
	// 605 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x52, 0xcd, 0x4e, 0xdb, 0x4c,
	0x14, 0xfd, 0x1c, 0xf2, 0xe3, 0x5c, 0x3b, 0x61, 0x18, 0xbe, 0x22, 0x8b, 0x45, 0x85, 0x2c, 0x55,
	0xaa, 0x52, 0x15, 0xa9, 0x6c, 0xba, 0xa6, 0xb4, 0x12, 0x8b, 0x82, 0x5a, 0xa0, 0x48, 0x5d, 0x45,
	0x43, 0x7c, 0xe3, 0xb8, 0xc5, 0x1e, 0x77, 0x66, 0x4c, 0xe1, 0xa5, 0xfa, 0x70, 0x7d, 0x82, 0xde,
	0x19, 0xdb, 0x31, 0xc9, 0x6e, 0xce, 0xfd, 0x3d, 0xe7, 0xdc, 0x81, 0x50, 0x89, 0xa5, 0x29, 0xef,
	0x8e, 0x4b, 0x25, 0x8d, 0xe4, 0xc3, 0x1a, 0xc5, 0xdf, 0x61, 0xf0, 0xa9, 0x30, 0xea, 0x89, 0xbf,
	0x02, 0x40, 0xfb, 0x98, 0x9b, 0xa7, 0x12, 0x23, 0xef, 0xc8, 0x7b, 0x3d, 0x3d, 0xd9, 0x3b, 0x6e,
	0x7a, 0x5c, 0xc9, 0x0d, 0x25, 0x78, 0x08, 0x7d, 0x83, 0x2a, 0x8f, 0x7a, 0x54, 0xd0, 0xe7, 0x13,
	0x18, 0x64, 0x45, 0x82, 0x8f, 0xd1, 0x8e, 0x83, 0x94, 0x4c, 0x84, 0x11, 0x51, 0x9f, 0x50, 0x18,
	0xdf, 0x02, 0xbb, 0x2e, 0x44, 0xa9, 0x57, 0xd2, 0x5c, 0xa0, 0x11, 0x36, 0x63, 0xb7, 0x2c, 0x64,
	0xb1, 0x9c, 0x6b, 0x23, 0x4c, 0xbd, 0x25, 0xe8, 0xb6, 0x9c, 0x51, 0xe6, 0xda, 0x26, 0xba, 0xb9,
	0xbd, 0x76, 0xae, 0x5b, 0xea, 0xb6, 0xc4, 0x1f, 0xc1, 0x6f, 0xe7, 0xae, 0x37, 0xda, 0x49, 0x21,
	0x9f, 0x81, 0x9f, 0x37, 0x9b, 0x5c, 0x67, 0x70, 0x12, 0xb5, 0xb3, 0xb7, 0x99, 0xc4, 0x7f, 0x3d,
	0x18, 0x5d, 0xa0, 0xd6, 0x22, 0x45, 0x62, 0xe5, 0xe7, 0x3a, 0x7d, 0xae, 0x7c, 0xbf, 0xed, 0x6b,
	0x4a, 0x9c, 0x76, 0x80, 0x9e, 0x91, 0x1d, 0xa5, 0xa5, 0x92, 0x79, 0x27, 0xdc, 0x11, 0xec, 0x3b,
	0xc4, 0xc0, 0xbf, 0x97, 0x34, 0xce, 0x46, 0x06, 0x9b, 0x3e, 0x0d, 0x1d, 0x7c, 0x09, 0x23, 0xeb,
	0x75, 0x86, 0x3a, 0x1a, 0x1d, 0xed, 0x10, 0xcd, 0xc9, 0x86, 0xd1, 0x7c, 0x0a, 0xc3, 0x85, 0xcc,
	0xf3, 0xcc, 0x44, 0xbe, 0xab, 0x8f, 0xc1, 0xd7, 0x0d, 0xff, 0x68, 0xec, 0x74, 0xb1, 0x6d, 0x5d,
	0xb6, 0x47, 0xe1, 0x0f, 0x5c, 0x98, 0x08, 0xa8, 0xc2, 0xe7, 0xfb, 0x10, 0xd4, 0x78, 0xbe, 0xca,
	0x0a, 0x13, 0x05, 0xce, 0xba, 0xf7, 0x30, 0x3e, 0x17, 0x2a, 0xa9, 0x4d, 0x6e, 0x49, 0x7b, 0xad,
	0x84, 0x07, 0x49, 0x37, 0xa9, 0xe5, 0x75, 0x0c, 0x6a, 0xcf, 0x0f, 0x61, 0xbc, 0x71, 0x9d, 0x42,
	0x26, 0x44, 0xde, 0x23, 0xf2, 0xfd, 0x38, 0x05, 0xb0, 0xb9, 0xb3, 0x95, 0x28, 0x52, 0x67, 0x52,
	0x96, 0x34, 0x33, 0xdf, 0x40, 0xb0, 0x70, 0xd1, 0xda, 0xda, 0x9e, 0xb3, 0xf6, 0xe0, 0xf9, 0xb9,
	0xeb, 0x26, 0xe7, 0xee, 0x2e, 0x8c, 0xec, 0xd4, 0x39, 0x75, 0xd7, 0xa6, 0x52, 0x80, 0xfe, 0x8a,
	0xc1, 0x47, 0x53, 0x7f, 0xa8, 0xd9, 0x3b, 0x18, 0x77, 0x1f, 0x71, 0x17, 0x02, 0x07, 0x2e, 0xa5,
	0xca, 0xc5, 0x3d, 0xfb, 0x8f, 0x04, 0xef, 0xba, 0x40, 0x37, 0x96, 0x79, 0xb3, 0x3f, 0x3d, 0x08,
	0x36, 0x4f, 0x38, 0xbc, 0xd0, 0xe9, 0x79, 0x55, 0x52, 0x43, 0x40, 0x1f, 0x40, 0xa7, 0x1f, 0x50,
	0x18, 0xe6, 0x91, 0x60, 0x20, 0xf0, 0x45, 0xc9, 0x52, 0x6a, 0x64, 0x3d, 0xd2, 0x38, 0x26, 0x7c,
	0x5a, 0x96, 0x58, 0x24, 0x6c, 0x87, 0xbf, 0x80, 0xbd, 0x35, 0xbc, 0x42, 0x5d, 0xca, 0x82, 0xaa,
	0xfa, 0x9c, 0xc3, 0x94, 0xc2, 0x57, 0xf8, 0xab, 0x42, 0x6d, 0x6e, 0xc9, 0x3e, 0x36, 0xe0, 0x87,
	0x70, 0xb0, 0x19, 0x5b, 0xd7, 0x0f, 0x2d, 0x69, 0xca, 0xb5, 0x37, 0x63, 0x23, 0xfa, 0x2a, 0xa1,
	0xe5, 0x83, 0x42, 0x99, 0x3b, 0x4b, 0xc4, 0xe7, 0x11, 0xfc, 0xff, 0x3c, 0xb2, 0x6e, 0x1e, 0x37,
	0xcb, 0xbe, 0x15, 0x0a, 0x05, 0xb9, 0x7a, 0x77, 0x8f, 0x0c, 0xf8, 0x1e, 0x4c, 0x9a, 0x81, 0xf6,
	0x34, 0x95, 0x66, 0x41, 0x53, 0x76, 0xb6, 0xc2, 0xc5, 0xcf, 0xaf, 0x95, 0x54, 0x55, 0xce, 0xc2,
	0x86, 0xfe, 0x8d, 0x12, 0x85, 0x5e, 0xa2, 0xfa, 0x8c, 0x22, 0x41, 0xc5, 0x26, 0x4d, 0xf7, 0x4d,
	0x96, 0xa3, 0xac, 0xcc, 0xa5, 0xfc, 0xcd, 0xa6, 0xb3, 0xb7, 0x30, 0xdd, 0xba, 0x0b, 0xd9, 0x74,
	0x9a, 0x24, 0x97, 0x74, 0x1a, 0xf2, 0x8c, 0x6c, 0xba, 0xc2, 0x5c, 0x3e, 0xa0, 0xc3, 0xde, 0xbf,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x2d, 0x25, 0x6c, 0xc8, 0x55, 0x04, 0x00, 0x00,
}
