package gmvc

import (
	"mime/multipart"
	"strconv"
)

type Form map[string][]string

func (f Form) get(key string) (string, bool) {
	if f == nil {
		return "", false
	}

	ss, ok := f[key]
	if !ok || len(ss) == 0 {
		return "", ok
	}

	return ss[0], ok
}

// integer

func (f Form) getInt(key string, def int64, bitsize int) (int64, bool, error) {
	s, ok := f.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	v, err := strconv.ParseInt(s, 10, bitsize)
	if err != nil {
		return def, ok, err
	}

	return v, ok, nil
}

func (f Form) getUint(key string, def uint64, bitsize int) (uint64, bool, error) {
	s, ok := f.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	v, err := strconv.ParseUint(s, 10, bitsize)
	if err != nil {
		return def, ok, err
	}

	return v, ok, nil
}

//Int

func (f Form) Int(key string) int {
	v, _, _ := f.getInt(key, 0, 0)
	return int(v)
}

func (f Form) GetInt(key string, def int) (int, bool, error) {
	v, ok, err := f.getInt(key, int64(def), 0)
	return int(v), ok, err
}

func (f Form) GetInts(key string) ([]int, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]int, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = int(v)
	}

	return vs, ok, nil
}

// Int8

func (f Form) Int8(key string) int8 {
	v, _, _ := f.getInt(key, 0, 8)
	return int8(v)
}

func (f Form) GetInt8(key string, def int8) (int8, bool, error) {
	v, ok, err := f.getInt(key, int64(def), 8)
	return int8(v), ok, err
}

func (f Form) GetInt8s(key string) ([]int8, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]int8, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = int8(v)
	}

	return vs, ok, nil
}

// Int16

func (f Form) Int16(key string) int16 {
	v, _, _ := f.getInt(key, 0, 16)
	return int16(v)
}

func (f Form) GetInt16(key string, def int16) (int16, bool, error) {
	v, ok, err := f.getInt(key, int64(def), 16)
	return int16(v), ok, err
}

func (f Form) GetInt16s(key string) ([]int16, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]int16, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = int16(v)
	}

	return vs, ok, nil
}

// Int32

func (f Form) Int32(key string) int32 {
	v, _, _ := f.getInt(key, 0, 32)
	return int32(v)
}

func (f Form) GetInt32(key string, def int32) (int32, bool, error) {
	v, ok, err := f.getInt(key, int64(def), 32)
	return int32(v), ok, err
}

func (f Form) GetInt32s(key string) ([]int32, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]int32, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = int32(v)
	}

	return vs, ok, nil
}

// Int64

func (f Form) Int64(key string) int64 {
	v, _, _ := f.getInt(key, 0, 64)
	return v
}

func (f Form) GetInt64(key string, def int64) (int64, bool, error) {
	v, ok, err := f.getInt(key, def, 64)
	return v, ok, err
}

func (f Form) GetInt64s(key string) ([]int64, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]int64, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = int64(v)
	}

	return vs, ok, nil
}

// uint

func (f Form) Uint(key string) uint {
	v, _, _ := f.getUint(key, 0, 0)
	return uint(v)
}

func (f Form) GetUint(key string, def uint) (uint, bool, error) {
	v, ok, err := f.getUint(key, uint64(def), 0)
	return uint(v), ok, err
}

func (f Form) GetUints(key string) ([]uint, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]uint, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseUint(s, 10, 0)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = uint(v)
	}

	return vs, ok, nil
}

// Uint8

func (f Form) Uint8(key string) uint8 {
	v, _, _ := f.getUint(key, 0, 8)
	return uint8(v)
}

func (f Form) GetUint8(key string, def uint8) (uint8, bool, error) {
	v, ok, err := f.getUint(key, uint64(def), 8)
	return uint8(v), ok, err
}

func (f Form) GetUint8s(key string) ([]uint16, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]uint16, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = uint16(v)
	}

	return vs, ok, nil
}

// Uint16

func (f Form) Uint16(key string) uint16 {
	v, _, _ := f.getUint(key, 0, 16)
	return uint16(v)
}

func (f Form) GetUint16(key string, def uint16) (uint16, bool, error) {
	v, ok, err := f.getUint(key, uint64(def), 16)
	return uint16(v), ok, err
}

func (f Form) GetUint16s(key string) ([]uint16, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]uint16, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = uint16(v)
	}

	return vs, ok, nil
}

// Uint32

func (f Form) Uint32(key string) uint32 {
	v, _, _ := f.getUint(key, 0, 32)
	return uint32(v)
}

func (f Form) GetUint32(key string, def uint32) (uint32, bool, error) {
	v, ok, err := f.getUint(key, uint64(def), 32)
	return uint32(v), ok, err
}

func (f Form) GetUint32s(key string) ([]uint32, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]uint32, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = uint32(v)
	}

	return vs, ok, nil
}

// Uint64

func (f Form) Uint64(key string) uint64 {
	v, _, _ := f.getUint(key, 0, 64)
	return v
}

func (f Form) GetUint64(key string, def uint64) (uint64, bool, error) {
	v, ok, err := f.getUint(key, def, 64)
	return v, ok, err
}

func (f Form) GetUint64s(key string) ([]uint64, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]uint64, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = uint64(v)
	}

	return vs, ok, nil
}

// float

func (f Form) getFloat(key string, def float64, bitsize int) (float64, bool, error) {
	s, ok := f.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	v, err := strconv.ParseFloat(s, bitsize)
	if err != nil {
		return def, ok, err
	}

	return v, ok, nil
}

// Float32

func (f Form) Float32(key string) float32 {
	v, _, _ := f.getFloat(key, 0, 32)
	return float32(v)
}

func (f Form) GetFloat32(key string, def float32) (float32, bool, error) {
	v, ok, err := f.getFloat(key, float64(def), 32)
	return float32(v), ok, err
}

func (f Form) GetFloat32s(key string) ([]float32, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]float32, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = float32(v)
	}

	return vs, ok, nil
}

// Float64

func (f Form) Float64(key string) float64 {
	v, _, _ := f.getFloat(key, 0, 64)
	return v
}

func (f Form) GetFloat64(key string, def float64) (float64, bool, error) {
	v, ok, err := f.getFloat(key, def, 64)
	return v, ok, err
}

func (f Form) GetFloat64s(key string) ([]float64, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]float64, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = v
	}

	return vs, ok, nil
}

// Bool

func (f Form) Bool(key string) bool {
	v, _, _ := f.GetBool(key, false)
	return v
}

func (f Form) GetBool(key string, def bool) (bool, bool, error) {
	s, ok := f.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return def, ok, err
	}

	return v, ok, nil
}

func (f Form) GetBools(key string) ([]bool, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]
	if !ok {
		return nil, ok, nil
	}

	vs := make([]bool, len(ss))
	for i, s := range ss {
		v, err := strconv.ParseBool(s)
		if err != nil {
			return nil, ok, err
		}
		vs[i] = v
	}

	return vs, ok, nil
}

// string

func (f Form) String(key string) string {
	v, _, _ := f.GetString(key, "")
	return v
}

func (f Form) GetString(key string, def string) (string, bool, error) {
	s, ok := f.get(key)
	if !ok {
		return def, ok, nil
	}

	return s, ok, nil
}

func (f Form) GetStrings(key string) ([]string, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]

	return ss, ok, nil
}

type MultipartForm struct {
	Form
	File map[string][]*multipart.FileHeader
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
