package main

import (
	. "github.com/fakewechat/lib/monitor"
	"time"
)

func main() {

	GMonitor.Init()
	GMonitor.InitServer("demo", "demo_001")
	go GMonitor.Report("127.0.0.1:1234")

	for i := 0; ; i++ {
		time.Sleep(3 * time.Second)
	}
}
