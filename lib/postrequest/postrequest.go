package postrequest

import (
	"errors"
	
	. "github.com/fakewechat/message"
)

func RecvfromQueue(queue *chan *GeneralMessage) (*GeneralMessage, error) {

	
	select {
	case v := <-*queue:
		return v, nil
	default:
		//fmt.Println("it is empty")
		return nil, errors.New("it is full")
	}
}

func RecvfromQueueBlock(queue *chan *GeneralMessage) (*GeneralMessage, error) {

	
	v := <-*queue
	return v, nil
}

func GetSomeMessageFromQueue(queue *chan *GeneralMessage, least int) []*GeneralMessage {

	postList := make([]*GeneralMessage, 0)

	v, _ := RecvfromQueueBlock(queue)
	postList = append(postList, v)
	for i := 1; i < least; i++ {

		v, err := RecvfromQueue(queue)
		if err != nil {
			return postList
		}
		postList = append(postList, v)
	}
	return postList
}

func PushtoQueue(queue *chan *GeneralMessage, req *GeneralMessage) error {
	
	select {
	case *queue <- req:
		
		return nil
	default:
		
		return errors.New("it is full")
	}
}
