package command

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"cardreader/common/status"
)

type Command struct {
	cla byte   // class
	ins byte   // instruction
	p1  byte   // param1
	p2  byte   // param2
	cmd []byte // Command
	ne  int    // expected response length
}

func (c *Command) lengthBytes(length int) []byte {
	if length < 0 {
		panic("invalid length")
	}
	var hexStrBytes = []byte(strconv.FormatUint(uint64(length), 16))
	if len(hexStrBytes)%2 > 0 {
		hexStrBytes = append([]byte{0x30}, hexStrBytes...)
	}
	var result = make([]byte, hex.DecodedLen(len(hexStrBytes)))
	if _, err := hex.Decode(result, hexStrBytes); err != nil {
		panic(err)
	}
	return result
}

func (c *Command) Bytes() ([]byte, error) {
	// @see iso7816-4-5.4.2
	// if c.ins&0x01 > 0 || c.ins&0x60 > 0 || c.ins&0x90 > 0 {
	// 	return nil, errors.New("invalid instruction")
	// }
	var cmd = []byte{c.cla, c.ins, c.p1, c.p2}
	if len(c.cmd) > 0 {
		var cmdLenBytes = c.lengthBytes(len(c.cmd))
		switch len(cmdLenBytes) {
		case 2:
			cmd = append(cmd, 0x00)
			fallthrough
		case 1, 3:
			cmd = append(cmd, cmdLenBytes...)
		default:
			return nil, errors.New("invalid Command length")
		}
		cmd = append(cmd, c.cmd...)
	}
	if c.ne > 0 {
		var neLenBytes = c.lengthBytes(c.ne)
		if neLenBytesSize := len(neLenBytes); neLenBytesSize == 0 || neLenBytesSize > 3 {
			return nil, errors.New("invalid ne length size")
		}
		cmd = append(cmd, neLenBytes...)
	} else {
		cmd = append(cmd, 0x00)
	}
	return cmd, nil
}

type handler func([]byte) ([]byte, error)

func (c *Command) Exec(h handler) (r *Response, err error) {
	var cBrief = []byte{c.cla, c.ins, c.p1, c.p2}
	var cBytes []byte
	defer func() {
		if r != nil {
			if sta := r.Status(); sta.Is(status.Nx61XXResponseBytesStillAvailable) {
				var newR *Response
				if newR, err = ReadMore(int(sta.SW2())).Exec(h); err == nil {
					r.AppendContent(newR.Content())
				}
			} else if sta.Is(status.Nx6CXXWrongLength) {
				c.ne = int(sta.SW2())
				r, err = c.Exec(h)
			}
		}
	}()
	if cBytes, err = c.Bytes(); err != nil {
		return nil, fmt.Errorf("cant build command: %02x, error: %w", cBrief, err)
	}
	var rBytes []byte
	if rBytes, err = h(cBytes); err != nil {
		return nil, fmt.Errorf("cant send command: %02x, error: %w", cBrief, err)
	}
	return ParseResponse(rBytes)
}
