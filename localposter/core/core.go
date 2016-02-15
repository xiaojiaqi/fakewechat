package core

import (
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
func ProcessClient_to_Local_Message(handler ProcessHandler, req *GeneralMessage) {

	senderId := req.SenderId

	result := 0
	needupdateOutbox := false
	for {
		handler.Watch2(GetUserInfoName(senderId), GetUserChatOutBoxName(senderId))

		user := handler.GetUserInfo(senderId)

		if user == nil {
			continue
		}

		result, needupdateOutbox = CoreClient_to_Local(user, req)
		CheckCoreClient_to_Local(result)

        if result == CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND {
            break
        }
		updated := false
		if needupdateOutbox {
			updated = handler.UpdateUserAndOutbox(senderId, user, req)
		} else {
			updated = handler.UpdateUser(senderId, user)
		}
		if !updated {
			continue
		} else {
			break
		}
	}

	if result == CLIENT_TO_LOCAL_SUCCESS_NOSEND {
		return
	}

	// send message
	req.MessageType = CHAT_LOCALPOST_TO_LOCALPOST
	handler.SendRequest(req)

	for {
		handler.Watch2(GetUserInfoName(senderId), GetUserChatOutBoxName(senderId))
		user := handler.GetUserInfo(senderId)
		if user == nil {
			continue
		}

		result, req := SyncClient_to_Local(user)
		if result == CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND {
			break
		} else if result == CLIENT_TO_LOCAL_SYNC_SUCCESS {
			if !handler.UpdateUserAndOutbox(senderId, user, req) {
				continue
			}
			req.MessageType = CHAT_LOCALPOST_TO_LOCALPOST
			handler.SendRequest(req)
		} else {
			panic("unexpected error")
		}
	}

}

func ProcessLocal_to_Local_Message(handler ProcessHandler, req *GeneralMessage) {

	receiverId := req.ReceiverId
	senderId := req.SenderId

	result := 0
	needupdateInbox := false
	for {
		handler.Watch2(GetUserInfoName(receiverId), GetUserChatInBoxName(receiverId))
		user := handler.GetUserInfo(receiverId)
		if user == nil {
			continue
		}

		result, needupdateInbox = CoreLocal_to_Local(user, req)
		CheckCoreLocal_to_Local(result)
		updated := false

		if needupdateInbox == true {
			updated = handler.UpdateUserAndInbox(receiverId, user, req)
		} else {
			updated = handler.UpdateUser(receiverId, user)
		}

		if !updated {
			continue
		} else {
			break
		}
	}

	if result == LOCAL_TO_LOCAL_SUCCESS_NOSEND {
		return
	}
	// send message
	// req need exchange sender and receiver
	tmpReceiverId := req.ReceiverId
	req.ReceiverId = req.SenderId
	req.SenderId = tmpReceiverId
	req.MessageType = CHAT_LOCALPOST_ACK

	handler.SendRequest(req)
	//
	for {
		handler.Watch2(GetUserInfoName(receiverId), GetUserChatInBoxName(receiverId))
		user := handler.GetUserInfo(receiverId)

		if user == nil {
			continue
		}

		result, req := SyncLocal_to_Local(user, senderId)

		CheckSyncLocal_to_Local(result)

		if result == LOCAL_TO_LOCAL_SYNC_SUCCESS_NOSEND {
			break
		} else if result == LOCAL_TO_LOCAL_SYNC_SUCCESS {
			if !handler.UpdateUserAndInbox(receiverId, user, req) {
				continue
			}

			tmpReceiverId := req.ReceiverId
			req.ReceiverId = req.SenderId
			req.SenderId = tmpReceiverId
			req.MessageType = CHAT_LOCALPOST_ACK
			handler.SendRequest(req)

		} else {
			panic("")
		}
	}
}

func ProcessLocal_ack_Message(handler ProcessHandler, req *GeneralMessage) {

	receiverId := req.ReceiverId
	senderid := req.SenderId

	result := 0
	for {
		handler.Watch(GetUserInfoName(receiverId))
		user := handler.GetUserInfo(receiverId)

		if user == nil {
			continue
		}

		result = CoreLocal_Ack(user, req)

		CheckCoreLocal_Ack(result)
		if !handler.UpdateUser(receiverId, user) {
			continue
		} else {
			break
		}
	}

	if result == LOCALACK_SUCCESS_NOSYNC {
		return
	}

	for {
		handler.Watch(GetUserInfoName(receiverId))
		user := handler.GetUserInfo(receiverId)

		if user == nil {
			continue
		}

		result := SyncLocal_Ack(user, senderid)
		CheckSyncLocal_Ack(result)

		if result == LOCALACK_SYNC_SUCCESS_NOSYNC {
			break
		} else if result == LOCALACK_SYNC_SUCCESS {
			if !handler.UpdateUser(receiverId, user) {
				continue
			} else {
				break
			}
		}
	}

}
