package ras_client

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/messages"
	"github.com/v8platform/rac/serialize"
)

//var _ types.Endpoint = (*endpoint)(nil)

func (c *Client) GetClusters(ctx context.Context) ([]*serialize.ClusterInfo, error) {

	req := &messages.GetClustersRequest{}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClustersResponse)

	return response.Clusters, err
}

func (e *Client) AuthenticateAgent(user, password string) {

	panic("TODO Надо понять когда это авторизация должна запускаться")
	//_, err := e.sendEndpointRequest(&messages.AuthenticateAgentRequest{
	//	User:     user,
	//	Password: password,
	//})

	//return err
}

func (c *Client) AuthenticateCluster(cluster uuid.UUID, user, password string) {

	c.pool.SetClusterAuth(cluster, user, password)

}

func (c *Client) AuthenticateInfobase(cluster uuid.UUID, user, password string) {

	c.pool.SetInfobaseAuth(cluster, user, password)

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

//
func (c *Client) GetClusterServices(ctx context.Context, id uuid.UUID) ([]*serialize.ServiceInfo, error) {

	req := &messages.GetClusterServicesRequest{ClusterID: id}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetClusterServicesResponse)
	//response.Infobases.Each(func(info *serialize.InfobaseSummaryInfo) {
	//	info.Cluster = id
	//})

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
		info.Cluster = id
	})

	return response.Infobases, err
}

//func (e *endpoint) GetClusterConnections(id uuid.UUID) (serialize.ConnectionShortInfoList, error) {
//
//	req := messages.GetConnectionsShortRequest{ID: id}
//	_, err := e.SendMessage(&req)
//
//	response := req.Response()
//
//	return response.Connections, err
//}
//
//func (e *endpoint) DisconnectConnection(cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) error {
//
//	req := messages.DisconnectConnectionRequest{
//		ClusterID:    cluster,
//		ProcessID:    process,
//		ConnectionID: connection,
//	}
//	_, err := e.SendMessage(&req)
//
//	return err
//}
//
//func (e *endpoint) GetInfobaseInfo(cluster uuid.UUID, infobase uuid.UUID) (serialize.InfobaseInfo, error) {
//
//	req := &messages.GetInfobaseInfoRequest{ClusterID: cluster, InfobaseID: infobase}
//	_, err := e.SendMessage(req)
//	response := req.Response()
//	return response.Infobase, err
//}
//
//func (e *endpoint) GetClusterInfo(cluster uuid.UUID) (serialize.ClusterInfo, error) {
//
//	req := &messages.GetClusterInfoRequest{ID: cluster}
//	_, err := e.SendMessage(req)
//	response := req.Response()
//	return response.Info, err
//}
//
//func (e *endpoint) CreateInfobase(cluster uuid.UUID, infobase serialize.InfobaseInfo, mode int) (serialize.InfobaseInfo, error) {
//
//	req := &messages.CreateInfobaseRequest{ClusterID: cluster, Infobase: &infobase, Mode: mode}
//	_, err := e.SendMessage(req)
//
//	if err != nil {
//		return serialize.InfobaseInfo{}, err
//	}
//	response := req.Response()
//	return e.GetInfobaseInfo(cluster, response.InfobaseID)
//}
//
//func (e *endpoint) DropInfobase(cluster uuid.UUID, infobase uuid.UUID) error {
//
//	req := &messages.DropInfobaseRequest{ClusterID: cluster, InfobaseId: infobase}
//	_, err := e.SendMessage(req)
//	return err
//
//}
//
//func (e *endpoint) UpdateSummaryInfobase(cluster uuid.UUID, infobase serialize.InfobaseSummaryInfo) error {
//	req := &messages.UpdateInfobaseShortRequest{ClusterID: cluster, Infobase: infobase}
//	_, err := e.SendMessage(req)
//	return err
//}
//
//func (e *endpoint) UpdateInfobase(cluster uuid.UUID, infobase serialize.InfobaseInfo) error {
//
//	req := &messages.UpdateInfobaseRequest{ClusterID: cluster, Infobase: infobase}
//	_, err := e.SendMessage(req)
//	return err
//
//}
//
func (c *Client) GetInfobaseConnections(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.ConnectionShortInfoList, error) {

	req := &messages.GetInfobaseConnectionsShortRequest{ClusterID: cluster, InfobaseID: infobase}
	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	response := resp.(*messages.GetInfobaseConnectionsShortResponse)

	response.Connections.Each(func(info *serialize.ConnectionShortInfo) {
		info.ClusterID = cluster
		info.InfobaseID = infobase
	})

	return response.Connections, nil
}

//func (e *endpoint) TerminateSession(cluster uuid.UUID, session uuid.UUID, msg string) error {
//
//	req := messages.TerminateSessionRequest{
//		ClusterID: cluster,
//		SessionID: session,
//		Message:   msg,
//	}
//	_, err := e.SendMessage(&req)
//
//	return err
//}
//
//func (e *endpoint) GetInfobaseSessions(cluster uuid.UUID, infobase uuid.UUID) (serialize.SessionInfoList, error) {
//
//	req := messages.GetInfobaseSessionsRequest{
//		ClusterID:  cluster,
//		InfobaseID: infobase,
//	}
//	_, err := e.SendMessage(&req)
//	response := req.Response()
//
//	return response.Sessions, err
//
//}
//
//func (e *endpoint) GetClusterSessions(cluster uuid.UUID) (serialize.SessionInfoList, error) {
//
//	req := messages.GetSessionsRequest{
//		ClusterID: cluster,
//	}
//	_, err := e.SendMessage(&req)
//	response := req.Response()
//
//	return response.Sessions, err
//
//}
//
//func (e *endpoint) GetClusterLocks(cluster uuid.UUID) (serialize.LocksList, error) {
//
//	req := messages.GetLocksRequest{
//		ClusterID: cluster,
//	}
//	_, err := e.SendMessage(&req)
//	response := req.Response()
//
//	return response.List, err
//
//}
//
//func (e *endpoint) GetInfobaseLocks(cluster uuid.UUID, infobase uuid.UUID) (serialize.LocksList, error) {
//
//	req := messages.GetInfobaseLockRequest{
//		ClusterID:  cluster,
//		InfobaseID: infobase,
//	}
//	_, err := e.SendMessage(&req)
//	response := req.Response()
//
//	return response.List, err
//
//}
//
//func (e *endpoint) GetSessionLocks(cluster uuid.UUID, infobase uuid.UUID, session uuid.UUID) (serialize.LocksList, error) {
//
//	req := messages.GetSessionLockRequest{
//		ClusterID:  cluster,
//		InfobaseID: infobase,
//		SessionID:  session,
//	}
//	_, err := e.SendMessage(&req)
//	response := req.Response()
//
//	return response.List, err
//
//}
//
//func (e *endpoint) GetConnectionLocks(cluster uuid.UUID, connection uuid.UUID) (serialize.LocksList, error) {
//
//	req := messages.GetConnectionLockRequest{
//		ClusterID:    cluster,
//		ConnectionID: connection,
//	}
//	_, err := e.SendMessage(&req)
//	response := req.Response()
//
//	return response.List, err
//
//}
