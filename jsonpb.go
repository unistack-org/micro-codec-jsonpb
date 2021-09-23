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
	DefaultMarshalOptions = jsonpb.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		UseProtoNames:   true,
		AllowPartial:    false,
	}

	DefaultUnmarshalOptions = jsonpb.UnmarshalOptions{
		DiscardUnknown: false,
		AllowPartial:   false,
	}
)

type jsonpbCodec struct {
	opts codec.Options
}

const (
	flattenTag = "flatten"
)

func (c *jsonpbCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	if m, ok := v.(*codec.Frame); ok {
		return m.Data, nil
	}

	if _, ok := v.(proto.Message); !ok {
		return nil, codec.ErrInvalidMessage
	}

	marshalOptions := DefaultMarshalOptions
	if options.Context != nil {
		if f, ok := options.Context.Value(marshalOptionsKey{}).(jsonpb.MarshalOptions); ok {
			marshalOptions = f
		}
	}

	return marshalOptions.Marshal(v.(proto.Message))
}

func (c *jsonpbCodec) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	if m, ok := v.(*codec.Frame); ok {
		m.Data = d
		return nil
	}

	if _, ok := v.(proto.Message); !ok {
		return codec.ErrInvalidMessage
	}

	unmarshalOptions := DefaultUnmarshalOptions
	if options.Context != nil {
		if f, ok := options.Context.Value(unmarshalOptionsKey{}).(jsonpb.UnmarshalOptions); ok {
			unmarshalOptions = f
		}
	}

	return unmarshalOptions.Unmarshal(d, v.(proto.Message))
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

func NewCodec(opts ...codec.Option) codec.Codec {
	return &jsonpbCodec{opts: codec.NewOptions(opts...)}
}
