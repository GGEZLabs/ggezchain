// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ggezchain/trade/stored_trade.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
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

type StoredTrade struct {
	TradeIndex           uint64      `protobuf:"varint,1,opt,name=trade_index,json=tradeIndex,proto3" json:"trade_index,omitempty"`
	TradeType            TradeType   `protobuf:"varint,2,opt,name=trade_type,json=tradeType,proto3,enum=ggezchain.trade.TradeType" json:"trade_type,omitempty"`
	Amount               *types.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty"`
	CoinMintingPriceUsd  string      `protobuf:"bytes,4,opt,name=coin_minting_price_usd,json=coinMintingPriceUsd,proto3" json:"coin_minting_price_usd,omitempty"`
	ReceiverAddress      string      `protobuf:"bytes,5,opt,name=receiver_address,json=receiverAddress,proto3" json:"receiver_address,omitempty"`
	Status               TradeStatus `protobuf:"varint,6,opt,name=status,proto3,enum=ggezchain.trade.TradeStatus" json:"status,omitempty"`
	Maker                string      `protobuf:"bytes,7,opt,name=maker,proto3" json:"maker,omitempty"`
	Checker              string      `protobuf:"bytes,8,opt,name=checker,proto3" json:"checker,omitempty"`
	TxDate               string      `protobuf:"bytes,9,opt,name=tx_date,json=txDate,proto3" json:"tx_date,omitempty"`
	CreateDate           string      `protobuf:"bytes,10,opt,name=create_date,json=createDate,proto3" json:"create_date,omitempty"`
	UpdateDate           string      `protobuf:"bytes,11,opt,name=update_date,json=updateDate,proto3" json:"update_date,omitempty"`
	ProcessDate          string      `protobuf:"bytes,12,opt,name=process_date,json=processDate,proto3" json:"process_date,omitempty"`
	TradeData            string      `protobuf:"bytes,13,opt,name=trade_data,json=tradeData,proto3" json:"trade_data,omitempty"`
	CoinMintingPriceJson string      `protobuf:"bytes,14,opt,name=coin_minting_price_json,json=coinMintingPriceJson,proto3" json:"coin_minting_price_json,omitempty"`
	ExchangeRateJson     string      `protobuf:"bytes,15,opt,name=exchange_rate_json,json=exchangeRateJson,proto3" json:"exchange_rate_json,omitempty"`
	BankingSystemData    string      `protobuf:"bytes,16,opt,name=banking_system_data,json=bankingSystemData,proto3" json:"banking_system_data,omitempty"`
	Result               string      `protobuf:"bytes,17,opt,name=result,proto3" json:"result,omitempty"`
}

func (m *StoredTrade) Reset()         { *m = StoredTrade{} }
func (m *StoredTrade) String() string { return proto.CompactTextString(m) }
func (*StoredTrade) ProtoMessage()    {}
func (*StoredTrade) Descriptor() ([]byte, []int) {
	return fileDescriptor_4cbf8bf5cc0260cd, []int{0}
}
func (m *StoredTrade) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StoredTrade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StoredTrade.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StoredTrade) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StoredTrade.Merge(m, src)
}
func (m *StoredTrade) XXX_Size() int {
	return m.Size()
}
func (m *StoredTrade) XXX_DiscardUnknown() {
	xxx_messageInfo_StoredTrade.DiscardUnknown(m)
}

var xxx_messageInfo_StoredTrade proto.InternalMessageInfo

func (m *StoredTrade) GetTradeIndex() uint64 {
	if m != nil {
		return m.TradeIndex
	}
	return 0
}

func (m *StoredTrade) GetTradeType() TradeType {
	if m != nil {
		return m.TradeType
	}
	return TradeType_TRADE_TYPE_UNSPECIFIED
}

func (m *StoredTrade) GetAmount() *types.Coin {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *StoredTrade) GetCoinMintingPriceUsd() string {
	if m != nil {
		return m.CoinMintingPriceUsd
	}
	return ""
}

func (m *StoredTrade) GetReceiverAddress() string {
	if m != nil {
		return m.ReceiverAddress
	}
	return ""
}

func (m *StoredTrade) GetStatus() TradeStatus {
	if m != nil {
		return m.Status
	}
	return TradeStatus_TRADE_STATUS_UNSPECIFIED
}

func (m *StoredTrade) GetMaker() string {
	if m != nil {
		return m.Maker
	}
	return ""
}

func (m *StoredTrade) GetChecker() string {
	if m != nil {
		return m.Checker
	}
	return ""
}

func (m *StoredTrade) GetTxDate() string {
	if m != nil {
		return m.TxDate
	}
	return ""
}

func (m *StoredTrade) GetCreateDate() string {
	if m != nil {
		return m.CreateDate
	}
	return ""
}

func (m *StoredTrade) GetUpdateDate() string {
	if m != nil {
		return m.UpdateDate
	}
	return ""
}

func (m *StoredTrade) GetProcessDate() string {
	if m != nil {
		return m.ProcessDate
	}
	return ""
}

func (m *StoredTrade) GetTradeData() string {
	if m != nil {
		return m.TradeData
	}
	return ""
}

func (m *StoredTrade) GetCoinMintingPriceJson() string {
	if m != nil {
		return m.CoinMintingPriceJson
	}
	return ""
}

func (m *StoredTrade) GetExchangeRateJson() string {
	if m != nil {
		return m.ExchangeRateJson
	}
	return ""
}

func (m *StoredTrade) GetBankingSystemData() string {
	if m != nil {
		return m.BankingSystemData
	}
	return ""
}

func (m *StoredTrade) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func init() {
	proto.RegisterType((*StoredTrade)(nil), "ggezchain.trade.StoredTrade")
}

func init() {
	proto.RegisterFile("ggezchain/trade/stored_trade.proto", fileDescriptor_4cbf8bf5cc0260cd)
}

var fileDescriptor_4cbf8bf5cc0260cd = []byte{
	// 527 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xc1, 0x6e, 0xd3, 0x4e,
	0x10, 0xc6, 0xe3, 0xff, 0xbf, 0x71, 0xc8, 0xa6, 0x34, 0xe9, 0xb6, 0x6a, 0x97, 0x00, 0x26, 0xf4,
	0x14, 0x24, 0x64, 0x2b, 0x29, 0x1c, 0x38, 0x02, 0x45, 0x05, 0x04, 0x12, 0x72, 0xca, 0xa5, 0x17,
	0x6b, 0x6d, 0x8f, 0x12, 0x53, 0xbc, 0x6b, 0xed, 0xae, 0xa3, 0x84, 0xa7, 0xe0, 0x41, 0x78, 0x10,
	0x8e, 0x3d, 0x72, 0x44, 0xc9, 0x8b, 0x20, 0xcf, 0xc6, 0x05, 0x45, 0xbd, 0x79, 0xbe, 0xef, 0xb7,
	0xe3, 0x6f, 0x67, 0x87, 0x9c, 0x4c, 0xa7, 0xf0, 0x2d, 0x99, 0xf1, 0x4c, 0x04, 0x46, 0xf1, 0x14,
	0x02, 0x6d, 0xa4, 0x82, 0x34, 0xc2, 0xc2, 0x2f, 0x94, 0x34, 0x92, 0x76, 0x6f, 0x18, 0x1f, 0xe5,
	0xbe, 0x97, 0x48, 0x9d, 0x4b, 0x1d, 0xc4, 0x5c, 0x43, 0x30, 0x1f, 0xc5, 0x60, 0xf8, 0x28, 0x48,
	0x64, 0x26, 0xec, 0x81, 0xfe, 0xfd, 0xed, 0xa6, 0xff, 0x74, 0x3b, 0xf9, 0xd1, 0x24, 0x9d, 0x09,
	0xfe, 0xe4, 0xa2, 0x52, 0xe9, 0x23, 0xd2, 0x41, 0x3b, 0xca, 0x44, 0x0a, 0x0b, 0xe6, 0x0c, 0x9c,
	0xe1, 0x4e, 0x48, 0x50, 0x7a, 0x57, 0x29, 0xf4, 0x05, 0xb1, 0x55, 0x64, 0x96, 0x05, 0xb0, 0xff,
	0x06, 0xce, 0x70, 0x6f, 0xdc, 0xf7, 0xb7, 0x32, 0xf9, 0xd8, 0xec, 0x62, 0x59, 0x40, 0xd8, 0x36,
	0xf5, 0x27, 0x1d, 0x11, 0x97, 0xe7, 0xb2, 0x14, 0x86, 0xfd, 0x3f, 0x70, 0x86, 0x9d, 0xf1, 0x3d,
	0xdf, 0x26, 0xf7, 0xab, 0xe4, 0xfe, 0x26, 0xb9, 0xff, 0x5a, 0x66, 0x22, 0xdc, 0x80, 0xf4, 0x94,
	0x1c, 0x55, 0x37, 0x89, 0xf2, 0x4c, 0x98, 0x4c, 0x4c, 0xa3, 0x42, 0x65, 0x09, 0x44, 0xa5, 0x4e,
	0xd9, 0xce, 0xc0, 0x19, 0xb6, 0xc3, 0x83, 0xca, 0xfd, 0x68, 0xcd, 0x4f, 0x95, 0xf7, 0x59, 0xa7,
	0xf4, 0x09, 0xe9, 0x29, 0x48, 0x20, 0x9b, 0x83, 0x8a, 0x78, 0x9a, 0x2a, 0xd0, 0x9a, 0x35, 0x11,
	0xef, 0xd6, 0xfa, 0x4b, 0x2b, 0xd3, 0x67, 0xc4, 0xd5, 0x86, 0x9b, 0x52, 0x33, 0x17, 0x6f, 0xf2,
	0xe0, 0xf6, 0x9b, 0x4c, 0x90, 0x09, 0x37, 0x2c, 0x3d, 0x24, 0xcd, 0x9c, 0x5f, 0x81, 0x62, 0x2d,
	0xec, 0x6a, 0x0b, 0xca, 0x48, 0x2b, 0x99, 0x41, 0x52, 0xe9, 0x77, 0x50, 0xaf, 0x4b, 0x7a, 0x4c,
	0x5a, 0x66, 0x11, 0xa5, 0xdc, 0x00, 0x6b, 0xa3, 0xe3, 0x9a, 0xc5, 0x19, 0x37, 0x38, 0xed, 0x44,
	0x01, 0x37, 0x60, 0x4d, 0x82, 0x26, 0xb1, 0x52, 0x0d, 0x94, 0x45, 0x7a, 0x03, 0x74, 0x2c, 0x60,
	0x25, 0x04, 0x1e, 0x93, 0xdd, 0x42, 0xc9, 0x04, 0xb4, 0xb6, 0xc4, 0x2e, 0x12, 0x9d, 0x8d, 0x86,
	0xc8, 0xc3, 0xfa, 0xc5, 0x52, 0x6e, 0x38, 0xbb, 0x8b, 0x80, 0x7d, 0x95, 0x33, 0x6e, 0x38, 0x7d,
	0x4e, 0x8e, 0x6f, 0x19, 0xf1, 0x17, 0x2d, 0x05, 0xdb, 0x43, 0xf6, 0x70, 0x7b, 0xc6, 0xef, 0xb5,
	0x14, 0xf4, 0x29, 0xa1, 0xb0, 0x48, 0x66, 0x5c, 0x4c, 0x21, 0x52, 0x55, 0x40, 0x3c, 0xd1, 0xc5,
	0x13, 0xbd, 0xda, 0x09, 0xb9, 0xb1, 0xb4, 0x4f, 0x0e, 0x62, 0x2e, 0xae, 0xaa, 0xfe, 0x7a, 0xa9,
	0x0d, 0xe4, 0x36, 0x4c, 0x0f, 0xf1, 0xfd, 0x8d, 0x35, 0x41, 0x07, 0x43, 0x1d, 0x11, 0x57, 0x81,
	0x2e, 0xbf, 0x1a, 0xb6, 0x6f, 0x07, 0x66, 0xab, 0x57, 0x6f, 0x7f, 0xae, 0x3c, 0xe7, 0x7a, 0xe5,
	0x39, 0xbf, 0x57, 0x9e, 0xf3, 0x7d, 0xed, 0x35, 0xae, 0xd7, 0x5e, 0xe3, 0xd7, 0xda, 0x6b, 0x5c,
	0xfa, 0xd3, 0xcc, 0xcc, 0xca, 0xd8, 0x4f, 0x64, 0x1e, 0x9c, 0x9f, 0xbf, 0xb9, 0xfc, 0xc0, 0x63,
	0x1d, 0xfc, 0xdd, 0xfc, 0xf9, 0x38, 0x58, 0xd4, 0xeb, 0xbf, 0x2c, 0x40, 0xc7, 0x2e, 0xee, 0xff,
	0xe9, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5e, 0xbb, 0x45, 0x4f, 0x73, 0x03, 0x00, 0x00,
}

func (m *StoredTrade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StoredTrade) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StoredTrade) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Result) > 0 {
		i -= len(m.Result)
		copy(dAtA[i:], m.Result)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.Result)))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x8a
	}
	if len(m.BankingSystemData) > 0 {
		i -= len(m.BankingSystemData)
		copy(dAtA[i:], m.BankingSystemData)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.BankingSystemData)))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x82
	}
	if len(m.ExchangeRateJson) > 0 {
		i -= len(m.ExchangeRateJson)
		copy(dAtA[i:], m.ExchangeRateJson)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.ExchangeRateJson)))
		i--
		dAtA[i] = 0x7a
	}
	if len(m.CoinMintingPriceJson) > 0 {
		i -= len(m.CoinMintingPriceJson)
		copy(dAtA[i:], m.CoinMintingPriceJson)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.CoinMintingPriceJson)))
		i--
		dAtA[i] = 0x72
	}
	if len(m.TradeData) > 0 {
		i -= len(m.TradeData)
		copy(dAtA[i:], m.TradeData)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.TradeData)))
		i--
		dAtA[i] = 0x6a
	}
	if len(m.ProcessDate) > 0 {
		i -= len(m.ProcessDate)
		copy(dAtA[i:], m.ProcessDate)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.ProcessDate)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.UpdateDate) > 0 {
		i -= len(m.UpdateDate)
		copy(dAtA[i:], m.UpdateDate)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.UpdateDate)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.CreateDate) > 0 {
		i -= len(m.CreateDate)
		copy(dAtA[i:], m.CreateDate)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.CreateDate)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.TxDate) > 0 {
		i -= len(m.TxDate)
		copy(dAtA[i:], m.TxDate)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.TxDate)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.Checker) > 0 {
		i -= len(m.Checker)
		copy(dAtA[i:], m.Checker)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.Checker)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.Maker) > 0 {
		i -= len(m.Maker)
		copy(dAtA[i:], m.Maker)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.Maker)))
		i--
		dAtA[i] = 0x3a
	}
	if m.Status != 0 {
		i = encodeVarintStoredTrade(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x30
	}
	if len(m.ReceiverAddress) > 0 {
		i -= len(m.ReceiverAddress)
		copy(dAtA[i:], m.ReceiverAddress)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.ReceiverAddress)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.CoinMintingPriceUsd) > 0 {
		i -= len(m.CoinMintingPriceUsd)
		copy(dAtA[i:], m.CoinMintingPriceUsd)
		i = encodeVarintStoredTrade(dAtA, i, uint64(len(m.CoinMintingPriceUsd)))
		i--
		dAtA[i] = 0x22
	}
	if m.Amount != nil {
		{
			size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintStoredTrade(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.TradeType != 0 {
		i = encodeVarintStoredTrade(dAtA, i, uint64(m.TradeType))
		i--
		dAtA[i] = 0x10
	}
	if m.TradeIndex != 0 {
		i = encodeVarintStoredTrade(dAtA, i, uint64(m.TradeIndex))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintStoredTrade(dAtA []byte, offset int, v uint64) int {
	offset -= sovStoredTrade(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StoredTrade) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TradeIndex != 0 {
		n += 1 + sovStoredTrade(uint64(m.TradeIndex))
	}
	if m.TradeType != 0 {
		n += 1 + sovStoredTrade(uint64(m.TradeType))
	}
	if m.Amount != nil {
		l = m.Amount.Size()
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.CoinMintingPriceUsd)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.ReceiverAddress)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovStoredTrade(uint64(m.Status))
	}
	l = len(m.Maker)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.Checker)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.TxDate)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.CreateDate)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.UpdateDate)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.ProcessDate)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.TradeData)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.CoinMintingPriceJson)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.ExchangeRateJson)
	if l > 0 {
		n += 1 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.BankingSystemData)
	if l > 0 {
		n += 2 + l + sovStoredTrade(uint64(l))
	}
	l = len(m.Result)
	if l > 0 {
		n += 2 + l + sovStoredTrade(uint64(l))
	}
	return n
}

func sovStoredTrade(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStoredTrade(x uint64) (n int) {
	return sovStoredTrade(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StoredTrade) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStoredTrade
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
			return fmt.Errorf("proto: StoredTrade: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StoredTrade: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeIndex", wireType)
			}
			m.TradeIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TradeIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeType", wireType)
			}
			m.TradeType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TradeType |= TradeType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
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
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Amount == nil {
				m.Amount = &types.Coin{}
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CoinMintingPriceUsd", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CoinMintingPriceUsd = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiverAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceiverAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= TradeStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Maker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Maker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Checker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Checker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxDate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxDate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreateDate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CreateDate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdateDate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UpdateDate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProcessDate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProcessDate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeData", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TradeData = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CoinMintingPriceJson", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CoinMintingPriceJson = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExchangeRateJson", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExchangeRateJson = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BankingSystemData", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BankingSystemData = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 17:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStoredTrade
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStoredTrade
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStoredTrade
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Result = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStoredTrade(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStoredTrade
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
func skipStoredTrade(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStoredTrade
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
					return 0, ErrIntOverflowStoredTrade
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
					return 0, ErrIntOverflowStoredTrade
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
				return 0, ErrInvalidLengthStoredTrade
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStoredTrade
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStoredTrade
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStoredTrade        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStoredTrade          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStoredTrade = fmt.Errorf("proto: unexpected end of group")
)
