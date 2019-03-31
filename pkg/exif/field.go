package exif

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

const (
	_                  = iota
	ifdFormatByte      // 1
	ifdFormatString    // 2
	ifdFormatShort     // 3
	ifdFormatULong     // 4
	ifdFormatURational // 5
	_
	ifdFormatUndefined // 7
	_
	ifdFormatSLong     // 9
	ifdFormatSRational //10
)

const (
	exifIfdPointer             = 34665
	gpsIfdPointer              = 0x8825
	interoperabilityIfdPointer = 40965
)

var formatTypeMap = map[uint16]uint32{
	ifdFormatByte:      1,
	ifdFormatString:    1,
	ifdFormatShort:     2,
	ifdFormatULong:     4,
	ifdFormatURational: 8,
	ifdFormatUndefined: 1,
	ifdFormatSLong:     4,
	ifdFormatSRational: 8,
}

type Field interface {
	GetTag() uint16
	GetType() uint16
	GetCount() uint32
	GetData() []byte
	SetData(*[]byte, uint32, binary.ByteOrder) error
	GetByteOrder() binary.ByteOrder
	String() string
}

type BaseField struct {
	Tag       uint16
	Type      uint16
	Count     uint32
	Data      []byte
	ByteOrder binary.ByteOrder
}

func (f *BaseField) GetTag() uint16 {
	return f.Tag
}

func (f *BaseField) GetType() uint16 {
	return f.Type
}

func (f *BaseField) GetCount() uint32 {
	return f.Count
}

func (f *BaseField) GetData() []byte {
	return f.Data
}

func (f *BaseField) GetByteOrder() binary.ByteOrder {
	return f.ByteOrder
}

func (f *BaseField) SetData(chunk *[]byte, offset uint32, bo binary.ByteOrder) error {
	data := (*chunk)[offset+8 : offset+12]
	byteCount := (f.Count * formatTypeMap[f.Type])
	if byteCount > 4 {
		dataOffset := bo.Uint32(data)
		data = (*chunk)[dataOffset : dataOffset+byteCount]
	}
	f.Data = data
	return nil
}

func (f *BaseField) String() string {
	return fmt.Sprintf("%v\n", f.Data)
}

func NewBaseField(chunk *[]byte, offset uint32, bo binary.ByteOrder) *BaseField {
	b := &BaseField{
		Tag:       bo.Uint16((*chunk)[offset : offset+2]),
		Type:      bo.Uint16((*chunk)[offset+2 : offset+4]),
		Count:     bo.Uint32((*chunk)[offset+4 : offset+8]),
		ByteOrder: bo,
	}

	return b
}

func NewField(chunk *[]byte, offset uint32, bo binary.ByteOrder) Field {
	b := NewBaseField(chunk, offset, bo)

	var f Field
	switch b.GetType() {
	case ifdFormatByte:
		f = &ByteField{b}
	case ifdFormatString:
		f = &StringField{b}
	case ifdFormatURational:
		f = &URational{
			BaseField: b,
			N:         make([]uint32, b.GetCount()),
			D:         make([]uint32, b.GetCount()),
		}
	case ifdFormatULong:
		f = &ULongField{b}
	case ifdFormatUndefined:
		f = &UndefinedField{b}
	case ifdFormatSRational:
		f = &SRational{BaseField: b}
	default:
		f = b
	}
	f.SetData(chunk, offset, bo)

	return f
}

// ByteField handles exif Byte Fields (type=1)
type ByteField struct {
	*BaseField
}

func (b *ByteField) String() string {
	var sb strings.Builder
	for _, by := range b.Data {
		sb.WriteString(strconv.Itoa(int(by)))
	}
	return sb.String()
}

// StringField handles exif ASCII Fields (type=2)
type StringField struct {
	*BaseField
}

func (s *StringField) String() string {
	return string(s.Data)
}

// ULongField handles exif Unsigned Long Fields (type=4)
type ULongField struct {
	*BaseField
}

func (s *ULongField) String() string {
	return fmt.Sprintf("%d", s.Data)
}

// URational handles exif Unsigned Rationals Fields (type=5)
type URational struct {
	*BaseField
	N []uint32
	D []uint32
}

func (f *URational) SetData(chunk *[]byte, offset uint32, bo binary.ByteOrder) error {
	f.BaseField.SetData(chunk, offset, bo)
	buf := bytes.NewBuffer(f.GetData())
	for i := uint32(0); i < f.GetCount(); i++ {
		binary.Read(buf, f.BaseField.ByteOrder, &f.N[i])
		binary.Read(buf, f.BaseField.ByteOrder, &f.D[i])
	}
	return nil
}

func (f *URational) String() string {
	return fmt.Sprintf("%d, %d", f.N, f.D)
}

// UndefinedField handles exif Undefined Fields (type=7)
type UndefinedField struct {
	*BaseField
}

func (f *UndefinedField) String() string {
	buf := bytes.NewBuffer(f.GetData())
	v := make([]byte, f.GetCount())
	binary.Read(buf, f.ByteOrder, v)
	return fmt.Sprintf("%02X", v)
}

// SRational handles exif Signed Rationals Fields (type=10)
type SRational struct {
	*BaseField
	N uint32
	D uint32
}

func (f *SRational) SetData(chunk *[]byte, offset uint32, bo binary.ByteOrder) error {
	data := (*chunk)[offset+8 : offset+12]
	byteCount := (f.Count * formatTypeMap[f.Type])
	if byteCount > 4 {
		offset = bo.Uint32(data)
		data = (*chunk)[offset : offset+byteCount]
	}
	f.Data = data
	f.N = bo.Uint32((*chunk)[offset : offset+4])
	f.D = bo.Uint32((*chunk)[offset+4 : offset+8])
	return nil
}

func (f *SRational) String() string {
	return fmt.Sprintf("%d, %d", f.N, f.D)
}
