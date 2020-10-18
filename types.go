package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/serialize"
	uuid "github.com/satori/go.uuid"
)

type Api interface {
	Version() string

	Close() error

	authApi
	clusterApi
	sessionApi
	locksApi
	connectionsApi
	infobaseApi
}

type authApi interface {
	AuthenticateAgent(user, password string)
	AuthenticateCluster(cluster uuid.UUID, user, password string)
	AuthenticateInfobase(cluster uuid.UUID, user, password string)
}

type clusterApi interface {
	GetClusters(ctx context.Context) ([]*serialize.ClusterInfo, error)
	GetClusterInfo(ctx context.Context, cluster uuid.UUID) (serialize.ClusterInfo, error)
	GetClusterInfobases(ctx context.Context, id uuid.UUID) (serialize.InfobaseSummaryList, error)
	GetClusterServices(ctx context.Context, id uuid.UUID) ([]*serialize.ServiceInfo, error)
	GetClusterManagers(ctx context.Context, id uuid.UUID) ([]*serialize.ManagerInfo, error)
}

type sessionApi interface {
	GetClusterSessions(ctx context.Context, cluster uuid.UUID) (serialize.SessionInfoList, error)
	GetInfobaseSessions(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.SessionInfoList, error)
	TerminateSession(ctx context.Context, cluster uuid.UUID, session uuid.UUID, msg string) error
}

type connectionsApi interface {
	GetClusterConnections(ctx context.Context, uuid uuid.UUID) (serialize.ConnectionShortInfoList, error)
	GetInfobaseConnections(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.ConnectionShortInfoList, error)
	DisconnectConnection(ctx context.Context, cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) error
}

type locksApi interface {
	GetClusterLocks(ctx context.Context, cluster uuid.UUID) (serialize.LocksList, error)
	GetInfobaseLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.LocksList, error)
	GetSessionLocks(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, session uuid.UUID) (serialize.LocksList, error)
	GetConnectionLocks(ctx context.Context, cluster uuid.UUID, connection uuid.UUID) (serialize.LocksList, error)
}

type infobaseApi interface {
	CreateInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseInfo, mode int) (serialize.InfobaseInfo, error)
	UpdateSummaryInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseSummaryInfo) error
	UpdateInfobase(ctx context.Context, cluster uuid.UUID, infobase serialize.InfobaseInfo) error
	DropInfobase(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID, mode int) error
	GetInfobaseInfo(ctx context.Context, cluster uuid.UUID, infobase uuid.UUID) (serialize.InfobaseInfo, error)
}
