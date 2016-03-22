package main

import (
	"fmt"
	. "github.com/fakewechat/lib/clientpool"
	. "github.com/fakewechat/lib/config"
	. "github.com/fakewechat/lib/contstant"
	"github.com/fakewechat/lib/flags"
	. "github.com/fakewechat/lib/monitor"
	. "github.com/fakewechat/lib/postrequest"
	"github.com/fakewechat/lib/serverstatus"
	. "github.com/fakewechat/lib/updateserverstatus"
	. "github.com/fakewechat/localposter/core"
	. "github.com/fakewechat/localposter/handler"
	. "github.com/fakewechat/message"

	. "github.com/fakewechat/lib/version"
	"math/rand"

	"net"
	"net/rpc"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var counter uint32

type LocalPosterAPI int

var channel [LocalPostChannelSize]chan *GeneralMessage

func (t *LocalPosterAPI) PosterMessage(Req *GeneralMessage, id *uint64) error {
	s := atomic.AddUint32(&counter, 1)
	channelId := s % LocalPostChannelSize
	err := PushtoQueue(&channel[channelId], Req)
	if err != nil {
		GMonitor.Add("DropedFromQueueFull", 1)
		return err
	}
	GMonitor.Add("LocalPostRequest", 1)
	return nil
}

func ProcessRequest(channel *chan *GeneralMessage, Channelindex int) {
	for {
		L := GetSomeMessageFromQueue(channel, LocalPostThreadQueueLength)
		for i := range L {
			// split user
			DispatchToServer(L[i], Channelindex)
		}
	}
}

func DispatchToServer(req *GeneralMessage, Channelindex int) {

	GMonitor.Add("LocalFowardPost", 1)
	if req.MessageType == CHAT_CLIENT_TO_LOCALPOST {
		ProcessClient_to_Local(req, Channelindex)
	} else if req.MessageType == CHAT_LOCALPOST_TO_LOCALPOST {
		ProcessLocal_to_Local(req, Channelindex)
	} else if req.MessageType == CHAT_LOCALPOST_ACK {
		ProcessLocal_ack(req, Channelindex)
	} else {
		panic("unsupport request " + strconv.Itoa(int(req.MessageType)))
	}
}

func ProcessClient_to_Local(req *GeneralMessage, Channelindex int) {

	GMonitor.Add("CLientSendRequest", 1)
	Redisclient := GloabalClientPool[Channelindex].GetClient(REDISSERVER, int(*flags.RgId))
	defer GloabalClientPool[Channelindex].ReturnClient(Redisclient)

	Rpcclient := GloabalClientPool[Channelindex].GetClient(POSTERSERVER, int(*flags.RgId))
	defer GloabalClientPool[Channelindex].ReturnClient(Rpcclient)

	var handler ProcessHandler
	localhandler := &LocalPosterHandler{}
	localhandler.Setup(Redisclient, Rpcclient)
	handler = localhandler

	ProcessClient_to_Local_Message(handler, req)
}

func ProcessLocal_to_Local(req *GeneralMessage, Channelindex int) {

	GMonitor.Add("LocalPosterRecvRequest", 1)
	Redisclient := GloabalClientPool[Channelindex].GetClient(REDISSERVER, int(*flags.RgId))
	defer GloabalClientPool[Channelindex].ReturnClient(Redisclient)

	Rpcclient := GloabalClientPool[Channelindex].GetClient(POSTERSERVER, int(*flags.RgId))
	defer GloabalClientPool[Channelindex].ReturnClient(Rpcclient)

	var handler ProcessHandler
	localhandler := &LocalPosterHandler{}
	localhandler.Setup(Redisclient, Rpcclient)
	handler = localhandler
	ProcessLocal_to_Local_Message(handler, req)
}

func ProcessLocal_ack(req *GeneralMessage, Channelindex int) {

	GMonitor.Add("LocalPosterRecvAck", 1)
	Redisclient := GloabalClientPool[Channelindex].GetClient(REDISSERVER, int(*flags.RgId))
	defer GloabalClientPool[Channelindex].ReturnClient(Redisclient)

	var handler ProcessHandler
	localhandler := &LocalPosterHandler{}
	// no rpc needed
	localhandler.Setup(Redisclient, nil)
	handler = localhandler
	ProcessLocal_ack_Message(handler, req)
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

	localposterAPI := new(LocalPosterAPI)
	rpc.Register(localposterAPI)

	tcpAddr, err := net.ResolveTCPAddr("tcp", *(flags.ListenAddress)+":"+strconv.Itoa(int(*(flags.ListenPort))))
	checkError(err)

	for i := 0; i < LocalPostChannelSize; i++ {
		channel[i] = make(chan *GeneralMessage, LocalPostChannelLength)
	}
	for i := 0; i < LocalPostChannelSize; i++ {
		go ProcessRequest(&channel[i], i)
		go CheckServerStatusUpdate(i, *flags.ServerType, int(*flags.RgId))
	}

	{
		servertype := make([]string, 0)
		rg := make([]int, 0)
		rgid := int(*flags.RgId)
		servertype = append(servertype, POSTERSERVER)
		rg = append(rg, rgid)

		servertype = append(servertype, REDISSERVER)
		rg = append(rg, rgid)

		servertype = append(servertype, MEMCACHEDSERVER)
		rg = append(rg, rgid)
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
