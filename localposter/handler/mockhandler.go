package handler

import (
	"fmt"
	. "github.com/fakewechat/lib/contstant"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/message"
)

type MockHandler struct {
	Users  map[uint64]*UserInfor
	Inbox  []*GeneralMessage
	Outbox []*GeneralMessage

	Send_to_Local  []*GeneralMessage
	Local_to_Local []*GeneralMessage
	Local_ack      []*GeneralMessage

	StoreMessages map[string]*GeneralMessage

	Watchstr string
}

func (handler *MockHandler) Init() {
	handler.Users = make(map[uint64]*UserInfor)
	handler.Inbox = make([]*GeneralMessage, 0)
	handler.Outbox = make([]*GeneralMessage, 0)
	handler.Send_to_Local = make([]*GeneralMessage, 0)
	handler.Local_to_Local = make([]*GeneralMessage, 0)
	handler.Local_ack = make([]*GeneralMessage, 0)
	handler.StoreMessages = make(map[string]*GeneralMessage)
}

func (handler *MockHandler) Check() bool {
	if handler.Watchstr != "" {
		return false
	}
	return true
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

func (handler *MockHandler) UpdateUser(userid uint64, u *UserInfor) (bool, error) {
	handler.Watchstr = ""
	return true, nil

}

func (handler *MockHandler) UpdateUserAndInbox(userid uint64, u *UserInfor, Req []*GeneralMessage) (bool, error) {

	for i := range Req {
		v := &GeneralMessage{}
		*v = *Req[i]
		//fmt.Println(v, userid)
		if v.MessageType != CHAT_LOCALPOST_TO_LOCALPOST {
			panic("")
		}
		if v.ReceiverId != userid {
			panic("")
		}

		handler.Inbox = append(handler.Inbox, v)

	}
	handler.Watchstr = ""
	return true, nil
}

func (handler *MockHandler) UpdateUserAndOutbox(userid uint64, u *UserInfor, Req []*GeneralMessage) (bool, error) {

	/*
		v := &GeneralMessage{}
		*v = *Req
		handler.Outbox = append(handler.Outbox, v)
	*/
	for i := range Req {
		v := &GeneralMessage{}
		*v = *Req[i]
		//fmt.Println(v, userid)
		if v.MessageType != CHAT_LOCALPOST_TO_LOCALPOST {
			fmt.Println(v)
			panic("")
		}
		if v.SenderId != userid {
			panic("")

		}
		handler.Outbox = append(handler.Outbox, v)

	}
	handler.Watchstr = ""
	return true, nil
}

func (handler *MockHandler) SendRequest(Req *GeneralMessage) {
	if Req.MessageType == CHAT_CLIENT_TO_LOCALPOST {
		handler.Send_to_Local = append(handler.Send_to_Local, Req)

	} else if Req.MessageType == CHAT_LOCALPOST_TO_LOCALPOST {
		//fmt.Println("handler.Local_to_Local = ", len (handler.Local_to_Local))
		handler.Local_to_Local = append(handler.Local_to_Local, Req)

	} else if Req.MessageType == CHAT_LOCALPOST_ACK {
		//fmt.Println("handler.Local_ack = ", len (handler.Local_ack))
		handler.Local_ack = append(handler.Local_ack, Req)
	} else {
		panic("wrong type")
	}
	//fmt.Println("sendout ", Req)
}

func (handler *MockHandler) Watch1(id string) {

	if handler.Watchstr != "" {
		fmt.Println("watch bug", handler.Watchstr, id)
		panic("wrong watch ")
	}
	handler.Watchstr = id
}

func (handler *MockHandler) Watch2(id1 string, id2 string) {
	if handler.Watchstr != "" {
		panic("wrong watch ")
	}
	handler.Watchstr = id1 + "," + id2
}

func (handler *MockHandler) UnWatch1(id string) {
	if handler.Watchstr == "" {
		panic("wrong unwatch ")
	}
	if handler.Watchstr != id {
		panic("wrong unwatch ")
	}
	handler.Watchstr = ""
}

func (handler *MockHandler) UnWatch2(id1 string, id2 string) {
	if handler.Watchstr == "" {
		panic("wrong unwatch ")
	}
	if handler.Watchstr != id1+","+id2 {
		panic("wrong unwatch ")
	}
	handler.Watchstr = ""
}

func (handler *MockHandler) GetRequest(userid uint64, key string) (Req *GeneralMessage) {
	newkey := GetMockStoreKey(userid, key)
	r, ok := handler.StoreMessages[newkey]
	if !ok {
		fmt.Println("want to get ", newkey)
		panic("(handler *MockHandler)GetRequest(userid uint64, key string)")
	}
	//fmt.Println("get:", newkey);

	return r
}
func (handler *MockHandler) StoreRequest(userid uint64, key string, Req *GeneralMessage) {
	newkey := GetMockStoreKey(userid, key)
	//fmt.Println("storage:", newkey);

	handler.StoreMessages[newkey] = Req
}
