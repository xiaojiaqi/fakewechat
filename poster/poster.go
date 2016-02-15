package main

import (
	"encoding/binary"
	"fmt"
	. "github.com/fakewechat/lib/clientpool"
	. "github.com/fakewechat/lib/config"
	. "github.com/fakewechat/lib/contstant"
	"github.com/fakewechat/lib/flags"
	. "github.com/fakewechat/lib/monitor"
	. "github.com/fakewechat/lib/postrequest"
	. "github.com/fakewechat/lib/rg"
	"github.com/fakewechat/lib/serverstatus"
	. "github.com/fakewechat/lib/updateserverstatus"
	. "github.com/fakewechat/lib/version"
	. "github.com/fakewechat/message"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var channel [PostChannelSize]chan *GeneralMessage

type PosterAPI int

var counter uint32

func (t *PosterAPI) PosterMessage(Req *GeneralMessage, id *uint64) error {
	fmt.Println("(t *PosterAPI) PosterMessage(Req *GeneralMessage, id *uint64)")
	s := atomic.AddUint32(&counter, 1)
	channelId := s % PostWorkThreadSize
	err := PushtoQueue(&channel[channelId], Req)
	if err != nil {

		fmt.Println("err := PushtoQueue(&channel[channelId], Req) failed")
		GMonitor.Add("DropedFromQueueFull", 1)
		return err
	}
	GMonitor.Add("PostRequest", 1)
	fmt.Println("PushtoQueue(&channel[channelId], Req) success")

	return nil
}

func ProcessRequest(channel *chan *GeneralMessage, Channelindex int) {
	for {
		L := GetSomeMessageFromQueue(channel, PostThreadQueueLength)
		for i := range L {
			// split user
			DispatchToServer(L[i], Channelindex)
		}
	}
}

func DispatchToServer(req *GeneralMessage, Channelindex int) {

	if req.MessageType == CHAT_CLIENT_TO_LOCALPOST {
		panic("get request  CHAT_CLIENT_TO_LOCALPOST")
	} else if req.MessageType == CHAT_LOCALPOST_TO_LOCALPOST || req.MessageType == CHAT_LOCALPOST_ACK {
		forwardChatMessage(req, Channelindex)
	} else {
		panic("unsupport request " + strconv.Itoa(int(req.MessageType)))
	}

}

func forwardChatMessage(req *GeneralMessage, Channelindex int) {

	receiverId := req.ReceiverId
	rgid := GetRg(receiverId)
	var client ClientInterface
	if rgid != int(*flags.RgId) {
		// send to other rg
		client = GloabalClientPool[Channelindex].GetClient(POSTERSERVER, rgid)
	} else {
		// to same rg
		client = GloabalClientPool[Channelindex].GetClient(LOCALPOSTSERVER, rgid)
	}

	if client == nil {
		fmt.Println("no client found xxx! send")
		return
	}

	defer GloabalClientPool[Channelindex].ReturnClient(client)
	fmt.Println("forwardChatMessage index ", Channelindex)
	rpcclient := client.GetRPCClient()
	if rpcclient == nil {
		panic("rpcclient is nil")
	}
	var ids uint64
	if rgid != int(*flags.RgId) {
		// send to other rg
		GMonitor.Add("PostFowardPost", 1)
		err := rpcclient.Call("PosterAPI.PosterMessage", req, &ids)
		if err != nil {
			fmt.Println("arith error:", err)
			return
		}

	} else {
		// to same rg
		GMonitor.Add("PostFowardLocal", 1)
		err := rpcclient.Call("LocalPosterAPI.PosterMessage", req, &ids)
		if err != nil {
			fmt.Println("arith error:", err)
			return
		}

	}

}

func GetLogBuff(list []*GeneralMessage) []byte {

	logbuf := make([]byte, 0)

	for i := range list {
		s := list[i].String()
		buffLen := len(s)

		header := make([]byte, 4)
		binary.BigEndian.PutUint16(header, uint16(buffLen))
		binary.BigEndian.PutUint16(header[2:], uint16(0x01))
		logbuf = append(logbuf, header...)
		logbuf = append(logbuf, []byte(s)...)
	}
	return logbuf
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

	posterAPI := new(PosterAPI)
	rpc.Register(posterAPI)

	tcpAddr, err := net.ResolveTCPAddr("tcp", *(flags.ListenAddress)+":"+strconv.Itoa(int(*(flags.ListenPort))))
	checkError(err)

	for i := 0; i < PostChannelSize; i++ {
		channel[i] = make(chan *GeneralMessage, PostChannelLength)
	}
	for i := 0; i < PostWorkThreadSize; i++ {
		go ProcessRequest(&channel[i], i)
		go CheckServerStatusUpdate(i, *flags.ServerType, int(*flags.RgId))
	}
	{
		servertype := make([]string, 0)
		rg := make([]int, 0)
		rgid := int(*flags.RgId)
		servertype = append(servertype, LOCALPOSTSERVER)
		rg = append(rg, rgid)

		servertype = append(servertype, POSTERSERVER)
		rg = append(rg, -1)

		go serverstatus.SyncServerStatusChanged(*flags.RouterServerURL, servertype, rg)
		go serverstatus.RegistServer(*flags.RouterRegistURL, *flags.ServerType, *flags.HostId, int(*flags.RgId), *flags.ListenAddress, *flags.ListenPort, *flags.CellId)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {

			continue
		}

		fmt.Println("local  ", conn.RemoteAddr().Network(), " ", conn.RemoteAddr().String())
		go rpcProcess(conn)
	}
}

func rpcProcess(conn net.Conn) {
	defer conn.Close()
	rpc.ServeConn(conn)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
