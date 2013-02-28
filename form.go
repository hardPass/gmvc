package gmvc

import (
	"mime/multipart"
	"strconv"
)

type Values map[string][]string

func (v Values) Get(key string) string {
	if v == nil {
		return ""
	}
	ss, ok := v[key]
	if !ok || len(ss) == 0 {
		return ""
	}
	return ss[0]
}

func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

func (v Values) Del(key string) {
	delete(v, key)
}

func (v Values) get(key string) (string, bool) {
	if v == nil {
		return "", false
	}

	ss, ok := v[key]
	if !ok || len(ss) == 0 {
		return "", ok
	}

	return ss[0], ok
}

// integer

func (v Values) getInt(key string, def int64, bitsize int) (int64, bool, error) {
	s, ok := v.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	n, err := strconv.ParseInt(s, 10, bitsize)
	if err != nil {
		return def, ok, err
	}

	return n, ok, nil
}

func (v Values) getUint(key string, def uint64, bitsize int) (uint64, bool, error) {
	s, ok := v.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	n, err := strconv.ParseUint(s, 10, bitsize)
	if err != nil {
		return def, ok, err
	}

	return n, ok, nil
}

//Int

func (v Values) Int(key string) int {
	n, _, _ := v.getInt(key, 0, 0)
	return int(n)
}

func (v Values) GetInt(key string, def int) (int, bool, error) {
	n, ok, err := v.getInt(key, int64(def), 0)
	return int(n), ok, err
}

func (v Values) GetInts(key string) ([]int, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]int, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = int(n)
	}

	return ns, ok, nil
}

// Int8

func (v Values) Int8(key string) int8 {
	n, _, _ := v.getInt(key, 0, 8)
	return int8(n)
}

func (v Values) GetInt8(key string, def int8) (int8, bool, error) {
	n, ok, err := v.getInt(key, int64(def), 8)
	return int8(n), ok, err
}

func (v Values) GetInt8s(key string) ([]int8, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]int8, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = int8(n)
	}

	return ns, ok, nil
}

// Int16

func (v Values) Int16(key string) int16 {
	n, _, _ := v.getInt(key, 0, 16)
	return int16(n)
}

func (v Values) GetInt16(key string, def int16) (int16, bool, error) {
	n, ok, err := v.getInt(key, int64(def), 16)
	return int16(n), ok, err
}

func (v Values) GetInt16s(key string) ([]int16, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]int16, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = int16(n)
	}

	return ns, ok, nil
}

// Int32

func (v Values) Int32(key string) int32 {
	n, _, _ := v.getInt(key, 0, 32)
	return int32(n)
}

func (v Values) GetInt32(key string, def int32) (int32, bool, error) {
	n, ok, err := v.getInt(key, int64(def), 32)
	return int32(n), ok, err
}

func (v Values) GetInt32s(key string) ([]int32, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]int32, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = int32(n)
	}

	return ns, ok, nil
}

// Int64

func (v Values) Int64(key string) int64 {
	n, _, _ := v.getInt(key, 0, 64)
	return n
}

func (v Values) GetInt64(key string, def int64) (int64, bool, error) {
	n, ok, err := v.getInt(key, def, 64)
	return n, ok, err
}

func (v Values) GetInt64s(key string) ([]int64, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]int64, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = int64(n)
	}

	return ns, ok, nil
}

// uint

func (v Values) Uint(key string) uint {
	n, _, _ := v.getUint(key, 0, 0)
	return uint(n)
}

func (v Values) GetUint(key string, def uint) (uint, bool, error) {
	n, ok, err := v.getUint(key, uint64(def), 0)
	return uint(n), ok, err
}

func (v Values) GetUints(key string) ([]uint, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]uint, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseUint(s, 10, 0)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = uint(n)
	}

	return ns, ok, nil
}

// Uint8

func (v Values) Uint8(key string) uint8 {
	n, _, _ := v.getUint(key, 0, 8)
	return uint8(n)
}

func (v Values) GetUint8(key string, def uint8) (uint8, bool, error) {
	n, ok, err := v.getUint(key, uint64(def), 8)
	return uint8(n), ok, err
}

func (v Values) GetUint8s(key string) ([]uint16, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]uint16, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = uint16(n)
	}

	return ns, ok, nil
}

// Uint16

func (v Values) Uint16(key string) uint16 {
	n, _, _ := v.getUint(key, 0, 16)
	return uint16(n)
}

func (v Values) GetUint16(key string, def uint16) (uint16, bool, error) {
	n, ok, err := v.getUint(key, uint64(def), 16)
	return uint16(n), ok, err
}

func (v Values) GetUint16s(key string) ([]uint16, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]uint16, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = uint16(n)
	}

	return ns, ok, nil
}

// Uint32

func (v Values) Uint32(key string) uint32 {
	n, _, _ := v.getUint(key, 0, 32)
	return uint32(n)
}

func (v Values) GetUint32(key string, def uint32) (uint32, bool, error) {
	n, ok, err := v.getUint(key, uint64(def), 32)
	return uint32(n), ok, err
}

func (v Values) GetUint32s(key string) ([]uint32, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]uint32, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = uint32(n)
	}

	return ns, ok, nil
}

// Uint64

func (v Values) Uint64(key string) uint64 {
	n, _, _ := v.getUint(key, 0, 64)
	return n
}

func (v Values) GetUint64(key string, def uint64) (uint64, bool, error) {
	n, ok, err := v.getUint(key, def, 64)
	return n, ok, err
}

func (v Values) GetUint64s(key string) ([]uint64, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	ns := make([]uint64, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, ok, err
		}
		ns[i] = uint64(n)
	}

	return ns, ok, nil
}

// float

func (v Values) getFloat(key string, def float64, bitsize int) (float64, bool, error) {
	s, ok := v.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	f, err := strconv.ParseFloat(s, bitsize)
	if err != nil {
		return def, ok, err
	}

	return f, ok, nil
}

// Float32

func (v Values) Float32(key string) float32 {
	f, _, _ := v.getFloat(key, 0, 32)
	return float32(f)
}

func (v Values) GetFloat32(key string, def float32) (float32, bool, error) {
	f, ok, err := v.getFloat(key, float64(def), 32)
	return float32(f), ok, err
}

func (v Values) GetFloat32s(key string) ([]float32, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	fs := make([]float32, len(ss))
	for i, s := range ss {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, ok, err
		}
		fs[i] = float32(f)
	}

	return fs, ok, nil
}

// Float64

func (v Values) Float64(key string) float64 {
	f, _, _ := v.getFloat(key, 0, 64)
	return f
}

func (v Values) GetFloat64(key string, def float64) (float64, bool, error) {
	f, ok, err := v.getFloat(key, def, 64)
	return f, ok, err
}

func (v Values) GetFloat64s(key string) ([]float64, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	fs := make([]float64, len(ss))
	for i, s := range ss {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, ok, err
		}
		fs[i] = f
	}

	return fs, ok, nil
}

// Bool

func (v Values) Bool(key string) bool {
	b, _, _ := v.GetBool(key, false)
	return b
}

func (v Values) GetBool(key string, def bool) (bool, bool, error) {
	s, ok := v.get(key)
	if !ok || s == "" {
		return def, ok, nil
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return def, ok, err
	}

	return b, ok, nil
}

func (v Values) GetBools(key string) ([]bool, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]
	if !ok {
		return nil, ok, nil
	}

	bs := make([]bool, len(ss))
	for i, s := range ss {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, ok, err
		}
		bs[i] = b
	}

	return bs, ok, nil
}

// string

func (v Values) String(key string) string {
	return v.Get(key)
}

func (v Values) GetString(key string, def string) (string, bool, error) {
	s, ok := v.get(key)
	if !ok {
		return def, ok, nil
	}

	return s, ok, nil
}

func (v Values) GetStrings(key string) ([]string, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	ss, ok := v[key]

	return ss, ok, nil
}

type MultipartForm struct {
	Values
	Files map[string][]*multipart.FileHeader
}

// file

func (f *MultipartForm) File(key string) *multipart.FileHeader {
	fh, _ := f.GetFile(key)
	return fh
}

func (f *MultipartForm) GetFile(key string) (*multipart.FileHeader, bool) {
	if f == nil {
		return nil, false
	}

	fhs, ok := f.Files[key]
	if !ok || len(fhs) == 0 {
		return nil, ok
	}
	return fhs[0], ok
}

func (f *MultipartForm) GetFiles(key string) ([]*multipart.FileHeader, bool) {
	if f == nil {
		return nil, false
	}

	fhs, ok := f.Files[key]
	return fhs, ok
}
