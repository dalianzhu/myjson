package myjson

/*
The MIT License (MIT)Copyright (C) <2019> <yinzihao>
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}

func TestMyJson2Simple(t *testing.T) {
	IsDebug = true
	as := assert.New(t)
	// 	jsStr := `
	// {"key1":{"key2":"123"}, "key3":{}}
	// `
	// 	js := NewJson(jsStr)
	// 	Debugf("TestMyJson2:%v", js.Get("key1"))

	// js = NewJson(longJsonVal)
	// Debugf("TestMyJson2:%s", js)

	jsStr := `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024,"Offset":0}`
	js := NewJson(jsStr)

	limit, err := js.Get("Limit").Int()
	as.Equal(1024, limit, "limit值为1024")
	Debugf("err:%v, %s", err, js.String())

	limit, err = ToInt(js.Get("Limit"))
	as.Equal(1024, limit, "limit值为1024")
	Debugf("err:%v, %s", err, js.String())
}

func TestMyJson2Example(t *testing.T) {
	const jsonStream = `
		{"Message": "Hello", "Array": [1, 2, 3], "Null": null, "Number": 1.234}
	`
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v", t, t)
		if dec.More() {
			fmt.Printf(" (more)")
		}
		fmt.Printf("\n")
	}
}

func TestMyJson2Float(t *testing.T) {
	as := assert.New(t)
	js := NewJson(`{"data":{"err":true},"env":{"float":123.321}}`)
	fmt.Println(js.String())
	floatJs := js.Get("env").Get("float")

	intVal, err := floatJs.Int()
	if err != nil {
		t.Error(err)
		return
	}
	as.Equal(123, intVal, "longInt的值不变")

	floatVal, err := js.Get("env").Get("float").Float64()
	if err != nil {
		t.Error(err)
		return
	}
	as.Equal(123.321, floatVal, "longInt的值不变")

}

func TestMyJson2LongInt(t *testing.T) {
	as := assert.New(t)
	js := NewJson(`{"data":{"err":0},"env":{"longInt":60365780445566778}}`)
	fmt.Println(js.String())
	longInt, _ := js.Get("env").Get("longInt").Int()
	as.Equal(60365780445566778, longInt, "longInt的值不变")
}

func TestStruct(t *testing.T) {
	as := assert.New(t)
	structObj := &testStruct{"honoryin", 18}
	js := NewJson(structObj)
	as.Equal("honoryin", js.Get("Name").String(), "结构体name")

	structObj2 := testStruct{"honoryin", 18}
	js = NewJson(structObj2)
	as.Equal("honoryin", js.Get("Name").String(), "结构体name")

	sliceObj := []*testStruct{
		{"haha1", 100},
		{"haha2", 200},
	}
	js = NewJson(sliceObj)
	if js.Index(0).Get("Name").String() != "haha1" {
		as.Fail("")
	}
}

func TestAppend(t *testing.T) {
	as := assert.New(t)
	js := NewJson(`{"sub":[]}`)
	log.Printf("TestAppend js:%v", js)
	if js.Get("sub").IsErrOrNil() {
		t.Error("append sub is not nil")
	}
	js.Get("sub").Append("hello")
	log.Printf("TestAppend: %s\n", js)

	if js.Get("sub").Len() != 1 {
		as.Fail("append failed")
	}
}

func TestMyJson2(t *testing.T) {
	IsDebug = true
	as := assert.New(t)
	// myjson，懒人专用
	/*
		exp1: 使用myjson解析json，基本操作
	*/
	jsStr := `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024,"Offset":0}`
	js := NewJson(jsStr)

	limit, err := js.Get("Limit").Int()
	as.Equal(1024, limit, "limit值为1024")

	// set操作
	js.Set("Limit", 2048)
	limit, _ = js.Get("Limit").Int()
	as.Equal(2048, limit)

	/* rm操作 */
	js = NewJson(`{"name":"yzh", "age":18}`)
	js.Rm("age")
	js.Rm("age123") // 删除一个不存在的key
	as.Equal(true, js.Get("age").IsErrOrNil(), "删除了age字段")

	// float操作
	jsStr = `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024.123,"Offset":0}`
	js = NewJson(jsStr)
	limitStr := js.Get("Limit").String()
	as.Equal("1024.123", limitStr, "limit值为1024.123")

	// append 操作
	js.Get("Filters").Index(0).Get("Values").Append("hello")

	v := js.Get("Filters").Index(0).Get("Values").Index(1).String()
	as.Equal("hello", v, "新添加的字段为hello")

	// insert 操作
	js.Get("Filters").Index(0).Get("Values").Insert(1, "world")
	v = js.Get("Filters").Index(0).Get("Values").Index(1).String()
	as.Equal("world", v, "新插入的字段为 world")

	/*
		exp3: 使用myjson解析slice
	*/
	js = NewJson(jsStr)

	value := js.Get("Filters").Index(0).Get("Values").Index(0).String()
	as.Equal("user_123", value, "value 值为user_123")

	/*
		exp5: 解析结构体为json（如果结构体内有不能json的字段，可能会出错哦）
	*/
	structObj := testStruct{"honoryin", 18}
	js = NewJson(structObj)
	as.Equal("honoryin", js.Get("Name").String(), "结构体name")

	/*
		exp6: 容错
	*/
	js = NewJson(jsStr)
	isnil := js.Get("hello").Get("world").Get("haha").IsErrOrNil()
	as.Equal(isnil, true, "错误的get，返回nil")

	isnil = js.Get("hello").Index(123).Get("haha").IsErrOrNil()
	as.Equal(isnil, true, "错误的get，返回nil 2")

	_, err = js.Get("hello").Index(123).Get("haha").Bool()
	if err == nil {
		t.Error(err)
	}

	_, err = js.Get("hello").Index(123).Get("haha").Int()
	if err == nil {
		t.Error(err)
	}

	v = js.Get("hello").Index(123).Get("haha").String()
	as.Equal("", v, "错误的值tostring返回空串")

	/* exp7: 长数字 */
	js = NewJson(`{"data":{"err":0},"env":{"longInt":60365780445566778}}`)
	fmt.Println(js.String())
	longInt, _ := js.Get("env").Get("longInt").Int()
	as.Equal(60365780445566778, longInt, "longInt的值不变")

	/* exp8: range函数 */

	js = NewJson(`{"data1": 123, "data2": 456}`)
	_ = js.RangeMap(func(key string, val MyJson2) (bool, error) {
		tp, _ := val.Int()
		if tp != 123 && tp != 456 {
			as.Fail("失败")
		}
		return true, nil
	})

	js = NewJson(`[0,1,2]`)
	jsIntArray := make([]int, 0)

	_ = js.RangeSlice(func(index int, val MyJson2) (bool, error) {
		tp, _ := val.Int()
		jsIntArray = append(jsIntArray, tp)
		return true, nil
	})

	as.Equal(jsIntArray[0], 0, "")
	as.Equal(jsIntArray[1], 1, "")
	as.Equal(jsIntArray[2], 2, "")

	/* exp9: float64 */
	js = NewJson(`{"data1": 123.321}`)
	floatVal, _ := js.Get("data1").Float64()
	as.Equal(123.321, floatVal, "转float64测试")

	// exp10: 测试乱码 错误的json
	js = NewJson(`asdf;vjaspoipewqurj`)
	as.Equal(true, js.IsErrOrNil(), fmt.Sprintf("json test is error:%v", js))

	js = NewJson(`{"key":123, "val":[1,2,3,4],}`)
	as.Equal(true, js.IsErrOrNil(), fmt.Sprintf("json test is error:%v", js))

	js = NewJson(`{"key":123, "val":[1,2,3,4]`)
	as.Equal(true, js.IsErrOrNil(), fmt.Sprintf("json test is error:%v", js))

	js = NewJson(`null`)
	as.Equal(false, js.IsErrOrNil(), "json test null val is error")

	// exp11: 测试clone
	js = NewJson(`{"data": 123.321, "arr":[{"name":"yzh"}]}`)
	jscopy := js.Clone()

	// exp11.1: 修改原数据的值
	js.Set("data", 111)
	js.Get("arr").Index(0).Set("name", "haha")

	floatVal, _ = jscopy.Get("data").Float64()
	as.Equal(123.321, floatVal, "被复制的参数不变")
	strVal := jscopy.Get("arr").Index(0).Get("name").String()
	as.Equal("yzh", strVal, "被复制的值不变")

	fmt.Println("origin js:", js)
	fmt.Println("copy js:", jscopy)

	log.Println("all success")

	// exp12: 测试 IsXXX
	js = NewJson(`{"null": null}`)
	as.Equal(false, js.Get("null").IsErrOrNil(), "测试 isNull")
	as.Equal(true, js.Get("null").IsNull(), "测试 isNull")

	js = NewJson("null")
	as.Equal(false, js.IsErrOrNil(), "测试 isNull2")
	as.Equal(true, js.IsNull(), "测试 isNull2")

	js = NewJson("")
	Debugf("empty string is:%v,%v", js.IsErrOrNil(), js.IsNull())
	as.Equal(true, js.IsErrOrNil(), "测试 isNull3")
	as.Equal(false, js.IsNull(), "测试 isNull3")

	js = NewJson(`{"null": null}`)
	as.Equal(true, js.IsMap(), "测试 isMap1")

	js = NewJson(`{}`)
	as.Equal(true, js.IsMap(), "测试 isMap2")
	as.Equal(false, js.IsErrOrNil(), "测试 isMap3")

	js = NewJson(`[]`)
	as.Equal(true, js.IsSlice(), "测试 isSlice1")
	as.Equal(false, js.IsErrOrNil(), "测试 isSlice1")

	err = js.RangeSlice(func(index int, val MyJson2) (bool, error) {
		return true, nil
	})
	as.Equal(err, nil, "测试 isSlice 3")

	err = js.RangeMap(func(key string, val MyJson2) (bool, error) {
		return true, nil
	})
	Debugf("slice json range map:%v", err)
	if err == nil {
		as.Fail("测试 isSlice 4")
	}
}

func BenchmarkTestMyjsonSysJson(b *testing.B) {
	bsVal := []byte(longJsonVal)
	for i := 0; i < b.N; i++ {
		tp := new(testLongJsonStruct)
		err := json.Unmarshal(bsVal, tp)
		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkTestMyjsonSysJsonMarshal(b *testing.B) {
	bsVal := []byte(longJsonVal)
	tp := new(testLongJsonStruct)
	json.Unmarshal(bsVal, tp)

	for i := 0; i < b.N; i++ {
		json.Marshal(tp)
	}
}

func BenchmarkTestMyjson(b *testing.B) {
	bsVal := []byte(longJsonVal)
	for i := 0; i < b.N; i++ {
		NewJson(bsVal)
	}
}

func BenchmarkTestMyjsonMarshal(b *testing.B) {
	bsVal := []byte(longJsonVal)
	js := NewJson(bsVal)
	for i := 0; i < b.N; i++ {
		js.Bytes()
	}
}
