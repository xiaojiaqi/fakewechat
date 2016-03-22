package config

const (
	GateWayChannelSize    = 4
	GateWayWorkThreadSize = GateWayChannelSize

	GateWayChannelLength = 100000

	MinLocalPostChannelSize = 3 /* 0 for client to local, 1 for local to local , 2 for local ack*/
	LocalPostChannelSize    = 18

	LocalPostWorkThreadSize = LocalPostChannelSize

	LocalPostChannelLength     = 100000
	LocalPostThreadQueueLength = 10 // each time, one thread get request from queue

	PostChannelSize    = 4
	PostWorkThreadSize = PostChannelSize

	PostChannelLength     = 300000
	PostThreadQueueLength = 10

	CacheServerChannelSize = 4
	CacheThreadSize        = CacheServerChannelSize

	ClientPoolSize = 64 // must bigger than each channelsize or each  workthreadsize

	MonitServerUpdateInterval = 10 // second
)
