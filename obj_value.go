package myjson

type SliceValue struct {
	sliceData []interface{}
}

func (s *SliceValue) GetValue() []interface{} {
	return s.sliceData
}

func (s *SliceValue) MarshalJSON() ([]byte, error) {
	Debugf("sliceWrap MarshalJson:")
	return jsonit.Marshal(s.sliceData)
}

type NullValue struct {
}

var globalNullWrap = &NullValue{}

func GetJsonNull() *NullValue {
	return globalNullWrap
}

var bytesNull = []byte("null")

func (s *NullValue) MarshalJSON() ([]byte, error) {
	// Debugf("nullWrap MarshalJson:")
	return bytesNull, nil
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
