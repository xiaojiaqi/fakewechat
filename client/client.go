package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fakewechat/message"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type JsonUserInfo struct {
	SendId          int      `json:"SendId"`          // how many message have been send
	ReceiveId       int      `json:"ReceiveId"`       // how many message have been received
	SendAckId       int      `json:SendAckId`         // how many message has been send to your friend
	ClientReceiveId int      `json:"ClientReceiveId"` // how many message from server by client
	Friends         []string `json:"friends"`         // friend list
}

var socket *net.UDPConn

func getHttp(url string) string {
	transport := http.Transport{
		DisableKeepAlives: true,
	}

	client := http.Client{
		Transport: &transport,
	}
	//fmt.Println(url)
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func toJson(jsonStr string) *JsonUserInfo {
	o := &JsonUserInfo{}
	if err := json.Unmarshal([]byte(jsonStr), o); err != nil {
		fmt.Println(err)
		return nil
	}
	return o
}

func getUserInfo(id int, host string, port string) *JsonUserInfo {
	url := "http://" + host + ":" + port + "/api/user/" + strconv.Itoa(id)

	r := getHttp(url)
	if r == "" {
		return nil
	}

	o := toJson(r)
	return o
}

func getRequestUrl(userid int, host string, port string) (maxsendid int, sendmessageurllist []urlmessage) {

	var o *JsonUserInfo
	for {
		o = getUserInfo(userid, host, port)
		if o == nil {
			time.Sleep(20 * time.Second)
			continue
		}
		break
	}

	// the send messageid
	sendmessageid := 0

	sendmessageurllist = make([]urlmessage, 0)
	urllist := make([]urlmessage, 0)

	for i := range o.Friends {

		sendmessageid, urllist = send10(strconv.Itoa(userid), o.Friends[i], sendmessageid, host, port)
		sendmessageurllist = append(sendmessageurllist, urllist...)
	}

	randomshuffle(sendmessageurllist)

	return sendmessageid, sendmessageurllist
}

func send10(userid string, friendid string, sendmessageid int, host string, port string) (maxsendid int, list []urlmessage) {
	list = make([]urlmessage, 0)
	for i := 1; i <= 5; i++ {
		sendmessageid += 1
		url := "http://" + host + ":" + port + "/api/c/" + userid + "/" + friendid + "/" + strconv.Itoa(sendmessageid) + "/?message=" + strconv.Itoa(i)

		req := urlmessage{}
		req.url = url
		req.messageid = sendmessageid
		list = append(list, req)

	}
	maxsendid = sendmessageid
	return maxsendid, list
}

func randomshuffle(list []urlmessage) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	length := len(list)
	for id := 0; id < length; id++ {
		newid := r1.Intn(int(length) - 1)
		orgid := list[id]
		list[id] = list[newid]
		list[newid] = orgid
	}

}

func SendRequest(id int, list []urlmessage, messageid int, host string, port string) {
	sendsucc := false
	SendAckId := 0

	for {

		// send all request
		for i := range list {

			if list[i].messageid >= SendAckId {

				getHttp(list[i].url)
			}
		}

		for i := 0; i < 10; i++ {
			o := getUserInfo(id, host, port)

			if o != nil {

				//fmt.Println("userid =", id, o.SendId, o.ReceiveId, o.SendAckId, messageid)
				SendAckId = o.SendAckId
				if o.SendId == o.SendAckId && o.SendAckId == messageid {
					sendsucc = true
					lock.Lock()
					send += 1

				        SendStatus(send)
					lock.Unlock()
					fmt.Println(send, "user  send")
					break
				}
			}

			time.Sleep(5 * time.Second)
		}

		if sendsucc == true {
			break
		}
	}

}

type urlRequest struct {
	host   string
	port   string
	userid int
}

type urlmessage struct {
	url       string
	messageid int
}

var host *string
var port *string
var minid *int
var maxid *int
var monitor *string

const (
	Concurrency = 300
)

func makeRequest(ch *chan urlRequest) {

	lock.Lock()
	requests = *maxid - *minid + 1
	lock.Unlock()

	for i := *minid; i <= *maxid; i++ {
		fmt.Println(i)
		a := urlRequest{}
		a.host = *host
		a.port = *port
		a.userid = i
		*ch <- a
	}

}

var lock sync.Mutex
var processed int
var requests int
var finished int
var send int

func Process(channel *chan urlRequest) {
	for {
		//   I know, it is a bug maybe dead lock
		//
		a := <-*channel
		fmt.Println("get ", a.userid)
		maxsendid, list := getRequestUrl(a.userid, a.host, a.port)
		SendRequest(a.userid, list, maxsendid, a.host, a.port)
		lock.Lock()
		processed += 1
		lock.Unlock()
	}
}

func CheckUserData(channel *chan urlRequest) {
	checked := 0
	for {
		a := <-*channel
		for {
			o := getUserInfo(a.userid, a.host, a.port)
			if o == nil {
				fmt.Println("o == nil")
				time.Sleep(10 * time.Second)
				continue
			}

			if (o.SendId == o.ReceiveId) && (o.SendId == o.SendAckId) {
				fmt.Println(a.userid, o.SendId, o.SendAckId, o.ReceiveId)
				lock.Lock()
				finished += 1
				fmt.Println(finished, "user finished", a.userid)


				lock.Unlock()
				break
			} else {
				fmt.Println("check user failed", a.userid, o.SendId, o.SendAckId, o.ReceiveId)
				time.Sleep(10 * time.Second)
				//panic("error")
			}

		}
		checked += 1

		if requests == checked {
			break
		}

	}
}

func SendStatus(finished int) {
	monitor := &message.ClientMonitorStatus{}
	monitor.Host = *host +"_" + *port

	monitor.Finished = int32(finished)

	buffer, err := proto.Marshal(monitor) //SerializeToOstream
	// 发送数据
	//fmt.Println(msg.String(), len(buffer), " ", buffer)
	_, err = socket.Write(buffer)
	if err != nil {
		fmt.Println("发送数据失败!", err)
		//return
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	host = flag.String("host", "127.0.0.1", "host")
	port = flag.String("port", "9501", "port")
	minid = flag.Int("minid", 1, "minid")
	maxid = flag.Int("maxid", 2500, "maxid")
	monitor = flag.String("monitor", "127.0.0.1:8002", "monitor")

	flag.Parse()

	for {
		var err error
		addr, err := net.ResolveUDPAddr("udp4", *monitor)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}

		socket, err = net.DialUDP("udp", nil, addr)

		if err != nil {
			fmt.Println("连接失败!", err)
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}

		// connection success
		break
	}


	t := time.Now().Unix()
	var channel chan urlRequest
	var checkchannel chan urlRequest
	channel = make(chan urlRequest, 999)
	checkchannel = make(chan urlRequest, 999)
	go makeRequest(&channel)
	for i := 0; i <= Concurrency; i++ {
		go Process(&channel)
	}

	for {
		time.Sleep(1 * time.Second)
		lock.Lock()
		fmt.Println("processed == requests", processed, requests)
		if processed == requests {

			lock.Unlock()
			go makeRequest(&checkchannel)
			break
		}
		lock.Unlock()
	}

	CheckUserData(&checkchannel)
	fmt.Println("spend ", time.Now().Unix()-t)
}
