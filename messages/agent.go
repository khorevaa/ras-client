package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize/esig"
	"io"
)

// GetAgentVersionRequest получение версии агента
//
//  type GET_AGENT_VERSION_REQUEST
//  respond GetAgentAdminsResponse
type GetAgentVersionRequest struct{}

func (r *GetAgentVersionRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (r *GetAgentVersionRequest) Format(_ codec.Encoder, _ int, _ io.Writer) {}

func (_ *GetAgentVersionRequest) Type() EndpointMessageType {
	return GET_AGENT_VERSION_REQUEST
}

// GetAgentVersionResponse ответ с версией агента кластера
//
//  type GET_AGENT_VERSION_RESPONSE
//  Users serialize.UsersList
type GetAgentVersionResponse struct {
	Version string
}

func (res *GetAgentVersionResponse) Parse(decoder codec.Decoder, _ int, r io.Reader) {

	decoder.StringPtr(&res.Version, r)
}

func (_ *GetAgentVersionResponse) Type() EndpointMessageType {
	return GET_AGENT_VERSION_RESPONSE
}
