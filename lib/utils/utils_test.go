package utils

import (
	"testing"
)

//
// normal test
//
func Test_fun1(t *testing.T) {

	if GetUserQueue(123) != "user#queue#123" {
		t.Error("GetUserQueue(123)  !=  user#queue#123")
	}

	if GetSendMessage(456) != "send_456" {
		t.Error("GetUserQueue(123)  ==  user#queue#123")
	}

	if GetLocalMessage(999, 44) != "local_999_44" {
		t.Error("GetUserQueue(123)  ==  user#queue#123")

	}

	if GetAckMessage(1111111111) != "ack_1111111111" {

		t.Error("GetAckMessage( 1111111111 ) !=  ack_1111111111")
	}

	if GetQueueId(0, 0, 9, 3) != 0 {
		t.Error("GetQueueId(0, 0) != 0 )")
	}

	if GetQueueId(1, 0, 9, 3) != 3 {
		t.Error("GetQueueId(1, 0) != 0 )")
	}

	if GetQueueId(2, 0, 9, 3) != 6 {
		t.Error("GetQueueId(2, 0) !=6 )")
	}

	if GetQueueId(3, 0, 9, 3) != 0 {
		t.Error("GetQueueId(3, 0) != 0 )")
	}

	if GetQueueId(0, 1, 9, 3) != 1 {
		t.Error("GetQueueId(0, 1) != 1 )")
	}

	if GetQueueId(1, 1, 9, 3) != 4 {
		t.Error("GetQueueId(1, 1) != 4 )")
	}

	if GetQueueId(2, 1, 9, 3) != 7 {
		t.Error("GetQueueId(2, 1) != 7 )")
	}

	if GetQueueId(3, 1, 9, 3) != 1 {
		t.Error("GetQueueId(3, 1) != 1 )")
	}
}
