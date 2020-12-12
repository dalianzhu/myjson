package myjson

import (
	"encoding/json"
	"log"
	"testing"
)

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
