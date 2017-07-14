// Code generated by protoc-gen-go.
// source: SecureBulkLoad.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/insionng/yougam/libraries/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

type SecureBulkLoadHFilesRequest struct {
	FamilyPath       []*BulkLoadHFileRequest_FamilyPath `protobuf:"bytes,1,rep,name=family_path" json:"family_path,omitempty"`
	AssignSeqNum     *bool                              `protobuf:"varint,2,opt,name=assign_seq_num" json:"assign_seq_num,omitempty"`
	FsToken          *DelegationToken                   `protobuf:"bytes,3,req,name=fs_token" json:"fs_token,omitempty"`
	BulkToken        *string                            `protobuf:"bytes,4,req,name=bulk_token" json:"bulk_token,omitempty"`
	XXX_unrecognized []byte                             `json:"-"`
}

func (m *SecureBulkLoadHFilesRequest) Reset()         { *m = SecureBulkLoadHFilesRequest{} }
func (m *SecureBulkLoadHFilesRequest) String() string { return proto1.CompactTextString(m) }
func (*SecureBulkLoadHFilesRequest) ProtoMessage()    {}

func (m *SecureBulkLoadHFilesRequest) GetFamilyPath() []*BulkLoadHFileRequest_FamilyPath {
	if m != nil {
		return m.FamilyPath
	}
	return nil
}

func (m *SecureBulkLoadHFilesRequest) GetAssignSeqNum() bool {
	if m != nil && m.AssignSeqNum != nil {
		return *m.AssignSeqNum
	}
	return false
}

func (m *SecureBulkLoadHFilesRequest) GetFsToken() *DelegationToken {
	if m != nil {
		return m.FsToken
	}
	return nil
}

func (m *SecureBulkLoadHFilesRequest) GetBulkToken() string {
	if m != nil && m.BulkToken != nil {
		return *m.BulkToken
	}
	return ""
}

type SecureBulkLoadHFilesResponse struct {
	Loaded           *bool  `protobuf:"varint,1,req,name=loaded" json:"loaded,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *SecureBulkLoadHFilesResponse) Reset()         { *m = SecureBulkLoadHFilesResponse{} }
func (m *SecureBulkLoadHFilesResponse) String() string { return proto1.CompactTextString(m) }
func (*SecureBulkLoadHFilesResponse) ProtoMessage()    {}

func (m *SecureBulkLoadHFilesResponse) GetLoaded() bool {
	if m != nil && m.Loaded != nil {
		return *m.Loaded
	}
	return false
}

type DelegationToken struct {
	Identifier       []byte  `protobuf:"bytes,1,opt,name=identifier" json:"identifier,omitempty"`
	Password         []byte  `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
	Kind             *string `protobuf:"bytes,3,opt,name=kind" json:"kind,omitempty"`
	Service          *string `protobuf:"bytes,4,opt,name=service" json:"service,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DelegationToken) Reset()         { *m = DelegationToken{} }
func (m *DelegationToken) String() string { return proto1.CompactTextString(m) }
func (*DelegationToken) ProtoMessage()    {}

func (m *DelegationToken) GetIdentifier() []byte {
	if m != nil {
		return m.Identifier
	}
	return nil
}

func (m *DelegationToken) GetPassword() []byte {
	if m != nil {
		return m.Password
	}
	return nil
}

func (m *DelegationToken) GetKind() string {
	if m != nil && m.Kind != nil {
		return *m.Kind
	}
	return ""
}

func (m *DelegationToken) GetService() string {
	if m != nil && m.Service != nil {
		return *m.Service
	}
	return ""
}

type PrepareBulkLoadRequest struct {
	TableName        *TableName `protobuf:"bytes,1,req,name=table_name" json:"table_name,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *PrepareBulkLoadRequest) Reset()         { *m = PrepareBulkLoadRequest{} }
func (m *PrepareBulkLoadRequest) String() string { return proto1.CompactTextString(m) }
func (*PrepareBulkLoadRequest) ProtoMessage()    {}

func (m *PrepareBulkLoadRequest) GetTableName() *TableName {
	if m != nil {
		return m.TableName
	}
	return nil
}

type PrepareBulkLoadResponse struct {
	BulkToken        *string `protobuf:"bytes,1,req,name=bulk_token" json:"bulk_token,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PrepareBulkLoadResponse) Reset()         { *m = PrepareBulkLoadResponse{} }
func (m *PrepareBulkLoadResponse) String() string { return proto1.CompactTextString(m) }
func (*PrepareBulkLoadResponse) ProtoMessage()    {}

func (m *PrepareBulkLoadResponse) GetBulkToken() string {
	if m != nil && m.BulkToken != nil {
		return *m.BulkToken
	}
	return ""
}

type CleanupBulkLoadRequest struct {
	BulkToken        *string `protobuf:"bytes,1,req,name=bulk_token" json:"bulk_token,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CleanupBulkLoadRequest) Reset()         { *m = CleanupBulkLoadRequest{} }
func (m *CleanupBulkLoadRequest) String() string { return proto1.CompactTextString(m) }
func (*CleanupBulkLoadRequest) ProtoMessage()    {}

func (m *CleanupBulkLoadRequest) GetBulkToken() string {
	if m != nil && m.BulkToken != nil {
		return *m.BulkToken
	}
	return ""
}

type CleanupBulkLoadResponse struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *CleanupBulkLoadResponse) Reset()         { *m = CleanupBulkLoadResponse{} }
func (m *CleanupBulkLoadResponse) String() string { return proto1.CompactTextString(m) }
func (*CleanupBulkLoadResponse) ProtoMessage()    {}

func init() {
}
