package tcp_stream

import (
	"github.com/pkg/errors"
)

var errBase = errors.New("[tcp_stream.tcp_pack]")
var (
	ErrParsePackInvalidHeaderLen       = errors.Wrap(errBase, "invalid header length")
	ErrParsePackUnsupportedProtocolVer = errors.Wrap(errBase, "unsupported protocol version")
	ErrParsePackBadData                = errors.Wrap(errBase, "bad data segment")
	ErrParsePackChecksum               = errors.Wrap(errBase, "checksum err")
	ErrParsePackReceived               = errors.Wrap(errBase, "parse received pack err")
)

const (
	protocolVer uint8 = 1 // 1byte 能表示 0~255

	protocolVerLen        = 1
	protocolDataLengthLen = 4
	protocolChecksumLen   = 2
	protocolHeaderLen     = protocolVerLen + protocolDataLengthLen + protocolChecksumLen
)

type Stream []byte

func (s Stream) Len() int {
	return len(s)
}
func (s Stream) Ver() uint8 {
	return s[0]
}
func (s Stream) DataLen() uint32 {
	return bin2Uint32(s[protocolVerLen : protocolVerLen+protocolDataLengthLen])
}
func (s Stream) Checksum() uint16 {
	return bin2Uint16(s[protocolVerLen+protocolDataLengthLen : protocolHeaderLen])
}
func (s Stream) CalChecksum() uint16 {
	return calculateChecksum(append(s[:protocolVerLen+protocolDataLengthLen], s[protocolHeaderLen:]...))
}

func (s Stream) Data() []byte {
	// for slice, here is no copy
	return s[protocolHeaderLen:]
}

func parsePack(stream Stream) ([]byte, error) {
	if stream.Len() < protocolHeaderLen {
		return nil, ErrParsePackInvalidHeaderLen
	}
	if stream.Ver() != protocolVer {
		return nil, ErrParsePackUnsupportedProtocolVer
	}
	dataLen := stream.DataLen()
	if uint32(len(stream)) != protocolHeaderLen+dataLen {
		return nil, ErrParsePackBadData
	}
	if stream.Checksum() != stream.CalChecksum() {
		return nil, ErrParsePackChecksum
	}
	return stream.Data(), nil
}
