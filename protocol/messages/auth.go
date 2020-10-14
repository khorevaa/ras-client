package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"io"
)

// ClusterAuthenticateRequest установка авторизации на кластере
//
//  type AUTHENTICATE_REQUEST = 10
//  kind MESSAGE_KIND = 1
//  respond nothing
type ClusterAuthenticateRequest struct {
	ClusterID      uuid.UUID
	User, Password string
}

func (r ClusterAuthenticateRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.String(r.User, w)
	encoder.String(r.Password, w)
}

func (_ ClusterAuthenticateRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ ClusterAuthenticateRequest) ResponseMessage() types.EndpointResponseMessage {
	return nullEndpointResponse()
}

func (_ ClusterAuthenticateRequest) Type() types.Typed {
	return AUTHENTICATE_REQUEST
}

// AuthenticateAgentRequest установка авторизации на агенте
//
//  type AUTHENTICATE_AGENT_REQUEST = 9
//  kind MESSAGE_KIND = 1
//  respond nothing
type AuthenticateAgentRequest struct {
	User, Password string
}

func (_ AuthenticateAgentRequest) ResponseMessage() types.EndpointResponseMessage {
	return nullEndpointResponse()
}

func (_ AuthenticateAgentRequest) Type() types.Typed {
	return AUTHENTICATE_AGENT_REQUEST
}

func (r AuthenticateAgentRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {

	encoder.String(r.User, w)
	encoder.String(r.Password, w)

}

// AuthenticateInfobaseRequest установка авторизации в информационной базе
//
//  type ADD_AUTHENTICATION_REQUEST = 11
//  kind MESSAGE_KIND = 1
//  respond nothing
type AuthenticateInfobaseRequest struct {
	ClusterID      uuid.UUID
	User, Password string
}

func (_ AuthenticateInfobaseRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ AuthenticateInfobaseRequest) ResponseMessage() types.EndpointResponseMessage {
	return nullEndpointResponse()
}

func (_ AuthenticateInfobaseRequest) Type() types.Typed {
	return ADD_AUTHENTICATION_REQUEST
}

func (r AuthenticateInfobaseRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)
	encoder.String(r.User, w)
	encoder.String(r.Password, w)

}
