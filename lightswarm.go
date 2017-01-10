package lightswarm

import (
	"encoding/binary"
	"io"
)

// Byte constants
const (
	END byte = 0xC0 // End byte
	ESC byte = 0xDB // Escape byte
)

// Command constants
const (
	ON                         byte = 0x20 // on
	OFF                        byte = 0x21 // off
	SET_LEVEL                  byte = 0x22 // set level
	FADE_TO_LEVEL              byte = 0x23 // fade to level
	SET_PSUEDO_ADDRESS         byte = 0x25 // set psuedo address
	ERASE_PSUEDO_ADDRESS_TABLE byte = 0x26 // erase psuedo address table
	SET_RGB_LEVELS             byte = 0x2C // set rgb levels
	TOGGLE                     byte = 0x2D // toggle
	FADE_MULTIPLE_TO_LEVEL     byte = 0x30 // fade multiple to level
	FADE_RGB_TO_LEVEL          byte = 0x31 // fade rgb to level
)

type Fade struct {
	Level    byte
	Interval byte
	Step     byte
}

func (f Fade) Args() []byte {
	return []byte{f.Level, f.Interval, f.Step}
}

type Frame struct {
	// Exported Fields
	Addr    uint16
	Cmd     byte
	CmdArgs []byte
}

// Returns the frame address broken into 2 bytes
func (f Frame) address() (b1, b2 byte) {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, f.Addr)
	return bs[0], bs[1]
}

// Calculates checksum for given bytes
func (f Frame) checksum(bs []byte) byte {
	var checksum byte
	for _, b := range bs {
		checksum ^= b
	}
	return checksum
}

// Wrap the given bytes in end bytes with escape sequences
func (f Frame) wrap(bs []byte) []byte {
	// create a new byte slice to store the wrapped frame, starting with the END byte
	nbs := []byte{END}
	// loop over the given byte slice, appending to the new byte slice
	// and performing any escapes required
	for _, b := range bs {
		switch b {
		case END:
			nbs = append(nbs, ESC, 0xDC)
		case ESC:
			nbs = append(nbs, ESC, 0xDD)
		default:
			nbs = append(nbs, b)
		}
	}
	// now add the end byte last
	nbs = append(nbs, END)
	return nbs
}

// Returns the frame in byte format for writing to lightswarm
func (f Frame) Bytes() []byte {
	bs := []byte{}
	// Add address bytes
	addr1, addr2 := f.address()
	bs = append(bs, addr1)
	bs = append(bs, addr2)
	// Add command byte
	bs = append(bs, f.Cmd)
	// Add command arg bytes
	bs = append(bs, f.CmdArgs...)
	// Add Checksum
	bs = append(bs, f.checksum(bs))
	// Wrap bytes in end bytes with escape sequences
	bs = f.wrap(bs)
	return bs
}

// Represents a single Lightswarm LED
type LED struct {
	// Exported Fields
	Addr   uint16
	Writer io.Writer
}

// Send the On command to the LED writer
func (led *LED) On() (int, error) {
	frame := Frame{Addr: led.Addr, Cmd: ON}
	return led.Writer.Write(frame.Bytes())
}

// Send the Off command to the LED writer
func (led *LED) Off() (int, error) {
	frame := Frame{Addr: led.Addr, Cmd: OFF}
	return led.Writer.Write(frame.Bytes())
}

// Fade to a light level
func (led *LED) Fade(f Fade) (int, error) {
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     FADE_TO_LEVEL,
		CmdArgs: f.Args(),
	}
	return led.Writer.Write(frame.Bytes())
}

// Set Red, Green and Blue levels
func (led *LED) RGB(r, g, b byte) (int, error) {
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     SET_RGB_LEVELS,
		CmdArgs: []byte{r, g, b},
	}
	return led.Writer.Write(frame.Bytes())
}

// Fade to a RGB level
func (led *LED) FadeRGB(r, g, b Fade) (int, error) {
	bs := []byte{}
	bs = append(bs, r.Args()...)
	bs = append(bs, g.Args()...)
	bs = append(bs, b.Args()...)
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     FADE_RGB_TO_LEVEL,
		CmdArgs: bs,
	}
	return led.Writer.Write(frame.Bytes())
}

// Constructs a new LED
func New(addr uint16, writer io.Writer) *LED {
	return &LED{
		Addr:   addr,
		Writer: writer,
	}
}
