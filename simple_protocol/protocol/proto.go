package protocol

import (
	"encoding/binary"
	"errors"

	"github.com/panjf2000/gnet/v2"
)

var ErrIncompletePacket = errors.New("incomplete packet")

const (
	messageTypeSize   = 1
	messageNumberSize = 2
	serverRouterSize  = 2
	bodySize          = 4

	MsgHeartBeat byte = 1
	MsgData      byte = 2
	MsgGzipData  byte = 3
	MsgZipData   byte = 4
)

type Header []byte

func NewHeader(bs []byte) Header {
	return Header(bs[0:9])
}

func (h Header) GetBytes() []byte {
	return h
}

func (h Header) SetMessageType(t byte) {
	h[0] = t
}

func (h Header) GetMessageType() byte {
	return h[0]
}

func (h Header) SetMessageNumber(num uint16) {
	binary.BigEndian.PutUint16(h[1:3], num)
}

func (h Header) GetMessageNumber() uint16 {
	return binary.BigEndian.Uint16(h[1:3])
}

func (h Header) SetServerRouter(router uint16) {
	binary.BigEndian.PutUint16(h[3:5], router)
}

func (h Header) GetSetServerRouterr() uint16 {
	return binary.BigEndian.Uint16(h[3:5])
}

func (h Header) SetMessageLength(num uint32) {
	binary.BigEndian.PutUint32(h[5:9], num)
}

func (h Header) GetMessageLength() uint32 {
	return binary.BigEndian.Uint32(h[5:9])
}

// SimpleCodec Protocol format:
//
// * 0    1        3        5                9
// * +----+--------+--------+----------------+
// * |type| number | server |    body len    |
// * +----+--------+--------+----------------+
// * |                                       |
// * +                                       +
// * |              body bytes               |
// * +                                       +
// * |               ... ...                 |
// * +---------------------------------------+
type SimpleCodec struct{}

func (codec SimpleCodec) Encode(buf []byte) ([]byte, error) {
	bodyOffset := messageTypeSize + messageNumberSize + serverRouterSize + bodySize
	msgLen := bodyOffset + len(buf)
	data := make([]byte, msgLen)
	header := make([]byte, bodyOffset)
	h := NewHeader(header)
	h.SetMessageType(MsgGzipData)
	h.SetMessageNumber(uint16(1314))
	h.SetServerRouter(uint16(1))
	h.SetMessageLength(uint32(len(buf)))
	copy(data[:bodyOffset], header)
	copy(data[bodyOffset:msgLen], buf)
	return data, nil
}

func (codec *SimpleCodec) Decode(c gnet.Conn) ([]byte, error) {
	bodyOffset := messageTypeSize + messageNumberSize + serverRouterSize + bodySize
	buf, _ := c.Peek(bodyOffset)
	if len(buf) < bodyOffset {
		return nil, ErrIncompletePacket
	}

	// if !bytes.Equal(magicNumberBytes, buf[:bodyOffset]) {
	// 	return nil, errors.New("invalid magic number")
	// }

	bodyLen := binary.BigEndian.Uint32(buf[5:9])
	msgLen := bodyOffset + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)

	return buf[bodyOffset:msgLen], nil
}

func (codec SimpleCodec) Unpack(buf []byte) ([]byte, error) {
	bodyOffset := messageTypeSize + messageNumberSize + serverRouterSize + bodySize
	if len(buf) < bodyOffset {
		return nil, ErrIncompletePacket
	}

	// if !bytes.Equal(magicNumberBytes, buf[:magicNumberSize]) {
	// 	return nil, errors.New("invalid magic number")
	// }

	bodyLen := binary.BigEndian.Uint32(buf[5:9])
	msgLen := bodyOffset + int(bodyLen)
	if len(buf) < msgLen {
		return nil, ErrIncompletePacket
	}

	return buf[bodyOffset:msgLen], nil
}
