package monitor

import (
	"testing"
)

func Test_monitor(t *testing.T) {
	a := &Monitor{}

	a.Init()

	a.Regist("aaa")
	a.Regist("bbb")
	a.Regist("ccc")

	a.Add("aaa", 1)
	a.Add("aaa", 1)
	a.Add("aaa", 1)

	a.Add("bbb", 1)
	a.Add("bbb", 1)
	a.Add("ccc", 1)
	a.Add("ccc", 2)
	a.Add("ccc", -1)

	a.Print()
}
