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
	FADE_DOWN                  byte = 0x24 // legacy fade down
	SET_PSUEDO_ADDRESS         byte = 0x25 // set psuedo address
	ERASE_PSUEDO_ADDRESS_TABLE byte = 0x26 // erase psuedo address table
	SET_RGB_LEVELS             byte = 0x2C // set rgb levels
	TOGGLE                     byte = 0x2D // toggle
	FADE_MULTIPLE_TO_LEVEL     byte = 0x30 // fade multiple to level
	FADE_RGB_TO_LEVEL          byte = 0x31 // fade rgb to level
)

// Helper for easily constructing Fade commands
type Fade struct {
	Level    int
	Interval int
	Step     int
}

// Returns the level to fade too, max value is 255
func (f Fade) level() int {
	if f.Level > 255 {
		return 255
	}
	return f.Level
}

// Returns the interval step for fading in 1/100's of a second
// 1 is the lowest allowed value and represents 1 1/100 (10ms)
func (f Fade) interval() int {
	if f.Interval == 0 { // 0 interval values are not supported
		return 1
	}
	return f.Interval
}

// Value to increment the light level by on each interval
// allowed range is 1-127
func (f Fade) step() int {
	if f.Step == 0 { // 0 interval values are not supported
		return 1
	}
	if f.Step > 127 {
		return 127
	}
	return f.Step
}

// Command arguments
func (f Fade) Args() []byte {
	return []byte{
		byte(f.level()),
		byte(f.interval()),
		byte(f.step()),
	}
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

// Write to the lightswarm writer
func (led *LED) write(frame Frame) (int, []byte, error) {
	b := frame.Bytes()
	n, err := led.Writer.Write(b)
	if err != nil {
		return 0, nil, err
	}
	return n, b, nil
}

// Send the On command to the LED writer
func (led *LED) On() (int, []byte, error) {
	frame := Frame{Addr: led.Addr, Cmd: ON}
	return led.write(frame)
}

// Send the Off command to the LED writer
func (led *LED) Off() (int, []byte, error) {
	frame := Frame{Addr: led.Addr, Cmd: OFF}
	return led.write(frame)
}

// Fade down legacy
func (led *LED) FadeDown(f Fade) (int, []byte, error) {
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     FADE_DOWN,
		CmdArgs: f.Args(),
	}
	return led.write(frame)
}

// Fade to a light level
func (led *LED) Fade(f Fade) (int, []byte, error) {
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     FADE_TO_LEVEL,
		CmdArgs: f.Args(),
	}
	return led.write(frame)
}

// Set Red, Green and Blue levels
func (led *LED) SetRGB(r, g, b byte) (int, []byte, error) {
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     SET_RGB_LEVELS,
		CmdArgs: []byte{r, g, b},
	}
	return led.write(frame)
}

// Fade to a RGB level
func (led *LED) FadeRGB(r, g, b Fade) (int, []byte, error) {
	args := []byte{}
	args = append(args, r.Args()...)
	args = append(args, g.Args()...)
	args = append(args, b.Args()...)
	frame := Frame{
		Addr:    led.Addr,
		Cmd:     FADE_RGB_TO_LEVEL,
		CmdArgs: args,
	}
	return led.write(frame)
}

// Constructs a new LED
func New(addr uint16, writer io.Writer) *LED {
	return &LED{
		Addr:   addr,
		Writer: writer,
	}
}
