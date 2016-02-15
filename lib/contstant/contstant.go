package constant

const (
	CHAT_CLIENT_TO_LOCALPOST    = 1000
	CHAT_LOCALPOST_TO_LOCALPOST = 2000
	CHAT_LOCALPOST_ACK          = 3000
)

const REDISSERVER string = "redis"
const MEMCACHEDSERVER string = "memcached"
const GWSEVER string = "gw"
const LOCALPOSTSERVER string = "localposter"
const POSTERSERVER string = "poster"
const CACHESERVER string = "cache"

func IsServer(s string) bool {
	// use a map instead? no, multithread
	if (s == REDISSERVER) || (s == MEMCACHEDSERVER) || (s == GWSEVER) || (s == LOCALPOSTSERVER) || (s == POSTERSERVER) || (s == CACHESERVER) {
		return true
	}
	return false
}

func IsAppServer(s string) bool {

	if (s == GWSEVER) || (s == LOCALPOSTSERVER) || (s == POSTERSERVER) || (s == CACHESERVER) {
		return true
	}
	return false
}

func IsRedis(s string) bool {

	if s == REDISSERVER {
		return true
	}
	return false
}
func IsMem(s string) bool {

	if s == MEMCACHEDSERVER {
		return true
	}
	return false

}
