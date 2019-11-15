package myjson

import (
	"fmt"
	"testing"

	"github.com/dalianzhu/logger"
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
	logger.Infoln("list l", l)
	as.Equal("4", l[3])

	reqData := NewDict()
	reqData["Offset"] = 0
	reqData["Limit"] = 1024

	filters := NewDict()
	filters["Name"] = "ServiceName"
	filters["Values"] = NewList().Add("user_123")

	reqData["Filters"] = NewList().Add(filters)
	//{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024,"Offset":0}
	logger.Debugf("TestLoopRouter_Call reqData %v", reqData.String())

	/*
		exp2: 使用myjson解析json，基本操作
	*/
	jsStr := `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024,"Offset":0}`
	js := NewJson(jsStr)

	limit, err := js.Get("Limit").Int()
	logger.Debugf("err %v", err)
	as.Equal(1024, limit, "limit值为1024")

	// set操作
	js.Set("Limit", 2048)
	limit, _ = js.Get("Limit").Int()
	as.Equal(2048, limit)

	/* rm操作 */
	js = NewJson(`{"name":"yzh", "age":18}`)
	js.Rm("age")
	js.Rm("age123") // 删除一个不存在的key
	logger.Debugf("rm age %v", js)
	as.Equal(true, js.Get("age").IsNil(), "删除了age字段")

	// float操作
	jsStr = `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024.123,"Offset":0}`
	js = NewJson(jsStr)
	limitStr := js.Get("Limit").String()
	logger.Debugf("js %v", js)
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
	logger.Debugf("err %v", err)
	if err == nil {
		t.Fail()
	}

	_, err = js.Get("hello").Index(123).Get("haha").Int()
	logger.Debugf("err %v", err)
	if err == nil {
		t.Fail()
	}

	v = js.Get("hello").Index(123).Get("haha").String()
	logger.Debugf("err value %v", v)
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
	logger.Debugf("Array err is %v, %v", err, arrayVal)

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
	logger.Debugf("nice string %v %v", js, js.ShortNiceJson())

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
	intVal, _ = js.Get("err").Int()
	as.Equal(1, intVal, "")

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
