package tcp_stream

import (
	"bytes"
	"encoding/binary"
)

func bin2Uint32(slice []byte) uint32 {
	var num uint32

	_ = binary.Read(bytes.NewReader(slice), binary.BigEndian, &num)
	return num
}

func bin2Uint16(slice []byte) uint16 {
	var num uint16

	_ = binary.Read(bytes.NewReader(slice), binary.BigEndian, &num)
	return num
}
