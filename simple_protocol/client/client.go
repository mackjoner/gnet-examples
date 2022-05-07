package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/gnet-io/gnet-examples/simple_protocol/protocol"
	jsoniter "github.com/json-iterator/go"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

func logErr(err error) {
	logging.Error(err)
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		network     string
		addr        string
		concurrency int
		packetSize  int
		packetBatch int
		packetCount int
	)

	// Example command: go run client.go --network tcp --address ":9000" --concurrency 100 --packet_size 1024 --packet_batch 20 --packet_count 1000
	flag.StringVar(&network, "network", "tcp", "--network tcp")
	flag.StringVar(&addr, "address", "127.0.0.1:9000", "--address 127.0.0.1:9000")
	flag.IntVar(&concurrency, "concurrency", 1024, "--concurrency 500")
	flag.IntVar(&packetSize, "packet_size", 1024, "--packe_size 256")
	flag.IntVar(&packetBatch, "packet_batch", 100, "--packe_batch 100")
	flag.IntVar(&packetCount, "packet_count", 10000, "--packe_count 10000")
	flag.Parse()

	logging.Infof("start %d clients...", concurrency)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			runClient(network, addr, packetSize, packetBatch, packetCount)
			wg.Done()
		}()
	}
	wg.Wait()
	// time.Sleep(15 * time.Second)
	logging.Infof("all %d clients are done", concurrency)
}

func runClient(network, addr string, packetSize, batch, count int) {
	rand.Seed(time.Now().UnixNano())
	c, err := net.Dial(network, addr)
	logErr(err)
	logging.Infof("connection=%s starts...", c.LocalAddr().String())
	defer func() {
		logging.Infof("connection=%s stops...", c.LocalAddr().String())
		c.Close()
	}()
	rd := bufio.NewReader(c)
	// msg, err := rd.ReadBytes('\n')
	// logErr(err)
	// expectMsg := "sweetness\r\n"
	// if string(msg) != expectMsg {
	// 	logging.Fatalf("the first response packet mismatches, expect: %s, but got: %s", expectMsg, msg)
	// }

	for i := 0; i < count; i++ {
		batchSendAndRecv(c, rd, packetSize, batch)
		// time.Sleep(3 * time.Second)
		sendHeartbeat(c)
		time.Sleep(5 * time.Second)
		// fmt.Println("===================")
	}
}

func sendHeartbeat(c net.Conn) {
	codec := protocol.SimpleCodec{}
	fmt.Println("============= 构造心跳消息 ===========")
	var buf, body []byte
	packet, _ := codec.Encode(body, protocol.MsgHeartBeat, uint16(rand.Intn(65535)), '1')
	buf = append(buf, packet...)
	fmt.Println("============= 客户端发送心跳消息 ===========")
	logging.Infof("client body: %s\n", string(body))
	writeLen, err := c.Write(buf)
	logErr(err)
	fmt.Printf("conn write msg length: %d\n", writeLen)
	respPacket := make([]byte, 1024)
	fmt.Println("============= 客户端读取回执消息 ===========")
	// respPacket, err := rd.()
	readLen, err := c.Read(respPacket)
	logErr(err)
	fmt.Printf("conn read msg length: %d, type: %s\n", readLen, string(respPacket[1]))
	// fmt.Println(string(respPacket[1]))
}

func batchSendAndRecv(c net.Conn, rd *bufio.Reader, packetSize, batch int) {
	codec := protocol.SimpleCodec{}
	// var (
	// 	requests  [][]byte
	// 	buf       []byte
	// 	packetLen int
	// )
	// for i := 0; i < batch; i++ {
	// 	req := make([]byte, packetSize)
	// 	_, err := rand.Read(req)
	// 	logErr(err)
	// 	requests = append(requests, req)
	// 	packet, _ := codec.Encode(req)
	// 	packetLen = len(packet)
	// 	buf = append(buf, packet...)
	// }

	// body := bytes.NewBuffer(nil)
	// writer := zlib.NewWriter(body)
	// writer.Write([]byte(`{"uri":"/reward/setting","data":"platform=1&package_name=com.gbwhatsapp&os_version=11&brand=samsung&model=SM-A022M&gaid=ea3de793-0df3-45b9-b31b-3fd8a407132e&mnc=&mcc=&network_type=9&network_str=&language=es-US&timezone=EST&ua=Mozilla%252F5.0%2B%2528Linux%253B%2BAndroid%2B11%253B%2BSM-A022M%2BBuild%252FRP1A.200720.012%253B%2Bwv%2529%2BAppleWebKit%252F537.36%2B%2528KHTML%252C%2Blike%2BGecko%2529%2BVersion%252F4.0%2BChrome%252F99.0.4844.58%2BMobile%2BSafari%252F537.36&gp_version=29.5.14-21%2B%255B0%255D%2B%255BPR%255D%2B430999422&sdk_version=MAL_16.0.17&app_version_name=2.21.24.22&orientation=1&screen_size=720x1600&dvi=4BztYrxBYFQ3%2BFQ3RUE0fk5QHal9GniAHnQUf7V2DZRsRr2tYg5rDkfTJ%2BzQh0R1RgftY%2Bf2Yrh0WozUhdVBRUE0D%2BzwHkc0LZRsRgxtHbi0G0zBHkeQD%2BfQWkwQ4%2Bi0Woz2hF5BRUE0HdSuR0M0hrc3LkI0G0zSiaRBn55o5nzo5VsTWjjMiUzf5Vz5i3z5ZAN0Woz0YFKTY7KtH75BRUE0NnvBi325NQVBNQ5WfoRsRrtthrxbD%2BzQRUE0Y%2BNFfAiPR0M0L7KAJoR1RURexjEbxavAR0M0DFK3HkPtYkV0G0zZxVM0WozuYrfBHk2QYgxtYoR1R3jMiUzf5Vz5i3z5ZAN0WozAH%2BzuDkM0G0z2Yrw%2FYbJ%2FR0M0H7QAh7et4ZR1RQzNiVj%2FiUvMfARMWUvei0PSiaRBn55o5nzo5VsTR0M0DrKthrN0G0z8iAQTJUc6DgfM%2BbxuJ7c%2F%2BFttY%2BfTH%2BR0WozT4%2BSQRUE0J%2BfQh0RsRgf2hdSXhgN0G0zthr2QDkzuW%2BDbDZRsRgzQY75thFV0G0ReiZRsRgf3LFQ%2FJoR1iAvsRrztJdxQhgQAJ7cTJ%2Bi0GUjsRrzthF5XhBR1iZM0L%2BiBfjl0G0Rei0RsRgfQYgfXh0R1injsRrQwHZR1RrwQ4kzBYFc3ijxuDbxtH7Ilh7KBRdHX40S3HZSdYFKgY7VlLF5PDgzXDkNe575UY7c3YBSnDk2AJkPgRoRsRgSEYFPQJdQMHZR1WnRsRgxXJ7cshrcwRUE0iTJoR0M0J7KTDkewHk2Xhg30G0RAWUvAMqSdN0zK&app_id=160670&unit_id=&data=key%253D2000088%2526state%253D1%2526network_type%253D9&m_sdk=msdk&channel=&band_width=0&open=0&country_code=PA"}`))
	// writer.Flush()
	// writer.Close()
	bodyData := protocol.BodyData{
		URI:  "https://configure.rayjump.com/setting",
		Data: []byte(`{"app_id":"118692","app_version_name":"7.0.6","att":"0","brt":"0.5","c":"k-ZAJ015J-17pdXh67gvLQ%3D%3D","charging":"0","country_code":"US","ct":"%5B7%5D","dmt":"8192","env":"3","f":"xY1BvB58V3L3xpM7MClscs/LNcl1rY5WLpdXsoTlIqXN5Jp6HXN0WQEE5rJFJstW","h":"250685.578125","h_model":"MacBookPro15%2C4","http_req":"2","i":"28242.382812","idfa":"00000000-0000-0000-0000-000000000000","idfv":"04A43F02-F155-4985-AD5F-F37DE3E98B6D","keyword":"4BztHFV0GUjeWozgHkP3H%2BR0GUjsRgSt4ZR1inRsRrf2hbxXYZR1RrcAH7HtP8kW1pm5R0M0HbSARUuORretJoR1RUjMWUjFiAhAi0RsRre/HBR1RUDMWUjFiAhAi0zK6N%3D%3D","language":"en","limit_trk":"1","lpm":"0","mcc":"","mnc":"","model":"x86_64","network_str":"","network_type":"1","open":"0","orientation":"-1","os_version":"14.5","package_name":"com.mobvista.ui.test20","platform":"2","power_rate":"0","screen_size":"1170.000000x2532.000000","sdk_version":"MI_7.0.6","sign":"ba51648b420376072683755a8ec190a2","simu":"1","skad":"%7B%22ver%22%3A%222.2%22%2C%22adnetids%22%3A%5B%22test1.skadnetwork%22%2C%22test1.skadnetwork%22%2C%22kbd757ywx3.skadnetwork%22%5D%2C%22tag%22%3A%222%22%7D","st":"9747f7822bf48ec1779e97e9d862371a","timezone":"GMT%2B08%3A00","ts":"1651826586","useragent":"Mozilla/5.0%20%28iPhone%3B%20CPU%20iPhone%20OS%2014_5%20like%20Mac%20OS%20X%29%20AppleWebKit/605.1.15%20%28KHTML%2C%20like%20Gecko%29%20Mobile/15E148","vol":"0.6"}`),
	}
	bodyByte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(bodyData)
	// header := make([]byte, 9)
	// h := protocol.NewHeader(header)
	// h.SetMessageType('4')
	// h.SetMessageNumber(uint16(1314))
	// h.SetServerRouter(uint16(1))
	// h.SetMessageLength(uint32(len(body.Bytes())))
	// fmt.Printf("header: %s\n", header)

	fmt.Println("============= 构造消息 ===========")
	var buf []byte
	packet, _ := codec.Encode(bodyByte, protocol.MsgData, uint16(rand.Intn(65535)), '1')
	buf = append(buf, packet...)
	fmt.Println("============= 客户端发送消息 ===========")
	logging.Infof("%s\n", string(bodyByte))
	writeLen, err := c.Write(buf)
	logErr(err)
	fmt.Printf("conn write msg length: %d\n", writeLen)
	respPacket := make([]byte, 1024)
	fmt.Println("============= 客户端读取消息 ===========")
	// respPacket, err := rd.()
	readLen, err := c.Read(respPacket)
	logErr(err)
	fmt.Printf("conn read msg length: %d\n", readLen)
	fmt.Println(string(respPacket))
	// time.Sleep(1000 * time.Millisecond)
	// for i, req := range requests {
	// 	rsp, err := codec.Unpack(respPacket[i*packetLen:])
	// 	logErr(err)
	// 	if !bytes.Equal(req, rsp) {
	// 		logging.Fatalf("request and response mismatch, conn=%s, packet size: %d, batch: %d",
	// 			c.LocalAddr().String(), packetSize, batch)
	// 	}
	// }
}
