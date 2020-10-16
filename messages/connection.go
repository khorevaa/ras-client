package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/serialize"
	"github.com/v8platform/rac/serialize/esig"
	"github.com/v8platform/rac/types"
	"io"
)

// GetConnectionsShortRequest получение списка соединений кластера
//
//  type GET_CONNECTIONS_SHORT_REQUEST = 51
//  kind MESSAGE_KIND = 1
//  respond GetConnectionsShortResponse
type GetConnectionsShortRequest struct {
	ID       uuid.UUID
	response *GetConnectionsShortResponse
}

func (_ *GetConnectionsShortRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetConnectionsShortRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetConnectionsShortResponse{
			ClusterID: r.ID,
		}
	}

	return r.response
}

func (_ *GetConnectionsShortRequest) Type() types.Typed {
	return GET_CONNECTIONS_SHORT_REQUEST
}

func (r *GetConnectionsShortRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ID, w)
}

func (r *GetConnectionsShortRequest) Response() *GetConnectionsShortResponse {
	return r.response
}

// GetConnectionsShortResponse ответ со списком соединений кластера
//
//  type GET_CONNECTIONS_SHORT_RESPONSE = 52
//  kind MESSAGE_KIND = 1
//  respond serialize.ConnectionShortInfoList
type GetConnectionsShortResponse struct {
	ClusterID   uuid.UUID
	Connections serialize.ConnectionShortInfoList
}

func (_ *GetConnectionsShortResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetConnectionsShortResponse) Type() types.Typed {
	return GET_CONNECTIONS_SHORT_RESPONSE
}

func (res *GetConnectionsShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.ConnectionShortInfoList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.ConnectionShortInfo) {
		info.ClusterID = res.ClusterID
	})

	res.Connections = list

}

// DisconnectConnectionRequest отключение соединения
//
//  type DISCONNECT_REQUEST = 59
//  kind MESSAGE_KIND = 1
//  respond nothing
type DisconnectConnectionRequest struct {
	ClusterID    uuid.UUID
	ProcessID    uuid.UUID
	ConnectionID uuid.UUID
}

func (_ *DisconnectConnectionRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *DisconnectConnectionRequest) Type() types.Typed {
	return DISCONNECT_REQUEST
}

func (_ DisconnectConnectionRequest) ResponseMessage() types.EndpointResponseMessage {
	return nullEndpointResponse()
}

func (r *DisconnectConnectionRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ProcessID, w)
	encoder.Uuid(r.ConnectionID, w)
}

// GetInfobaseConnectionsShortRequest получение списка соединений кластера
//
//  type GET_INFOBASE_CONNECTIONS_SHORT_REQUEST = 52
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseConnectionsShortResponse
type GetInfobaseConnectionsShortRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	response   *GetInfobaseConnectionsShortResponse
}

func (r *GetInfobaseConnectionsShortRequest) Sig() esig.ESIG {
	return esig.From2Uuid(r.ClusterID, r.InfobaseID)
}

func (_ *GetInfobaseConnectionsShortRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetInfobaseConnectionsShortRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetInfobaseConnectionsShortResponse{}
	}

	r.response.ClusterID = r.ClusterID
	r.response.InfobaseID = r.InfobaseID

	return r.response
}

func (_ *GetInfobaseConnectionsShortRequest) Type() types.Typed {
	return GET_INFOBASE_CONNECTIONS_SHORT_REQUEST
}

func (r *GetInfobaseConnectionsShortRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

func (r *GetInfobaseConnectionsShortRequest) Response() *GetInfobaseConnectionsShortResponse {
	return r.response
}

// GetConnectionsShortResponse ответ со списком соединений кластера
//
//  type GET_INFOBASE_CONNECTIONS_SHORT_RESPONSE = 53
//  kind MESSAGE_KIND = 1
//  respond Connections serialize.ConnectionShortInfoList
type GetInfobaseConnectionsShortResponse struct {
	ClusterID   uuid.UUID
	InfobaseID  uuid.UUID
	Connections serialize.ConnectionShortInfoList
}

func (_ *GetInfobaseConnectionsShortResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetInfobaseConnectionsShortResponse) Type() types.Typed {
	return GET_INFOBASE_CONNECTIONS_SHORT_RESPONSE
}

func (res *GetInfobaseConnectionsShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.ConnectionShortInfoList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.ConnectionShortInfo) {
		info.ClusterID = res.ClusterID
		info.InfobaseID = res.InfobaseID
	})

	res.Connections = list

}
