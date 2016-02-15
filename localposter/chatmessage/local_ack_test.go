package chatmessage

import (
	. "github.com/fakewechat/message"
	"testing"
)

//
// normal test
//
func Test_ack100(t *testing.T) {

	user := &UserInfor{}
	user.SendAckId = 100
	user.SendId = 101

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 30

	user.UserMap[1] = onefriend

	req := &GeneralMessage{}
	req.SendId = 101
	req.ReceiverId = 1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result := CoreLocal_Ack(user, req)
	if result != LOCALACK_SUCCESS {
		t.Error("Test_ack100  send 100 error, result")
	}

}
