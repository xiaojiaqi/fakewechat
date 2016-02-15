package handler

import (
	. "github.com/fakewechat/lib/clientpool"
	. "github.com/fakewechat/lib/monitor"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/message"
	"github.com/garyburd/redigo/redis"
)

type ProcessHandler interface {
	GetUserInfo(userid uint64) *UserInfor
	UpdateUser(userid uint64, u *UserInfor) bool
	UpdateUserAndInbox(userid uint64, u *UserInfor, Req *GeneralMessage) bool
	UpdateUserAndOutbox(userid uint64, u *UserInfor, Req *GeneralMessage) bool
	SendRequest(Req *GeneralMessage)
	Watch(id string)
	Watch2(id string, id2 string)
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
	n, err := conn.Do("GET", GetUserInfoName(userid))
	if err != nil {

	}
	GMonitor.Add("DBload", 1)
	user := UserFromRedis(n.([]byte))
	return user
}

func (handler *LocalPosterHandler) UpdateUser(userid uint64, u *UserInfor) bool {

	userInfoName := GetUserInfoName(userid)
	
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	conn.Send("MULTI")
	conn.Send("SET", userInfoName, UserToRedis(u))

	r, err := conn.Do("EXEC")
	if err != nil {
	} else {
		if r == nil {
		} else {
			GMonitor.Add("DBwriteSuccess", 1)
			return true
		}
	}
	GMonitor.Add("DBconflict", 1)
	return false

}

func (handler *LocalPosterHandler) UpdateUserAndInbox(userid uint64, u *UserInfor, Req *GeneralMessage) bool {

	userInfoName := GetUserInfoName(userid)
	userInboxName := GetUserChatInBoxName(userid)

	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	conn.Send("MULTI")
	conn.Send("SET", userInfoName, UserToRedis(u))
	conn.Send("RPUSH", userInboxName, MessageToRedis(Req))

	r, err := conn.Do("EXEC")
	if err != nil {
	} else {
		if r == nil {

		} else {
			GMonitor.Add("DBwriteSuccess", 1)
			return true
		}
	}
	GMonitor.Add("DBconflict", 1)
	return false
}

func (handler *LocalPosterHandler) UpdateUserAndOutbox(userid uint64, u *UserInfor, Req *GeneralMessage) bool {

	userInfoName := GetUserInfoName(userid)
	useroutboxName := GetUserChatOutBoxName(userid)
	conn := (redis.Conn)(*handler.Redisclient.GetRedisClient())
	conn.Send("MULTI")
	conn.Send("SET", userInfoName, UserToRedis(u))
	conn.Send("RPUSH", useroutboxName, MessageToRedis(Req))

	r, err := conn.Do("EXEC")
	if err != nil {
	} else {
		if r == nil {
		} else {
			GMonitor.Add("DBwriteSuccess", 1)
			return true
		}
	}
	GMonitor.Add("DBconflict", 1)
	return false
}

func (handler *LocalPosterHandler) SendRequest(Req *GeneralMessage) {
	var ids uint64
	handler.Rpcclient.GetRPCClient().Call("PosterAPI.PosterMessage", Req, &ids)

	GMonitor.Add("LocalFowardPost", 1)
}

func (handler *LocalPosterHandler) Watch(id string) {

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
