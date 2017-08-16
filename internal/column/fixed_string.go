package column

import (
	"fmt"
	"reflect"

	"github.com/kshvakov/clickhouse/internal/binary"
)

type FixedString struct {
	base
	len      int
	scanType reflect.Type
}

func (str *FixedString) Read(decoder *binary.Decoder) (interface{}, error) {
	v, err := decoder.Fixed(str.len)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (str *FixedString) Write(encoder *binary.Encoder, v interface{}) error {
	var fixedString []byte
	switch v := v.(type) {
	case string:
		fixedString = binary.Str2Bytes(v)
	case []byte:
		fixedString = v
	default:
		return fmt.Errorf("unexpected type %T", v)
	}
	switch {
	case len(fixedString) > str.len:
		return fmt.Errorf("too large value '%s' (expected %d, got %d)", fixedString, str.len, len(fixedString))
	case len(fixedString) < str.len:
		tmp := make([]byte, str.len)
		copy(tmp, fixedString)
		fixedString = tmp
	}
	if _, err := encoder.Write(fixedString); err != nil {
		return err
	}
	return nil
}

func parseFixedString(name, chType string) (*FixedString, error) {
	var strLen int
	if _, err := fmt.Sscanf(chType, "FixedString(%d)", &strLen); err != nil {
		return nil, err
	}
	return &FixedString{
		base: base{
			name:     name,
			chType:   chType,
			scanType: scanTypes[string("")],
		},
		len: strLen,
	}, nil
}
