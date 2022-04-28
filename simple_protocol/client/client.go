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
	}
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
	body := []byte(`{"uri":"/reward/setting","data":"platform=1&package_name=com.gbwhatsapp&os_version=11&brand=samsung&model=SM-A022M&gaid=ea3de793-0df3-45b9-b31b-3fd8a407132e&mnc=&mcc=&network_type=9&network_str=&language=es-US&timezone=EST&ua=Mozilla%252F5.0%2B%2528Linux%253B%2BAndroid%2B11%253B%2BSM-A022M%2BBuild%252FRP1A.200720.012%253B%2Bwv%2529%2BAppleWebKit%252F537.36%2B%2528KHTML%252C%2Blike%2BGecko%2529%2BVersion%252F4.0%2BChrome%252F99.0.4844.58%2BMobile%2BSafari%252F537.36&gp_version=29.5.14-21%2B%255B0%255D%2B%255BPR%255D%2B430999422&sdk_version=MAL_16.0.17&app_version_name=2.21.24.22&orientation=1&screen_size=720x1600&dvi=4BztYrxBYFQ3%2BFQ3RUE0fk5QHal9GniAHnQUf7V2DZRsRr2tYg5rDkfTJ%2BzQh0R1RgftY%2Bf2Yrh0WozUhdVBRUE0D%2BzwHkc0LZRsRgxtHbi0G0zBHkeQD%2BfQWkwQ4%2Bi0Woz2hF5BRUE0HdSuR0M0hrc3LkI0G0zSiaRBn55o5nzo5VsTWjjMiUzf5Vz5i3z5ZAN0Woz0YFKTY7KtH75BRUE0NnvBi325NQVBNQ5WfoRsRrtthrxbD%2BzQRUE0Y%2BNFfAiPR0M0L7KAJoR1RURexjEbxavAR0M0DFK3HkPtYkV0G0zZxVM0WozuYrfBHk2QYgxtYoR1R3jMiUzf5Vz5i3z5ZAN0WozAH%2BzuDkM0G0z2Yrw%2FYbJ%2FR0M0H7QAh7et4ZR1RQzNiVj%2FiUvMfARMWUvei0PSiaRBn55o5nzo5VsTR0M0DrKthrN0G0z8iAQTJUc6DgfM%2BbxuJ7c%2F%2BFttY%2BfTH%2BR0WozT4%2BSQRUE0J%2BfQh0RsRgf2hdSXhgN0G0zthr2QDkzuW%2BDbDZRsRgzQY75thFV0G0ReiZRsRgf3LFQ%2FJoR1iAvsRrztJdxQhgQAJ7cTJ%2Bi0GUjsRrzthF5XhBR1iZM0L%2BiBfjl0G0Rei0RsRgfQYgfXh0R1injsRrQwHZR1RrwQ4kzBYFc3ijxuDbxtH7Ilh7KBRdHX40S3HZSdYFKgY7VlLF5PDgzXDkNe575UY7c3YBSnDk2AJkPgRoRsRgSEYFPQJdQMHZR1WnRsRgxXJ7cshrcwRUE0iTJoR0M0J7KTDkewHk2Xhg30G0RAWUvAMqSdN0zK&app_id=160670&unit_id=&data=key%253D2000088%2526state%253D1%2526network_type%253D9&m_sdk=msdk&channel=&band_width=0&open=0&country_code=PA"}`)
	// header := make([]byte, 9)
	// h := protocol.NewHeader(header)
	// h.SetMessageType('4')
	// h.SetMessageNumber(uint16(1314))
	// h.SetServerRouter(uint16(1))
	// h.SetMessageLength(uint32(len(body.Bytes())))
	// fmt.Printf("header: %s\n", header)

	fmt.Println("============= 构造消息 ===========")
	var buf []byte
	packet, _ := codec.Encode(body, protocol.MsgZipData, uint16(1314), uint16(1))
	buf = append(buf, packet...)
	fmt.Println("============= 客户端发送消息 ===========")
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
	time.Sleep(1000 * time.Millisecond)
	// for i, req := range requests {
	// 	rsp, err := codec.Unpack(respPacket[i*packetLen:])
	// 	logErr(err)
	// 	if !bytes.Equal(req, rsp) {
	// 		logging.Fatalf("request and response mismatch, conn=%s, packet size: %d, batch: %d",
	// 			c.LocalAddr().String(), packetSize, batch)
	// 	}
	// }
}
