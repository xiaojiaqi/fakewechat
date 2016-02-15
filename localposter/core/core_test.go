package core

import (
	. "github.com/fakewechat/localposter/handler"
	. "github.com/fakewechat/message"
	"testing"

	//"fmt"
	"math/rand"
	"time"
)

func Test_Core_001(t *testing.T) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	mhandler := &MockHandler{}

	mhandler.Init()
	{
		user := &UserInfor{} //userid = 1
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)
		user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
		user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
		user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

		onefriend := &User{}
		onefriend.UserId = 9
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[9] = onefriend

		mhandler.Users[1] = user

	}

	{
		user := &UserInfor{} //userid = 9
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)
		user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
		user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
		user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

		onefriend := &User{}
		onefriend.UserId = 1
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[1] = onefriend

		mhandler.Users[9] = user

	}
	var id, length uint64
	length = 300

	array := make([]*GeneralMessage, 0)
	for id = 0; id < length; id++ {
		req := &GeneralMessage{}
		req.SendId = 1 + id
		req.ReceiverId = 9
		req.SenderId = 1

		chat := &ChatMessage{}
		chat.ReceiverId = 9
		chat.SenderId = 1 + id
		chat.SendId = 1 + id

		req.Chatmessage = chat
		array = append(array, req)
	}

	for id = 0; id < length; id++ {
		newid := r1.Intn(int(length) - 1)
		req1 := array[id]
		array[id] = array[newid]
		array[newid] = req1
	}
	u1 := mhandler.GetUserInfo(1)
	u9 := mhandler.GetUserInfo(9)

	for id = 0; id < length; id++ {
		ProcessClient_to_Local_Message(mhandler, array[id])
	}

	var tmplength uint64
	for i := range mhandler.Local_to_Local {
		//fmt.Println(mhandler.Local_to_Local[i].String())
		array[tmplength] = mhandler.Local_to_Local[i]
		tmplength += 1
	}
	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())

	//fmt.Println("tmplength := ", tmplength)

	for id = 0; id < tmplength; id++ {
		ProcessLocal_to_Local_Message(mhandler, array[id])
	}
	tmplength = 0
	for i := range mhandler.Local_ack {
		//fmt.Println(mhandler.Local_to_Local[i].String())
		array[tmplength] = mhandler.Local_ack[i]
		tmplength += 1
	}
	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())
	//fmt.Println("tmplength := ", tmplength)

	for id = 0; id < tmplength; id++ {
		ProcessLocal_ack_Message(mhandler, array[id])
	}

	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())
	if u1.SendId != u1.SendAckId && u1.SendAckId != length {

		panic("u.SendId != u.SendAckId && u.SendAckId != length")
	}

	if u9.ReceiveId != length {
		panic("u.ReceiveId != length")
	}
	id = 1
	for i := range mhandler.Inbox {
		if mhandler.Inbox[i].SendId != id {
		//	fmt.Println(id, mhandler.Inbox[i].String())
			panic("mhandler.Inbox[i].SendId != id")
		}
		id += 1
	}

	id = 1
	for i := range mhandler.Outbox {
		if mhandler.Outbox[i].SendId != id {
		//	fmt.Println(id, mhandler.Outbox[i].String())
			panic("mhandler.Inbox[i].SendId != id")
		}
		id += 1
	}
 
}

func Test_Core_002(t *testing.T) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	mhandler := &MockHandler{}

	mhandler.Init()
	{
		user := &UserInfor{} //userid = 1
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)
		user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
		user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
		user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

		onefriend := &User{}
		onefriend.UserId = 9
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[9] = onefriend

		mhandler.Users[1] = user

	}

	{
		user := &UserInfor{} //userid = 9
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)
		user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
		user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
		user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

		onefriend := &User{}
		onefriend.UserId = 1
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[1] = onefriend

		mhandler.Users[9] = user

	}
	var id, length, times, j uint64
	length = 400
	times = 3

	array := make([]*GeneralMessage, 0)

	for j = 0; j < times; j++ {
		for id = 0; id < length; id++ {
			req := &GeneralMessage{}
			req.SendId = 1 + id
			req.ReceiverId = 9
			req.SenderId = 1

			chat := &ChatMessage{}
			chat.ReceiverId = 9
			chat.SenderId = 1 + id
			chat.SendId = 1 + id

			req.Chatmessage = chat
			array = append(array, req)
		}
	}

	for id = 0; id < length*times; id++ {
		newid := r1.Intn(int(length*times) - 1)
		req1 := array[id]
		array[id] = array[newid]
		array[newid] = req1
	}
	u1 := mhandler.GetUserInfo(1)
	u9 := mhandler.GetUserInfo(9)

	for id = 0; id < length*times; id++ {
		ProcessClient_to_Local_Message(mhandler, array[id])
	}

	var tmplength uint64
	for i := range mhandler.Local_to_Local {
		//fmt.Println(mhandler.Local_to_Local[i].String())
		array[tmplength] = mhandler.Local_to_Local[i]
		tmplength += 1
	}
	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())

	//fmt.Println("tmplength := ", tmplength)

	for id = 0; id < tmplength; id++ {
		ProcessLocal_to_Local_Message(mhandler, array[id])
	}
	tmplength = 0
	for i := range mhandler.Local_ack {
		//fmt.Println(mhandler.Local_to_Local[i].String())
		array[tmplength] = mhandler.Local_ack[i]
		tmplength += 1
	}
	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())
	//fmt.Println("tmplength := ", tmplength)

	for id = 0; id < tmplength; id++ {
		ProcessLocal_ack_Message(mhandler, array[id])
	}

	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())
	if u1.SendId != u1.SendAckId && u1.SendAckId != length {

		panic("u.SendId != u.SendAckId && u.SendAckId != length")
	}

	if u9.ReceiveId != length {
		panic("u.ReceiveId != length")
	}
	id = 1
	for i := range mhandler.Inbox {
		if mhandler.Inbox[i].SendId != id {
	//		fmt.Println(id, mhandler.Inbox[i].String())
			panic("mhandler.Inbox[i].SendId != id")
		}
		id += 1
	}

	id = 1
	for i := range mhandler.Outbox {
		if mhandler.Outbox[i].SendId != id {
		//	fmt.Println(id, mhandler.Outbox[i].String())
			panic("mhandler.Inbox[i].SendId != id")
		}
		id += 1
	}

}

func Test_Core_003(t *testing.T) {

	mhandler := &MockHandler{}

	mhandler.Init()
	{
		user := &UserInfor{} //userid = 1
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)
		user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
		user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
		user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

		onefriend := &User{}
		onefriend.UserId = 9
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[9] = onefriend

		mhandler.Users[1] = user

	}

	{
		user := &UserInfor{} //userid = 9
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)
		user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
		user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
		user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

		onefriend := &User{}
		onefriend.UserId = 1
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[1] = onefriend

		mhandler.Users[9] = user

	}
	var id, length uint64
	length = 10

	array := make([]*GeneralMessage, 0)
	for id = length; id > 0; id-- {
		req := &GeneralMessage{}
		req.SendId = 1 + id - 1
		req.ReceiverId = 9
		req.SenderId = 1

		chat := &ChatMessage{}
		chat.ReceiverId = 9
		chat.SenderId = 1 + id - 1
		chat.SendId = 1 + id - 1

		req.Chatmessage = chat
		array = append(array, req)
	}

	u1 := mhandler.GetUserInfo(1)
	u9 := mhandler.GetUserInfo(9)

	for id = 0; id < length; id++ {
		ProcessClient_to_Local_Message(mhandler, array[id])
	}

	var tmplength uint64
	for i := range mhandler.Local_to_Local {
		//fmt.Println(mhandler.Local_to_Local[i].String())
		array[tmplength] = mhandler.Local_to_Local[i]
		tmplength += 1
	}
	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())

	//fmt.Println("tmplength := ", tmplength)

	for id = 0; id < tmplength; id++ {
		ProcessLocal_to_Local_Message(mhandler, array[id])
	}
	tmplength = 0
	for i := range mhandler.Local_ack {
		//fmt.Println(mhandler.Local_to_Local[i].String())
		array[tmplength] = mhandler.Local_ack[i]
		tmplength += 1
	}
	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())
	//fmt.Println("tmplength := ", tmplength)

	for id = 0; id < tmplength; id++ {
		ProcessLocal_ack_Message(mhandler, array[id])
	}

	//fmt.Println("u1:= ", u1.String())
	//fmt.Println("u9:= ", u9.String())
	if u1.SendId != u1.SendAckId && u1.SendAckId != length {

		panic("u.SendId != u.SendAckId && u.SendAckId != length")
	}

	if u9.ReceiveId != length {
		panic("u.ReceiveId != length")
	}
	id = 1
	for i := range mhandler.Inbox {
		if mhandler.Inbox[i].SendId != id {
	//		fmt.Println(id, mhandler.Inbox[i].String())
			panic("mhandler.Inbox[i].SendId != id")
		}
		id += 1
	}

	id = 1
	for i := range mhandler.Outbox {
		if mhandler.Outbox[i].SendId != id {
		//	fmt.Println(id, mhandler.Outbox[i].String())
			panic("mhandler.Inbox[i].SendId != id")
		}
		id += 1
	}
 
}
