package clientpool

import (
	"fmt"
	. "github.com/fakewechat/lib/config"
	. "github.com/fakewechat/lib/contstant"
	. "github.com/fakewechat/lib/serverstatus"
	"net/rpc"
	. "github.com/fakewechat/lib/version"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
)


type ClientInterface interface {
	GetVersion() Version
	GetServerType() string
	GetRG() int
	GetKeyinclientPool() string //  type + rg
	GetKeyinRGPool() string     //  hostid + sid
	GetCellId() int             // getHashIndex
	GetRPCClient() *rpc.Client
	GetMemCachedClient() *memcache.Client
	GetRedisClient() *redis.Conn
}

type ClientPool struct {
	version        Version
	lock           sync.Mutex
	sid            int
	pools          map[string]*RGPool // key mean GW+rg, Post+RG,
	memcachedpools [1024]ClientInterface
	redispools     [1024]ClientInterface
}

type RGPool struct {
	maps map[string]ClientInterface
}

type baseClient struct {
	version           Version
	serverType        string
	rgid              int    // for localposter poster gateway it mean rgid
	cellid            int    // for redis & memcache only
	key_in_clientPool string // type + rg
	key_in_RGPool     string // hostid+sid
}

type RPCClient struct {
	base      baseClient
	Rawclient *rpc.Client
}

type RedisClient struct {
	base      baseClient
	Rawclient *redis.Conn
}

type MemcachedClient struct {
	base      baseClient
	Rawclient *memcache.Client
}

var GloabalClientPool [ClientPoolSize]*ClientPool

func InitClientPool() {
	for i := 0; i < ClientPoolSize; i++ {
		p := &ClientPool{}
		p.Init()
		GloabalClientPool[i] = p
	}
}

func (v *ClientPool) Init() {
	v.pools = make(map[string]*RGPool)
}

func (v *ClientPool) Show() {
}

func (v *ClientPool) AddNew(newversion Version, serverList map[string]ServerStatus, localtype string, rg int) {

	fmt.Println("AddNew(newversion Version, serverList map[string]ServerStatus, rg, cellid)", " localtype: ", localtype, " rg: ", rg)
	v.lock.Lock()
	defer v.lock.Unlock()

	localServer_key := localtype + "-" + strconv.Itoa(rg)
	for _, newserver := range serverList {

		service := newserver.Host + ":" + strconv.Itoa(newserver.Port)
		servertype := newserver.ServerType
		server_key := servertype + "-" + strconv.Itoa(newserver.Rg)
		fmt.Println("server_key:", server_key, " localServer_key:", localServer_key)
		if localServer_key == server_key {
			continue
		}
		var base baseClient
		v.sid += 1
		base.key_in_clientPool = servertype + strconv.Itoa(newserver.Rg)
		base.version = newversion
		base.serverType = servertype
		base.rgid = newserver.Rg
		base.cellid = newserver.CellId
		base.key_in_RGPool = newserver.Name + strconv.Itoa(v.sid)

		// switch, refactor me
		if IsRedis(newserver.ServerType) {
			fmt.Println("connect redis: ", service, " rg ", base.rgid, " cellid ", base.cellid)
			conn, err := redis.Dial("tcp", service)
			if err != nil {
				fmt.Println(err)
				continue
			}

			newClient := &RedisClient{}
			newClient.base = base
			newClient.Rawclient = &conn

			var point ClientInterface
			point = newClient
			v.redispools[base.rgid] = point

		} else if IsMem(newserver.ServerType) {

			fmt.Println("mem: ", service, " rg ", base.rgid, " cellid ", base.cellid)

			mc := memcache.New(service)
			newClient := &MemcachedClient{}
			mc.Set(&memcache.Item{Key: "hello", Value: []byte("world")})
			newClient.base = base
			newClient.Rawclient = mc
			var point ClientInterface
			point = newClient
			v.memcachedpools[base.rgid] = point
			// hard code , refactor me
		} else if IsAppServer(newserver.ServerType) {

			newClient := &RPCClient{}
			newRPCClient, err := rpc.Dial("tcp", service)
			if err != nil {
				panic(err)
				continue
			}

			fmt.Println("connect app: ", newserver.ServerType, " ", service, " rg ", base.rgid, " cellid ", base.cellid)
			//fmt.Println("Connect to ", service)
			newClient.base = base
			newClient.Rawclient = newRPCClient

			pool, ok := v.pools[base.key_in_clientPool]
			if ok != true {
				newPool := &RGPool{maps: make(map[string]ClientInterface)}

				v.pools[base.key_in_clientPool] = newPool
				pool = newPool
			}

			var point ClientInterface

			point = newClient
			pool.maps[base.key_in_RGPool] = point
		} else {
			//unsupport type
			fmt.Println("%v", newserver)

			panic("unsupport type")
		}

	}
	v.Show()
	v.version = newversion
}

func (v *ClientPool) GetClient(servertype string, rg int) ClientInterface {

	key := servertype + strconv.Itoa(rg)
	v.lock.Lock()
	defer v.lock.Unlock()

	if IsRedis(servertype) {
		p := v.redispools[rg]
		v.redispools[rg] = nil
		//fmt.Println(servertype, p, v)
		return p
	} else if IsMem(servertype) {
		p := v.memcachedpools[rg]
		v.memcachedpools[rg] = nil
		return p
	} else if IsAppServer(servertype) {
		onepool, ok := v.pools[key]
		if ok != true {
			return nil
		}

		for k, rpcclient := range onepool.maps {
			delete(onepool.maps, k)
			v.Show()
			return rpcclient
		}

	} else {

		panic("unsupport type")
	}
	return nil
}

func (v *ClientPool) ReturnClient(p ClientInterface) {

	v.lock.Lock()
	defer v.lock.Unlock()
	if IsRedis(p.GetServerType()) {
		v.redispools[p.GetRG()] = p
	} else if IsMem(p.GetServerType()) {
		v.memcachedpools[p.GetRG()] = p
	} else if IsAppServer(p.GetServerType()) {
		onepool, ok := v.pools[p.GetKeyinclientPool()]
		if ok != true {
			panic("can't found client pool")
			return
		}
		onepool.maps[p.GetKeyinRGPool()] = p
	} else {

		panic("unsupport type")
	}
}

func (v *RPCClient) GetVersion() Version {
	return v.base.version
}
func (v *RPCClient) GetServerType() string {
	return v.base.serverType
}
func (v *RPCClient) GetRG() int {
	return v.base.rgid
}

func (v *RPCClient) GetKeyinclientPool() string { //  type + rg
	return v.base.key_in_clientPool
}

func (v *RPCClient) GetKeyinRGPool() string { //  hostid + sid
	return v.base.key_in_RGPool
}

func (v *RPCClient) GetCellId() int { // getHashIndex

	return -1
}

func (v *RPCClient) GetRPCClient() *rpc.Client {
	return v.Rawclient
}

func (v *RPCClient) GetMemCachedClient() *memcache.Client {
	return nil
}

func (v *RPCClient) GetRedisClient() *redis.Conn {
	return nil
}

func (v *RedisClient) GetVersion() Version {
	return v.base.version
}
func (v *RedisClient) GetServerType() string {
	return v.base.serverType
}
func (v *RedisClient) GetRG() int {
	return v.base.rgid
}

func (v *RedisClient) GetKeyinclientPool() string { //  type + rg
	return v.base.key_in_clientPool
}

func (v *RedisClient) GetKeyinRGPool() string { //  hostid + sid
	return v.base.key_in_RGPool
}

func (v *RedisClient) GetCellId() int { // getHashIndex

	return v.base.cellid
}

func (v *RedisClient) GetRPCClient() *rpc.Client {

	return nil
}

func (v *RedisClient) GetMemCachedClient() *memcache.Client {
	return nil
}

func (v *RedisClient) GetRedisClient() *redis.Conn {
	return v.Rawclient
}

func (v *MemcachedClient) GetVersion() Version {
	return v.base.version
}
func (v *MemcachedClient) GetServerType() string {
	return v.base.serverType
}
func (v *MemcachedClient) GetRG() int {
	return v.base.rgid
}

func (v *MemcachedClient) GetKeyinclientPool() string { //  type + rg
	return v.base.key_in_clientPool
}

func (v *MemcachedClient) GetKeyinRGPool() string { //  hostid + sid
	return v.base.key_in_RGPool
}

func (v *MemcachedClient) GetCellId() int { // getHashIndex

	return v.base.cellid
}

func (v *MemcachedClient) GetRPCClient() *rpc.Client {

	return nil
}

func (v *MemcachedClient) GetMemCachedClient() *memcache.Client {
	return v.Rawclient
}

func (v *MemcachedClient) GetRedisClient() *redis.Conn {

	return nil
}
