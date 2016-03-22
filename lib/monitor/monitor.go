package monitor

import "github.com/golang/protobuf/proto"
import "sync/atomic"
import "fmt"
import "time"
import "net"
import "github.com/fakewechat/message"

type Monitor struct {
	data       [1024][2]int64
	index      int64
	mymaps     map[string]int64
	mymaps2    map[int64]string
	servertype string
	servername string
}

var GMonitor Monitor

const (
	Reportfrequency = 10 // report every 10 second
)

func (m *Monitor) Init() {

	m.mymaps = make(map[string]int64)
	m.mymaps2 = make(map[int64]string)

	m.InitNormalResign()

}

func (m *Monitor) InitServer(servertype string, servername string) {
	m.servertype = servertype
	m.servername = servername
}

func (m *Monitor) InitNormalResign() {

	m.Regist("RecvMessage")    // 获取的请求数目
	m.Regist("SendMessage")    // 发送的请求数目
	m.Regist("Connections")    // 建立的连接数
	m.Regist("Disconnections") // 断开的连接数
	m.Regist("RecvBytes")      // 接收字节数
	m.Regist("SendBytes")      // 发送字节数

	m.Regist("GateWayRequestFromClient") // 从客户端获取的请求 gateway only
	m.Regist("RPCRequestToServer")       // 从服务器端获取的请求

	m.Regist("RPCResponseFromServerSucc") // 请求服务器成功

	m.Regist("RPCResponseFromServerFailedSum")       // 请求服务器失败
	m.Regist("RPCResponseFromServerFailedNetIsue")   // 请求服务器失败
	m.Regist("RPCResponseFromServerFailedQueueFull") // 请求服务器失败 服务队列满了

	m.Regist("DropedFromQueueTimeout") // 在队列里时间太久，被直接丢掉

	m.Regist("GateWay200") // gateway 200

	m.Regist("GateWay503") // gateway local queue is full
	m.Regist("GateWay404") // gateway not found
	m.Regist("GateWay302") // gateway 301/302 redirect

	m.Regist("Router302") // gateway 301/302 redirect

	m.Regist("GateWayUserInfo")    // gateway 200
	m.Regist("GateWayChatMessage") // gateway 200

	// poster

	m.Regist("PostRequest")

	m.Regist("PostFowardPost")      // forward to another poster
	m.Regist("PostFowardLocal")     // forward to localposter
	m.Regist("DropedFromQueueFull") // 在队列里时间太多，被直接丢掉

	// localposter
	m.Regist("LocalPostRequest")
	m.Regist("LocalFowardPost") // forward to poster

	m.Regist("DBconflict") // redis write conflict
	m.Regist("DBload")     // redis write conflict
	m.Regist("DBwriteSuccess")
	m.Regist("CLientSendRequest")      // gateway to localposter
	m.Regist("LocalPosterRecvRequest") //  localposter get message

	m.Regist("CLientSendRequest_time")
	m.Regist("LocalPosterRecvRequest_time")
	m.Regist("LocalPosterRecvAck_time")
	m.Regist("UserFromRedis_time")
	m.Regist("UserToRedis_time")
	m.Regist("MessageToRedis_time")
	m.Regist("MessageFromRedis_time")
	m.Regist("GetUserInfo_time")

	m.Regist("UpdateUserAndInbox_time")
	m.Regist("UpdateUserAndOutbox_time")
	m.Regist("GetRequest_time")
	m.Regist("StoreRequest_time")
	m.Regist("rpc_send_time")

	m.Regist("CLientSendRequest")
	m.Regist("LocalPosterRecvRequest")
	m.Regist("LocalPosterRecvAck")
	m.Regist("UserFromRedis")
	m.Regist("UserToRedis")
	m.Regist("MessageToRedis")
	m.Regist("MessageFromRedis")
	m.Regist("GetUserInfo")
	m.Regist("UpdateUserAndInbox")
	m.Regist("UpdateUserAndOutbox")
	m.Regist("GetRequest")
	m.Regist("StoreRequest")
	m.Regist("rpc_send")

}

func (m *Monitor) Report(host string) {

	var socket *net.UDPConn
	for {
		var err error
		addr, err := net.ResolveUDPAddr("udp4", host)
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

	defer socket.Close()
	fmt.Println("begin report")
	for {

		tnow := time.Now().Second()
		sleep := Reportfrequency - tnow%Reportfrequency
		//fmt.Println("sleep ", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)

		tnow = time.Now().Second()
		tnow = tnow / Reportfrequency * Reportfrequency

		var i int64

		monitor := &message.MonitorStatus{
			Info: &message.RequestInfo{

				Snid:         0,
				Connectionid: 0,
				Index:        0,
			},
		}
		for i = 0; i < m.index; i++ {
			if v, ok := m.mymaps2[i]; ok {
				//	fmt.Println(v, m.data[i][0], m.data[i][1])

				msg := &message.MonitorServerStatus{
					Timestamp:     time.Now().Unix(),
					Name:          v,
					Servername:    m.servername,
					Servertype:    m.servertype,
					Absolutevalue: m.data[i][0],
					Changesvalue:  m.data[i][0] - m.data[i][1],
				}
				m.data[i][1] = m.data[i][0]
				monitor.Status = append(monitor.Status, msg)

			}
		}

		buffer, err := proto.Marshal(monitor) //SerializeToOstream
		// 发送数据
		m.Add("SendBytes", len(buffer))
		m.Add("ResponseSend", 1)
		//fmt.Println(msg.String(), len(buffer), " ", buffer)
		_, err = socket.Write(buffer)
		if err != nil {
			fmt.Println("发送数据失败!", err)
			//return
		}

	}

}

func (m *Monitor) Regist(k string) {
	_, ok := m.mymaps[k]
	if ok {
		return
	}
	m.mymaps[k] = m.index
	m.mymaps2[m.index] = k
	m.index++
}

func (m *Monitor) Add(k string, i int) int64 {
	if v, ok := m.mymaps[k]; ok {
		return atomic.AddInt64(&m.data[v][0], int64(i))
	}
	return 0
}

func (m *Monitor) Add64(k string, i int64) int64 {
	if v, ok := m.mymaps[k]; ok {
		return atomic.AddInt64(&m.data[v][0], i)
	}
	return 0
}

func (m *Monitor) Print() {
	var i int64
	for i = 0; i < m.index; i++ {
		if v, ok := m.mymaps2[i]; ok {
			fmt.Println(v, m.data[i][0], m.data[i][1])
		}
	}
}
