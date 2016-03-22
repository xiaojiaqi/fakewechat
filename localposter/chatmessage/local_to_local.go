package chatmessage

import (
	//"fmt"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/message"
)

const (
	// response of ProcessClient_to_LocalMessgae
	LOCAL_TO_LOCAL_SUCCESS        = 60 // success, send data
	LOCAL_TO_LOCAL_SUCCESS_NOSEND = 70 // success, drop data

	// sync

	LOCAL_TO_LOCAL_SYNC_SUCCESS        = 80 // SUCCESS, send ack
	LOCAL_TO_LOCAL_SYNC_SUCCESS_NOSEND = 90 // SUCCESS, no send , no send

)

/*
func ProcessLocal_to_Local(req *GeneralMessage, Channelindex int) {

	/*

	   	message sendid                                  client.recvid


	         <= client.recvid                               do nothing, get client recvid, send ack



	         -----------------------------------------------

	   	 = client.recvid +1                           client.recv + 1,  receivedid +1 , send ack , and sync2


	   	 ----------------------------------------------

	   	                                              message lost on net,  restore to queue,
	   	> client.recvid +1

	     sync2
	   	if message == client.recvid +1 and send ack
	   	    for ;;
	   	     	check storage queue :
	   	        	if  message.id == client.recvid +1 :
	   	           		client.recv + 1,  receivedid +1 , send ack ,
	   	           	else:
	   	           	    break


*/
var client_cmd2 int

// need refactor

//  无需存储消息
//  无需更新userinfo
//  只需要简单回复即可

//  情况 1, 消息的 chatmessage 的sendid <= friend.ReceiveId,

func CheckResponOnlyLocal_to_Local(user *UserInfor, req *GeneralMessage) bool {

	//
	SendIdofclient := req.Chatmessage.SendId

	friend, ok := user.UserMap[req.SenderId]
	if !ok {
		panic("not exist friend") // we should drop it
	}

	if friend.ReceiveId >= SendIdofclient {
		return true
	} else {
		return false
	}
}

//  无需存储消息
//  无需更新userinfo
//  无需回应

// 情况  LocalMessage 里面已经有了
func CheckStoreagedLocal_to_Local(user *UserInfor, req *GeneralMessage) bool {

	_, ok := user.UserMap[req.SenderId]
	if !ok {
		panic("not exist friend") // we should drop it
	}

	recvQ, ok := user.LocalMessage[req.SenderId] // find the friend
	if !ok {
		return false
	}
	SendIdofclient := req.Chatmessage.SendId

	_, ok = recvQ.MessageMap[SendIdofclient]
	if ok {
		return true
	} else {
		return false
	}
}

// 更新userinfor， 在recvQ 里面加一条记录
//
// 因为时序问题，可能存在 下次重做 别的线程已经完成的情况, needstorage = false 表明 不需要保存了
// nosync 表明不需要去刷新了

func StoreageLocal_to_Local(user *UserInfor, req *GeneralMessage) (needstorage bool, needsync bool) {
	friend, ok := user.UserMap[req.SenderId]
	if !ok {
		panic("not exist friend") // we should drop it
	}

	recvQ, ok := user.LocalMessage[req.SenderId] // find the
	if !ok {
		recvQ = &RecvQueue{}
		recvQ.MessageMap = make(map[uint64]uint64)
		user.LocalMessage[req.SenderId] = recvQ
	}

	if recvQ.MessageMap == nil {

		recvQ.MessageMap = make(map[uint64]uint64)
	}

	SendIdofclient := req.Chatmessage.SendId

	_, ok = recvQ.MessageMap[SendIdofclient]
	if !ok {

		recvQ.MessageMap[SendIdofclient] = SendIdofclient
		needstorage = true

	} else {
		needstorage = false
	}

	if SendIdofclient == friend.ReceiveId+1 {
		needsync = true
	}
	return needstorage, needsync
}

//  如果结果不空, update 用户信息
//  inbox update
//  返回

func SyncLocal_to_Local(user *UserInfor, senderId uint64) []string {
	Result := make([]string, 0)

	friend, ok := user.UserMap[senderId]
	if !ok {
		panic("not exist friend") // we should drop it
	}
	nextrecvid := friend.ReceiveId + 1

	recvQ, ok := user.LocalMessage[senderId]
	if !ok {
		return Result
	}
	for {
		if recvQ.MessageMap == nil {
			recvQ.MessageMap = make(map[uint64]uint64)
		}

		_, ok := recvQ.MessageMap[nextrecvid]

		if ok {
			friend.ReceiveId += 1
			delete(recvQ.MessageMap, nextrecvid)

			id := GetLocalMessage(senderId, nextrecvid)
			Result = append(Result, id)
			nextrecvid += 1
			user.ReceiveId += 1
		} else {
			break
		}

	}
	return Result

}
