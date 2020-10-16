package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/esig"
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

func (r GetInfobasesShortRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (r GetInfobasesShortRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ GetInfobasesShortRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetInfobasesShortRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetInfobasesShortResponse{}
	}
	r.response.ClusterID = r.ClusterID

	return r.response
}

func (_ GetInfobasesShortRequest) Type() types.Typed {
	return GET_INFOBASES_SHORT_REQUEST
}

func (r *GetInfobasesShortRequest) Response() *GetInfobasesShortResponse {
	return r.response
}

// GetInfobasesShortResponse
// type GET_INFOBASES_SHORT_RESPONSE = 44
type GetInfobasesShortResponse struct {
	Infobases serialize.InfobaseSummaryList
	ClusterID uuid.UUID
}

func (res *GetInfobasesShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	list := serialize.InfobaseSummaryList{}

	list.Parse(decoder, version, r)

	list.Each(func(info *serialize.InfobaseSummaryInfo) {
		info.Cluster = res.ClusterID
	})

	res.Infobases = list
	//
	//count := decoder.Size(r)
	//
	//for i := 0; i < count; i++ {
	//
	//	info := &serialize.InfobaseSummaryInfo{}
	//	decoder.UuidPtr(&info.UUID, r)
	//	decoder.StringPtr(&info.Description, r)
	//	decoder.StringPtr(&info.Name, r)
	//
	//	res.Infobases = append(res.Infobases, info)
	//}

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
	ClusterID uuid.UUID
	Infobase  *serialize.InfobaseInfo
	response  *CreateInfobaseResponse
	Mode      int // Mode 1 - создавать базу на сервере, 0 - не создавать

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
	encoder.Uuid(r.ClusterID, w)

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
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
}

func (_ *CreateInfobaseResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *CreateInfobaseResponse) Type() types.Typed {
	return CREATE_INFOBASE_RESPONSE
}

func (res *CreateInfobaseResponse) Parse(decoder codec.Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&res.InfobaseID, r)

}

// GetInfobaseInfoRequest запрос получение информации по информационной базе
//
//  type GET_INFOBASE_INFO_REQUEST = 49
//  kind MESSAGE_KIND = 1
//  respond GetInfobaseInfoResponse
type GetInfobaseInfoRequest struct {
	ClusterID  uuid.UUID
	InfobaseID uuid.UUID
	response   *GetInfobaseInfoResponse
}

func (_ *GetInfobaseInfoRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetInfobaseInfoRequest) ResponseMessage() types.EndpointResponseMessage {

	if r.response == nil {
		r.response = &GetInfobaseInfoResponse{}
	}

	r.response.ClusterID = r.ClusterID

	return r.response
}

func (_ *GetInfobaseInfoRequest) Type() types.Typed {
	return GET_INFOBASE_INFO_REQUEST
}

func (r *GetInfobaseInfoRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseID, w)
}

func (r *GetInfobaseInfoRequest) Response() *GetInfobaseInfoResponse {
	return r.response
}

// GetInfobaseInfoResponse ответ с информацией о информационной базы
//  type GET_INFOBASE_INFO_RESPONSE = 50
//  return serialize.InfobaseInfo
type GetInfobaseInfoResponse struct {
	ClusterID uuid.UUID
	Infobase  serialize.InfobaseInfo
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
	info.ClusterID = res.ClusterID

	res.Infobase = *info

}

// DropInfobaseRequest запрос удаление информационной базы
//
//  type DROP_INFOBASE_REQUEST = 42
//  kind MESSAGE_KIND = 1
//  respond nothing
type DropInfobaseRequest struct {
	ClusterID  uuid.UUID
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

func (r *DropInfobaseRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
	encoder.Uuid(r.InfobaseId, w)
	encoder.Int(r.Mode, w)
}

// UpdateInfobaseRequest запрос обновление данных по информационной базы
//
//  type UPDATE_INFOBASE_REQUEST = 40
//  kind MESSAGE_KIND = 1
//  respond nothing
type UpdateInfobaseRequest struct {
	ClusterID uuid.UUID
	Infobase  serialize.InfobaseInfo
}

func (_ *UpdateInfobaseRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *UpdateInfobaseRequest) ResponseMessage() types.EndpointResponseMessage {

	return nullEndpointResponse()
}

func (_ *UpdateInfobaseRequest) Type() types.Typed {
	return UPDATE_INFOBASE_REQUEST
}

func (r *UpdateInfobaseRequest) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)
	r.Infobase.Format(encoder, version, w)

}

// UpdateInfobaseShortRequest запрос обновление данных по информационной базы
//
//  type UPDATE_INFOBASE_REQUEST = 40
//  kind MESSAGE_KIND = 1
//  respond nothing
type UpdateInfobaseShortRequest struct {
	ClusterID uuid.UUID
	Infobase  serialize.InfobaseSummaryInfo
}

func (_ *UpdateInfobaseShortRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *UpdateInfobaseShortRequest) ResponseMessage() types.EndpointResponseMessage {

	return nullEndpointResponse()
}

func (_ *UpdateInfobaseShortRequest) Type() types.Typed {
	return UPDATE_INFOBASE_SHORT_REQUEST
}

func (r *UpdateInfobaseShortRequest) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)
	r.Infobase.Format(encoder, version, w)

}
