package status

var Nx9000NormalProcessing = Status{
	level:   Normal,
	sw:      [2]byte{0x90, 0x00},
	comment: "Normal processing",
}

var Nx61XXResponseBytesStillAvailable = Status{
	arg2:    true,
	level:   Retry,
	sw:      [2]byte{0x61, 0x00},
	comment: "Response Bytes Still Available",
}

var Nx6CXXWrongLength = Status{
	arg2:    true,
	level:   Retry,
	sw:      [2]byte{0x6C, 0x00},
	comment: "Wrong length Le: SW2 indicates the exact length",
}
