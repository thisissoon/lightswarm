package lightswarm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrameAddress(t *testing.T) {
	tt := []struct {
		address uint16
		b1      byte
		b2      byte
	}{
		{
			1,
			0,
			1,
		},
		{
			690,
			2,
			178,
		},
		{
			65535,
			255,
			255,
		},
	}
	for i, tc := range tt {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			f := Frame{Addr: tc.address}
			b1, b2 := f.address()
			assert.Equal(t, tc.b1, b1)
			assert.Equal(t, tc.b2, b2)
		})
	}
}

func TestFrameChecksum(t *testing.T) {
	tt := []struct {
		bs       []byte
		expected byte
	}{
		{
			[]byte{0, 1, 2, 3, 4},
			4,
		},
		{
			[]byte{2, 178, ON},
			144,
		},
	}
	for i, tc := range tt {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			f := Frame{}
			checksum := f.checksum(tc.bs)
			assert.Equal(t, tc.expected, checksum)
		})
	}
}
