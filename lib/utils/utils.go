package utils

import (
	"fmt"
	. "github.com/fakewechat/message"
	"github.com/golang/protobuf/proto"
	//"gopkg.in/redis.v3"
	"net/http"
	"strconv"
	//"strings"
)

func GetUserFriendOutBoxName(i uint64) string {
	//"F#outbox_#userid
	return "F#outbox" + strconv.Itoa(int(i))
}

func GetUserFriendInBoxName(i uint64) string {
	return "F#inbox" + strconv.Itoa(int(i))
}

func GetUserChatOutBoxName(i uint64) string {

	return "C#outbox" + strconv.Itoa(int(i))
}

func GetUserChatInBoxName(i uint64) string {

	return "C#inbox" + strconv.Itoa(int(i))
}

func GetUserChatInBoxQueueName(i uint64) string {
	return "C#inboxqueue" + strconv.Itoa(int(i))
}

func GetUserInfoName(i uint64) string {
	//fmt.Println("GetUserInfoName :", i)
	return "user#" + strconv.Itoa(int(i))
}

func GetReqdate(req *http.Request, key string) string {
	value := ""
	if len(req.Form[key]) > 0 {
		value = req.Form[key][0]
		//fmt.Println(value)
	}
	return value
}

func GetReqIntdata(req *http.Request, key string) int {

	value := ""
	if len(req.Form[key]) > 0 {
		value = req.Form[key][0]
		fmt.Println(value)
	}
	return GetIntValue(value)
}

func GetIntValue(key string) int {

	result, err := strconv.Atoi(key)
	if err != nil {
		result = 0
	}
	return result
}

func UserFromRedis(in []byte) *UserInfor {
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

	if user.SendedQueue == nil || user.SendedQueue.MessageMap == nil {
		user.SendedQueue = &SendQueue{}
		user.SendedQueue.MessageMap = make(map[uint64]*GeneralMessage)
	}

	if user.RecvedQueue == nil || user.RecvedQueue.MessageMap == nil {
		user.RecvedQueue = &RecvQueue{}
		user.RecvedQueue.MessageMap = make(map[uint64]*SendQueue)
	}

	if user.AckedQueue == nil || user.AckedQueue.MessageMap == nil {
		user.AckedQueue = &AckQueue{}
		user.AckedQueue.MessageMap = make(map[uint64]*GeneralMessage)
	}
	//fmt.Println(ToStrU(user))
	//user.Userlist = &UserList{Uidlist: make([]*User, 0)}
	return user

	/*
		    for k, v := range val.Val() {

				fmt.Println("k, v =", k, v)
				if k%2 == 1 {
					continue
				}
				if v == "Sendid" {
					fmt.Println("get sendindex", v, val.Val()[k+1])
					user.Sendid = proto.Uint64(uint64(GetIntValue(val.Val()[k+1])))
				} else if v == "Receivedid" {
					user.Receivedid = proto.Uint64(uint64(GetIntValue(val.Val()[k+1])))
				} else if v == "Clientreceivedid" {
					user.Clientreceivedid = proto.Uint64(uint64(GetIntValue(val.Val()[k+1])))
				} else if v == "Sendackid" {
					user.Sendackid = proto.Uint64(uint64(GetIntValue(val.Val()[k+1])))

				} else {
					fmt.Println("userid ", k, v)
					client := &User{}
					value := val.Val()[k+1]
					values := strings.Split(value, ",")
					fmt.Println("value:", value, " values:", values)
					client.Usrid = proto.Uint64(uint64(k))
					client.Rgid = proto.Uint32(0)
					client.Sendid = proto.Uint64(uint64(GetIntValue(values[0])))
					client.Recvid = proto.Uint64(uint64(GetIntValue(values[1])))

					user.Userlist.Uidlist = append(user.Userlist.Uidlist, client)
				}
			}
	*/
	//return user

}

func UserToRedis(u *UserInfor) []byte {
	out, err := proto.Marshal(u)
	if err != nil {
		panic("UserToRedis failed")
		return nil
	}
	return out

	/*
	   	result := make([]string, 0)




	       result = append(result, "Sendid")
	   	result = append(result, strconv.Itoa(int(*u.Sendid)))

	   	result = append(result, "Sendackid")
	   	result = append(result, strconv.Itoa(int(*u.Sendackid)))

	   	result = append(result, "Receivedid")
	   	result = append(result, strconv.Itoa(int(*u.Receivedid)))

	   	result = append(result, "Clientreceivedid")
	   	result = append(result, strconv.Itoa(int(*u.Clientreceivedid)))

	   	for _, v := range u.Userlist.Uidlist {
	   		fmt.Println(v)
	   		userid := *v.Usrid
	   		sendid := *v.Sendid
	   		Recvid := *v.Recvid
	   		result = append(result, strconv.Itoa(int(userid)))
	   		result = append(result, strconv.Itoa(int(sendid))+","+strconv.Itoa(int(Recvid)))

	   	}
	   	fmt.Println("==")
	   	for _, value := range result {
	   		fmt.Println(value)

	   	}

	   	return result
	*/

}

func MessageToRedis(u *GeneralMessage) []byte {
	out, err := proto.Marshal(u)
	if err != nil {
		panic("MessageToRedis failed")
		return nil
	}
	return out
}

func ToStr(M *GeneralMessage) string {
	s := " Type:" + strconv.Itoa(int(M.MessageType))
	s += " ID:" + strconv.Itoa(int(M.SendId))
	s += " From:" + strconv.Itoa(int(M.SenderId))
	s += " To:" + strconv.Itoa(int(M.ReceiverId))
	return s
}

func ToStrU(M *UserInfor) string {

	s := " SendId:" + strconv.Itoa(int(M.SendId))
	s += " SendAckId:" + strconv.Itoa(int(M.SendAckId))
	s += " ReceiveId:" + strconv.Itoa(int(M.ReceiveId))
	t := ""
	for i := range M.SendedQueue.MessageMap {
		t += strconv.Itoa(int(i)) + " "
	}
	s += " Send:[" + t + "]"

	t = ""
	for i := range M.RecvedQueue.MessageMap {
		for k := range M.RecvedQueue.MessageMap[i].MessageMap {
			t += strconv.Itoa(int(k)) + " "
		}
	}
	s += " RecvedQueue:[" + t + "]"

	t = ""
	for i := range M.AckedQueue.MessageMap {
		t += strconv.Itoa(int(i)) + " "
	}
	s += " Acked:[" + t + "]"

	return s
}
