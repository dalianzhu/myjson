## myjson前言
因为不想在go里面定义json的对应结构，很久很久以前，我撸了一个偷懒的库。没有什么技术含量，核心技术点就是，使用interface{}作为json的承载，用自己的结构体对json进行了包装。性能大概是原生json的一半，但是，它很好用呀！
## 使用方法
### quick start
更多用法参考用例
```
jsStr := `{"Filters":[{"Name":"ServiceName","Values":["user_123"]}],"Limit":1024,"Offset":0}`
js := NewJson(jsStr)

// 获取limit字段，并转换为int类型
limit, err := js.Get("Limit").Int()
// limit == 1024

// 修改limit为2048
js.Set("Limit", 2048)

limit, err := js.Get("Limit").Int()
// limit == 2048

// 删除limit字段
js.Rm("Limit")
jsStr = js.String() // 此时jsStr里面没有Limit字段

// 组合两个myjson   
js1 := NewJson("{}")
js2 := NewJson("{}")
js2.Set("name", "pig")

js1.Set("js2", js2)
fmt.Println(js1.String()) // {"js2":{"name":"pig"}}
```

### json数组处理
myjson也能方便的构造一个数组

使用myjson的Append方法
```
js:=NewJson("[]")
js.Append("hello")
js.Append("world")
jsStr:=js.String() // jsStr == '["hello", "world"]'
```
遍历数组，如果return false，循环将立刻结束
```
js:=NewJson("[]")
js.Append("hello")
js.Append("world")
jsArr.RangeSlice(func(index int, val interface{}) bool {
	fmt.Println(index, val)
	return true
})
// output:0 hello 
// 1 world

// myjson的对象可以组合，比如下面这个复杂一点的数组
jsArr := NewJson("[]")
jsArr.Append(NewJson(`{"name":"yzh"}`))
jsArr.Append(NewJson(`{"name":"zhh"}`))
jsArr.RangeSlice(func(index int, val interface{}) bool {
	name := NewJson(val).Get("name").String()
	fmt.Println(name)
	return true
})
// output:yzh
// zhh
```
### myjson 字典处理
字典对象除了像quick start中那样增删改查，还能像数组一样遍历
```
jsDict := NewJson("{}")
jsDict.Set("name", "yzh")
jsDict.Set("age", 18)
jsDict.RangeMap(func(key string, val interface{}) bool {
	fmt.Println(key, val)
	return true
})
// output:
// name yzh
// age 18
```

### myjson 获取json值
每个myjson对象，都包含尝试 Int Bool Float64的方法，使用Value可以获取原始数据(一般为[]interface{}, map[string]interface{})
```
jsDict := NewJson("{}")
tpDict := NewJson("{}")
tpDict.Set("name", "yzh")
tpDict.Set("age", 18)
jsDict.Set("user", tpDict) // 两个myjson对象可以组合

age, err := jsDict.Get("user").Get("age").Int()
fmt.Println(age, err)

realVal := jsDict.Value() // type map[string]interface{} 获取原始值
```
判断一个值是否存在，必须使用`IsNil`方法：
```
jsDict := NewJson("{}")
jsDict.Set("name", "yzh")

val := jsDict.Get("haha")
if val.IsNil() {
	fmt.Println("not exist")
}
```

### Clone and 裁剪
使用Clone可以把js深复制一份，对新的jscopy对象操作，不会影响原来的json。注：性能影响
```
js = NewJson(`{"data": 123.321, "arr":[{"name":"yzh"}]}`)
jscopy := js.Clone()
```
有时在打日志的时候，可能需要忽略长的字段，打印前面几位即可，此时可以使用
```
intVal, _ := js.ShortNiceJson()
```
它会深复制这个json，并把所有超过120个字符的值裁剪为"前120位xxxxx......"
