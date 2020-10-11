package messages

import (
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
	"io"
)

// GetClustersRequest получение списка кластеров
//
//  type GET_CLUSTERS_REQUEST = 12
//  kind MESSAGE_KIND = 1
//  respond GetClustersResponse
type GetClustersRequest struct {
	response *GetClustersResponse
}

func (_ *GetClustersRequest) Format(encoder codec.Encoder, version int, w io.Writer) {}

func (_ *GetClustersRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetClustersRequest) ResponseMessage() types.EndpointResponseMessage {

	if r.response == nil {
		r.response = &GetClustersResponse{}
	}

	return r.response
}

func (r *GetClustersRequest) Response() *GetClustersResponse {
	return r.response
}

func (_ GetClustersRequest) Type() types.Typed {
	return GET_CLUSTERS_REQUEST
}

// GetClustersResponse ответ со списком кластеров
//
//  type GET_CLUSTERS_RESPONSE = 13
//  kind MESSAGE_KIND = 1
//  Clusters []*serialize.ClusterInfo
type GetClustersResponse struct {
	Clusters []*serialize.ClusterInfo
}

func (res *GetClustersResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &serialize.ClusterInfo{}
		info.Parse(decoder, version, r)
		res.Clusters = append(res.Clusters, info)
	}

}

func (_ *GetClustersResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetClustersResponse) Type() types.Typed {
	return GET_CLUSTERS_RESPONSE
}
