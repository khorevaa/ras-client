package protocol

import (
	"errors"
	"github.com/k0kubun/pp"
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/types"
)

const DefaultFormat = 256

var serviceVersions = []string{"3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0"}

func (m *RASConn) NewEndpoint() (types.Endpoint, error) {

	if len(m.serviceVersion) > 0 {
		return m.OpenEndpoint(m.serviceVersion)
	}

	version := serviceVersions[len(serviceVersions)-1]

	end, err := m.OpenEndpoint(version)

	if err != nil {

		supportedVersion := detectSupportedVersion(err)
		if len(supportedVersion) > 0 {
			return nil, errors.New(pp.Sprint("ras no supported service version", serviceVersions))
		}

		m.serviceVersion = supportedVersion

		return m.OpenEndpoint(m.serviceVersion)

	}

	m.serviceVersion = version

	return end, err
}

func (m *RASConn) OpenEndpoint(version string) (types.Endpoint, error) {

	ack := &OpenEndpointMessageAck{}
	_, err := m.SendRequest(&OpenEndpointMessage{
		Version: version,
		ack:     ack,
	})

	if err != nil {
		return nil, err
	}

	end := m.registryNewEndpoint(ack)

	return end, nil
}

func (e *endpoint) GetClusters() ([]*ClusterInfo, error) {

	req := &GetClustersRequest{}
	_, err := e.SendMessage(req)

	if err != nil {
		return nil, err
	}

	r := req.Response()

	return r.Clusters, err
}

func (e *endpoint) AuthenticateAgent(user, password string) error {

	_, err := e.SendMessage(&AuthenticateAgentRequest{
		user:     user,
		password: password,
	})

	return err
}

func (e *endpoint) AuthenticateCluster(uuid uuid.UUID, user, password string) error {

	_, err := e.SendMessage(&ClusterAuthenticateRequest{
		ID:       uuid,
		user:     user,
		password: password,
	})

	return err
}

func (e *endpoint) GetClusterManagers(id uuid.UUID) ([]*ManagerInfo, error) {

	req := &GetClusterManagersRequest{ID: id}
	_, err := e.SendMessage(req)

	response := req.Response()

	return response.Managers, err
}

func (e *endpoint) GetClusterServices(id uuid.UUID) ([]*ServiceInfo, error) {

	req := GetClusterServicesRequest{ID: id}
	_, err := e.SendMessage(&req)

	response := req.Response()

	return response.Services, err
}

func (e *endpoint) GetClusterInfobases(id uuid.UUID) ([]*InfobaseInfo, error) {

	req := GetInfobasesShortRequest{ID: id}
	_, err := e.SendMessage(req)
	response := req.Response()
	return response.Infobases, err
}

func (e *endpoint) GetClusterConnections(id uuid.UUID) ([]*ConnectionInfo, error) {

	req := GetConnectionsShortRequest{ID: id}
	_, err := e.SendMessage(&req)

	response := req.Response()

	return response.Connections, err
}

type ClusterInfo struct {
	UUID                       uuid.UUID `rac:"cluster"` // UUID cluster                    : 6d6958e1-a96c-4999-a995-698a0298161e
	Host                       string    // Host                          : Sport2
	Port                       int16     // Port                          : 1541
	Name                       string    // Name                          : "Новый кластер"
	ExpirationTimeout          int       // ExpirationTimeout expiration-timeout            : 0
	LifetimeLimit              int       // LifetimeLimit lifetime-limit                : 0
	MaxMemorySize              int       // MaxMemorySize max-memory-count               : 0
	MaxMemoryTimeLimit         int       // MaxMemoryTimeLimit max-memory-time-limit         : 0
	SecurityLevel              int       // SecurityLevel security-level                : 0
	SessionFaultToleranceLevel int       // SessionFaultToleranceLevel session-fault-tolerance-level : 0
	LoadBalancingMode          int       // LoadBalancingMode load-balancing-Mode           : performance
	ErrorsCountThreshold       int       // ErrorsCountThreshold errors-count-threshold        : 0
	KillProblemProcesses       bool      // KillProblemProcesses kill-problem-processes        : 0
	KillByMemoryWithDump       bool      // KillByMemoryWithDump kill-by-memory-with-dump      : 0
	LifeTimeLimit              int
}
