// Package jsonpb provides a json codec
package jsonpb

import (
	"bytes"
	"io"
	"io/ioutil"

	// nolint: staticcheck
	oldjsonpb "github.com/golang/protobuf/jsonpb"
	// nolint: staticcheck
	oldproto "github.com/golang/protobuf/proto"
	"github.com/unistack-org/micro/v3/codec"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	JsonpbMarshaler = jsonpb.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		UseProtoNames:   true,
		AllowPartial:    false,
	}

	JsonpbUnmarshaler = jsonpb.UnmarshalOptions{
		DiscardUnknown: false,
		AllowPartial:   false,
	}

	OldJsonpbMarshaler = oldjsonpb.Marshaler{
		OrigName:     true,
		EmitDefaults: false,
	}

	OldJsonpbUnmarshaler = oldjsonpb.Unmarshaler{
		AllowUnknownFields: false,
	}
)

type jsonpbCodec struct{}

func (c *jsonpbCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	case proto.Message:
		return JsonpbMarshaler.Marshal(m)
	case oldproto.Message:
		buf, err := OldJsonpbMarshaler.MarshalToString(m)
		return []byte(buf), err
	}
	return nil, codec.ErrInvalidMessage
}

func (c *jsonpbCodec) Unmarshal(d []byte, v interface{}) error {
	if d == nil {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = d
		return nil
	case proto.Message:
		return JsonpbUnmarshaler.Unmarshal(d, m)
	case oldproto.Message:
		return OldJsonpbUnmarshaler.Unmarshal(bytes.NewReader(d), m)
	}
	return codec.ErrInvalidMessage
}
func (c *jsonpbCodec) ReadHeader(conn io.ReadWriter, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *jsonpbCodec) ReadBody(conn io.ReadWriter, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		m.Data = buf
		return nil
	case oldproto.Message:
		return OldJsonpbUnmarshaler.Unmarshal(conn, m)
	case proto.Message:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		return JsonpbUnmarshaler.Unmarshal(buf, m)
	}
	return codec.ErrInvalidMessage
}

func (c *jsonpbCodec) Write(conn io.ReadWriter, m *codec.Message, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	case oldproto.Message:
		return OldJsonpbMarshaler.Marshal(conn, m)
	case proto.Message:
		buf, err := JsonpbMarshaler.Marshal(m)
		if err != nil {
			return err
		}
		_, err = conn.Write(buf)
		return err
	}
	return codec.ErrInvalidMessage
}

func (c *jsonpbCodec) String() string {
	return "jsonpb"
}

func NewCodec() codec.Codec {
	return &jsonpbCodec{}
}
