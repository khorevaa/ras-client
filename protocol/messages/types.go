package messages

import (
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"io"
)

type EndpointMessageKind int

func (e EndpointMessageKind) Type() int {
	return int(e)
}

const (
	VOID_MESSAGE_KIND EndpointMessageKind = 0
	MESSAGE_KIND      EndpointMessageKind = 1
	EXCEPTION_KIND    EndpointMessageKind = 0xff
)

var _ types.EndpointResponseMessage = (*nEndpointResponse)(nil)

func nullEndpointResponse() *nEndpointResponse {
	return &nEndpointResponse{}
}

type nEndpointResponse struct{}

func (n *nEndpointResponse) Type() types.Typed { return NULL_ENDPOINT_RESPONSE }

func (n *nEndpointResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {}

func (n *nEndpointResponse) Kind() types.Typed { return VOID_MESSAGE_KIND }
