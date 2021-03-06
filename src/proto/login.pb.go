// Code generated by protoc-gen-go.
// source: src/netpb/login.proto
// DO NOT EDIT!

/*
Package netpb is a generated protocol buffer package.

It is generated from these files:
	src/netpb/login.proto

It has these top-level messages:
	NetMsg
	CSNetMsg
	SCNetMsg
	CSLoginReq
	SCLoginRsp
*/
package netpb

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type MsgID int32

const (
	MsgID_EM_CS_LOGIN MsgID = 0
	MsgID_EM_SC_LOGIN MsgID = 1
)

var MsgID_name = map[int32]string{
	0: "EM_CS_LOGIN",
	1: "EM_SC_LOGIN",
}
var MsgID_value = map[string]int32{
	"EM_CS_LOGIN": 0,
	"EM_SC_LOGIN": 1,
}

func (x MsgID) Enum() *MsgID {
	p := new(MsgID)
	*p = x
	return p
}
func (x MsgID) String() string {
	return proto.EnumName(MsgID_name, int32(x))
}
func (x *MsgID) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MsgID_value, data, "MsgID")
	if err != nil {
		return err
	}
	*x = MsgID(value)
	return nil
}

type SCLoginRsp_LoginErrorType int32

const (
	SCLoginRsp_SUCCESS   SCLoginRsp_LoginErrorType = 0
	SCLoginRsp_PWD_ERROR SCLoginRsp_LoginErrorType = 1
)

var SCLoginRsp_LoginErrorType_name = map[int32]string{
	0: "SUCCESS",
	1: "PWD_ERROR",
}
var SCLoginRsp_LoginErrorType_value = map[string]int32{
	"SUCCESS":   0,
	"PWD_ERROR": 1,
}

func (x SCLoginRsp_LoginErrorType) Enum() *SCLoginRsp_LoginErrorType {
	p := new(SCLoginRsp_LoginErrorType)
	*p = x
	return p
}
func (x SCLoginRsp_LoginErrorType) String() string {
	return proto.EnumName(SCLoginRsp_LoginErrorType_name, int32(x))
}
func (x *SCLoginRsp_LoginErrorType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SCLoginRsp_LoginErrorType_value, data, "SCLoginRsp_LoginErrorType")
	if err != nil {
		return err
	}
	*x = SCLoginRsp_LoginErrorType(value)
	return nil
}

type NetMsg struct {
	M_ID             *MsgID    `protobuf:"varint,1,req,name=m_ID,enum=netpb.MsgID" json:"m_ID,omitempty"`
	M_CSNetMsg       *CSNetMsg `protobuf:"bytes,2,opt,name=m_CSNetMsg" json:"m_CSNetMsg,omitempty"`
	M_SCNetMsg       *SCNetMsg `protobuf:"bytes,3,opt,name=m_SCNetMsg" json:"m_SCNetMsg,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *NetMsg) Reset()         { *m = NetMsg{} }
func (m *NetMsg) String() string { return proto.CompactTextString(m) }
func (*NetMsg) ProtoMessage()    {}

func (m *NetMsg) GetM_ID() MsgID {
	if m != nil && m.M_ID != nil {
		return *m.M_ID
	}
	return MsgID_EM_CS_LOGIN
}

func (m *NetMsg) GetM_CSNetMsg() *CSNetMsg {
	if m != nil {
		return m.M_CSNetMsg
	}
	return nil
}

func (m *NetMsg) GetM_SCNetMsg() *SCNetMsg {
	if m != nil {
		return m.M_SCNetMsg
	}
	return nil
}

type CSNetMsg struct {
	M_CSLoginReq     *CSLoginReq `protobuf:"bytes,1,opt,name=m_CSLoginReq" json:"m_CSLoginReq,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *CSNetMsg) Reset()         { *m = CSNetMsg{} }
func (m *CSNetMsg) String() string { return proto.CompactTextString(m) }
func (*CSNetMsg) ProtoMessage()    {}

func (m *CSNetMsg) GetM_CSLoginReq() *CSLoginReq {
	if m != nil {
		return m.M_CSLoginReq
	}
	return nil
}

type SCNetMsg struct {
	M_SCLoginRsp     *SCLoginRsp `protobuf:"bytes,1,opt,name=m_SCLoginRsp" json:"m_SCLoginRsp,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *SCNetMsg) Reset()         { *m = SCNetMsg{} }
func (m *SCNetMsg) String() string { return proto.CompactTextString(m) }
func (*SCNetMsg) ProtoMessage()    {}

func (m *SCNetMsg) GetM_SCLoginRsp() *SCLoginRsp {
	if m != nil {
		return m.M_SCLoginRsp
	}
	return nil
}

type CSLoginReq struct {
	M_Name           *string `protobuf:"bytes,1,req,name=m_Name" json:"m_Name,omitempty"`
	M_Pwd            *string `protobuf:"bytes,2,req,name=m_Pwd" json:"m_Pwd,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CSLoginReq) Reset()         { *m = CSLoginReq{} }
func (m *CSLoginReq) String() string { return proto.CompactTextString(m) }
func (*CSLoginReq) ProtoMessage()    {}

func (m *CSLoginReq) GetM_Name() string {
	if m != nil && m.M_Name != nil {
		return *m.M_Name
	}
	return ""
}

func (m *CSLoginReq) GetM_Pwd() string {
	if m != nil && m.M_Pwd != nil {
		return *m.M_Pwd
	}
	return ""
}

type SCLoginRsp struct {
	M_Ret            *SCLoginRsp_LoginErrorType `protobuf:"varint,1,req,name=m_Ret,enum=netpb.SCLoginRsp_LoginErrorType" json:"m_Ret,omitempty"`
	M_Roles          []*SCLoginRsp_Result       `protobuf:"bytes,2,rep,name=m_Roles" json:"m_Roles,omitempty"`
	XXX_unrecognized []byte                     `json:"-"`
}

func (m *SCLoginRsp) Reset()         { *m = SCLoginRsp{} }
func (m *SCLoginRsp) String() string { return proto.CompactTextString(m) }
func (*SCLoginRsp) ProtoMessage()    {}

func (m *SCLoginRsp) GetM_Ret() SCLoginRsp_LoginErrorType {
	if m != nil && m.M_Ret != nil {
		return *m.M_Ret
	}
	return SCLoginRsp_SUCCESS
}

func (m *SCLoginRsp) GetM_Roles() []*SCLoginRsp_Result {
	if m != nil {
		return m.M_Roles
	}
	return nil
}

type SCLoginRsp_Result struct {
	M_Guid           *int64  `protobuf:"varint,1,req,name=m_Guid" json:"m_Guid,omitempty"`
	M_Name           *string `protobuf:"bytes,2,req,name=m_Name" json:"m_Name,omitempty"`
	M_Hp             *string `protobuf:"bytes,3,req,name=m_Hp" json:"m_Hp,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SCLoginRsp_Result) Reset()         { *m = SCLoginRsp_Result{} }
func (m *SCLoginRsp_Result) String() string { return proto.CompactTextString(m) }
func (*SCLoginRsp_Result) ProtoMessage()    {}

func (m *SCLoginRsp_Result) GetM_Guid() int64 {
	if m != nil && m.M_Guid != nil {
		return *m.M_Guid
	}
	return 0
}

func (m *SCLoginRsp_Result) GetM_Name() string {
	if m != nil && m.M_Name != nil {
		return *m.M_Name
	}
	return ""
}

func (m *SCLoginRsp_Result) GetM_Hp() string {
	if m != nil && m.M_Hp != nil {
		return *m.M_Hp
	}
	return ""
}

func init() {
	proto.RegisterEnum("netpb.MsgID", MsgID_name, MsgID_value)
	proto.RegisterEnum("netpb.SCLoginRsp_LoginErrorType", SCLoginRsp_LoginErrorType_name, SCLoginRsp_LoginErrorType_value)
}
