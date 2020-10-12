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

func (e *endpoint) GetClusterInfobases(id uuid.UUID) ([]*serialize.InfobaseSummaryInfo, error) {

	req := messages.GetInfobasesShortRequest{ClusterID: id}
	_, err := e.SendMessage(req)
	response := req.Response()
	return response.Infobases, err
}

func (e *endpoint) GetClusterConnections(id uuid.UUID) (serialize.ConnectionInfoList, error) {

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

	req := &messages.GetInfobaseInfoRequest{ClusterID: cluster, InfobaseId: infobase}
	_, err := e.SendMessage(req)
	response := req.Response()
	return response.Infobase, err
}
