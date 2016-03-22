package main

import (
	//"encoding/binary"
	"fmt"

	. "github.com/fakewechat/lib/config"
	. "github.com/fakewechat/lib/contstant"

	"github.com/fakewechat/lib/flags"

	. "github.com/fakewechat/lib/monitor"
	//. "github.com/fakewechat/lib/postrequest"
	//"github.com/fakewechat/lib/rpc"
	. "github.com/fakewechat/lib/clientpool"
	. "github.com/fakewechat/lib/rg"
	"github.com/fakewechat/lib/serverstatus"
	. "github.com/fakewechat/lib/updateserverstatus"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/lib/version"
	. "github.com/fakewechat/message"
	"github.com/garyburd/redigo/redis"
	"net/rpc"

	//"github.com/golang/protobuf/proto"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var counter uint32

type CacheAPI int

var mutex [CacheServerChannelSize]sync.Mutex

func (t *CacheAPI) GetUserInfo(id *uint64, u *UserInfor) error {

	//fmt.Println("GetUserInfo")
	s := atomic.AddUint32(&counter, 1)
	Channelindex := s % CacheServerChannelSize
	mutex[Channelindex].Lock()
	defer mutex[Channelindex].Unlock()

	pool := GloabalClientPool[Channelindex]
	Rgid := GetRg(*id)
	client := pool.GetClient("redis", Rgid)
	defer pool.ReturnClient(client)
	if client == nil {
		panic("no redis is ready")
	}

	conn := (redis.Conn)(*client.GetRedisClient())

	userkey := GetUserInfoName(*id)

	n, err := conn.Do("GET", userkey)
	if err != nil {
		fmt.Println(err)

	}

	*u = *UserFromRedis(n.([]byte))

	return nil

}

func main() {
	// print version
	PrintVersion()

	runtime.GOMAXPROCS(runtime.NumCPU())

	flags.Parse()

	GMonitor.Init()
	GMonitor.InitServer(*flags.ServerType, *flags.HostId)
	go GMonitor.Report(*flags.MonitorHost)

	InitClientPool()

	rand.Seed(time.Now().UTC().UnixNano())

	cacheAPI := new(CacheAPI)
	rpc.Register(cacheAPI)

	tcpAddr, err := net.ResolveTCPAddr("tcp", *(flags.ListenAddress)+":"+strconv.Itoa(int(*(flags.ListenPort))))
	checkError(err)

	for i := 0; i < CacheThreadSize; i++ {
		go CheckServerStatusUpdate(i, *flags.ServerType, int(*flags.RgId))
	}

	{
		servertype := make([]string, 0)
		rg := make([]int, 0)
		rgid := int(*flags.RgId)

		servertype = append(servertype, REDISSERVER)
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
