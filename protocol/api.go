package protocol

import (
	uuid "github.com/satori/go.uuid"
)

const DefaultFormat = 256

func (m *RASConn) OpenEndpointByVersion(version string) (RespondMessage, error) {

	resp, err := m.SendRequest(&OpenEndpointMessage{
		Version: version,
	}, &OpenEndpointMessageAck{}, &EndpointFailure{})

	return resp, err
}

func (m *RASConn) OpenEndpoint(version string) (RespondMessage, error) {

	resp, err := m.SendRequest(&OpenEndpointMessage{
		Version: version,
	}, &OpenEndpointMessageAck{}, &EndpointFailure{})

	switch r := resp.(type) {

	case *OpenEndpointMessageAck:

		m.Endpoint = &Endpoint{
			r.EndpointID,
			r.ServiceID,
			r.Version,
			DefaultFormat,
		}

	}

	return resp, err
}

func (m *RASConn) GetClusters() ([]*ClusterInfo, error) {

	response := GetClustersResponse{}
	resp, err := m.SendEndpointMessage(&GetClustersRequest{}, &response)

	if err != nil {
		return nil, err
	}

	//dry.PanicIfErr(err)

	//pp.Println(resp)

	message := resp.waitResponse

	r := message.(*GetClustersResponse)

	return r.Clusters, err
}

func (m *RASConn) AuthenticateAgent(user, password string) error {

	_, err := m.SendEndpointMessage(&AuthenticateAgentRequest{
		user:     user,
		password: password,
	})

	return err
}

func (m *RASConn) AuthenticateCluster(uuid uuid.UUID, user, password string) error {

	_, err := m.SendEndpointMessage(&ClusterAuthenticateRequest{
		ID:       uuid,
		user:     user,
		password: password,
	})

	return err
}

func (m *RASConn) GetClusterManagers(id uuid.UUID) ([]*ManagerInfo, error) {

	response := GetClusterManagersResponse{}
	_, err := m.SendEndpointMessage(&GetClusterManagersRequest{
		ID: id,
	}, &response)

	return response.Managers, err
}

func (m *RASConn) GetClusterServices(id uuid.UUID) ([]*ServiceInfo, error) {

	response := GetClusterServicesResponse{}
	_, err := m.SendEndpointMessage(GetClusterServicesRequest{
		ID: id,
	}, &response)

	return response.Services, err
}

func (m *RASConn) GetClusterInfobases(id uuid.UUID) ([]*InfobaseInfo, error) {

	response := GetInfobasesShortResponse{}
	_, err := m.SendEndpointMessage(GetInfobasesShortRequest{
		ID: id,
	}, &response)

	return response.Infobases, err
}

func (m *RASConn) GetClusterConnections(id uuid.UUID) ([]*ConnectionInfo, error) {

	response := GetConnectionsShortResponse{}
	_, err := m.SendEndpointMessage(GetConnectionsShortRequest{
		ID: id,
	}, &response)

	return response.Connections, err
}

type ClusterInfo struct {
	UUID                       string `rac:"cluster"` // UUID cluster                    : 6d6958e1-a96c-4999-a995-698a0298161e
	Host                       string // Host                          : Sport2
	Port                       int    // Port                          : 1541
	Name                       string // Name                          : "Новый кластер"
	ExpirationTimeout          int    // ExpirationTimeout expiration-timeout            : 0
	LifetimeLimit              int    // LifetimeLimit lifetime-limit                : 0
	MaxMemorySize              int    // MaxMemorySize max-memory-count               : 0
	MaxMemoryTimeLimit         int    // MaxMemoryTimeLimit max-memory-time-limit         : 0
	SecurityLevel              int    // SecurityLevel security-level                : 0
	SessionFaultToleranceLevel int    // SessionFaultToleranceLevel session-fault-tolerance-level : 0
	LoadBalancingMode          int    // LoadBalancingMode load-balancing-mode           : performance
	ErrorsCountThreshold       int    // ErrorsCountThreshold errors-count-threshold        : 0
	KillProblemProcesses       bool   // KillProblemProcesses kill-problem-processes        : 0
	KillByMemoryWithDump       bool   // KillByMemoryWithDump kill-by-memory-with-dump      : 0
}
