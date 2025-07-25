// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ggezchain/acl/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the acl module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params         Params         `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	AclAuthorities []AclAuthority `protobuf:"bytes,2,rep,name=acl_authorities,json=aclAuthorities,proto3" json:"acl_authorities"`
	AclAdmins      []AclAdmin     `protobuf:"bytes,3,rep,name=acl_admins,json=aclAdmins,proto3" json:"acl_admins"`
	SuperAdmin     *SuperAdmin    `protobuf:"bytes,4,opt,name=super_admin,json=superAdmin,proto3" json:"super_admin,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f4037a4949b2fa8, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetAclAuthorities() []AclAuthority {
	if m != nil {
		return m.AclAuthorities
	}
	return nil
}

func (m *GenesisState) GetAclAdmins() []AclAdmin {
	if m != nil {
		return m.AclAdmins
	}
	return nil
}

func (m *GenesisState) GetSuperAdmin() *SuperAdmin {
	if m != nil {
		return m.SuperAdmin
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "ggezchain.acl.GenesisState")
}

func init() { proto.RegisterFile("ggezchain/acl/genesis.proto", fileDescriptor_3f4037a4949b2fa8) }

var fileDescriptor_3f4037a4949b2fa8 = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4e, 0x4f, 0x4f, 0xad,
	0x4a, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x4f, 0x4c, 0xce, 0xd1, 0x4f, 0x4f, 0xcd, 0x4b, 0x2d, 0xce,
	0x2c, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x85, 0x4b, 0xea, 0x25, 0x26, 0xe7, 0x48,
	0x09, 0x26, 0xe6, 0x66, 0xe6, 0xe5, 0xeb, 0x83, 0x49, 0x88, 0x0a, 0x29, 0x91, 0xf4, 0xfc, 0xf4,
	0x7c, 0x30, 0x53, 0x1f, 0xc4, 0x82, 0x8a, 0x4a, 0xa1, 0x1a, 0x5a, 0x90, 0x58, 0x94, 0x98, 0x0b,
	0x35, 0x53, 0x4a, 0x11, 0x55, 0x2e, 0x31, 0x39, 0x27, 0x3e, 0xb1, 0xb4, 0x24, 0x23, 0xbf, 0x28,
	0xb3, 0xa4, 0x12, 0xaa, 0x44, 0x16, 0x8b, 0x92, 0x94, 0xdc, 0xcc, 0x3c, 0xa8, 0xb4, 0x3c, 0xaa,
	0x74, 0x71, 0x69, 0x41, 0x6a, 0x11, 0xb2, 0x02, 0xa5, 0x19, 0x4c, 0x5c, 0x3c, 0xee, 0x10, 0x8f,
	0x04, 0x97, 0x24, 0x96, 0xa4, 0x0a, 0x59, 0x70, 0xb1, 0x41, 0xdc, 0x20, 0xc1, 0xa8, 0xc0, 0xa8,
	0xc1, 0x6d, 0x24, 0xaa, 0x87, 0xe2, 0x31, 0xbd, 0x00, 0xb0, 0xa4, 0x13, 0xe7, 0x89, 0x7b, 0xf2,
	0x0c, 0x2b, 0x9e, 0x6f, 0xd0, 0x62, 0x0c, 0x82, 0xaa, 0x17, 0xf2, 0xe7, 0xe2, 0x47, 0x76, 0x61,
	0x66, 0x6a, 0xb1, 0x04, 0x93, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0x34, 0x9a, 0x11, 0x8e, 0xc9, 0x39,
	0x8e, 0x30, 0x6f, 0x20, 0x1b, 0xc4, 0x97, 0x88, 0x90, 0xc8, 0x4c, 0x2d, 0x16, 0x72, 0xe4, 0xe2,
	0x82, 0xfb, 0xa7, 0x58, 0x82, 0x19, 0x6c, 0x96, 0x38, 0x16, 0xb3, 0x40, 0xf2, 0xc8, 0xe6, 0x70,
	0x26, 0x42, 0x05, 0x8b, 0x85, 0xac, 0xb8, 0xb8, 0x91, 0xfc, 0x2c, 0xc1, 0x02, 0xf6, 0x92, 0x24,
	0x9a, 0x19, 0xc1, 0x20, 0x15, 0x60, 0x0d, 0x41, 0x5c, 0xc5, 0x70, 0xb6, 0x93, 0xdb, 0x89, 0x47,
	0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85,
	0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0xe9, 0xa4, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9,
	0x25, 0xe7, 0xe7, 0xea, 0xbb, 0xbb, 0xbb, 0x46, 0xf9, 0x24, 0x26, 0x15, 0xeb, 0x23, 0x42, 0xba,
	0xcc, 0x48, 0xbf, 0x02, 0x1c, 0xdc, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x49, 0x6c, 0xe0, 0x90, 0x36,
	0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x86, 0xd6, 0xe8, 0xa6, 0x3f, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SuperAdmin != nil {
		{
			size, err := m.SuperAdmin.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.AclAdmins) > 0 {
		for iNdEx := len(m.AclAdmins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AclAdmins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.AclAuthorities) > 0 {
		for iNdEx := len(m.AclAuthorities) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AclAuthorities[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.AclAuthorities) > 0 {
		for _, e := range m.AclAuthorities {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.AclAdmins) > 0 {
		for _, e := range m.AclAdmins {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.SuperAdmin != nil {
		l = m.SuperAdmin.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AclAuthorities", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AclAuthorities = append(m.AclAuthorities, AclAuthority{})
			if err := m.AclAuthorities[len(m.AclAuthorities)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AclAdmins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AclAdmins = append(m.AclAdmins, AclAdmin{})
			if err := m.AclAdmins[len(m.AclAdmins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SuperAdmin", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SuperAdmin == nil {
				m.SuperAdmin = &SuperAdmin{}
			}
			if err := m.SuperAdmin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
