package myjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type JsonVal struct {
	Kind JsonValKind
}

func decodeSlice(dec *json.Decoder, j *JsonVal) error {
	kind := &SliceJsonValKind{make([]*JsonVal, 0, 10)}
	j.Kind = kind
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
		switch v := t.(type) {
		case string:
			kind.val = append(kind.val, &JsonVal{&StrJsonValKind{v}})
		case bool:
			kind.val = append(kind.val, &JsonVal{&BoolJsonValKind{v}})
		case json.Number:
			kind.val = append(kind.val, &JsonVal{&NumberJsonValKind{v}})
		case nil:
			kind.val = append(kind.val, &JsonVal{&NullJsonValKind{}})
		// 还可能是一个deli
		case json.Delim:
			// 如果是一个map，则转入下轮
			deli := v.String()
			if deli == "{" {
				newJsonVal := &JsonVal{}
				err = decodeMap(dec, newJsonVal)
				if err != nil {
					return err
				}
				kind.val = append(kind.val, newJsonVal)
			} else if deli == "[" {
				newJsonVal := &JsonVal{}
				err = decodeSlice(dec, newJsonVal)
				if err != nil {
					return err
				}
				kind.val = append(kind.val, newJsonVal)
			} else if deli == "]" {
				return nil
			}
		}
	}
	return nil
}

func decodeMap(dec *json.Decoder, j *JsonVal) error {
	kind := &MapJsonValKind{make(map[string]*JsonVal, 10)}
	j.Kind = kind
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
				break
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
			// 此时val有两种情况，一种是普通值
			switch v := t.(type) {
			case string:
				kind.val[key] = &JsonVal{&StrJsonValKind{v}}
			case bool:
				kind.val[key] = &JsonVal{&BoolJsonValKind{v}}
			case json.Number:
				kind.val[key] = &JsonVal{&NumberJsonValKind{v}}
			case nil:
				kind.val[key] = &JsonVal{&NullJsonValKind{}}
			// 还可能是一个deli
			case json.Delim:
				deli := v.String()
				// 如果是一个map，则转入下轮
				if deli == "{" {
					newJsonVal := &JsonVal{}
					err = decodeMap(dec, newJsonVal)
					if err != nil {
						return err
					}
					kind.val[key] = newJsonVal
				} else if deli == "[" {
					newJsonVal := &JsonVal{}
					err = decodeSlice(dec, newJsonVal)
					if err != nil {
						return err
					}
					kind.val[key] = newJsonVal
				} else if deli == "}" {
					return nil
				}
			}
		}
	}

	return nil
}

func (j *JsonVal) UnmarshalJSON(bytesVal []byte) (err error) {
	dec := json.NewDecoder(bytes.NewReader(bytesVal))
	dec.UseNumber()
	for {
		t, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// 此时val有两种情况，一种是普通值
		switch v := t.(type) {
		case string:
			j.Kind = &JsonVal{&StrJsonValKind{v}}
		case bool:
			j.Kind = &JsonVal{&BoolJsonValKind{v}}
		case json.Number:
			j.Kind = &JsonVal{&NumberJsonValKind{v}}
		case nil:
			j.Kind = &JsonVal{&NullJsonValKind{}}
		// 还可能是一个deli
		case json.Delim:
			// 如果是一个map，则转入下轮
			deli := v.String()
			if deli == "{" {
				err = decodeMap(dec, j)
				if err != nil {
					return err
				}
			} else if deli == "[" {
				err = decodeSlice(dec, j)
				if err != nil {
					return err
				}
			} else if deli == "]" || deli == "}" {
				return nil
			}
		}
	}
	return nil
}

func (j *JsonVal) MarshalJSON() ([]byte, error) {
	return j.Kind.MarshalJSON()
}

type JsonValKind interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type MapJsonValKind struct {
	val map[string]*JsonVal
}

func (m *MapJsonValKind) UnmarshalJSON(bytesVal []byte) error {
	return nil
}

func (m *MapJsonValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.val)
}

type SliceJsonValKind struct {
	val []*JsonVal
}

func (s *SliceJsonValKind) UnmarshalJSON(bytesVal []byte) error {
	return nil
}

func (s *SliceJsonValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.val)
}

// ------------------------------------------------

type NumberJsonValKind struct {
	val json.Number
}

func (n *NumberJsonValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.val)
}

func (n *NumberJsonValKind) UnmarshalJSON(bytes []byte) error {
	Debugf("NumberJsonValKind UnmarshalJSON:%s", bytes)
	return json.Unmarshal(bytes, &n.val)
}

type StrJsonValKind struct {
	val string
}

func (s *StrJsonValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.val)
}

func (s *StrJsonValKind) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &s.val)
}

type BoolJsonValKind struct {
	val bool
}

func (b *BoolJsonValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.val)
}

func (b *BoolJsonValKind) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &b.val)
}

type NullJsonValKind struct {
}

func (n *NullJsonValKind) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (n *NullJsonValKind) UnmarshalJSON([]byte) error {
	return nil
}

func toJsonValKind(i interface{}) (JsonValKind, error) {
	switch v := i.(type) {
	case string:
		return &StrJsonValKind{v}, nil
	case bool:
		return &BoolJsonValKind{v}, nil
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return &NumberJsonValKind{json.Number(ToStr(v))}, nil
	case nil:
		return &NullJsonValKind{}, nil
	}
	return nil, fmt.Errorf("%v is not json", i)
}
