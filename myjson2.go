package myjson

import (
	"bytes"
	"encoding/json"
	sysjson "encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/structpb"
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

	PbValue() *structpb.Value
}

func NewJsonFromBytes(bytesVal []byte) MyJson2 {
	val := &ValueJson{}
	err := val.UnmarshalJSON(bytesVal)
	if err != nil {
		return &NilOrErrJson{}
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
			return &NilOrErrJson{}
		}
		return NewJsonFromBytes(bytesVal)

	case reflect.Ptr:
		if refVal.Elem().Kind() == reflect.Struct {
			bytesVal, err := json.Marshal(val)
			if err != nil {
				// errstr := fmt.Sprintf("js解析失败：%v", err)
				return &NilOrErrJson{}
			}
			return NewJsonFromBytes(bytesVal)
		}
	}
	return &NilOrErrJson{}
}

type ValueJson struct {
	data interface{}
}

func (v *ValueJson) SetData(i interface{}) {
	if i == nil {
		v.data = globalNullWrap
		return
	}
	switch val := i.(type) {
	case []interface{}:
		sliceVal := &sliceWrap{}
		sliceVal.sliceData = val
		v.data = sliceVal
		return
	default:
		v.data = i
	}
}

func (v *ValueJson) GetValue() interface{} {
	return v.data
}

func (v *ValueJson) PbValue() *structpb.Value {
	bytesVal, err := json.Marshal(v.data)
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
		return new(NilOrErrJson)
	}

	structVal, ok := v.data.(map[string]interface{})
	if !ok {
		Debugf("ValueJson Get:%v", reflect.TypeOf(v.data))
		return new(NilOrErrJson)
	}

	jsonValData, ok := structVal[key]
	if ok {
		return &ValueJson{jsonValData}
	}
	return new(NilOrErrJson)
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
}

func valueToJsonGoVal(val interface{}) (interface{}, error) {
	switch setData := val.(type) {
	case MyJson2:
		return setData.GetValue(), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return val, nil
	case string:
		return val, nil
	case nil, *nullWrap:
		return GetJsonNull(), nil
	case *sliceWrap:
		return val, nil
	case time.Time:
		return setData.Format("2006-01-02 15:04:05"), nil
	case bool:
		return val, nil
	default:
		return nil, fmt.Errorf("val:%v cannot set to json", val)
	}
}

func (v *ValueJson) Set(key string, val interface{}) error {
	if v.IsErrOrNil() {
		return fmt.Errorf("json is nil")
	}
	structVal, ok := v.data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("json is not map json")
	}

	goVal, err := valueToJsonGoVal(val)
	if err != nil {
		return err
	}
	structVal[key] = goVal
	return nil
}

func (v *ValueJson) Rm(key string) {
	if v.IsErrOrNil() {
		return
	}
	structVal, ok := v.data.(map[string]interface{})
	if !ok {
		return
	}
	delete(structVal, key)
}

func (v *ValueJson) Index(i int) MyJson2 {
	listVal, ok := v.data.(*sliceWrap)
	if !ok {
		return &NilOrErrJson{}
	}

	if i >= len(listVal.sliceData) {
		return new(NilOrErrJson)
	}

	valueData := listVal.sliceData[i]
	return &ValueJson{valueData}
}

func insertValue(sliceBody *sliceWrap, index int, val interface{}) {
	// 把尾巴弄出来
	rear := append([]interface{}{}, sliceBody.sliceData[index:]...)
	tpSlice := append(sliceBody.sliceData[0:index], val)
	sliceBody.sliceData = append(tpSlice, rear...)
}

func (v *ValueJson) Insert(i int, val interface{}) error {
	listVal, ok := v.data.(*sliceWrap)
	if !ok {
		return fmt.Errorf("json is not slice json")
	}

	tpValue, err := valueToJsonGoVal(val)
	if err != nil {
		return err
	}
	insertValue(listVal, i, tpValue)
	return nil
}

func (v *ValueJson) Append(val interface{}) error {
	l, ok := v.data.(*sliceWrap)
	if !ok {
		return fmt.Errorf("json is not slice json")
	}
	tpValue, err := valueToJsonGoVal(val)
	if err != nil {
		return err
	}
	// Debugf("append:%v\n", jsonValKind)
	l.sliceData = append(l.sliceData, tpValue)
	return nil
}

func (v *ValueJson) Len() int {
	l, ok := v.data.(*sliceWrap)
	if !ok {
		return 0
	}
	return len(l.sliceData)
}

func (v *ValueJson) String() string {
	switch objValue := v.data.(type) {
	case string:
		return objValue
	case *sliceWrap, map[string]interface{}:
		ret, _ := v.MarshalJSON()
		return string(ret)
	default:
		return ToStr(v.data)
	}
}

func (v *ValueJson) Bytes() []byte {
	switch v.data.(type) {
	case *sliceWrap, map[string]interface{}:
		ret, _ := v.MarshalJSON()
		return ret
	default:
		return []byte(ToStr(v.data))
	}
}

func (v *ValueJson) Int() (int, error) {
	return ToInt(v.data)
}

func (v *ValueJson) Float64() (float64, error) {
	return ToFloat64(v.data)
}

func (v *ValueJson) Bool() (bool, error) {
	return ToBool(v.data)
}

func (v *ValueJson) Keys() ([]string, error) {
	mapVal, ok := v.data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("json is not map, has no keys")
	}
	ret := make([]string, 0, 10)
	for key := range mapVal {
		ret = append(ret, key)
	}
	return ret, nil
}

func (n *ValueJson) Items() ([]interface{}, error) {
	l, ok := n.data.(*sliceWrap)
	if !ok {
		return nil, fmt.Errorf("%v is not slice, has no items", n.data)
	}
	ret := make([]interface{}, 0, 10)
	for _, tpVal := range l.sliceData {
		ret = append(ret, tpVal)
	}
	return ret, nil
}

func (v *ValueJson) RangeSlice(f func(index int, val MyJson2) (bool, error)) error {
	l, ok := v.data.(*sliceWrap)
	if !ok {
		return fmt.Errorf("%v is not slice", v.data)
	}
	for i, tpVal := range l.sliceData {
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
	mapVal, ok := v.data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%v is not map", v.data)
	}
	for key, tpVal := range mapVal {
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
	return NewJson(v.Bytes())
}

func (v *ValueJson) IsErrOrNil() bool {
	_, ok := v.data.(*NilOrErrJson)
	if ok {
		return true
	}
	if v.data == nil {
		return true
	}
	return false
}

func (v *ValueJson) IsSlice() bool {
	if _, ok := v.data.(*sliceWrap); ok {
		return true
	}
	return false
}

func (v *ValueJson) IsMap() bool {
	if _, ok := v.data.(map[string]interface{}); ok {
		return true
	}
	return false
}

func (v *ValueJson) IsNull() bool {
	_, ok := v.data.(*nullWrap)
	if ok {
		return true
	}
	return false
}

func ToStr(obj interface{}) string {
	switch v := obj.(type) {
	case string:
		return v
	case *nullWrap:
		return string(bytesNull)
	case *sliceWrap:
		ret, err := v.MarshalJSON()
		if err != nil {
			return ""
		}
		return string(ret)
	case []byte:
		return string(v)
	case MyJson2:
		return v.String()
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
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
	case MyJson2:
		return v.Int()
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

func (v *ValueJson) MarshalJSON() ([]byte, error) {
	return jsonit.Marshal(v.data)
}

func (v *ValueJson) UnmarshalJSON(bytesVal []byte) error {
	Debugf("ValueJson UnmarshalJSON run")
	dec := json.NewDecoder(bytes.NewReader(bytesVal))
	dec.UseNumber()
	for {
		t, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		Debugf("ValueJson UnmarshalJSON token: %T %v", t, t)
		// 此时val有两种情况，一种是普通值
		switch typeVal := t.(type) {
		default:
			v.data = typeVal
		case nil:
			v.data = globalNullWrap
		// 还可能是一个deli
		case json.Delim:
			// 如果是一个map，则转入下轮
			deli := typeVal.String()
			if deli == "{" {
				Debugf("ValueJson UnmarshalJSON map case")
				m := make(map[string]interface{}, 10)
				err = decodeMap(dec, m)
				if err != nil {
					return err
				}
				v.data = m
			} else if deli == "[" {
				sliceVal := &sliceWrap{}
				sliceVal.sliceData = make([]interface{}, 0, 10)
				err = decodeSlice(dec, sliceVal)
				if err != nil {
					return err
				}
				v.data = sliceVal
			} else if deli == "]" || deli == "}" {
				return nil
			}
		}
	}
}
