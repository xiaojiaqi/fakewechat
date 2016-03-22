package chatmessage

import (
	"fmt"
	. "github.com/fakewechat/message"
	"testing"
	//"math/rand"
	//"time"
)

//
// normal test
//

func Test_CheckResponOnlyLocal_to_Local(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 100
	user.SendAckId = 100
	var senderid uint64
	senderid = 5
	user.UserMap = make(map[uint64]*User)
	user.SendedMessage = make(map[uint64]uint64)
	user.LocalMessage = make(map[uint64]*RecvQueue)

	onefriend := &User{}
	onefriend.UserId = Receiverid
	onefriend.SendId = 30
	onefriend.ReceiveId = 100

	user.UserMap[senderid] = onefriend

	req := &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat := &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 29

	req.Chatmessage = chat

	result := CheckResponOnlyLocal_to_Local(user, req)

	if result != true {
		t.Error("CheckResponOnlyLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 111

	req.Chatmessage = chat

	result = CheckResponOnlyLocal_to_Local(user, req)

	if result != false {
		t.Error("CheckResponOnlyLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = 1234

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 111

	req.Chatmessage = chat

	defer func() {
		if r := recover(); r != nil {

		} else {
			t.Error("CheckResponOnlyLocal_to_Local   ")
		}
	}()

	result = CheckResponOnlyLocal_to_Local(user, req)

	if result != false {
		t.Error("CheckResponOnlyLocal_to_Local   ")
	}
}

func Test_CheckStoreagedLocal_to_Local(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 100
	user.SendAckId = 100
	var senderid uint64
	senderid = 5

	user.UserMap = make(map[uint64]*User)
	user.SendedMessage = make(map[uint64]uint64)
	user.LocalMessage = make(map[uint64]*RecvQueue)

	onefriend := &User{}
	onefriend.UserId = Receiverid
	onefriend.SendId = 30
	onefriend.ReceiveId = 100

	user.UserMap[senderid] = onefriend

	req := &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat := &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 101

	req.Chatmessage = chat

	result := CheckStoreagedLocal_to_Local(user, req)

	if result != false {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 101

	req.Chatmessage = chat

	result = CheckStoreagedLocal_to_Local(user, req)

	if result != false {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 101

	req.Chatmessage = chat

	r := &RecvQueue{}
	r.MessageMap = make(map[uint64]uint64)
	r.MessageMap[101] = 101

	user.LocalMessage[senderid] = r

	result = CheckStoreagedLocal_to_Local(user, req)

	if result != true {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 103

	req.Chatmessage = chat

	result = CheckStoreagedLocal_to_Local(user, req)

	if result != false {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = 1234

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 101

	req.Chatmessage = chat

	defer func() {
		if r := recover(); r != nil {

		} else {
			t.Error("StoreageLocal_to_Local   ")
		}
	}()
	result = CheckStoreagedLocal_to_Local(user, req)

	if result != false {
		t.Error("StoreageLocal_to_Local   ")
	}
}

func Test_StoreageLocal_to_Local(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 100
	user.SendAckId = 100
	var senderid uint64
	senderid = 5
	user.UserMap = make(map[uint64]*User)
	user.SendedMessage = make(map[uint64]uint64)
	user.LocalMessage = make(map[uint64]*RecvQueue)

	onefriend := &User{}
	onefriend.UserId = Receiverid
	onefriend.SendId = 30
	onefriend.ReceiveId = 100

	user.UserMap[senderid] = onefriend

	req := &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat := &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 101

	req.Chatmessage = chat

	storage, sync := StoreageLocal_to_Local(user, req)

	if storage != true && sync != false {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 101

	req.Chatmessage = chat

	storage, sync = StoreageLocal_to_Local(user, req)

	if storage != false && sync != false {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 102

	req.Chatmessage = chat

	storage, sync = StoreageLocal_to_Local(user, req)

	if storage != true && sync != false {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = senderid

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 103

	req.Chatmessage = chat

	storage, sync = StoreageLocal_to_Local(user, req)

	if storage != true && sync != false {
		t.Error("StoreageLocal_to_Local   ")
	}

	req = &GeneralMessage{}
	req.SendId = 1011
	req.ReceiverId = Receiverid
	req.SenderId = 1234

	chat = &ChatMessage{}
	chat.ReceiverId = Receiverid
	chat.SendId = 101

	req.Chatmessage = chat

	defer func() {
		if r := recover(); r != nil {

		} else {
			t.Error("StoreageLocal_to_Local   ")
		}
	}()
	storage, sync = StoreageLocal_to_Local(user, req)

	if storage != false && sync != false {
		t.Error("StoreageLocal_to_Local   ")
	}

}

func Test_SyncLocal_to_Local(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 100
	user.SendAckId = 100
	var senderid, senderid2, senderid3 uint64
	senderid = 5
	senderid2 = 6
	senderid3 = 7

	user.UserMap = make(map[uint64]*User)
	user.SendedMessage = make(map[uint64]uint64)
	user.LocalMessage = make(map[uint64]*RecvQueue)

	onefriend := &User{}
	onefriend.UserId = Receiverid
	onefriend.SendId = 30
	onefriend.ReceiveId = 55

	user.UserMap[senderid] = onefriend

	onefriend2 := &User{}
	onefriend2.UserId = Receiverid
	onefriend2.SendId = 30
	onefriend2.ReceiveId = 55

	user.UserMap[senderid2] = onefriend2

	r := &RecvQueue{}
	r.MessageMap = make(map[uint64]uint64)
	r.MessageMap[56] = 56
	r.MessageMap[57] = 57
	r.MessageMap[58] = 58
	r.MessageMap[59] = 59

	user.LocalMessage[senderid] = r

	list := SyncLocal_to_Local(user, senderid)
	fmt.Println(list)
	if len(list) != 4 {
		t.Error("Test_SyncLocal_to_Local   ")
	}
	if onefriend.ReceiveId != 59 {
		t.Error("Test_SyncLocal_to_Local   ")
	}
	if len(r.MessageMap) != 0 {
		t.Error("Test_SyncLocal_to_Local   ")
	}

	list = SyncLocal_to_Local(user, senderid2)
	fmt.Println(list)
	if len(list) != 0 {
		t.Error("Test_SyncLocal_to_Local   ")
	}
	if onefriend.ReceiveId != 59 {
		t.Error("Test_SyncLocal_to_Local   ")
	}
	if len(r.MessageMap) != 0 {
		t.Error("Test_SyncLocal_to_Local   ")
	}

	defer func() {
		if r := recover(); r != nil {

		} else {
			t.Error("Test_SyncLocal_to_Local   ")
		}
	}()

	list = SyncLocal_to_Local(user, senderid3)
}
