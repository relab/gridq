// Code generated by protoc-gen-gogo.
// source: gorums.proto
// DO NOT EDIT!

/*
Package gorums is a generated protocol buffer package.

It is generated from these files:
	gorums.proto

It has these top-level messages:
*/
package gorums

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

var E_Qc = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50000,
	Name:          "gorums.qc",
	Tag:           "varint,50000,opt,name=qc",
}

var E_Correctable = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50001,
	Name:          "gorums.correctable",
	Tag:           "varint,50001,opt,name=correctable",
}

var E_CorrectablePr = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50002,
	Name:          "gorums.correctable_pr",
	Tag:           "varint,50002,opt,name=correctable_pr,json=correctablePr",
}

var E_Multicast = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50003,
	Name:          "gorums.multicast",
	Tag:           "varint,50003,opt,name=multicast",
}

var E_QcFuture = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50004,
	Name:          "gorums.qc_future",
	Tag:           "varint,50004,opt,name=qc_future,json=qcFuture",
}

var E_QfWithReq = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50005,
	Name:          "gorums.qf_with_req",
	Tag:           "varint,50005,opt,name=qf_with_req,json=qfWithReq",
}

var E_PerNodeArg = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50006,
	Name:          "gorums.per_node_arg",
	Tag:           "varint,50006,opt,name=per_node_arg,json=perNodeArg",
}

var E_CustomReturnType = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         51000,
	Name:          "gorums.custom_return_type",
	Tag:           "bytes,51000,opt,name=custom_return_type,json=customReturnType",
}

func init() {
	proto.RegisterExtension(E_Qc)
	proto.RegisterExtension(E_Correctable)
	proto.RegisterExtension(E_CorrectablePr)
	proto.RegisterExtension(E_Multicast)
	proto.RegisterExtension(E_QcFuture)
	proto.RegisterExtension(E_QfWithReq)
	proto.RegisterExtension(E_PerNodeArg)
	proto.RegisterExtension(E_CustomReturnType)
}

func init() { proto.RegisterFile("gorums.proto", fileDescriptorGorums) }

var fileDescriptorGorums = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x3d, 0x6b, 0xdb, 0x40,
	0x1c, 0xc6, 0x75, 0x2e, 0x18, 0xfb, 0xec, 0x96, 0xe2, 0xa5, 0xa5, 0xc3, 0xe1, 0xb1, 0x93, 0x54,
	0xe8, 0x76, 0xd0, 0xd2, 0x1a, 0xda, 0x4e, 0x75, 0x8b, 0x08, 0x04, 0xb2, 0x1c, 0xd2, 0xe9, 0xf4,
	0x02, 0x92, 0xef, 0xee, 0xaf, 0x3b, 0x82, 0xb7, 0x4c, 0x19, 0x43, 0x3e, 0x46, 0x3e, 0x42, 0x3e,
	0x42, 0x46, 0xe7, 0x95, 0x8c, 0xb6, 0xb2, 0x64, 0xcc, 0x47, 0x08, 0x91, 0x64, 0xe2, 0x4d, 0xe3,
	0x73, 0x3c, 0xbf, 0xdf, 0xdd, 0x03, 0x87, 0xc7, 0x89, 0x04, 0x5b, 0x94, 0xae, 0x02, 0x69, 0xe4,
	0xa4, 0xdf, 0xa4, 0x4f, 0xd3, 0x44, 0xca, 0x24, 0x17, 0x5e, 0x7d, 0x1a, 0xda, 0xd8, 0x8b, 0x44,
	0xc9, 0x21, 0x53, 0x46, 0x42, 0xd3, 0xa4, 0x5f, 0x70, 0x4f, 0xf3, 0x09, 0x71, 0x9b, 0xa2, 0xbb,
	0x2d, 0xba, 0x7f, 0x85, 0x49, 0x65, 0xf4, 0x4f, 0x99, 0x4c, 0x2e, 0xca, 0x8f, 0xab, 0xe3, 0x37,
	0x53, 0xf4, 0x79, 0xe0, 0xf7, 0x34, 0xa7, 0x33, 0x3c, 0xe2, 0x12, 0x40, 0x70, 0x13, 0x84, 0xb9,
	0xe8, 0x44, 0x2f, 0x5b, 0x74, 0x17, 0xa2, 0x7f, 0xf0, 0xbb, 0x9d, 0xc8, 0x14, 0x74, 0x6a, 0xae,
	0x5a, 0xcd, 0xdb, 0x1d, 0xee, 0x3f, 0xd0, 0xef, 0x78, 0x58, 0xd8, 0xdc, 0x64, 0x3c, 0x28, 0x4d,
	0xa7, 0xe3, 0xba, 0x75, 0xbc, 0x22, 0xf4, 0x1b, 0x1e, 0x6a, 0xce, 0x62, 0x6b, 0x2c, 0x74, 0x4f,
	0xb9, 0x69, 0xf9, 0x81, 0xe6, 0xbf, 0x6b, 0x82, 0xfe, 0xc0, 0x23, 0x1d, 0xb3, 0xc3, 0xcc, 0xa4,
	0x0c, 0x84, 0xee, 0x14, 0xdc, 0x6e, 0x1f, 0xa0, 0xe3, 0xfd, 0xcc, 0xa4, 0xbe, 0xd0, 0x74, 0x86,
	0xc7, 0x4a, 0x00, 0x5b, 0xc8, 0x48, 0xb0, 0x00, 0x92, 0x4e, 0xc5, 0x5d, 0xab, 0xc0, 0x4a, 0xc0,
	0x5c, 0x46, 0xe2, 0x27, 0x24, 0x74, 0x8e, 0x27, 0xdc, 0x96, 0x46, 0x16, 0x0c, 0x84, 0xb1, 0xb0,
	0x60, 0x66, 0xa9, 0xba, 0xd7, 0x9c, 0x9f, 0xbc, 0x98, 0x86, 0xfe, 0xfb, 0x86, 0xf5, 0x6b, 0x74,
	0x6f, 0xa9, 0xc4, 0xec, 0xd7, 0x6a, 0x43, 0x9c, 0xfb, 0x0d, 0x71, 0xd6, 0x1b, 0x82, 0x8e, 0x2a,
	0x82, 0xce, 0x2a, 0x82, 0x2e, 0x2a, 0x82, 0x56, 0x15, 0x41, 0xeb, 0x8a, 0xa0, 0xc7, 0x8a, 0x38,
	0x4f, 0x15, 0x41, 0xa7, 0x0f, 0xc4, 0x39, 0xf8, 0x90, 0x64, 0x26, 0xb5, 0xa1, 0xcb, 0x65, 0xe1,
	0x81, 0xc8, 0x83, 0xd0, 0x6b, 0x3e, 0x5f, 0xd8, 0xaf, 0x2f, 0xfe, 0xfa, 0x1c, 0x00, 0x00, 0xff,
	0xff, 0x54, 0xe8, 0xa7, 0xd7, 0x9b, 0x02, 0x00, 0x00,
}