package types

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/esig"
	"github.com/v8platform/rac/serialize"
	"io"
)

type Endpoint interface {
	Version() int

	SendMessage(req EndpointRequestMessage) (interface{}, error)
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

type RequestMessage interface {
	Type() Typed
	Format(codec codec.Encoder, w io.Writer)
}

type ResponseMessage interface {
	Type() Typed
	Parse(codec codec.Decoder, r io.Reader)
}

type EndpointRequestMessage interface {
	Type() Typed
	Format(encoder codec.Encoder, version int, w io.Writer)
	Sig() esig.ESIG
}

type EndpointResponseMessage interface {
	Type() Typed
	Parse(decoder codec.Decoder, version int, r io.Reader)
}

type Typed interface {
	Type() int
}
