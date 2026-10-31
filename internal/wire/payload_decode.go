package wire

import "math"

func DecodeValue(raw []byte) (any, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	r := &reader{b: raw}
	return decodeValue(r)
}

func DecodeMap(raw []byte) (map[string]any, error) {
	v, err := DecodeValue(raw)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return map[string]any{}, nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return map[string]any{}, nil
	}
	return m, nil
}

func decodeValue(r *reader) (any, error) {
	for !r.done() {
		field, wt, err := r.tag()
		if err != nil {
			return nil, err
		}
		switch field {
		case valNull:
			if _, err := r.varint(); err != nil {
				return nil, err
			}
			return nil, nil
		case valNumber:
			u, err := r.fixed64()
			if err != nil {
				return nil, err
			}
			return normalizeNumber(math.Float64frombits(u)), nil
		case valString:
			raw, err := r.bytes()
			if err != nil {
				return nil, err
			}
			return string(raw), nil
		case valBool:
			n, err := r.varint()
			if err != nil {
				return nil, err
			}
			return n != 0, nil
		case valStruct:
			raw, err := r.bytes()
			if err != nil {
				return nil, err
			}
			return decodeStruct(raw)
		case valList:
			raw, err := r.bytes()
			if err != nil {
				return nil, err
			}
			return decodeList(raw)
		default:
			if err := r.skip(wt); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func normalizeNumber(f float64) any {
	if f == math.Trunc(f) && math.Abs(f) < maxExactInt {
		return int64(f)
	}
	return f
}
