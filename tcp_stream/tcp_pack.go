package tcp_stream

import (
	"encoding/binary"
	"github.com/pkg/errors"
)

var errBase = errors.New("[tcp_stream->tcp_pack]")
var (
	errBaseForParsePack = errors.Wrap(errBase, "parsePack()")

	ErrParsePackInvalidHeaderLen = errors.Wrap(errBaseForParsePack, "invalid header length")
	ErrParsePackUnsupportedpVer  = errors.Wrap(errBaseForParsePack, "unsupported p version")
	ErrParsePackBadData          = errors.Wrap(errBaseForParsePack, "bad data segment")
	ErrParsePackChecksum         = errors.Wrap(errBaseForParsePack, "checksum err")

	errBaseForWrapPack     = errors.Wrap(errBase, "wrapPack()")
	ErrWrapPackNoData      = errors.Wrap(errBaseForWrapPack, "no data input")
	ErrWrapPackDataTooLong = errors.Wrap(errBaseForWrapPack, "data too long")
	//ErrParsePackReceived   = errors.Wrap(errBase, "parse received pack err")
)

const (
	pVer uint8 = 1 // 1byte 能表示 0~255

	pDataLengthLen = 4

	pVerLen      = 1
	pChecksumLen = 2
	pFixedLen    = pDataLengthLen + pVerLen + pChecksumLen
)

const (
	dataLengthLimit = 4096 // 4KB
)

type Stream []byte

func (s Stream) Len() int {
	return len(s)
}
func (s Stream) Ver() uint8 {
	return s[pDataLengthLen]
}
func (s Stream) RawDataLen() int {
	return len(s[pFixedLen:])
}
func (s Stream) Checksum() uint16 {
	return bin2Uint16(s[pDataLengthLen+pVerLen : pFixedLen])
}
func (s Stream) CalChecksum() uint16 {
	return calculateChecksum(s[pFixedLen:])
}
func (s *Stream) Write(b []byte) {
	*s = append(*s, b...)
}
func (s Stream) RawData() []byte {
	// for slice, this is no copy
	return s[pFixedLen:]
}

func parsePack(stream Stream) ([]byte, error) {
	if stream.Len() < pFixedLen {
		return nil, ErrParsePackInvalidHeaderLen
	}
	if stream.Ver() != pVer {
		return nil, ErrParsePackUnsupportedpVer
	}
	rawLen := stream.RawDataLen()
	if stream.Len() != pFixedLen+rawLen {
		return nil, ErrParsePackBadData
	}
	if stream.Checksum() != stream.CalChecksum() {
		//println(111, fmt.Sprintf("%X %X", stream.Checksum(), stream.CalChecksum()))
		return nil, ErrParsePackChecksum
	}
	return stream.RawData(), nil
}

func wrapPack(rawdata []byte) (Stream, error) {
	rawdataLen := len(rawdata)
	if rawdataLen == 0 {
		return nil, ErrWrapPackNoData
	}
	if rawdataLen > dataLengthLimit {
		return nil, ErrWrapPackDataTooLong
	}
	stream := make(Stream, pFixedLen, pFixedLen+rawdataLen)

	// 依次写 数据长度、版本号、校验和、裸数据
	binary.BigEndian.PutUint32(stream, uint32(pFixedLen+rawdataLen))
	stream[pDataLengthLen] = pVer
	binary.BigEndian.PutUint16(stream[pDataLengthLen+pVerLen:], calculateChecksum(rawdata))
	stream.Write(rawdata)

	return stream, nil
}
