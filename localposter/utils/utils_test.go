package utils

import (
	"testing"
)

func Test_utls_001(t *testing.T) {
	i := 0
	i = getChannelId(0, 2, 3, 0)
	if i != 0 {
		t.Error("getChannelId(2, 0, 3, 0) !=0", i)
	}
	i = getChannelId(1, 2, 3, 0)
	if i != 3 {
		t.Error("getChannelId(2, 1, 3,0) != 3", i)
	}
	i = getChannelId(2, 2, 3, 0)
	if i != 0 {
		t.Error("getChannelId(1, 2, 3, 0) != 0", i)
	}
	i = getChannelId(3, 2, 3, 0)
	if i != 3 {
		t.Error("getChannelId(1, 2, 3, 0) != 3", i)
	}

	i = getChannelId(0, 2, 3, 1)
	if i != 1 {
		t.Error("getChannelId(2, 0, 3, 1) !=1", i)
	}
	i = getChannelId(1, 2, 3, 1)
	if i != 4 {
		t.Error("getChannelId(2, 1, 3,1) != 3", i)
	}
	i = getChannelId(2, 2, 3, 1)
	if i != 1 {
		t.Error("getChannelId(1, 2, 3, 1) != 0", i)
	}
	i = getChannelId(3, 2, 3, 1)
	if i != 4 {
		t.Error("getChannelId(1, 2, 3, 1) != 3", i)
	}
}
