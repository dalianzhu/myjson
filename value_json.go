package myjson

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type ValueJson struct {
	data interface{}
}

func (v *ValueJson) SetData(i interface{}) {
	data, err := valueToJsonGoVal(i)
	if err != nil {
		return
	}
	v.data = data
	return
}

func (v *ValueJson) GetValue() interface{} {
	return v.data
}

func (v *ValueJson) Get(key string) MyJson2 {
	if v.IsErrOrNil() {
		return v
	}

	structVal, ok := v.data.(map[string]interface{})
	if !ok {
		Debugf("ValueJson Get:%v", reflect.TypeOf(v.data))
		return NewErrorJson(fmt.Errorf("The data is not a map, you cannot use Get"))
	}

	jsonValData, ok := structVal[key]
	if ok {
		return &ValueJson{jsonValData}
	}
	return NewErrorJson(nil)
}

func valueToJsonGoVal(val interface{}) (interface{}, error) {
	switch setData := val.(type) {
	case MyJson2:
		return setData.GetValue(), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return val, nil
	case string:
		return val, nil
	case nil, *NullValue:
		return GetJsonValNull(), nil
	case *SliceValue:
		return val, nil
	case time.Time:
		return setData.Format("2006-01-02 15:04:05"), nil
	case bool:
		return val, nil
	case []interface{}:
		sliceVal := &SliceValue{}
		sliceVal.sliceData = setData
		return sliceVal, nil
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
	listVal, ok := v.data.(*SliceValue)
	if !ok {
		return NewErrorJson(fmt.Errorf("The data is not a slice, you cannot use Index"))
	}

	if i >= len(listVal.sliceData) {
		return NewErrorJson(fmt.Errorf("Index exceeds the maximum length"))
	}

	valueData := listVal.sliceData[i]
	return &ValueJson{valueData}
}

func insertValue(sliceBody *SliceValue, index int, val interface{}) {
	// 把尾巴弄出来
	rear := append([]interface{}{}, sliceBody.sliceData[index:]...)
	tpSlice := append(sliceBody.sliceData[0:index], val)
	sliceBody.sliceData = append(tpSlice, rear...)
}

func (v *ValueJson) Insert(i int, val interface{}) error {
	listVal, ok := v.data.(*SliceValue)
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
	l, ok := v.data.(*SliceValue)
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
	l, ok := v.data.(*SliceValue)
	if !ok {
		return 0
	}
	return len(l.sliceData)
}

func (v *ValueJson) String() string {
	switch objValue := v.data.(type) {
	case string:
		return objValue
	case *SliceValue, map[string]interface{}:
		ret, err := v.MarshalJSON()
		if err != nil {
			return err.Error()
		}
		return string(ret)
	default:
		return ToStr(v.data)
	}
}

func (v *ValueJson) Bytes() []byte {
	switch v.data.(type) {
	case *SliceValue, map[string]interface{}:
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
	l, ok := n.data.(*SliceValue)
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
	l, ok := v.data.(*SliceValue)
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
	_, ok := v.data.(*ErrorValue)
	if ok {
		return true
	}
	if v.data == nil {
		return true
	}
	return false
}

func (v *ValueJson) IsSlice() bool {
	if _, ok := v.data.(*SliceValue); ok {
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
	_, ok := v.data.(*NullValue)
	if ok {
		return true
	}
	return false
}

func ToStr(obj interface{}) string {
	switch v := obj.(type) {
	case string:
		return v
	case *NullValue:
		return "null"
	case *SliceValue:
		ret, err := v.MarshalJSON()
		if err != nil {
			return ""
		}
		return string(ret)
	case []byte:
		return string(v)
	case MyJson2:
		return v.String()
	case *ErrorValue:
		return v.Err.Error()
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

func switchValue(iter *jsoniter.Iterator) interface{} {
	switch iter.WhatIsNext() {
	case jsoniter.ArrayValue:
		sliceVal := &SliceValue{}
		sliceVal.sliceData = make([]interface{}, 0, 10)
		ret := iter.ReadArrayCB(func(i *jsoniter.Iterator) bool {
			sliceVal.sliceData = append(sliceVal.sliceData, switchValue(i))
			return true
		})
		if ret == false {
			return NewErrorValue(nil)
		}
		return sliceVal
	case jsoniter.ObjectValue:
		mapVal := make(map[string]interface{})
		ret := iter.ReadMapCB(func(i *jsoniter.Iterator, s string) bool {
			mapVal[s] = switchValue(i)
			return true
		})
		if ret == false {
			return NewErrorValue(nil)
		}
		return mapVal
	case jsoniter.NilValue:
		iter.Read()
		return globalNullWrap
	case jsoniter.NumberValue:
		return iter.ReadNumber()
	case jsoniter.StringValue:
		return iter.ReadString()
	case jsoniter.BoolValue:
		return iter.ReadBool()
	case jsoniter.InvalidValue:
		Debugf("find InvalidValue")
		iter.Read()
		return NewErrorValue(nil)
	}
	return NewErrorValue(nil)
}

func (v *ValueJson) UnmarshalJSON(bytesVal []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, bytesVal)
	v.data = switchValue(iter)
	return nil
}

func (v *ValueJson) MarshalJSON() ([]byte, error) {
	return jsonit.Marshal(v.data)
}
