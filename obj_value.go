package myjson

import "fmt"

type SliceValue struct {
	sliceData []interface{}
}

func (s *SliceValue) GetValue() []interface{} {
	return s.sliceData
}

func (s *SliceValue) String() string {
	return fmt.Sprintf("%v", s.sliceData)
}

func (s *SliceValue) MarshalJSON() ([]byte, error) {
	Debugf("sliceWrap MarshalJson:")
	return jsonit.Marshal(s.sliceData)
}

type NullValue struct {
}

var globalNullWrap = &NullValue{}

func GetJsonValNull() *NullValue {
	return globalNullWrap
}

var bytesNull = []byte("null")

func (s *NullValue) MarshalJSON() ([]byte, error) {
	// Debugf("nullWrap MarshalJson:")
	return bytesNull, nil
}

func (s *NullValue) String() string {
	return "null"
}

func NewErrorJson(err error) MyJson2 {
	return &ValueJson{data: NewErrorValue(err)}
}

func NewErrorValue(err error) *ErrorValue {
	return &ErrorValue{Err: err}
}

type ErrorValue struct {
	Err error
}

func (e *ErrorValue) String() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}
