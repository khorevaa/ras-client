package types

import (
	"github.com/v8platform/rac/protocol/codec"
	"io"
)

type Endpoint interface {
	Version() int

	SendMessage(req EndpointRequestMessage) (interface{}, error)
	Close()
}

type RequestMessage interface {
	Type() Typed
	Format(codec codec.Encoder, w io.Writer)
	ResponseMessage() ResponseMessage
}

type ResponseMessage interface {
	Type() Typed
	Parse(codec codec.Decoder, r io.Reader)
}

type EndpointRequestMessage interface {
	Type() Typed
	Format(encoder codec.Encoder, version int, w io.Writer)
	Kind() Typed
	ResponseMessage() EndpointResponseMessage
}

type EndpointResponseMessage interface {
	Type() Typed
	Parse(decoder codec.Decoder, version int, r io.Reader)
	Kind() Typed
}

type Typed interface {
	Type() int
}
