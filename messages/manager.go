package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/serialize"
	"github.com/v8platform/rac/serialize/esig"
	"github.com/v8platform/rac/types"
	"io"
)

// GetClusterManagersRequest получение списка менеджеров кластера
//
//  type GET_CLUSTER_MANAGERS_REQUEST = 19
//  kind MESSAGE_KIND = 1
//  respond GetClusterManagersResponse
type GetClusterManagersRequest struct {
	ClusterID uuid.UUID
	response  *GetClusterManagersResponse
}

func (r *GetClusterManagersRequest) Sig() esig.ESIG {
	return esig.FromUuid(r.ClusterID)
}

func (_ *GetClusterManagersRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetClusterManagersRequest) ResponseMessage() types.EndpointResponseMessage {

	if r.response == nil {
		r.response = &GetClusterManagersResponse{}
	}

	return r.response
}

func (_ *GetClusterManagersRequest) Type() types.Typed {
	return GET_CLUSTER_MANAGERS_REQUEST
}

func (r *GetClusterManagersRequest) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(r.ClusterID, w)

}

func (r *GetClusterManagersRequest) Response() *GetClusterManagersResponse {
	return r.response
}

// GetClusterManagersResponse содержит список менеджеров кластера
//  type GET_CLUSTER_MANAGERS_RESPONSE = 20
//  Managers serialize.ManagerInfo
type GetClusterManagersResponse struct {
	Managers []*serialize.ManagerInfo
}

func (res *GetClusterManagersResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.ManagerInfo{}
		info.Parse(decoder, version, r)

		res.Managers = append(res.Managers, info)
	}
}

func (_ *GetClusterManagersResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetClusterManagersResponse) Type() types.Typed {
	return GET_CLUSTER_MANAGERS_RESPONSE
}
