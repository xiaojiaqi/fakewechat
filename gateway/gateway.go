package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/fakewechat/lib/clientpool"
	. "github.com/fakewechat/lib/config"
	. "github.com/fakewechat/lib/contstant"
	"github.com/fakewechat/lib/flags"
	. "github.com/fakewechat/lib/monitor"
	"github.com/fakewechat/lib/net/http"
	. "github.com/fakewechat/lib/rg"
	"github.com/fakewechat/lib/serverstatus"
	. "github.com/fakewechat/lib/sid"
	. "github.com/fakewechat/lib/updateserverstatus"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/lib/version"
	. "github.com/fakewechat/message"
	//"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type JsonUserInfo struct {
	SendId          int      `json:"SendId"`          // how many message have been send
	ReceiveId       int      `json:"ReceiveId"`       // how many message have been received
	SendAckId       int      `json:SendAckId`         // how many message has been send to your friend
	ClientReceiveId int      `json:"ClientReceiveId"` // how many message from server by client
	Friends         []string `json:"friends"`         // friend list
}

type HttpContext struct {
	response           *http.ResponseWriter
	request            *http.Request
	requestType        int
	getRequestUnixTime uint64
	conn               net.Conn

	bufferIO *bufio.ReadWriter
	//	result 				chan HttpContext
}

var channel [GateWayChannelSize]chan *HttpContext
var sleeptime time.Duration
var monitor Monitor
var sid Sid

const (
	userinfoRequest     = 0
	friendinboxRequest  = 1
	friendoutboxRequest = 2

	chatinboxRequest  = 3
	chatoutboxRequest = 4

	postfriendRequest      = 5
	postchatmessageRequest = 6
)

type route struct {
	pattern *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request)
}

type RegexpHandler struct {
	routes []*route
}

func getReqdate(req *http.Request, key string) string {
	value := ""
	fmt.Println(req.Form)
	if len(req.Form[key]) > 0 {
		value = req.Form[key][0]
		fmt.Println(value)
	}
	return value
}

func recvfromQueue(queue chan *HttpContext) (*HttpContext, error) {
	v := <-queue
	return v, nil
}

func pushtoQueue(queue *chan *HttpContext, req *HttpContext) error {
	//fmt.Println("len :", len(*queue))
	select {

	case *queue <- req:
		//fmt.Println("no full")
		return nil
	default:
		//fmt.Println("it is full")
		return errors.New("it is full")
	}
}

var serversion Version
var counter uint32

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("handle ", r.URL.Path, " ", time.Now().Format("2006-01-02 15:04:05"))
	for _, route := range h.routes {

		if route.pattern.MatchString(r.URL.Path) {
			//fmt.Println("  found", r.URL.Path, " ", time.Now().Format("2006-01-02 15:04:05"))
			route.handler(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func sendtoQue(types int, rw http.ResponseWriter, req *http.Request) {
	hj, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	netconn, bio, err := hj.Hijack()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	q := &HttpContext{
		response:    &rw,
		request:     req,
		requestType: types,
		conn:        netconn,
		bufferIO:    bio,
		//httpconn:     http.Conn(hj),
		//result: make(chan HttpContext),
	}
	s := atomic.AddUint32(&counter, 1)

	channelId := s % GateWayChannelSize
	if nil != pushtoQueue(&channel[channelId], q) {
		fmt.Println("too many request, close")
		bio.WriteString("HTTP/1.1 555 OK\r\n")
		bio.WriteString("server: gateway\r\n")
		bio.WriteString("Connection: close\r\n\r\n")
		bio.Flush()
		netconn.Close()
	}
}

func userinfo(rw http.ResponseWriter, req *http.Request) {
	sendtoQue(userinfoRequest, rw, req)
}

func friendgroupinbox(rw http.ResponseWriter, req *http.Request) {
	sendtoQue(friendinboxRequest, rw, req)
}

func friendgroupoutbox(rw http.ResponseWriter, req *http.Request) {
	sendtoQue(friendoutboxRequest, rw, req)
}

func chatmessageinbox(rw http.ResponseWriter, req *http.Request) {
	sendtoQue(chatinboxRequest, rw, req)
}

func chatmessageoutbox(rw http.ResponseWriter, req *http.Request) {
	sendtoQue(chatoutboxRequest, rw, req)
}

func postfriend(rw http.ResponseWriter, req *http.Request) {
	sendtoQue(postfriendRequest, rw, req)
}

func postchatmessage(rw http.ResponseWriter, req *http.Request) {
	sendtoQue(postchatmessageRequest, rw, req)
}

func file(rw http.ResponseWriter, req *http.Request) {

	hj, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	netconn, bio, err := hj.Hijack()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	url := req.URL
	fmt.Println(url.Path)
	b, _ := fileCache[url.Path]
	bio.WriteString("HTTP/1.1 200 OK\r\n")

	bio.WriteString("server: gateway\r\n")

	bio.WriteString(string("Content-Length: " + strconv.Itoa(len(b)) + "\r\n\r\n"))
	bio.Write(b)
	bio.Flush()

	netconn.Close()
}

func ProcessHttpRequest(Channelindex int) {
	for {
		// only one channel of gateway
		a := <-channel[Channelindex]
		GMonitor.Add("GateWayRequestFromClient", 1)
		types := a.requestType
		//fmt.Println("ProcessHttpRequest ", Channelindex)
		if types == userinfoRequest {
			processuserinfoRequest(Channelindex, a)
		} else if types == friendinboxRequest {
			processfriendinboxRequest(Channelindex, a)
		} else if types == friendoutboxRequest {
			processfriendoutboxRequest(Channelindex, a)
		} else if types == chatinboxRequest {
			processchatinboxRequest(Channelindex, a)
		} else if types == chatoutboxRequest {
			processchatoutboxRequest(Channelindex, a)
		} else if types == postfriendRequest {
			processpostfriendRequest(Channelindex, a)
		} else if types == postchatmessageRequest {
			processpostchatmessageRequest(Channelindex, a)

		}

		conn := a.conn
		bufrw := a.bufferIO

		bufrw.Flush()
		conn.Close()
		/*
			        fmt.Println("Get a call ", len(channel[0]), " ", time.Now().Format("2006-01-02 15:04:05"))
					conn := a.conn
					{
						bufrw := a.bufio
						index += 1


						bufrw.WriteString("HTTP/1.1 200 OK\r\n")
						bufrw.WriteString("server: gateway\r\n")
						sendbuff := ""
						if types == friendlistRequest {
							sendbuff = "Friend." + strconv.Itoa(index) + "\n"
						} else if types == inboxRequest {
							sendbuff = "inbox" + strconv.Itoa(index)
						} else if types == outboxRequest {
							sendbuff = "outbox" + strconv.Itoa(index)
						} else if types == postRequest {
							sendbuff = "post" + strconv.Itoa(index)
						}
						bufrw.WriteString("Content-Length: " + strconv.Itoa(len(sendbuff)) + "\r\n\r\n")
						bufrw.WriteString(sendbuff)
						//a.result <- a
						fmt.Println("write done")
						bufrw.Flush()
						conn.Close()
						fmt.Println("conn close")
					}
		*/
	}
}

/*
   userinfoRequest     = 0
   friendinboxRequest  = 1
   friendoutboxRequest = 2

   chatinboxRequest  = 3
   chatoutboxRequest = 4

   postfriendRequest      = 5
   postchatmessageRequest = 6
*/

/*
func getReqdate(req *http.Request, key string) string {
	value := ""
	if len(req.Form[key]) > 0 {
		value = req.Form[key][0]
		fmt.Println(value)
	}
	return value
}

func getReqIntdata(req *http.Request, key string) int {

	value := ""
	if len(req.Form[key]) > 0 {
		value = req.Form[key][0]
		fmt.Println(value)
	}
	return getIntValue(value)
}

func getIntValue(key string) int {

	result, err := strconv.Atoi(key)
	if err != nil {
		result = 0
	}
	return result
}
*/

type LocalPosterAPI int

func processuserinfoRequest(Channelindex int, req *HttpContext) {

	url := req.request.URL
	paths := strings.Split(url.Path, "/")
	//fmt.Println(paths, " Channelindex:", Channelindex)

	var userid uint64
	userid = uint64(GetIntValue(paths[3]))
	Rgid := GetRg(userid)
	//fmt.Println("rgid", Rgid)
	if Rgid != int(*flags.RgId) {
		return
	}
	client := GloabalClientPool[Channelindex].GetClient(CACHESERVER, Rgid)
	//fmt.Println("get client success", client)
	if client == nil {
		fmt.Println("no client found xxx! send")
		return
	}

	defer GloabalClientPool[Channelindex].ReturnClient(client)

	//fmt.Println("send index ", Channelindex)

	rpcclient := client.GetRPCClient()
	u := UserInfor{}
	//fmt.Println("send Request now ", Channelindex)
	err := rpcclient.Call("CacheAPI.GetUserInfo", &userid, &u)
	if err != nil {
		fmt.Println("arith error:", err)
		return
	}
	//fmt.Println("CacheAPI.GetUserInfo done ")
	GMonitor.Add("GateWayUserInfo", 1)

	bufrw := req.bufferIO

	bufrw.WriteString("HTTP/1.1 200 OK\r\n")
	bufrw.WriteString("server: gateway\r\n")
	sendbuff := User2Json(&u)
	bufrw.WriteString("Content-Length: " + strconv.Itoa(len(sendbuff)) + "\r\n\r\n")
	bufrw.WriteString(sendbuff)
	//a.result <- a
	fmt.Println("write done")
	GMonitor.Add("GateWay200", 1)

}

func processfriendinboxRequest(Channelindex int, req *HttpContext) {

}

func processfriendoutboxRequest(Channelindex int, req *HttpContext) {

}

func processchatinboxRequest(Channelindex int, req *HttpContext) {

}

func processchatoutboxRequest(Channelindex int, req *HttpContext) {

}

func processpostfriendRequest(Channelindex int, req *HttpContext) {

}

func processpostchatmessageRequest(Channelindex int, context *HttpContext) {

	url := context.request.URL
	paths := strings.Split(url.Path, "/")
	fmt.Println(paths, "processpostchatmessageRequest  Channelindex:", Channelindex)

	var senduserid uint64
	var recvuserid uint64
	var sendid uint64
	senduserid = uint64(GetIntValue(paths[3]))
	recvuserid = uint64(GetIntValue(paths[4]))
	sendid = uint64(GetIntValue(paths[5]))
	context.request.ParseForm()
	message := getReqdate(context.request, "message")
	//fmt.Println("message = ", message)
	sendRgid := GetRg(senduserid)

	//fmt.Println("rgid", sendRgid)
	if sendRgid != int(*flags.RgId) {
		return
	}

	client := GloabalClientPool[Channelindex].GetClient(LOCALPOSTSERVER, sendRgid)
	fmt.Println("get client success", client)
	if client == nil {
		fmt.Println("no client found xxx! send")
		return
	}

	defer GloabalClientPool[Channelindex].ReturnClient(client)

	//fmt.Println("send index ", Channelindex)

	rpcclient := client.GetRPCClient()

	postmessage := GeneralMessage{}
	postmessage.MessageType = CHAT_CLIENT_TO_LOCALPOST
	postmessage.SendId = sendid

	postmessage.SenderId = senduserid
	postmessage.ReceiverId = recvuserid

	postmessage.RequestTimeStamp = uint64(time.Now().Unix())

	chatmessage := &ChatMessage{}
	chatmessage.SId = sid.GetId()
	chatmessage.SenderId = senduserid
	chatmessage.ReceiverId = recvuserid

	chatmessage.MessageBody = message
	chatmessage.RequesTtimeStamp = uint64(time.Now().Unix())
	chatmessage.SendId = 0

	/*
		f := &FriendGroupMessage{}

		f.Sid = proto.Uint64(1011111111110) // machine sid
		f.Rootid = proto.Uint64(1011111111110)
		f.Parent = proto.Uint64(1011111111110)
		f.Senderuserid = proto.Uint64(1011111111110) // sender
		f.Recvuserid = proto.Uint64(1011111111110)   // sender
		f.Body = proto.String("we are the world")
		f.Posttime = proto.Uint64(20110)
	*/
	postmessage.Chatmessage = chatmessage

	//	u := UserInfor{}
	var u uint64
	//fmt.Println("send Request now ", Channelindex)
	err := rpcclient.Call("LocalPosterAPI.PosterMessage", &postmessage, &u)
	if err != nil {
		fmt.Println("arith error:", err)
		return
	}
	GMonitor.Add("GateWayChatMessage", 1)
	//fmt.Println("LocalPosterAPI.PosterMessage done ")
	bufrw := context.bufferIO

	bufrw.WriteString("HTTP/1.1 200 OK\r\n")
	bufrw.WriteString("server: gateway\r\n")
	sendbuff := "OK"
	bufrw.WriteString("Content-Length: " + strconv.Itoa(len(sendbuff)) + "\r\n\r\n")
	bufrw.WriteString(sendbuff)
	//a.result <- a
	fmt.Println("write done")
	GMonitor.Add("GateWay200", 1)
}

/*
func checkMon(index int) {
	var LocalVersion Version
	for {
		time.Sleep(10 * time.Second)
		if LocalVersion != serverstatus.GlobalVersion {
			LocalVersion = serverstatus.GlobalVersion
			copy := serverstatus.CopyLocalServerStatus()
			fmt.Println("get new copy", copy)
			GloabalClientPool[index].AddNew(LocalVersion, copy, *flags.HostType,  int(*flags.RgId) )
		} else {
			fmt.Println("equal")
		}
	}
}
*/

func getFilelist(path string) string {
	var strRet string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		b, _ := ioutil.ReadFile(path)
		fileCache[path[4:len(path)]] = b
		fmt.Println(path)
		strRet += path + "\r\n"
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return strRet
}

var fileCache map[string][]byte

func loadFile() {
	fileCache = make(map[string][]byte)
	getFilelist("./site/")
	for i, _ := range fileCache {
		fmt.Println(i)
	}
}

func main() {

	PrintVersion()

	runtime.GOMAXPROCS(runtime.NumCPU())

	flags.Parse()

	GMonitor.Init()
	GMonitor.InitServer(*flags.ServerType, *flags.HostId)
	go GMonitor.Report(*flags.MonitorHost)

	InitClientPool()

	rand.Seed(time.Now().UTC().UnixNano())

	// set message id
	sid.SetIndex(*flags.SnId)

	loadFile()

	sleeptime = 20

	for i := 0; i < GateWayChannelSize; i++ {
		channel[i] = make(chan *HttpContext, GateWayChannelLength)
	}
	for i := 0; i < GateWayWorkThreadSize; i++ {
		go ProcessHttpRequest(i)
		go CheckServerStatusUpdate(i, *flags.ServerType, int(*flags.RgId))
	}

	{
		servertype := make([]string, 0)
		rg := make([]int, 0)
		servertype = append(servertype, LOCALPOSTSERVER)
		rgid := int(*flags.RgId)
		rg = append(rg, rgid)
		servertype = append(servertype, CACHESERVER)
		rg = append(rg, rgid)
		go serverstatus.SyncServerStatusChanged(*flags.RouterServerURL, servertype, rg)

		go serverstatus.RegistServer(*flags.RouterRegistURL, *flags.ServerType, *flags.HostId, int(*flags.RgId), *flags.ListenAddress, *flags.ListenPort, *flags.CellId)
	}
	var regexpHandler RegexpHandler
	xp, err := regexp.Compile("/api/user/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, userinfo) // user infor
	}
	xp, err = regexp.Compile("/api/f/inbox/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, friendgroupinbox) // friend group timeline
	}
	xp, err = regexp.Compile("/api/f/outbox/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, friendgroupoutbox) // friend group personal page
	}

	xp, err = regexp.Compile("/api/c/inbox/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, chatmessageinbox) // chat message recevice
	}
	xp, err = regexp.Compile("/api/c/outbox/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, chatmessageoutbox) // chat message sended
	}
	xp, err = regexp.Compile("/api/f/[0-9]*/[0-9]*/[0-9]*/") // post a message to friend group  #api/f/$userid/$indexid/$messageid
	if err == nil {
		regexpHandler.HandleFunc(xp, postfriend) //设定访问的路径
	}

	xp, err = regexp.Compile("/api/c/[0-9]*/[0-9]*/[0-9]*/") // post a message to chat message
	if err == nil {
		regexpHandler.HandleFunc(xp, postchatmessage) //设定访问的路径
	}

	xp, err = regexp.Compile("/(css|js|images)/.*")
	if err == nil {
		regexpHandler.HandleFunc(xp, file) //设定访问的路径
	}

	xp, err = regexp.Compile("/*html")
	if err == nil {
		regexpHandler.HandleFunc(xp, file) //设定访问的路径
	}

	server := &http.Server{
		Addr:           *(flags.ListenAddress) + ":" + strconv.Itoa(int(*(flags.ListenPort))),
		Handler:        regexpHandler,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//hostid string, serverType string, rg int, host string, port uint)
	log.Fatal(server.ListenAndServe())
}

/*
type Args struct {
	A, B int
}

func send(index int) {
	for {
		//sendReq(index)
	}
}
*/

/*
func sendReq(index int) {

	fmt.Println("send ", index)
	time.Sleep(10 * time.Second)

	rpcClient := GloabalClientPool[index].GetClient("cache", 1)
	if rpcClient == nil {
		fmt.Println("no client found xxx! send")
		return
	}

	defer GloabalClientPool[index].ReturnClient(rpcClient)

	fmt.Println("send index ", index)

	client := rpcClient.GetRPCClient()
	var userid uint64
	userid = 1
	u := UserInfor{}
	err := client.Call("CacheAPI.GetUserInfo", &userid, &u)
	if err != nil {
		fmt.Println("arith error:", err)
		return
	}
	fmt.Printf("Arith: %d %d %d %v   %d\n", *u.Sendindex, *u.Recvindex, *u.Clientrecvindex, u, time.Now().Unix())
}
*/

func User2Json(u *UserInfor) string {
	var o JsonUserInfo
	/*

			SendId           int      `json:"SendId"`           // how many message have been send
		ReceiveId       int      `json:"ReceiveId"`       // how many message have been received
		SendAckId        int      `json:SendAckId`          // how many message has been send to your friend
		ClientReceiveId int      `json:"ClientReceiveId"` // how many message from server by client

	*/
	//fmt.Println("%v", *u)
	// if the data is 0, the point will been null ? bug?

	o.SendId = int(u.SendId) //

	o.ReceiveId = int(u.ReceiveId)
	o.ClientReceiveId = int(u.ClientReceiveId) //
	o.SendAckId = int(u.SendAckId)             // int

	for k, _ := range u.UserMap {
		o.Friends = append(o.Friends, strconv.Itoa(int(k)))
	}

	if b, err := json.Marshal(o); err == nil {
		//fmt.Println("================struct 到json str==")
		return string(b)
	}
	return string("")
}
