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
	Name string
	Age  int
}

func TestMyJson(t *testing.T) {
	as := assert.New(t)
	// myjson，懒人专用
	/*
		exp1: 使用List,Dict 构造一个json字符串，http请求专用，妈妈再也不用担心我定义一堆json结构了
	*/
	l := NewList().Add("1").Add("2").Add("3")
	l = l.Add("4")
	log.Println("list l", l)
	as.Equal("4", l[3])

	reqData := NewDict()
	reqData["Offset"] = 0
	reqData["Limit"] = 1024

	filters := NewDict()
	filters["Name"] = "ServiceName"
	filters["Values"] = NewList().Add("user_123")

	reqData["Filters"] = NewList().Add(filters)
	//{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024,"Offset":0}
	log.Printf("TestLoopRouter_Call reqData %v", reqData.String())

	/*
		exp2: 使用myjson解析json，基本操作
	*/
	jsStr := `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024,"Offset":0}`
	js := NewJson(jsStr)

	limit, err := js.Get("Limit").Int()
	log.Printf("err %v", err)
	as.Equal(1024, limit, "limit值为1024")

	// set操作
	js.Set("Limit", 2048)
	limit, _ = js.Get("Limit").Int()
	as.Equal(2048, limit)

	/* rm操作 */
	js = NewJson(`{"name":"yzh", "age":18}`)
	js.Rm("age")
	js.Rm("age123") // 删除一个不存在的key
	log.Printf("rm age %v", js)
	as.Equal(true, js.Get("age").IsNil(), "删除了age字段")

	// float操作
	jsStr = `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024.123,"Offset":0}`
	js = NewJson(jsStr)
	limitStr := js.Get("Limit").String()
	log.Printf("js %v", js)
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
		exp4: 混用List, Dict, MyJson
	*/
	list := NewList().Add("hello").Add("world")
	js = NewJson(list)
	as.Equal("hello", js.Index(0).String(), "第0个元素为hello")

	/*
		exp5: 解析结构体为json（如果结构体内有不能json的字段，可能会出错哦）
	*/
	structObj := testStruct{"honoryin", 18}
	js = NewJsonFromStruct(structObj)
	as.Equal("honoryin", js.Get("Name").String(), "结构体name")

	/*
		exp6: 容错
	*/
	js = NewJson(jsStr)
	isnil := js.Get("hello").Get("world").Get("haha").IsNil()
	as.Equal(isnil, true, "错误的get，返回nil")

	isnil = js.Get("hello").Index(123).Get("haha").IsNil()
	as.Equal(isnil, true, "错误的get，返回nil 2")

	_, err = js.Get("hello").Index(123).Get("haha").Bool()
	log.Printf("err %v", err)
	if err == nil {
		t.Fail()
	}

	_, err = js.Get("hello").Index(123).Get("haha").Int()
	log.Printf("err %v", err)
	if err == nil {
		t.Fail()
	}

	v = js.Get("hello").Index(123).Get("haha").String()
	log.Printf("err value %v", v)
	as.Equal("", v, "错误的值tostring返回空串")

	/* 长数字 */
	js = NewJson(`{"data":{"err":0},"env":{"longInt":60365780445566778}}`)
	fmt.Println(js.String())
	longInt, _ := js.Get("env").Get("longInt").Int()
	as.Equal(60365780445566778, longInt, "longInt的值不变")

	/* range函数 */

	js = NewJson(`{"data1": 123, "data2": 456}`)
	_ = js.RangeMap(func(key string, val interface{}) bool {
		tp, _ := ToInt(val)
		if tp != 123 && tp != 456 {
			as.Fail("失败")
		}
		return true
	})

	array := NewList()
	array = append(array, 0)
	array = append(array, 1)
	array = append(array, 2)
	js = NewJson(array)
	arrayVal, err := js.Array()
	log.Printf("Array err is %v, %v", err, arrayVal)

	intArray := make([]int, 0)
	as.Equal(nil, err, "Array 无错误")
	_ = js.RangeSlice(func(index int, val interface{}) bool {
		tp, _ := ToInt(val)
		intArray = append(intArray, tp)
		return true
	})
	as.Equal(0, intArray[0])
	as.Equal(1, intArray[1])
	as.Equal(2, intArray[2])

	js = NewJson(`[0,1,2]`)
	jsIntArray := make([]int, 0)

	_ = js.RangeSlice(func(index int, val interface{}) bool {
		tp, _ := ToInt(val)
		jsIntArray = append(jsIntArray, tp)
		return true
	})

	as.Equal(jsIntArray[0], 0, "")
	as.Equal(jsIntArray[1], 1, "")
	as.Equal(jsIntArray[2], 2, "")

	/* 返回一个裁剪过的json，它把长行给缩短了。nice short json */
	js = NewJson(`{"data1": 123, "data2": [{"name": "111223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445512233445566778899001122334455"}]}`)
	array = NewList()
	array = append(array, 0)
	array = append(array, 1)
	array = append(array, 2)

	js.Set("data3", array)
	log.Printf("nice string %v %v", js, js.ShortNiceJson())

	as.Equal("111223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445......",
		js.ShortNiceJson().Get("data2").Index(0).Get("name").String(),
		"short js长字符串截断")

	as.Equal("111223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445511223344556677889900112233445512233445566778899001122334455",
		js.Get("data2").Index(0).Get("name").String(),
		"js 长字符串 原字符串")

	intVal, _ := js.ShortNiceJson().Get("data3").Index(0).Int()
	as.Equal(0, intVal, "short js 其它不变")

	/* float64 */
	js = NewJson(`{"data1": 123.321}`)
	floatVal, _ := js.Get("data1").Float64()
	as.Equal(123.321, floatVal, "转float64测试")

	// 测试乱码
	js = NewJson(`asdf;vjaspoipewqurj`)
	fmt.Println("error json string", js)
	as.Equal(true, js.IsNil())

	// 测试clone
	js = NewJson(`{"data": 123.321, "arr":[{"name":"yzh"}]}`)
	jscopy := js.Clone()

	// 修改原数据的值
	js.Set("data", 111)
	js.Get("arr").Set(0, NewJson("{}"))

	floatVal, _ = jscopy.Get("data").Float64()
	as.Equal(123.321, floatVal, "被复制的参数不变")

	strVal := jscopy.Get("arr").Index(0).Get("name").String()
	as.Equal("yzh", strVal, "被复制的值不变")
	fmt.Println("origin js", js)
	fmt.Println("copy js", jscopy)
}

func TestValidate(t *testing.T) {
	rules := `
{
	"name": "gt=5,required;user name must greater than 5, and not empty",
	"info": [{
			"years": "gt=2;years must greater than 2"
		},
		"_required;info must has value"
	],
	"attr": {
		"name": "gt=5;attr name must greater than 5",
        "vals": ["gt=10;vals must greater than 10"]
    },
    "willRequired":{
        "_required":"willRequired must has value"
    }
}
`
	rulesJs := NewJson(rules)
	debugf("rulesJs is:%s", rulesJs)

	check := func(origin string, errMsg string) {
		originJs := NewJson(origin)
		err := Validate(originJs, rulesJs)
		log.Printf("CHECK!!@%v@\n", err)
		if errMsg != "" {
			if err == nil {
				log.Println("CHECKFAILED err cant be nil")
				t.Fail()
			}

			if err.Error() != errMsg {
				log.Println("CHECKFAILED err msg not correct")
				t.Fail()
			}
		} else {
			if err != nil {
				log.Println("CHECKFAILED err is not nil")
				t.Fail()
			}
		}
	}

	log.Println("check 1")
	origin := `
{
	"name": "yzhyzh",
	"info": [{
		"years": 10,
		"age": 18
    }],
    "willRequired":{}
}`
	check(origin, "")

	log.Println("check 2")
	origin = `
{
	"info": [{
		"years": 3,
		"age": 18
	}],
    "willRequired":{}
}`
	check(origin, "user name must greater than 5, and not empty")

	log.Println("check 3")
	origin = `
{
	"name": "hello world",
	"info": [{
		"years": 10,
		"age": 18
	}],
	"attr": {
		"name": "hi",
		"vals": [15]
	},
    "willRequired":{}
}`
	check(origin, "attr name must greater than 5")

	log.Println("check 4")
	origin = `
{
	"name": "hello world",
	"info": [{
		"years": 10,
		"age": 18
	}],
	"attr": {
		"name": "hello world",
		"vals": [1, 2, 3, 4]
    },
    "willRequired":{}
}`
	check(origin, "vals must greater than 10")

	log.Println("check 5")
	origin = `
{
	"name": "hello world",
	"info": [{
		"years": 10,
		"age": 18
	}],
	"attr": {
		"name": "hello world",
		"vals": [15,22]
	},
    "willRequired":{}
}`
	check(origin, "")

	log.Println("check 6")
	origin = `
{
	"name": "hello world",
	"attr": {
		"name": "hello world",
		"vals": [15,22]
	},
    "willRequired":{}
}`
	check(origin, "info must has value")

	log.Println("check 7")
	origin = `
{
    "name":"hello hello",
	"info": [{
		"years": 3,
		"age": 18
	},{
		"years": 1,
		"age": 18
	}],
    "willRequired":{}
}`
	check(origin, "years must greater than 2")

}
