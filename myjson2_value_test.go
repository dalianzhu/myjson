package myjson

import (
	"encoding/json"
	"log"
	"testing"
)

func TestKeysItems(t *testing.T) {
	js := NewJson("{}")
	keys, err := js.Keys()
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("keys:%v", keys)
	js.Set("hello", "world")
	js.Set("hello1", NewJson(`{"sub":"haha"}`))
	keys, err = js.Keys()
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("keys:%v", keys)
	if !stringIn("hello", keys) {
		t.Fail()
		return
	}
	if !stringIn("hello1", keys) {
		t.Fail()
		return
	}
	_, err = js.Items()
	if err == nil {
		t.Fail()
		return
	}

	js = NewJson("[]")
	items, err := js.Items()
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("items:%v", items)
	js.Append("hello")
	items, err = js.Items()
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("items:%v", items)
	if items[0].(string) != "hello" {
		t.Fail()
		return
	}
}

func stringIn(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}

func TestMapJsonVal_UnmarshalJSON(t *testing.T) {
	js2 := ` [1,2,3,4,{"key":123,
"sub":
[1,2,3,   12341234125123412357890987]},null,1,true,
false]`
	var m = &ValueJson{}
	err := json.Unmarshal([]byte(js2), m)
	if err != nil {
		t.Error(err)
	}
	Debugf("Unmarshal json:%v", m.data)

	v, err := json.Marshal(m)
	log.Printf("mapVal:%s, %v\n", v, err)
}

func TestMapJsonVal_MarshalJSON(t *testing.T) {
	// mapVal:= map[string]interface{}{}
	// err = dec.Decode(&mapVal)
	// log.Printf("i:%v, err:%v\n", mapVal, err)

	// more := dec.More()
	// offset := dec.InputOffset()
	// log.Printf("token:%v, err:%v, more:%v, offset:%v\n", token, err, more, offset)

}
