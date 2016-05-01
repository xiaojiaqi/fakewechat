package main

import (
	"fmt"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/message"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"strings"
	//"bufio"
	//"io"
	"flag"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
)

func readfile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

var host *string
var port *string
var minid *int
var maxid *int
var filespath *string
var method *string

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	host = flag.String("host", "127.0.0.1", "host")
	port = flag.String("port", "1500", "port")
	minid = flag.Int("minid", 1, "minid")
	maxid = flag.Int("maxid", 2500, "maxid")
	filespath = flag.String("file", "db.txt", "file")
	method = flag.String("method", "save", "method")
	flag.Parse()
	if *method == "save" {
		saveFile()
	} else {
		checkDB()
	}

}

func saveFile() {
	result := readfile(*filespath)
	conn, err := redis.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println(err)
		return
	}

	v := strings.Split(result, "\n")
	var userid uint64
	userid = 0
	userInfoName := ""
	user := &UserInfor{}
	if user.UserMap == nil {
		user.UserMap = make(map[uint64]*User)
	}

	for _, i := range v {

		items := strings.Split(i, ",")

		if len(items) < 3 {
			continue
		}
		itemid := GetIntValue(items[0])
		friendid := GetIntValue(items[1])

		if userid == 0 {
			userid = uint64(itemid)
			user.UserId = uint32(userid)
			userInfoName = "user#" + strconv.Itoa(itemid)
			s := &User{}
			s.UserId = uint64(friendid)
			user.UserMap[uint64(friendid)] = s
		} else if userid == uint64(itemid) {
			s := &User{}
			s.UserId = uint64(friendid)
			user.UserMap[uint64(friendid)] = s
		} else {
			conn.Send("MULTI")
			//fmt.Println(*user)
			conn.Send("SET", userInfoName, UserToRedis(user))
			r, err := conn.Do("EXEC")
			if err != nil {
				fmt.Println(r)
			}
			userid = 0
			user = &UserInfor{}
			if user.UserMap == nil {
				user.UserMap = make(map[uint64]*User)
			}

			userid = uint64(itemid)
			user.UserId = uint32(userid)
			userInfoName = "user#" + strconv.Itoa(itemid)
			s := &User{}
			s.UserId = uint64(friendid)
			user.UserMap[uint64(friendid)] = s
		}
	}
	if userid != 0 {
		conn.Send("MULTI")
		//fmt.Println(*user)
		conn.Send("SET", userInfoName, UserToRedis(user))

		r, err := conn.Do("EXEC")
		if err != nil {
			fmt.Println(r)
		}
	}

}

func checkDB() {

	conn, err := redis.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println(err)
		return
	}
	sumsize := 0
	for i := *minid; i <= *maxid; i++ {
		ui := uint64(i)
		userInfoName := GetUserInfoName(ui)
		cin := GetUserChatInBoxName(ui)
		cout := GetUserChatOutBoxName(ui)

		n, err := conn.Do("GET", userInfoName)

		if err != nil {
			panic(err)

		}
		if n == nil {
			panic(n)
		}

		user := UserFromRedis1(n.([]byte))
		n, err = redis.Int(conn.Do("LLEN", cin))

		if err != nil {
			panic(err)

		}
		insize := n.(int)

		n, err = redis.Int(conn.Do("LLEN", cout))

		if err != nil {
			panic(err)

		}
		outsize := n.(int)

		if insize != outsize {
			panic("insize != outsize")
		}
		if 5*len(user.UserMap) != outsize {
			panic("insize != outsize")
		}
		//fmt.Println(outsize, insize, 5 * len(user.UserMap))
		sumsize += insize
	}

	fmt.Println("success", sumsize)

}

func UserFromRedis1(in []byte) *UserInfor {

	user := &UserInfor{}
	err := proto.Unmarshal(in, user)
	if err != nil {
		fmt.Println(err)
		panic("proto.Unmarshal(rawbyte, user)")
		return nil
	}

	if user.UserMap == nil {
		user.UserMap = make(map[uint64]*User)
	}

	if user.SendedMessage == nil {
		user.SendedMessage = make(map[uint64]uint64)
	}

	if user.LocalMessage == nil {
		user.LocalMessage = make(map[uint64]*RecvQueue)
	}

	if user.AckMessage == nil {
		user.AckMessage = make(map[uint64]uint64)
	}
	return user

}
