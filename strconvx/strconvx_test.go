package strconvx_test

import (
	"math"
	"testing"

	"github.com/yu1ec/go-pkg/strconvx"
)

func TestToNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType any
		want     any
		wantErr  bool
	}{
		{"Valid int", "123", int(0), 123, false},
		{"Valid int8 max", "127", int8(0), int8(127), false},
		{"Valid int8 min", "-128", int8(0), int8(-128), false},
		{"Valid int16", "32767", int16(0), int16(32767), false},
		{"Valid int32", "2147483647", int32(0), int32(2147483647), false},
		{"Valid int64", "9223372036854775807", int64(0), int64(9223372036854775807), false},
		{"Valid uint", "123", uint(0), uint(123), false},
		{"Valid uint8", "255", uint8(0), uint8(255), false},
		{"Valid uint16", "65535", uint16(0), uint16(65535), false},
		{"Valid uint32", "4294967295", uint32(0), uint32(4294967295), false},
		{"Valid uint64", "18446744073709551615", uint64(0), uint64(18446744073709551615), false},
		{"Valid float32", "3.14", float32(0), float32(3.14), false},
		{"Valid float64", "3.141592653589793", float64(0), 3.141592653589793, false},
		{"Invalid number", "not a number", int(0), 0, true},
		{"Empty string", "", int(0), 0, true},
		{"Overflow int8", "128", int8(0), int8(0), true},
		{"Underflow int8", "-129", int8(0), int8(0), true},
		{"Overflow uint8", "256", uint8(0), uint8(0), true},
		{"Negative for uint", "-1", uint(0), uint(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.wantType.(type) {
			case int:
				got, err := strconvx.ToNumber[int](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case int8:
				got, err := strconvx.ToNumber[int8](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case int16:
				got, err := strconvx.ToNumber[int16](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case int32:
				got, err := strconvx.ToNumber[int32](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case int64:
				got, err := strconvx.ToNumber[int64](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case uint:
				got, err := strconvx.ToNumber[uint](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case uint8:
				got, err := strconvx.ToNumber[uint8](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case uint16:
				got, err := strconvx.ToNumber[uint16](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case uint32:
				got, err := strconvx.ToNumber[uint32](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case uint64:
				got, err := strconvx.ToNumber[uint64](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case float32:
				got, err := strconvx.ToNumber[float32](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			case float64:
				got, err := strconvx.ToNumber[float64](tt.input)
				checkResult(t, got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func checkResult(t *testing.T, got any, err error, want any, wantErr bool) {
	if (err != nil) != wantErr {
		t.Errorf("ToNumber() error = %v, wantErr %v", err, wantErr)
		return
	}
	if !wantErr && got != want {
		t.Errorf("ToNumber() = %v, want %v", got, want)
	}
}

func TestToNumberOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue any
		want         any
	}{
		{"Valid int", "123", 0, 123},
		{"Invalid int", "not a number", 42, 42},
		{"Valid float64", "3.14", 0.0, 3.14},
		{"Invalid float64", "not a number", 3.14, 3.14},
		{"Empty string for int", "", 100, 100},
		{"Empty string for float", "", 1.23, 1.23},
		{"Overflow int8", "128", int8(127), int8(127)},
		{"Underflow int8", "-129", int8(-128), int8(-128)},
		{"Max int8", "127", int8(0), int8(127)},
		{"Min int8", "-128", int8(0), int8(-128)},
		{"Underflow uint", "-1", uint(0), uint(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.defaultValue.(type) {
			case int:
				got := strconvx.ToNumberOrDefault(tt.input, tt.defaultValue.(int))
				if got != tt.want {
					t.Errorf("ToNumberOrDefault() = %v, want %v", got, tt.want)
				}
			case float64:
				got := strconvx.ToNumberOrDefault(tt.input, tt.defaultValue.(float64))
				if got != tt.want {
					t.Errorf("ToNumberOrDefault() = %v, want %v", got, tt.want)
				}
			case int8:
				got := strconvx.ToNumberOrDefault(tt.input, tt.defaultValue.(int8))
				if got != tt.want {
					t.Errorf("ToNumberOrDefault() = %v, want %v", got, tt.want)
				}
			case uint:
				got := strconvx.ToNumberOrDefault(tt.input, tt.defaultValue.(uint))
				if got != tt.want {
					t.Errorf("ToNumberOrDefault() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestToNumberEdgeCases(t *testing.T) {
	maxInt8, err := strconvx.ToNumber[int8]("127")
	if err != nil || maxInt8 != math.MaxInt8 {
		t.Errorf("Failed to convert max int8")
	}

	minInt8, err := strconvx.ToNumber[int8]("-128")
	if err != nil || minInt8 != math.MinInt8 {
		t.Errorf("Failed to convert min int8")
	}

	maxInt64, err := strconvx.ToNumber[int64]("9223372036854775807")
	if err != nil || maxInt64 != math.MaxInt64 {
		t.Errorf("Failed to convert max int64")
	}

	minInt64, err := strconvx.ToNumber[int64]("-9223372036854775808")
	if err != nil || minInt64 != math.MinInt64 {
		t.Errorf("Failed to convert min int64")
	}

	maxUint64, err := strconvx.ToNumber[uint64]("18446744073709551615")
	if err != nil || maxUint64 != math.MaxUint64 {
		t.Errorf("Failed to convert max uint64")
	}

	maxFloat64, err := strconvx.ToNumber[float64]("1.7976931348623157e+308")
	if err != nil || maxFloat64 != math.MaxFloat64 {
		t.Errorf("Failed to convert max float64")
	}
}
