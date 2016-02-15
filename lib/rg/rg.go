package rg

const (
	rgsize        = 100
	rediscellsize = 2500
	memcellsize   = 10
)

func GetRg(userid uint64) int {
	//  1-50      1
	//  51-100    2
	//  101-150   3

	// for Irregular we should use bitmap
	id := int(userid) - 1
	rg := id/rgsize + 1
	return rg
}

/*
func GetRedisCellId(userid uint64) int {
	id := int(userid) - 1
	rg := id / rgsize
	id = id - rg * rgsize
	rg = id / rediscellsize + 1

	return rg
}

func GetMemcachedCellId(messageid uint64) int {

	id := int(messageid)
	return id % + memcellsize
}
*/
