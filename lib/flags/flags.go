package flags

import (
	"flag"
	"fmt"
)

//

// gw, post,
var ProcessName *string

// listen address
var ListenAddress *string

// listen port
var ListenPort *uint

// hostid, such like  gw_r01_001
var HostId *string

// gw
var ServerType *string

// 1, 2, 3, 4
var RgId *uint

// rg size ,default  25000
var RgSize *int

// monitor server
var MonitorHost *string

//
var SnId *uint

// cell id
var CellId *int

// default is 45 second
var GatewayMessageTTL *int

// default route url
var RouterServerURL *string

// default route url
var RouterRegistURL *string

func Parse() {
	ProcessName = flag.String("process", "process", "process name")
	ListenAddress = flag.String("listenaddress", "0.0.0.0", "listen address")
	ListenPort = flag.Uint("listenport", 1000, "list port")
	HostId = flag.String("hostid", "host", "host name")
	ServerType = flag.String("servertype", "gw", "server type")
	RgId = flag.Uint("rgid", 1, "RG ID")
	MonitorHost = flag.String("monitorhost", "127.0.0.1:8000", "monitor service")
	SnId = flag.Uint("snid", 1, "snid")

	CellId = flag.Int("cellid", 1, "cellid")

	RouterServerURL = flag.String("routeserverurl", "http://0.0.0.0.:8080/server", "route server url")
	RouterRegistURL = flag.String("routeregisturl", "http://0.0.0.0.:8080/regist/", "route regist url")
	flag.Parse()
	Print()
}

func Print() {

	fmt.Println("ProcessName: ", *ProcessName)
	fmt.Println("ListenAddress: ", *ListenAddress)
	fmt.Println("ListenPort: ", *ListenPort)
	fmt.Println("HostId: ", *HostId)
	fmt.Println("RgId: ", *RgId)
	fmt.Println("MonitorHost: ", *MonitorHost)
	fmt.Println("SnId: ", *SnId)
	fmt.Println("CellId: ", *CellId)
	fmt.Println("Route Url: ", *RouterServerURL)
	fmt.Println("Route Url: ", *RouterRegistURL)

}
