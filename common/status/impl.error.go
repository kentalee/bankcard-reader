package status

import (
	"errors"
	"fmt"
)

var UnknownError = errors.New("unknown error")

func NewUnknownError(sw1, sw2 byte) (Status, error) {
	return Status{
		level:   ExecutionError,
		sw:      [2]byte{sw1, sw2},
		comment: "Unknown Error",
	}, fmt.Errorf("%w: %02x%02x", UnknownError, sw1, sw2)
}

var Ex6A82FileNotFound = Status{
	sw:      [2]byte{0x6a, 0x82},
	comment: "File Not Found",
}

var Ex6400Unchanged = Status{
	level:   ExecutionError,
	sw:      [2]byte{0x64, 0x00},
	comment: "State of non-volatile memory unchanged",
}

var Ex64XXRefused = Status{
	arg2:    true,
	level:   ExecutionError,
	sw:      [2]byte{0x64, 0x00},
	comment: "State of non-volatile memory unchanged(RFU)",
}

var Ex66XXSecurityReserved = Status{
	arg2:    true,
	level:   ExecutionError,
	sw:      [2]byte{0x66, 0x00},
	comment: "Reserved for security-related issues",
}

var Ex6700WrongLength = Status{
	level:   CheckingError,
	sw:      [2]byte{0x67, 0x00},
	comment: "Wrong length",
}

var Ex6B00WrongParameter = Status{
	level:   CheckingError,
	sw:      [2]byte{0x6B, 0x00},
	comment: "Wrong parameter(s) P1-P2",
}

var Ex6D00InstructionCodeErr = Status{
	level:   CheckingError,
	sw:      [2]byte{0x6D, 0x00},
	comment: "Instruction code not supported or invalid",
}

var Ex6E00ClassNotSupported = Status{
	level:   CheckingError,
	sw:      [2]byte{0x6E, 0x00},
	comment: "Class not supported",
}

var Ex6F00NoPreciseDiagnosis = Status{
	level:   CheckingError,
	sw:      [2]byte{0x6F, 0x00},
	comment: "No precise diagnosis",
}
