package main

import (
	"fmt"

	"encoding/binary"
	"flag"
	. "github.com/fakewechat/message"
	"github.com/golang/protobuf/proto"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

var LogName *string

func main() {

	LogName = flag.String("log", "log.txt", "file log")

	flag.Parse()
	f, err := os.Open(*LogName)
	if err != nil {
		return
	}
	defer f.Close()

	l := make([]*MonitorStatus, 0)
	for {
		head := make([]byte, 2)
		n2, err := f.Read(head)
		if err != nil {
			break
		}
		if n2 != 2 {
			break
		}
		needread := binary.BigEndian.Uint16(head[0:2])
		//fmt.Println("need read ", needread)
		buff := make([]byte, needread)
		n2, err = f.Read(buff)
		if err != nil {
			break
		}
		if n2 != int(needread) {
			break
		}

		msg := &MonitorStatus{}
		err = proto.Unmarshal(buff, msg) //unSerialize
		if err != nil {
			break
		}
		l = append(l, msg)
		//fmt.Printf("read: %s\n len: %d\n", msg.String(),  needread)
	}

	mapall := getResult(l, "GateWayRequestFromClient", "gw", "*")
	map1 := getResult(l, "GateWayRequestFromClient", "gw", "gw5")
	map2 := getResult(l, "GateWayRequestFromClient", "gw", "gw6")
	map3 := getResult(l, "GateWayRequestFromClient", "gw", "gw7")
	map4 := getResult(l, "GateWayRequestFromClient", "gw", "gw8")

	merge("gatewayandclient.txt", mapall, map1, map2, map3, map4)

	mapall = getResult(l, "GateWayRequestFromClient", "gw", "*")
	map1 = getResult(l, "GateWayUserInfo", "gw", "*")
	map2 = getResult(l, "GateWayChatMessage", "gw", "*")
	map3 = getResult(l, "GateWay200", "gw", "*")
	map4 = getResult(l, "GateWay503", "gw", "*")

	merge("gateway.txt", mapall, map1, map2, map3, map4)

	mapall = getResult(l, "PostRequest", "poster", "*")
	map1 = getResult(l, "PostFowardPost", "poster", "*")
	map2 = getResult(l, "PostFowardLocal", "poster", "*")
	map3 = getResult(l, "DropedFromQueueFull", "poster", "*")

	merge("poster.txt", mapall, map1, map2, map3)

	mapall = getResult(l, "LocalPostRequest", "localposter", "*")
	map1 = getResult(l, "LocalFowardPost", "localposter", "*")

	merge("localposter.txt", mapall, map1)

	map1 = getResult(l, "DBconflict", "localposter", "*")
	map2 = getResult(l, "DBload", "localposter", "*")
	map3 = getResult(l, "DBwriteSuccess", "localposter", "*")

	merge("db.txt", map1, map2, map3)

	map1 = getResult(l, "CLientSendRequest", "localposter", "*")
	map2 = getResult(l, "LocalPosterRecvRequest", "localposter", "*")
	map3 = getResult(l, "LocalPosterRecvAck", "localposter", "*")

	merge("chatmessage.txt", map1, map2, map3)

}

type record struct {
	Timestamp int64
	Value     int64
}

func getResult(l []*MonitorStatus, name string, types string, host string) map[string][]record {

	maps := make(map[string][]record)

	for i := range l {

		statuslist := l[i].Status
		for k := range statuslist {

			r := statuslist[k]

			if r.Servertype == types && r.Name == name {

			} else {
				continue
			}

			key := types + host
			if host == "*" {
				key = types
			}

			if host != "" && host != "*" {
				if r.Servername != host {
					continue
				}
			}

			fmt.Println(r.String())
			locallist, ok := maps[key]
			if !ok {
				locallist = make([]record, 0)
				maps[key] = locallist
			}
			gotit := false
			for iter := range locallist {
				if locallist[iter].Timestamp == r.Timestamp {
					gotit = true
					locallist[iter].Value += r.Changesvalue
				} else {

				}
			}
			if gotit == false {
				tmp := record{}
				tmp.Timestamp = r.Timestamp
				tmp.Value = r.Changesvalue
				locallist = append(locallist, tmp)
				maps[key] = locallist
			}
		}
	}

	for mapkey := range maps {
		fmt.Println("mapkey:", mapkey)
		f, _ := os.Create(mapkey + ".txt")
		defer f.Close()
		locallist, _ := maps[mapkey]
		for iter := range locallist {
			r := locallist[iter]
			str_time := time.Unix(r.Timestamp, 0).Format("15:04:05")
			fmt.Println(str_time, r.Value)
			str_time += " " + strconv.Itoa(int(r.Value)) + "\n"
			io.WriteString(f, str_time)
		}

	}
	return maps
}

func merge(filename string, arg ...map[string][]record) {

	resultmap := make(map[int64]string)
	fmt.Println(len(arg))
	for onemap := range arg {
		currentmap := arg[onemap]
		for mapkey := range currentmap {
			locallist, _ := currentmap[mapkey]
			for iter := range locallist {
				r := locallist[iter]
				line, ok := resultmap[r.Timestamp]
				if !ok {
					line = time.Unix(r.Timestamp, 0).Format("15:04:05")
				}
				line += " " + strconv.Itoa(int(r.Value))
				fmt.Println(line)
				resultmap[r.Timestamp] = line
			}
		}
	}
	sorted_keys := make([]int, 0)
	for k, _ := range resultmap {
		sorted_keys = append(sorted_keys, int(k))
	}
	sort.Ints(sorted_keys)
	fmt.Println(filename)
	f, _ := os.Create(filename)
	defer f.Close()
	for i := range sorted_keys {
		io.WriteString(f, resultmap[int64(sorted_keys[i])]+"\n")
		fmt.Println(resultmap[int64(sorted_keys[i])])
	}
}
