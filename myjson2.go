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

	"google.golang.org/protobuf/types/known/structpb"
)

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
	Insert(i int, val interface{}) (MyJson2, error)
	Append(val interface{}) (MyJson2, error)

	Len() int
	String() string
	Bytes() []byte
	Int() (int, error)
	Float64() (float64, error)
	Bool() (bool, error)
	Clone() MyJson2

	RangeSlice(f func(index int, val MyJson2) (bool, error)) error
	RangeMap(f func(key string, val MyJson2) (bool, error)) error

	IsErrOrNil() bool
	IsSlice() bool
	IsMap() bool

	PbValue() *structpb.Value
}

func NewJsonFromBytes(bytesVal []byte) MyJson2 {
	val := &JsonVal{}
	err := val.UnmarshalJSON(bytesVal)
	if err != nil {
		return &NilJson{}
	}
	v := &ValueJson{val}
	return v
}

func NewJson(val interface{}) MyJson2 {
	switch v := val.(type) {
	case string:
		val := NewJsonFromBytes([]byte(v))
		Debugf("NewJson1:%v", val)
		return val
	case []byte:
		val := NewJsonFromBytes(v)
		Debugf("NewJson2:%v", val)
		return val
	}

	refVal := reflect.ValueOf(val)
	switch refVal.Kind() {
	case reflect.Struct, reflect.Slice:
		bytesVal, err := json.Marshal(val)
		if err != nil {
			// errstr := fmt.Sprintf("js解析失败：%v", err)
			return &NilJson{}
		}
		return NewJsonFromBytes(bytesVal)

	case reflect.Ptr:
		if refVal.Elem().Kind() == reflect.Struct {
			bytesVal, err := json.Marshal(val)
			if err != nil {
				// errstr := fmt.Sprintf("js解析失败：%v", err)
				return &NilJson{}
			}
			return NewJsonFromBytes(bytesVal)
		}
	}
	return &NilJson{}
}

type ValueJson struct {
	data *JsonVal
}

func (v *ValueJson) PbValue() *structpb.Value {
	bytesVal, err := v.data.Kind.MarshalJSON()
	if err != nil {
		return nil
	}

	pbVal := new(structpb.Value)
	err = pbVal.UnmarshalJSON(bytesVal)
	if err != nil {
		return nil
	}
	return pbVal
}

func (v *ValueJson) Get(key string) MyJson2 {
	if v.IsErrOrNil() {
		return new(NilJson)
	}

	structVal, ok := v.data.Kind.(*MapJsonValKind)
	if !ok {
		return new(NilJson)
	}

	jsonValData, ok := structVal.val[key]
	if ok {
		return &ValueJson{jsonValData}
	}
	return new(NilJson)
}

func iterToPbValue(val interface{}) *structpb.Value {
	switch realVal := val.(type) {
	case MyJson2:
		return realVal.PbValue()
	default:
		tpVal, err := structpb.NewValue(val)
		if err != nil {
			return structpb.NewNullValue()
		}
		return tpVal
	}
	return structpb.NewNullValue()
}

func (v *ValueJson) Set(key string, val interface{}) error {
	if v.IsErrOrNil() {
		return fmt.Errorf("json is nil")
	}
	structVal, ok := v.data.Kind.(*MapJsonValKind)
	if !ok {
		return fmt.Errorf("json is not map json")
	}
	jsonValData, err := toJsonValKind(val)
	if err != nil {
		return nil
	}
	structVal.val[key] = &JsonVal{jsonValData}
	return nil
}

func (v *ValueJson) Rm(key string) {
	if v.IsErrOrNil() {
		return
	}
	structVal, ok := v.data.Kind.(*MapJsonValKind)
	if !ok {
		return
	}
	delete(structVal.val, key)
}

func (v *ValueJson) Index(i int) MyJson2 {
	listVal, ok := v.data.Kind.(*SliceJsonValKind)
	if !ok {
		return &NilJson{}
	}

	if i >= len(listVal.val) {
		return new(NilJson)
	}

	valueData := listVal.val[i]
	return &ValueJson{valueData}
}

func insertValue(sliceBody []*JsonVal, index int, val *JsonVal) []*JsonVal {
	// 把尾巴弄出来
	rear := append([]*JsonVal{}, sliceBody[index:]...)
	tpSlice := append(sliceBody[0:index], val)
	return append(tpSlice, rear...)
}

func (v *ValueJson) Insert(i int, val interface{}) (MyJson2, error) {
	listVal, ok := v.data.Kind.(*SliceJsonValKind)
	if !ok {
		return &NilJson{}, fmt.Errorf("json is not slice json")
	}

	jsonValKind, err := toJsonValKind(val)
	if err != nil {
		return v, err
	}
	listVal.val = insertValue(listVal.val, i, &JsonVal{jsonValKind})
	return v, nil
}

func (v *ValueJson) Append(val interface{}) (MyJson2, error) {
	l, ok := v.data.Kind.(*SliceJsonValKind)
	if !ok {
		return &NilJson{}, fmt.Errorf("json is not slice json")
	}
	jsonValKind, err := toJsonValKind(val)
	if err != nil {
		return v, err
	}
	// Debugf("append:%v\n", jsonValKind)
	l.val = append(l.val, &JsonVal{jsonValKind})
	return v, nil
}

func (v *ValueJson) Len() int {
	l, ok := v.data.Kind.(*SliceJsonValKind)
	if !ok {
		return 0
	}
	return len(l.val)
}

func (v *ValueJson) String() string {
	return string(v.Bytes())
}

func (v *ValueJson) Bytes() []byte {
	switch realVal := v.data.Kind.(type) {
	case *BoolJsonValKind:
		return []byte(ToStr(realVal.val))
	case *NullJsonValKind:
		return []byte("null")
	case *NumberJsonValKind:
		return []byte(realVal.val)
	case *StrJsonValKind:
		return []byte(realVal.val)
	case *MapJsonValKind, *SliceJsonValKind:
		jsBytes, err := v.data.MarshalJSON()
		if err != nil {
			Debugf("MapJsonValKind MarshalJSON find err:%v", err)
			return []byte("")
		}
		return jsBytes
	default:
	}
	return []byte("")
}

func (v *ValueJson) Int() (int, error) {
	numberVal, ok := v.data.Kind.(*NumberJsonValKind)
	if !ok {
		return 0, fmt.Errorf("%v is not number", v.data)
	}
	return ToInt(numberVal.val)
}

func (v *ValueJson) Float64() (float64, error) {
	floatVal, ok := v.data.Kind.(*NumberJsonValKind)
	if !ok {
		return 0, fmt.Errorf("%v is not number", v.data)
	}
	return ToFloat64(floatVal.val)
}

func (v *ValueJson) Bool() (bool, error) {
	boolVal, ok := v.data.Kind.(*BoolJsonValKind)
	if !ok {
		return false, fmt.Errorf("%v is not bool", v.data)
	}
	return boolVal.val, nil
}

func (v *ValueJson) RangeSlice(f func(index int, val MyJson2) (bool, error)) error {
	l := v.data.Kind.(*SliceJsonValKind)
	if l == nil {
		return fmt.Errorf("%v is not slice", v.data)
	}
	for i, tpVal := range l.val {
		tpJs := &ValueJson{tpVal}
		ret, err := f(i, tpJs)
		if err != nil {
			return err
		}
		if ret == false {
			return nil
		}
	}
	return nil
}

func (v *ValueJson) RangeMap(f func(key string, val MyJson2) (bool, error)) error {
	mapVal := v.data.Kind.(*MapJsonValKind)
	if mapVal == nil {
		return fmt.Errorf("%v is not map", v.data)
	}
	for key, tpVal := range mapVal.val {
		ret, err := f(key, &ValueJson{tpVal})
		if err != nil {
			return err
		}
		if ret == false {
			return nil
		}
	}
	return nil
}

func (v *ValueJson) Clone() MyJson2 {
	return NewJson(v.String())
}

func (v *ValueJson) IsErrOrNil() bool {
	if v.data == nil {
		return true
	}
	return false
}

func (v *ValueJson) IsSlice() bool {
	if _, ok := v.data.Kind.(*SliceJsonValKind); ok {
		return true
	}
	return false
}

func (v *ValueJson) IsMap() bool {
	if _, ok := v.data.Kind.(*MapJsonValKind); ok {
		return true
	}
	return false
}

func ToStr(obj interface{}) string {
	switch v := obj.(type) {
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", obj)
	}
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
			info := fmt.Sprintf("ToInt, error, overflowd %v", v)
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
			strv = strings.Split(v, ".")[0]
		}
		if strv == "" {
			return 0, nil
		}
		if intv, err := strconv.Atoi(strv); err == nil {
			return intv, nil
		}
	}
	return 0, fmt.Errorf("%v cannot convert to int", intObj)
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
	}
	return 0, fmt.Errorf("%v cannot convert to float", item)
}

func ToBool(item interface{}) (bool, error) {

	switch v := item.(type) {
	case bool:
		return v, nil
	default:
		boolValue, err := strconv.ParseBool(ToStr(item))
		if err != nil {
			return false, fmt.Errorf("%v cannot convert to bool", item)
		}
		return boolValue, nil
	}
}
