package myjson

/*
myJson是对json的封装，用interface{} 屏蔽了对结构体的使用依赖。常见操作看用例，看用例
*/

import (
	"bytes"
	sysjson "encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/dalianzhu/logger"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type MyJson struct {
	prev      *MyJson
	prevkey   string
	previndex int
	data      interface{}
}

// NewJson 从其它对象创建myJson对象。
func NewJson(data interface{}) *MyJson {
	switch v := data.(type) {
	case string:
		return NewJsonFromStr(v)
	case []byte:
		return NewJsonFromBytes(v)
	default:
		return NewJsonFromData(v)
	}
}
func decodeJson(r io.Reader) (*MyJson, error) {
	var f interface{}
	decoder := json.NewDecoder(r)
	decoder.UseNumber() // UseNumber causes the Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
	if err := decoder.Decode(&f); err != nil {
		return nil, err
	}

	return &MyJson{data: f}, nil
}

// NewJsonFromBytes 从bytes对象创建myJson对象。bytes对象必须是标准的json格式。
func NewJsonFromBytes(b []byte) *MyJson {
	js, err := decodeJson(bytes.NewReader(b))
	if err != nil {
		errstr := fmt.Sprintf("js解析失败：%s", b)
		return NewErrJson(1, errstr)
	}
	return js
}

// NewJsonFromStr 从一个字符串对象创建myJson对象
func NewJsonFromStr(str string) *MyJson {
	js, err := decodeJson(strings.NewReader(str))
	if err != nil {
		errstr := fmt.Sprintf("js解析失败：%s", str)
		return NewErrJson(1, errstr)
	}
	return js
}

func NewErrJson(errcode int, errmsg string) *MyJson {
	result := NewJsonFromStr("{}")
	result.Set("err_msg", errmsg)
	result.Set("err", errcode)

	return result
}

// NewJsonFromStruct 从一个结构体对象创建myJson对象
func NewJsonFromStruct(b interface{}) *MyJson {
	var f interface{}
	bytesArr, err := json.Marshal(b)
	if err != nil {
		errstr := fmt.Sprintf("js解析失败：%v", err)
		return NewErrJson(1, errstr)
	}

	err = json.Unmarshal(bytesArr, &f)
	if err != nil {
		errstr := fmt.Sprintf("js解析失败：%v", err)
		return NewErrJson(1, errstr)
	}

	return &MyJson{data: f}
}

func NewJsonFromData(d interface{}) *MyJson {
	return &MyJson{data: d}
}

func getMap(key string, mapBody interface{}) (interface{}, bool) {
	switch v := mapBody.(type) {
	case map[string]interface{}:
		return v[key], true
	case Dict:
		return v[key], true
	default:
		return nil, false
	}
}

func setMap(key string, mapBody, data interface{}) bool {
	var val interface{}
	if value, ok := data.(*MyJson); ok {
		val = value.Value()
	} else {
		val = data
	}

	switch v := mapBody.(type) {
	case map[string]interface{}:
		v[key] = val
		return true
	case Dict:
		v[key] = val
		return true
	default:
		return false
	}
}

func getSlice(key int, sliceBody interface{}) (interface{}, bool) {
	switch v := sliceBody.(type) {
	case []interface{}:
		return v[key], true
	case List:
		return v[key], true
	default:
		return nil, false
	}
}

func appendSlice(sliceBody, data interface{}) (interface{}, bool) {
	var val interface{}
	if value, ok := data.(*MyJson); ok {
		val = value.Value()
	} else {
		val = data
	}

	switch v := sliceBody.(type) {
	case []interface{}:
		v = append(v, val)
		return v, true
	case List:
		v = append(v, val)
		return v, true
	default:
		return nil, false
	}
}

func setSlice(index int, sliceBody, data interface{}) bool {
	var val interface{}
	if value, ok := data.(*MyJson); ok {
		val = value.Value()
	} else {
		val = data
	}

	switch v := sliceBody.(type) {
	case []interface{}:
		v[index] = val
		return true
	case List:
		v[index] = val
		return true
	default:
		return false
	}
}

func insertSlice(index int, sliceBody, data interface{}) (interface{}, bool) {
	var val interface{}
	if value, ok := data.(*MyJson); ok {
		val = value.Value()
	} else {
		val = data
	}

	switch v := sliceBody.(type) {
	case []interface{}:
		rear := append([]interface{}{}, v[index:]...)
		v = append(v[0:index], val)
		return append(v, rear...), true
	case List:
		rear := append([]interface{}{}, v[index:]...)
		v = append(v[0:index], val)
		return append(v, rear...), true
	default:
		return nil, false
	}
}

// Get 获取一个key值。
func (j *MyJson) Get(key string) *MyJson {
	m, ok := getMap(key, j.data)
	if !ok {
		return &MyJson{
			prev:    j,
			prevkey: key,
			data:    nil,
		}
	}

	return &MyJson{
		prev:    j,
		prevkey: key,
		data:    m,
	}
}

// maintainParent 维护这个节点与父节点的关系
func maintainParent(child *MyJson) {
	if child.prev == nil {
		return
	}

	switch child.prev.Value().(type) {
	case map[string]interface{}:
		child.prev.Set(child.prevkey, child)
	case []interface{}:
		child.prev.Set(child.previndex, child)
	}
}

// Append 往数组中添加值，当json不为slice，则panic
func (j *MyJson) Append(val interface{}) *MyJson {
	var v interface{}
	if value, ok := val.(*MyJson); ok {
		v = value.Value()
	} else {
		v = val
	}

	data, ok := appendSlice(j.data, v)
	if !ok {
		panic(fmt.Sprintf("%v not slice cannot append", j.data))
	}
	j.data = data

	maintainParent(j)
	return j
}

// Insert 往数组中添加值，当json不为slice，则panic
func (j *MyJson) Insert(index int, val interface{}) *MyJson {
	v, ok := insertSlice(index, j.data, val)
	if !ok {
		panic(fmt.Sprintf("%v not slice cannot insert", j.data))
	}
	j.data = v
	maintainParent(j)
	return j
}

// IsNil 判定data最不是空
func (j *MyJson) IsNil() bool {
	if j.data == nil {
		return true
	}
	return false
}

func (j *MyJson) IsSlice() bool {
	switch j.data.(type) {
	case List:
		return true
	case []interface{}:
		return true
	default:
		return false
	}
}

func (j *MyJson) IsMap() bool {
	switch j.data.(type) {
	case Dict, map[string]interface{}:
		return true
	default:
		return false
	}
}

// Index 传入位置获取slice对应位置的myjson对象
func (j *MyJson) Index(key int) *MyJson {
	v, ok := getSlice(key, j.data)
	if !ok {
		return &MyJson{
			prev:      j,
			previndex: key,
			data:      nil,
		}
	}

	return &MyJson{
		prev:      j,
		previndex: key,
		data:      v,
	}
}

// Set 对当前的myjson对象对应key设置值
func (j *MyJson) Set(key interface{}, val interface{}) *MyJson {
	switch v := key.(type) {
	case string:
		ok := setMap(v, j.data, val)
		if !ok {
			panic(fmt.Sprintf("%v not map cannot set", j.data))
		}
	case int:
		ok := setSlice(v, j.data, val)
		if !ok {
			panic(fmt.Sprintf("%v not map cannot set", j.data))
		}
	}
	return j
}

func (j *MyJson) Rm(key interface{}) *MyJson {
	switch keyVal := key.(type) {
	case string:
		if j.IsMap() {
			switch v := j.data.(type) {
			case map[string]interface{}:
				delete(v, keyVal)
			case Dict:
				delete(v, keyVal)
			}
		}
	case int:
		// 暂时未实现rm slice功能
	}
	return j
}

// Value 返回myjson对象的真实数据
func (j *MyJson) Value() interface{} {
	v := j.data
	return v
}

// Len 返回数组对象的长度
func (j *MyJson) Len() int {
	switch v := j.data.(type) {
	case []interface{}:
		return len(v)
	case List:
		return len(v)
	default:
		return 0
	}
}

// String方法返回myjson对象的字符串值
func (j *MyJson) String() string {
	if j.data == nil {
		return ""
	}
	switch j.data.(type) {
	case map[string]interface{}, []interface{}, Dict, List:
		result, err := json.Marshal(j.data)
		if err != nil {
			logger.Errorln("json to string error" + err.Error())
			return ""
		}
		return string(result)
	default:
		return ToStr(j.data)
	}
}

// Int 返回myjson对象的真实数据
func (j *MyJson) Int() (int, error) {
	v := j.data
	if v == nil {
		return 0, errors.New(fmt.Sprintf("%v not int", j.data))
	}
	return ToInt(v)
}

func (j *MyJson) Float64() (float64, error) {
	v := j.data
	if v == nil {
		return 0, errors.New(fmt.Sprintf("%v not float64", j.data))
	}
	return ToFloat64(v)
}

// Bool 返回myjson对象的真实数据
func (j *MyJson) Bool() (bool, error) {
	v, ok := j.data.(bool)
	if ok {
		return v, nil
	}
	return false, fmt.Errorf("%v not bool", j.data)
}

// Array 返回数组对象的真实数据
func (j *MyJson) Array() ([]interface{}, error) {
	if j.IsSlice() == false {
		return nil, fmt.Errorf("%v not array", j.data)
	}

	switch v := j.data.(type) {
	case List:
		return v, nil
	case []interface{}:
		return v, nil
	default:
		return nil, fmt.Errorf("%v not array", j.data)
	}
}

func (j *MyJson) RangeMap(f func(key string, val interface{}) bool) error {
	if j.IsMap() == false {
		return fmt.Errorf("%v not map", j.data)
	}
	switch v := j.data.(type) {
	case Dict:
		for key, val := range v {
			ret := f(key, val)
			if !ret {
				break
			}
		}
	case map[string]interface{}:
		for key, val := range v {
			ret := f(key, val)
			if !ret {
				break
			}
		}
	}
	return nil
}

func (j *MyJson) RangeSlice(f func(index int, val interface{}) bool) error {
	if j.IsSlice() == false {
		return fmt.Errorf("%v not Slice", j.data)
	}
	switch v := j.data.(type) {
	case List:
		for key, val := range v {
			ret := f(key, val)
			if !ret {
				break
			}
		}
	case []interface{}:
		for key, val := range v {
			ret := f(key, val)
			if !ret {
				break
			}
		}
	}
	return nil
}

func handlerVal(val interface{}, cutLongStr bool) interface{} {
	switch valV := val.(type) {
	case Dict:
		return handlerMap(valV, cutLongStr)
	case map[string]interface{}:
		return handlerMap(valV, cutLongStr)
	case List:
		return handlerSlice(valV, cutLongStr)
	case []interface{}:
		return handlerSlice(valV, cutLongStr)
	case string:
		return handlerString(valV, cutLongStr)
	default:
		return valV
	}
}

func handlerMap(js interface{}, cutLongStr bool) Dict {
	ret := NewDict()
	switch v := js.(type) {
	case Dict:
		for key, val := range v {
			ret[key] = handlerVal(val, cutLongStr)
		}
	case map[string]interface{}:
		for key, val := range v {
			ret[key] = handlerVal(val, cutLongStr)
		}
	}
	return ret
}

func handlerSlice(js interface{}, cutLongStr bool) List {
	ret := NewList()
	switch v := js.(type) {
	case List:
		for _, val := range v {
			ret = append(ret, handlerVal(val, cutLongStr))
		}
	case []interface{}:
		for _, val := range v {
			ret = append(ret, handlerVal(val, cutLongStr))
		}
	}
	return ret
}

func handlerString(js string, cutLongStr bool) string {
	if cutLongStr && len(js) > 120 {
		return js[:120] + "......"
	}
	return js
}

// ShortNiceJson 用来打印的json，把长的string省略成不超过20个字符的数据
func (j *MyJson) ShortNiceJson() *MyJson {
	cutLongStr := true
	if j.IsSlice() {
		return NewJson(handlerSlice(j.data, cutLongStr))
	}
	if j.IsMap() {
		return NewJson(handlerMap(j.data, cutLongStr))
	}
	return NewJson(j.data)
}

// Clone 把这个json对象clone一份，深复制，性能杀手
func (j *MyJson) Clone() *MyJson {
	cutLongStr := false
	if j.IsSlice() {
		return NewJson(handlerSlice(j.data, cutLongStr))
	}
	if j.IsMap() {
		return NewJson(handlerMap(j.data, cutLongStr))
	}
	return NewJson(j.data)
}

func NewList() List {
	l := make([]interface{}, 0)
	return List(l)
}

type List []interface{}

func (l List) Add(i interface{}) List {
	tp, _ := appendSlice(l, i)
	return tp.(List)
}

func (l *List) String() string {
	result, err := json.Marshal(l)
	if err != nil {
		panic("json to string error" + err.Error())
	}
	return string(result)
}

func NewDict() Dict {
	d := make(map[string]interface{})
	return Dict(d)
}

type Dict map[string]interface{}

func (d *Dict) String() string {
	result, err := json.Marshal(d)
	if err != nil {
		logger.Errorf("json to string error, %v", err)
		return ""
	}
	return string(result)
}

func ToStr(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
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
	return 0, fmt.Errorf("%v cannot to int", intObj)
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
	return 0, fmt.Errorf("%v cannot to float", item)
}
