package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

var _ clusterApi = (*Client)(nil)

func (c *Client) GetClusters(ctx context.Context) ([]*serialize.ClusterInfo, error) {

	req := &messages.GetClustersRequest{}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClustersResponse)

	return response.Clusters, err
}

func (c *Client) GetClusterInfo(ctx context.Context, cluster uuid.UUID) (serialize.ClusterInfo, error) {

	req := &messages.GetClusterInfoRequest{ClusterID: cluster}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return serialize.ClusterInfo{}, err
	}

	response := resp.(*messages.GetClusterInfoResponse)
	return response.Info, nil
}

func (c *Client) GetClusterManagers(ctx context.Context, id uuid.UUID) ([]*serialize.ManagerInfo, error) {

	req := &messages.GetClusterManagersRequest{ClusterID: id}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClusterManagersResponse)

	return response.Managers, err
}

func (c *Client) GetClusterServices(ctx context.Context, id uuid.UUID) ([]*serialize.ServiceInfo, error) {

	req := &messages.GetClusterServicesRequest{ClusterID: id}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClusterServicesResponse)

	return response.Services, err
}

func (c *Client) GetClusterInfobases(ctx context.Context, id uuid.UUID) (serialize.InfobaseSummaryList, error) {

	req := &messages.GetInfobasesShortRequest{ClusterID: id}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetInfobasesShortResponse)

	response.Infobases.Each(func(info *serialize.InfobaseSummaryInfo) {
		info.ClusterID = id
	})

	return response.Infobases, err
}