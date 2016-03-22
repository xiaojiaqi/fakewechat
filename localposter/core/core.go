package core

import (
	"fmt"
	. "github.com/fakewechat/lib/contstant"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/localposter/chatmessage"
	. "github.com/fakewechat/localposter/handler"
	. "github.com/fakewechat/message"
)

/*
  client => local
  sender = 100, recverid = 200, type = 100  // handle user 100

  local => local
  sender = 100, recverid = 200, type = 200  // handle user 200

  localack
  sender = 200,  recverid = 100, type = 300 // handle user 100


*/

const (
	CHAT_CLIENT_DROPED                 = 1001
	CHAT_CLIENT_RESEND                 = 1002
	CHAT_CLIENT_NONEEDSTORE            = 1003
	CHAT_CLIENT_STORE                  = 1004
	CHAT_CLIENT_STORE_NOSYNC           = 1005
	CHAT_CLIENT_SYNC_NOREQUESTNEEDSEND = 1006
	CHAT_CLIENT_SYNC_ANDSEND           = 1007

	//CHAT_LOCALPOST_TO_LOCALPOST = 2000
	CHAT_LOCAL_RESPONSEONLY      = 2001
	CHAT_LOCAL_DROP              = 2002
	CHAT_LOCAL_STOREONLY         = 2003
	CHAT_LOCAL_NOMESSAGENEEDSEND = 2004
	CHAT_LOCAL_MESSAGENEEDSEND   = 2005
	//CHAT_LOCALPOST_ACK          = 3000

	CHAT_ACK_NOSYNC = 3001
	CHAT_ACK_SYNC   = 3002
)

func ProcessClient_to_Local_Message(handler ProcessHandler, req *GeneralMessage) int {

	senderId := req.SenderId
	needsync := false

	//fmt.Println("ProcessClient_to_Local_Message ", req)
	for {
		handler.Watch1(GetUserInfoName(senderId))
		user := handler.GetUserInfo(senderId)

		if CheckDropMessageClient_to_Local(user, req) {
			handler.UnWatch1(GetUserInfoName(senderId))

			//fmt.Println("CheckDropMessageClient_to_Local ", req)
			return CHAT_CLIENT_DROPED
		}

		if CheckResendMessageClient_to_Local(user, req) {
			handler.UnWatch1(GetUserInfoName(senderId))

			// get Req
			mewreq := handler.GetRequest(senderId, GetSendMessage(req.SendId))
			// send

			//fmt.Println("CheckResendMessageClient_to_Local(user, req)")
			handler.SendRequest(mewreq)

			return CHAT_CLIENT_RESEND
		}
		if CheckStoragedMessageClient_to_Local(user, req) {
			handler.UnWatch1(GetUserInfoName(senderId))

			//fmt.Println("!CheckStoragedMessageClient_to_Local(user, req)")
			return CHAT_CLIENT_NONEEDSTORE
		}

		_, needsync = StorageMessageClient_to_Local(user, req)

		// store the message to map
		key := GetRawMessage(req.SendId)
		handler.StoreRequest(senderId, key, req)

		b, _ := handler.UpdateUser(senderId, user)
		if !b {
			continue
		} else {
			break
		}

	}

	//fmt.Println("needsync", needsync)
	if !needsync {

		return CHAT_CLIENT_STORE_NOSYNC
	}

	//fmt.Println("sync: ")

	// for less
	for {
		handler.Watch2(GetUserInfoName(senderId), GetUserChatOutBoxName(senderId))
		user := handler.GetUserInfo(senderId)
		if user == nil {
			continue
		}

		list := SyncClient_to_Local(user)
		if len(list) == 0 {
			handler.UnWatch2(GetUserInfoName(senderId), GetUserChatOutBoxName(senderId))

			return CHAT_CLIENT_SYNC_NOREQUESTNEEDSEND
		}

		reqlist := make([]*GeneralMessage, 0)

		for i := range list {
			key := list[i].Rawname
			newkey := list[i].Sendname
			newReq := handler.GetRequest(senderId, key)
			//fmt.Println(list[i])
			newReq.Chatmessage.SendId = list[i].Chatmessagesendid
			newReq.MessageType = CHAT_LOCALPOST_TO_LOCALPOST
			reqlist = append(reqlist, newReq)
			handler.StoreRequest(senderId, newkey, newReq)

		}
		b, _ := handler.UpdateUserAndOutbox(senderId, user, reqlist)
		if b {

			for i := range reqlist {
				handler.SendRequest(reqlist[i])

			}
			return CHAT_CLIENT_SYNC_ANDSEND
		} else {
			continue
		}

	}

}

func makeAckMessage(req *GeneralMessage) *GeneralMessage {

	ack := &GeneralMessage{}
	*ack = *req

	tmpReceiverId := ack.ReceiverId
	ack.ReceiverId = ack.SenderId
	ack.SenderId = tmpReceiverId
	ack.MessageType = CHAT_LOCALPOST_ACK
	return ack
}

func ProcessLocal_to_Local_Message(handler ProcessHandler, req *GeneralMessage) int {

	receiverId := req.ReceiverId
	senderId := req.SenderId
	needsync := false
	needstorage := false
	for {
		handler.Watch1(GetUserInfoName(receiverId))
		user := handler.GetUserInfo(receiverId)

		if CheckResponOnlyLocal_to_Local(user, req) {

			ack := makeAckMessage(req)
			handler.SendRequest(ack)

			handler.UnWatch1(GetUserInfoName(receiverId))
			return CHAT_LOCAL_RESPONSEONLY

		}

		if CheckStoreagedLocal_to_Local(user, req) {
			handler.UnWatch1(GetUserInfoName(receiverId))
			return CHAT_LOCAL_DROP
		}

		needstorage, needsync = StoreageLocal_to_Local(user, req)

		if needstorage {
			// store the message to map
			key := GetLocalMessage(senderId, req.Chatmessage.SendId)
			handler.StoreRequest(receiverId, key, req)

			b, _ := handler.UpdateUser(receiverId, user)
			if !b {
				continue
			} else {
				break
			}

		} else {
			handler.UnWatch1(GetUserInfoName(receiverId))
			break
		}
	}

	if !needsync {
		return CHAT_LOCAL_STOREONLY
	}

	for {
		handler.Watch2(GetUserInfoName(receiverId), GetUserChatInBoxName(receiverId))
		user := handler.GetUserInfo(receiverId)
		if user == nil {
			continue
		}

		list := SyncLocal_to_Local(user, senderId)
		if len(list) == 0 {
			handler.UnWatch2(GetUserInfoName(receiverId), GetUserChatInBoxName(receiverId))

			return CHAT_LOCAL_NOMESSAGENEEDSEND
		}

		reqlist := make([]*GeneralMessage, 0)

		for i := range list {
			key := list[i]
			newReq := handler.GetRequest(receiverId, key)

			reqlist = append(reqlist, newReq)
		}
		b, _ := handler.UpdateUserAndInbox(receiverId, user, reqlist)
		if b {

			for i := range reqlist {
				newReq := reqlist[i]
				newReq = makeAckMessage(newReq)

				handler.SendRequest(newReq)

			}

			return CHAT_LOCAL_MESSAGENEEDSEND
		} else {
			continue
		}

	}
	return CHAT_LOCAL_MESSAGENEEDSEND

}

func ProcessLocal_ack_Message(handler ProcessHandler, req *GeneralMessage) int {

	receiverId := req.ReceiverId
	senderid := req.SenderId

	result := 0

	for {
		handler.Watch1(GetUserInfoName(receiverId))
		user := handler.GetUserInfo(receiverId)

		if user == nil {
			panic(" user == nil")
			continue
		}

		result = CoreLocal_Ack(user, req)
		//CheckCoreLocal_Ack(result)

		succ, _ := handler.UpdateUser(receiverId, user)
		if !succ {
			fmt.Println("handler.UpdateUser(receiverId, user)")
			continue
		} else {
			break
		}
	}

	if result == LOCALACK_SUCCESS_NOSYNC {
		return CHAT_ACK_NOSYNC
	}

	var user *UserInfor
	for {
		handler.Watch1(GetUserInfoName(receiverId))
		user = handler.GetUserInfo(receiverId)

		if user == nil {
			panic(" user == nil")
			continue
		}

		result := SyncLocal_Ack(user, senderid)
		//CheckSyncLocal_Ack(result)

		if result == LOCALACK_SYNC_SUCCESS_NOSYNC {
			handler.UnWatch1(GetUserInfoName(receiverId))
			break
		} else if result == LOCALACK_SYNC_SUCCESS {
			succ, _ := handler.UpdateUser(receiverId, user)
			if !succ {
				fmt.Println("succ, _ := handler.UpdateUser(receiverId, user)")
				continue
			} else {
				break
			}
		}
	}
	return CHAT_ACK_SYNC

}
