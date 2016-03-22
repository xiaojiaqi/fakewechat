package chatmessage

import (
	//"fmt"
	//. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/message"
)

const (
	// response of ProcessClient_to_LocalMessgae
	LOCALACK_SUCCESS        = 100 // success, send data
	LOCALACK_SUCCESS_NOSYNC = 110
	// send data

	// sync

	LOCALACK_SYNC_SUCCESS        = 120 // SUCCESS, send ack
	LOCALACK_SYNC_SUCCESS_NOSYNC = 130 // SUCCESS, no send , break

)

/*
func ProcessLocal_ack(req *GeneralMessage, Channelindex int) {

	/*

	   	message sendid                                     ackid


	         <=  ackid                                         do nothing



	         -----------------------------------------------

	   	 =  ackid +1                                        ackid + 1,   and sync3


	   	 ----------------------------------------------

	   	                                              message lost on net,  restore to queue,
	   	> client.ackid +1
*/
var client_cmd3 int

func CoreLocal_Ack(user *UserInfor, req *GeneralMessage) (result int) {

	result = 0
	client_cmd3 += 1

	_, ok := user.UserMap[req.SenderId]
	if !ok {
		panic("not exist friend") // we should drop it
	}

	//fmt.Println("client ", client_cmd3, ToStr(req), req.SendId, user.SendAckId )
	if req.SendId <= user.SendAckId {
		// do nothing
		result = LOCALACK_SUCCESS_NOSYNC

	} else if req.SendId == user.SendAckId+1 {
		user.SendAckId += 1
		//fmt.Println("Delete 3000 ", user.SendId, user.SendAckId, user.ReceiveId)
		delete(user.AckMessage, req.SendId)

		result = LOCALACK_SUCCESS

	} else {
		_, ok := user.AckMessage[req.SendId]
		if !ok {
			user.AckMessage[req.SendId] = req.SendId // restore a message is too heavy, id is enough
		}
		result = LOCALACK_SUCCESS_NOSYNC
	}
	return result
}

func SyncLocal_Ack(user *UserInfor, userid uint64) (result int) {

	result = LOCALACK_SYNC_SUCCESS_NOSYNC

	for {
		newid := user.SendAckId + 1
		_, ok := user.AckMessage[newid]
		if ok {

			delete(user.AckMessage, newid)
			// also remove send queue
			delete(user.SendedMessage, newid)
			//fmt.Println("delete ", newid)
			//fmt.Println(user.AckMessage)
			user.SendAckId = newid
			//fmt.Println("3000 ", user.SendId, user.SendAckId, user.ReceiveId)
			result = LOCALACK_SYNC_SUCCESS
		} else {
			//fmt.Println("delete ", newid, " not found")
			break
		}
	}
	return result
}

/*
func CheckCoreLocal_Ack(r int) bool {

	if (r == LOCALACK_SUCCESS) || (r == LOCALACK_SUCCESS_NOSYNC) {
		return true
	} else {
		panic("CheckCoreLocal_to_Local")
	}
}

// sync
func CheckSyncLocal_Ack(r int) bool {

	if (r == LOCALACK_SYNC_SUCCESS) || (r == LOCALACK_SYNC_SUCCESS_NOSYNC) {
		return true
	} else {

		panic("CheckCoreLocal_to_Local")
	}
}
*/
