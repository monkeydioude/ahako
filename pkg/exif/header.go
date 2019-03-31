package exif

import (
	"encoding/binary"
	"errors"
)

type Header struct {
	ByteOrder binary.ByteOrder
	Offset    uint32
}

func GetHeader(chunk *[]byte) (uint32, binary.ByteOrder, error) {
	bo, err := getByteOrder(chunk)
	if err != nil {
		return 0, nil, err
	}

	io, err := getOffset(chunk, bo)
	if err != nil {
		return 0, nil, err
	}

	return io, bo, nil
}

func getOffset(chunk *[]byte, bo binary.ByteOrder) (uint32, error) {
	if len((*chunk)[4:8]) < 4 {
		return 0, errors.New("Can not get Offset, provided chunk length was < 4")
	}

	return bo.Uint32((*chunk)[4:8]), nil
}

func getByteOrder(chunk *[]byte) (binary.ByteOrder, error) {
	if len((*chunk)[:2]) < 2 {
		return nil, errors.New("Can not get ByteOrder, provided chunk length was < 2")
	}

	if (*chunk)[0] == LittleEndian[0] && (*chunk)[1] == LittleEndian[1] {
		return binary.LittleEndian, nil
	}

	if (*chunk)[0] == BigEndian[0] && (*chunk)[1] == BigEndian[1] {
		return binary.BigEndian, nil
	}

	return nil, errors.New("Could not define ByteOrder")
}
