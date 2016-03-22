package chatmessage

import (
	"fmt"
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

func CoreLocal_to_Local(user *UserInfor, req *GeneralMessage) (result int, needupdateInBox bool) {

	result = 0
	needupdateInBox = false
	// find out the user
	SendIdofclient := req.Chatmessage.SendId
	client_cmd2 += 1
	fmt.Println("client ", client_cmd2, ToStr(req))

	friend, ok := user.UserMap[req.SenderId]
	if !ok {
		panic("not exist friend") // we should drop it
	}
	/* for debug
	if SendIdofclient > 5 {

		fmt.Println(SendIdofclient, friend.ReceiveId)
		panic("SendIdofclient > 5")
	}*/

	if friend.ReceiveId >= SendIdofclient {
		// we need send a ack,
		result = LOCAL_TO_LOCAL_SUCCESS
	} else if friend.ReceiveId+1 == SendIdofclient { // good ! it it what we need
		friend.ReceiveId += 1
		needupdateInBox = true
		user.ReceiveId += 1
		/*
			if friend.ReceiveId > 5 {
				fmt.Println(SendIdofclient, friend.ReceiveId+1)
				panic("firend.ReceiveId > 5")
			}
		*/
		fmt.Println("2000 ", user.SendId, user.SendAckId, user.ReceiveId)

		result = LOCAL_TO_LOCAL_SUCCESS
	} else { //
		recvQ, ok := user.RecvedQueue.MessageMap[req.SenderId] // find the
		if !ok {
			recvQ = &SendQueue{}
			recvQ.MessageMap = make(map[uint64]*GeneralMessage)
			user.RecvedQueue.MessageMap[req.SenderId] = recvQ
		}
		if recvQ.MessageMap == nil {
			recvQ.MessageMap = make(map[uint64]*GeneralMessage)
		}
		_, reqexist := recvQ.MessageMap[SendIdofclient]
		if !reqexist {

			recvQ.MessageMap[SendIdofclient] = req
		}

		result = LOCAL_TO_LOCAL_SUCCESS_NOSEND
	}

	return result, needupdateInBox
}

func SyncLocal_to_Local(user *UserInfor, userid uint64) (int, *GeneralMessage) {
	result := 0
	var req *GeneralMessage
	friend, ok := user.UserMap[userid]
	if !ok {
		panic("not exist friend") // we should drop it
	}

	nextrecvid := friend.ReceiveId + 1

	recvQ, ok := user.RecvedQueue.MessageMap[userid]
	if !ok {
		result = LOCAL_TO_LOCAL_SYNC_SUCCESS_NOSEND
	} else {
		ok := false
		req, ok = recvQ.MessageMap[nextrecvid]
		if ok {

			delete(recvQ.MessageMap, nextrecvid)

			friend.ReceiveId += 1
			/* for debug
			if friend.ReceiveId > 5 {
				fmt.Println(friend.ReceiveId, nextrecvid)
				panic("firend.ReceiveId > 5")
			} */

			user.ReceiveId += 1
			fmt.Println("2000 ", user.SendId, user.SendAckId, user.ReceiveId)
			result = LOCAL_TO_LOCAL_SYNC_SUCCESS
		} else {
			req = nil
			result = LOCAL_TO_LOCAL_SYNC_SUCCESS_NOSEND
		}
	}
	return result, req

}

func CheckCoreLocal_to_Local(r int) bool {
	if r == LOCAL_TO_LOCAL_SUCCESS || r == LOCAL_TO_LOCAL_SUCCESS_NOSEND {
		return true
	} else {
		panic("CheckCoreLocal_to_Local")
	}

}

// sync
func CheckSyncLocal_to_Local(r int) bool {
	if r == LOCAL_TO_LOCAL_SYNC_SUCCESS || r == LOCAL_TO_LOCAL_SYNC_SUCCESS_NOSEND {
		return true
	} else {
		panic("CheckSyncLocal_to_Local")
	}

}
