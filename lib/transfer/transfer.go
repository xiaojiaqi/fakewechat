package transfer

import (
	"fmt"
	"github.com/smartproxy/lib/buff"
	"github.com/smartproxy/lib/conn"
)

func Transfer_Data(ClientConn conn.ConnInterface, ServerConn conn.ConnInterface) {
	// 从 server 获取结果
	defer ClientConn.Close()
	defer fmt.Printf("close ClientConn connection conn %T\n", ClientConn)

	defer ServerConn.Close()
	defer fmt.Printf("close ServerConn connection conn %T\n", ServerConn)

	var ClientBuff buff.Buffer
	ClientBuff.Init(65 * 1024)

	for {
		if !ClientConn.Read(&ClientBuff) {
			return
		}

		if !ServerConn.Write(&ClientBuff, ClientBuff.BuffDataLen) {
			return
		}

	}
	// 加密后发给客户端
}
