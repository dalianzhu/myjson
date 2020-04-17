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

### 检验json的合法性
引入了validator来检验json的合法性，具体的规则可以参考	"gopkg.in/go-playground/validator.v9" 项目。

简单示例：

每个Key的规则用`;`与errmsg分隔。如果不传errmsg将传出默认的错误
```
rules := `
{
	"name": "gt=5,required;user name must greater than 5, and not empty",
	"age": "gt=20;age must greater than 20"
}
`
data:= `
{
	"name":"hi",
	"age": 25
}

err := Validate(NewJson(data), NewValidateRules(rules))
// err: user name must greater than 5, and not empty
`
```

子key info 是一个k-v结构, 则会深入检测到它的内部:
```
rules := `
{
	"name": "gt=5,required;user name must greater than 5, and not empty",
	"info": {
		"age": "gt=20;age must greater than 20"
	}
}
`
data:= `
{
	"name":"helloworld",
	"info": {
		"age": 15
	}
}

err := Validate(NewJson(data), NewValidateRules(rules))
// err: age must greater than 20
```

子key info是一个KV结构，而且不能不传值，在info的描述下加上`_required`，它的value是返回的错误信息
```
rules := `
{
	"name": "gt=5,required;user name must greater than 5, and not empty",
	"info": {
		"age": "gt=20;age must greater than 20",
		"_required": "info not exists"
	}
}
`
data:= `
{
	"name":"helloworld",
}

err := Validate(NewJson(data), NewValidateRules(rules))
// err: info not exists
```

子key info 是一个数组，可以定义子key规则
```go
rules := `
{
	"name": "gt=5,required;user name must greater than 5, and not empty",
	"info": ["gt=20;info must greater than 20"]
}
`
data:= `
{
	"name":"helloworld",
	"info": [120,25,3]
}

err := Validate(NewJson(data), NewValidateRules(rules))
// err: gt=20;info must greater than 20
```
注意，子key是数组，数组的value可以为任何json类型，遵从上面的规则
```
rules := `
{
	"name": "gt=5,required;user name must greater than 5, and not empty",
	"info": [
		{
			"age": "gt=20"
		}
	]
}
```
如果子key是数组，且这个key不能不传。则将它的第二个value设为`_required;ERRMSG`
```
rules := `
{
	"name": "gt=5,required;user name must greater than 5, and not empty",
	"info": [
		{
			"age": "gt=20"
		},
		"_required;info cannot be empty"
	]
}
```

附带一个复杂的例子:
```json
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
	"company": {
		"_required": "company is empty"
	}
}
```