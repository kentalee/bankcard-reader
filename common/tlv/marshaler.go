package tlv

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
)

func packId(id string) (raw []byte, err error) {
	if raw, err = hex.DecodeString(id); err != nil {
		return nil, err
	}
	for offset := 0; offset < len(raw); offset++ {
		currentByte := raw[offset]
		if offset == 0 { // first tag byte
			if currentByte&0b00011111 != 0b00011111 {
				if len(raw) != 1 {
					return nil, errors.New("invalid tlv tag: unknown extra bytes")
				}
			} else if len(raw) == 1 {
				return nil, errors.New("invalid tlv tag: insufficient bytes")
			}
		} else if currentByte&0b10000000 == 0 {
			if offset < len(raw)-1 {
				return nil, errors.New("invalid tlv tag: unknown extra bytes")
			}
		} else if currentByte&0b10000000 > 0 && offset == len(raw)-1 {
			return nil, errors.New("invalid tlv tag: insufficient bytes")
		}
	}
	return raw, nil
}

func packLength(length int) (raw []byte, err error) {
	if length <= 0b01111111 {
		return []byte{byte(length)}, nil
	}
	defer func() {
		if _pErr := recover(); _pErr != nil {
			err = fmt.Errorf("panic: %v", _pErr)
		}
	}()
	switch {
	case math.MinInt8 <= length && length <= math.MaxInt8:
		raw = make([]byte, 1)
	case math.MinInt16 <= length && length <= math.MaxInt16:
		raw = make([]byte, 2)
	case math.MinInt32 <= length && length <= math.MaxInt32:
		raw = make([]byte, 4)
	default:
		raw = make([]byte, 8)
	}
	binary.PutUvarint(raw, uint64(length))
	raw = append([]byte{1<<7 | byte(len(raw))}, raw...)
	return raw, nil
}
