package jester

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"
)

var ErrTypeMismatch = errors.New("jester: type assertion failed (type mismatch)")

type Data struct {
	data any
}

// MarshalJSON implements the json.Marshaler interface.
func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Data) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(&d.data)
}

// New creates a new Data instance with an empty map.
func New(data any) *Data {
	return &Data{data: data}
}

// NewEmpty creates a new Data instance with an empty map.
func NewEmpty() *Data {
	return New(make(map[string]any))
}

// NewJson creates a new Data instance from JSON data.
func NewJson(data []byte) (d *Data, err error) {
	d = &Data{}
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	err = dec.Decode(&d.data)
	return d, err
}

// NewReader creates a new Data instance from an io.Reader.
func NewReader(r io.Reader) (d *Data, err error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	d = NewEmpty()
	err = dec.Decode(&d.data)
	return d, err
}

// Interface returns the underlying data.
func (d *Data) Interface() any {
	return d.data
}

// Set modifies the data structure by setting the value for the specified key.
func (d *Data) Set(key string, val any) {
	m, err := d.Map()
	if err != nil {
		return
	}
	m[key] = val
}

// SetPath modifies the data structure by setting the value for the specified path.
func (d *Data) SetPath(branch []any, val any) {
	if len(branch) == 0 {
		d.data = val
		return
	}

	current := d

	for i := range len(branch) - 1 {
		key := branch[i]

		next := current.get(key)
		if next.data == nil {
			switch k := key.(type) {
			case string:
				current.Set(k, make(map[string]any))
			case int:
				// Need to create a slice large enough to hold the index
				slice, err := current.Slice()
				if err != nil || int(k) >= len(slice) {
					current.Set(strconv.Itoa(k), make(map[string]any))
				} else {
					// Ensure slice has capacity
					index := k
					for len(slice) <= index {
						slice = append(slice, nil)
					}
					slice[index] = make(map[string]any)
					current.data = slice
				}
			}
			current = current.get(key)
			continue
		}

		// If next data is not a map or slice, convert it to a map
		if _, errMap := next.Map(); errMap != nil {
			if _, errSlice := next.Slice(); errSlice != nil {
				// Force convert primitive to map
				switch k := key.(type) {
				case string:
					current.Set(k, make(map[string]any))
				case int:
					current.Set(strconv.Itoa(k), make(map[string]any))
				}
			}
		}

		current = current.get(key)
	}

	lastKey := branch[len(branch)-1]

	if k, ok := lastKey.(string); ok {
		current.Set(k, val)
		return
	}

	k, ok := lastKey.(int)
	if !ok {
		return
	}

	slice, err := current.Slice()
	if err != nil {
		// Not a slice, convert to map and use string key
		current.Set(strconv.Itoa(k), val)
		return
	}

	for len(slice) <= k {
		slice = append(slice, nil)
	}
	slice[k] = val
	current.data = slice
}

// Delete deletes a key from the data structure.
func (d *Data) Delete(key string) {
	m, err := d.Map()
	if err != nil {
		return
	}
	delete(m, key)
}

// Get retrieves a value from the data structure at the specified path.
func (d *Data) Get(keys ...any) *Data {
	data := d

	for _, key := range keys {
		data = data.get(key)

		// If the data is nil, return it immediately.
		if data.data == nil {
			return data
		}
	}

	return data
}

func (d *Data) get(key any) *Data {
	if d.data == nil {
		return New(nil)
	}

	// Try as map with string key
	if dataMap, ok := d.data.(map[string]any); ok {
		if keyStr, ok := key.(string); ok {
			if v, ok := dataMap[keyStr]; ok {
				return New(v)
			}
			return New(nil)
		}
		// Try to convert int key to string for maps
		if keyInt, ok := key.(int); ok {
			keyStr := strconv.Itoa(keyInt)
			if v, ok := dataMap[keyStr]; ok {
				return New(v)
			}
		}
	}

	// Try as slice with int key
	if dataSlice, ok := d.data.([]any); ok {
		if keyInt, ok := key.(int); ok {
			if keyInt >= 0 && keyInt < len(dataSlice) {
				return New(dataSlice[keyInt])
			}
		}
	}

	return New(nil)
}

// Map returns the underlying data as a map[string]any.
func (d *Data) Map() (map[string]any, error) {
	if m, ok := d.data.(map[string]any); ok {
		return m, nil
	}
	return nil, ErrTypeMismatch
}

// MustMap returns the underlying data as a map[string]any with optional default value.
func (d *Data) MustMap(args ...map[string]any) map[string]any {
	var value map[string]any

	if m, ok := d.data.(map[string]any); ok {
		value = m
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// Slice returns the underlying data as a []any.
func (d *Data) Slice() ([]any, error) {
	if s, ok := d.data.([]any); ok {
		return s, nil
	}
	return nil, ErrTypeMismatch
}

// MustSlice returns the underlying data as a []any with optional default value.
func (d *Data) MustSlice(args ...[]any) []any {
	var value []any

	if s, ok := d.data.([]any); ok {
		value = s
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// Bool returns the underlying data as a bool.
func (d *Data) Bool() (bool, error) {
	if b, ok := d.data.(bool); ok {
		return b, nil
	}
	return false, ErrTypeMismatch
}

// MustBool returns the underlying data as a bool with optional default value.
func (d *Data) MustBool(args ...bool) bool {
	var value bool

	if b, ok := d.data.(bool); ok {
		value = b
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// String returns the underlying data as a string.
func (d *Data) String() (string, error) {
	if s, ok := d.data.(string); ok {
		return s, nil
	}
	return "", ErrTypeMismatch
}

// MustString returns the underlying data as a string with optional default value.
func (d *Data) MustString(args ...string) string {
	var value string

	if s, ok := d.data.(string); ok {
		value = s
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// Bytes returns the underlying data as a []byte.
func (d *Data) Bytes() ([]byte, error) {
	if b, ok := d.data.([]byte); ok {
		return b, nil
	}
	return nil, ErrTypeMismatch
}

// MustBytes returns the underlying data as a []byte with optional default value.
func (d *Data) MustBytes(args ...[]byte) []byte {
	var value []byte

	if b, ok := d.data.([]byte); ok {
		value = b
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// StringSlice returns the underlying data as a []string.
func (d *Data) StringSlice() ([]string, error) {
	s, err := d.Slice()
	if err != nil {
		return nil, err
	}

	strs := make([]string, 0, len(s))

	for _, v := range s {
		if v == nil {
			strs = append(strs, "")
			continue
		}

		str, ok := v.(string)
		if !ok {
			return nil, ErrTypeMismatch
		}

		strs = append(strs, str)
	}

	return strs, nil
}

// MustStringSlice returns the underlying data as a []string with optional default value.
func (d *Data) MustStringSlice(args ...[]string) []string {
	var value []string

	if s, err := d.StringSlice(); err == nil {
		value = s
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// Int returns the underlying data as an int.
func (d *Data) Int() (int, error) {
	switch v := d.data.(type) {
	case json.Number:
		i, err := v.Int64()
		return int(i), err
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(d.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(d.data).Uint()), nil
	case float32, float64:
		return int(reflect.ValueOf(d.data).Float()), nil
	default:
		return 0, ErrTypeMismatch
	}
}

// MustInt returns the underlying data as an int with optional default value.
func (d *Data) MustInt(args ...int) int {
	var value int

	if i, err := d.Int(); err == nil {
		value = i
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// Int64 returns the underlying data as an int64.
func (d *Data) Int64() (int64, error) {
	switch v := d.data.(type) {
	case json.Number:
		return v.Int64()
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(d.data).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(d.data).Uint()), nil
	case float32, float64:
		return int64(reflect.ValueOf(d.data).Float()), nil
	default:
		return 0, ErrTypeMismatch
	}
}

// MustInt64 returns the underlying data as an int64 with optional default value.
func (d *Data) MustInt64(args ...int64) int64 {
	var value int64

	if i, err := d.Int64(); err == nil {
		value = i
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// Uint64 returns the underlying data as a uint64.
func (d *Data) Uint64() (uint64, error) {
	switch v := d.data.(type) {
	case json.Number:
		return strconv.ParseUint(v.String(), 10, 64)
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(d.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(d.data).Uint(), nil
	case float32, float64:
		return uint64(reflect.ValueOf(d.data).Float()), nil
	default:
		return 0, ErrTypeMismatch
	}
}

// MustUint64 returns the underlying data as a uint64 with optional default value.
func (d *Data) MustUint64(args ...uint64) uint64 {
	var value uint64

	if i, err := d.Uint64(); err == nil {
		value = i
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}

// Float64 returns the underlying data as a float64.
func (d *Data) Float64() (float64, error) {
	switch v := d.data.(type) {
	case json.Number:
		return v.Float64()
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(d.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(d.data).Uint()), nil
	case float32, float64:
		return reflect.ValueOf(d.data).Float(), nil
	default:
		return 0, ErrTypeMismatch
	}
}

// MustFloat64 returns the underlying data as a float64 with optional default value.
func (d *Data) MustFloat64(args ...float64) float64 {
	var value float64

	if f, err := d.Float64(); err == nil {
		value = f
	} else if len(args) > 0 {
		value = args[0]
	}

	return value
}
