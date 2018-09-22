// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: por/sigchain.proto

package por

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strconv "strconv"

import bytes "bytes"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type SigAlgo int32

const (
	ECDSA SigAlgo = 0
)

var SigAlgo_name = map[int32]string{
	0: "ECDSA",
}
var SigAlgo_value = map[string]int32{
	"ECDSA": 0,
}

func (SigAlgo) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_sigchain_7c92fb6ec426cd37, []int{0}
}

type SigChainElem struct {
	Addr       []byte  `protobuf:"bytes,1,opt,name=Addr,proto3" json:"Addr,omitempty"`
	NextPubkey []byte  `protobuf:"bytes,2,opt,name=NextPubkey,proto3" json:"NextPubkey,omitempty"`
	Mining     bool    `protobuf:"varint,3,opt,name=Mining,proto3" json:"Mining,omitempty"`
	SigAlgo    SigAlgo `protobuf:"varint,4,opt,name=SigAlgo,proto3,enum=por.SigAlgo" json:"SigAlgo,omitempty"`
	Signature  []byte  `protobuf:"bytes,5,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (m *SigChainElem) Reset()      { *m = SigChainElem{} }
func (*SigChainElem) ProtoMessage() {}
func (*SigChainElem) Descriptor() ([]byte, []int) {
	return fileDescriptor_sigchain_7c92fb6ec426cd37, []int{0}
}
func (m *SigChainElem) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SigChainElem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SigChainElem.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *SigChainElem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigChainElem.Merge(dst, src)
}
func (m *SigChainElem) XXX_Size() int {
	return m.Size()
}
func (m *SigChainElem) XXX_DiscardUnknown() {
	xxx_messageInfo_SigChainElem.DiscardUnknown(m)
}

var xxx_messageInfo_SigChainElem proto.InternalMessageInfo

func (m *SigChainElem) GetAddr() []byte {
	if m != nil {
		return m.Addr
	}
	return nil
}

func (m *SigChainElem) GetNextPubkey() []byte {
	if m != nil {
		return m.NextPubkey
	}
	return nil
}

func (m *SigChainElem) GetMining() bool {
	if m != nil {
		return m.Mining
	}
	return false
}

func (m *SigChainElem) GetSigAlgo() SigAlgo {
	if m != nil {
		return m.SigAlgo
	}
	return ECDSA
}

func (m *SigChainElem) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type SigChain struct {
	Nonce      uint32          `protobuf:"varint,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	DataSize   uint32          `protobuf:"varint,2,opt,name=DataSize,proto3" json:"DataSize,omitempty"`
	DataHash   []byte          `protobuf:"bytes,3,opt,name=DataHash,proto3" json:"DataHash,omitempty"`
	BlockHash  []byte          `protobuf:"bytes,4,opt,name=BlockHash,proto3" json:"BlockHash,omitempty"`
	SrcPubkey  []byte          `protobuf:"bytes,5,opt,name=SrcPubkey,proto3" json:"SrcPubkey,omitempty"`
	DestPubkey []byte          `protobuf:"bytes,6,opt,name=DestPubkey,proto3" json:"DestPubkey,omitempty"`
	Elems      []*SigChainElem `protobuf:"bytes,7,rep,name=Elems" json:"Elems,omitempty"`
}

func (m *SigChain) Reset()      { *m = SigChain{} }
func (*SigChain) ProtoMessage() {}
func (*SigChain) Descriptor() ([]byte, []int) {
	return fileDescriptor_sigchain_7c92fb6ec426cd37, []int{1}
}
func (m *SigChain) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SigChain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SigChain.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *SigChain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigChain.Merge(dst, src)
}
func (m *SigChain) XXX_Size() int {
	return m.Size()
}
func (m *SigChain) XXX_DiscardUnknown() {
	xxx_messageInfo_SigChain.DiscardUnknown(m)
}

var xxx_messageInfo_SigChain proto.InternalMessageInfo

func (m *SigChain) GetNonce() uint32 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *SigChain) GetDataSize() uint32 {
	if m != nil {
		return m.DataSize
	}
	return 0
}

func (m *SigChain) GetDataHash() []byte {
	if m != nil {
		return m.DataHash
	}
	return nil
}

func (m *SigChain) GetBlockHash() []byte {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

func (m *SigChain) GetSrcPubkey() []byte {
	if m != nil {
		return m.SrcPubkey
	}
	return nil
}

func (m *SigChain) GetDestPubkey() []byte {
	if m != nil {
		return m.DestPubkey
	}
	return nil
}

func (m *SigChain) GetElems() []*SigChainElem {
	if m != nil {
		return m.Elems
	}
	return nil
}

func init() {
	proto.RegisterType((*SigChainElem)(nil), "por.SigChainElem")
	proto.RegisterType((*SigChain)(nil), "por.SigChain")
	proto.RegisterEnum("por.SigAlgo", SigAlgo_name, SigAlgo_value)
}
func (x SigAlgo) String() string {
	s, ok := SigAlgo_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *SigChainElem) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SigChainElem)
	if !ok {
		that2, ok := that.(SigChainElem)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Addr, that1.Addr) {
		return false
	}
	if !bytes.Equal(this.NextPubkey, that1.NextPubkey) {
		return false
	}
	if this.Mining != that1.Mining {
		return false
	}
	if this.SigAlgo != that1.SigAlgo {
		return false
	}
	if !bytes.Equal(this.Signature, that1.Signature) {
		return false
	}
	return true
}
func (this *SigChain) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SigChain)
	if !ok {
		that2, ok := that.(SigChain)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Nonce != that1.Nonce {
		return false
	}
	if this.DataSize != that1.DataSize {
		return false
	}
	if !bytes.Equal(this.DataHash, that1.DataHash) {
		return false
	}
	if !bytes.Equal(this.BlockHash, that1.BlockHash) {
		return false
	}
	if !bytes.Equal(this.SrcPubkey, that1.SrcPubkey) {
		return false
	}
	if !bytes.Equal(this.DestPubkey, that1.DestPubkey) {
		return false
	}
	if len(this.Elems) != len(that1.Elems) {
		return false
	}
	for i := range this.Elems {
		if !this.Elems[i].Equal(that1.Elems[i]) {
			return false
		}
	}
	return true
}
func (this *SigChainElem) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&por.SigChainElem{")
	s = append(s, "Addr: "+fmt.Sprintf("%#v", this.Addr)+",\n")
	s = append(s, "NextPubkey: "+fmt.Sprintf("%#v", this.NextPubkey)+",\n")
	s = append(s, "Mining: "+fmt.Sprintf("%#v", this.Mining)+",\n")
	s = append(s, "SigAlgo: "+fmt.Sprintf("%#v", this.SigAlgo)+",\n")
	s = append(s, "Signature: "+fmt.Sprintf("%#v", this.Signature)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *SigChain) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 11)
	s = append(s, "&por.SigChain{")
	s = append(s, "Nonce: "+fmt.Sprintf("%#v", this.Nonce)+",\n")
	s = append(s, "DataSize: "+fmt.Sprintf("%#v", this.DataSize)+",\n")
	s = append(s, "DataHash: "+fmt.Sprintf("%#v", this.DataHash)+",\n")
	s = append(s, "BlockHash: "+fmt.Sprintf("%#v", this.BlockHash)+",\n")
	s = append(s, "SrcPubkey: "+fmt.Sprintf("%#v", this.SrcPubkey)+",\n")
	s = append(s, "DestPubkey: "+fmt.Sprintf("%#v", this.DestPubkey)+",\n")
	if this.Elems != nil {
		s = append(s, "Elems: "+fmt.Sprintf("%#v", this.Elems)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringSigchain(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *SigChainElem) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SigChainElem) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Addr) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(len(m.Addr)))
		i += copy(dAtA[i:], m.Addr)
	}
	if len(m.NextPubkey) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(len(m.NextPubkey)))
		i += copy(dAtA[i:], m.NextPubkey)
	}
	if m.Mining {
		dAtA[i] = 0x18
		i++
		if m.Mining {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.SigAlgo != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(m.SigAlgo))
	}
	if len(m.Signature) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(len(m.Signature)))
		i += copy(dAtA[i:], m.Signature)
	}
	return i, nil
}

func (m *SigChain) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SigChain) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Nonce != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(m.Nonce))
	}
	if m.DataSize != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(m.DataSize))
	}
	if len(m.DataHash) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(len(m.DataHash)))
		i += copy(dAtA[i:], m.DataHash)
	}
	if len(m.BlockHash) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(len(m.BlockHash)))
		i += copy(dAtA[i:], m.BlockHash)
	}
	if len(m.SrcPubkey) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(len(m.SrcPubkey)))
		i += copy(dAtA[i:], m.SrcPubkey)
	}
	if len(m.DestPubkey) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintSigchain(dAtA, i, uint64(len(m.DestPubkey)))
		i += copy(dAtA[i:], m.DestPubkey)
	}
	if len(m.Elems) > 0 {
		for _, msg := range m.Elems {
			dAtA[i] = 0x3a
			i++
			i = encodeVarintSigchain(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintSigchain(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedSigChainElem(r randySigchain, easy bool) *SigChainElem {
	this := &SigChainElem{}
	v1 := r.Intn(100)
	this.Addr = make([]byte, v1)
	for i := 0; i < v1; i++ {
		this.Addr[i] = byte(r.Intn(256))
	}
	v2 := r.Intn(100)
	this.NextPubkey = make([]byte, v2)
	for i := 0; i < v2; i++ {
		this.NextPubkey[i] = byte(r.Intn(256))
	}
	this.Mining = bool(bool(r.Intn(2) == 0))
	this.SigAlgo = SigAlgo([]int32{0}[r.Intn(1)])
	v3 := r.Intn(100)
	this.Signature = make([]byte, v3)
	for i := 0; i < v3; i++ {
		this.Signature[i] = byte(r.Intn(256))
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedSigChain(r randySigchain, easy bool) *SigChain {
	this := &SigChain{}
	this.Nonce = uint32(r.Uint32())
	this.DataSize = uint32(r.Uint32())
	v4 := r.Intn(100)
	this.DataHash = make([]byte, v4)
	for i := 0; i < v4; i++ {
		this.DataHash[i] = byte(r.Intn(256))
	}
	v5 := r.Intn(100)
	this.BlockHash = make([]byte, v5)
	for i := 0; i < v5; i++ {
		this.BlockHash[i] = byte(r.Intn(256))
	}
	v6 := r.Intn(100)
	this.SrcPubkey = make([]byte, v6)
	for i := 0; i < v6; i++ {
		this.SrcPubkey[i] = byte(r.Intn(256))
	}
	v7 := r.Intn(100)
	this.DestPubkey = make([]byte, v7)
	for i := 0; i < v7; i++ {
		this.DestPubkey[i] = byte(r.Intn(256))
	}
	if r.Intn(10) != 0 {
		v8 := r.Intn(5)
		this.Elems = make([]*SigChainElem, v8)
		for i := 0; i < v8; i++ {
			this.Elems[i] = NewPopulatedSigChainElem(r, easy)
		}
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randySigchain interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneSigchain(r randySigchain) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringSigchain(r randySigchain) string {
	v9 := r.Intn(100)
	tmps := make([]rune, v9)
	for i := 0; i < v9; i++ {
		tmps[i] = randUTF8RuneSigchain(r)
	}
	return string(tmps)
}
func randUnrecognizedSigchain(r randySigchain, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldSigchain(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldSigchain(dAtA []byte, r randySigchain, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateSigchain(dAtA, uint64(key))
		v10 := r.Int63()
		if r.Intn(2) == 0 {
			v10 *= -1
		}
		dAtA = encodeVarintPopulateSigchain(dAtA, uint64(v10))
	case 1:
		dAtA = encodeVarintPopulateSigchain(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateSigchain(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateSigchain(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateSigchain(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateSigchain(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *SigChainElem) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Addr)
	if l > 0 {
		n += 1 + l + sovSigchain(uint64(l))
	}
	l = len(m.NextPubkey)
	if l > 0 {
		n += 1 + l + sovSigchain(uint64(l))
	}
	if m.Mining {
		n += 2
	}
	if m.SigAlgo != 0 {
		n += 1 + sovSigchain(uint64(m.SigAlgo))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovSigchain(uint64(l))
	}
	return n
}

func (m *SigChain) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Nonce != 0 {
		n += 1 + sovSigchain(uint64(m.Nonce))
	}
	if m.DataSize != 0 {
		n += 1 + sovSigchain(uint64(m.DataSize))
	}
	l = len(m.DataHash)
	if l > 0 {
		n += 1 + l + sovSigchain(uint64(l))
	}
	l = len(m.BlockHash)
	if l > 0 {
		n += 1 + l + sovSigchain(uint64(l))
	}
	l = len(m.SrcPubkey)
	if l > 0 {
		n += 1 + l + sovSigchain(uint64(l))
	}
	l = len(m.DestPubkey)
	if l > 0 {
		n += 1 + l + sovSigchain(uint64(l))
	}
	if len(m.Elems) > 0 {
		for _, e := range m.Elems {
			l = e.Size()
			n += 1 + l + sovSigchain(uint64(l))
		}
	}
	return n
}

func sovSigchain(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSigchain(x uint64) (n int) {
	return sovSigchain(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *SigChainElem) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SigChainElem{`,
		`Addr:` + fmt.Sprintf("%v", this.Addr) + `,`,
		`NextPubkey:` + fmt.Sprintf("%v", this.NextPubkey) + `,`,
		`Mining:` + fmt.Sprintf("%v", this.Mining) + `,`,
		`SigAlgo:` + fmt.Sprintf("%v", this.SigAlgo) + `,`,
		`Signature:` + fmt.Sprintf("%v", this.Signature) + `,`,
		`}`,
	}, "")
	return s
}
func (this *SigChain) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SigChain{`,
		`Nonce:` + fmt.Sprintf("%v", this.Nonce) + `,`,
		`DataSize:` + fmt.Sprintf("%v", this.DataSize) + `,`,
		`DataHash:` + fmt.Sprintf("%v", this.DataHash) + `,`,
		`BlockHash:` + fmt.Sprintf("%v", this.BlockHash) + `,`,
		`SrcPubkey:` + fmt.Sprintf("%v", this.SrcPubkey) + `,`,
		`DestPubkey:` + fmt.Sprintf("%v", this.DestPubkey) + `,`,
		`Elems:` + strings.Replace(fmt.Sprintf("%v", this.Elems), "SigChainElem", "SigChainElem", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringSigchain(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *SigChainElem) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSigchain
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SigChainElem: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SigChainElem: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addr = append(m.Addr[:0], dAtA[iNdEx:postIndex]...)
			if m.Addr == nil {
				m.Addr = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextPubkey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NextPubkey = append(m.NextPubkey[:0], dAtA[iNdEx:postIndex]...)
			if m.NextPubkey == nil {
				m.NextPubkey = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mining", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Mining = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigAlgo", wireType)
			}
			m.SigAlgo = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigAlgo |= (SigAlgo(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSigchain(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSigchain
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
func (m *SigChain) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSigchain
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SigChain: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SigChain: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataSize", wireType)
			}
			m.DataSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DataSize |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataHash = append(m.DataHash[:0], dAtA[iNdEx:postIndex]...)
			if m.DataHash == nil {
				m.DataHash = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BlockHash = append(m.BlockHash[:0], dAtA[iNdEx:postIndex]...)
			if m.BlockHash == nil {
				m.BlockHash = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrcPubkey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SrcPubkey = append(m.SrcPubkey[:0], dAtA[iNdEx:postIndex]...)
			if m.SrcPubkey == nil {
				m.SrcPubkey = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestPubkey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestPubkey = append(m.DestPubkey[:0], dAtA[iNdEx:postIndex]...)
			if m.DestPubkey == nil {
				m.DestPubkey = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Elems", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigchain
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Elems = append(m.Elems, &SigChainElem{})
			if err := m.Elems[len(m.Elems)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSigchain(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSigchain
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
func skipSigchain(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSigchain
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
					return 0, ErrIntOverflowSigchain
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSigchain
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthSigchain
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSigchain
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipSigchain(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthSigchain = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSigchain   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("por/sigchain.proto", fileDescriptor_sigchain_7c92fb6ec426cd37) }

var fileDescriptor_sigchain_7c92fb6ec426cd37 = []byte{
	// 387 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xbd, 0x8a, 0xdb, 0x40,
	0x14, 0x85, 0x35, 0xb1, 0xe5, 0x9f, 0x89, 0x1c, 0x92, 0xc1, 0x04, 0x61, 0xc2, 0x45, 0xb8, 0x48,
	0x44, 0x20, 0x32, 0x38, 0x6d, 0x1a, 0xff, 0x41, 0x9a, 0x98, 0x20, 0x3d, 0x81, 0x24, 0x2b, 0xe3,
	0xc1, 0xb6, 0x46, 0xe8, 0x07, 0x92, 0x54, 0xfb, 0x08, 0xfb, 0x08, 0x5b, 0xee, 0x23, 0xec, 0x23,
	0x6c, 0xe9, 0xd2, 0xcd, 0xc2, 0x6a, 0xdc, 0x6c, 0xe9, 0x72, 0xcb, 0x45, 0x23, 0x4b, 0xde, 0x6e,
	0xce, 0x39, 0xe8, 0xde, 0x73, 0x3f, 0x61, 0x12, 0xf1, 0x78, 0x94, 0x30, 0xea, 0xaf, 0x5d, 0x16,
	0x5a, 0x51, 0xcc, 0x53, 0x4e, 0x1a, 0x11, 0x8f, 0x07, 0xdf, 0x28, 0x4b, 0xd7, 0x99, 0x67, 0xf9,
	0x7c, 0x37, 0xa2, 0x9c, 0xf2, 0x91, 0xcc, 0xbc, 0xec, 0x8f, 0x54, 0x52, 0xc8, 0x57, 0xf9, 0xcd,
	0xf0, 0x06, 0x61, 0xcd, 0x61, 0x74, 0x56, 0x8c, 0x59, 0x6c, 0x83, 0x1d, 0x21, 0xb8, 0x39, 0x59,
	0xad, 0x62, 0x1d, 0x19, 0xc8, 0xd4, 0x6c, 0xf9, 0x26, 0x80, 0xf1, 0x32, 0xf8, 0x9b, 0xfe, 0xce,
	0xbc, 0x4d, 0xf0, 0x4f, 0x7f, 0x23, 0x93, 0x57, 0x0e, 0xf9, 0x88, 0x5b, 0xbf, 0x58, 0xc8, 0x42,
	0xaa, 0x37, 0x0c, 0x64, 0x76, 0xec, 0xb3, 0x22, 0x9f, 0x71, 0xdb, 0x61, 0x74, 0xb2, 0xa5, 0x5c,
	0x6f, 0x1a, 0xc8, 0x7c, 0x37, 0xd6, 0xac, 0x88, 0xc7, 0xd6, 0xd9, 0xb3, 0xab, 0x90, 0x7c, 0xc2,
	0x5d, 0x87, 0xd1, 0xd0, 0x4d, 0xb3, 0x38, 0xd0, 0x55, 0x39, 0xfe, 0x62, 0x0c, 0x1f, 0x10, 0xee,
	0x54, 0x15, 0x49, 0x1f, 0xab, 0x4b, 0x1e, 0xfa, 0x81, 0xec, 0xd7, 0xb3, 0x4b, 0x41, 0x06, 0xb8,
	0x33, 0x77, 0x53, 0xd7, 0x61, 0xff, 0x03, 0x59, 0xaf, 0x67, 0xd7, 0xba, 0xca, 0x7e, 0xba, 0xc9,
	0x5a, 0xd6, 0xd3, 0xec, 0x5a, 0x17, 0x8b, 0xa7, 0x5b, 0xee, 0x6f, 0x64, 0xd8, 0x2c, 0x17, 0xd7,
	0x86, 0xac, 0x15, 0xfb, 0xe7, 0xab, 0xab, 0x5a, 0x95, 0x51, 0x40, 0x99, 0x07, 0x49, 0x05, 0xa5,
	0x55, 0x42, 0xb9, 0x38, 0xe4, 0x0b, 0x56, 0x0b, 0xa0, 0x89, 0xde, 0x36, 0x1a, 0xe6, 0xdb, 0xf1,
	0x87, 0xea, 0xf4, 0x1a, 0xb5, 0x5d, 0xe6, 0x5f, 0xfb, 0x35, 0x25, 0xd2, 0xc5, 0xea, 0x62, 0x36,
	0x77, 0x26, 0xef, 0x95, 0xe9, 0x8f, 0x7d, 0x0e, 0xca, 0x21, 0x07, 0xe5, 0x94, 0x03, 0x7a, 0xce,
	0x01, 0x5d, 0x09, 0x40, 0xb7, 0x02, 0xd0, 0x9d, 0x00, 0x74, 0x2f, 0x00, 0xed, 0x05, 0xa0, 0x47,
	0x01, 0xe8, 0x49, 0x80, 0x72, 0x12, 0x80, 0xae, 0x8f, 0xa0, 0xec, 0x8f, 0xa0, 0x1c, 0x8e, 0xa0,
	0x78, 0x2d, 0xf9, 0x77, 0xbf, 0xbf, 0x04, 0x00, 0x00, 0xff, 0xff, 0x4a, 0xe4, 0x21, 0x0f, 0x27,
	0x02, 0x00, 0x00,
}
