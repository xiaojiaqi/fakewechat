package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/fakewechat/message"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
)

const (
	UDP_PORT = 8000
)

var LogName *string
var buffque chan ([]byte)

func main() {

	LogName = flag.String("log", "log.txt", "file log")
	flag.Parse()
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: UDP_PORT,
	})
	buffque = make(chan []byte, 1024)
	if err != nil {
		fmt.Println("listen failed", err)
		return
	}

	defer socket.Close()

	go forward()
	for {

		data := make([]byte, 65536)
		read, _, err := socket.ReadFromUDP(data)

		if err != nil {
			fmt.Println("read data failed", err)
			continue
		}
		buffque <- data[:read]

	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func forward() {
	f, err := os.Create("./" + *LogName) //create log file
	check(err)
	defer f.Close()
	for {
		buff := <-buffque

		msg := &message.MonitorStatus{}
		err := proto.Unmarshal(buff, msg) //unSerialize

		header := make([]byte, 2)
		buffLen := len(buff)
		binary.BigEndian.PutUint16(header, uint16(buffLen))

		f.Write(header)
		f.Write(buff)
		f.Sync()
		if err == nil {
			fmt.Printf("read: %s\n len: %d\n", msg.String(), buffLen)
		} else {
			fmt.Println(err)
			fmt.Printf("wrong data")
			os.Exit(0)
		}
		buff = nil
	}
}
