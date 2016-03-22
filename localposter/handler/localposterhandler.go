package handler

import (
	//"fmt"
	. "github.com/fakewechat/lib/clientpool"
	. "github.com/fakewechat/lib/monitor"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/message"
	"github.com/garyburd/redigo/redis"
	"time"
)

type ProcessHandler interface {
	GetUserInfo(userid uint64) *UserInfor
	UpdateUser(userid uint64, u *UserInfor) (bool, error)
	UpdateUserAndInbox(userid uint64, u *UserInfor, Req []*GeneralMessage) (bool, error)
	UpdateUserAndOutbox(userid uint64, u *UserInfor, Req []*GeneralMessage) (bool, error)
	SendRequest(Req *GeneralMessage)
	GetRequest(userid uint64, key string) (Req *GeneralMessage)
	StoreRequest(userid uint64, key string, Req *GeneralMessage)

	Watch1(id string)
	Watch2(id string, id2 string)
	//Watch3(id string, id2 string, id3 string)
	UnWatch1(id string)
	UnWatch2(id string, id2 string)
	//UnWatch3(id string, id2 string, id3 string)

}

type LocalPosterHandler struct {
	Redisclient ClientInterface
	Rpcclient   ClientInterface
}

func (handler *LocalPosterHandler) Setup(Redisclient ClientInterface, Rpcclient ClientInterface) {
	handler.Redisclient = Redisclient
	handler.Rpcclient = Rpcclient

}

func (handler *LocalPosterHandler) GetUserInfo(userid uint64) *UserInfor {

	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())

	GMonitor.Add("GetUserInfo", 1)
	t1 := time.Now().UnixNano()

	n, err := conn.Do("GET", GetUserInfoName(userid))
	t2 := time.Now().UnixNano()
	GMonitor.Add64("GetUserInfo_time", t2-t1)
	if err != nil {

	}
	GMonitor.Add("DBload", 1)
	user := UserFromRedis(n.([]byte))
	return user
}

func (handler *LocalPosterHandler) UpdateUser(userid uint64, u *UserInfor) (bool, error) {

	userInfoName := GetUserInfoName(userid)

	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	conn.Send("MULTI")
	conn.Send("SET", userInfoName, UserToRedis(u))
	GMonitor.Add("UpdateUser", 1)
	t1 := time.Now().UnixNano()
	r, err := conn.Do("EXEC")
	t2 := time.Now().UnixNano()
	GMonitor.Add64("UpdateUser_time", t2-t1)
	if err != nil {
	} else {
		if r == nil {
		} else {
			GMonitor.Add("DBwriteSuccess", 1)
			return true, nil
		}
	}
	GMonitor.Add("DBconflict", 1)
	return false, nil

}

func (handler *LocalPosterHandler) UpdateUserAndInbox(userid uint64, u *UserInfor, Req []*GeneralMessage) (bool, error) {

	userInfoName := GetUserInfoName(userid)
	userInboxName := GetUserChatInBoxName(userid)

	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	conn.Send("MULTI")
	conn.Send("SET", userInfoName, UserToRedis(u))
	for i := range Req {
		conn.Send("RPUSH", userInboxName, MessageToRedis(Req[i]))
	}
	GMonitor.Add("UpdateUserAndInbox", 1)
	t1 := time.Now().UnixNano()
	r, err := conn.Do("EXEC")
	t2 := time.Now().UnixNano()
	GMonitor.Add64("UpdateUserAndInbox_time", t2-t1)
	if err != nil {
	} else {
		if r == nil {

		} else {
			GMonitor.Add("DBwriteSuccess", 1)
			return true, nil
		}
	}
	GMonitor.Add("DBconflict", 1)
	return false, nil
}

func (handler *LocalPosterHandler) UpdateUserAndOutbox(userid uint64, u *UserInfor, Req []*GeneralMessage) (bool, error) {

	userInfoName := GetUserInfoName(userid)
	useroutboxName := GetUserChatOutBoxName(userid)
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	conn.Send("MULTI")
	conn.Send("SET", userInfoName, UserToRedis(u))
	for i := range Req {
		conn.Send("RPUSH", useroutboxName, MessageToRedis(Req[i]))
	}
	GMonitor.Add("UpdateUserAndOutbox", 1)
	t1 := time.Now().UnixNano()
	r, err := conn.Do("EXEC")
	t2 := time.Now().UnixNano()
	GMonitor.Add64("UpdateUserAndOutbox_time", t2-t1)

	if err != nil {

		panic("err != nil")
	} else {
		if r == nil {
		} else {
			GMonitor.Add("DBwriteSuccess", 1)
			return true, nil
		}
	}
	GMonitor.Add("DBconflict", 1)
	return false, nil
}

func (handler *LocalPosterHandler) SendRequest(Req *GeneralMessage) {
	var ids uint64

	GMonitor.Add("rpc_send", 1)
	t1 := time.Now().UnixNano()
	handler.Rpcclient.GetRPCClient().Call("PosterAPI.PosterMessage", Req, &ids)
	t2 := time.Now().UnixNano()
	GMonitor.Add64("rpc_send_time", t2-t1)
	GMonitor.Add("LocalFowardPost", 1)
}

func (handler *LocalPosterHandler) GetRequest(userid uint64, key string) (Req *GeneralMessage) {
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	GMonitor.Add("GetRequest", 1)
	t1 := time.Now().UnixNano()

	r, err := conn.Do("HGET", GetUserQueue(userid), key)
	t2 := time.Now().UnixNano()
	GMonitor.Add64("GetRequest_time", t2-t1)

	if err != nil {

		panic("err != nil")
	} else {
		if r == nil {

		} else {

			return MessageFromRedis(r.([]byte))
		}
	}
	panic("GetRequest faile")
}

func (handler *LocalPosterHandler) StoreRequest(userid uint64, key string, Req *GeneralMessage) {
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())

	GMonitor.Add("StoreRequest", 1)
	t1 := time.Now().UnixNano()
	r, err := conn.Do("HSET", GetUserQueue(userid), key, MessageToRedis(Req))

	t2 := time.Now().UnixNano()
	GMonitor.Add64("StoreRequest_time", t2-t1)
	if err != nil {

		panic("err != nil")
	} else {
		if r == nil {

		} else {
			return
			//fmt.Println(r)
		}
	}
	panic("storeRequest faile")

}

func (handler *LocalPosterHandler) Watch1(id string) {

	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	_, err := conn.Do("WATCH", id)
	if err != nil {
	}
}

func (handler *LocalPosterHandler) Watch2(id1 string, id2 string) {
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())

	_, err := conn.Do("WATCH", id1, id2)
	if err != nil {
	}
}

func (handler *LocalPosterHandler) UnWatch1(id string) {
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())

	_, err := conn.Do("WATCH", id)
	if err != nil {
	}
}

func (handler *LocalPosterHandler) UnWatch2(id1 string, id2 string) {
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())

	_, err := conn.Do("WATCH", id1, id2)
	if err != nil {
	}
}
