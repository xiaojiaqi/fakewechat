package main

import "time"
import "net"
import "github.com/fakewechat/message"
import "math/rand"
import "strconv"
import "fmt"
import "github.com/golang/protobuf/proto"

func main()  {
    host := "127.0.0.1:8002"
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
    s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for {

		  

		monitor := &message.ClientMonitorStatus{
			 
		}
		monitor.Host = "192.168.0." + strconv.Itoa(r1.Intn(40))
        monitor.Finished =  int32(r1.Intn(12500))
        
		buffer, err := proto.Marshal(monitor) //SerializeToOstream
		// 发送数据
		//fmt.Println(msg.String(), len(buffer), " ", buffer)
		_, err = socket.Write(buffer)
		if err != nil {
			fmt.Println("发送数据失败!", err)
			//return
		}
		
		
		time.Sleep(time.Duration(3) * time.Second)
	}

}
