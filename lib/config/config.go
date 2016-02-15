package config

const (
	GateWayChannelSize    = 4
	GateWayWorkThreadSize = GateWayChannelSize

	GateWayChannelLength = 1000000

	LocalPostChannelSize    = 4
	LocalPostWorkThreadSize = LocalPostChannelSize

	LocalPostChannelLength     = 1000000
	LocalPostThreadQueueLength = 10 // each time, one thread get request from queue

	PostChannelSize    = 4
	PostWorkThreadSize = PostChannelSize

	PostChannelLength     = 1000000
	PostThreadQueueLength = 10

	CacheServerChannelSize = 4
	CacheThreadSize        = CacheServerChannelSize

	ClientPoolSize = 64 // must bigger than each channelsize or each  workthreadsize

	MonitServerUpdateInterval = 10 // second
)
