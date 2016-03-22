package chatmessage

import (
	//"fmt"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/message"
)

type Raw_to_client struct {
	Rawname           string
	Sendname          string
	Chatmessagesendid uint64
}

const (
	// response of ProcessClient_to_LocalMessgae
	CLIENT_TO_LOCAL_SUCCESS        = 20 // success, send data
	CLIENT_TO_LOCAL_SUCCESS_NOSEND = 30 // success, don't send data

	// send data

	// sync

	CLIENT_TO_LOCAL_SYNC_SUCCESS        = 40 // SUCCESS, send data
	CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND = 50 // SUCCESS, no send , break

)

/*

		message sendid              sendid        ackid

	                                           1) messageid < ackid    :  drop message.  ack recved
	      <  sendid + 1                         2) messageid = ackid    :  drop message.  ack recved
	                                           3) messageid > ackid    :  send to local again, maybe ack loss


	      -----------------------------------------------
		                                       1)  messageid < ackid :    ?? bug
		 = sendid +1                           2)  messageid = ackid :    send to local again and  sendid +1, call sync
		                                       3)  messageid > ackid :    send to local again and  sendid +1, call sync

		 ----------------------------------------------

		                                       1)  messageid < ackid :   storage message,
		> sendid + 1                           2)  messageid = ackid :   storage message,
		                                       3)  messageid > ackid :   storage message,

		sync:

		if message == sendid +1 and send to local
		    for ;;
		     	check storage queue :
		        	if  message.id == sendid :
		           		send to local and sendid +1
		           	else:
		           	    break
*/

//refactor me

var client_cmd int

func CheckDropMessageClient_to_Local(user *UserInfor, req *GeneralMessage) bool {

	if user.SendAckId >= req.SendId {

		return true
	}
	return false
}

func CheckResendMessageClient_to_Local(user *UserInfor, req *GeneralMessage) bool {

	if req.SendId <= user.SendId { // 已经送过一次了
		return true
	}
	return false
}

func CheckStoragedMessageClient_to_Local(user *UserInfor, req *GeneralMessage) bool {
	_, ok := user.SendedMessage[req.SendId]
	if ok {
		return true
	} else {
		return false
	}

}

func StorageMessageClient_to_Local(user *UserInfor, req *GeneralMessage) (needstorage bool, needsync bool) {
	needstorage = false
	needsync = false
	_, ok := user.SendedMessage[req.SendId]
	if ok {
		return false, false
	}

	user.SendedMessage[req.SendId] = req.ReceiverId
	needstorage = true

	if user.SendId+1 == req.SendId {
		needsync = true
	}

	return needstorage, needsync

}

func SyncClient_to_Local(user *UserInfor) []Raw_to_client {
	nextsendid := user.SendId + 1
	list := make([]Raw_to_client, 0)

	for {
		recverId, ok := user.SendedMessage[nextsendid]
		if !ok {
			break
		}

		//find user
		friend, ok := user.UserMap[recverId]
		if !ok {
			//fmt.Println("no such friend", recverId)
			panic("not exist friend") // we should drop it
		}
		a := Raw_to_client{}
		a.Rawname = GetRawMessage(nextsendid)
		a.Sendname = GetSendMessage(nextsendid)

		friend.SendId += 1
		a.Chatmessagesendid = friend.SendId
		user.SendId += 1

		list = append(list, a)
		delete(user.SendedMessage, nextsendid)

		nextsendid += 1

	}
	return list
}
