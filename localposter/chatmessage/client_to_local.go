package chatmessage

import (
	. "github.com/fakewechat/message"
	. "github.com/fakewechat/lib/utils"
	"fmt"
	)


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
func CoreClient_to_Local(user *UserInfor, req *GeneralMessage) (result int, needupdateOutBox bool) {
	result = 0
    client_cmd += 1
    fmt.Println("client ", client_cmd, ToStr(req))
	needupdateOutBox = false
	if req.SendId <= user.SendId { // sended meesgae
		

		if user.SendAckId >= req.SendId {
			result = CLIENT_TO_LOCAL_SUCCESS_NOSEND
		} else { // has been send to local , but didn't recv the ack, may been loss, need retry
			r, ok := user.SendedQueue.MessageMap[req.SendId]
		if !ok {
			panic("no data")
		} else {
			*req = *r
		}
			result = CLIENT_TO_LOCAL_SUCCESS
		}
	} else if req.SendId == user.SendId+1 { // good ! it it what we need
		user.SendId += 1
		needupdateOutBox = true

		//find user
		friend, ok := user.UserMap[req.ReceiverId]
		if !ok {
			panic("not exist friend") // we should drop it
		}
		friend.SendId += 1

		req.Chatmessage.SendId = friend.SendId
		/*
		if req.Chatmessage.MessageBody != strconv.Itoa(int(friend.SendId)) {
			fmt.Println(req.Chatmessage.MessageBody, friend.SendId)
			panic("req.Chatmessage.MessageBody != strconv.Itoa( friend.SendId)")
		}*/
		user.SendedQueue.MessageMap[req.SendId] = req
		fmt.Println("1000 ",  user.SendId, user.SendAckId, user.ReceiveId)
		
		result = CLIENT_TO_LOCAL_SUCCESS
	} else { // the
		_, ok := user.SendedQueue.MessageMap[req.SendId]
		if !ok {
			user.SendedQueue.MessageMap[req.SendId] = req
		} else {
			// do noting
		}
		result = CLIENT_TO_LOCAL_SUCCESS_NOSEND
	}

	return result, needupdateOutBox
}

func SyncClient_to_Local(user *UserInfor) (int, *GeneralMessage) {
	result := 0
	nextsendid := user.SendId + 1

	var req *GeneralMessage

	mess, ok := user.SendedQueue.MessageMap[nextsendid]
	if ok {
		result = CLIENT_TO_LOCAL_SYNC_SUCCESS
		req = mess

		//find user
		friend, ok := user.UserMap[req.ReceiverId]
		if !ok {
			panic("not exist friend") // we should drop it
		}
		friend.SendId += 1
		req.Chatmessage.SendId = friend.SendId
		result = CLIENT_TO_LOCAL_SYNC_SUCCESS
		user.SendId += 1
		/* for test only
		if req.Chatmessage.MessageBody != strconv.Itoa(int(friend.SendId)) {
			fmt.Println(req.Chatmessage.MessageBody, friend.SendId)
			panic("req.Chatmessage.MessageBody != strconv.Itoa( friend.SendId)")
		}
		*/
	} else {
		// do noting
		result = CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND
	}
	return result, req

}

func CheckCoreClient_to_Local(r int) bool {
	if r == CLIENT_TO_LOCAL_SUCCESS || r == CLIENT_TO_LOCAL_SUCCESS_NOSEND {
		return true
	} else {
		panic("CheckCoreClient_to_Local")
	}

}


// sync
func CheckSyncClient_to_Local(r int) bool {
	if r == CLIENT_TO_LOCAL_SYNC_SUCCESS || r == CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND {
		return true
	} else {
		panic("CheckCoreClient_to_Local")
	}

}
