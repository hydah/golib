package xcjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type JSON struct {
	data interface{}
}

func NewJSON(body []byte) (*JSON, error) {
	json := new(JSON)
	err := json.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return json, nil
}

func New(data interface{}) *JSON {
	return &JSON{data}
}

func (j *JSON) Encode() ([]byte, error) {
	return j.MarshalJSON()
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(&j.data)
}

func (j *JSON) UnmarshalJSON(body []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.UseNumber()
	return dec.Decode(&j.data)
}

func (j *JSON) Map() (map[string]interface{}, error) {
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

func (j *JSON) MustMap() map[string]interface{} {
	return j.data.(map[string]interface{})
}

func (j *JSON) MustPhpMap() map[string]interface{} {
	m, err := j.Map()
	if err != nil {
		arr, err := j.Slice()
		if err != nil {
			m = map[string]interface{}{}
			return m
		}
		m = map[string]interface{}{}
		for idx, s := range arr {
			m[fmt.Sprintf("%d", idx)] = s
		}
		return m
	}
	return m
}

func (j *JSON) Slice() ([]interface{}, error) {
	if s, ok := (j.data).([]interface{}); ok {
		return s, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

func (j *JSON) MustSlice() []interface{} {
	return j.data.([]interface{})
}

func (j *JSON) Bool() (bool, error) {
	if b, ok := (j.data).(bool); ok {
		return b, nil
	}
	return false, errors.New("type assertion to bool failed")
}

func (j *JSON) MustBool() bool {
	return j.data.(bool)
}

func (j *JSON) MustPhpBool() bool {
	b, ok := j.data.(bool)
	if !ok {
		str := fmt.Sprintf("%v", j.data)
		if str == "true" || str == "True" {
			return true
		}
		if str == "false" || str == "False" {
			return false
		}
	}
	return b
}

func (j *JSON) String() (string, error) {
	if s, ok := (j.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

func (j *JSON) MustString() string {
	return j.data.(string)
}

func (j *JSON) MustPhpString() string {
	switch val := j.data.(type) {
	case string:
		return val
	case int, int64, uint64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%f", val)
	case json.Number:
		return val.String()
	}
	if j.data == nil {
		return ""
	}
	return fmt.Sprintf("%v", j.data)
}

func (j *JSON) MustPhpInt() int64 {
	//t := reflect.TypeOf(j.data)
	//logger.Trace("%v", t)

	switch val := j.data.(type) {
	case float64:
		return int64(val)
	case string:
		id, _ := strconv.ParseInt(val, 10, 64)
		return id
	case int64:
		return val
	case int:
		return int64(val)
	case uint64:
		return int64(j.data.(uint64))
	case json.Number:
		id, _ := val.Int64()
		return id
	}
	return 0
}

func (j *JSON) MustPhpArray() (arr []interface{}) {
	var ok bool
	if arr, ok = j.data.([]interface{}); ok {
		return
	}
	mp := j.MustPhpMap()
	arr = make([]interface{}, 0)
	for _, v := range mp {
		arr = append(arr, v)
	}
	return
}

func (j *JSON) MustPhpStringArray() (sarr []string) {
	arr := j.MustPhpArray()
	sarr = make([]string, 0)
	for _, v := range arr {
		if s, ok := v.(string); ok {
			sarr = append(sarr, s)
		}
	}
	return
}

func (j *JSON) Float64() (float64, error) {
	switch val := j.data.(type) {
	case json.Number:
		return val.Float64()
	default:
		if f, ok := (j.data).(float64); ok {
			return f, nil
		}
	}
	return -1, errors.New("type assertion to float64 failed")
}

func (j *JSON) MustFloat() float64 {
	switch val := j.data.(type) {
	case json.Number:
		if f, err := val.Float64(); err == nil {
			return f
		}
	}
	return j.data.(float64)
}

func (j *JSON) Int() (int, error) {
	switch val := j.data.(type) {
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return int(i), nil
		}
	}
	if f, ok := (j.data).(float64); ok {
		return int(f), nil
	}
	return -1, errors.New("type assertion to int failed")
}

func (j *JSON) Int64() (int64, error) {
	switch val := j.data.(type) {
	case json.Number:
		return val.Int64()
	}
	if f, ok := (j.data).(float64); ok {
		return int64(f), nil
	}
	return -1, errors.New("type assertion to int64 failed")
}

func (j *JSON) MustInt64() int64 {
	switch val := j.data.(type) {
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return i
		}
	}
	num := j.data.(float64)
	return int64(num)
}

func (j *JSON) Bytes() ([]byte, error) {
	if s, ok := (j.data).(string); ok {
		return []byte(s), nil
	}
	return nil, errors.New("type assertion to []byte failed")
}

func (j *JSON) StringSlice() ([]string, error) {
	slice, err := j.Slice()
	if err != nil {
		return nil, err
	}
	var stringSlice = make([]string, 0)
	for _, v := range slice {
		ss, ok := v.(string)
		if !ok {
			return nil, errors.New("type assertion to []string failed")
		}
		stringSlice = append(stringSlice, ss)
	}
	return stringSlice, nil
}

func (j *JSON) Set(key string, val interface{}) (err error) {
	m, err := j.Map()
	if err != nil {
		return
	}
	m[key] = val
	return
}

func (j *JSON) Get(key string) *JSON {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &JSON{val}
		}
		return &JSON{nil}
	}
	return &JSON{nil}
}

func (j *JSON) Del(key string) (err error) {
	m, err := j.Map()
	if err != nil {
		return
	}
	delete(m, key)
	return
}

func (j *JSON) GetFromMap(key string, def interface{}) (*JSON, error) {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &JSON{val}, nil
		}
	}
	err = fmt.Errorf("key: %s not found in map", key)
	return &JSON{def}, err
}

func (j *JSON) GetIndex(index int) *JSON {
	s, err := j.Slice()
	if err == nil {
		if len(s) > index {
			return &JSON{s[index]}
		}
		return &JSON{nil}
	}
	return &JSON{nil}
}

func (j *JSON) CheckGet(key string) (*JSON, bool) {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &JSON{val}, true
		}
	}
	return nil, false
}

func (j *JSON) GetLot(key ...string) *JSON {
	jin := j
	for i := range key {
		m, err := jin.Map()
		if err != nil {
			return &JSON{nil}
		}
		if val, ok := m[key[i]]; ok {
			jin = &JSON{val}
		} else {
			return &JSON{nil}
		}
	}
	return jin
}

func (j *JSON) Data() interface{} {
	return j.data
}

func (j *JSON) SetData(data interface{}) {
	j.data = data
}
