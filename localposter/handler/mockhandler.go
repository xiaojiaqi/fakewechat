package handler

import (
	"fmt"
	. "github.com/fakewechat/lib/contstant"
	. "github.com/fakewechat/message"
)

type MockHandler struct {
	Users  map[uint64]*UserInfor
	Inbox  []*GeneralMessage
	Outbox []*GeneralMessage

	Send_to_Local  []*GeneralMessage
	Local_to_Local []*GeneralMessage
	Local_ack      []*GeneralMessage
}

func (handler *MockHandler) Init() {
	handler.Users = make(map[uint64]*UserInfor)
	handler.Inbox = make([]*GeneralMessage, 0)
	handler.Outbox = make([]*GeneralMessage, 0)
	handler.Send_to_Local = make([]*GeneralMessage, 0)
	handler.Local_to_Local = make([]*GeneralMessage, 0)
	handler.Local_ack = make([]*GeneralMessage, 0)

}

func (handler *MockHandler) GetUserInfo(userid uint64) *UserInfor {

	u, ok := handler.Users[userid]
	if ok {
		return u
	} else {
		fmt.Println("can't found such user ", userid)
		panic("no such user")
		return nil
	}
}

func (handler *MockHandler) UpdateUser(userid uint64, u *UserInfor) bool {

	return true

}

func (handler *MockHandler) UpdateUserAndInbox(userid uint64, u *UserInfor, Req *GeneralMessage) bool {

	v := &GeneralMessage{}
	*v = *Req
	handler.Inbox = append(handler.Inbox, v)
	return true
}

func (handler *MockHandler) UpdateUserAndOutbox(userid uint64, u *UserInfor, Req *GeneralMessage) bool {

	v := &GeneralMessage{}
	*v = *Req
	handler.Outbox = append(handler.Outbox, v)
	return true
}

func (handler *MockHandler) SendRequest(Req *GeneralMessage) {
	if Req.MessageType == CHAT_CLIENT_TO_LOCALPOST {
		handler.Send_to_Local = append(handler.Send_to_Local, Req)

	} else if Req.MessageType == CHAT_LOCALPOST_TO_LOCALPOST {
		handler.Local_to_Local = append(handler.Local_to_Local, Req)

	} else if Req.MessageType == CHAT_LOCALPOST_ACK {
		handler.Local_ack = append(handler.Local_ack, Req)
	} else {
		panic("wrong type")
	}
}

func (handler *MockHandler) Watch(id string) {

}

func (handler *MockHandler) Watch2(id1 string, id2 string) {

}
