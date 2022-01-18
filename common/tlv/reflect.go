package tlv

import (
	"sync"
)

var typeCache sync.Map

//
// func ParseRfType(t reflect.Type) error {
// 	switch t.Kind() {
// 	// case reflect.Bool:
// 	// 	return boolEncoder
// 	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 	// 	return intEncoder
// 	// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 	// 	return uintEncoder
// 	// case reflect.Float32:
// 	// 	return float32Encoder
// 	// case reflect.Float64:
// 	// 	return float64Encoder
// 	case reflect.String:
// 		return stringEncoder
// 	case reflect.Interface:
// 		return interfaceEncoder
// 	case reflect.Struct:
// 		return newStructEncoder(t)
// 	case reflect.Map:
// 		return newMapEncoder(t)
// 	// case reflect.Slice:
// 	// 	return newSliceEncoder(t)
// 	// case reflect.Array:
// 	// 	return newArrayEncoder(t)
// 	case reflect.Ptr:
// 		return newPtrEncoder(t)
// 	default:
// 		return unsupportedTypeEncoder
// 	}
// }
