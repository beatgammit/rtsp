package rtp

import (
	"testing"
)

func TestToUint(t *testing.T) {
	tests := []struct {
		arr []byte
		exp uint
	}{
		{[]byte{1, 2}, 0x102},
		{[]byte{3, 2, 1, 0}, 0x3020100},
	}
	for _, tst := range tests {
		val := toUint(tst.arr)
		if val != tst.exp {
			t.Errorf("%d != %d for % x", val, tst.exp, tst.arr)
		}
	}
}
