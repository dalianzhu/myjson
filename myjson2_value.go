package myjson

import (
	"bytes"
	"errors"
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

func (s *sliceWrap) GetValue() []interface{} {
	return s.sliceData
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
