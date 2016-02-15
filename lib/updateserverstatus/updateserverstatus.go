package updateserverstatus

import (
	//	"fmt"
	. "github.com/fakewechat/lib/clientpool"
	"github.com/fakewechat/lib/serverstatus"
	. "github.com/fakewechat/lib/version"
	"time"
)

func CheckServerStatusUpdate(index int, servertype string, rg int) {
	var LocalVersion Version
	for {
		time.Sleep(10 * time.Second)
		if LocalVersion != serverstatus.GlobalVersion {
			LocalVersion = serverstatus.GlobalVersion
			copy := serverstatus.CopyLocalServerStatus()
			//fmt.Println("get new copy", copy)
			GloabalClientPool[index].AddNew(LocalVersion, copy, servertype, rg)

		} else {
			//			fmt.Println("equal")
		}
	}
}
