//
//     Copyright 2020 The Dragonfly Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: pkg/rpc/manager/manager.proto

package manager

import (
	base "github.com/dragonflyoss/Dragonfly2/pkg/rpc/base"
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type NavigatorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// client ip
	Ip string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	// client host name
	HostName string `protobuf:"bytes,2,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	// json format: {vpcId:xxx,sn:xxx,group:xxx,...}
	HostTag string `protobuf:"bytes,3,opt,name=host_tag,json=hostTag,proto3" json:"host_tag,omitempty"`
}

func (x *NavigatorRequest) Reset() {
	*x = NavigatorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NavigatorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NavigatorRequest) ProtoMessage() {}

func (x *NavigatorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NavigatorRequest.ProtoReflect.Descriptor instead.
func (*NavigatorRequest) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{0}
}

func (x *NavigatorRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *NavigatorRequest) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *NavigatorRequest) GetHostTag() string {
	if x != nil {
		return x.HostTag
	}
	return ""
}

type SchedulerNodes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State *base.ResponseState `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	// ip:port
	Addrs      []string  `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	ClientHost *HostInfo `protobuf:"bytes,3,opt,name=client_host,json=clientHost,proto3" json:"client_host,omitempty"`
}

func (x *SchedulerNodes) Reset() {
	*x = SchedulerNodes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchedulerNodes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchedulerNodes) ProtoMessage() {}

func (x *SchedulerNodes) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchedulerNodes.ProtoReflect.Descriptor instead.
func (*SchedulerNodes) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{1}
}

func (x *SchedulerNodes) GetState() *base.ResponseState {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *SchedulerNodes) GetAddrs() []string {
	if x != nil {
		return x.Addrs
	}
	return nil
}

func (x *SchedulerNodes) GetClientHost() *HostInfo {
	if x != nil {
		return x.ClientHost
	}
	return nil
}

type HeartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// identify servers with hostname
	HostName string `protobuf:"bytes,1,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	// Types that are assignable to From:
	//	*HeartRequest_Scheduler
	//	*HeartRequest_Cdn
	From isHeartRequest_From `protobuf_oneof:"from"`
}

func (x *HeartRequest) Reset() {
	*x = HeartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartRequest) ProtoMessage() {}

func (x *HeartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartRequest.ProtoReflect.Descriptor instead.
func (*HeartRequest) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{2}
}

func (x *HeartRequest) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (m *HeartRequest) GetFrom() isHeartRequest_From {
	if m != nil {
		return m.From
	}
	return nil
}

func (x *HeartRequest) GetScheduler() bool {
	if x, ok := x.GetFrom().(*HeartRequest_Scheduler); ok {
		return x.Scheduler
	}
	return false
}

func (x *HeartRequest) GetCdn() bool {
	if x, ok := x.GetFrom().(*HeartRequest_Cdn); ok {
		return x.Cdn
	}
	return false
}

type isHeartRequest_From interface {
	isHeartRequest_From()
}

type HeartRequest_Scheduler struct {
	Scheduler bool `protobuf:"varint,2,opt,name=scheduler,proto3,oneof"`
}

type HeartRequest_Cdn struct {
	Cdn bool `protobuf:"varint,3,opt,name=cdn,proto3,oneof"`
}

func (*HeartRequest_Scheduler) isHeartRequest_From() {}

func (*HeartRequest_Cdn) isHeartRequest_From() {}

type ClientConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClientConfig) Reset() {
	*x = ClientConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientConfig) ProtoMessage() {}

func (x *ClientConfig) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientConfig.ProtoReflect.Descriptor instead.
func (*ClientConfig) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{3}
}

type ManagementConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State *base.ResponseState `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	// Types that are assignable to Config:
	//	*ManagementConfig_SchedulerConfig_
	//	*ManagementConfig_CdnConfig_
	Config isManagementConfig_Config `protobuf_oneof:"config"`
}

func (x *ManagementConfig) Reset() {
	*x = ManagementConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManagementConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManagementConfig) ProtoMessage() {}

func (x *ManagementConfig) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManagementConfig.ProtoReflect.Descriptor instead.
func (*ManagementConfig) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{4}
}

func (x *ManagementConfig) GetState() *base.ResponseState {
	if x != nil {
		return x.State
	}
	return nil
}

func (m *ManagementConfig) GetConfig() isManagementConfig_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (x *ManagementConfig) GetSchedulerConfig() *ManagementConfig_SchedulerConfig {
	if x, ok := x.GetConfig().(*ManagementConfig_SchedulerConfig_); ok {
		return x.SchedulerConfig
	}
	return nil
}

func (x *ManagementConfig) GetCdnConfig() *ManagementConfig_CdnConfig {
	if x, ok := x.GetConfig().(*ManagementConfig_CdnConfig_); ok {
		return x.CdnConfig
	}
	return nil
}

type isManagementConfig_Config interface {
	isManagementConfig_Config()
}

type ManagementConfig_SchedulerConfig_ struct {
	SchedulerConfig *ManagementConfig_SchedulerConfig `protobuf:"bytes,2,opt,name=scheduler_config,json=schedulerConfig,proto3,oneof"`
}

type ManagementConfig_CdnConfig_ struct {
	CdnConfig *ManagementConfig_CdnConfig `protobuf:"bytes,3,opt,name=cdn_config,json=cdnConfig,proto3,oneof"`
}

func (*ManagementConfig_SchedulerConfig_) isManagementConfig_Config() {}

func (*ManagementConfig_CdnConfig_) isManagementConfig_Config() {}

type ServerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HostInfo *HostInfo `protobuf:"bytes,1,opt,name=host_info,json=hostInfo,proto3" json:"host_info,omitempty"`
	RpcPort  int32     `protobuf:"varint,2,opt,name=rpc_port,json=rpcPort,proto3" json:"rpc_port,omitempty"`
	DownPort int32     `protobuf:"varint,3,opt,name=down_port,json=downPort,proto3" json:"down_port,omitempty"`
}

func (x *ServerInfo) Reset() {
	*x = ServerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerInfo) ProtoMessage() {}

func (x *ServerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerInfo.ProtoReflect.Descriptor instead.
func (*ServerInfo) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{5}
}

func (x *ServerInfo) GetHostInfo() *HostInfo {
	if x != nil {
		return x.HostInfo
	}
	return nil
}

func (x *ServerInfo) GetRpcPort() int32 {
	if x != nil {
		return x.RpcPort
	}
	return 0
}

func (x *ServerInfo) GetDownPort() int32 {
	if x != nil {
		return x.DownPort
	}
	return 0
}

type HostInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip       string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	HostName string `protobuf:"bytes,2,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	// security isolation domain for network
	SecurityDomain string `protobuf:"bytes,3,opt,name=security_domain,json=securityDomain,proto3" json:"security_domain,omitempty"`
	// area|country|province|city|...
	Location    string `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	Idc         string `protobuf:"bytes,5,opt,name=idc,proto3" json:"idc,omitempty"`
	NetTopology string `protobuf:"bytes,6,opt,name=net_topology,json=netTopology,proto3" json:"net_topology,omitempty"`
}

func (x *HostInfo) Reset() {
	*x = HostInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HostInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostInfo) ProtoMessage() {}

func (x *HostInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostInfo.ProtoReflect.Descriptor instead.
func (*HostInfo) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{6}
}

func (x *HostInfo) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *HostInfo) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *HostInfo) GetSecurityDomain() string {
	if x != nil {
		return x.SecurityDomain
	}
	return ""
}

func (x *HostInfo) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *HostInfo) GetIdc() string {
	if x != nil {
		return x.Idc
	}
	return ""
}

func (x *HostInfo) GetNetTopology() string {
	if x != nil {
		return x.NetTopology
	}
	return ""
}

type ManagementConfig_SchedulerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientConfig *ClientConfig `protobuf:"bytes,1,opt,name=client_config,json=clientConfig,proto3" json:"client_config,omitempty"`
	CdnHosts     []*ServerInfo `protobuf:"bytes,2,rep,name=cdn_hosts,json=cdnHosts,proto3" json:"cdn_hosts,omitempty"` //......
}

func (x *ManagementConfig_SchedulerConfig) Reset() {
	*x = ManagementConfig_SchedulerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManagementConfig_SchedulerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManagementConfig_SchedulerConfig) ProtoMessage() {}

func (x *ManagementConfig_SchedulerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManagementConfig_SchedulerConfig.ProtoReflect.Descriptor instead.
func (*ManagementConfig_SchedulerConfig) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{4, 0}
}

func (x *ManagementConfig_SchedulerConfig) GetClientConfig() *ClientConfig {
	if x != nil {
		return x.ClientConfig
	}
	return nil
}

func (x *ManagementConfig_SchedulerConfig) GetCdnHosts() []*ServerInfo {
	if x != nil {
		return x.CdnHosts
	}
	return nil
}

type ManagementConfig_CdnConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ManagementConfig_CdnConfig) Reset() {
	*x = ManagementConfig_CdnConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_manager_manager_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManagementConfig_CdnConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManagementConfig_CdnConfig) ProtoMessage() {}

func (x *ManagementConfig_CdnConfig) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_manager_manager_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManagementConfig_CdnConfig.ProtoReflect.Descriptor instead.
func (*ManagementConfig_CdnConfig) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_manager_manager_proto_rawDescGZIP(), []int{4, 1}
}

var File_pkg_rpc_manager_manager_proto protoreflect.FileDescriptor

var file_pkg_rpc_manager_manager_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x1a, 0x17, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70,
	0x63, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x5a, 0x0a, 0x10, 0x4e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x74, 0x61, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x68, 0x6f, 0x73, 0x74, 0x54, 0x61, 0x67, 0x22, 0x85, 0x01,
	0x0a, 0x0e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x73,
	0x12, 0x29, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61,
	0x64, 0x64, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x61, 0x64, 0x64, 0x72,
	0x73, 0x12, 0x32, 0x0a, 0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x6f, 0x73, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x48, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x48, 0x6f, 0x73, 0x74, 0x22, 0x67, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x72, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x72, 0x12, 0x12, 0x0a, 0x03, 0x63, 0x64, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48,
	0x00, 0x52, 0x03, 0x63, 0x64, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x22, 0x0e,
	0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xf3,
	0x02, 0x0a, 0x10, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x29, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x56,
	0x0a, 0x10, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2e, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x44, 0x0a, 0x0a, 0x63, 0x64, 0x6e, 0x5f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x64, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x48,
	0x00, 0x52, 0x09, 0x63, 0x64, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x7f, 0x0a, 0x0f,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x3a, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0c, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x30, 0x0a, 0x09, 0x63,
	0x64, 0x6e, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x08, 0x63, 0x64, 0x6e, 0x48, 0x6f, 0x73, 0x74, 0x73, 0x1a, 0x0b, 0x0a,
	0x09, 0x43, 0x64, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x08, 0x0a, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x22, 0x74, 0x0a, 0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x2e, 0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e,
	0x48, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x70, 0x63, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x72, 0x70, 0x63, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x64, 0x6f, 0x77, 0x6e, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x64, 0x6f, 0x77, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x22, 0xb1, 0x01, 0x0a, 0x08, 0x48,
	0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79,
	0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73,
	0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x63,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x63, 0x12, 0x21, 0x0a, 0x0c, 0x6e,
	0x65, 0x74, 0x5f, 0x74, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6e, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x6f, 0x6c, 0x6f, 0x67, 0x79, 0x32, 0x91,
	0x01, 0x0a, 0x07, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x45, 0x0a, 0x0d, 0x47, 0x65,
	0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x73, 0x12, 0x19, 0x2e, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x4e, 0x61, 0x76, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x22,
	0x00, 0x12, 0x3f, 0x0a, 0x09, 0x4b, 0x65, 0x65, 0x70, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x15,
	0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e,
	0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x22, 0x00, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x64, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x66, 0x6c, 0x79, 0x6f, 0x73, 0x73, 0x2f, 0x44, 0x72,
	0x61, 0x67, 0x6f, 0x6e, 0x66, 0x6c, 0x79, 0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63,
	0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_rpc_manager_manager_proto_rawDescOnce sync.Once
	file_pkg_rpc_manager_manager_proto_rawDescData = file_pkg_rpc_manager_manager_proto_rawDesc
)

func file_pkg_rpc_manager_manager_proto_rawDescGZIP() []byte {
	file_pkg_rpc_manager_manager_proto_rawDescOnce.Do(func() {
		file_pkg_rpc_manager_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_rpc_manager_manager_proto_rawDescData)
	})
	return file_pkg_rpc_manager_manager_proto_rawDescData
}

var file_pkg_rpc_manager_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_pkg_rpc_manager_manager_proto_goTypes = []interface{}{
	(*NavigatorRequest)(nil),                 // 0: manager.NavigatorRequest
	(*SchedulerNodes)(nil),                   // 1: manager.SchedulerNodes
	(*HeartRequest)(nil),                     // 2: manager.HeartRequest
	(*ClientConfig)(nil),                     // 3: manager.ClientConfig
	(*ManagementConfig)(nil),                 // 4: manager.ManagementConfig
	(*ServerInfo)(nil),                       // 5: manager.ServerInfo
	(*HostInfo)(nil),                         // 6: manager.HostInfo
	(*ManagementConfig_SchedulerConfig)(nil), // 7: manager.ManagementConfig.SchedulerConfig
	(*ManagementConfig_CdnConfig)(nil),       // 8: manager.ManagementConfig.CdnConfig
	(*base.ResponseState)(nil),               // 9: base.ResponseState
}
var file_pkg_rpc_manager_manager_proto_depIdxs = []int32{
	9,  // 0: manager.SchedulerNodes.state:type_name -> base.ResponseState
	6,  // 1: manager.SchedulerNodes.client_host:type_name -> manager.HostInfo
	9,  // 2: manager.ManagementConfig.state:type_name -> base.ResponseState
	7,  // 3: manager.ManagementConfig.scheduler_config:type_name -> manager.ManagementConfig.SchedulerConfig
	8,  // 4: manager.ManagementConfig.cdn_config:type_name -> manager.ManagementConfig.CdnConfig
	6,  // 5: manager.ServerInfo.host_info:type_name -> manager.HostInfo
	3,  // 6: manager.ManagementConfig.SchedulerConfig.client_config:type_name -> manager.ClientConfig
	5,  // 7: manager.ManagementConfig.SchedulerConfig.cdn_hosts:type_name -> manager.ServerInfo
	0,  // 8: manager.Manager.GetSchedulers:input_type -> manager.NavigatorRequest
	2,  // 9: manager.Manager.KeepAlive:input_type -> manager.HeartRequest
	1,  // 10: manager.Manager.GetSchedulers:output_type -> manager.SchedulerNodes
	4,  // 11: manager.Manager.KeepAlive:output_type -> manager.ManagementConfig
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_pkg_rpc_manager_manager_proto_init() }
func file_pkg_rpc_manager_manager_proto_init() {
	if File_pkg_rpc_manager_manager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_rpc_manager_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NavigatorRequest); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchedulerNodes); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartRequest); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientConfig); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManagementConfig); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerInfo); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HostInfo); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManagementConfig_SchedulerConfig); i {
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
		file_pkg_rpc_manager_manager_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManagementConfig_CdnConfig); i {
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
	file_pkg_rpc_manager_manager_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*HeartRequest_Scheduler)(nil),
		(*HeartRequest_Cdn)(nil),
	}
	file_pkg_rpc_manager_manager_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*ManagementConfig_SchedulerConfig_)(nil),
		(*ManagementConfig_CdnConfig_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_rpc_manager_manager_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_rpc_manager_manager_proto_goTypes,
		DependencyIndexes: file_pkg_rpc_manager_manager_proto_depIdxs,
		MessageInfos:      file_pkg_rpc_manager_manager_proto_msgTypes,
	}.Build()
	File_pkg_rpc_manager_manager_proto = out.File
	file_pkg_rpc_manager_manager_proto_rawDesc = nil
	file_pkg_rpc_manager_manager_proto_goTypes = nil
	file_pkg_rpc_manager_manager_proto_depIdxs = nil
}
