package gmvc

import (
	"strconv"
)

type PathVars map[string]string

func (p PathVars) Get(key string) string {
	if p == nil {
		return ""
	}
	return p[key]
}

// string

func (p PathVars) String(key string) string {
	return p.Get(key)
}

// integer

func (p PathVars) getInt(key string, bitsize int) int64 {
	s, ok := p[key]
	if !ok || s == "" {
		return 0
	}

	v, err := strconv.ParseInt(s, 10, bitsize)
	if err != nil {
		return 0
	}

	return v
}

func (p PathVars) getUint(key string, bitsize int) uint64 {
	s, ok := p[key]
	if !ok || s == "" {
		return 0
	}

	v, err := strconv.ParseUint(s, 10, bitsize)
	if err != nil {
		return 0
	}

	return v
}

//Int

func (p PathVars) Int(key string) int {
	return int(p.getInt(key, 0))
}

// Int8

func (p PathVars) Int8(key string) int8 {
	return int8(p.getInt(key, 8))
}

// Int16

func (p PathVars) Int16(key string) int16 {
	return int16(p.getInt(key, 16))
}

// Int32

func (p PathVars) Int32(key string) int32 {
	return int32(p.getInt(key, 32))
}

// Int64

func (p PathVars) Int64(key string) int64 {
	return p.getInt(key, 64)
}

// uint

func (p PathVars) Uint(key string) uint {
	return uint(p.getUint(key, 0))
}

// Uint8

func (p PathVars) Uint8(key string) uint8 {
	return uint8(p.getUint(key, 8))
}

// Uint16

func (p PathVars) Uint16(key string) uint16 {
	return uint16(p.getUint(key, 16))
}

// Uint32

func (p PathVars) Uint32(key string) uint32 {
	return uint32(p.getUint(key, 32))
}

// Uint64

func (p PathVars) Uint64(key string) uint64 {
	return p.getUint(key, 64)
}

// float

func (p PathVars) getFloat(key string, bitsize int) float64 {
	s, ok := p[key]
	if !ok || s == "" {
		return 0
	}

	v, err := strconv.ParseFloat(s, bitsize)
	if err != nil {
		return 0
	}

	return v
}

// Float32

func (p PathVars) Float32(key string) float32 {
	return float32(p.getFloat(key, 32))
}

// Float64

func (p PathVars) Float64(key string) float64 {
	return p.getFloat(key, 64)
}

// bool

func (p PathVars) Bool(key string) bool {
	s, ok := p[key]
	if !ok || s == "" {
		return false
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return v
}
