package myjson

/*
The MIT License (MIT)Copyright (C) <2019> <yinzihao>
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
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
	js.Get("sub").Append("hello")
	log.Printf("TestAppend: %s\n", js)
	if js.Get("sub").Len() != 1 {
		as.Fail("")
	}
}

func TestMyJson2(t *testing.T) {
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

	/* 长数字 */
	// js = NewJson2(`{"data":{"err":0},"env":{"longInt":60365780445566778}}`)
	// fmt.Println(js.String())
	// longInt, _ := js.Get("env").Get("longInt").Int()
	// as.Equal(60365780445566778, longInt, "longInt的值不变")

	/* range函数 */

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

	/* float64 */
	js = NewJson(`{"data1": 123.321}`)
	floatVal, _ := js.Get("data1").Float64()
	as.Equal(123.321, floatVal, "转float64测试")

	// 测试乱码 错误的json
	js = NewJson(`asdf;vjaspoipewqurj`)
	as.Equal(true, js.IsErrOrNil(), fmt.Sprintf("json test is error:%v", js))

	js = NewJson(`null`)
	as.Equal(false, js.IsErrOrNil(), "json test null val is error")

	// 测试clone
	js = NewJson(`{"data": 123.321, "arr":[{"name":"yzh"}]}`)
	jscopy := js.Clone()

	// 修改原数据的值
	js.Set("data", 111)
	js.Get("arr").Index(0).Set("name", "haha")

	floatVal, _ = jscopy.Get("data").Float64()
	as.Equal(123.321, floatVal, "被复制的参数不变")
	strVal := jscopy.Get("arr").Index(0).Get("name").String()
	as.Equal("yzh", strVal, "被复制的值不变")

	fmt.Println("origin js:", js)
	fmt.Println("copy js:", jscopy)

	log.Println("all success")
}
