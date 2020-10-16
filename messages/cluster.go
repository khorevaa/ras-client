package messages

import (
	"github.com/khorevaa/ras-client/protocol/codec"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"io"
)

// GetClustersRequest получение списка кластеров
//
//  type GET_CLUSTERS_REQUEST = 11
//  kind MESSAGE_KIND = 1
//  respond GetClustersResponse
type GetClustersRequest struct{}

func (_ *GetClustersRequest) Sig() esig.ESIG {
	return esig.Nil
}

func (_ *GetClustersRequest) Format(_ codec.Encoder, _ int, _ io.Writer) {}

func (_ GetClustersRequest) Type() EndpointMessageType {
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

func (_ *GetClustersResponse) Type() EndpointMessageType {
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

func (_ *GetClusterInfoRequest) Type() EndpointMessageType {
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

func (_ *GetClusterInfoResponse) Type() EndpointMessageType {
	return GET_CLUSTER_INFO_RESPONSE
}
