// Code generated by protoc-gen-go. DO NOT EDIT.
// source: snapshot.proto

package snapshot

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StorageItem struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Val                  []byte   `protobuf:"bytes,2,opt,name=val,proto3" json:"val,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StorageItem) Reset()         { *m = StorageItem{} }
func (m *StorageItem) String() string { return proto.CompactTextString(m) }
func (*StorageItem) ProtoMessage()    {}
func (*StorageItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_snapshot_1383887adb7ef841, []int{0}
}
func (m *StorageItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StorageItem.Unmarshal(m, b)
}
func (m *StorageItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StorageItem.Marshal(b, m, deterministic)
}
func (dst *StorageItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageItem.Merge(dst, src)
}
func (m *StorageItem) XXX_Size() int {
	return xxx_messageInfo_StorageItem.Size(m)
}
func (m *StorageItem) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageItem.DiscardUnknown(m)
}

var xxx_messageInfo_StorageItem proto.InternalMessageInfo

func (m *StorageItem) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *StorageItem) GetVal() []byte {
	if m != nil {
		return m.Val
	}
	return nil
}

type StorageItems struct {
	Items                []*StorageItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *StorageItems) Reset()         { *m = StorageItems{} }
func (m *StorageItems) String() string { return proto.CompactTextString(m) }
func (*StorageItems) ProtoMessage()    {}
func (*StorageItems) Descriptor() ([]byte, []int) {
	return fileDescriptor_snapshot_1383887adb7ef841, []int{1}
}
func (m *StorageItems) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StorageItems.Unmarshal(m, b)
}
func (m *StorageItems) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StorageItems.Marshal(b, m, deterministic)
}
func (dst *StorageItems) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageItems.Merge(dst, src)
}
func (m *StorageItems) XXX_Size() int {
	return xxx_messageInfo_StorageItems.Size(m)
}
func (m *StorageItems) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageItems.DiscardUnknown(m)
}

var xxx_messageInfo_StorageItems proto.InternalMessageInfo

func (m *StorageItems) GetItems() []*StorageItem {
	if m != nil {
		return m.Items
	}
	return nil
}

type AccountDetail struct {
	Nonce                []byte        `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Balance              []byte        `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	StorageRoot          []byte        `protobuf:"bytes,3,opt,name=storageRoot,proto3" json:"storageRoot,omitempty"`
	CodeHash             []byte        `protobuf:"bytes,4,opt,name=codeHash,proto3" json:"codeHash,omitempty"`
	Rlp                  []byte        `protobuf:"bytes,5,opt,name=rlp,proto3" json:"rlp,omitempty"`
	Storage              *StorageItems `protobuf:"bytes,6,opt,name=storage,proto3" json:"storage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *AccountDetail) Reset()         { *m = AccountDetail{} }
func (m *AccountDetail) String() string { return proto.CompactTextString(m) }
func (*AccountDetail) ProtoMessage()    {}
func (*AccountDetail) Descriptor() ([]byte, []int) {
	return fileDescriptor_snapshot_1383887adb7ef841, []int{2}
}
func (m *AccountDetail) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountDetail.Unmarshal(m, b)
}
func (m *AccountDetail) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountDetail.Marshal(b, m, deterministic)
}
func (dst *AccountDetail) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountDetail.Merge(dst, src)
}
func (m *AccountDetail) XXX_Size() int {
	return xxx_messageInfo_AccountDetail.Size(m)
}
func (m *AccountDetail) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountDetail.DiscardUnknown(m)
}

var xxx_messageInfo_AccountDetail proto.InternalMessageInfo

func (m *AccountDetail) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

func (m *AccountDetail) GetBalance() []byte {
	if m != nil {
		return m.Balance
	}
	return nil
}

func (m *AccountDetail) GetStorageRoot() []byte {
	if m != nil {
		return m.StorageRoot
	}
	return nil
}

func (m *AccountDetail) GetCodeHash() []byte {
	if m != nil {
		return m.CodeHash
	}
	return nil
}

func (m *AccountDetail) GetRlp() []byte {
	if m != nil {
		return m.Rlp
	}
	return nil
}

func (m *AccountDetail) GetStorage() *StorageItems {
	if m != nil {
		return m.Storage
	}
	return nil
}

type Account struct {
	Addr                 []byte         `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Detail               *AccountDetail `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_snapshot_1383887adb7ef841, []int{3}
}
func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (dst *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(dst, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetAddr() []byte {
	if m != nil {
		return m.Addr
	}
	return nil
}

func (m *Account) GetDetail() *AccountDetail {
	if m != nil {
		return m.Detail
	}
	return nil
}

func init() {
	proto.RegisterType((*StorageItem)(nil), "snapshot.StorageItem")
	proto.RegisterType((*StorageItems)(nil), "snapshot.StorageItems")
	proto.RegisterType((*AccountDetail)(nil), "snapshot.AccountDetail")
	proto.RegisterType((*Account)(nil), "snapshot.Account")
}

func init() { proto.RegisterFile("snapshot.proto", fileDescriptor_snapshot_1383887adb7ef841) }

var fileDescriptor_snapshot_1383887adb7ef841 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcb, 0x4a, 0xfc, 0x30,
	0x14, 0xc6, 0xc9, 0xbf, 0xd3, 0x76, 0x38, 0x99, 0xbf, 0xc8, 0xc1, 0x4b, 0x70, 0x55, 0xba, 0x2a,
	0x08, 0xa3, 0xd6, 0xa5, 0x2b, 0xc1, 0x85, 0x6e, 0x5c, 0xc4, 0x27, 0xc8, 0xb4, 0xc1, 0x19, 0x8c,
	0x4d, 0x49, 0xe2, 0x80, 0x8f, 0xe7, 0x9b, 0x49, 0x2e, 0x1d, 0x2b, 0xb8, 0x3b, 0xdf, 0x97, 0xef,
	0x5c, 0x7e, 0x04, 0x8e, 0xec, 0x20, 0x46, 0xbb, 0xd5, 0x6e, 0x3d, 0x1a, 0xed, 0x34, 0x2e, 0x27,
	0x5d, 0xdf, 0x00, 0x7d, 0x71, 0xda, 0x88, 0x57, 0xf9, 0xe4, 0xe4, 0x3b, 0x1e, 0x43, 0xf6, 0x26,
	0x3f, 0x19, 0xa9, 0x48, 0xb3, 0xe2, 0xbe, 0xf4, 0xce, 0x5e, 0x28, 0xf6, 0x2f, 0x3a, 0x7b, 0xa1,
	0xea, 0x3b, 0x58, 0xcd, 0x5a, 0x2c, 0x5e, 0x42, 0xbe, 0xf3, 0x05, 0x23, 0x55, 0xd6, 0xd0, 0xf6,
	0x74, 0x7d, 0x58, 0x36, 0x8b, 0xf1, 0x98, 0xa9, 0xbf, 0x08, 0xfc, 0xbf, 0xef, 0x3a, 0xfd, 0x31,
	0xb8, 0x07, 0xe9, 0xc4, 0x4e, 0xe1, 0x09, 0xe4, 0x83, 0x1e, 0x3a, 0x99, 0x96, 0x46, 0x81, 0x0c,
	0xca, 0x8d, 0x50, 0xc2, 0xfb, 0x71, 0xf5, 0x24, 0xb1, 0x02, 0x6a, 0xe3, 0x5c, 0xae, 0xb5, 0x63,
	0x59, 0x78, 0x9d, 0x5b, 0x78, 0x01, 0xcb, 0x4e, 0xf7, 0xf2, 0x51, 0xd8, 0x2d, 0x5b, 0x84, 0xe7,
	0x83, 0xf6, 0x38, 0x46, 0x8d, 0x2c, 0x8f, 0x38, 0x46, 0x8d, 0x78, 0x0d, 0x65, 0x6a, 0x66, 0x45,
	0x45, 0x1a, 0xda, 0x9e, 0xfd, 0x09, 0x60, 0xf9, 0x14, 0xab, 0x9f, 0xa1, 0x4c, 0x08, 0x88, 0xb0,
	0x10, 0x7d, 0x6f, 0xd2, 0xed, 0xa1, 0xc6, 0x2b, 0x28, 0xfa, 0x80, 0x16, 0x2e, 0xa7, 0xed, 0xf9,
	0xcf, 0xbc, 0x5f, 0xe4, 0x3c, 0xc5, 0x36, 0x45, 0xf8, 0x94, 0xdb, 0xef, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x12, 0x1b, 0x18, 0xd1, 0xa6, 0x01, 0x00, 0x00,
}
