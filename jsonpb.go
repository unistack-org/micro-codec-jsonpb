// Package jsonpb provides a json codec
package jsonpb

import (
	"io"

	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
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
)

type jsonpbCodec struct{}

const (
	flattenTag = "flatten"
)

func (c *jsonpbCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	case proto.Message:
		if nv, nerr := rutil.StructFieldByTag(m, codec.DefaultTagName, flattenTag); nerr == nil {
			if nm, ok := nv.(proto.Message); ok {
				m = nm
			}
		}
		return JsonpbMarshaler.Marshal(m)
	}
	return nil, codec.ErrInvalidMessage
}

func (c *jsonpbCodec) Unmarshal(d []byte, v interface{}) error {
	if len(d) == 0 {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = d
		return nil
	case proto.Message:
		if nv, nerr := rutil.StructFieldByTag(m, codec.DefaultTagName, flattenTag); nerr == nil {
			if nm, ok := nv.(proto.Message); ok {
				m = nm
			}
		}
		return JsonpbUnmarshaler.Unmarshal(d, m)
	}
	return codec.ErrInvalidMessage
}

func (c *jsonpbCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *jsonpbCodec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}
	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}
	return c.Unmarshal(buf, v)
}

func (c *jsonpbCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return codec.ErrInvalidMessage
	}
	_, err = conn.Write(buf)
	return err
}

func (c *jsonpbCodec) String() string {
	return "jsonpb"
}

func NewCodec() codec.Codec {
	return &jsonpbCodec{}
}
