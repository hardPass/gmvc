package gmvc

import (
	"strconv"
)

var (
	values *formvalues = &formvalues{}
)

type formvalues struct{}

func (fv *formvalues) get(f map[string][]string, key string) (string, bool) {
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

func (fv *formvalues) getInt(f map[string][]string, key string, def int64, bitsize int) (int64, bool, error) {
	s, ok := fv.get(f, key)
	if !ok || s == "" {
		return def, ok, nil
	}

	v, err := strconv.ParseInt(s, 10, bitsize)
	if err != nil {
		return def, ok, err
	}

	return v, ok, nil
}

func (fv *formvalues) getUint(f map[string][]string, key string, def uint64, bitsize int) (uint64, bool, error) {
	s, ok := fv.get(f, key)
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

func (fv *formvalues) Int(f map[string][]string, key string) int {
	v, _, _ := fv.getInt(f, key, 0, 0)
	return int(v)
}

func (fv *formvalues) GetInt(f map[string][]string, key string, def int) (int, bool, error) {
	v, ok, err := fv.getInt(f, key, int64(def), 0)
	return int(v), ok, err
}

func (fv *formvalues) GetInts(f map[string][]string, key string) ([]int, bool, error) {
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

func (fv *formvalues) Int8(f map[string][]string, key string) int8 {
	v, _, _ := fv.getInt(f, key, 0, 8)
	return int8(v)
}

func (fv *formvalues) GetInt8(f map[string][]string, key string, def int8) (int8, bool, error) {
	v, ok, err := fv.getInt(f, key, int64(def), 8)
	return int8(v), ok, err
}

func (fv *formvalues) GetInt8s(f map[string][]string, key string) ([]int8, bool, error) {
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

func (fv *formvalues) Int16(f map[string][]string, key string) int16 {
	v, _, _ := fv.getInt(f, key, 0, 16)
	return int16(v)
}

func (fv *formvalues) GetInt16(f map[string][]string, key string, def int16) (int16, bool, error) {
	v, ok, err := fv.getInt(f, key, int64(def), 16)
	return int16(v), ok, err
}

func (fv *formvalues) GetInt16s(f map[string][]string, key string) ([]int16, bool, error) {
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

func (fv *formvalues) Int32(f map[string][]string, key string) int32 {
	v, _, _ := fv.getInt(f, key, 0, 32)
	return int32(v)
}

func (fv *formvalues) GetInt32(f map[string][]string, key string, def int32) (int32, bool, error) {
	v, ok, err := fv.getInt(f, key, int64(def), 32)
	return int32(v), ok, err
}

func (fv *formvalues) GetInt32s(f map[string][]string, key string) ([]int32, bool, error) {
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

func (fv *formvalues) Int64(f map[string][]string, key string) int64 {
	v, _, _ := fv.getInt(f, key, 0, 64)
	return v
}

func (fv *formvalues) GetInt64(f map[string][]string, key string, def int64) (int64, bool, error) {
	v, ok, err := fv.getInt(f, key, def, 64)
	return v, ok, err
}

func (fv *formvalues) GetInt64s(f map[string][]string, key string) ([]int64, bool, error) {
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

func (fv *formvalues) Uint(f map[string][]string, key string) uint {
	v, _, _ := fv.getUint(f, key, 0, 0)
	return uint(v)
}

func (fv *formvalues) GetUint(f map[string][]string, key string, def uint) (uint, bool, error) {
	v, ok, err := fv.getUint(f, key, uint64(def), 0)
	return uint(v), ok, err
}

func (fv *formvalues) GetUints(f map[string][]string, key string) ([]uint, bool, error) {
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

func (fv *formvalues) Uint8(f map[string][]string, key string) uint8 {
	v, _, _ := fv.getUint(f, key, 0, 8)
	return uint8(v)
}

func (fv *formvalues) GetUint8(f map[string][]string, key string, def uint8) (uint8, bool, error) {
	v, ok, err := fv.getUint(f, key, uint64(def), 8)
	return uint8(v), ok, err
}

func (fv *formvalues) GetUint8s(f map[string][]string, key string) ([]uint16, bool, error) {
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

func (fv *formvalues) Uint16(f map[string][]string, key string) uint16 {
	v, _, _ := fv.getUint(f, key, 0, 16)
	return uint16(v)
}

func (fv *formvalues) GetUint16(f map[string][]string, key string, def uint16) (uint16, bool, error) {
	v, ok, err := fv.getUint(f, key, uint64(def), 16)
	return uint16(v), ok, err
}

func (fv *formvalues) GetUint16s(f map[string][]string, key string) ([]uint16, bool, error) {
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

func (fv *formvalues) Uint32(f map[string][]string, key string) uint32 {
	v, _, _ := fv.getUint(f, key, 0, 32)
	return uint32(v)
}

func (fv *formvalues) GetUint32(f map[string][]string, key string, def uint32) (uint32, bool, error) {
	v, ok, err := fv.getUint(f, key, uint64(def), 32)
	return uint32(v), ok, err
}

func (fv *formvalues) GetUint32s(f map[string][]string, key string) ([]uint32, bool, error) {
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

func (fv *formvalues) Uint64(f map[string][]string, key string) uint64 {
	v, _, _ := fv.getUint(f, key, 0, 64)
	return v
}

func (fv *formvalues) GetUint64(f map[string][]string, key string, def uint64) (uint64, bool, error) {
	v, ok, err := fv.getUint(f, key, def, 64)
	return v, ok, err
}

func (fv *formvalues) GetUint64s(f map[string][]string, key string) ([]uint64, bool, error) {
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

func (fv *formvalues) getFloat(f map[string][]string, key string, def float64, bitsize int) (float64, bool, error) {
	s, ok := fv.get(f, key)
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

func (fv *formvalues) Float32(f map[string][]string, key string) float32 {
	v, _, _ := fv.getFloat(f, key, 0, 32)
	return float32(v)
}

func (fv *formvalues) GetFloat32(f map[string][]string, key string, def float32) (float32, bool, error) {
	v, ok, err := fv.getFloat(f, key, float64(def), 32)
	return float32(v), ok, err
}

func (fv *formvalues) GetFloat32s(f map[string][]string, key string) ([]float32, bool, error) {
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

func (fv *formvalues) Float64(f map[string][]string, key string) float64 {
	v, _, _ := fv.getFloat(f, key, 0, 64)
	return v
}

func (fv *formvalues) GetFloat64(f map[string][]string, key string, def float64) (float64, bool, error) {
	v, ok, err := fv.getFloat(f, key, def, 64)
	return v, ok, err
}

func (fv *formvalues) GetFloat64s(f map[string][]string, key string) ([]float64, bool, error) {
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

func (fv *formvalues) Bool(f map[string][]string, key string) bool {
	v, _, _ := fv.GetBool(f, key, false)
	return v
}

func (fv *formvalues) GetBool(f map[string][]string, key string, def bool) (bool, bool, error) {
	s, ok := fv.get(f, key)
	if !ok || s == "" {
		return def, ok, nil
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return def, ok, err
	}

	return v, ok, nil
}

func (fv *formvalues) GetBools(f map[string][]string, key string) ([]bool, bool, error) {
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

func (fv *formvalues) String(f map[string][]string, key string) string {
	v, _, _ := fv.GetString(f, key, "")
	return v
}

func (fv *formvalues) GetString(f map[string][]string, key string, def string) (string, bool, error) {
	s, ok := fv.get(f, key)
	if !ok {
		return def, ok, nil
	}

	return s, ok, nil
}

func (fv *formvalues) GetStrings(f map[string][]string, key string) ([]string, bool, error) {
	if f == nil {
		return nil, false, nil
	}

	ss, ok := f[key]

	return ss, ok, nil
}
