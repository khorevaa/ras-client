package protocol

import (
	"errors"
	"github.com/k0kubun/pp"
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/messages"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
)

var serviceVersions = []string{"3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0"}

var _ types.Endpoint = (*endpoint)(nil)

func (c *Client) NewEndpoint() (types.Endpoint, error) {

	if len(c.serviceVersion) > 0 {
		return c.OpenEndpoint(c.serviceVersion)
	}

	version := serviceVersions[len(serviceVersions)-1]

	end, err := c.OpenEndpoint(version)

	if err != nil {

		supportedVersion := detectSupportedVersion(err)
		if len(supportedVersion) > 0 {
			return nil, errors.New(pp.Sprint("ras no supported service version", serviceVersions))
		}

		c.serviceVersion = supportedVersion

		return c.OpenEndpoint(c.serviceVersion)

	}

	c.serviceVersion = version

	return end, err
}

func (c *Client) OpenEndpoint(version string) (types.Endpoint, error) {

	ack := &OpenEndpointMessageAck{}
	_, err := c.SendRequest(&OpenEndpointMessage{
		Version: version,
		ack:     ack,
	})

	if err != nil {
		return nil, err
	}

	end := c.registryNewEndpoint(ack)

	return end, nil
}

func (e *endpoint) GetClusters() ([]*serialize.ClusterInfo, error) {

	req := &messages.GetClustersRequest{}
	_, err := e.SendMessage(req)

	if err != nil {
		return nil, err
	}

	r := req.Response()

	return r.Clusters, err
}

func (e *endpoint) AuthenticateAgent(user, password string) error {

	_, err := e.SendMessage(&messages.AuthenticateAgentRequest{
		User:     user,
		Password: password,
	})

	return err
}

func (e *endpoint) AuthenticateCluster(uuid uuid.UUID, user, password string) error {

	_, err := e.SendMessage(&messages.ClusterAuthenticateRequest{
		ClusterID: uuid,
		User:      user,
		Password:  password,
	})

	return err
}

func (e *endpoint) AuthenticateInfobase(cluster uuid.UUID, user, password string) error {

	_, err := e.SendMessage(&messages.AuthenticateInfobaseRequest{
		ClusterID: cluster,
		User:      user,
		Password:  password,
	})

	return err
}

func (e *endpoint) GetClusterManagers(id uuid.UUID) ([]*serialize.ManagerInfo, error) {

	req := &messages.GetClusterManagersRequest{ClusterID: id}
	_, err := e.SendMessage(req)

	response := req.Response()

	return response.Managers, err
}

func (e *endpoint) GetClusterServices(id uuid.UUID) ([]*serialize.ServiceInfo, error) {

	req := messages.GetClusterServicesRequest{ClusterID: id}
	_, err := e.SendMessage(&req)

	response := req.Response()

	return response.Services, err
}

func (e *endpoint) GetClusterInfobases(id uuid.UUID) (serialize.InfobaseSummaryList, error) {

	req := &messages.GetInfobasesShortRequest{ClusterID: id}
	_, err := e.SendMessage(req)
	response := req.Response()
	return response.Infobases, err
}

func (e *endpoint) GetClusterConnections(id uuid.UUID) (serialize.ConnectionShortInfoList, error) {

	req := messages.GetConnectionsShortRequest{ID: id}
	_, err := e.SendMessage(&req)

	response := req.Response()

	return response.Connections, err
}

func (e *endpoint) DisconnectConnection(cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) error {

	req := messages.DisconnectConnectionRequest{
		ClusterID:    cluster,
		ProcessID:    process,
		ConnectionID: connection,
	}
	_, err := e.SendMessage(&req)

	return err
}

func (e *endpoint) GetInfobaseInfo(cluster uuid.UUID, infobase uuid.UUID) (serialize.InfobaseInfo, error) {

	req := &messages.GetInfobaseInfoRequest{ClusterID: cluster, InfobaseID: infobase}
	_, err := e.SendMessage(req)
	response := req.Response()
	return response.Infobase, err
}

func (e *endpoint) GetClusterInfo(cluster uuid.UUID) (serialize.ClusterInfo, error) {

	req := &messages.GetClusterInfoRequest{ID: cluster}
	_, err := e.SendMessage(req)
	response := req.Response()
	return response.Info, err
}

func (e *endpoint) CreateInfobase(cluster uuid.UUID, infobase serialize.InfobaseInfo, mode int) (serialize.InfobaseInfo, error) {

	req := &messages.CreateInfobaseRequest{ClusterID: cluster, Infobase: &infobase, Mode: mode}
	_, err := e.SendMessage(req)

	if err != nil {
		return serialize.InfobaseInfo{}, err
	}
	response := req.Response()
	return e.GetInfobaseInfo(cluster, response.InfobaseID)
}

func (e *endpoint) DropInfobase(cluster uuid.UUID, infobase uuid.UUID) error {

	req := &messages.DropInfobaseRequest{ClusterID: cluster, InfobaseId: infobase}
	_, err := e.SendMessage(req)
	return err

}

func (e *endpoint) UpdateSummaryInfobase(cluster uuid.UUID, infobase serialize.InfobaseSummaryInfo) error {
	req := &messages.UpdateInfobaseShortRequest{ClusterID: cluster, Infobase: infobase}
	_, err := e.SendMessage(req)
	return err
}

func (e *endpoint) UpdateInfobase(cluster uuid.UUID, infobase serialize.InfobaseInfo) error {

	req := &messages.UpdateInfobaseRequest{ClusterID: cluster, Infobase: infobase}
	_, err := e.SendMessage(req)
	return err

}

func (e *endpoint) GetInfobaseConnections(cluster uuid.UUID, infobase uuid.UUID) (serialize.ConnectionShortInfoList, error) {

	req := messages.GetInfobaseConnectionsShortRequest{ClusterID: cluster, InfobaseID: infobase}
	_, err := e.SendMessage(&req)

	response := req.Response()

	return response.Connections, err
}

func (e *endpoint) TerminateSession(cluster uuid.UUID, session uuid.UUID, msg string) error {

	req := messages.TerminateSessionRequest{
		ClusterID: cluster,
		SessionID: session,
		Message:   msg,
	}
	_, err := e.SendMessage(&req)

	return err
}

func (e *endpoint) GetInfobaseSessions(cluster uuid.UUID, infobase uuid.UUID) (serialize.SessionInfoList, error) {

	req := messages.GetInfobaseSessionsRequest{
		ClusterID:  cluster,
		InfobaseID: infobase,
	}
	_, err := e.SendMessage(&req)
	response := req.Response()

	return response.Sessions, err

}

func (e *endpoint) GetClusterSessions(cluster uuid.UUID) (serialize.SessionInfoList, error) {

	req := messages.GetSessionsRequest{
		ClusterID: cluster,
	}
	_, err := e.SendMessage(&req)
	response := req.Response()

	return response.Sessions, err

}

func (e *endpoint) GetClusterLocks(cluster uuid.UUID) (serialize.LocksList, error) {

	req := messages.GetLocksRequest{
		ClusterID: cluster,
	}
	_, err := e.SendMessage(&req)
	response := req.Response()

	return response.List, err

}

func (e *endpoint) GetInfobaseLocks(cluster uuid.UUID, infobase uuid.UUID) (serialize.LocksList, error) {

	req := messages.GetInfobaseLockRequest{
		ClusterID:  cluster,
		InfobaseID: infobase,
	}
	_, err := e.SendMessage(&req)
	response := req.Response()

	return response.List, err

}

func (e *endpoint) GetSessionLocks(cluster uuid.UUID, infobase uuid.UUID, session uuid.UUID) (serialize.LocksList, error) {

	req := messages.GetSessionLockRequest{
		ClusterID:  cluster,
		InfobaseID: infobase,
		SessionID:  session,
	}
	_, err := e.SendMessage(&req)
	response := req.Response()

	return response.List, err

}

func (e *endpoint) GetConnectionLocks(cluster uuid.UUID, connection uuid.UUID) (serialize.LocksList, error) {

	req := messages.GetConnectionLockRequest{
		ClusterID:    cluster,
		ConnectionID: connection,
	}
	_, err := e.SendMessage(&req)
	response := req.Response()

	return response.List, err

}
