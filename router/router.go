package main

import (
	"encoding/json"
	"fmt"
	. "github.com/fakewechat/lib/contstant"
	"github.com/fakewechat/lib/flags"
	. "github.com/fakewechat/lib/serverstatus"
	. "github.com/fakewechat/lib/utils"
	. "github.com/fakewechat/lib/version"
	"log"
	"net/http"
	"regexp"
	//"strconv"
	"strings"
	"sync"
	"time"
)

var lock sync.Mutex

var servermap map[string]ServerStatus

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {

		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func getserver(rw http.ResponseWriter, req *http.Request) {
	servers := make([]ServerStatus, 0)

	lock.Lock()
	defer lock.Unlock()
	for _, v := range servermap {
		servers = append(servers, v)
	}
	lang2, err := json.Marshal(servers)
	if err == nil {
		rw.Write(lang2)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func regist(rw http.ResponseWriter, req *http.Request) {

	lock.Lock()
	defer lock.Unlock()
	url := req.URL

	fmt.Println("url = ", url.Path)
	paths := strings.Split(url.Path, "/")
	fmt.Println("paths ", paths)
	req.ParseForm()
	serv := ServerStatus{}

	serv.Name = GetReqdate(req, "id")
	serv.ServerType = paths[2]
	serv.Rg = GetIntValue(paths[3])
	serv.Host = GetReqdate(req, "host")
	serv.Port = GetReqIntdata(req, "port")
	serv.CellId = GetReqIntdata(req, "cellid")
	if !IsServer(serv.ServerType) {
		panic("unsupport server type :" + serv.ServerType)
	}
	fmt.Println(serv)
	servermap[serv.ServerType+serv.Name] = serv
	fmt.Println("====regist begin ======")
	for k, v := range servermap {
		fmt.Println(k, v)
	}
	fmt.Println("===== end =======")
	rw.Write([]byte("Done"))
}

func redirect(rw http.ResponseWriter, req *http.Request) {

	lock.Lock()
	defer lock.Unlock()
	header := rw.Header()
	header.Add("Location", "http://www.163.com")

	rw.WriteHeader(http.StatusFound)
}

func main() {

	PrintVersion()
	 
	servermap = make(map[string]ServerStatus)

	flags.Parse()

	var regexpHandler RegexpHandler
	xp, err := regexp.Compile("/server")
	if err == nil {
		regexpHandler.HandleFunc(xp, getserver)  
	}
	xp, err = regexp.Compile("/regist/(cache|gw|poster|localposter|redis|memcached)/[0-9]*/.*")
	if err == nil {
		regexpHandler.HandleFunc(xp, regist)  
	}

	xp, err = regexp.Compile("/friend/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, redirect)  
	}

	xp, err = regexp.Compile("/inbox/[0-9]*/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, redirect)  
	}
	xp, err = regexp.Compile("/outbox/[0-9]*/[0-9]*")
	if err == nil {
		regexpHandler.HandleFunc(xp, redirect)  
	}

	server := &http.Server{
		Addr:           ":8080",
		Handler:        regexpHandler,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())

}
