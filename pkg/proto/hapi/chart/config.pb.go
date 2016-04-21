// Code generated by protoc-gen-go.
// source: hapi/chart/config.proto
// DO NOT EDIT!

package chart

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

//
// Config:
//
// 		A config supplies values to the parametrizable templates of a chart.
//
type Config struct {
	Raw    string            `protobuf:"bytes,1,opt,name=raw" json:"raw,omitempty"`
	Values map[string]*Value `protobuf:"bytes,2,rep,name=values" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Config) GetValues() map[string]*Value {
	if m != nil {
		return m.Values
	}
	return nil
}

//
// Value:
//
// 		TODO
//
type Value struct {
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *Value) Reset()                    { *m = Value{} }
func (m *Value) String() string            { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()               {}
func (*Value) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func init() {
	proto.RegisterType((*Config)(nil), "hapi.chart.Config")
	proto.RegisterType((*Value)(nil), "hapi.chart.Value")
}

var fileDescriptor1 = []byte{
	// 179 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0xcf, 0x48, 0x2c, 0xc8,
	0xd4, 0x4f, 0xce, 0x48, 0x2c, 0x2a, 0xd1, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x02, 0x49, 0xe8, 0x81, 0x25, 0x94, 0x16, 0x30, 0x72, 0xb1, 0x39,
	0x83, 0x25, 0x85, 0x04, 0xb8, 0x98, 0x8b, 0x12, 0xcb, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83,
	0x40, 0x4c, 0x21, 0x33, 0x2e, 0xb6, 0xb2, 0xc4, 0x9c, 0xd2, 0xd4, 0x62, 0x09, 0x26, 0x05, 0x66,
	0x0d, 0x6e, 0x23, 0x39, 0x3d, 0x84, 0x4e, 0x3d, 0x88, 0x2e, 0xbd, 0x30, 0xb0, 0x02, 0xd7, 0xbc,
	0x92, 0xa2, 0xca, 0x20, 0xa8, 0x6a, 0x29, 0x1f, 0x2e, 0x6e, 0x24, 0x61, 0x90, 0xc1, 0xd9, 0xa9,
	0x95, 0x30, 0x83, 0x81, 0x4c, 0x21, 0x75, 0x2e, 0x56, 0xb0, 0x52, 0xa0, 0xb9, 0x8c, 0x40, 0x73,
	0x05, 0x91, 0xcd, 0x05, 0xeb, 0x0c, 0x82, 0xc8, 0x5b, 0x31, 0x59, 0x30, 0x2a, 0xc9, 0x72, 0xb1,
	0x82, 0xc5, 0x84, 0x44, 0x60, 0xba, 0x20, 0x26, 0x41, 0x38, 0x4e, 0xec, 0x51, 0xac, 0x60, 0x8d,
	0x49, 0x6c, 0x60, 0xdf, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xe1, 0x12, 0x60, 0xda, 0xf8,
	0x00, 0x00, 0x00,
}
