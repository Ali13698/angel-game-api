package wire

import (
	"fmt"
	"reflect"
)

const (
	valNull   = 1
	valNumber = 2
	valString = 3
	valBool   = 4
	valStruct = 5
	valList   = 6
)

const maxExactInt = 1 << 53

func EncodeValue(v any) []byte { return appendValue(nil, v) }

func EncodeValueField(b []byte, field int, v any) []byte {
	return appendBytesField(b, field, appendValue(nil, v))
}

func appendValue(b []byte, v any) []byte {
	if v == nil {
		return appendVarintField(b, valNull, 0)
	}
	if raw, ok := v.([]byte); ok {
		return appendStringField(b, valString, string(raw))
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		var n uint64
		if rv.Bool() {
			n = 1
		}
		return appendVarintField(b, valBool, n)
	case reflect.String:
		return appendStringField(b, valString, rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return appendDoubleField(b, valNumber, float64(rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return appendDoubleField(b, valNumber, float64(rv.Uint()))
	case reflect.Float32, reflect.Float64:
		return appendDoubleField(b, valNumber, rv.Float())
	case reflect.Slice, reflect.Array:
		return appendList(b, rv)
	case reflect.Map:
		return appendStruct(b, rv)
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return appendVarintField(b, valNull, 0)
		}
		return appendValue(b, rv.Elem().Interface())
	default:
		return appendStringField(b, valString, fmt.Sprint(v))
	}
}

func appendStruct(b []byte, rv reflect.Value) []byte {
	var inner []byte
	iter := rv.MapRange()
	for iter.Next() {
		k := iter.Key()
		if k.Kind() != reflect.String {
			continue
		}
		var entry []byte
		entry = appendStringField(entry, 1, k.String())
		entry = appendBytesField(entry, 2, appendValue(nil, iter.Value().Interface()))
		inner = appendBytesField(inner, 1, entry)
	}
	return appendBytesField(b, valStruct, inner)
}

func appendList(b []byte, rv reflect.Value) []byte {
	var inner []byte
	n := rv.Len()
	for i := 0; i < n; i++ {
		inner = appendBytesField(inner, 1, appendValue(nil, rv.Index(i).Interface()))
	}
	return appendBytesField(b, valList, inner)
}
