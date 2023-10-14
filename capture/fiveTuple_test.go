package capture

import (
	"fmt"
	"testing"
)

func TestFiveTuple_Hash(t *testing.T) {
	fiveTuple := FiveTuple{
		SrcIP:        [4]byte{111,22,11,22},
		DstIP:        [4]byte{222,12,11,22},
		SrcPort:      22,
		DstPort:      888,
		ProtocolType: 1,
	}

	fmt.Println(fiveTuple.Hash())
}
