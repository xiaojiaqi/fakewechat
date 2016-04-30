package main

import (
	"strings"
	"fmt"
	. "github.com/fakewechat/message"
	//"github.com/golang/protobuf/proto"
    . "github.com/fakewechat/lib/utils"
    "github.com/garyburd/redigo/redis"
	//"bufio"
	//"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
		"flag"
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

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	host = flag.String("host", "127.0.0.1", "host")
	port = flag.String("port", "1500", "port")
	minid = flag.Int("minid", 1, "minid")
	maxid = flag.Int("maxid", 2500, "maxid")
	filespath = flag.String("file", "db.txt", "file")

	flag.Parse()

	result := readfile(*filespath)
	conn, err := redis.Dial("tcp", *host + ":" +  *port  )
	if err != nil {
		fmt.Println(err)
		return
	}

	v := strings.Split(result, "\n")
    var userid uint64
	userid  = 0
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
			s.UserId = uint64( friendid)
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
