package myjson

import (
	"encoding/json"
	sysjson "encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"
)

type MyJson2 interface {
	Get(key string) MyJson2
	Set(key string, val interface{})
	Rm(key string)
	Index(i int) MyJson2
	Insert(i int, val interface{}) MyJson2
	Append(val interface{}) MyJson2

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

func NewJson(val interface{}) MyJson2 {
	switch v := val.(type) {
	case string:
		pbVal := new(structpb.Value)
		err := pbVal.UnmarshalJSON([]byte(v))
		if err != nil {
			return &NilJson{}
		}
		return &ValueJson{pbVal}
	case []byte:
		pbVal := new(structpb.Value)
		err := pbVal.UnmarshalJSON(v)
		if err != nil {
			return &NilJson{}
		}
		return &ValueJson{pbVal}
	}

	refVal := reflect.ValueOf(val)
	switch refVal.Kind() {
	case reflect.Struct, reflect.Slice:
		bytesArr, err := json.Marshal(val)
		if err != nil {
			// errstr := fmt.Sprintf("js解析失败：%v", err)
			return &NilJson{}
		}
		pbVal := new(structpb.Value)
		pbVal.UnmarshalJSON(bytesArr)
		return &ValueJson{pbVal}
	case reflect.Ptr:
		if refVal.Elem().Kind() == reflect.Struct {
			bytesArr, err := json.Marshal(val)
			if err != nil {
				// errstr := fmt.Sprintf("js解析失败：%v", err)
				return &NilJson{}
			}
			pbVal := new(structpb.Value)
			pbVal.UnmarshalJSON(bytesArr)
			return &ValueJson{pbVal}
		}
	}
	return &NilJson{}
}

type ValueJson struct {
	data *structpb.Value
}

func (v *ValueJson) PbValue() *structpb.Value {
	return v.data
}

func (v *ValueJson) Get(key string) MyJson2 {
	if v.IsErrOrNil() {
		return new(NilJson)
	}

	structVal := v.data.GetStructValue()
	if structVal != nil {
		v, ok := structVal.Fields[key]
		if ok {
			return &ValueJson{v}
		}
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

func (v *ValueJson) Set(key string, val interface{}) {
	if v.IsErrOrNil() {
		return
	}
	structVal := v.data.GetStructValue()
	if structVal != nil {
		structVal.Fields[key] = iterToPbValue(val)
	}
}

func (v *ValueJson) Rm(key string) {
	if v.IsErrOrNil() {
		return
	}
	structVal := v.data.GetStructValue()
	if structVal != nil {
		delete(v.data.GetStructValue().Fields, key)
	}
}

func (v *ValueJson) Index(i int) MyJson2 {
	l := v.data.GetListValue()
	if l == nil {
		return new(NilJson)
	}
	if i >= len(l.Values) {
		return new(NilJson)
	}

	valueData := l.Values[i]
	return &ValueJson{valueData}
}

func insertValue(sliceBody []*structpb.Value, index int, val *structpb.Value) []*structpb.Value {
	// 把尾巴弄出来
	rear := append([]*structpb.Value{}, sliceBody[index:]...)
	tpSlice := append(sliceBody[0:index], val)
	return append(tpSlice, rear...)
}

func (v *ValueJson) Insert(i int, val interface{}) MyJson2 {
	l := v.data.GetListValue()
	if l == nil {
		return new(NilJson)
	}
	l.Values = insertValue(l.Values, i, iterToPbValue(val))
	return v
}

func (v *ValueJson) Append(val interface{}) MyJson2 {
	l := v.data.GetListValue()
	if l == nil {
		return new(NilJson)
	}
	l.Values = append(l.Values, iterToPbValue(val))
	return v
}

func (v *ValueJson) Len() int {
	l := v.data.GetListValue()
	if l == nil {
		return 0
	}
	return len(l.Values)
}

func (v *ValueJson) String() string {
	return string(v.Bytes())
}

func (v *ValueJson) Bytes() []byte {
	switch realVal := v.data.Kind.(type) {
	case *structpb.Value_BoolValue:
		return []byte(ToStr(realVal.BoolValue))
	case *structpb.Value_NullValue:
		return []byte("")
	case *structpb.Value_NumberValue:
		return []byte(ToStr(realVal.NumberValue))
	case *structpb.Value_StringValue:
		return []byte(realVal.StringValue)

	case *structpb.Value_ListValue, *structpb.Value_StructValue:
		jsBytes, err := v.data.MarshalJSON()
		if err != nil {
			return []byte("")
		}
		return jsBytes
	}
	return []byte("")
}

func (v *ValueJson) Int() (int, error) {
	floatVal, ok := v.data.GetKind().(*structpb.Value_NumberValue)
	if !ok {
		return 0, fmt.Errorf("%v is not number", v.data)
	}
	intVal, _ := ToInt(floatVal.NumberValue)
	return intVal, nil
}

func (v *ValueJson) Float64() (float64, error) {
	floatVal, ok := v.data.GetKind().(*structpb.Value_NumberValue)
	if !ok {
		return 0, fmt.Errorf("%v is not number", v.data)
	}
	return floatVal.NumberValue, nil
}

func (v *ValueJson) Bool() (bool, error) {
	boolVal, ok := v.data.GetKind().(*structpb.Value_BoolValue)
	if !ok {
		return false, fmt.Errorf("%v is not bool", v.data)
	}
	return boolVal.BoolValue, nil
}

func (v *ValueJson) RangeSlice(f func(index int, val MyJson2) (bool, error)) error {
	l := v.data.GetListValue()
	if l == nil {
		return fmt.Errorf("%v is not slice", v.data)
	}
	for i, tpVal := range l.Values {
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
	mapVal := v.data.GetStructValue()
	if mapVal == nil {
		return fmt.Errorf("%v is not map", v.data)
	}
	for key, tpVal := range mapVal.Fields {
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
	if v.data.GetListValue() == nil {
		return false
	}
	return true
}

func (v *ValueJson) IsMap() bool {
	if v.data.GetStructValue() == nil {
		return false
	}
	return true
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
	switch v := intObj.(type) {
	case sysjson.Number:
		vint64, err := v.Int64()
		return int(vint64), err
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
