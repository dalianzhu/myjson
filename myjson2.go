package myjson

import (
	"encoding/json"
	sysjson "encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var jsonit = jsoniter.ConfigCompatibleWithStandardLibrary

var IsDebug = false

func Debugf(format string, i ...interface{}) {
	if IsDebug {
		log.Printf(format+"\n", i...)
	}
}

type MyJson2 interface {
	Get(key string) MyJson2
	Set(key string, val interface{}) error
	Rm(key string)
	Index(i int) MyJson2
	Insert(i int, val interface{}) error
	Append(val interface{}) error

	GetValue() interface{}
	Len() int
	String() string
	Bytes() []byte
	Int() (int, error)
	Float64() (float64, error)
	Bool() (bool, error)
	Clone() MyJson2

	Keys() ([]string, error)
	Items() ([]interface{}, error)

	RangeSlice(f func(index int, val MyJson2) (bool, error)) error
	RangeMap(f func(key string, val MyJson2) (bool, error)) error

	IsErrOrNil() bool
	IsSlice() bool
	IsMap() bool
	IsNull() bool
}

func NewJsonFromBytes(bytesVal []byte) MyJson2 {
	val := &ValueJson{}
	err := val.UnmarshalJSON(bytesVal)
	if err != nil {
		return NewErrorJson(err)
	}
	return val
}

func NewJson(val interface{}) MyJson2 {
	switch v := val.(type) {
	case []byte:
		val := NewJsonFromBytes(v)
		return val
	case string:
		val := NewJsonFromBytes([]byte(v))
		return val
	case MyJson2:
		return v
	}

	refVal := reflect.ValueOf(val)
	switch refVal.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map:
		bytesVal, err := json.Marshal(val)
		if err != nil {
			// errstr := fmt.Sprintf("js解析失败：%v", err)
			return NewErrorJson(err)
		}
		return NewJsonFromBytes(bytesVal)

	case reflect.Ptr:
		if refVal.Elem().Kind() == reflect.Struct {
			bytesVal, err := json.Marshal(val)
			if err != nil {
				// errstr := fmt.Sprintf("js解析失败：%v", err)
				return NewErrorJson(err)
			}
			return NewJsonFromBytes(bytesVal)
		}
	}
	return NewErrorJson(fmt.Errorf("json val is invalid"))
}

func ToInt(intObj interface{}) (int, error) {
	// 假定int == int64，运行在64位机
	// Debugf("ToInt:%v", intObj)
	switch v := intObj.(type) {
	case sysjson.Number:
		strVal := string(v)
		return ToInt(strVal)
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		if v > math.MaxInt64 {
			info := fmt.Sprintf("ToInt error, overflowd %v", v)
			return 0, errors.New(info)
		}
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		strv := v
		if strings.Contains(v, ".") {
			strvArr := strings.Split(v, ".")
			if len(strvArr) != 2 {
				return 0, fmt.Errorf("ToInt error, invalid format:%v", v)
			}
			strv = strvArr[0]
		}
		if strv == "" {
			return 0, fmt.Errorf("ToInt error, empty string cannot convert to int")
		}
		if intv, err := strconv.Atoi(strv); err == nil {
			return intv, nil
		}
	case MyJson2:
		return v.Int()
	}
	return 0, fmt.Errorf("ToInt error, %v cannot convert to int", intObj)
}

func ToFloat64(item interface{}) (float64, error) {
	switch v := item.(type) {
	case sysjson.Number:
		return v.Float64()
	case int, int8, int16, int64, uint, uint8, uint16, uint32, uint64:
		intVal, err := ToInt(item)
		return float64(intVal), err
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case string:
		if floatNum, err := strconv.ParseFloat(v, 64); err == nil {
			return floatNum, nil
		}
	case MyJson2:
		return v.Float64()
	}
	return 0, fmt.Errorf("%v cannot convert to float", item)
}

func ToBool(item interface{}) (bool, error) {
	switch v := item.(type) {
	case bool:
		return v, nil
	case MyJson2:
		return v.Bool()
	default:
		boolValue, err := strconv.ParseBool(ToStr(item))
		if err != nil {
			return false, fmt.Errorf("%v cannot convert to bool", item)
		}
		return boolValue, nil
	}
}
