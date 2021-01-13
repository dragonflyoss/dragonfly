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

// check and make seeds

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: pkg/rpc/cdnsystem/cdnsystem.proto

package cdnsystem

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

type SeedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Url    string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Filter string `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
	// check content downloaded from remote url
	UrlMeta *base.UrlMeta `protobuf:"bytes,4,opt,name=url_meta,json=urlMeta,proto3" json:"url_meta,omitempty"`
}

func (x *SeedRequest) Reset() {
	*x = SeedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeedRequest) ProtoMessage() {}

func (x *SeedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeedRequest.ProtoReflect.Descriptor instead.
func (*SeedRequest) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescGZIP(), []int{0}
}

func (x *SeedRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *SeedRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *SeedRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

func (x *SeedRequest) GetUrlMeta() *base.UrlMeta {
	if x != nil {
		return x.UrlMeta
	}
	return nil
}

// keep piece meta and data separately
// check piece md5 and total content length, no more file md5 verification
// scheduler can also analyze cdn traffic from origin
type PieceSeed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State         *base.ResponseState `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	PieceStyle    base.PieceStyle     `protobuf:"varint,2,opt,name=piece_style,json=pieceStyle,proto3,enum=base.PieceStyle" json:"piece_style,omitempty"`
	PieceNum      int32               `protobuf:"varint,3,opt,name=piece_num,json=pieceNum,proto3" json:"piece_num,omitempty"`
	SeedAddr      string              `protobuf:"bytes,4,opt,name=seed_addr,json=seedAddr,proto3" json:"seed_addr,omitempty"`
	Done          bool                `protobuf:"varint,5,opt,name=done,proto3" json:"done,omitempty"`                                        // whether or not all seeds are downloaded
	ContentLength int64               `protobuf:"varint,6,opt,name=content_length,json=contentLength,proto3" json:"content_length,omitempty"` // content total length for the url
}

func (x *PieceSeed) Reset() {
	*x = PieceSeed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PieceSeed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PieceSeed) ProtoMessage() {}

func (x *PieceSeed) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PieceSeed.ProtoReflect.Descriptor instead.
func (*PieceSeed) Descriptor() ([]byte, []int) {
	return file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescGZIP(), []int{1}
}

func (x *PieceSeed) GetState() *base.ResponseState {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *PieceSeed) GetPieceStyle() base.PieceStyle {
	if x != nil {
		return x.PieceStyle
	}
	return base.PieceStyle_PLAIN_UNSPECIFIED
}

func (x *PieceSeed) GetPieceNum() int32 {
	if x != nil {
		return x.PieceNum
	}
	return 0
}

func (x *PieceSeed) GetSeedAddr() string {
	if x != nil {
		return x.SeedAddr
	}
	return ""
}

func (x *PieceSeed) GetDone() bool {
	if x != nil {
		return x.Done
	}
	return false
}

func (x *PieceSeed) GetContentLength() int64 {
	if x != nil {
		return x.ContentLength
	}
	return 0
}

var File_pkg_rpc_cdnsystem_cdnsystem_proto protoreflect.FileDescriptor

var file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x64, 0x6e, 0x73, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x2f, 0x63, 0x64, 0x6e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x64, 0x6e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x1a, 0x17,
	0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x0b, 0x53, 0x65, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x08, 0x75, 0x72, 0x6c,
	0x5f, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x55, 0x72, 0x6c, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x07, 0x75, 0x72, 0x6c, 0x4d,
	0x65, 0x74, 0x61, 0x22, 0xde, 0x01, 0x0a, 0x09, 0x50, 0x69, 0x65, 0x63, 0x65, 0x53, 0x65, 0x65,
	0x64, 0x12, 0x29, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x31, 0x0a, 0x0b,
	0x70, 0x69, 0x65, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x10, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x69, 0x65, 0x63, 0x65, 0x53, 0x74,
	0x79, 0x6c, 0x65, 0x52, 0x0a, 0x70, 0x69, 0x65, 0x63, 0x65, 0x53, 0x74, 0x79, 0x6c, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x70, 0x69, 0x65, 0x63, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x70, 0x69, 0x65, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x12, 0x1b, 0x0a, 0x09,
	0x73, 0x65, 0x65, 0x64, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x73, 0x65, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x6f, 0x6e,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4c, 0x65,
	0x6e, 0x67, 0x74, 0x68, 0x32, 0x87, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x65, 0x64, 0x65, 0x72, 0x12,
	0x3f, 0x0a, 0x0b, 0x4f, 0x62, 0x74, 0x61, 0x69, 0x6e, 0x53, 0x65, 0x65, 0x64, 0x73, 0x12, 0x16,
	0x2e, 0x63, 0x64, 0x6e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x53, 0x65, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x64, 0x6e, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x2e, 0x50, 0x69, 0x65, 0x63, 0x65, 0x53, 0x65, 0x65, 0x64, 0x22, 0x00, 0x30, 0x01,
	0x12, 0x3c, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x50, 0x69, 0x65, 0x63, 0x65, 0x54, 0x61, 0x73, 0x6b,
	0x73, 0x12, 0x16, 0x2e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x50, 0x69, 0x65, 0x63, 0x65, 0x54, 0x61,
	0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x50, 0x69, 0x65, 0x63, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x00, 0x42, 0x36,
	0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x72, 0x61,
	0x67, 0x6f, 0x6e, 0x66, 0x6c, 0x79, 0x6f, 0x73, 0x73, 0x2f, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e,
	0x66, 0x6c, 0x79, 0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x64, 0x6e,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescOnce sync.Once
	file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescData = file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDesc
)

func file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescGZIP() []byte {
	file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescOnce.Do(func() {
		file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescData)
	})
	return file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDescData
}

var file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_rpc_cdnsystem_cdnsystem_proto_goTypes = []interface{}{
	(*SeedRequest)(nil),           // 0: cdnsystem.SeedRequest
	(*PieceSeed)(nil),             // 1: cdnsystem.PieceSeed
	(*base.UrlMeta)(nil),          // 2: base.UrlMeta
	(*base.ResponseState)(nil),    // 3: base.ResponseState
	(base.PieceStyle)(0),          // 4: base.PieceStyle
	(*base.PieceTaskRequest)(nil), // 5: base.PieceTaskRequest
	(*base.PiecePacket)(nil),      // 6: base.PiecePacket
}
var file_pkg_rpc_cdnsystem_cdnsystem_proto_depIdxs = []int32{
	2, // 0: cdnsystem.SeedRequest.url_meta:type_name -> base.UrlMeta
	3, // 1: cdnsystem.PieceSeed.state:type_name -> base.ResponseState
	4, // 2: cdnsystem.PieceSeed.piece_style:type_name -> base.PieceStyle
	0, // 3: cdnsystem.Seeder.ObtainSeeds:input_type -> cdnsystem.SeedRequest
	5, // 4: cdnsystem.Seeder.GetPieceTasks:input_type -> base.PieceTaskRequest
	1, // 5: cdnsystem.Seeder.ObtainSeeds:output_type -> cdnsystem.PieceSeed
	6, // 6: cdnsystem.Seeder.GetPieceTasks:output_type -> base.PiecePacket
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pkg_rpc_cdnsystem_cdnsystem_proto_init() }
func file_pkg_rpc_cdnsystem_cdnsystem_proto_init() {
	if File_pkg_rpc_cdnsystem_cdnsystem_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SeedRequest); i {
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
		file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PieceSeed); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_rpc_cdnsystem_cdnsystem_proto_goTypes,
		DependencyIndexes: file_pkg_rpc_cdnsystem_cdnsystem_proto_depIdxs,
		MessageInfos:      file_pkg_rpc_cdnsystem_cdnsystem_proto_msgTypes,
	}.Build()
	File_pkg_rpc_cdnsystem_cdnsystem_proto = out.File
	file_pkg_rpc_cdnsystem_cdnsystem_proto_rawDesc = nil
	file_pkg_rpc_cdnsystem_cdnsystem_proto_goTypes = nil
	file_pkg_rpc_cdnsystem_cdnsystem_proto_depIdxs = nil
}
