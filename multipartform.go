package gmvc

import (
	"mime/multipart"
)

type MultipartForm multipart.Form

//Int

func (f *MultipartForm) Int(key string) int {
	return values.Int(f.Value, key)
}

func (f *MultipartForm) GetInt(key string, def int) (int, bool, error) {
	return values.GetInt(f.Value, key, def)
}

func (f *MultipartForm) GetInts(key string) ([]int, bool, error) {
	return values.GetInts(f.Value, key)
}

// Int8

func (f *MultipartForm) Int8(key string) int8 {
	return values.Int8(f.Value, key)
}

func (f *MultipartForm) GetInt8(key string, def int8) (int8, bool, error) {
	return values.GetInt8(f.Value, key, def)
}

func (f *MultipartForm) GetInt8s(key string) ([]int8, bool, error) {
	return values.GetInt8s(f.Value, key)
}

// Int16

func (f *MultipartForm) Int16(key string) int16 {
	return values.Int16(f.Value, key)
}

func (f *MultipartForm) GetInt16(key string, def int16) (int16, bool, error) {
	return values.GetInt16(f.Value, key, def)
}

func (f *MultipartForm) GetInt16s(key string) ([]int16, bool, error) {
	return values.GetInt16s(f.Value, key)
}

// Int32

func (f *MultipartForm) Int32(key string) int32 {
	return values.Int32(f.Value, key)
}

func (f *MultipartForm) GetInt32(key string, def int32) (int32, bool, error) {
	return values.GetInt32(f.Value, key, def)
}

func (f *MultipartForm) GetInt32s(key string) ([]int32, bool, error) {
	return values.GetInt32s(f.Value, key)
}

// Int64

func (f *MultipartForm) Int64(key string) int64 {
	return values.Int64(f.Value, key)
}

func (f *MultipartForm) GetInt64(key string, def int64) (int64, bool, error) {
	return values.GetInt64(f.Value, key, def)
}

func (f *MultipartForm) GetInt64s(key string) ([]int64, bool, error) {
	return values.GetInt64s(f.Value, key)
}

// uint

func (f *MultipartForm) Uint(key string) uint {
	return values.Uint(f.Value, key)
}

func (f *MultipartForm) GetUint(key string, def uint) (uint, bool, error) {
	return values.GetUint(f.Value, key, def)
}

func (f *MultipartForm) GetUints(key string) ([]uint, bool, error) {
	return values.GetUints(f.Value, key)
}

// Uint8

func (f *MultipartForm) Uint8(key string) uint8 {
	return values.Uint8(f.Value, key)
}

func (f *MultipartForm) GetUint8(key string, def uint8) (uint8, bool, error) {
	return values.GetUint8(f.Value, key, def)
}

func (f *MultipartForm) GetUint8s(key string) ([]uint16, bool, error) {
	return values.GetUint8s(f.Value, key)
}

// Uint16

func (f *MultipartForm) Uint16(key string) uint16 {
	return values.Uint16(f.Value, key)
}

func (f *MultipartForm) GetUint16(key string, def uint16) (uint16, bool, error) {
	return values.GetUint16(f.Value, key, def)
}

func (f *MultipartForm) GetUint16s(key string) ([]uint16, bool, error) {
	return values.GetUint16s(f.Value, key)
}

// Uint32

func (f *MultipartForm) Uint32(key string) uint32 {
	return values.Uint32(f.Value, key)
}

func (f *MultipartForm) GetUint32(key string, def uint32) (uint32, bool, error) {
	return values.GetUint32(f.Value, key, def)
}

func (f *MultipartForm) GetUint32s(key string) ([]uint32, bool, error) {
	return values.GetUint32s(f.Value, key)
}

// Uint64

func (f *MultipartForm) Uint64(key string) uint64 {
	return values.Uint64(f.Value, key)
}

func (f *MultipartForm) GetUint64(key string, def uint64) (uint64, bool, error) {
	return values.GetUint64(f.Value, key, def)
}

func (f *MultipartForm) GetUint64s(key string) ([]uint64, bool, error) {
	return values.GetUint64s(f.Value, key)
}

// Float32

func (f *MultipartForm) Float32(key string) float32 {
	return values.Float32(f.Value, key)
}

func (f *MultipartForm) GetFloat32(key string, def float32) (float32, bool, error) {
	return values.GetFloat32(f.Value, key, def)
}

func (f *MultipartForm) GetFloat32s(key string) ([]float32, bool, error) {
	return values.GetFloat32s(f.Value, key)
}

// Float64

func (f *MultipartForm) Float64(key string) float64 {
	return values.Float64(f.Value, key)
}

func (f *MultipartForm) GetFloat64(key string, def float64) (float64, bool, error) {
	return values.GetFloat64(f.Value, key, def)
}

func (f *MultipartForm) GetFloat64s(key string) ([]float64, bool, error) {
	return values.GetFloat64s(f.Value, key)
}

// Bool

func (f *MultipartForm) Bool(key string) bool {
	return values.Bool(f.Value, key)
}

func (f *MultipartForm) GetBool(key string, def bool) (bool, bool, error) {
	return values.GetBool(f.Value, key, def)
}

func (f *MultipartForm) GetBools(key string) ([]bool, bool, error) {
	return values.GetBools(f.Value, key)
}

// string

func (f *MultipartForm) String(key string) string {
	return values.String(f.Value, key)
}

func (f *MultipartForm) GetString(key string, def string) (string, bool, error) {
	return values.GetString(f.Value, key, def)
}

func (f *MultipartForm) GetStrings(key string) ([]string, bool, error) {
	return values.GetStrings(f.Value, key)
}

// file

func (f *MultipartForm) GetFile(key string) (*multipart.FileHeader, bool) {
	if f == nil {
		return nil, false
	}

	fhs, ok := f.File[key]
	if !ok || len(fhs) == 0 {
		return nil, ok
	}

	return fhs[0], ok
}

func (f *MultipartForm) GetFiles(key string) ([]*multipart.FileHeader, bool) {
	if f == nil {
		return nil, false
	}

	fhs, ok := f.File[key]

	return fhs, ok
}
