package utils

func GetChannelId(index uint32, times int, channelsize int, offset int) int {
	id := int(index%uint32(times))*channelsize + offset
	return id
}
