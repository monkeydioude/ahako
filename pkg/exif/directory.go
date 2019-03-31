package exif

import (
	"encoding/binary"
	"fmt"
)

const (
	// IFD0th contains label for 0th IFD
	IFD0th = "IFD_0TH"
	// IFDEXIF contains label for EXIF IFD
	IFDEXIF = "IFD_EXIF"
	// IFDGPS contains label for GPS IFD
	IFDGPS = "IFD_GPS"
	// IFDInteroperability contains label for Interoperability IFD
	IFDInteroperability = "IFD_INTEROPERABILITY"
)

// Directory matches IFD structure of Metadata
type Directory struct {
	Name   string
	Fields map[uint16]Field
	NextID uint32
}

// NewDirectory generates a pointer to Directory. A name should be given its IFD label
func NewDirectory(name string, chunk *[]byte, offset uint32, byteOrder binary.ByteOrder) *Directory {
	d := &Directory{
		Name:   name,
		Fields: make(map[uint16]Field),
	}
	count := int(byteOrder.Uint16((*chunk)[offset : offset+2]))
	offset += 2

	for i := 0; i < count; i++ {
		t := byteOrder.Uint16((*chunk)[offset : offset+2])
		d.Fields[t] = NewField(chunk, offset, byteOrder)
		offset += 12
	}

	d.NextID = byteOrder.Uint32((*chunk)[offset : offset+4])
	return d
}

// GetField finds a field from a IFD directory using its tag
func (d *Directory) GetField(tag uint16) (Field, error) {
	if _, ok := d.Fields[tag]; !ok {
		return nil, fmt.Errorf("Unknown field tag %d in directory %s", tag, d.Name)
	}

	return d.Fields[tag], nil
}
