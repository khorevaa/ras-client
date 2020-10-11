package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
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
		r.response = &GetConnectionsShortResponse{}
	}

	return r.response
}

func (_ *GetConnectionsShortRequest) Type() types.Typed {
	return GET_CONNECTIONS_SHORT_REQUEST
}

func (r *GetConnectionsShortRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
}

func (r *GetConnectionsShortRequest) Response() *GetConnectionsShortResponse {
	return r.response
}

// GetConnectionsShortResponse ответ со списком соединений кластера
//
//  type GET_CONNECTIONS_SHORT_RESPONSE = 52
//  kind MESSAGE_KIND = 1
//  respond serialize.ConnectionInfoList
type GetConnectionsShortResponse struct {
	Connections serialize.ConnectionInfoList
}

func (_ *GetConnectionsShortResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetConnectionsShortResponse) Type() types.Typed {
	return GET_CONNECTIONS_SHORT_RESPONSE
}

func (res *GetConnectionsShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.ConnectionInfoList{}
	list.Parse(decoder, version, r)

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

func (r *DisconnectConnectionRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ProcessID, w)
	encoder.Uuid(r.ConnectionID, w)
}
