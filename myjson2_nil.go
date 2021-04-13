package myjson

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

// NilOrErrJson 用来做占位，不等于null，而是真正的空内存
type NilOrErrJson struct {
}

func (n *NilOrErrJson) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("json is nil pointer")
}

func (n *NilOrErrJson) UnmarshalJSON(bytes []byte) error {
	return nil
}

func (n *NilOrErrJson) GetValue() interface{} {
	return nil
}

func (n *NilOrErrJson) Get(key string) MyJson2 {
	return new(NilOrErrJson)
}

func (n *NilOrErrJson) Set(key string, val interface{}) error {
	return fmt.Errorf("nil cannot set")
}

func (n *NilOrErrJson) PbValue() *structpb.Value {
	return nil
}

func (n *NilOrErrJson) Rm(key string) {
}

func (n *NilOrErrJson) Index(i int) MyJson2 {
	return new(NilOrErrJson)
}

func (n *NilOrErrJson) Insert(i int, val interface{}) error {
	return fmt.Errorf("nil cannot insert")
}

func (n *NilOrErrJson) Append(val interface{}) error {
	return fmt.Errorf("nil cannot append")
}

func (n *NilOrErrJson) Len() int {
	return 0
}

func (n *NilOrErrJson) String() string {
	return ""
}

func (n *NilOrErrJson) Bytes() []byte {
	return []byte("")
}

func (n *NilOrErrJson) Int() (int, error) {
	return 0, fmt.Errorf("json is nil")
}

func (n *NilOrErrJson) Float64() (float64, error) {
	return 0, fmt.Errorf("json is nil")
}

func (n *NilOrErrJson) Bool() (bool, error) {
	return false, fmt.Errorf("json is nil")
}

func (n *NilOrErrJson) Clone() MyJson2 {
	return &NilOrErrJson{}
}

func (n *NilOrErrJson) RangeSlice(f func(index int, val MyJson2) (bool, error)) error {
	return fmt.Errorf("json is nil")
}

func (n *NilOrErrJson) RangeMap(f func(key string, val MyJson2) (bool, error)) error {
	return fmt.Errorf("json is nil")
}

func (n *NilOrErrJson) IsErrOrNil() bool {
	return true
}

func (n *NilOrErrJson) IsSlice() bool {
	return false
}

func (n *NilOrErrJson) IsMap() bool {
	return false
}

func (n *NilOrErrJson) IsNull() bool {
	return false
}

func (n *NilOrErrJson) Keys() ([]string, error) {
	return nil, fmt.Errorf("Nil json has no keys")
}

func (n *NilOrErrJson) Items() ([]interface{}, error) {
	return nil, fmt.Errorf("nil json has no items")
}
