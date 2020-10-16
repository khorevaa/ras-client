package messages

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/serialize"
	"github.com/v8platform/rac/serialize/esig"
	"github.com/v8platform/rac/types"
	"io"
)

// GetClustersRequest получение списка кластеров
//
//  type GET_CLUSTERS_REQUEST = 11
//  kind MESSAGE_KIND = 1
//  respond GetClustersResponse
type GetClustersRequest struct {
	response *GetClustersResponse
}

func (_ *GetClustersRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (_ *GetClustersRequest) Format(_ codec.Encoder, _ int, _ io.Writer) {}

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
//  type GET_CLUSTERS_RESPONSE = 12
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

// GetClusterInfoRequest получение информации о кластере
//
//  type GET_CLUSTER_INFO_REQUEST = 13
//  kind MESSAGE_KIND = 1
//  respond GetClustersResponse
type GetClusterInfoRequest struct {
	ID       uuid.UUID
	response *GetClusterInfoResponse
}

func (r *GetClusterInfoRequest) Format(encoder codec.Encoder, _ int, w io.Writer) {
	encoder.Uuid(r.ID, w)
}

func (_ *GetClusterInfoRequest) Kind() types.Typed {
	return MESSAGE_KIND
}

func (r *GetClusterInfoRequest) ResponseMessage() types.EndpointResponseMessage {

	if r.response == nil {
		r.response = &GetClusterInfoResponse{}
	}

	return r.response
}

func (r *GetClusterInfoRequest) Response() *GetClusterInfoResponse {
	return r.response
}

func (_ GetClusterInfoRequest) Type() types.Typed {
	return GET_CLUSTER_INFO_REQUEST
}

// GetClustersResponse ответ со списком кластеров
//
//  type GET_CLUSTERS_RESPONSE = 14
//  kind MESSAGE_KIND = 1
//  Info serialize.ClusterInfo
type GetClusterInfoResponse struct {
	Info serialize.ClusterInfo
}

func (res *GetClusterInfoResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	info := serialize.ClusterInfo{}
	info.Parse(decoder, version, r)
	res.Info = info
}

func (_ *GetClusterInfoResponse) Kind() types.Typed {
	return MESSAGE_KIND
}

func (_ *GetClusterInfoResponse) Type() types.Typed {
	return GET_CLUSTER_INFO_RESPONSE
}
