package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/serialize"
	"github.com/v8platform/rac/types"
	"io"
)

// GetLocksRequest получение списка блокировок кластера
//
//  type GET_LOCKS_REQUEST = 66
//  kind MESSAGE_KIND = 1
//  respond GetSessionsResponse
type GetLocksRequest struct {
	ClusterID uuid.UUID
	response  *GetLocksResponse
}

func (_ *GetLocksRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetLocksRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetLocksResponse{}
	}

	r.response.ClusterID = r.ClusterID

	return r.response
}

func (_ *GetLocksRequest) Type() types.Typed {
	return GET_LOCKS_REQUEST
}

func (r *GetLocksRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (r *GetLocksRequest) Response() *GetLocksResponse {
	return r.response
}

// GetLocksResponse ответ со списком блокировок кластера
//
//  type GET_LOCKS_RESPONSE = 67
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetLocksResponse struct {
	ClusterID uuid.UUID
	List      serialize.LocksList
}

func (_ *GetLocksResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetLocksResponse) Type() types.Typed {
	return GET_LOCKS_RESPONSE
}

func (res *GetLocksResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.LockInfo) {
		info.ClusterID = res.ClusterID
	})

	res.List = list

}

// GetInfobaseLockRequest получение списка блокировок информационной базы кластера
//
//  type GET_INFOBASE_LOCKS_REQUEST = 68
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseSessionsResponse
type GetInfobaseLockRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	response   *GetInfobaseLockResponse
}

func (_ *GetInfobaseLockRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetInfobaseLockRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetInfobaseLockResponse{}
	}

	r.response.ClusterID = r.ClusterID
	r.response.InfobaseID = r.InfobaseID

	return r.response
}

func (_ *GetInfobaseLockRequest) Type() types.Typed {
	return GET_INFOBASE_LOCKS_REQUEST
}

func (r *GetInfobaseLockRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

func (r *GetInfobaseLockRequest) Response() *GetInfobaseLockResponse {
	return r.response
}

// GetInfobaseLockResponse ответ со списком сблокировок иб
//
//  type GET_INFOBASE_LOCKS_RESPONSE = 69
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetInfobaseLockResponse struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	List       serialize.LocksList
}

func (_ *GetInfobaseLockResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetInfobaseLockResponse) Type() types.Typed {
	return GET_INFOBASE_LOCKS_RESPONSE
}

func (res *GetInfobaseLockResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.LockInfo) {
		info.ClusterID = res.ClusterID
		info.InfobaseID = res.InfobaseID
	})

	res.List = list

}

// GetSessionLockRequest получение списка блокировок сессии информационной базы кластера
//
//  type GET_SESSION_LOCKS_REQUEST = 72
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseSessionsResponse
type GetSessionLockRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	SessionID  uuid.UUID
	response   *GetSessionLockResponse
}

func (_ *GetSessionLockRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetSessionLockRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetSessionLockResponse{}
	}

	r.response.ClusterID = r.ClusterID
	r.response.InfobaseID = r.InfobaseID
	r.response.SessionID = r.SessionID

	return r.response
}

func (_ *GetSessionLockRequest) Type() types.Typed {
	return GET_SESSION_LOCKS_REQUEST
}

func (r *GetSessionLockRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
	encoder.Uuid(r.SessionID, w)
}

func (r *GetSessionLockRequest) Response() *GetSessionLockResponse {
	return r.response
}

// GetSessionLockResponse ответ со списком блокировок сессии иб
//
//  type GET_SESSION_LOCKS_RESPONSE = 73
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetSessionLockResponse struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	SessionID  uuid.UUID
	List       serialize.LocksList
}

func (_ *GetSessionLockResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetSessionLockResponse) Type() types.Typed {
	return GET_SESSION_LOCKS_RESPONSE
}

func (res *GetSessionLockResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.LockInfo) {
		info.ClusterID = res.ClusterID
		info.InfobaseID = res.InfobaseID
		//info.SessionID = res.SessionID
	})

	res.List = list

}

// GetSessionLockRequest получение списка блокировок сессии информационной базы кластера
//
//  type GET_CONNECTION_LOCKS_REQUEST = 70
//  kind MESSAGE_KIND = 1
//  respond GetConnectionLockResponse
type GetConnectionLockRequest struct {
	ClusterID    uuid.UUID
	ConnectionID uuid.UUID
	response     *GetConnectionLockResponse
}

func (_ *GetConnectionLockRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetConnectionLockRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetConnectionLockResponse{}
	}

	r.response.ClusterID = r.ClusterID
	r.response.ConnectionID = r.ConnectionID

	return r.response
}

func (_ *GetConnectionLockRequest) Type() types.Typed {
	return GET_CONNECTION_LOCKS_REQUEST
}

func (r *GetConnectionLockRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.ConnectionID, w)
}

func (r *GetConnectionLockRequest) Response() *GetConnectionLockResponse {
	return r.response
}

// GetSessionLockResponse ответ со списком блокировок сессии иб
//
//  type GET_CONNECTION_LOCKS_RESPONSE = 71
//  kind MESSAGE_KIND = 1
//  respond Sessions serialize.SessionInfoList
type GetConnectionLockResponse struct {
	ClusterID    uuid.UUID
	ConnectionID uuid.UUID
	List         serialize.LocksList
}

func (_ *GetConnectionLockResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetConnectionLockResponse) Type() types.Typed {
	return GET_CONNECTION_LOCKS_RESPONSE
}

func (res *GetConnectionLockResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.LocksList{}
	list.Parse(decoder, version, r)
	list.Each(func(info *serialize.LockInfo) {
		info.ClusterID = res.ClusterID
	})

	res.List = list

}
