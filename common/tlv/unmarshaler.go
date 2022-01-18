package tlv

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func parseId(raw []byte) (id string, size int, err error) {
	var offset int
	var tagBytes []byte
	var tagMaxBytes = len(raw)
	if tagMaxBytes == 0 {
		return "", 0, ErrNoSufficientLength
	}
	for offset = 0; offset < tagMaxBytes; offset++ {
		currentByte := raw[offset]
		tagBytes = append(tagBytes, currentByte)
		if offset == 0 { // first tag byte
			if currentByte&0b00011111 != 0b00011111 {
				break
			} else if tagMaxBytes == 1 {
				return "", 0, ErrNoSufficientLength
			}
		} else if currentByte&0b10000000 == 0 {
			break
		}
	}
	return hex.EncodeToString(tagBytes), offset + 1, nil
}

func parseLength(raw []byte) (length int, size int, err error) {
	var lenBytes []byte
	var lenMaxBytes = len(raw)
	if lenMaxBytes == 0 {
		return 0, 0, ErrNoSufficientLength
	}
	if currentByte := raw[0]; currentByte&0b10000000 == 0 {
		size = 1
		lenBytes = []byte{currentByte}
	} else {
		var lengthSize = int(currentByte & 0b01111111)
		if lenMaxBytes < 1+lengthSize {
			return 0, 0, ErrNoSufficientLength
		}
		var lengthBytes = raw[1:lengthSize]
		if _size, _read := binary.Uvarint(lengthBytes); _read <= 0 {
			return 0, 0, ErrNoSufficientLength
			// todo len(decoded byte) == len(lengthBytes)
		} else {
			size = 1 + int(_size)
			lenBytes = raw[1 : 1+int(_size)]
		}
	}
	if _length, _read := binary.Uvarint(lenBytes); _read <= 0 {
		return 0, 0, errors.New("cant decode LData bytes")
		// todo len(decoded byte) == len(lengthBytes)
	} else {
		length = int(_length)
	}
	return length, size, nil
}

func Unmarshal(raw []byte, target interface{}) (n int, err error) {
	var rfValue reflect.Value
	if _rfValue, ok := target.(reflect.Value); ok {
		rfValue = _rfValue
	} else {
		rfValue = reflect.Indirect(reflect.ValueOf(target))
	}
	var rfType = rfValue.Type()
	var totalMsgLen = len(raw)
	for n < totalMsgLen {
		var _size int
		// `T` - field id
		var fieldTag string
		{
			if fieldTag, _size, err = parseId(raw[n:]); err != nil {
				return n, fmt.Errorf("%w: index: %d", err, n)
			} else if n += _size; n > totalMsgLen {
				return n, fmt.Errorf("%w: index: %d, expected: %d, raw: %02x", ErrNoSufficientLength, n, _size, raw)
			}
		}
		// `L` - field length
		var fieldLen int
		{
			if fieldLen, _size, err = parseLength(raw[n:]); err != nil {
				return n, fmt.Errorf("%w: index: %d", err, n)
			}
			if n += _size; n+fieldLen > totalMsgLen {
				return n, fmt.Errorf("%w: index: %d, expected: %d, raw: %02x", ErrNoSufficientLength, n, fieldLen, raw)
			}
		}
		// `V` - field value
		var fieldData []byte
		{
			fieldData = raw[n : n+fieldLen]
			n += fieldLen
		}
		if rfType.Kind() == reflect.Struct {
			for j := 0; j < rfValue.NumField(); j++ {
				if strings.EqualFold(rfType.Field(j).Tag.Get("tlv"), fieldTag) {
					var fieldValue = rfValue.Field(j)
					if fieldValue.Kind() == reflect.Ptr {
						fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
						fieldValue = fieldValue.Elem()
					}
					if err = store(fieldData, fieldValue); err != nil {
						return n, fmt.Errorf("%w: tag: %s", err, fieldTag)
					}
					break
				}
			}
		} else if err = store(fieldData, rfValue); err != nil {
			return n, fmt.Errorf("%w: tag: %s", err, fieldTag)
		}
	}
	return n, nil
}

func store(fieldData []byte, rfValue reflect.Value) (err error) {
	if len(fieldData) == 0 {
		return nil
	}
	rfValue = reflect.Indirect(rfValue)
	var rfType = rfValue.Type()
	switch {
	case rfType.Kind() == reflect.Ptr:
		rfValue.Elem().Set(reflect.New(rfType.Elem()))
		return store(fieldData, rfValue.Elem())
	case rfType.Kind() == reflect.Struct:
		if _, err = Unmarshal(fieldData, rfValue); err != nil {
			return err
		}
	case rfType.Kind() == reflect.String:
		rfValue.SetString(string(fieldData))
	case rfType.Kind() == reflect.Slice && rfType.Elem().Kind() == reflect.Uint8:
		rfValue.SetBytes(fieldData)
	case rfType.Kind() == reflect.Uint8:
		if len(fieldData) > 1 {
			return fmt.Errorf("data too long")
		}
		rfValue.SetUint(uint64(fieldData[0]))
	default:
		return fmt.Errorf("unsupported field type: %s", rfType)
	}
	return nil
}
