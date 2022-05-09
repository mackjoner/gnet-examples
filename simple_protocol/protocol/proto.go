package protocol

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

var ErrIncompletePacket = errors.New("incomplete packet")

const (
	protoVersionSize  = 1
	messageTypeSize   = 1
	messageNumberSize = 2
	bodySize          = 4

	HeaderSize = 8

	MsgHeartBeat byte = 1
	MsgData      byte = 2
	MsgGzipData  byte = 3
	MsgZipData   byte = 4

	DefaultProtoVersion byte = 1
)

// SimpleCodec Protocol format:
//
// * 0    1    2        4                8
// * +----+----+--------+----------------+
// * |ver |type| number |    body len    |
// * +----+----+--------+----------------+
// * |                                   |
// * +                                   +
// * |               body bytes          |
// * +                                   +
// * |                ... ...            |
// * +-----------------------------------+
type SimpleCodec struct{}

type BodyData struct {
	URI  string `json:"uri"`
	Data []byte `json:"data"`
}

func (codec SimpleCodec) Encode(header []byte, body []byte) ([]byte, error) {
	bodyOffset := protoVersionSize + messageTypeSize + messageNumberSize + bodySize
	if header[1] == MsgZipData {
		body = EncodeZlib(body)
	}
	if header[1] == MsgGzipData {
		body = EncodeGzip(body)
	}
	msgLen := bodyOffset + len(body)
	data := make([]byte, msgLen)
	//header := make([]byte, bodyOffset)
	//header[0] = DefaultProtoVersion
	//header[1] = msgType
	//binary.BigEndian.PutUint16(header[2:4], msgNumber)
	//binary.BigEndian.PutUint32(header[4:8], uint32(len(body)))
	copy(data[:bodyOffset], header)
	copy(data[bodyOffset:msgLen], body)
	return data, nil
}

func (codec *SimpleCodec) Decode(c gnet.Conn) ([]byte, error) {
	bodyOffset := protoVersionSize + messageTypeSize + messageNumberSize + bodySize
	fmt.Printf("headerSize: %d\n", bodyOffset)
	// 消息头
	buf, _ := c.Peek(bodyOffset)
	if len(buf) < bodyOffset {
		fmt.Println("协议的消息头长度不符合")
		return nil, ErrIncompletePacket
	}

	// if !bytes.Equal(magicNumberBytes, buf[:bodyOffset]) {
	// 	return nil, errors.New("invalid magic number")
	// }

	bodyLen := binary.BigEndian.Uint32(buf[4:8])
	// messageType := buf[0]
	fmt.Println("header messageType: ", buf[1])
	fmt.Println("bodySize: ", bodyLen)
	msgLen := bodyOffset + int(bodyLen)
	fmt.Printf("msgSize: %d\n", msgLen)
	if c.InboundBuffered() < msgLen {
		fmt.Println("协议的内容长度不符合")
		return nil, ErrIncompletePacket
	}
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)
	msgData, err := decodeMsg(buf[1], buf, bodyOffset, msgLen)
	if err != nil {
		fmt.Printf("decodeMsg err: %s\n", err.Error())
		return nil, err
	}

	// TODO
	// 数据解包处理
	// 组装 http.Request
	// 调用 http 接口的 handler
	// 返回 http 接口的 response
	if buf[1] == MsgData || buf[1] == MsgZipData || buf[1] == MsgGzipData {
		//var bodyData BodyData
		_, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(msgData)
		//err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(msgData, &bodyData)
		if err != nil {
			logging.Errorf("byte unmarshal err: %s\n", err.Error())
		}
		//logging.Infof("uri: %s, request data: %+v\n", bodyData.URI, bodyData.Data)
		//url, err := url.Parse(bodyData.URI)
		//if err != nil {
		//logging.Infof("ParseRequestURI err: %s\n", err.Error())
		//}
		//logging.Infof("scheme: %s,host: %s, path: %s\n", url.Scheme, url.Host, url.Path)
		// onMessage(bodyData.URI, bodyData.Data)
	}

	// return buf[bodyOffset:msgLen], nil
	// return msg, nil
	return buf, nil
}

func onMessage(url string, buffer []byte) error {
	// bearer := os.Getenv("TOKEN")
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(buffer))
	// req.Header.Add("Authorization", "Bearer "+bearer)
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}
	// else {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		// } else {
		s := fmt.Sprintf("Response status: %s", resp.Status)
		log.Println(s)
		return errors.New(s)
	}
	// var result map[string]interface{}
	// json.NewDecoder(resp.Body).Decode(&result)
	// log.Println(result["status"])
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	// jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &result)
	log.Printf("resp: %+v\n", string(body))
	return nil
	// }
}

func decodeMsg(t byte, data []byte, bodyOffset, msgLen int) ([]byte, error) {
	switch t {
	case MsgHeartBeat:
		fmt.Println("======== MsgHeartBeat ========")
		fallthrough
	case MsgData:
		fmt.Println("======== MsgData ========")
		return data[bodyOffset:msgLen], nil
	case MsgGzipData:
		fmt.Println("======== MsgGzipData ========")
		return DecodeGzip(data[bodyOffset:msgLen])
	case MsgZipData:
		fmt.Println("======== MsgZipData ========")
		return DecodeZlib(data[bodyOffset:msgLen])
	}
	fmt.Println("======== Unknown Msg Type ========")
	return data, nil
}

func DecodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func EncodeGzip(data []byte) []byte {
	buf := bytes.NewBuffer(nil)
	writer := gzip.NewWriter(buf)
	writer.Write(data)
	writer.Flush()
	writer.Close()
	return buf.Bytes()
}

func DecodeZlib(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func EncodeZlib(data []byte) []byte {
	buf := bytes.NewBuffer(nil)
	writer := zlib.NewWriter(buf)
	writer.Write(data)
	writer.Flush()
	writer.Close()
	return buf.Bytes()
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
