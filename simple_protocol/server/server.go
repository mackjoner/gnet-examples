package main

import (
	"flag"
	"fmt"
	"sync/atomic"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"

	"github.com/gnet-io/gnet-examples/simple_protocol/protocol"
)

type simpleServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	network   string
	addr      string
	multicore bool
	connected int64
}

// OnBoot tcp server run
func (s *simpleServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running server on %s with multi-core=%t", fmt.Sprintf("%s://%s", s.network, s.addr), s.multicore)
	s.eng = eng
	return
}

// OnOpen tcp connected
func (s *simpleServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("============= 客户端服务端建立连接 ===========")
	c.SetContext(new(protocol.SimpleCodec))
	atomic.AddInt64(&s.connected, 1)
	return
}

// OnClose tcp connect closeed
func (s *simpleServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("============= 客户端服务端连接关闭 ===========")
	if err != nil {
		logging.Infof("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt64(&s.connected, -1)
	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return
}

// OnTracffice processing received data
func (s *simpleServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	fmt.Println("============= 服务端数据接收 ===========")
	// protocol
	codec := c.Context().(*protocol.SimpleCodec)
	var packets [][]byte
	for {
		// decode receive data
		fmt.Println("============= 服务端循环处理收到的数据 ===========")
		data, err := codec.Decode(c)
		logging.Infof("data length: %d", len(data))
		if err == protocol.ErrIncompletePacket {
			fmt.Println("============= 服务端退出循环 ===========")
			break
		}
		if err != nil {
			fmt.Println("============= 服务端处理数据有错误 ===========")
			logging.Errorf("invalid packet: %v", err)
			return gnet.Close
		}
		// TODO
		// receive data conversion to http request
		// http response write to packet
		// packet, _ := codec.Encode(data)
		packet := []byte(`{"foo":"bar"}`)
		packets = append(packets, packet)
	}
	// write data packet
	if n := len(packets); n > 1 {
		fmt.Println("============= 服务端回写数据 packets ===========")
		_, _ = c.Writev(packets)
	} else if n == 1 {
		fmt.Println("============= 服务端回写数据 packets[0] ===========")
		_, _ = c.Write(packets[0])
	}
	return
}

func main() {
	var port int
	var multicore bool

	// Example command: go run server.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore=true")
	flag.Parse()
	ss := &simpleServer{
		network:   "tcp",
		addr:      fmt.Sprintf(":%d", port),
		multicore: multicore,
	}
	err := gnet.Run(ss, ss.network+"://"+ss.addr, gnet.WithMulticore(multicore))
	logging.Infof("server exits with error: %v", err)
}
