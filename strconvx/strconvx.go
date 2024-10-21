package strconvx

import (
	"fmt"
	"strconv"
)

// Number 是一个约束，限制了可以使用的数字类型
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// ToNumber 尝试将字符串转换为指定的数字类型
func ToNumber[T Number](s string) (T, error) {
	var result T
	switch any(result).(type) {
	case int8:
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return result, fmt.Errorf("failed to convert to int8: %w", err)
		}
		return T(v), nil
	case int16:
		v, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return result, fmt.Errorf("failed to convert to int16: %w", err)
		}
		return T(v), nil
	case int32:
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return result, fmt.Errorf("failed to convert to int32: %w", err)
		}
		return T(v), nil
	case int, int64:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return result, fmt.Errorf("failed to convert to int64: %w", err)
		}
		return T(v), nil
	case uint8:
		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return result, fmt.Errorf("failed to convert to uint8: %w", err)
		}
		return T(v), nil
	case uint16:
		v, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return result, fmt.Errorf("failed to convert to uint16: %w", err)
		}
		return T(v), nil
	case uint32:
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return result, fmt.Errorf("failed to convert to uint32: %w", err)
		}
		return T(v), nil
	case uint, uint64:
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return result, fmt.Errorf("failed to convert to uint64: %w", err)
		}
		return T(v), nil
	case float32:
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return result, fmt.Errorf("failed to convert to float32: %w", err)
		}
		return T(v), nil
	case float64:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return result, fmt.Errorf("failed to convert to float64: %w", err)
		}
		return T(v), nil
	default:
		return result, fmt.Errorf("unsupported number type")
	}
}

// ToNumberOrDefault 尝试将字符串转换为指定的数字类型
// 如果转换失败，返回提供的默认值
func ToNumberOrDefault[T Number](s string, defaultValue T) T {
	result, err := ToNumber[T](s)
	if err != nil {
		return defaultValue
	}
	return result
}
