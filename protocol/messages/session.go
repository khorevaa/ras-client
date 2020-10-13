package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
	"io"
)

// TerminateSessionRequest отключение сеанса
//
//  type DISCONNECT_REQUEST = 71
//  kind MESSAGE_KIND = 1
//  respond nothing
type TerminateSessionRequest struct {
	ClusterID uuid.UUID
	SessionID uuid.UUID
	Message   string
}

func (_ *TerminateSessionRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *TerminateSessionRequest) Type() types.Typed {
	return TERMINATE_SESSION_REQUEST
}

func (_ TerminateSessionRequest) ResponseMessage() types.EndpointResponseMessage {
	return nullEndpointResponse()
}

func (r *TerminateSessionRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.SessionID, w)
	encoder.String(r.Message, w)
}

// GetInfobaseSessionsRequest получение списка сессий информационной базы кластера
//
//  type GET_INFOBASE_SESSIONS_REQUEST = 61
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseSessionsResponse
type GetInfobaseSessionsRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	response   *GetInfobaseSessionsResponse
}

func (_ *GetInfobaseSessionsRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetInfobaseSessionsRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetInfobaseSessionsResponse{}
	}

	r.response.ClusterID = r.ClusterID
	r.response.InfobaseID = r.InfobaseID

	return r.response
}

func (_ *GetInfobaseSessionsRequest) Type() types.Typed {
	return GET_INFOBASE_SESSIONS_REQUEST
}

func (r *GetInfobaseSessionsRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

func (r *GetInfobaseSessionsRequest) Response() *GetInfobaseSessionsResponse {
	return r.response
}

// GetInfobaseSessionsResponse ответ со списком сессий кластера
//
//  type GET_INFOBASE_SESSIONS_RESPONSE = 62
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetInfobaseSessionsResponse struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	Sessions   serialize.SessionInfoList
}

func (_ *GetInfobaseSessionsResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetInfobaseSessionsResponse) Type() types.Typed {
	return GET_INFOBASE_SESSIONS_RESPONSE
}

func (res *GetInfobaseSessionsResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.SessionInfoList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.SessionInfo) {
		info.ClusterID = res.ClusterID
		info.InfobaseID = res.InfobaseID
	})

	res.Sessions = list

}

// GetSessionsRequest получение списка сессий кластера
//
//  type GET_SESSIONS_REQUEST = 59
//  kind MESSAGE_KIND = 1
//  respond GetSessionsResponse
type GetSessionsRequest struct {
	ClusterID uuid.UUID
	response  *GetSessionsResponse
}

func (_ *GetSessionsRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetSessionsRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetSessionsResponse{}
	}

	r.response.ClusterID = r.ClusterID

	return r.response
}

func (_ *GetSessionsRequest) Type() types.Typed {
	return GET_SESSIONS_REQUEST
}

func (r *GetSessionsRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (r *GetSessionsRequest) Response() *GetSessionsResponse {
	return r.response
}

// GetInfobaseSessionsResponse ответ со списком сессий кластера
//
//  type GET_SESSIONS_RESPONSE = 60
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetSessionsResponse struct {
	ClusterID uuid.UUID
	Sessions  serialize.SessionInfoList
}

func (_ *GetSessionsResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetSessionsResponse) Type() types.Typed {
	return GET_SESSIONS_RESPONSE
}

func (res *GetSessionsResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.SessionInfoList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.SessionInfo) {
		info.ClusterID = res.ClusterID
	})

	res.Sessions = list

}
