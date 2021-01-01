package myjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var bytesTrue = []byte("true")
var bytesFalse = []byte("false")

var bytesQuoto = []byte(`"`)
var bytesQuotoReplaced = []byte(`\"`)

var bytesSlash = []byte(`\`)
var bytesSlashReplaced = []byte(`\\`)

func bytesToJsBytes(btVal []byte) []byte {
	tp := bytes.Replace(btVal, bytesSlash, bytesSlashReplaced, -1)
	tp = bytes.Replace(tp, bytesQuoto, bytesQuotoReplaced, -1)
	return tp
}

type sliceWrap struct {
	sliceData []interface{}
}

func (s *sliceWrap) MarshalJSON() ([]byte, error) {
	Debugf("sliceWrap MarshalJson:")
	return jsonit.Marshal(s.sliceData)
}

type nullWrap struct {
}

var globalNullWrap = &nullWrap{}

func GetJsonNull() *nullWrap {
	return globalNullWrap
}

var bytesNull = []byte("null")

func (s *nullWrap) MarshalJSON() ([]byte, error) {
	// Debugf("nullWrap MarshalJson:")
	return bytesNull, nil
}

var StopIter = errors.New("iter stop")

func getDecodeVal(dec *json.Decoder, t json.Token, op func(interface{})) error {
	switch v := t.(type) {
	default:
		op(t)
	case nil:
		op(globalNullWrap)
	// 还可能是一个deli
	case json.Delim:
		// 如果是一个map，则转入下轮
		deli := v.String()
		if deli == "{" {
			mapVal := make(map[string]interface{}, 10)
			err := decodeMap(dec, mapVal)
			if err != nil {
				return err
			}
			op(mapVal)
		} else if deli == "[" {
			sliceVal := &sliceWrap{}
			sliceVal.sliceData = make([]interface{}, 0, 10)
			err := decodeSlice(dec, sliceVal)
			if err != nil {
				return err
			}
			op(sliceVal)
		} else if deli == "]" || deli == "}" {
			return StopIter
		}
	}
	return nil
}

func decodeSlice(dec *json.Decoder, sliceVal *sliceWrap) error {
	Debugf("decodeSlice Run")
	for {
		// slice 则只有值
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// 此时val有两种情况，一种是普通值
		err = getDecodeVal(dec, t, func(i interface{}) {
			sliceVal.sliceData = append(sliceVal.sliceData, i)
			Debugf("slice append:%v", i)
		})
		if errors.Is(err, StopIter) {
			// 正常退出，否则json不正确
			return nil
		}
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("json is error")
}

func decodeMap(dec *json.Decoder, m map[string]interface{}) error {
	isSetKey := false
	key := ""
	//  map 轮换key，val
	for {
		if isSetKey == false {
			// 此时获取key
			t, err := dec.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			v, ok := t.(json.Delim)
			if ok && v.String() == "}" {
				return nil
			}

			key = t.(string)
			isSetKey = true
		} else {
			// 已经设置了key
			isSetKey = false
			t, err := dec.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			err = getDecodeVal(dec, t, func(i interface{}) {
				m[key] = i
				Debugf("map set %v:%v", key, i)
			})
			// 本个map发现了 }，处理完了，正常退出
			if errors.Is(err, StopIter) {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("json is error")
}
