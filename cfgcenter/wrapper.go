package cfgcenter

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ini "github.com/hydah/go-ini-v1"

	"github.com/hydah/golib/cfgcenter/cfgstruct"
	"github.com/hydah/golib/logger"
)

var (
	ErrNotInited      = errors.New("DefaultCfgCenter not inited.")
	ErrEtcdAddrsNil   = errors.New("Etcd addrs is null.")
	ErrServiceKeyNil  = errors.New("Etcd service key is null.")
	ErrFileNameKeyNil = errors.New("Etcd filename key is null.")
	ErrSecNameNil     = errors.New("Sec name is null.")
	ErrKeyNameNil     = errors.New("Key name is null.")
	ErrUnknownMethod  = errors.New("Unknown method.")
)

const (
	DEFAULT_FILENAME_IN_ETCD = "default"
)

func LoadConfig(filename string, cfg interface{}) (cfgIni *ini.File, err error) {

	cfgIni, err = loadFromFile(filename, cfg)
	if err != nil {
		logger.Error(err)
		return
	}
	// check if cfg parsed valid, and log error
	compare(cfgIni, cfg)
	return
}

func getBaseCfg(filename string) (baseCfg *cfgstruct.BaseCfgSt, err error) {
	baseCfg = new(cfgstruct.BaseCfgSt)
	_, err = loadFromFile(filename, baseCfg)
	if err != nil {
		return
	}
	return
}

func loadFromFile(filename string, cfg interface{}) (cfgIni *ini.File, err error) {
	cfgIni, err = ini.Load(filename)
	if err != nil {
		return
	}
	err = cfgIni.StrictMapTo(cfg)
	if err != nil {
		return
	}
	return
}

// 还不是那么可靠，只打error作为参考，暂未返回error
func compare(cfgIni *ini.File, cfgSt interface{}) (err error) {
	iniMap := make(map[string]map[string]string)
	for _, sec := range cfgIni.Sections() {
		_, ok := iniMap[sec.Name()]
		if !ok {
			iniMap[sec.Name()] = make(map[string]string)
		}
		for _, kv := range sec.Keys() {
			iniMap[sec.Name()][kv.Name()] = kv.Value()
		}
	}

	typ := reflect.TypeOf(cfgSt)
	val := reflect.ValueOf(cfgSt)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	rangeField(val, "", iniMap)
	for secName, sec := range iniMap {
		for keyName, keyValue := range sec {
			if keyValue == "" {
				continue
			}
			logger.Error("sec %s, key %s, ini value %s, not contain in struct", secName, keyName, keyValue)
		}

	}
	return
}

// modified from gopkg.in/ini.v1 struct.go reflectFrom
func rangeField(val reflect.Value, secName string, iniMap map[string]map[string]string) (err error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		tpField := typ.Field(i)

		tag := tpField.Tag.Get("ini")
		if tag == "-" {
			continue
		}

		opts := strings.SplitN(tag, ",", 2)
		var fieldName, rawName string
		if len(opts[0]) > 0 {
			fieldName = opts[0]
			rawName = opts[0]
		} else {
			fieldName = tpField.Name
		}
		if len(fieldName) == 0 || !field.CanSet() {
			continue
		}

		// 只支持结构体指针
		if tpField.Type.Kind() == reflect.Ptr && tpField.Type.Elem().Kind() != reflect.Struct {
			logger.Error("unsupported type: %s", tpField.Type)
			continue
		}

		isNormalStruct := tpField.Type.Kind() == reflect.Struct && !tpField.Anonymous
		isNormalStructPtr := tpField.Type.Kind() == reflect.Ptr && !tpField.Anonymous && tpField.Type.Elem().Kind() == reflect.Struct
		if isNormalStruct || isNormalStructPtr {
			// 结构体成员或结构体指针成员，独立出一个section
			err = rangeField(field, fieldName, iniMap)
			if err != nil {
				// return
			}
			continue
		} else if tpField.Anonymous && rawName != "" {
			// 继承，如果有ini tag，按照tag独立出一个section
			err = rangeField(field, rawName, iniMap)
			if err != nil {
				// return
			}
			continue
		} else if tpField.Anonymous && rawName == "" {
			// 继承，如果没有ini tag，展开到上一层section
			err = rangeField(field, secName, iniMap)
			if err != nil {
				// return
			}
			continue
		}

		delim := tpField.Tag.Get("delim")
		if len(delim) == 0 {
			delim = ","
		}

		keyName := fieldName
		keyValue, err := marshalWithProperType(tpField.Type, field, delim)
		if err != nil {
			// return fmt.Errorf("error marshaling field (%s), sec name (%s), error : %v", fieldName, secName, err)
			logger.Error("error marshaling field (%s), sec name (%s), error : %v", fieldName, secName, err)
			continue
		}

		var iniValue string
		_, ok := iniMap[secName]
		if ok {
			iniValue, ok = iniMap[secName][keyName]
			if ok && iniValue == keyValue {
				delete(iniMap[secName], keyName)
				continue
			}
		}
		if isEmptyValue(field) {
			continue
		}
		// return fmt.Errorf("sec %s, key %s, struct value %s, ini value %s, not equal", secName, keyName, keyValue, iniValue)
		logger.Error("sec %s, key %s, struct value %s, ini value %s, not equal", secName, keyName, keyValue, iniValue)
	}
	return nil
}

// modified from gopkg.in/ini.v1 struct.go reflectWithProperType
func marshalWithProperType(t reflect.Type, field reflect.Value, delim string) (value string, err error) {
	switch t.Kind() {
	case reflect.String:
		value = field.String()
	case reflect.Bool:
		value = fmt.Sprint(field.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = fmt.Sprint(field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = fmt.Sprint(field.Uint())
	case reflect.Float32, reflect.Float64:
		value = fmt.Sprint(field.Float())
	case reflect.Slice:
		return marshalSliceWithProperType(field, delim)
	default:
		err = fmt.Errorf("unsupported type '%s'", t)
	}
	return
}

// modified from gopkg.in/ini.v1 struct.go reflectSliceWithProperType
func marshalSliceWithProperType(field reflect.Value, delim string) (value string, err error) {
	slice := field.Slice(0, field.Len())
	if field.Len() == 0 {
		return
	}

	var buf bytes.Buffer
	sliceOf := field.Type().Elem().Kind()
	for i := 0; i < field.Len(); i++ {
		switch sliceOf {
		case reflect.String:
			buf.WriteString(slice.Index(i).String())
		case reflect.Int, reflect.Int64:
			buf.WriteString(fmt.Sprint(slice.Index(i).Int()))
		case reflect.Uint, reflect.Uint64:
			buf.WriteString(fmt.Sprint(slice.Index(i).Uint()))
		case reflect.Float64:
			buf.WriteString(fmt.Sprint(slice.Index(i).Float()))
		default:
			err = fmt.Errorf("unsupported type '[]%s'", sliceOf)
			return
		}
		buf.WriteString(delim)
	}
	value = buf.String()[:buf.Len()-1]
	return
}

// modified from gopkg.in/ini.v1 struct.go isEmptyValue
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}