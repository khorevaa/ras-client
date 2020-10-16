package rclient

import (
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

type ClientApi interface {
	Version() int

	Close()

	AuthenticateAgent(user, password string) error
	AuthenticateCluster(cluster uuid.UUID, user, password string) error
	AuthenticateInfobase(cluster uuid.UUID, user, password string) error

	GetClusters() ([]*serialize.ClusterInfo, error)
	GetClusterInfo(cluster uuid.UUID) (serialize.ClusterInfo, error)
	GetClusterInfobases(id uuid.UUID) (serialize.InfobaseSummaryList, error)
	GetClusterServices(id uuid.UUID) ([]*serialize.ServiceInfo, error)
	GetClusterManagers(id uuid.UUID) ([]*serialize.ManagerInfo, error)

	GetClusterConnections(uuid uuid.UUID) (serialize.ConnectionShortInfoList, error)
	GetClusterSessions(cluster uuid.UUID) (serialize.SessionInfoList, error)
	GetClusterLocks(cluster uuid.UUID) (serialize.LocksList, error)
	DisconnectConnection(cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) error
	TerminateSession(cluster uuid.UUID, session uuid.UUID, msg string) error

	GetInfobaseInfo(cluster uuid.UUID, infobase uuid.UUID) (serialize.InfobaseInfo, error)
	GetInfobaseConnections(cluster uuid.UUID, infobase uuid.UUID) (serialize.ConnectionShortInfoList, error)
	GetInfobaseSessions(cluster uuid.UUID, infobase uuid.UUID) (serialize.SessionInfoList, error)
	GetInfobaseLocks(cluster uuid.UUID, infobase uuid.UUID) (serialize.LocksList, error)

	CreateInfobase(cluster uuid.UUID, infobase serialize.InfobaseInfo, mode int) (serialize.InfobaseInfo, error)
	UpdateSummaryInfobase(cluster uuid.UUID, infobase serialize.InfobaseSummaryInfo) error
	UpdateInfobase(cluster uuid.UUID, infobase serialize.InfobaseInfo) error
	DropInfobase(cluster uuid.UUID, infobase uuid.UUID) error

	GetSessionLocks(cluster uuid.UUID, infobase uuid.UUID, session uuid.UUID) (serialize.LocksList, error)
	GetConnectionLocks(cluster uuid.UUID, connection uuid.UUID) (serialize.LocksList, error)
}
