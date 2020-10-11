package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
	"io"
)

// GetClusterServicesRequest получение списка сервисов кластера
//
//  type GET_CLUSTER_SERVICES_REQUEST = 38
//  kind MESSAGE_KIND = 1
//  respond GetClusterServicesResponse
type GetClusterServicesRequest struct {
	ClusterID uuid.UUID
	response  *GetClusterServicesResponse
}

func (r *GetClusterServicesRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ClusterID, w)
}

func (_ *GetClusterServicesRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetClusterServicesRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetClusterServicesResponse{}
	}

	return r.response
}

func (r *GetClusterServicesRequest) Response() *GetClusterServicesResponse {

	return r.response
}

func (_ *GetClusterServicesRequest) Type() types.Typed {
	return GET_CLUSTER_SERVICES_REQUEST
}

// GetClusterServicesResponse содержит список сервисов кластера
//  type GET_CLUSTER_SERVICES_RESPONSE = 37
//  Services serialize.ManagerInfo
type GetClusterServicesResponse struct {
	Services []*serialize.ServiceInfo
}

func (res *GetClusterServicesResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.ServiceInfo{}
		info.Parse(decoder, version, r)
		res.Services = append(res.Services, info)
	}

}

func (_ *GetClusterServicesResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetClusterServicesResponse) Type() types.Typed {
	return GET_CLUSTER_SERVICES_RESPONSE
}
