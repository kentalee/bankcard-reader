package status

import (
	"bytes"
	"errors"
	"fmt"
)

type toBytes interface {
	Content() []byte
}

type Level uint8

const (
	Normal Level = iota
	Retry
	Warning
	ExecutionError
	CheckingError
)

type Status struct {
	sw      [2]byte
	arg2    bool
	level   Level
	comment string
}

func (e Status) Content() []byte {
	return e.sw[:]
}

func (e Status) Error() string {
	return fmt.Sprintf("%02x", e.sw)
}

func (e Status) Level() Level {
	return e.level
}

func (e Status) Is(ne error) bool {
	if nex, ok := ne.(toBytes); ok {
		if e.arg2 {
			return bytes.Equal(nex.Content()[:1], e.sw[:1])
		} else {
			return bytes.Equal(nex.Content(), e.sw[:])
		}
	}
	return false
}

func (e Status) SW1() byte {
	return e.sw[0]
}

func (e Status) SW2() byte {
	return e.sw[1]
}

func Parse(raw []byte) (e Status, err error) {
	if len(raw) != 2 {
		return Status{}, errors.New("invalid length")
	}
	defer func() {
		if err == nil && e.arg2 {
			e.sw[1] = raw[1]
		}
	}()
	switch sw1, sw2 := raw[0], raw[1]; {
	case sw1 == 0x90 && sw2 == 0x00: // No further qualification
		return Nx9000NormalProcessing, nil
	case sw1 == 0x61: // SW2 indicates the number of response bytes still available
		return Nx61XXResponseBytesStillAvailable, nil
	case sw1 == 0x6C: // Wrong length Le: SW2 indicates the exact length (see text below)
		return Nx6CXXWrongLength, nil
	case sw1 == 0x62: // State of non-volatile memory unchanged (further qualification in SW2, see table 13)
		return parse62(sw2)
	case sw1 == 0x63: // State of non-volatile memory changed (further qualification in SW2, see table 14)
		return parse63(sw2)
	case sw1 == 0x64: // State of non-volatile memory unchanged (SW2 = ’00’, other values are RFU)
		return parse64(sw2)
	case sw1 == 0x65: // State of non-volatile memory changed (further qualification in SW2, see table 15)
		return parse65(sw2)
	case sw1 == 0x66: // Reserved for security-related issues (not defined in this part of ISO/IEC 7816)
		return Ex66XXSecurityReserved, nil
	case sw1 == 0x67 && sw2 == 0x00: // Wrong length
		return Ex6700WrongLength, nil
	case sw1 == 0x68: // Functions in CLA not supported (further qualification in SW2, see table 16)
		return parse68(sw2)
	case sw1 == 0x69: // Command not allowed (further qualification in SW2, see table 17)
		return parse69(sw2)
	case sw1 == 0x6A: // Wrong parameter(s) P1-P2 (further qualification in SW2, see table 18)
		return parse6A(sw2)
	case sw1 == 0x6B && sw2 == 0x00: // Wrong parameter(s) P1-P2
		return Ex6B00WrongParameter, nil
	case sw1 == 0x6D && sw2 == 0x00: // Instruction code not supported or invalid
		return Ex6D00InstructionCodeErr, nil
	case sw1 == 0x6E && sw2 == 0x00: // Class not supported
		return Ex6E00ClassNotSupported, nil
	case sw1 == 0x6F && sw2 == 0x00: // No precise diagnosis
		return Ex6F00NoPreciseDiagnosis, nil
	default:
		return NewUnknownError(sw1, sw2)
	}
}

func parse62(sw2 byte) (e Status, err error) {
	switch sw2 {
	// ’00’	No information given
	// ’81’	Part of returned data may be corrupted
	// ’82’	End of file/record reached before reading Le bytes
	// ’83’	Selected file invalidated
	// ’84’	FCI not formatted according to 1.1.5
	}
	return NewUnknownError(0x62, sw2)
}

func parse63(sw2 byte) (e Status, err error) {
	switch sw2 {
	// ’00’	No information given
	// ’81’	File filled up by the last write
	// ‘CX’	Counter provided by ‘X’ (valued from 0 to 15) (exact meaning depending on the command)
	}
	return NewUnknownError(0x63, sw2)
}

func parse64(sw2 byte) (e Status, err error) {
	if sw2 == 0x00 {
		return Ex6400Unchanged, nil
	}
	return Ex64XXRefused, nil
}

func parse65(sw2 byte) (e Status, err error) {
	switch sw2 {
	// ’00’	No information given
	// ’81’	Memory failure
	}
	return NewUnknownError(0x65, sw2)
}

func parse68(sw2 byte) (e Status, err error) {
	switch sw2 {
	// ’00’	No information given
	// ’81’	Logical channel not supported
	// ’82’	Secure messaging not supported
	}
	return NewUnknownError(0x68, sw2)
}

func parse69(sw2 byte) (e Status, err error) {
	switch sw2 {
	// ’00’	No information given
	// ’81’	Command incompatible with file structure
	// ’82’	Security status not satisfied
	// ’83’	Authentication method blocked
	// ’84’	Referenced data invalidated
	// ’85’	Conditions of use not satisfied
	// ’86’	Command not allowed (no current EF)
	// ’87’	Expected SM data objects missing
	// ’88’	SM data objects incorrect
	}
	return NewUnknownError(0x69, sw2)
}

func parse6A(sw2 byte) (e Status, err error) {
	switch sw2 {
	// ’00’	No information given
	// ’80’	Incorrect parameters in the data field
	// ’81’	Function not supported
	// ’82’	File not found
	// ’83’	Record not found
	// ’84’	Not enough memory space in the file
	// ’85’	Lc inconsistent with TLV structure
	// ’86’	Incorrect parameters P1-P2
	// ’87’	Lc inconsistent with P1-P2
	// ’88’
	// Referenced data not found
	}
	return NewUnknownError(0x6A, sw2)
}
