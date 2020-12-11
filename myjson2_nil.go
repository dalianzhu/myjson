package myjson

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

type NilJson struct {
}

func (n *NilJson) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (n *NilJson) UnmarshalJSON(bytes []byte) error {
	return nil
}

func (n *NilJson) GetValue() interface{} {
	panic("implement me")
}

func (n *NilJson) Get(key string) MyJson2 {
	return new(NilJson)
}

func (n *NilJson) Set(key string, val interface{}) error {
	return nil
}

func (n *NilJson) PbValue() *structpb.Value {
	return nil
}

func (n *NilJson) Rm(key string) {
}

func (n *NilJson) Index(i int) MyJson2 {
	return new(NilJson)
}

func (n *NilJson) Insert(i int, val interface{}) (MyJson2, error) {
	return new(NilJson), nil
}

func (n *NilJson) Append(val interface{}) (MyJson2, error) {
	return new(NilJson), nil
}

func (n *NilJson) Len() int {
	return 0
}

func (n *NilJson) String() string {
	return ""
}

func (n *NilJson) Bytes() []byte {
	return []byte("")
}

func (n *NilJson) Int() (int, error) {
	return 0, fmt.Errorf("json is nil")
}

func (n *NilJson) Float64() (float64, error) {
	return 0, fmt.Errorf("json is nil")
}

func (n *NilJson) Bool() (bool, error) {
	return false, fmt.Errorf("json is nil")
}

func (n *NilJson) Clone() MyJson2 {
	return &NilJson{}
}

func (n *NilJson) RangeSlice(f func(index int, val MyJson2) (bool, error)) error {
	return fmt.Errorf("json is nil")
}

func (n *NilJson) RangeMap(f func(key string, val MyJson2) (bool, error)) error {
	return fmt.Errorf("json is nil")
}

func (n *NilJson) IsErrOrNil() bool {
	return true
}

func (n *NilJson) IsSlice() bool {
	return false
}

func (n *NilJson) IsMap() bool {
	return false
}
