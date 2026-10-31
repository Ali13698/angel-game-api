package wire

import (
	"encoding/binary"
	"errors"
	"math"
)

const (
	wtVarint  = 0
	wtFixed64 = 1
	wtBytes   = 2
	wtFixed32 = 5
)

var ErrCorrupt = errors.New("wire: corrupt payload")

func appendVarint(b []byte, v uint64) []byte {
	for v >= 0x80 {
		b = append(b, byte(v)|0x80)
		v >>= 7
	}
	return append(b, byte(v))
}

func appendTag(b []byte, field, wt int) []byte {
	return appendVarint(b, uint64(field)<<3|uint64(wt))
}

func appendBytesField(b []byte, field int, raw []byte) []byte {
	b = appendTag(b, field, wtBytes)
	b = appendVarint(b, uint64(len(raw)))
	return append(b, raw...)
}

func appendStringField(b []byte, field int, s string) []byte {
	b = appendTag(b, field, wtBytes)
	b = appendVarint(b, uint64(len(s)))
	return append(b, s...)
}

func appendVarintField(b []byte, field int, v uint64) []byte {
	b = appendTag(b, field, wtVarint)
	return appendVarint(b, v)
}

func appendDoubleField(b []byte, field int, f float64) []byte {
	b = appendTag(b, field, wtFixed64)
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], math.Float64bits(f))
	return append(b, tmp[:]...)
}

type reader struct {
	b []byte
	i int
}

func (r *reader) done() bool { return r.i >= len(r.b) }

func (r *reader) varint() (uint64, error) {
	var v uint64
	for shift := uint(0); shift < 64; shift += 7 {
		if r.i >= len(r.b) {
			return 0, ErrCorrupt
		}
		c := r.b[r.i]
		r.i++
		v |= uint64(c&0x7f) << shift
		if c < 0x80 {
			return v, nil
		}
	}
	return 0, ErrCorrupt
}

func (r *reader) tag() (field, wt int, err error) {
	k, err := r.varint()
	if err != nil {
		return 0, 0, err
	}
	field = int(k >> 3)
	if field <= 0 {
		return 0, 0, ErrCorrupt
	}
	return field, int(k & 7), nil
}

func (r *reader) bytes() ([]byte, error) {
	n, err := r.varint()
	if err != nil {
		return nil, err
	}
	if n > uint64(len(r.b)-r.i) {
		return nil, ErrCorrupt
	}
	out := r.b[r.i : r.i+int(n)]
	r.i += int(n)
	return out, nil
}

func (r *reader) fixed64() (uint64, error) {
	if len(r.b)-r.i < 8 {
		return 0, ErrCorrupt
	}
	v := binary.LittleEndian.Uint64(r.b[r.i : r.i+8])
	r.i += 8
	return v, nil
}

func (r *reader) skip(wt int) error {
	switch wt {
	case wtVarint:
		_, err := r.varint()
		return err
	case wtFixed64:
		_, err := r.fixed64()
		return err
	case wtBytes:
		_, err := r.bytes()
		return err
	case wtFixed32:
		if len(r.b)-r.i < 4 {
			return ErrCorrupt
		}
		r.i += 4
		return nil
	default:
		return ErrCorrupt
	}
}
