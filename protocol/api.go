package protocol

import (
	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
)

const DefaultFormat = 256

func (m *RASConn) OpenEndpointByVersion(version string) (RespondMessage, error) {

	resp, err := m.SendRequest(&OpenEndpointMessage{
		Version: version,
	}, &OpenEndpointMessageAck{}, &EndpointFeature{})

	return resp, err
}

func (m *RASConn) OpenEndpoint(version string) (RespondMessage, error) {

	resp, err := m.SendRequest(&OpenEndpointMessage{
		Version: version,
	}, &OpenEndpointMessageAck{}, &EndpointFeature{})

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

func (m *RASConn) GetClusters() ([]ClusterInfo, error) {

	resp, err := m.SendEndpointRequest(&GetClustersRequest{}, &GetClustersResponse{})

	dry.PanicIfErr(err)

	pp.Println(resp)

	return resp.(GetClustersResponse).Clusters, err
}

func (m *RASConn) AuthenticateAgent(user, password string) error {

	_, err := m.SendEndpointRequest(&AuthenticateAgentRequest{
		user:     user,
		password: password,
	})

	dry.PanicIfErr(err)

	//pp.Println(resp)

	return err
}

type ClusterInfo struct {
	UUID                       string `rac:"cluster"` // UUID cluster                    : 6d6958e1-a96c-4999-a995-698a0298161e
	Host                       string // Host                          : Sport2
	Port                       int    // Port                          : 1541
	Name                       string // Name                          : "Новый кластер"
	ExpirationTimeout          int    // ExpirationTimeout expiration-timeout            : 0
	LifetimeLimit              int    // LifetimeLimit lifetime-limit                : 0
	MaxMemorySize              int    // MaxMemorySize max-memory-size               : 0
	MaxMemoryTimeLimit         int    // MaxMemoryTimeLimit max-memory-time-limit         : 0
	SecurityLevel              int    // SecurityLevel security-level                : 0
	SessionFaultToleranceLevel int    // SessionFaultToleranceLevel session-fault-tolerance-level : 0
	LoadBalancingMode          int    // LoadBalancingMode load-balancing-mode           : performance
	ErrorsCountThreshold       int    // ErrorsCountThreshold errors-count-threshold        : 0
	KillProblemProcesses       bool   // KillProblemProcesses kill-problem-processes        : 0
	KillByMemoryWithDump       bool   // KillByMemoryWithDump kill-by-memory-with-dump      : 0
}
