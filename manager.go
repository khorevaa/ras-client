package rac

import (
	"context"
	"fmt"
	"github.com/v8platform/find"
	"net"
	"strings"
	"time"
)

const defaultTimeout = 15 * time.Second

type Runner interface {
	RunCtx(ctx context.Context, command string, args []string) (respond []byte, err error)
}

type Manager struct {
	Host string
	Port string

	options *ManagerOptions

	defCluster     ClusterInfo
	defClusterAuth Auth

	idxServers map[string]ServerInfo
	idxCluster map[string]ClusterInfo

	autoUpdate      bool
	updateAt        time.Time
	updateInterval  time.Duration
	lastUpdateError error

	runner  Runner
	log     Logger
	racPath string
}

func (m *Manager) ClusterAuth() (user string, pwd string) {

	return m.defClusterAuth.Sig()
}

func (m *Manager) ClusterSig() string {

	if len(m.defCluster.UUID) == 0 {
		return ""
	}

	return m.defCluster.UUID
}

func (m *Manager) GetDefCluster() (ClusterInfo, Auth) {

	return m.defCluster, m.defClusterAuth

}

func (m *Manager) SetDefCluster(cluster ClusterSig, auth AuthSig) error {

	uuid := cluster.ClusterSig()

	resp, err := m.Clusters(ClustersInfo{UUID: uuid})

	if err != nil {
		return err
	}

	m.defCluster = resp.Info

	if auth != nil {

		user, pwd := auth.Auth()

		m.defClusterAuth = Auth{
			User: user,
			Pwd:  pwd,
		}
	}

	return nil

}

func newManager(hostPort string, options *ManagerOptions) (*Manager, error) {

	host, port, _ := net.SplitHostPort(hostPort)

	if options == nil {
		options = &ManagerOptions{
			V8Version:      "8.3",
			UpdateInterval: time.Hour,
			Timeout:        defaultTimeout,
			DetectCluster:  true,
			ClusterAuth: struct {
				User string
				Pwd  string
			}{User: "", Pwd: ""},
			Logger: nullLogger{},
		}
	}

	manager := &Manager{
		Host:    host,
		Port:    port,
		options: options,
		log:     options.Logger,
		defClusterAuth: Auth{
			User: options.ClusterAuth.User,
			Pwd:  options.ClusterAuth.Pwd,
		},
		updateInterval: options.UpdateInterval,
		idxCluster:     make(map[string]ClusterInfo),
		idxServers:     make(map[string]ServerInfo),
		runner:         options.Runner,
	}

	err := manager.init()

	return manager, err

}

func extractManagerOptions(opts []interface{}) *ManagerOptions {

	var options *ManagerOptions

	for _, opt := range opts {

		switch o := opt.(type) {
		case *ManagerOptions:
			options = o
		case ManagerOptions:
			options = &o
		}

	}

	return options

}

func NewManager(hostPort string, opts ...interface{}) (*Manager, error) {

	options := extractManagerOptions(opts)

	return newManager(hostPort, options)

}

func (m *Manager) init() error {

	if m.runner == nil {
		m.runner = initDefaultRunner(m.options)
	}

	if len(m.racPath) == 0 {
		var err error
		m.racPath, err = find.RAC(find.WithVersion(m.options.V8Version))
		if err != nil {
			return err
		}
	}

	err := m.updateCluster()

	if err != nil {
		return err
	}

	m.detectDefCluster()

	m.process()

	return nil

}

func initDefaultRunner(options *ManagerOptions) Runner {

	return &runner{
		Timeout:         options.Timeout,
		TryTimeoutCount: options.TryTimeoutCount,
	}

}

func (m *Manager) detectDefCluster() {

	if m.options.DetectCluster {

		var cluster ClusterInfo
		for _, info := range m.idxCluster {

			cluster = info
			break
		}

		m.defCluster = cluster
	}

}

func (m *Manager) process() {

	go m.pullUpdater()

}

func (m *Manager) updateCluster() error {

	resp, err := m.Clusters(ClustersList{})

	if err != nil {
		return err
	}

	m.idxCluster = make(map[string]ClusterInfo)
	for _, info := range resp.List {
		m.idxCluster[info.UUID] = info
	}

	m.updateAt = time.Now()

	return nil

}

func (m *Manager) Do(what Valued, opts ...interface{}) (*RawRespond, error) {

	return m.DoCtx(context.Background(), what, opts...)

}

func (m *Manager) DoCtx(ctx context.Context, what Valued, opts ...interface{}) (respond *RawRespond, err error) {

	command := what.Command()
	params := what.Values()

	doOptions := extractOptions(opts)
	err = doOptions.Options(opts...)
	if err != nil {
		return nil, err
	}
	paramsOpts := extractParams(opts)
	mParams := m.embedParams()
	raw := m.doCtx(ctx, command, paramsOpts, params, mParams, doOptions.Values())

	if raw.Error != nil {
		return raw, raw.Error
	}

	err = what.Parse(raw)
	if err != nil {
		return raw, err
	}

	return raw, nil

}

func (m *Manager) doCtx(ctx context.Context, command DoCommand, setParams ...map[string]string) *RawRespond {

	var args []string

	params := mergeParams(setParams...)
	if err := command.Check(params); err != nil {
		return newRawRespond(nil, err)
	}

	args = append(args, strings.Fields(command.Command())...)

	for key, value := range params {

		switch key {
		case "--licenses", "--drop-database", "--clear-database": // TODO Заглушка для булевных значений, которые подставляется без значения
			v, _ := parseBool(value)
			if v {
				args = append(args, key)
			}

		default:
			args = append(args, fmt.Sprintf("%s=%s", key, value))
		}

	}

	rac := m.racPath

	args = append(args, m.ServerSig())

	raw, err := m.runner.RunCtx(ctx, rac, args)

	res := newRawRespond(raw, err)

	return res
}

func mergeParams(params ...map[string]string) map[string]string {

	merged := make(map[string]string)

	for _, s := range params {

		for key, value := range s {

			if len(value) == 0 {
				continue
			}

			merged[key] = value
		}

	}

	return merged
}

func (m *Manager) ServerSig() string {

	if len(m.Host) == 0 {
		return ""
	}

	return net.JoinHostPort(m.Host, m.Port)

}

func (m *Manager) pullUpdater() {

	if m.updateInterval == 0 {
		return
	}

	ticker := time.NewTicker(m.updateInterval)

	for {
		select {
		// handle incoming updates
		case <-ticker.C:

			err := m.updateCluster()

			if err != nil {
				m.lastUpdateError = err
			}
		}
		//// call to stop polling
		//case <-m.options.ctx.Done():
		//	return
		//}
	}

}
