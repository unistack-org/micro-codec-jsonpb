// Package jsonpb provides a json protobuf codec
package jsonpb

import (
	pb "go.unistack.org/micro-proto/v4/codec"
	"go.unistack.org/micro/v4/codec"
	rutil "go.unistack.org/micro/v4/util/reflect"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	DefaultMarshalOptions = jsonpb.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		UseProtoNames:   false,
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

type jsonpbCodecV2 struct {
	opts codec.Options
}

var (
	_ codec.Codec   = (*jsonpbCodec)(nil)
	_ codec.CodecV2 = (*jsonpbCodecV2)(nil)
)

func (c *jsonpbCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case proto.Message:
		marshalOptions := DefaultMarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(jsonpb.MarshalOptions); ok {
				marshalOptions = f
			}
		}
		return marshalOptions.Marshal(m)
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	default:
		return nil, codec.ErrInvalidMessage
	}
}

func (c *jsonpbCodec) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	case proto.Message:
		unmarshalOptions := DefaultUnmarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(unmarshalOptionsKey{}).(jsonpb.UnmarshalOptions); ok {
				unmarshalOptions = f
			}
		}
		return unmarshalOptions.Unmarshal(d, m)
	case *codec.RawMessage:
		*m = append((*m)[0:0], d...)
		return nil
	case codec.RawMessage:
		copy(m, d)
		return nil
	default:
		return codec.ErrInvalidMessage
	}
}

func (c *jsonpbCodec) String() string {
	return "jsonpb"
}

func NewCodec(opts ...codec.Option) codec.Codec {
	return &jsonpbCodec{opts: codec.NewOptions(opts...)}
}

func (c *jsonpbCodecV2) Marshal(d []byte, v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case proto.Message:
		marshalOptions := DefaultMarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(marshalOptionsKey{}).(jsonpb.MarshalOptions); ok {
				marshalOptions = f
			}
		}
		return marshalOptions.MarshalAppend(d[:0], m)
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	default:
		return nil, codec.ErrInvalidMessage
	}
}

func (c *jsonpbCodecV2) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	case proto.Message:
		unmarshalOptions := DefaultUnmarshalOptions
		if options.Context != nil {
			if f, ok := options.Context.Value(unmarshalOptionsKey{}).(jsonpb.UnmarshalOptions); ok {
				unmarshalOptions = f
			}
		}
		return unmarshalOptions.Unmarshal(d, m)
	case *codec.RawMessage:
		*m = append((*m)[0:0], d...)
		return nil
	case codec.RawMessage:
		copy(m, d)
		return nil
	default:
		return codec.ErrInvalidMessage
	}
}

func (c *jsonpbCodecV2) String() string {
	return "jsonpb"
}

func NewCodecV2(opts ...codec.Option) codec.CodecV2 {
	return &jsonpbCodecV2{opts: codec.NewOptions(opts...)}
}
