package chatmessage

import (
	. "github.com/fakewechat/message"
	"math/rand"
	"testing"
	"time"
)

//
// normal test
//

func Test_Local_to_local100(t *testing.T) {

	user := &UserInfor{} //userid = 1
	user.SendId = 300

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 400
	onefriend.SendId = 30
	onefriend.ReceiveId = 60

	user.UserMap[400] = onefriend

	req := &GeneralMessage{} // message sender == 400
	req.SendId = 1101
	req.ReceiverId = 1
	req.SendId = 100
	req.SenderId = 400

	chat := &ChatMessage{}
	chat.ReceiverId = 1
	chat.SendId = 61
	req.Chatmessage = chat

	result, needupdatRecvQueue := CoreLocal_to_Local(user, req)
	if result != LOCAL_TO_LOCAL_SUCCESS {
		t.Error("CoreClient_to_Local  Test_Local_to_local100 error, result")
	}
	if needupdatRecvQueue != true {
		t.Error("CoreClient_to_Local  Test_Local_to_local100 error, needupdatRecvQueue ")
	}
}

func Test_Local_to_local2002(t *testing.T) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	user := &UserInfor{} //userid = 1
	user.SendId = 300

	user.UserMap = make(map[uint64]*User)
	user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
	user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
	user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

	onefriend := &User{}
	onefriend.UserId = 400
	onefriend.SendId = 30
	onefriend.ReceiveId = 60

	user.UserMap[400] = onefriend

	var id uint64
	var array [100]*GeneralMessage
	for id = 0; id < 100; id++ {
		req := &GeneralMessage{}
		req.SendId = 101 + id
		req.ReceiverId = 1
		req.SenderId = 400

		chat := &ChatMessage{}
		chat.ReceiverId = 1
		chat.SenderId = 400
		chat.SendId = 61 + id

		req.Chatmessage = chat
		array[id] = req
	}

	for id = 0; id < 100; id++ {
		newid := r1.Intn(99)
		req1 := array[id]
		array[id] = array[newid]
		array[newid] = req1
	}

	var newid uint64
	newid = 101

	for id = 0; id < 100; id++ {
		result, needupdate := CoreLocal_to_Local(user, array[id])
		if result == LOCAL_TO_LOCAL_SUCCESS {

			//fmt.Println("now id = ", array[id].SendId)
			if needupdate != true {
				t.Error("Test_Local_to_local2002    error, needupdate ")
			}

			if array[id].SendId != newid {
				t.Error("Test_Local_to_local2002   error, array[id].SendId != newid")
			}
			newid += 1

		}

	}

	for id = 0; id < 100; id++ {
		result, req := SyncLocal_to_Local(user, 400)
		if result == LOCAL_TO_LOCAL_SYNC_SUCCESS_NOSEND {
			break
		} else if result == LOCAL_TO_LOCAL_SYNC_SUCCESS {
			if req.SendId != newid {
				t.Error("Test_Local_to_local2002   error, req.SendId != newid ")
			}
			// fmt.Println( req.SendId, req.Chatmessage.SendId )
			newid += 1
		}

	}
	if newid != 201 {
		t.Error("Test_Local_to_local2002   error, newid != 201 ")
	}

	if len(user.RecvedQueue.MessageMap[400].MessageMap) != 0 {
		t.Error("Test_Local_to_local2002   error, user.RecvedQueue.MessageMap[400].MessageMap)  != 0")
	}

	if user.UserMap[400].ReceiveId != 160 {
		t.Error("Test_Local_to_local2002   error, user.Sendid != 1230")
	}
}
