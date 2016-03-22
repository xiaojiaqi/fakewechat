package serverstatus

import (
	"encoding/json"
	"fmt"
	. "github.com/fakewechat/lib/config"
	. "github.com/fakewechat/lib/version"
	"strconv"

	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type ServerStatus struct {
	Name       string `json:"name"`
	ServerType string `json:"srvtype"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Rg         int    `json:"rg"`
	CellId     int    `json:"cellid"`
}

var GlobalVersion Version
var LocalServerStatus map[string]ServerStatus
var Mutex sync.Mutex

func CheckVersionChanged(oldversion Version) bool {

	if oldversion != GlobalVersion {
		return false
	}
	return true
}

func CopyLocalServerStatus() map[string]ServerStatus {
	newMap := make(map[string]ServerStatus)
	Mutex.Lock()
	for k, v := range LocalServerStatus {
		newMap[k] = v
	}
	Mutex.Unlock()

	return newMap
}

func SyncServerStatusChanged(url string, servertype []string, rg []int) {

	for {
		time.Sleep(MonitServerUpdateInterval * time.Second)
		newmap := getServerStatusFromUrl(url)

		destmap := Filter(newmap, servertype, rg)

		bresult := false
		if len(destmap) == 0 {
			continue
		}
		Mutex.Lock()
		bresult = checkServerStatusEqual(LocalServerStatus, destmap)
		Mutex.Unlock()
		if !bresult {
			fmt.Println("not equal")
			copytoLocalServerStatus(destmap)
		} else {

		}

	}
}

func copytoLocalServerStatus(newMap map[string]ServerStatus) {
	Mutex.Lock()

	LocalServerStatus = make(map[string]ServerStatus)
	for k, v := range newMap {
		LocalServerStatus[k] = v
	}

	GlobalVersion += 1

	Mutex.Unlock()
	fmt.Println("CopytoServerList: ", GlobalVersion, " ", LocalServerStatus)
}

func RegistServer(url string, serverType string, hostid string, rg int, host string, port uint, cellid int) {
	//use vector join refactor me
	url = url + serverType + "/" + strconv.Itoa(rg) + "/" + "?"
	url = url + "id=" + hostid + "&" + "host=" + host + "&" + "port=" + strconv.Itoa(int(port))
	url = url + "&cellid=" + strconv.Itoa(int(cellid))
	for {
		fmt.Println(url)
		_, err := http.Get(url)
		if err == nil {
			return
		}
		time.Sleep(MonitServerUpdateInterval * time.Second)
	}

}

func getServerStatusFromUrl(url string) map[string]ServerStatus {
	result := make(map[string]ServerStatus)
	//fmt.Println("getServerStatusFromUrl ", url)

	transport := http.Transport{
		DisableKeepAlives: true,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(url)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result
	}
	vServer := make([]ServerStatus, 0)
	json.Unmarshal([]byte(body), &vServer)
	for _, ser := range vServer {
		key := ser.ServerType + ser.Name
		result[key] = ser
	}
	return result
}

func Filter(org map[string]ServerStatus, servertype []string, rg []int) map[string]ServerStatus {

	dest := make(map[string]ServerStatus)

	len := len(servertype)
	for id := 0; id < len; id++ {
		for key, ser := range org {
			mservertype := servertype[id]
			mrg := rg[id]
			if ser.ServerType == mservertype {
				if mrg == -1 || ser.Rg == mrg {
					dest[key] = ser
				}
			}
		}
	}
	return dest

}

func checkServerStatusEqual(oldMap map[string]ServerStatus, newMap map[string]ServerStatus) bool {
	if len(newMap) != len(oldMap) {
		fmt.Println("len(newMap) != len(oldMap)", len(newMap), len(oldMap))
		return false
	}

	for k, _ := range newMap {
		_, ok := oldMap[k]
		if !ok {
			fmt.Println(k, ", not exit")
			fmt.Println("new", newMap)
			fmt.Println("old", oldMap)

			return false

		}
	}
	return true
}
