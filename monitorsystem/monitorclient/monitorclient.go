package main

import (
	//"encoding/binary"
	"flag"
	"fmt"
	"github.com/fakewechat/message"
	"github.com/golang/protobuf/proto"
	"net"
	//"os"
	"strconv"
	"sort"
	"time"
)

const (
	UDP_PORT = 8002
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
	buffque = make(chan []byte, 10240)
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
	
	dic := make(map[string] int)
	keys := make([]string,0)
	starttime := time.Now().Unix()
	for {
		buff := <-buffque

		msg := &message.ClientMonitorStatus{}
		err := proto.Unmarshal(buff, msg) //unSerialize
		if err !=nil {
			fmt.Println(err)
			continue
		}
		finished := int(msg.Finished)
		value, ok := dic[msg.Host]
		if !ok {
			dic[msg.Host] = finished
			keys = append(keys, msg.Host)
			sort.Strings(keys)
		} else if value ==  finished {
			continue
		} else {
			dic[msg.Host] = finished
		}
		
		DrawScreen(dic, keys, starttime)

		buff = nil
	}
	
}

func CleanScreen() {
	s := fmt.Sprintf("\x1b[0;0H")
	
	fmt.Println(s)
	s = fmt.Sprint("\x1b[2J")
	fmt.Println(s)
		
}

func WriteString (color int, x int, y int, str string) {
		s1 := fmt.Sprintf("\x1b[%d;%dH", y, x)
		
		s2 := fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str  )
		fmt.Println(s1+s2)

}

func DrawScreen( dict map[string] int, keys[]string, starttime int64) {
	CleanScreen()
	
	Red := 31
	Yellow := 33
	sum := 0
	for _,v := range dict {
		sum += v
	}
	spendtime := time.Now().Unix() - starttime
	WriteString (Red, 30, 0, "sum: " +  strconv.Itoa(sum) + " sepend " +  strconv.Itoa(int(spendtime)) + " seconds"  )
	
	index := 0
	for k  := range keys {
		
		x := 20*(index % 4)
		y := 2 + index / 4
		WriteString(Yellow, x, y, keys[k] + ":" + strconv.Itoa(  dict[ keys[k]  ]  )   )
        index += 1
	}
	
}
