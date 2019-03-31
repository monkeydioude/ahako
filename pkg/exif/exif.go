package exif

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
)

const (
	Marker_SOI   = 0xFFD8
	Marker_APP0  = 0xFFE0
	Marker_APP1  = 0xFFe1
	Marker_APP15 = 0xFFEF

	Tag_Orientation = 0x0112
)

var LittleEndian = []byte{'I', 'I'}
var BigEndian = []byte{'M', 'M'}

type Exif struct {
	Directories map[string]*Directory
}

func NewExif(segment *[]byte) (*Exif, error) {
	return &Exif{
		Directories: make(map[string]*Directory),
	}, nil
}

func isExifSign(file *os.File, chunk *[]byte) bool {
	return bytes.Equal(*chunk, []byte{'E', 'x', 'i', 'f', 0x00, 0x00})
}

func findSegment(file *os.File) ([]byte, error) {
	var segmentBuf []byte
	var marker uint16
	var len uint16

	if err := binary.Read(file, binary.BigEndian, &marker); err != nil {
		return nil, err
	}

	if marker != Marker_SOI {
		return nil, errors.New("Unhandled file type for EXIF retrieving")
	}

	for marker <= Marker_APP1 {
		err := binary.Read(file, binary.BigEndian, &marker)

		if err != nil {
			return nil, err
		}

		if marker < Marker_APP0 || Marker_APP15 < marker {
			return nil, errors.New("Could not find APP1 marker")
		}

		if marker != Marker_APP1 {
			continue
		}

		if err := binary.Read(file, binary.BigEndian, &len); err != nil {
			return nil, err
		}

		len -= 2
		segmentBuf = make([]byte, len)

		if err := binary.Read(file, binary.BigEndian, segmentBuf); err != nil {
			return nil, err
		}

		c := segmentBuf[:6]
		if isExifSign(file, &c) {
			return segmentBuf[6:], nil
		}
	}

	return nil, errors.New("Could not find EXIF segment")
}

func GetExifFromFile(file *os.File) (*Exif, error) {
	file.Seek(0, 0)
	segment, err := findSegment(file)
	if err != nil {
		return nil, err
	}

	exif, err := NewExif(nil)
	if err != nil {
		return nil, err
	}

	offset, byteOrder, err := GetHeader(&segment)
	if err != nil {
		return nil, err
	}

	exif.Directories[IFD0th] = NewDirectory(IFD0th, &segment, offset, byteOrder)

	exifIfd, err := exif.Directories[IFD0th].GetField(exifIfdPointer)
	if err != nil {
		return nil, err
	}

	exif.Directories[IFDEXIF] = NewDirectory(IFDEXIF, &segment, byteOrder.Uint32(exifIfd.GetData()), byteOrder)

	gpsIfd, err := exif.Directories[IFD0th].GetField(gpsIfdPointer)
	if err != nil {
		return nil, err
	}

	exif.Directories[IFDGPS] = NewDirectory(IFDGPS, &segment, byteOrder.Uint32(gpsIfd.GetData()), byteOrder)

	interopIfd, err := exif.Directories[IFDEXIF].GetField(interoperabilityIfdPointer)
	if err != nil {
		return nil, err
	}

	exif.Directories[IFDInteroperability] = NewDirectory(IFDInteroperability, &segment, byteOrder.Uint32(interopIfd.GetData()), byteOrder)

	return exif, nil
}

func (e *Exif) GetOrientation() uint16 {
	f, err := e.Directories[IFD0th].GetField(0x112)

	if err != nil {
		log.Printf("[WARN] %s\n", err)
		return 1
	}

	return f.GetByteOrder().Uint16(f.GetData())
}
