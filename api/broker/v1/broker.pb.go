// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.26.1
// source: broker/v1/broker.proto

package v1

import (
	v1 "github.com/knita-io/knita/api/executor/v1"
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

type RuntimeTender struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenderId string   `protobuf:"bytes,1,opt,name=tender_id,json=tenderId,proto3" json:"tender_id,omitempty"`
	Opts     *v1.Opts `protobuf:"bytes,2,opt,name=opts,proto3" json:"opts,omitempty"`
}

func (x *RuntimeTender) Reset() {
	*x = RuntimeTender{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_v1_broker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuntimeTender) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuntimeTender) ProtoMessage() {}

func (x *RuntimeTender) ProtoReflect() protoreflect.Message {
	mi := &file_broker_v1_broker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuntimeTender.ProtoReflect.Descriptor instead.
func (*RuntimeTender) Descriptor() ([]byte, []int) {
	return file_broker_v1_broker_proto_rawDescGZIP(), []int{0}
}

func (x *RuntimeTender) GetTenderId() string {
	if x != nil {
		return x.TenderId
	}
	return ""
}

func (x *RuntimeTender) GetOpts() *v1.Opts {
	if x != nil {
		return x.Opts
	}
	return nil
}

type RuntimeContracts struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Contracts []*RuntimeContract `protobuf:"bytes,1,rep,name=contracts,proto3" json:"contracts,omitempty"`
}

func (x *RuntimeContracts) Reset() {
	*x = RuntimeContracts{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_v1_broker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuntimeContracts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuntimeContracts) ProtoMessage() {}

func (x *RuntimeContracts) ProtoReflect() protoreflect.Message {
	mi := &file_broker_v1_broker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuntimeContracts.ProtoReflect.Descriptor instead.
func (*RuntimeContracts) Descriptor() ([]byte, []int) {
	return file_broker_v1_broker_proto_rawDescGZIP(), []int{1}
}

func (x *RuntimeContracts) GetContracts() []*RuntimeContract {
	if x != nil {
		return x.Contracts
	}
	return nil
}

type RuntimeContract struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenderId     string           `protobuf:"bytes,1,opt,name=tender_id,json=tenderId,proto3" json:"tender_id,omitempty"`
	ContractId   string           `protobuf:"bytes,2,opt,name=contract_id,json=contractId,proto3" json:"contract_id,omitempty"`
	RuntimeId    string           `protobuf:"bytes,3,opt,name=runtime_id,json=runtimeId,proto3" json:"runtime_id,omitempty"`
	Opts         *v1.Opts         `protobuf:"bytes,4,opt,name=opts,proto3" json:"opts,omitempty"`
	SysInfo      *v1.SystemInfo   `protobuf:"bytes,5,opt,name=sys_info,json=sysInfo,proto3" json:"sys_info,omitempty"`
	ExecutorInfo *v1.ExecutorInfo `protobuf:"bytes,6,opt,name=executor_info,json=executorInfo,proto3" json:"executor_info,omitempty"`
}

func (x *RuntimeContract) Reset() {
	*x = RuntimeContract{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_v1_broker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuntimeContract) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuntimeContract) ProtoMessage() {}

func (x *RuntimeContract) ProtoReflect() protoreflect.Message {
	mi := &file_broker_v1_broker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuntimeContract.ProtoReflect.Descriptor instead.
func (*RuntimeContract) Descriptor() ([]byte, []int) {
	return file_broker_v1_broker_proto_rawDescGZIP(), []int{2}
}

func (x *RuntimeContract) GetTenderId() string {
	if x != nil {
		return x.TenderId
	}
	return ""
}

func (x *RuntimeContract) GetContractId() string {
	if x != nil {
		return x.ContractId
	}
	return ""
}

func (x *RuntimeContract) GetRuntimeId() string {
	if x != nil {
		return x.RuntimeId
	}
	return ""
}

func (x *RuntimeContract) GetOpts() *v1.Opts {
	if x != nil {
		return x.Opts
	}
	return nil
}

func (x *RuntimeContract) GetSysInfo() *v1.SystemInfo {
	if x != nil {
		return x.SysInfo
	}
	return nil
}

func (x *RuntimeContract) GetExecutorInfo() *v1.ExecutorInfo {
	if x != nil {
		return x.ExecutorInfo
	}
	return nil
}

type RuntimeSettlement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnectionInfo *RuntimeConnectionInfo `protobuf:"bytes,1,opt,name=connection_info,json=connectionInfo,proto3" json:"connection_info,omitempty"`
}

func (x *RuntimeSettlement) Reset() {
	*x = RuntimeSettlement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_v1_broker_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuntimeSettlement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuntimeSettlement) ProtoMessage() {}

func (x *RuntimeSettlement) ProtoReflect() protoreflect.Message {
	mi := &file_broker_v1_broker_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuntimeSettlement.ProtoReflect.Descriptor instead.
func (*RuntimeSettlement) Descriptor() ([]byte, []int) {
	return file_broker_v1_broker_proto_rawDescGZIP(), []int{3}
}

func (x *RuntimeSettlement) GetConnectionInfo() *RuntimeConnectionInfo {
	if x != nil {
		return x.ConnectionInfo
	}
	return nil
}

type RuntimeConnectionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Transport:
	//
	//	*RuntimeConnectionInfo_Unix
	//	*RuntimeConnectionInfo_Tcp
	Transport isRuntimeConnectionInfo_Transport `protobuf_oneof:"transport"`
}

func (x *RuntimeConnectionInfo) Reset() {
	*x = RuntimeConnectionInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_v1_broker_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuntimeConnectionInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuntimeConnectionInfo) ProtoMessage() {}

func (x *RuntimeConnectionInfo) ProtoReflect() protoreflect.Message {
	mi := &file_broker_v1_broker_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuntimeConnectionInfo.ProtoReflect.Descriptor instead.
func (*RuntimeConnectionInfo) Descriptor() ([]byte, []int) {
	return file_broker_v1_broker_proto_rawDescGZIP(), []int{4}
}

func (m *RuntimeConnectionInfo) GetTransport() isRuntimeConnectionInfo_Transport {
	if m != nil {
		return m.Transport
	}
	return nil
}

func (x *RuntimeConnectionInfo) GetUnix() *RuntimeTransportUnix {
	if x, ok := x.GetTransport().(*RuntimeConnectionInfo_Unix); ok {
		return x.Unix
	}
	return nil
}

func (x *RuntimeConnectionInfo) GetTcp() *RuntimeTransportTCP {
	if x, ok := x.GetTransport().(*RuntimeConnectionInfo_Tcp); ok {
		return x.Tcp
	}
	return nil
}

type isRuntimeConnectionInfo_Transport interface {
	isRuntimeConnectionInfo_Transport()
}

type RuntimeConnectionInfo_Unix struct {
	Unix *RuntimeTransportUnix `protobuf:"bytes,1,opt,name=unix,proto3,oneof"`
}

type RuntimeConnectionInfo_Tcp struct {
	Tcp *RuntimeTransportTCP `protobuf:"bytes,2,opt,name=tcp,proto3,oneof"`
}

func (*RuntimeConnectionInfo_Unix) isRuntimeConnectionInfo_Transport() {}

func (*RuntimeConnectionInfo_Tcp) isRuntimeConnectionInfo_Transport() {}

type RuntimeTransportUnix struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SocketPath string `protobuf:"bytes,1,opt,name=socket_path,json=socketPath,proto3" json:"socket_path,omitempty"`
}

func (x *RuntimeTransportUnix) Reset() {
	*x = RuntimeTransportUnix{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_v1_broker_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuntimeTransportUnix) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuntimeTransportUnix) ProtoMessage() {}

func (x *RuntimeTransportUnix) ProtoReflect() protoreflect.Message {
	mi := &file_broker_v1_broker_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuntimeTransportUnix.ProtoReflect.Descriptor instead.
func (*RuntimeTransportUnix) Descriptor() ([]byte, []int) {
	return file_broker_v1_broker_proto_rawDescGZIP(), []int{5}
}

func (x *RuntimeTransportUnix) GetSocketPath() string {
	if x != nil {
		return x.SocketPath
	}
	return ""
}

type RuntimeTransportTCP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *RuntimeTransportTCP) Reset() {
	*x = RuntimeTransportTCP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broker_v1_broker_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuntimeTransportTCP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuntimeTransportTCP) ProtoMessage() {}

func (x *RuntimeTransportTCP) ProtoReflect() protoreflect.Message {
	mi := &file_broker_v1_broker_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuntimeTransportTCP.ProtoReflect.Descriptor instead.
func (*RuntimeTransportTCP) Descriptor() ([]byte, []int) {
	return file_broker_v1_broker_proto_rawDescGZIP(), []int{6}
}

func (x *RuntimeTransportTCP) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

var File_broker_v1_broker_proto protoreflect.FileDescriptor

var file_broker_v1_broker_proto_rawDesc = []byte{
	0x0a, 0x16, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x72, 0x6f, 0x6b,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72,
	0x1a, 0x1a, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x50, 0x0a, 0x0d,
	0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x54, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1b, 0x0a,
	0x09, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x04, 0x6f, 0x70,
	0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70, 0x74, 0x73, 0x52, 0x04, 0x6f, 0x70, 0x74, 0x73, 0x22, 0x49,
	0x0a, 0x10, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x73, 0x12, 0x35, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x52,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x09,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x22, 0x80, 0x02, 0x0a, 0x0f, 0x52, 0x75,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x04, 0x6f, 0x70,
	0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70, 0x74, 0x73, 0x52, 0x04, 0x6f, 0x70, 0x74, 0x73, 0x12, 0x2f,
	0x0a, 0x08, 0x73, 0x79, 0x73, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x73, 0x79, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x3b, 0x0a, 0x0d, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f,
	0x72, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0c,
	0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x5b, 0x0a, 0x11,
	0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x74, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x46, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x62, 0x72, 0x6f,
	0x6b, 0x65, 0x72, 0x2e, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x89, 0x01, 0x0a, 0x15, 0x52, 0x75,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x32, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x52, 0x75, 0x6e, 0x74, 0x69,
	0x6d, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x55, 0x6e, 0x69, 0x78, 0x48,
	0x00, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x78, 0x12, 0x2f, 0x0a, 0x03, 0x74, 0x63, 0x70, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x52, 0x75,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x43,
	0x50, 0x48, 0x00, 0x52, 0x03, 0x74, 0x63, 0x70, 0x42, 0x0b, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x37, 0x0a, 0x14, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x55, 0x6e, 0x69, 0x78, 0x12, 0x1f, 0x0a,
	0x0b, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x50, 0x61, 0x74, 0x68, 0x22, 0x2f,
	0x0a, 0x13, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x54, 0x43, 0x50, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x32,
	0x88, 0x01, 0x0a, 0x0d, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x72, 0x6f, 0x6b, 0x65,
	0x72, 0x12, 0x39, 0x0a, 0x06, 0x54, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x62, 0x72,
	0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x54, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x1a, 0x18, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x52, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x12, 0x3c, 0x0a, 0x06,
	0x53, 0x65, 0x74, 0x74, 0x6c, 0x65, 0x12, 0x17, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x1a,
	0x19, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x53, 0x65, 0x74, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6e, 0x69, 0x74, 0x61, 0x2d, 0x69,
	0x6f, 0x2f, 0x6b, 0x6e, 0x69, 0x74, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x72, 0x6f, 0x6b,
	0x65, 0x72, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_broker_v1_broker_proto_rawDescOnce sync.Once
	file_broker_v1_broker_proto_rawDescData = file_broker_v1_broker_proto_rawDesc
)

func file_broker_v1_broker_proto_rawDescGZIP() []byte {
	file_broker_v1_broker_proto_rawDescOnce.Do(func() {
		file_broker_v1_broker_proto_rawDescData = protoimpl.X.CompressGZIP(file_broker_v1_broker_proto_rawDescData)
	})
	return file_broker_v1_broker_proto_rawDescData
}

var file_broker_v1_broker_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_broker_v1_broker_proto_goTypes = []interface{}{
	(*RuntimeTender)(nil),         // 0: broker.RuntimeTender
	(*RuntimeContracts)(nil),      // 1: broker.RuntimeContracts
	(*RuntimeContract)(nil),       // 2: broker.RuntimeContract
	(*RuntimeSettlement)(nil),     // 3: broker.RuntimeSettlement
	(*RuntimeConnectionInfo)(nil), // 4: broker.RuntimeConnectionInfo
	(*RuntimeTransportUnix)(nil),  // 5: broker.RuntimeTransportUnix
	(*RuntimeTransportTCP)(nil),   // 6: broker.RuntimeTransportTCP
	(*v1.Opts)(nil),               // 7: executor.Opts
	(*v1.SystemInfo)(nil),         // 8: executor.SystemInfo
	(*v1.ExecutorInfo)(nil),       // 9: executor.ExecutorInfo
}
var file_broker_v1_broker_proto_depIdxs = []int32{
	7,  // 0: broker.RuntimeTender.opts:type_name -> executor.Opts
	2,  // 1: broker.RuntimeContracts.contracts:type_name -> broker.RuntimeContract
	7,  // 2: broker.RuntimeContract.opts:type_name -> executor.Opts
	8,  // 3: broker.RuntimeContract.sys_info:type_name -> executor.SystemInfo
	9,  // 4: broker.RuntimeContract.executor_info:type_name -> executor.ExecutorInfo
	4,  // 5: broker.RuntimeSettlement.connection_info:type_name -> broker.RuntimeConnectionInfo
	5,  // 6: broker.RuntimeConnectionInfo.unix:type_name -> broker.RuntimeTransportUnix
	6,  // 7: broker.RuntimeConnectionInfo.tcp:type_name -> broker.RuntimeTransportTCP
	0,  // 8: broker.RuntimeBroker.Tender:input_type -> broker.RuntimeTender
	2,  // 9: broker.RuntimeBroker.Settle:input_type -> broker.RuntimeContract
	1,  // 10: broker.RuntimeBroker.Tender:output_type -> broker.RuntimeContracts
	3,  // 11: broker.RuntimeBroker.Settle:output_type -> broker.RuntimeSettlement
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_broker_v1_broker_proto_init() }
func file_broker_v1_broker_proto_init() {
	if File_broker_v1_broker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_broker_v1_broker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuntimeTender); i {
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
		file_broker_v1_broker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuntimeContracts); i {
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
		file_broker_v1_broker_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuntimeContract); i {
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
		file_broker_v1_broker_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuntimeSettlement); i {
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
		file_broker_v1_broker_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuntimeConnectionInfo); i {
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
		file_broker_v1_broker_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuntimeTransportUnix); i {
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
		file_broker_v1_broker_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuntimeTransportTCP); i {
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
	file_broker_v1_broker_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*RuntimeConnectionInfo_Unix)(nil),
		(*RuntimeConnectionInfo_Tcp)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_broker_v1_broker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_broker_v1_broker_proto_goTypes,
		DependencyIndexes: file_broker_v1_broker_proto_depIdxs,
		MessageInfos:      file_broker_v1_broker_proto_msgTypes,
	}.Build()
	File_broker_v1_broker_proto = out.File
	file_broker_v1_broker_proto_rawDesc = nil
	file_broker_v1_broker_proto_goTypes = nil
	file_broker_v1_broker_proto_depIdxs = nil
}
