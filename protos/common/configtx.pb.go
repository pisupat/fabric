// Code generated by protoc-gen-go.
// source: common/configtx.proto
// DO NOT EDIT!

package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ConfigItem_ConfigType int32

const (
	ConfigItem_POLICY  ConfigItem_ConfigType = 0
	ConfigItem_CHAIN   ConfigItem_ConfigType = 1
	ConfigItem_ORDERER ConfigItem_ConfigType = 2
	ConfigItem_PEER    ConfigItem_ConfigType = 3
	ConfigItem_MSP     ConfigItem_ConfigType = 4
)

var ConfigItem_ConfigType_name = map[int32]string{
	0: "POLICY",
	1: "CHAIN",
	2: "ORDERER",
	3: "PEER",
	4: "MSP",
}
var ConfigItem_ConfigType_value = map[string]int32{
	"POLICY":  0,
	"CHAIN":   1,
	"ORDERER": 2,
	"PEER":    3,
	"MSP":     4,
}

func (x ConfigItem_ConfigType) String() string {
	return proto.EnumName(ConfigItem_ConfigType_name, int32(x))
}
func (ConfigItem_ConfigType) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{7, 0} }

// ConfigEnvelope is designed to contain _all_ configuration for a chain with no dependency
// on previous configuration transactions.
//
// It is generated with the following scheme:
//   1. Retrieve the existing configuration
//   2. Note the highest configuration sequence number, store it and increment it by one
//   3. Modify desired ConfigItems, setting each LastModified to the stored and incremented sequence number
//     a) Note that the ConfigItem has a ChannelHeader header attached to it, who's type is set to CONFIGURATION_ITEM
//   4. Create Config message containing the new configuration, marshal it into ConfigEnvelope.config and encode the required signatures
//     a) Each signature is of type ConfigSignature
//     b) The ConfigSignature signature is over the concatenation of signatureHeader and the Config bytes (which includes a ChannelHeader)
//   5. Submit new Config for ordering in Envelope signed by submitter
//     a) The Envelope Payload has data set to the marshaled ConfigEnvelope
//     b) The Envelope Payload has a header of type Header.Type.CONFIGURATION_TRANSACTION
//
// The configuration manager will verify:
//   1. All configuration items and the envelope refer to the correct chain
//   2. Some configuration item has been added or modified
//   3. No existing configuration item has been ommitted
//   4. All configuration changes have a LastModification of one more than the last configuration's highest LastModification number
//   5. All configuration changes satisfy the corresponding modification policy
type ConfigEnvelope struct {
	Config     []byte             `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	Signatures []*ConfigSignature `protobuf:"bytes,2,rep,name=signatures" json:"signatures,omitempty"`
}

func (m *ConfigEnvelope) Reset()                    { *m = ConfigEnvelope{} }
func (m *ConfigEnvelope) String() string            { return proto.CompactTextString(m) }
func (*ConfigEnvelope) ProtoMessage()               {}
func (*ConfigEnvelope) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *ConfigEnvelope) GetSignatures() []*ConfigSignature {
	if m != nil {
		return m.Signatures
	}
	return nil
}

// ConfigTemplate is used as a serialization format to share configuration templates
// The orderer supplies a configuration template to the user to use when constructing a new
// chain creation transaction, so this is used to facilitate that.
type ConfigTemplate struct {
	Items []*ConfigItem `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
}

func (m *ConfigTemplate) Reset()                    { *m = ConfigTemplate{} }
func (m *ConfigTemplate) String() string            { return proto.CompactTextString(m) }
func (*ConfigTemplate) ProtoMessage()               {}
func (*ConfigTemplate) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *ConfigTemplate) GetItems() []*ConfigItem {
	if m != nil {
		return m.Items
	}
	return nil
}

// This message may change slightly depending on the finalization of signature schemes for transactions
type Config struct {
	Header *ChannelHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Items  []*ConfigItem  `protobuf:"bytes,2,rep,name=items" json:"items,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *Config) GetHeader() *ChannelHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Config) GetItems() []*ConfigItem {
	if m != nil {
		return m.Items
	}
	return nil
}

// XXX this structure is to allow us to minimize the diffs in this change series
// it will be renamed Config once the original is ready to be removed
type ConfigNext struct {
	Header  *ChannelHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Channel *ConfigGroup   `protobuf:"bytes,2,opt,name=channel" json:"channel,omitempty"`
}

func (m *ConfigNext) Reset()                    { *m = ConfigNext{} }
func (m *ConfigNext) String() string            { return proto.CompactTextString(m) }
func (*ConfigNext) ProtoMessage()               {}
func (*ConfigNext) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *ConfigNext) GetHeader() *ChannelHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ConfigNext) GetChannel() *ConfigGroup {
	if m != nil {
		return m.Channel
	}
	return nil
}

// ConfigGroup is the hierarchical data structure for holding config
type ConfigGroup struct {
	Version   uint64                   `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Groups    map[string]*ConfigGroup  `protobuf:"bytes,2,rep,name=groups" json:"groups,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Values    map[string]*ConfigValue  `protobuf:"bytes,3,rep,name=values" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Policies  map[string]*ConfigPolicy `protobuf:"bytes,4,rep,name=policies" json:"policies,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ModPolicy string                   `protobuf:"bytes,5,opt,name=mod_policy,json=modPolicy" json:"mod_policy,omitempty"`
}

func (m *ConfigGroup) Reset()                    { *m = ConfigGroup{} }
func (m *ConfigGroup) String() string            { return proto.CompactTextString(m) }
func (*ConfigGroup) ProtoMessage()               {}
func (*ConfigGroup) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *ConfigGroup) GetGroups() map[string]*ConfigGroup {
	if m != nil {
		return m.Groups
	}
	return nil
}

func (m *ConfigGroup) GetValues() map[string]*ConfigValue {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *ConfigGroup) GetPolicies() map[string]*ConfigPolicy {
	if m != nil {
		return m.Policies
	}
	return nil
}

// ConfigValue represents an individual piece of config data
type ConfigValue struct {
	Version   uint64 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Value     []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	ModPolicy string `protobuf:"bytes,3,opt,name=mod_policy,json=modPolicy" json:"mod_policy,omitempty"`
}

func (m *ConfigValue) Reset()                    { *m = ConfigValue{} }
func (m *ConfigValue) String() string            { return proto.CompactTextString(m) }
func (*ConfigValue) ProtoMessage()               {}
func (*ConfigValue) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

type ConfigPolicy struct {
	Version   uint64  `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Policy    *Policy `protobuf:"bytes,2,opt,name=policy" json:"policy,omitempty"`
	ModPolicy string  `protobuf:"bytes,3,opt,name=mod_policy,json=modPolicy" json:"mod_policy,omitempty"`
}

func (m *ConfigPolicy) Reset()                    { *m = ConfigPolicy{} }
func (m *ConfigPolicy) String() string            { return proto.CompactTextString(m) }
func (*ConfigPolicy) ProtoMessage()               {}
func (*ConfigPolicy) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *ConfigPolicy) GetPolicy() *Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

type ConfigItem struct {
	Type               ConfigItem_ConfigType `protobuf:"varint,1,opt,name=type,enum=common.ConfigItem_ConfigType" json:"type,omitempty"`
	LastModified       uint64                `protobuf:"varint,2,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
	ModificationPolicy string                `protobuf:"bytes,3,opt,name=modification_policy,json=modificationPolicy" json:"modification_policy,omitempty"`
	Key                string                `protobuf:"bytes,4,opt,name=key" json:"key,omitempty"`
	Value              []byte                `protobuf:"bytes,5,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *ConfigItem) Reset()                    { *m = ConfigItem{} }
func (m *ConfigItem) String() string            { return proto.CompactTextString(m) }
func (*ConfigItem) ProtoMessage()               {}
func (*ConfigItem) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

type ConfigSignature struct {
	SignatureHeader []byte `protobuf:"bytes,1,opt,name=signature_header,json=signatureHeader,proto3" json:"signature_header,omitempty"`
	Signature       []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *ConfigSignature) Reset()                    { *m = ConfigSignature{} }
func (m *ConfigSignature) String() string            { return proto.CompactTextString(m) }
func (*ConfigSignature) ProtoMessage()               {}
func (*ConfigSignature) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func init() {
	proto.RegisterType((*ConfigEnvelope)(nil), "common.ConfigEnvelope")
	proto.RegisterType((*ConfigTemplate)(nil), "common.ConfigTemplate")
	proto.RegisterType((*Config)(nil), "common.Config")
	proto.RegisterType((*ConfigNext)(nil), "common.ConfigNext")
	proto.RegisterType((*ConfigGroup)(nil), "common.ConfigGroup")
	proto.RegisterType((*ConfigValue)(nil), "common.ConfigValue")
	proto.RegisterType((*ConfigPolicy)(nil), "common.ConfigPolicy")
	proto.RegisterType((*ConfigItem)(nil), "common.ConfigItem")
	proto.RegisterType((*ConfigSignature)(nil), "common.ConfigSignature")
	proto.RegisterEnum("common.ConfigItem_ConfigType", ConfigItem_ConfigType_name, ConfigItem_ConfigType_value)
}

func init() { proto.RegisterFile("common/configtx.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 641 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x54, 0xdf, 0x4f, 0xd4, 0x40,
	0x10, 0xb6, 0x3f, 0xae, 0xc7, 0xcd, 0x1d, 0xd0, 0x2c, 0xa0, 0x0d, 0x91, 0x88, 0x35, 0x31, 0x87,
	0x06, 0x2e, 0xe2, 0x03, 0x86, 0xc4, 0x07, 0x3d, 0x1b, 0x21, 0x11, 0x38, 0x17, 0x62, 0x22, 0x31,
	0x21, 0xa5, 0x5d, 0xee, 0xaa, 0x6d, 0xb7, 0x69, 0xf7, 0x08, 0x7d, 0xf5, 0xcf, 0xf5, 0xaf, 0x30,
	0xdd, 0xdd, 0x96, 0x16, 0xcf, 0x23, 0xbc, 0x40, 0x77, 0xe6, 0x9b, 0xef, 0xfb, 0x66, 0xf6, 0x76,
	0x60, 0xcd, 0xa3, 0x51, 0x44, 0xe3, 0x81, 0x47, 0xe3, 0xab, 0x60, 0xcc, 0x6e, 0x76, 0x92, 0x94,
	0x32, 0x8a, 0x0c, 0x11, 0x5e, 0x5f, 0xa9, 0xd2, 0xc5, 0x3f, 0x91, 0x5c, 0x2f, 0x6b, 0x12, 0x1a,
	0x06, 0x5e, 0x40, 0x32, 0x11, 0xb6, 0x5d, 0x58, 0x1a, 0x72, 0x16, 0x27, 0xbe, 0x26, 0x21, 0x4d,
	0x08, 0x7a, 0x0c, 0x86, 0xe0, 0xb5, 0x94, 0x4d, 0xa5, 0xdf, 0xc3, 0xf2, 0x84, 0xf6, 0x00, 0xb2,
	0x60, 0x1c, 0xbb, 0x6c, 0x9a, 0x92, 0xcc, 0x52, 0x37, 0xb5, 0x7e, 0x77, 0xf7, 0xc9, 0x8e, 0xd4,
	0x10, 0x1c, 0xa7, 0x65, 0x1e, 0xd7, 0xa0, 0xf6, 0x7e, 0x29, 0x71, 0x46, 0xa2, 0x24, 0x74, 0x19,
	0x41, 0x7d, 0x68, 0x05, 0x8c, 0x44, 0x99, 0xa5, 0x70, 0x16, 0xd4, 0x64, 0x39, 0x64, 0x24, 0xc2,
	0x02, 0x60, 0xbb, 0x60, 0x88, 0x20, 0xda, 0x06, 0x63, 0x42, 0x5c, 0x9f, 0xa4, 0xdc, 0x56, 0x77,
	0x77, 0xad, 0x2a, 0x9a, 0xb8, 0x71, 0x4c, 0xc2, 0x03, 0x9e, 0xc4, 0x12, 0x74, 0x2b, 0xa1, 0xde,
	0x27, 0xf1, 0x13, 0x40, 0x04, 0x8f, 0xc9, 0x0d, 0x7b, 0xa8, 0xcc, 0x36, 0xb4, 0x3d, 0x91, 0xb0,
	0x54, 0x8e, 0x5f, 0x69, 0x0a, 0x7d, 0x4e, 0xe9, 0x34, 0xc1, 0x25, 0xc6, 0xfe, 0xa3, 0x41, 0xb7,
	0x96, 0x40, 0x16, 0xb4, 0xaf, 0x49, 0x9a, 0x05, 0x34, 0xe6, 0x72, 0x3a, 0x2e, 0x8f, 0x68, 0x0f,
	0x8c, 0x71, 0x01, 0x29, 0x1b, 0x78, 0x36, 0x83, 0x77, 0x87, 0xff, 0xcd, 0x9c, 0x98, 0xa5, 0x39,
	0x96, 0xf0, 0xa2, 0xf0, 0xda, 0x0d, 0xa7, 0x24, 0xb3, 0xb4, 0xff, 0x17, 0x7e, 0xe3, 0x08, 0x59,
	0x28, 0xe0, 0xe8, 0x3d, 0x2c, 0x94, 0xbf, 0x0d, 0x4b, 0xe7, 0xa5, 0xcf, 0x67, 0x95, 0x8e, 0x24,
	0x46, 0x14, 0x57, 0x25, 0x68, 0x03, 0x20, 0xa2, 0xfe, 0x05, 0x3f, 0xe7, 0x56, 0x6b, 0x53, 0xe9,
	0x77, 0x70, 0x27, 0xa2, 0x3e, 0xc7, 0xe7, 0xeb, 0xc7, 0xd0, 0xad, 0xb9, 0x45, 0x26, 0x68, 0xbf,
	0x48, 0xce, 0x9b, 0xee, 0xe0, 0xe2, 0x13, 0x6d, 0x41, 0x8b, 0x1b, 0x99, 0x37, 0x47, 0x81, 0xd8,
	0x57, 0xdf, 0x29, 0x05, 0x5f, 0xad, 0x89, 0x07, 0xf3, 0xf1, 0xda, 0x3a, 0xdf, 0x57, 0x58, 0x6c,
	0x74, 0x36, 0x83, 0xf1, 0x55, 0x93, 0x71, 0xb5, 0xc9, 0x28, 0xfa, 0xac, 0x51, 0xda, 0x3f, 0xca,
	0xbb, 0xe6, 0x62, 0x73, 0xee, 0x7a, 0xb5, 0x4e, 0xdc, 0x93, 0x14, 0x77, 0x06, 0xaa, 0xdd, 0x19,
	0xa8, 0x4d, 0xa1, 0x57, 0x17, 0x9e, 0x43, 0xff, 0x12, 0x0c, 0x49, 0x22, 0x8c, 0x2f, 0x95, 0xc6,
	0xa5, 0x65, 0x99, 0xbd, 0x4f, 0xf0, 0xb7, 0x5a, 0x3e, 0x94, 0xe2, 0xf5, 0xa0, 0x37, 0xa0, 0xb3,
	0x3c, 0x21, 0x5c, 0x6c, 0x69, 0x77, 0xe3, 0xdf, 0xf7, 0x25, 0x3f, 0xcf, 0xf2, 0x84, 0x60, 0x0e,
	0x45, 0x2f, 0x60, 0x31, 0x74, 0x33, 0x76, 0x11, 0x51, 0x3f, 0xb8, 0x0a, 0x88, 0xcf, 0xfd, 0xe8,
	0xb8, 0x57, 0x04, 0x8f, 0x64, 0x0c, 0x0d, 0x60, 0x45, 0xe4, 0x3d, 0x97, 0x05, 0x34, 0x6e, 0xda,
	0x41, 0xf5, 0x94, 0x6c, 0x5c, 0x5e, 0x94, 0x7e, 0x7b, 0x51, 0xd5, 0x3c, 0x5b, 0xb5, 0x79, 0xda,
	0xc3, 0xd2, 0x7e, 0xe1, 0x08, 0x01, 0x18, 0xa3, 0x93, 0x2f, 0x87, 0xc3, 0xef, 0xe6, 0x23, 0xd4,
	0x81, 0xd6, 0xf0, 0xe0, 0xc3, 0xe1, 0xb1, 0xa9, 0xa0, 0x2e, 0xb4, 0x4f, 0xf0, 0x27, 0x07, 0x3b,
	0xd8, 0x54, 0xd1, 0x02, 0xe8, 0x23, 0xc7, 0xc1, 0xa6, 0x86, 0xda, 0xa0, 0x1d, 0x9d, 0x8e, 0x4c,
	0xdd, 0x3e, 0x87, 0xe5, 0x3b, 0xab, 0x0e, 0x6d, 0x81, 0x59, 0x2d, 0xbb, 0x8b, 0xda, 0xee, 0xe8,
	0xe1, 0xe5, 0x2a, 0x2e, 0xb6, 0x06, 0x7a, 0x0a, 0x9d, 0x2a, 0x24, 0x2f, 0xfb, 0x36, 0xf0, 0x71,
	0xfb, 0xfc, 0xf5, 0x38, 0x60, 0x93, 0xe9, 0x65, 0x31, 0xcb, 0xc1, 0x24, 0x4f, 0x48, 0x1a, 0x12,
	0x7f, 0x4c, 0xd2, 0xc1, 0x95, 0x7b, 0x99, 0x06, 0xde, 0x80, 0x6f, 0xec, 0x4c, 0xae, 0xf5, 0x4b,
	0x83, 0x1f, 0xdf, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x67, 0x29, 0x46, 0xff, 0x0d, 0x06, 0x00,
	0x00,
}
