// Code generated by protoc-gen-go.
// source: RPC.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/insionng/yougam/libraries/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

// User Information proto.  Included in ConnectionHeader on connection setup
type UserInformation struct {
	EffectiveUser    *string `protobuf:"bytes,1,req,name=effective_user" json:"effective_user,omitempty"`
	RealUser         *string `protobuf:"bytes,2,opt,name=real_user" json:"real_user,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *UserInformation) Reset()         { *m = UserInformation{} }
func (m *UserInformation) String() string { return proto1.CompactTextString(m) }
func (*UserInformation) ProtoMessage()    {}

func (m *UserInformation) GetEffectiveUser() string {
	if m != nil && m.EffectiveUser != nil {
		return *m.EffectiveUser
	}
	return ""
}

func (m *UserInformation) GetRealUser() string {
	if m != nil && m.RealUser != nil {
		return *m.RealUser
	}
	return ""
}

// Rpc client version info proto. Included in ConnectionHeader on connection setup
type VersionInfo struct {
	Version          *string `protobuf:"bytes,1,req,name=version" json:"version,omitempty"`
	Url              *string `protobuf:"bytes,2,req,name=url" json:"url,omitempty"`
	Revision         *string `protobuf:"bytes,3,req,name=revision" json:"revision,omitempty"`
	User             *string `protobuf:"bytes,4,req,name=user" json:"user,omitempty"`
	Date             *string `protobuf:"bytes,5,req,name=date" json:"date,omitempty"`
	SrcChecksum      *string `protobuf:"bytes,6,req,name=src_checksum" json:"src_checksum,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *VersionInfo) Reset()         { *m = VersionInfo{} }
func (m *VersionInfo) String() string { return proto1.CompactTextString(m) }
func (*VersionInfo) ProtoMessage()    {}

func (m *VersionInfo) GetVersion() string {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return ""
}

func (m *VersionInfo) GetUrl() string {
	if m != nil && m.Url != nil {
		return *m.Url
	}
	return ""
}

func (m *VersionInfo) GetRevision() string {
	if m != nil && m.Revision != nil {
		return *m.Revision
	}
	return ""
}

func (m *VersionInfo) GetUser() string {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

func (m *VersionInfo) GetDate() string {
	if m != nil && m.Date != nil {
		return *m.Date
	}
	return ""
}

func (m *VersionInfo) GetSrcChecksum() string {
	if m != nil && m.SrcChecksum != nil {
		return *m.SrcChecksum
	}
	return ""
}

// This is sent on connection setup after the connection preamble is sent.
type ConnectionHeader struct {
	UserInfo    *UserInformation `protobuf:"bytes,1,opt,name=user_info" json:"user_info,omitempty"`
	ServiceName *string          `protobuf:"bytes,2,opt,name=service_name" json:"service_name,omitempty"`
	// Cell block codec we will use sending over optional cell blocks.  Server throws exception
	// if cannot deal.  Null means no codec'ing going on so we are pb all the time (SLOW!!!)
	CellBlockCodecClass *string `protobuf:"bytes,3,opt,name=cell_block_codec_class" json:"cell_block_codec_class,omitempty"`
	// Compressor we will use if cell block is compressed.  Server will throw exception if not supported.
	// Class must implement hadoop's CompressionCodec Interface.  Can't compress if no codec.
	CellBlockCompressorClass *string      `protobuf:"bytes,4,opt,name=cell_block_compressor_class" json:"cell_block_compressor_class,omitempty"`
	VersionInfo              *VersionInfo `protobuf:"bytes,5,opt,name=version_info" json:"version_info,omitempty"`
	XXX_unrecognized         []byte       `json:"-"`
}

func (m *ConnectionHeader) Reset()         { *m = ConnectionHeader{} }
func (m *ConnectionHeader) String() string { return proto1.CompactTextString(m) }
func (*ConnectionHeader) ProtoMessage()    {}

func (m *ConnectionHeader) GetUserInfo() *UserInformation {
	if m != nil {
		return m.UserInfo
	}
	return nil
}

func (m *ConnectionHeader) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *ConnectionHeader) GetCellBlockCodecClass() string {
	if m != nil && m.CellBlockCodecClass != nil {
		return *m.CellBlockCodecClass
	}
	return ""
}

func (m *ConnectionHeader) GetCellBlockCompressorClass() string {
	if m != nil && m.CellBlockCompressorClass != nil {
		return *m.CellBlockCompressorClass
	}
	return ""
}

func (m *ConnectionHeader) GetVersionInfo() *VersionInfo {
	if m != nil {
		return m.VersionInfo
	}
	return nil
}

// Optional Cell block Message.  Included in client RequestHeader
type CellBlockMeta struct {
	// Length of the following cell block.  Could calculate it but convenient having it too hand.
	Length           *uint32 `protobuf:"varint,1,opt,name=length" json:"length,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CellBlockMeta) Reset()         { *m = CellBlockMeta{} }
func (m *CellBlockMeta) String() string { return proto1.CompactTextString(m) }
func (*CellBlockMeta) ProtoMessage()    {}

func (m *CellBlockMeta) GetLength() uint32 {
	if m != nil && m.Length != nil {
		return *m.Length
	}
	return 0
}

// At the RPC layer, this message is used to carry
// the server side exception to the RPC client.
type ExceptionResponse struct {
	// Class name of the exception thrown from the server
	ExceptionClassName *string `protobuf:"bytes,1,opt,name=exception_class_name" json:"exception_class_name,omitempty"`
	// Exception stack trace from the server side
	StackTrace *string `protobuf:"bytes,2,opt,name=stack_trace" json:"stack_trace,omitempty"`
	// Optional hostname.  Filled in for some exceptions such as region moved
	// where exception gives clue on where the region may have moved.
	Hostname *string `protobuf:"bytes,3,opt,name=hostname" json:"hostname,omitempty"`
	Port     *int32  `protobuf:"varint,4,opt,name=port" json:"port,omitempty"`
	// Set if we are NOT to retry on receipt of this exception
	DoNotRetry       *bool  `protobuf:"varint,5,opt,name=do_not_retry" json:"do_not_retry,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ExceptionResponse) Reset()         { *m = ExceptionResponse{} }
func (m *ExceptionResponse) String() string { return proto1.CompactTextString(m) }
func (*ExceptionResponse) ProtoMessage()    {}

func (m *ExceptionResponse) GetExceptionClassName() string {
	if m != nil && m.ExceptionClassName != nil {
		return *m.ExceptionClassName
	}
	return ""
}

func (m *ExceptionResponse) GetStackTrace() string {
	if m != nil && m.StackTrace != nil {
		return *m.StackTrace
	}
	return ""
}

func (m *ExceptionResponse) GetHostname() string {
	if m != nil && m.Hostname != nil {
		return *m.Hostname
	}
	return ""
}

func (m *ExceptionResponse) GetPort() int32 {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return 0
}

func (m *ExceptionResponse) GetDoNotRetry() bool {
	if m != nil && m.DoNotRetry != nil {
		return *m.DoNotRetry
	}
	return false
}

// Header sent making a request.
type RequestHeader struct {
	// Monotonically increasing call_id to keep track of RPC requests and their response
	CallId     *uint32   `protobuf:"varint,1,opt,name=call_id" json:"call_id,omitempty"`
	TraceInfo  *RPCTInfo `protobuf:"bytes,2,opt,name=trace_info" json:"trace_info,omitempty"`
	MethodName *string   `protobuf:"bytes,3,opt,name=method_name" json:"method_name,omitempty"`
	// If true, then a pb Message param follows.
	RequestParam *bool `protobuf:"varint,4,opt,name=request_param" json:"request_param,omitempty"`
	// If present, then an encoded data block follows.
	CellBlockMeta *CellBlockMeta `protobuf:"bytes,5,opt,name=cell_block_meta" json:"cell_block_meta,omitempty"`
	// 0 is NORMAL priority.  100 is HIGH.  If no priority, treat it as NORMAL.
	// See HConstants.
	Priority         *uint32 `protobuf:"varint,6,opt,name=priority" json:"priority,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RequestHeader) Reset()         { *m = RequestHeader{} }
func (m *RequestHeader) String() string { return proto1.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()    {}

func (m *RequestHeader) GetCallId() uint32 {
	if m != nil && m.CallId != nil {
		return *m.CallId
	}
	return 0
}

func (m *RequestHeader) GetTraceInfo() *RPCTInfo {
	if m != nil {
		return m.TraceInfo
	}
	return nil
}

func (m *RequestHeader) GetMethodName() string {
	if m != nil && m.MethodName != nil {
		return *m.MethodName
	}
	return ""
}

func (m *RequestHeader) GetRequestParam() bool {
	if m != nil && m.RequestParam != nil {
		return *m.RequestParam
	}
	return false
}

func (m *RequestHeader) GetCellBlockMeta() *CellBlockMeta {
	if m != nil {
		return m.CellBlockMeta
	}
	return nil
}

func (m *RequestHeader) GetPriority() uint32 {
	if m != nil && m.Priority != nil {
		return *m.Priority
	}
	return 0
}

type ResponseHeader struct {
	CallId *uint32 `protobuf:"varint,1,opt,name=call_id" json:"call_id,omitempty"`
	// If present, then request threw an exception and no response message (else we presume one)
	Exception *ExceptionResponse `protobuf:"bytes,2,opt,name=exception" json:"exception,omitempty"`
	// If present, then an encoded data block follows.
	CellBlockMeta    *CellBlockMeta `protobuf:"bytes,3,opt,name=cell_block_meta" json:"cell_block_meta,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *ResponseHeader) Reset()         { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string { return proto1.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()    {}

func (m *ResponseHeader) GetCallId() uint32 {
	if m != nil && m.CallId != nil {
		return *m.CallId
	}
	return 0
}

func (m *ResponseHeader) GetException() *ExceptionResponse {
	if m != nil {
		return m.Exception
	}
	return nil
}

func (m *ResponseHeader) GetCellBlockMeta() *CellBlockMeta {
	if m != nil {
		return m.CellBlockMeta
	}
	return nil
}

func init() {
}
