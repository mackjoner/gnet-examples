package protocol

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
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

// type Header []byte

// func NewHeader(bs []byte) Header {
// 	return Header(bs)
// }

// func (h Header) GetBytes() []byte {
// 	return h
// }

// func (h Header) SetMessageType(t byte) {
// 	h[0] = t
// }

// func (h Header) GetMessageType() byte {
// 	return h[0]
// }

// func (h Header) SetMessageNumber(num uint16) {
// 	binary.BigEndian.PutUint16(h[1:3], num)
// }

// func (h Header) GetMessageNumber() uint16 {
// 	return binary.BigEndian.Uint16(h[1:3])
// }

// func (h Header) SetServerRouter(router uint16) {
// 	binary.BigEndian.PutUint16(h[3:5], router)
// }

// func (h Header) GetSetServerRouterr() uint16 {
// 	return binary.BigEndian.Uint16(h[3:5])
// }

// func (h Header) SetMessageLength(num uint32) {
// 	binary.BigEndian.PutUint32(h[5:9], num)
// }

// func (h Header) GetMessageLength() uint32 {
// 	return binary.BigEndian.Uint32(h[5:9])
// }

func (codec SimpleCodec) Encode(body []byte, msgType byte, msgNumber uint16, serverRouter uint16) ([]byte, error) {
	bodyOffset := messageTypeSize + messageNumberSize + serverRouterSize + bodySize
	msgLen := bodyOffset + len(body)
	data := make([]byte, msgLen)
	header := make([]byte, bodyOffset)
	header[0] = msgType
	binary.BigEndian.PutUint16(header[1:3], msgNumber)
	binary.BigEndian.PutUint16(header[3:5], serverRouter)
	binary.BigEndian.PutUint32(header[5:9], uint32(len(body)))
	copy(data[:bodyOffset], header)
	copy(data[bodyOffset:msgLen], body)
	return data, nil
}

func (codec *SimpleCodec) Decode(c gnet.Conn) ([]byte, error) {
	bodyOffset := messageTypeSize + messageNumberSize + serverRouterSize + bodySize
	fmt.Printf("headerSize: %d\n", bodyOffset)
	// 消息头
	buf, _ := c.Peek(bodyOffset)
	if len(buf) < bodyOffset {
		fmt.Println("协议的消息头长度不符合")
		return nil, ErrIncompletePacket
	}
	h := NewHeader(buf)
	logging.Infof("%s", string(h.GetMessageType()))

	// if !bytes.Equal(magicNumberBytes, buf[:bodyOffset]) {
	// 	return nil, errors.New("invalid magic number")
	// }

	bodyLen := binary.BigEndian.Uint32(buf[5:9])
	messageType := buf[0]
	fmt.Printf("header messageType: %s\n", string(messageType))
	fmt.Printf("bodySize: %d\n", bodyLen)
	msgLen := bodyOffset + int(bodyLen)
	fmt.Printf("msgSize: %d\n", msgLen)
	if c.InboundBuffered() < msgLen {
		fmt.Println("协议的内容长度不符合")
		return nil, ErrIncompletePacket
	}
	buf, _ = c.Peek(msgLen)
	msg, err := decodeMsg(messageType, buf)
	if err != nil {
		return nil, err
	}
	_, _ = c.Discard(msgLen)
	// return buf[bodyOffset:msgLen], nil
	return msg, nil
}

func decodeMsg(t byte, data []byte) ([]byte, error) {
	fmt.Printf("messageType: %s\n", string(t))
	// switch t {
	// case MsgHeartBeat:
	// 	break
	// case MsgData:
	// 	return data, nil
	// case MsgGzipData:
	// 	return decodeGzip(data)
	// case MsgZipData:
	// 	return decodeZlib(data)
	// }
	return data, nil
}

func decodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func decodeZlib(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

// func (codec SimpleCodec) Unpack(buf []byte) ([]byte, error) {
// 	bodyOffset := messageTypeSize + messageNumberSize + serverRouterSize + bodySize
// 	if len(buf) < bodyOffset {
// 		return nil, ErrIncompletePacket
// 	}

// 	// if !bytes.Equal(magicNumberBytes, buf[:magicNumberSize]) {
// 	// 	return nil, errors.New("invalid magic number")
// 	// }

// 	bodyLen := binary.BigEndian.Uint32(buf[5:9])
// 	msgLen := bodyOffset + int(bodyLen)
// 	if len(buf) < msgLen {
// 		return nil, ErrIncompletePacket
// 	}

// 	return buf[bodyOffset:msgLen], nil
// }
