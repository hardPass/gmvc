package gmvc

import (
	"net/url"
)

type Form url.Values

//Int

func (f Form) Int(key string) int {
	return values.Int(f, key)
}

func (f Form) GetInt(key string, def int) (int, bool, error) {
	return values.GetInt(f, key, def)
}

func (f Form) GetInts(key string) ([]int, bool, error) {
	return values.GetInts(f, key)
}

// Int8

func (f Form) Int8(key string) int8 {
	return values.Int8(f, key)
}

func (f Form) GetInt8(key string, def int8) (int8, bool, error) {
	return values.GetInt8(f, key, def)
}

func (f Form) GetInt8s(key string) ([]int8, bool, error) {
	return values.GetInt8s(f, key)
}

// Int16

func (f Form) Int16(key string) int16 {
	return values.Int16(f, key)
}

func (f Form) GetInt16(key string, def int16) (int16, bool, error) {
	return values.GetInt16(f, key, def)
}

func (f Form) GetInt16s(key string) ([]int16, bool, error) {
	return values.GetInt16s(f, key)
}

// Int32

func (f Form) Int32(key string) int32 {
	return values.Int32(f, key)
}

func (f Form) GetInt32(key string, def int32) (int32, bool, error) {
	return values.GetInt32(f, key, def)
}

func (f Form) GetInt32s(key string) ([]int32, bool, error) {
	return values.GetInt32s(f, key)
}

// Int64

func (f Form) Int64(key string) int64 {
	return values.Int64(f, key)
}

func (f Form) GetInt64(key string, def int64) (int64, bool, error) {
	return values.GetInt64(f, key, def)
}

func (f Form) GetInt64s(key string) ([]int64, bool, error) {
	return values.GetInt64s(f, key)
}

// uint

func (f Form) Uint(key string) uint {
	return values.Uint(f, key)
}

func (f Form) GetUint(key string, def uint) (uint, bool, error) {
	return values.GetUint(f, key, def)
}

func (f Form) GetUints(key string) ([]uint, bool, error) {
	return values.GetUints(f, key)
}

// Uint8

func (f Form) Uint8(key string) uint8 {
	return values.Uint8(f, key)
}

func (f Form) GetUint8(key string, def uint8) (uint8, bool, error) {
	return values.GetUint8(f, key, def)
}

func (f Form) GetUint8s(key string) ([]uint16, bool, error) {
	return values.GetUint8s(f, key)
}

// Uint16

func (f Form) Uint16(key string) uint16 {
	return values.Uint16(f, key)
}

func (f Form) GetUint16(key string, def uint16) (uint16, bool, error) {
	return values.GetUint16(f, key, def)
}

func (f Form) GetUint16s(key string) ([]uint16, bool, error) {
	return values.GetUint16s(f, key)
}

// Uint32

func (f Form) Uint32(key string) uint32 {
	return values.Uint32(f, key)
}

func (f Form) GetUint32(key string, def uint32) (uint32, bool, error) {
	return values.GetUint32(f, key, def)
}

func (f Form) GetUint32s(key string) ([]uint32, bool, error) {
	return values.GetUint32s(f, key)
}

// Uint64

func (f Form) Uint64(key string) uint64 {
	return values.Uint64(f, key)
}

func (f Form) GetUint64(key string, def uint64) (uint64, bool, error) {
	return values.GetUint64(f, key, def)
}

func (f Form) GetUint64s(key string) ([]uint64, bool, error) {
	return values.GetUint64s(f, key)
}

// Float32

func (f Form) Float32(key string) float32 {
	return values.Float32(f, key)
}

func (f Form) GetFloat32(key string, def float32) (float32, bool, error) {
	return values.GetFloat32(f, key, def)
}

func (f Form) GetFloat32s(key string) ([]float32, bool, error) {
	return values.GetFloat32s(f, key)
}

// Float64

func (f Form) Float64(key string) float64 {
	return values.Float64(f, key)
}

func (f Form) GetFloat64(key string, def float64) (float64, bool, error) {
	return values.GetFloat64(f, key, def)
}

func (f Form) GetFloat64s(key string) ([]float64, bool, error) {
	return values.GetFloat64s(f, key)
}

// Bool

func (f Form) Bool(key string) bool {
	return values.Bool(f, key)
}

func (f Form) GetBool(key string, def bool) (bool, bool, error) {
	return values.GetBool(f, key, def)
}

func (f Form) GetBools(key string) ([]bool, bool, error) {
	return values.GetBools(f, key)
}

// string

func (f Form) String(key string) string {
	return values.String(f, key)
}

func (f Form) GetString(key string, def string) (string, bool, error) {
	return values.GetString(f, key, def)
}

func (f Form) GetStrings(key string) ([]string, bool, error) {
	return values.GetStrings(f, key)
}
