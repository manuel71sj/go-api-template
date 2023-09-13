package str

import (
	"encoding/json"
	"strconv"
	"unsafe"
)

type S string

func (s S) String() string {
	return string(s)
}

func (s S) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func (s S) Bool() (bool, error) {
	b, err := strconv.ParseBool(s.String())
	if err != nil {
		return false, err
	}

	return b, nil
}

func (s S) DefaultBool(defaultVal bool) bool {
	b, err := s.Bool()
	if err != nil {
		return defaultVal
	}

	return b
}

func (s S) Int64() (int64, error) {
	i, err := strconv.ParseInt(s.String(), 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (s S) DefaultInt64(defaultVal int64) int64 {
	i, err := s.Int64()
	if err != nil {
		return defaultVal
	}

	return i
}

func (s S) Int() (int, error) {
	i, err := s.Int64()
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func (s S) DefaultInt(defaultVal int) int {
	i, err := s.Int()
	if err != nil {
		return defaultVal
	}

	return i
}

func (s S) Uint64() (uint64, error) {
	i, err := strconv.ParseUint(s.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (s S) DefaultUint64(defaultVal uint64) uint64 {
	i, err := s.Uint64()
	if err != nil {
		return defaultVal
	}
	return i
}

func (s S) Uint() (uint, error) {
	i, err := s.Uint64()
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func (s S) DefaultUint(defaultVal uint) uint {
	i, err := s.Uint()
	if err != nil {
		return defaultVal
	}
	return uint(i)
}

func (s S) Float64() (float64, error) {
	f, err := strconv.ParseFloat(s.String(), 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func (s S) DefaultFloat64(defaultVal float64) float64 {
	f, err := s.Float64()
	if err != nil {
		return defaultVal
	}
	return f
}

func (s S) Float32() (float32, error) {
	f, err := s.Float64()
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func (s S) DefaultFloat32(defaultVal float32) float32 {
	f, err := s.Float32()
	if err != nil {
		return defaultVal
	}
	return f
}

func (s S) ToJSON(v interface{}) error {
	return json.Unmarshal(s.Bytes(), v)
}

func NewWithByte(b []byte) S {
	return *(*S)(unsafe.Pointer(&b))
}
