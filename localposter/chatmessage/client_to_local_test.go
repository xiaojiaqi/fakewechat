package chatmessage

import (
	. "github.com/fakewechat/message"
	"testing"

	//"fmt"
	"math/rand"
	"time"
)

//
// normal test
//
func Test_send100(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 100

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 30
	
	user.SendedQueue = &SendQueue{}
	user.SendedQueue.MessageMap = make(map[uint64]*GeneralMessage)



	user.UserMap[1] = onefriend

	req := &GeneralMessage{}
	req.SendId = 101
	req.ReceiverId = 1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result, needupdate := CoreClient_to_Local(user, req)
	if result != CLIENT_TO_LOCAL_SUCCESS {
		t.Error("CoreClient_to_Local  send 100 error, result")
	}
	if needupdate != true {
		t.Error("CoreClient_to_Local  send 100 error, needupdate ")
	}

}

//
// don't send data,
//
func Test_send100_notsend(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 110
	user.SendAckId = 109

	user.UserMap = make(map[uint64]*User)
	
	
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 30

	user.UserMap[1] = onefriend

    user.SendedQueue = &SendQueue{}
	user.SendedQueue.MessageMap = make(map[uint64]*GeneralMessage)


	req := &GeneralMessage{}
	req.SendId = 109
	req.ReceiverId = 1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

    user.SendedQueue.MessageMap[109] = req

	result, needupdate := CoreClient_to_Local(user, req)
	if result != CLIENT_TO_LOCAL_SUCCESS_NOSEND {
		t.Error("CoreClient_to_Local  send 100 notsend error, result")
	}
	if needupdate != false {
		t.Error("CoreClient_to_Local  send 100 notsend error, needupdate ")
	}

}

//
//  send again
//
func Test_send100_sendagain(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 110
	user.SendAckId = 66

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 30

	user.UserMap[1] = onefriend
	
	user.SendedQueue = &SendQueue{}
	user.SendedQueue.MessageMap = make(map[uint64]*GeneralMessage)

    
	req := &GeneralMessage{}
	req.SendId = 101
	req.ReceiverId = 1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat
    user.SendedQueue.MessageMap[101] = req
	result, needupdate := CoreClient_to_Local(user, req)
	if result != CLIENT_TO_LOCAL_SUCCESS {
		t.Error("CoreClient_to_Local  send 100 sendagain error, result")
	}
	if needupdate != false {
		t.Error("CoreClient_to_Local  send 100 sendagain error, needupdate ")
	}

}

//
//  send again
//
func Test_send100_storeRequest(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 110
	user.SendAckId = 66

	user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
	user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
	user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 30

	user.UserMap[1] = onefriend

	req := &GeneralMessage{}
	req.SendId = 121
	req.ReceiverId = 1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	result, needupdate := CoreClient_to_Local(user, req)
	if result != CLIENT_TO_LOCAL_SUCCESS_NOSEND {
		t.Error("CoreClient_to_Local  send 100 storeRequest error, result")
	}
	if needupdate != false {
		t.Error("CoreClient_to_Local  send 100 storeRequest error, needupdate ")
	}
	_, ok := user.SendedQueue.MessageMap[121]
	if !ok {
		t.Error("CoreClient_to_Local  send 100 storeRequest error, no request restoraged ")
	}
}

func Test_sync_nosend(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 110
	user.SendAckId = 66

	user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
	user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
	user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 30

	user.UserMap[1] = onefriend

	req := &GeneralMessage{}
	req.SendId = 121
	req.ReceiverId = 1

	chat := &ChatMessage{}
	chat.ReceiverId = 1

	req.Chatmessage = chat

	user.SendedQueue.MessageMap[121] = req

	result, sendreq := SyncClient_to_Local(user)
	if result != CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND {
		t.Error("CoreClient_to_Local  Test_sync_nosend error, result")
	}
	if sendreq != nil {
		t.Error("CoreClient_to_Local  Test_sync_nosend error, sendreq ")
	}
	_, ok := user.SendedQueue.MessageMap[121]
	if !ok {
		t.Error("CoreClient_to_Local Test_sync_nosend error, no request restoraged ")
	}
}

func Test_sync_end3times(t *testing.T) {

	user := &UserInfor{}
	user.SendId = 111
	user.SendAckId = 66

	user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
	user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
	user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

	user.UserMap = make(map[uint64]*User)
	{
		onefriend := &User{}
		onefriend.UserId = 1
		onefriend.SendId = 30

		user.UserMap[1] = onefriend
	}

	{
		onefriend := &User{}
		onefriend.UserId = 3
		onefriend.SendId = 10

		user.UserMap[3] = onefriend
	}

	//  112 for user 3
	{
		req := &GeneralMessage{}
		req.SendId = 112
		req.ReceiverId = 3

		chat := &ChatMessage{}
		chat.ReceiverId = 3

		req.Chatmessage = chat

		user.SendedQueue.MessageMap[112] = req

	}

	//113 for user 3
	{
		req := &GeneralMessage{}
		req.SendId = 113
		req.ReceiverId = 3

		chat := &ChatMessage{}
		chat.ReceiverId = 3

		req.Chatmessage = chat

		user.SendedQueue.MessageMap[113] = req

	}

	//114 for user 1
	{
		req := &GeneralMessage{}
		req.SendId = 114
		req.ReceiverId = 1

		chat := &ChatMessage{}
		chat.ReceiverId = 1

		req.Chatmessage = chat

		user.SendedQueue.MessageMap[114] = req

	}
	//116 for user 3
	{
		req := &GeneralMessage{}
		req.SendId = 116
		req.ReceiverId = 3

		chat := &ChatMessage{}
		chat.ReceiverId = 3

		req.Chatmessage = chat

		user.SendedQueue.MessageMap[116] = req

	}

	// 112
	{
		result, sendreq := SyncClient_to_Local(user)
		if result != CLIENT_TO_LOCAL_SYNC_SUCCESS {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 1 error, result")
		}
		if sendreq == nil {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 1 error, sendreq ")
		}
		if sendreq.SendId != 112 || sendreq.Chatmessage.SendId != 11 {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 1 error, sendreq.SendId != 112 || sendreq.Chatmessage.SendId != 11 ")

		}

		{
			u, ok := user.UserMap[3]
			if ok {
				if u.SendId != 11 {
					t.Error("CoreClient_to_Local  Test_sync_end3times case 1 error, u.SendId != 11 ")

				}
			} else {
				t.Error("CoreClient_to_Local  Test_sync_end3times case 1 error,  no such user 3?? ")

			}

		}
    /* leak this bug
		_, ok := user.SendedQueue.MessageMap[112]
		if ok {
			t.Error("CoreClient_to_Local Test_sync_end3times case 2 error, still in queue")
		}
	*/

	}

	// 113
	{
		result, sendreq := SyncClient_to_Local(user)
		if result != CLIENT_TO_LOCAL_SYNC_SUCCESS {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 2 error, result")
		}
		if sendreq == nil {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 2 error, sendreq ")
		}
		if sendreq.SendId != 113 || sendreq.Chatmessage.SendId != 12 {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 2 error, sendreq.SendId != 113 || sendreq.Chatmessage.SendId != 31 ")

		}

		{
			u, ok := user.UserMap[3]
			if ok {
				if u.SendId != 12 {
					t.Error("CoreClient_to_Local  Test_sync_end3times case 2 error, u.SendId != 12 ")

				}
			} else {
				t.Error("CoreClient_to_Local  Test_sync_end3times case 2 error,  no such user 3?? ")
			}
		}
		/*
		_, ok := user.SendedQueue.MessageMap[113]
		if ok {
			t.Error("CoreClient_to_Local Test_sync_end3times case 2 error, still in queue")
		}
		*/
	}

	// 114
	{
		result, sendreq := SyncClient_to_Local(user)
		if result != CLIENT_TO_LOCAL_SYNC_SUCCESS {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 3 error, result")
		}
		if sendreq == nil {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 3 error, sendreq ")
		}
		if sendreq.SendId != 114 || sendreq.Chatmessage.SendId != 31 {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 3 error, sendreq.SendId != 114 || sendreq.Chatmessage.SendId !=31 ")

		}

		u, ok := user.UserMap[1]
		if ok {
			if u.SendId != 31 {
				t.Error("CoreClient_to_Local  Test_sync_end3times case 3 error, u.SendId != 31 ")

			}
		} else {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 3 error,  no such user 3?? ")

		}

	}
    /*
	_, ok := user.SendedQueue.MessageMap[114]
	if ok {
		t.Error("CoreClient_to_Local Test_sync_end3times case 3 error, still in queue")
	}
	*/
	// 116
	{
		result, sendreq := SyncClient_to_Local(user)
		if result != CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 4 error, result")
		}
		if sendreq != nil {
			t.Error("CoreClient_to_Local  Test_sync_end3times case 4 error, sendreq ")
		}

		{
			{
				u, ok := user.UserMap[3]
				if ok {
					if u.SendId != 12 {
						t.Error("CoreClient_to_Local  Test_sync_end3times case 4 error, u.SendId != 11 ")

					}
				} else {
					t.Error("CoreClient_to_Local  Test_sync_end3times case 4 error,  no such user 3?? ")

				}
			}

			{
				u, ok := user.UserMap[1]
				if ok {
					if u.SendId != 31 {
						t.Error("CoreClient_to_Local  Test_sync_end3times case 4 error, u.SendId != 31 ")

					}
				} else {
					t.Error("CoreClient_to_Local  Test_sync_end3times case 4 error,  no such user 1?? ")

				}
			}

		}

		/*
		_, ok := user.SendedQueue.MessageMap[116]
		if !ok {
			t.Error("CoreClient_to_Local Test_sync_end3times case 4 error, still in queue")
		}
		*/
	}
}

//
// normal random test
//
func Test_randsend100(t *testing.T) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	user := &UserInfor{}
	user.SendId = 100

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 1130

	user.UserMap[1] = onefriend

	user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
	user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
	user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

	var array [100]*GeneralMessage
	var id uint64
	for id = 0; id < 100; id++ {
		req := &GeneralMessage{}
		req.SendId = 101 + id
		req.ReceiverId = 1

		chat := &ChatMessage{}
		chat.ReceiverId = 1

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
		result, needupdate := CoreClient_to_Local(user, array[id])
		if result == CLIENT_TO_LOCAL_SUCCESS {
			if needupdate != true {
				t.Error("CoreClient_to_Local  randsend100 error, needupdate ")
			}

			if array[id].SendId != newid {
				t.Error("CoreClient_to_Local randsend100 error, array[id].SendId != newid")
			}
			newid += 1

		}

	}
	//fmt.Println("got ", newid, " alreay")
	for id = 0; id < 100; id++ {
		result, req := SyncClient_to_Local(user)
		if result == CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND {
			break

		} else if result == CLIENT_TO_LOCAL_SYNC_SUCCESS {
			if req.SendId != newid {

				t.Error("CoreClient_to_Local randsend100 error, req.SendId != newid ")
			}
			newid += 1
		}

	}
	if newid != 201 {
		t.Error("CoreClient_to_Local randsend100 error, newid != 201 ")
	}

    /*
	if len(user.SendedQueue.MessageMap) != 0 {
		t.Error("CoreClient_to_Local randsend100 error, user.SendedQueue)  != 0")
	}
	*/
	if user.SendId != 200 {
		t.Error("CoreClient_to_Local randsend100 error, user.Sendid != 200")
	}

	if user.UserMap[1].SendId != 1230 {
		t.Error("CoreClient_to_Local randsend100 error, user.Sendid != 1230")
	}

}

//
// normal random test 2
//
func Test_randsend100_2(t *testing.T) {
	//s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)

	user := &UserInfor{}
	user.SendId = 100

	user.UserMap = make(map[uint64]*User)
	onefriend := &User{}
	onefriend.UserId = 1
	onefriend.SendId = 30

	user.UserMap[1] = onefriend

	user.SendedQueue = &SendQueue{MessageMap: make(map[uint64]*GeneralMessage)} //    `protobuf:"bytes,6,opt" json:"SendedQueue,omitempty"`
	user.RecvedQueue = &RecvQueue{MessageMap: make(map[uint64]*SendQueue)}      //  `protobuf:"bytes,7,opt" json:"RecvedQueue,omitempty"`
	user.AckedQueue = &AckQueue{MessageMap: make(map[uint64]*GeneralMessage)}   // `protobuf:"bytes,8,opt" json:"AckedQueue,omitempty"`

	var array [100]*GeneralMessage
	var id uint64
	for id = 0; id < 100; id++ {
		req := &GeneralMessage{}
		req.SendId = 101 + id
		req.ReceiverId = 1

		chat := &ChatMessage{}
		chat.ReceiverId = 1

		req.Chatmessage = chat
		array[id] = req
	}

	var newid, clientsendid uint64
	newid = 101
	clientsendid = 31
	for id = 0; id < 100; id++ {
		result, needupdate := CoreClient_to_Local(user, array[id])
		if result == CLIENT_TO_LOCAL_SUCCESS {
			if needupdate != true {
				t.Error("CoreClient_to_Local  randsend100_2 error, needupdate ")
			}

			if array[id].SendId != newid {
				t.Error("CoreClient_to_Local randsend100_2 error, array[id].SendId != newid")
			}
			if array[id].Chatmessage.SendId != clientsendid {
				t.Error("CoreClient_to_Local randsend100_2 error, array[id].SendId != newid")
			}
			newid += 1
			clientsendid += 1

		}

	}
	//fmt.Println("got ", newid, " alreay")
	for id = 0; id < 100; id++ {
		result, req := SyncClient_to_Local(user)
		if result == CLIENT_TO_LOCAL_SYNC_SUCCESS_NOSEND {
			break

		} else if result == CLIENT_TO_LOCAL_SYNC_SUCCESS {
			if req.SendId != newid {

				t.Error("CoreClient_to_Local randsend100_2 error, req.SendId != newid ")
			}
			newid += 1
			if array[id].Chatmessage.SendId != clientsendid {
				t.Error("CoreClient_to_Local randsend100_2 error, array[id].SendId != newid")
			}

			clientsendid += 1
		}

	}
	if newid != 201 {
		t.Error("CoreClient_to_Local randsend100_2 error, newid != 201 ")
	}
	/*
	if len(user.SendedQueue.MessageMap) != 0 {
		t.Error("CoreClient_to_Local randsend100_2 error, user.SendedQueue)  != 0")
	}
	*/
	if user.SendId != 200 {
		t.Error("CoreClient_to_Local randsend100_2 error, user.Sendid != 200")
	}

	if user.UserMap[1].SendId != 130 {
		t.Error("CoreClient_to_Local randsend100_2 error, user.Sendid != 201")
	}

}
