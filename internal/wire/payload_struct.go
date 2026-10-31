package wire

func decodeStruct(raw []byte) (map[string]any, error) {
	out := make(map[string]any)
	r := &reader{b: raw}
	for !r.done() {
		field, wt, err := r.tag()
		if err != nil {
			return nil, err
		}
		if field != 1 {
			if err := r.skip(wt); err != nil {
				return nil, err
			}
			continue
		}
		entry, err := r.bytes()
		if err != nil {
			return nil, err
		}
		k, v, err := decodeEntry(entry)
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}

func decodeEntry(raw []byte) (string, any, error) {
	var key string
	var val any
	r := &reader{b: raw}
	for !r.done() {
		field, wt, err := r.tag()
		if err != nil {
			return "", nil, err
		}
		switch field {
		case 1:
			k, err := r.bytes()
			if err != nil {
				return "", nil, err
			}
			key = string(k)
		case 2:
			vb, err := r.bytes()
			if err != nil {
				return "", nil, err
			}
			v, err := decodeValue(&reader{b: vb})
			if err != nil {
				return "", nil, err
			}
			val = v
		default:
			if err := r.skip(wt); err != nil {
				return "", nil, err
			}
		}
	}
	return key, val, nil
}

func decodeList(raw []byte) ([]any, error) {
	out := make([]any, 0, 4)
	r := &reader{b: raw}
	for !r.done() {
		field, wt, err := r.tag()
		if err != nil {
			return nil, err
		}
		if field != 1 {
			if err := r.skip(wt); err != nil {
				return nil, err
			}
			continue
		}
		vb, err := r.bytes()
		if err != nil {
			return nil, err
		}
		v, err := decodeValue(&reader{b: vb})
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}
