package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
	"io"
)

// GetInfobasesShortRequest получение списка инфорамационных баз кластера
//
//  type GET_INFOBASES_SHORT_REQUEST = 43
//  kind MESSAGE_KIND = 1
//  respond GetInfobasesShortResponse
type GetInfobasesShortRequest struct {
	ClusterID uuid.UUID
	response  *GetInfobasesShortResponse
}

func (r GetInfobasesShortRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ GetInfobasesShortRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r GetInfobasesShortRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetInfobasesShortResponse{}
	}

	return r.response
}

func (_ GetInfobasesShortRequest) Type() types.Typed {
	return GET_INFOBASES_SHORT_REQUEST
}

func (r GetInfobasesShortRequest) Response() *GetInfobasesShortResponse {
	return r.response
}

// GetInfobasesShortResponse
// type GET_INFOBASES_SHORT_RESPONSE = 44
type GetInfobasesShortResponse struct {
	Infobases []*serialize.InfobaseSummaryInfo
}

func (res *GetInfobasesShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.InfobaseSummaryInfo{}
		decoder.UuidPtr(&info.UUID, r)
		decoder.StringPtr(&info.Description, r)
		decoder.StringPtr(&info.Name, r)

		res.Infobases = append(res.Infobases, info)
	}

}

func (_ *GetInfobasesShortResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetInfobasesShortResponse) Type() types.Typed {
	return GET_INFOBASES_SHORT_RESPONSE
}

// CreateInfobaseRequest запрос на создание новой базы
//
//  type CREATE_INFOBASE_REQUEST = 38
//  kind MESSAGE_KIND = 1
//  respond CreateInfobaseResponse
type CreateInfobaseRequest struct {
	ID       uuid.UUID
	Infobase *serialize.InfobaseInfo
	response *CreateInfobaseResponse
	Mode     int // Mode 1 - создавать базу на сервере, 0 - не создавать

}

func (_ *CreateInfobaseRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *CreateInfobaseRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &CreateInfobaseResponse{}
	}

	return r.response
}

func (_ *CreateInfobaseRequest) Type() types.Typed {
	return CREATE_INFOBASE_REQUEST
}

func (r *CreateInfobaseRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)

	r.Infobase.Format(encoder, version, w)

	encoder.Int(r.Mode, w)
}

func (r *CreateInfobaseRequest) Response() *CreateInfobaseResponse {
	return r.response
}

// CreateInfobaseResponse ответ создания новой информационной базы
//  type CREATE_INFOBASE_RESPONSE = 39
//  return uuid.UUID созданной базы
type CreateInfobaseResponse struct {
	ID uuid.UUID
}

func (_ *CreateInfobaseResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *CreateInfobaseResponse) Type() types.Typed {
	return CREATE_INFOBASE_RESPONSE
}

func (res *CreateInfobaseResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&res.ID, r)

}

// GetInfobaseInfoRequest запрос получение информации по информационной базе
//
//  type GET_INFOBASE_INFO_REQUEST = 49
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseInfoResponse
type GetInfobaseInfoRequest struct {
	ID         uuid.UUID
	InfobaseId uuid.UUID
	response   *GetInfobaseInfoResponse
}

func (_ *GetInfobaseInfoRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetInfobaseInfoRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetInfobaseInfoResponse{}
	}

	return r.response
}

func (_ *GetInfobaseInfoRequest) Type() types.Typed {
	return GET_INFOBASE_INFO_REQUEST
}

func (r *GetInfobaseInfoRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
	encoder.Uuid(r.InfobaseId, w)
}

func (r *GetInfobaseInfoRequest) Response() *GetInfobaseInfoResponse {
	return r.response
}

// GetInfobaseInfoResponse ответ с информацией о информационной базы
//  type GET_INFOBASE_INFO_RESPONSE = 50
//  return serialize.InfobaseInfo
type GetInfobaseInfoResponse struct {
	infobase *serialize.InfobaseInfo
}

func (_ *GetInfobaseInfoResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetInfobaseInfoResponse) Type() types.Typed {
	return GET_INFOBASE_INFO_RESPONSE
}

func (res *GetInfobaseInfoResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	info := &serialize.InfobaseInfo{}
	info.Parse(decoder, version, r)
	res.infobase = info

}

// DropInfobaseRequest запрос удаление информационной базы
//
//  type DROP_INFOBASE_REQUEST = 42
//  kind MESSAGE_KIND = 1
//  respond nothing
type DropInfobaseRequest struct {
	ID         uuid.UUID
	InfobaseId uuid.UUID
	Mode       int
}

func (_ *DropInfobaseRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *DropInfobaseRequest) ResponseMessage() types.EndpointResponseMessage {

	return nullEndpointResponse()
}

func (_ *DropInfobaseRequest) Type() types.Typed {
	return DROP_INFOBASE_REQUEST
}

func (r *DropInfobaseRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
	encoder.Uuid(r.InfobaseId, w)
	encoder.Int(r.Mode, w)
}
