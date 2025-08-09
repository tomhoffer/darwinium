package utils

import (
	"fmt"
	"strconv"
)

func ConvertToFloat64(v any) (float64, error) {
	switch n := any(v).(type) {
	case float32:
		return float64(n), nil
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint8:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case uintptr:
		return float64(n), nil
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string %q to float64: %w", n, err)
		}
		return f, nil
	default:
		s := fmt.Sprint(v)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %T to float64: %w", v, err)
		}
		return f, nil
	}
}
