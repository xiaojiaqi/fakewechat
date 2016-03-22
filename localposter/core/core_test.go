package core

import (
	. "github.com/fakewechat/lib/contstant"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/localposter/handler"
	. "github.com/fakewechat/message"
	"testing"

	"fmt"
	"math/rand"
	"time"
)

var length uint64
var userid, friendid uint64
var initSendId, friendSendId uint64

var s1 rand.Source
var r1 *rand.Rand

var mhandler, mhandler1, mhandler2 *MockHandler
var times int
var u1, u9 *UserInfor

func InitSystem() {
	length = 10000
	s1 = rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
	times = 30
	userid = 100
	friendid = 200

	initSendId = uint64(r1.Intn(int(length) + 100))
	friendSendId = uint64(r1.Intn(int(length) + 100))
	userid = uint64(r1.Intn(int(length) + 100))
	friendid = uint64(r1.Intn(int(length) + 100))

	if userid == friendid {
		friendid += 10
	}

	mhandler = nil
	mhandler = &MockHandler{}
	mhandler.Init()

	mhandler1 = nil
	mhandler1 = &MockHandler{}
	mhandler1.Init()

	mhandler2 = nil
	mhandler2 = &MockHandler{}
	mhandler2.Init()
	//fmt.Println("InitSystem()")
	u1 = nil
	u9 = nil

}
func init() {
	InitSystem()
}

func makeuser() *UserInfor {

	user := &UserInfor{}
	user.UserMap = make(map[uint64]*User)
	user.SendedMessage = make(map[uint64]uint64)    // `protobuf:"bytes,8,rep,name=SendedMessage" json:"SendedMessage,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	user.LocalMessage = make(map[uint64]*RecvQueue) //`protobuf:"bytes,9,rep,name=localMessage" json:"localMessage,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	user.AckMessage = make(map[uint64]uint64)       //`protobuf:"bytes,10,rep,name=ackMessage" json:"ackMessage,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	return user
}

func Test_Client_to_local_01(t *testing.T) {
	InitSystem()
	user := makeuser() //userid = 1
	user.SendAckId = 100

	user.UserMap = make(map[uint64]*User)

	onefriend := &User{}
	onefriend.UserId = friendid
	onefriend.SendId = 0
	onefriend.ReceiveId = 100
	user.UserMap[friendid] = onefriend
	mhandler.Users[userid] = user

	req := &GeneralMessage{}
	req.SendId = 33
	req.ReceiverId = friendid
	req.SenderId = userid
	req.MessageType = CHAT_CLIENT_TO_LOCALPOST
	chat := &ChatMessage{}
	chat.ReceiverId = friendid
	chat.SenderId = userid
	chat.SendId = 0

	req.Chatmessage = chat

	result := ProcessClient_to_Local_Message(mhandler, req)
	if result != CHAT_CLIENT_DROPED && mhandler.Check() {

		t.Error("ProcessClient_to_Local_Message   ")
	}

}

func Test_Client_to_local_02(t *testing.T) {

	InitSystem()

	user := makeuser() //userid = 1
	user.SendAckId = 40
	user.SendId = 80
	user.UserMap = make(map[uint64]*User)

	onefriend := &User{}
	onefriend.UserId = friendid
	onefriend.SendId = 0
	onefriend.ReceiveId = 100
	user.UserMap[friendid] = onefriend
	mhandler.Users[userid] = user

	req := &GeneralMessage{}
	req.SendId = 60
	req.ReceiverId = friendid
	req.SenderId = userid
	req.MessageType = CHAT_CLIENT_TO_LOCALPOST
	chat := &ChatMessage{}
	chat.ReceiverId = friendid
	chat.SenderId = userid
	chat.SendId = 0
	req.Chatmessage = chat

	mhandler.StoreMessages[GetMockStoreKey(userid, GetSendMessage(60))] = req

	result := ProcessClient_to_Local_Message(mhandler, req)
	if result != CHAT_CLIENT_RESEND && mhandler.Check() {

		t.Error("ProcessClient_to_Local_Message   ")
	}
}

func Test_Client_to_local_03(t *testing.T) {
	InitSystem()

	user := makeuser() //userid = 1

	user.SendAckId = 40
	user.SendId = 100
	user.UserMap = make(map[uint64]*User)

	onefriend := &User{}

	onefriend.UserId = friendid
	onefriend.SendId = 0
	onefriend.ReceiveId = 100
	user.UserMap[friendid] = onefriend
	mhandler.Users[userid] = user

	user.SendedMessage[120] = friendid

	req := &GeneralMessage{}
	req.SendId = 120
	req.ReceiverId = friendid
	req.SenderId = userid
	req.MessageType = CHAT_CLIENT_TO_LOCALPOST
	chat := &ChatMessage{}
	chat.ReceiverId = friendid
	chat.SenderId = userid
	chat.SendId = 0
	req.Chatmessage = chat

	//mhandler.StoreMessages[GetSendMessage(60)] = req

	result := ProcessClient_to_Local_Message(mhandler, req)
	if result != CHAT_CLIENT_NONEEDSTORE && mhandler.Check() {

		t.Error("ProcessClient_to_Local_Message   ")
	}
}

func Test_Client_to_local_04(t *testing.T) {
	InitSystem()

	user := makeuser() //userid = 1

	user.SendAckId = 40
	user.SendId = 100
	user.UserMap = make(map[uint64]*User)

	onefriend := &User{}

	onefriend.UserId = friendid
	onefriend.SendId = 0
	onefriend.ReceiveId = 100
	user.UserMap[friendid] = onefriend
	mhandler.Users[userid] = user

	req := &GeneralMessage{}
	req.SendId = 120
	req.ReceiverId = friendid
	req.SenderId = userid
	req.MessageType = CHAT_CLIENT_TO_LOCALPOST
	chat := &ChatMessage{}
	chat.ReceiverId = friendid
	chat.SenderId = userid
	chat.SendId = 0
	req.Chatmessage = chat

	//mhandler.StoreMessages[GetSendMessage(60)] = req

	result := ProcessClient_to_Local_Message(mhandler, req)
	if result != CHAT_CLIENT_STORE_NOSYNC && mhandler.Check() {

		t.Error("ProcessClient_to_Local_Message   ")
	}
}

// sync
func Test_Client_to_local_05(t *testing.T) {
	InitSystem()
	user := makeuser() //userid = 1

	user.SendAckId = 100
	user.SendId = 100
	user.UserMap = make(map[uint64]*User)

	onefriend := &User{}

	onefriend.UserId = friendid
	onefriend.SendId = 0
	onefriend.ReceiveId = 100
	user.UserMap[friendid] = onefriend
	mhandler.Users[userid] = user

	req := &GeneralMessage{}
	req.SendId = 101
	req.ReceiverId = friendid
	req.SenderId = userid
	req.MessageType = CHAT_CLIENT_TO_LOCALPOST
	chat := &ChatMessage{}
	chat.ReceiverId = friendid
	chat.SenderId = userid
	chat.SendId = 0
	req.Chatmessage = chat

	//mhandler.StoreMessages[GetSendMessage(60)] = req

	result := ProcessClient_to_Local_Message(mhandler, req)
	if result != CHAT_CLIENT_SYNC_ANDSEND && mhandler.Check() {

		t.Error("ProcessClient_to_Local_Message   ")
	}

}

// random

func Test_Core_001(t *testing.T) {
	InitSystem()

	{
		user := makeuser() //userid = 1
		user.SendId = initSendId
		user.SendAckId = initSendId

		user.UserMap = make(map[uint64]*User)

		onefriend := &User{}
		onefriend.UserId = userid
		onefriend.SendId = friendSendId
		onefriend.ReceiveId = 0

		user.UserMap[friendid] = onefriend

		mhandler1.Users[userid] = user
		u1 = user

	}

	{
		user := makeuser() //userid = 9
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)

		onefriend := &User{}
		onefriend.UserId = friendid
		onefriend.SendId = 0
		onefriend.ReceiveId = friendSendId

		user.UserMap[userid] = onefriend

		mhandler2.Users[friendid] = user
		u9 = user
	}

	var id uint64
	//length = 10

	array := make([]*GeneralMessage, 0)
	for k := 0; k < 2; k++ {
		for id = 1; id <= length; id++ {
			req := &GeneralMessage{}
			req.SendId = initSendId + id
			req.ReceiverId = friendid
			req.SenderId = userid

			chat := &ChatMessage{}
			chat.ReceiverId = friendid
			chat.SenderId = userid
			chat.SendId = friendSendId + id

			req.Chatmessage = chat
			array = append(array, req)
		}
	}

	for id = 0; id < length*2; id++ {
		newid := r1.Intn(int(length) - 1)
		req1 := array[id]
		array[id] = array[newid]
		array[newid] = req1
	}

	for id = 0; id < length*2; id++ {
		_ = ProcessClient_to_Local_Message(mhandler1, array[id])
	}

	if len(u1.SendedMessage) != 0 {
		t.Error("len(u1.SendedMessage) != 0  ")

	}

	if u1.SendId != initSendId+length {
		fmt.Println(u1.SendId)
		t.Error("u1.SendId != initSendId + length   ")
	}

	if len(u1.SendedMessage) != 0 {
		t.Error("len(u1.SendedMessage) != 0 ")
	}

	myfriend, ok := u1.UserMap[friendid]

	if !ok {
		t.Error("maps, ok := u1.UserMap[friendid]   ")
	}
	if myfriend.SendId != friendSendId+length {
		t.Error("myfriend.SendId != friendSendId + length  ")
	}

	if uint64(len(mhandler1.Send_to_Local)) != 0 {
		t.Error("len(mhandler1.Send_to_Local) != 0 ")
	}

	if uint64(len(mhandler1.Local_to_Local)) < length {
		t.Error("uint64(len(mhandler1.Local_to_Local)) < length ")
	}

	var i uint64
	/*
		for i = 0; i < length; i++ {
			if (mhandler1.Local_to_Local[i].SendId != initSendId+uint64(i)+uint64(1)) &&
				(mhandler1.Local_to_Local[i].Chatmessage.SendId != friendSendId+uint64(i)+uint64(1)) {
				t.Error("mhandler1.Send_to_Local[i].SendId != initSendId + uint64(i)  +1 ")
			}
		}
	*/

	for i = 1; i <= length; i++ {
		key := GetRawMessage(i + initSendId)
		newkey := GetMockStoreKey(userid, key)

		_, ok := mhandler1.StoreMessages[newkey]
		if !ok {
			t.Error("mhandler1.StoreMessages[newkey] ", newkey)
		}

		key = GetSendMessage(i + initSendId)
		newkey = GetMockStoreKey(userid, key)
		_, ok = mhandler1.StoreMessages[newkey]
		if !ok {
			t.Error("mhandler1.StoreMessages[newkey] ", newkey)
		}
	}

	if len(mhandler1.StoreMessages) != int(2*length) {
		t.Error("len(mhandler1.StoreMessages) != 2* length ")
	}

	for i := range mhandler1.Local_to_Local {
		req := mhandler1.Local_to_Local[i]

		_ = ProcessLocal_to_Local_Message(mhandler2, req)
		//fmt.Println(r)
	}

	//fmt.Println(u9)
	//fmt.Println(mhandler2.Local_ack)
	if uint64(len(mhandler2.Local_ack)) < length {

		t.Error("len(mhandler2.Local_ack) < length ")
	}
	/*
		for i = 0; i < length; i++ {
			//fmt.Println("i = ", i, mhandler1.Local_to_Local[i])
			if (mhandler2.Local_ack[i].SendId != initSendId+uint64(i)+uint64(1)) &&
				(mhandler2.Local_ack[i].Chatmessage.SendId != friendSendId+uint64(i)+uint64(1)) {
				t.Error("mhandler2.Send_to_Local[i].SendId != initSendId + uint64(i)  +1 ")
			}
		}
	*/

	for i := range mhandler2.Local_ack {
		req := mhandler2.Local_ack[i]

		_ = ProcessLocal_ack_Message(mhandler1, req)
		//r = 0
		//fmt.Println("ProcessLocal_ack_Message: ", r)
	}

	if u1.SendAckId != initSendId+length {
		fmt.Println(u1, u1.SendAckId, initSendId+length)
		t.Error("u1.Sendackid != initSendId +length ")
	}

	if len(u1.SendedMessage) != 0 {
		t.Error("len( u1.SendedMessage) != 0 ")
	}

	if len(u1.AckMessage) != 0 {
		fmt.Println(u1.AckMessage)
		t.Error("len( u1.AckMessage) != 0")
	}

	if len(u9.LocalMessage[userid].MessageMap) != 0 {
		t.Error("len(u9.LocalMessage[userid]) != 0  ")

	}

	/*
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
	*/
}

func Test_100time(t *testing.T) {

	for i := 0; i < times; i++ {
		InitSystem()
		Test_Local_to_local_001(t)
		Test_Core_001(t)
	}
}

func Test_Local_to_local_001(t *testing.T) {
	InitSystem()

	{
		user := makeuser() //userid = 9
		user.SendId = 0
		user.SendAckId = 0

		user.UserMap = make(map[uint64]*User)

		onefriend := &User{}
		onefriend.UserId = friendid
		onefriend.SendId = 0
		onefriend.ReceiveId = initSendId

		user.UserMap[userid] = onefriend

		mhandler2.Users[friendid] = user
		u9 = user
	}
	var id uint64
	array := make([]*GeneralMessage, 0)
	for id = 1; id <= length; id++ {
		req := &GeneralMessage{}
		req.SendId = initSendId + id
		req.ReceiverId = friendid
		req.SenderId = userid
		req.MessageType = CHAT_LOCALPOST_TO_LOCALPOST

		chat := &ChatMessage{}
		chat.ReceiverId = friendid
		chat.SenderId = userid
		chat.SendId = initSendId + id

		req.Chatmessage = chat
		array = append(array, req)
	}

	for id = 0; id < length; id++ {
		newid := r1.Intn(int(length) - 1)
		req1 := array[id]
		array[id] = array[newid]
		array[newid] = req1
	}

	for id = 0; id < length; id++ {
		//fmt.Println(array[id])
		_ = ProcessLocal_to_Local_Message(mhandler2, array[id])
		//fmt.Println(u9)
		//fmt.Println(r)
	}
	//fmt.Println(u9)
	if u9.UserMap[userid].ReceiveId != initSendId+length && mhandler2.Check() {
		fmt.Println(u9.UserMap[userid].ReceiveId, initSendId+length)
		t.Error("u9.UserMap[userid].ReceiveId!= initSendId + length  && mhandler2.Check()  ")
	}
	if len(u9.LocalMessage[userid].MessageMap) != 0 {
		t.Error("len(u9.LocalMessage[userid].MessageMap) != 0")
	}
	//fmt.Println(length, len(mhandler2.Local_ack), &mhandler2)

	for id = 0; id < uint64(len(mhandler2.Local_ack)); id++ {
		if mhandler2.Local_ack[id].SendId != initSendId+id+1 {
			fmt.Println(id, mhandler2.Local_ack[id].SendId, initSendId+id+1)
			t.Error("mhandler2.GeneralMessage[i].SendId != initSendId")
		}

	}

}

func Test_Local_ack_001(t *testing.T) {
	InitSystem()

	{
		user := makeuser() //userid = 1
		user.SendId = initSendId
		user.SendAckId = initSendId - 5

		user.UserMap = make(map[uint64]*User)

		onefriend := &User{}
		onefriend.UserId = userid
		onefriend.SendId = friendSendId
		onefriend.ReceiveId = 0

		user.UserMap[friendid] = onefriend
		mhandler1.Users[userid] = user

		u1 = user
		user.AckMessage[initSendId-0] = initSendId - 0
		user.AckMessage[initSendId-1] = initSendId - 1
		user.AckMessage[initSendId-2] = initSendId - 2

		user.AckMessage[initSendId-3] = initSendId - 3
		//user.GetAckMessage[initSendId - 4] = initSendId - 4

	}

	{
		user := makeuser() //userid = 9
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)

		onefriend := &User{}
		onefriend.UserId = friendid
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[userid] = onefriend

		mhandler2.Users[friendid] = user
		u9 = user
	}

	req := &GeneralMessage{}
	req.SendId = initSendId
	req.ReceiverId = userid
	req.SenderId = friendid

	chat := &ChatMessage{}
	chat.ReceiverId = userid
	chat.SenderId = friendid
	chat.SendId = initSendId

	ProcessLocal_ack_Message(mhandler1, req)

	req = &GeneralMessage{}
	req.SendId = initSendId - 4
	req.ReceiverId = userid
	req.SenderId = friendid

	chat = &ChatMessage{}
	chat.ReceiverId = userid
	chat.SenderId = friendid
	chat.SendId = initSendId

	ProcessLocal_ack_Message(mhandler1, req)
	if u1.SendAckId != initSendId && len(u1.AckMessage) != 0 && mhandler1.Check() {
		t.Error("u1.SendAckId != initSendId && len(u1.AckMessage) != 0  ")
	}
	if u9.SendAckId != 0 {

	}
}

func Test_Local_ack_002(t *testing.T) {
	InitSystem()
	{
		user := makeuser() //userid = 1
		user.SendId = initSendId
		user.SendAckId = initSendId - length

		user.UserMap = make(map[uint64]*User)

		onefriend := &User{}
		onefriend.UserId = userid
		onefriend.SendId = friendSendId
		onefriend.ReceiveId = 0

		user.UserMap[friendid] = onefriend
		mhandler1.Users[userid] = user

		u1 = user

		//user.GetAckMessage[initSendId - 4] = initSendId - 4

	}

	{
		user := makeuser() //userid = 9
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)

		onefriend := &User{}
		onefriend.UserId = friendid
		onefriend.SendId = 0
		onefriend.ReceiveId = 0

		user.UserMap[userid] = onefriend

		mhandler2.Users[friendid] = user
		u9 = user
	}
	var id uint64
	array := make([]*GeneralMessage, 0)
	for id = 0; id < length; id++ {
		req := &GeneralMessage{}
		req.SendId = initSendId - uint64(id)
		req.ReceiverId = userid
		req.SenderId = friendid

		chat := &ChatMessage{}
		chat.ReceiverId = userid
		chat.SenderId = friendid
		chat.SendId = initSendId

		req.Chatmessage = chat
		array = append(array, req)
	}

	for id = 0; id < length; id++ {
		newid := r1.Intn(int(length) - 1)
		req1 := array[id]
		array[id] = array[newid]
		array[newid] = req1
	}
	for id = 0; id < length; id++ {
		ProcessLocal_ack_Message(mhandler1, array[id])
	}
	if u1.SendAckId != initSendId && len(u1.AckMessage) != 0 && mhandler1.Check() {
		t.Error("u1.SendAckId != initSendId && len(u1.AckMessage) != 0  ")
	}
	if u9.SendAckId != 0 {

	}
}

/*
func Test_Core_002(t *testing.T) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	mhandler := &MockHandler{}

	mhandler.Init()
	{
		user := &UserInfor{} //userid = 1
		user.SendId = 0

		user.UserMap = make(map[uint64]*User)

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
*/
