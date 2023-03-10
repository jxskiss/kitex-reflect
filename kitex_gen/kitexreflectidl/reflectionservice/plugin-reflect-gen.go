// Code generated by thrift-gen-kitex-reflect. DO NOT EDIT.
package reflectionservice

import (
	idl "github.com/jxskiss/kitex-reflect/kitex_gen/kitexreflectidl"
)

func init() {
	SetupReflectPlugin()
}

func SetupReflectPlugin() {
	serviceInfo().Methods["ReflectService"] = PluginReflect.NewMethodInfo("ReflectService")
}

func GetReflectServiceRespPayload() *idl.ReflectServiceRespPayload {
	return PluginReflect.NewReflectServiceRespPayload()
}

func GetIDLRawBytes() []byte {
	return idlRawBytes
}

var PluginReflect = &idl.PluginImpl{
	Version:          "20230120022144Z",
	IsCombineService: false,

	GetIDLBytes: func() []byte { return idlRawBytes },
}

var idlRawBytes = []byte{
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xbc, 0x94, 0xcb, 0x4e, 0xc3, 0x30,
	0x10, 0x45, 0x6f, 0x1e, 0x4d, 0x42, 0xdd, 0x12, 0x40, 0x02, 0xc1, 0x4f, 0xf0, 0x13, 0x94, 0x45,
	0xa5, 0x2e, 0x50, 0x91, 0xd8, 0xf7, 0xe1, 0x96, 0x11, 0xc1, 0x06, 0xdb, 0x45, 0xed, 0xdf, 0xa3,
	0x74, 0x5c, 0xc8, 0xab, 0x1b, 0x04, 0xa8, 0x8b, 0x2b, 0x37, 0xf6, 0xc9, 0x99, 0x19, 0x39, 0x02,
	0x01, 0x80, 0x3e, 0x2d, 0x8b, 0x5b, 0xf7, 0x6c, 0x68, 0xe5, 0x72, 0x84, 0x03, 0x00, 0xc8, 0x11,
	0x09, 0xce, 0xb8, 0x5c, 0x07, 0xbc, 0x31, 0x5c, 0x6b, 0x81, 0xb0, 0xfc, 0xfb, 0x85, 0x9c, 0xdc,
	0x1a, 0xb9, 0x2a, 0xe4, 0xc2, 0xd1, 0xb2, 0xc8, 0x11, 0xed, 0x8f, 0x21, 0x47, 0xcf, 0x9f, 0x4f,
	0x7c, 0xa6, 0x3e, 0xb3, 0x32, 0x63, 0xe6, 0x24, 0xd6, 0x99, 0xcd, 0xc2, 0x31, 0xeb, 0x7a, 0xca,
	0x98, 0x47, 0x69, 0x3e, 0x68, 0x21, 0xa7, 0xf2, 0xfd, 0x61, 0xb6, 0x2b, 0xf4, 0x6c, 0xe9, 0xa1,
	0x41, 0xb6, 0x3f, 0x13, 0xf0, 0xee, 0x8b, 0xfb, 0x2d, 0x59, 0x47, 0x6a, 0x3d, 0x1e, 0x4d, 0x9e,
	0xa4, 0xb1, 0xa4, 0x55, 0x86, 0xa8, 0x7c, 0xc5, 0xa0, 0x4a, 0x27, 0xb5, 0x16, 0x88, 0xab, 0x46,
	0x19, 0x12, 0x00, 0xe9, 0x97, 0x99, 0x40, 0x8a, 0x4a, 0x85, 0x10, 0xe8, 0x81, 0xb3, 0x65, 0x78,
	0xd3, 0x34, 0xb4, 0x6f, 0x75, 0xc5, 0xa8, 0xa6, 0x98, 0xfe, 0x8a, 0x17, 0x23, 0x43, 0x46, 0x9e,
	0x8d, 0xed, 0x9d, 0x7e, 0x9d, 0x93, 0x92, 0xde, 0xa1, 0xc9, 0x8e, 0xe7, 0x5a, 0x17, 0x9d, 0xe4,
	0xe0, 0x08, 0x39, 0x67, 0x72, 0x34, 0x1e, 0x4d, 0x5a, 0xa2, 0x73, 0x52, 0x33, 0xb3, 0xeb, 0xc4,
	0x65, 0x3f, 0x69, 0xe0, 0x65, 0x6b, 0xc4, 0x1b, 0x69, 0x5d, 0xe7, 0x7c, 0x53, 0xdf, 0xd9, 0x3f,
	0x77, 0xba, 0x6a, 0x0f, 0x55, 0x2b, 0x2b, 0xff, 0x51, 0x2a, 0xc7, 0x89, 0xbf, 0x1b, 0x7d, 0x9f,
	0xa2, 0x72, 0xd7, 0xce, 0xbd, 0x1f, 0x69, 0xe5, 0x15, 0xd9, 0x05, 0x07, 0x43, 0xde, 0x76, 0x5a,
	0x2f, 0x23, 0x44, 0xf9, 0x8b, 0xbe, 0x15, 0x8f, 0x94, 0xd9, 0xa9, 0x3c, 0x3c, 0xac, 0x1b, 0xd5,
	0x1b, 0x9e, 0x57, 0xb3, 0xfa, 0xee, 0xa9, 0x1e, 0x25, 0x37, 0x9b, 0x51, 0xff, 0x40, 0x08, 0x64,
	0x8d, 0x26, 0x25, 0xfb, 0xf5, 0x10, 0x03, 0xc1, 0x4f, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc9,
	0xb4, 0x52, 0xad, 0xac, 0x04, 0x00, 0x00,
}
