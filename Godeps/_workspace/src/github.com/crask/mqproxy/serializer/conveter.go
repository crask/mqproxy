package serializer

import (
	"bytes"
	"encoding/json"
	"errors"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Serializer struct {
	Converter string
}

func (s *Serializer) Unmarshal(data []byte, v interface{}) error {
	if s.Converter == "json" {
		return unmarshalJson(data, v)
	} else if s.Converter == "msgpack" {
		return unmarshalMsgpack(data, v)
	}

	return errors.New("unkonwn converter")
}

func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	if s.Converter == "json" {
		return marshalJson(v)
	} else if s.Converter == "msgpack" {
		return marshalMsgpack(v)
	}

	buf := &bytes.Buffer{}
	return buf.Bytes(), errors.New("unkonwn converter")
}

func marshalJson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func marshalMsgpack(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func unmarshalJson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func unmarshalMsgpack(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
