package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gnet-io/gnet-examples/simple_protocol/protocol"
	"github.com/gorilla/mux"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

var mockResp []byte = []byte(`{"status":1,"msg":"success","data":{"GDPR_area":false,"aa":true,"ab_id":"4_2,7_4,8_1,42_3","activeAppStatus":2,"activeAppTime":7200,"adchoice_icon":"https://cdn-adn-https.rayjump.com/cdn-adn/v2/portal/18/12/21/17/35/5c1cb3f63ce2a.png","adchoice_link":"https://www.mintegral.com/en/privacy/","adchoice_size":"35x35","adct":864000,"alrbs":0,"apk_toast":"","atf":[],"ath":0,"atrqt":0,"ats_c":0,"awct":5400,"cc":"CN","cdai":"LdxThdi1WBK/WgfPhbxQYkeXHBPwHZKsYFh=","cdnate_cfg":{"c.gdt.qq.com":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]},"c1.gdt.qq.com":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]},"c2.gdt.qq.com":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]},"c3.gdt.qq.com":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]},"c4.gdt.qq.com":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]},"event.inmobi.cn":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]},"gdt.qq.com":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]},"sc.gdt.qq.com":{"height":["__HEIGHT__"],"width":["__WIDTH__"],"x":["__DOWN_X__","__UP_X__"],"y":["__DOWN_Y__","__UP_Y__"]}},"cdt":3,"cfb":true,"cfc":1,"cnt":0,"confirm_c_play":"继续试玩","confirm_c_rv":"继续观看","confirm_description":"关闭后您将不会获得任何奖励","confirm_t":"确认关闭","confirm_title":"确认关闭？","cou":0,"country_code":"CN","crash_ct":0,"csdai":"LdxThdi1WBK/WgfPhbxQYkeXHBPwHZKAJ7eXHM==","csl":0,"cspn":0,"csw":1,"ct":120,"cud":0,"cudl":"k2T=","dlapk":1,"dlrf":1,"dlrfct":604800,"ercd":[-1,-10,-1201,-1202,-1203,-1205,-1206,-1208,-1301,-1302,-1305,-1306,-1307,-1915,10602,10603,10604,10609,10610,10616],"fbk_swt":1,"getpf":86400,"getpf_rv":43200,"hl":2,"hst":"4BzED0R1RrtTJdSAG0IX4b2ED0PBD+QqJk2MWrfXYZRsRgfTRUE0LdxThdi1WBKAH+xTLkPgWgzt4ku2Y+v/DFKwR0M0Y7h0G0zEJdxMhAEXWFc/DkePJ7QUhB2uYbi/hrcPLg5whoPUYFT0Woz3H0R1RrtTJdSAG0IXH75THkfTWgzt4ku2Y+v/DFKwWF2th73XHrQ/HoRsRrxBRUE0LdxThdi1WBK3H+xQDbN/hrcPLg5whoPUYFTXYkcMLZKBH+f2YdN0WozMY7uARUE0LdxThdi1WBKMY7cPWkR/hrcPLg5whoPUYFT0WozFRUE0LdxThdi1WBK/H+N/hrcPLg5whoPUYFT06N==","http_type":1,"hv":[],"ils":2,"is_startup_crashsystem":1,"iseu":0,"iupdid":0,"jt":[{"domain":"app.appsflyer.com","format":""},{"domain":"app.adjust.com","format":""},{"domain":"app.adjust.io","format":""},{"domain":"control.kochava.com","format":""},{"domain":"url.haloapps.com","format":""},{"domain":"ad.apsalar.com","format":""},{"domain":"td.lenzmx.com","format":""},{"domain":"cd.lenzmx.com","format":""},{"domain":"rd.lenzmx.com","format":""},{"domain":"tracking.lenzmx.com","format":""},{"domain":"measurementapi.com","format":""},{"domain":"app-adforce","format":""},{"domain":"uri6.com","format":""},{"domain":"lnk8.cn","format":""},{"domain":"lnk0.com","format":""},{"domain":"sc.adsensor.org","format":""},{"domain":"onelink.me","format":""},{"domain":"42trck.com","format":""},{"domain":"wadogo.go2cloud.org","format":""},{"domain":"tlnk.io","format":""},{"domain":"hastrk1.com","format":""},{"domain":"hastrk2.com","format":""},{"domain":"hastrk3.com","format":""},{"domain":"api-01.com","format":""},{"domain":"api-02.com","format":""},{"domain":"api-03.com","format":""},{"domain":"api-04.com","format":""},{"domain":"api-05.com","format":""}],"jumpt":1,"jumpt_bw":0,"k":0,"k_a":1,"k_j":1,"k_s":1,"l":0,"lqcnt":30,"lqpt":"GnibfM==","lqswt":0,"lqto":5,"lqtype":2,"mapping_cache_rate":0,"mcs":150,"mcto":10,"n2":86400,"n3":86400,"n4":1800,"offercacheRate":0,"offercachepacing":21600,"omsdkjs_url":"https://cdn-adn-https.rayjump.com/cdn-adn/v2/portal/19/08/20/11/06/5d5b63cb457e2.js","opent":1,"pcct":43200,"pcrn":100,"pcto":20,"platform_logo":"","platform_name":"Mintegral","plct":3600,"plctb":3600,"protect":0,"publisher_id":5731,"pw":"k2T=","refactor_switch":[{"2":true}],"retryoffer":0,"ruct":5400,"rurl":false,"sc":0,"sdk_info":"MI_7.0.6","sfct":1800,"sfzg":0,"skreld_tm":3,"spct":5400,"storekit":1,"t_vba":[],"tcct":21600,"tcto":10,"tokencachetime":604800,"uct":5400,"ujds":true,"up_tips":0,"up_tips_url":"https://hybird.rayjump.com/rv/authoriztion.html","upaid":1,"upal":86400,"updevid":1,"upgd":1,"uplc":0,"upmi":0,"upsrl":1,"useexpriedcacheoffer":2,"version_type":2,"wcus":1,"web_env_url":"","wicon":1,"wreq":2}}`)

type simpleServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	network   string
	addr      string
	multicore bool
	connected int64
	handler   http.Handler
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
	// codec := c.Context().(*protocol.SimpleCodec)
	var (
		packets [][]byte
		//heartbeatResp []byte
	)

	for {
		// decode receive data
		fmt.Println("============= 服务端循环处理收到的数据 ===========")
		// data, err := codec.Decode(c)
		// logging.Infof("data length: %d", len(data))
		// if err == protocol.ErrIncompletePacket {
		// 	fmt.Println("============= 服务端退出循环 ===========")
		// 	break
		// }
		// if err != nil {
		// 	fmt.Println("============= 服务端处理数据有错误 ===========")
		// 	logging.Errorf("invalid packet: %v", err)
		// 	return gnet.Close
		// }
		// TODO
		// receive data conversion to http request
		// http response write to packet
		// packet, _ := codec.Encode(data[:protocol.HeaderSize], mockResp)
		// var resp http.ResponseWriter
		req, _ := http.NewRequest(http.MethodGet, "/foo", nil)
		// w := httptest.NewRecorder()
		w := NewCustomResponseWriter()
		s.handler.ServeHTTP(w, req)
		// b, _ := ioutil.ReadAll(req.Response.Body)
		fmt.Println("foo handler response ", w.Result())
		// req, _ = http.NewRequest(http.MethodGet, "/bar", nil)
		// s.handler.ServeHTTP(http.Response, req)
		// b, _ = ioutil.ReadAll(req.Response.Body)
		// fmt.Println("bar handler response ", b)
		//if data[1] == protocol.MsgHeartBeat {
		//	fmt.Printf("%s 心跳包 %s\n", c.RemoteAddr().String(), c.LocalAddr().String())
		//	packets = append(packets, heartbeatResp)
		//} else {
		packet := []byte(`{"foo":"bar"}`)
		packets = append(packets, packet)
		//}
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

func fooHandler(w http.ResponseWriter, r *http.Request) {
	// r.Response.Write([]byte("Gorilla Foo!\n"))
	w.Write([]byte("Gorilla Foo!\n"))
	w.WriteHeader(http.StatusOK)
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	// r.Response.Write([]byte("Gorilla Bar!\n"))
	w.Write([]byte("Gorilla Bar!\n"))
	w.WriteHeader(http.StatusOK)
}

type CustomResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func NewCustomResponseWriter() *CustomResponseWriter {
	return &CustomResponseWriter{
		header: http.Header{},
	}
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.header
}

func (w *CustomResponseWriter) Result() []byte {
	return w.body
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	// implement it as per your requirement
	return 0, nil
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func main() {
	var port int
	var multicore bool

	// Example command: go run server.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore=true")
	flag.Parse()

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/foo", fooHandler)
	r.HandleFunc("/bar", barHandler)

	ss := &simpleServer{
		network:   "tcp",
		addr:      fmt.Sprintf(":%d", port),
		multicore: multicore,
		handler:   r,
	}
	err := gnet.Run(ss, ss.network+"://"+ss.addr, gnet.WithMulticore(multicore))
	logging.Infof("server exits with error: %v", err)
}
