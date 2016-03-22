package chatmessage

import (
	. "github.com/fakewechat/message"
	"testing"
)

//
// normal test
//
func Test_CoreLocal_Ack(t *testing.T) {

	user := &UserInfor{}
	user.SendAckId = 100
	user.SendId = 105

	var senderid1, senderid2, senderid3, senderid4 uint64
	senderid1 = 30
	senderid2 = 40
	senderid3 = 50
	senderid4 = 60

	user.UserMap = make(map[uint64]*User)
	user.AckMessage = make(map[uint64]uint64)

	onefriend1 := &User{}
	onefriend1.UserId = senderid1
	onefriend1.SendId = 111

	onefriend2 := &User{}
	onefriend2.UserId = senderid2
	onefriend2.SendId = 30

	onefriend3 := &User{}
	onefriend3.UserId = senderid3
	onefriend3.SendId = 30

	user.UserMap[senderid1] = onefriend1
	user.UserMap[senderid2] = onefriend2
	user.UserMap[senderid3] = onefriend3

	req := &GeneralMessage{}
	req.SendId = 101
	req.ReceiverId = 1
	req.SenderId = senderid1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result := CoreLocal_Ack(user, req)
	if result != LOCALACK_SUCCESS {
		t.Error("Test_CoreLocal_Ack result")
	}

	req = &GeneralMessage{}
	req.SendId = 101
	req.ReceiverId = 1
	req.SenderId = senderid1

	chat = &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result = CoreLocal_Ack(user, req)
	if result != LOCALACK_SUCCESS_NOSYNC {
		t.Error("Test_CoreLocal_Ack ")
	}

	req = &GeneralMessage{}
	req.SendId = 103
	req.ReceiverId = 1
	req.SenderId = senderid1

	chat = &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result = CoreLocal_Ack(user, req)
	if result != LOCALACK_SUCCESS_NOSYNC {
		t.Error("Test_CoreLocal_Ack ")
	}

	req = &GeneralMessage{}
	req.SendId = 103
	req.ReceiverId = 1
	req.SenderId = senderid1

	chat = &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result = CoreLocal_Ack(user, req)
	if result != LOCALACK_SUCCESS_NOSYNC {
		t.Error("Test_CoreLocal_Ack ")
	}

	req = &GeneralMessage{}
	req.SendId = 103
	req.ReceiverId = 1
	req.SenderId = senderid4

	chat = &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	defer func() {
		if r := recover(); r != nil {

		} else {
			t.Error("Test_SyncLocal_to_Local   ")
		}
	}()

	result = CoreLocal_Ack(user, req)
	if result != LOCALACK_SUCCESS_NOSYNC {
		t.Error("Test_CoreLocal_Ack")
	}
}

func Test_SyncLocal_Ack(t *testing.T) {

	user := &UserInfor{}
	user.SendAckId = 100
	user.SendId = 110

	var senderid1, senderid2, senderid3 uint64
	senderid1 = 30
	senderid2 = 40
	senderid3 = 50

	user.UserMap = make(map[uint64]*User)
	user.AckMessage = make(map[uint64]uint64)
	user.SendedMessage = make(map[uint64]uint64)

	onefriend1 := &User{}
	onefriend1.UserId = senderid1
	onefriend1.SendId = 111

	onefriend2 := &User{}
	onefriend2.UserId = senderid2
	onefriend2.SendId = 30

	onefriend3 := &User{}
	onefriend3.UserId = senderid3
	onefriend3.SendId = 30

	user.UserMap[senderid1] = onefriend1
	user.UserMap[senderid2] = onefriend2
	user.UserMap[senderid3] = onefriend3

	user.AckMessage[101] = 101
	user.AckMessage[102] = 102
	user.AckMessage[103] = 103
	user.AckMessage[104] = 104
	user.AckMessage[106] = 106
	user.AckMessage[107] = 107

	user.SendedMessage[101] = senderid1
	user.SendedMessage[102] = senderid1
	user.SendedMessage[103] = senderid1
	user.SendedMessage[104] = senderid1
	user.SendedMessage[105] = senderid1

	user.SendedMessage[106] = senderid1
	user.SendedMessage[107] = senderid1

	list := SyncLocal_Ack(user, senderid1)
	if list != LOCALACK_SYNC_SUCCESS && user.SendAckId != 104 && len(user.AckMessage) != 2 && len(user.SendedMessage) != 3 {
		t.Error("Test_SyncLocal_Ack")
	}

	list = SyncLocal_Ack(user, senderid1)
	if list != LOCALACK_SYNC_SUCCESS_NOSYNC && user.SendAckId != 104 && len(user.AckMessage) != 2 && len(user.SendedMessage) != 3 {
		t.Error("Test_SyncLocal_Ack")
	}

	req := &GeneralMessage{}
	req.SendId = 105
	req.ReceiverId = 1
	req.SenderId = senderid1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result := CoreLocal_Ack(user, req)
	if result != LOCALACK_SUCCESS {
		t.Error("Test_CoreLocal_Ack ")
	}

	list = SyncLocal_Ack(user, senderid1)
	if list != LOCALACK_SYNC_SUCCESS && user.SendAckId != 107 && len(user.AckMessage) != 0 && len(user.SendedMessage) != 0 {
		t.Error("Test_SyncLocal_Ack")
	}
}
