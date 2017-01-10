package lightswarm

import (
	"encoding/binary"
	"io"
)

// Byte constants
const (
	END byte = 192 // End byte
	ESC byte = 219 // Escape byte
)

// Command constants
const (
	ON                         byte = 32 // on
	OFF                        byte = 33 // off
	SET_LEVEL                  byte = 34 // set level
	FADE_TO_LEVEL              byte = 35 // fade to level
	SET_PSUEDO_ADDRESS         byte = 37 // set psuedo address
	ERASE_PSUEDO_ADDRESS_TABLE byte = 38 // erase psuedo address table
	SET_RGB_LEVELS             byte = 44 // set rgb levels
	TOGGLE                     byte = 45 // toggle
	FADE_MULTIPLE_TO_LEVEL     byte = 48 // fade multiple to level
	FADE_RGB_TO_LEVEL          byte = 49 // fade rgb to level
)

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
			nbs = append(nbs, ESC, 220)
		case ESC:
			nbs = append(nbs, ESC, 221)
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

// Constructs a new LED
func New(addr uint16, writer io.Writer) *LED {
	return &LED{
		Addr:   addr,
		Writer: writer,
	}
}
