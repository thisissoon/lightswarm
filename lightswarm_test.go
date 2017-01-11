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

func TestFrameWrap(t *testing.T) {
	tt := []struct {
		name     string
		bs       []byte
		expected []byte
	}{
		{
			"no bytes",
			[]byte{},
			[]byte{END, END},
		},
		{
			"end byte in bytes",
			[]byte{END},
			[]byte{END, ESC, 0xDC, END},
		},
		{
			"esc byte in bytes",
			[]byte{ESC},
			[]byte{END, ESC, 0xDD, END},
		},
		{
			"turn 690 on",
			[]byte{2, 178, ON, 144},
			[]byte{END, 2, 178, ON, 144, END},
		},
		{
			"turn 738 on",
			[]byte{2, 226, ON, 192},
			[]byte{END, 2, 226, ON, ESC, 0xDC, END},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			f := Frame{}
			frame := f.wrap(tc.bs)
			assert.Equal(t, tc.expected, frame)
		})
	}
}
