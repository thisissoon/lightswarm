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

// Escape Sequences
var (
	ENDSEQ = []byte{ESC, 0xDC} // END byte escape sequence
	ESCSEQ = []byte{ESC, 0xDD} // ESC byte escape sequence
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

// Helper for easily constructing Fade commands
type Fade struct {
	Level    byte
	Interval byte
	Step     byte
}

// Command arguments
func (f Fade) Args() []byte {
	return []byte{f.Level, f.Interval, f.Step}
}

// Builds a data frame to be sent on the serial connection
// A frame follows this format:
// 1: Start Byte
// 2: Destination Address byte 1
// 3: Destination Address byte 2
// 4: Command byte
// 5: Information bytes (0 or more bytes)
// 6: Checksum byte
// 7: End byte
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
	frame := []byte{END}
	// loop over the given byte slice, appending to the new byte slice
	// and performing any escapes required
	for _, b := range bs {
		switch b {
		case END:
			frame = append(frame, ENDSEQ...)
		case ESC:
			frame = append(frame, ESCSEQ...)
		default:
			frame = append(frame, b)
		}
	}
	// now add the end byte last
	frame = append(frame, END)
	return frame
}

// Returns the frame in byte format for writing to lightswarm
func (f Frame) Bytes() []byte {
	// Create the data frame
	frame := []byte{}
	// Add address bytes
	addr1, addr2 := f.address()
	frame = append(frame, addr1)
	frame = append(frame, addr2)
	// Add command byte
	frame = append(frame, f.Cmd)
	// Add command arg bytes
	frame = append(frame, f.CmdArgs...)
	// Calculate & Add Checksum
	checksum := f.checksum(frame)
	frame = append(frame, checksum)
	// Wrap bytes in end bytes with escape sequences
	frame = f.wrap(frame)
	return frame
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
func (led *LED) SetRGB(r, g, b byte) (int, error) {
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     SET_RGB_LEVELS,
		CmdArgs: []byte{r, g, b},
	}
	return led.Writer.Write(frame.Bytes())
}

// Fade to a RGB level
func (led *LED) FadeRGB(r, g, b Fade) (int, error) {
	args := []byte{}
	args = append(args, r.Args()...)
	args = append(args, g.Args()...)
	args = append(args, b.Args()...)
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     FADE_RGB_TO_LEVEL,
		CmdArgs: args,
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
