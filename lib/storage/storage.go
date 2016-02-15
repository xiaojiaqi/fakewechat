package storage

import (
	"errors"
	. "github.com/fakewechat/lib/monitor"
	"github.com/fakewechat/lib/rpc"
	. "github.com/fakewechat/message"
)

type Storage int

func (t *Storage) GetUnackMessageId(args *Args, reply *int) error {
	return nil
}

func (t *Storage) GetUnackMessage(args *Args, quo *Quotient) error {
	return nil
}

func (t *Storage) AckMessageId(args *Args, reply *int) error {
	return nil
}
