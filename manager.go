package rac

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

const defaultTimeout = 15 * time.Second

type clusterAuth struct {
	User     string
	Password string
}

func (auth clusterAuth) Args() []string {

	if len(auth.User) == 0 {
		return []string{}
	}

	args := []string{
		"--cluster-user=" + auth.User,
	}

	if len(auth.Password) > 0 {
		args = append(args, "--cluster-pwd="+auth.Password)
	}

	return args

}

type Logger interface {
	Errorf(msg string, args ...interface{})
}

type nullLogger struct{}

func (n nullLogger) Errorf(msg string, args ...interface{}) {}

var _ Logger = (*nullLogger)(nil)

type Manager struct {
	Host    string
	Port    string
	options *Options

	defCluster     ClusterInfo
	defClusterAuth AuthSig

	idxServers map[string]ServerInfo
	idxCluster map[string]ClusterInfo

	autoUpdate      bool
	updateAt        time.Time
	updateInterval  time.Duration
	lastUpdateError error

	log Logger
}

func (m *Manager) ClusterSig() (string, string, string) {

	if len(m.defCluster.UUID) == 0 {
		return "", "", ""
	}

	var (
		user, pwd string
	)

	if m.defClusterAuth != nil {
		user, pwd = m.defClusterAuth.Sig()
	}

	return m.defCluster.UUID, user, pwd
}

func (m *Manager) GetDefCluster() (ClusterInfo, AuthSig) {

	return m.defCluster, m.defClusterAuth

}

func (m *Manager) SetCluster(cluster Clusterable) error {

	uuid, user, pwd := cluster.ClusterSig()

	resp, err := m.Clusters(ClustersInfo{UUID: uuid})

	if err != nil {
		return err
	}

	m.defCluster = resp.Info

	if len(user) > 0 {
		m.defClusterAuth = Auth{user, pwd}
	}

	return nil

}

func NewManager(hostPort string, opts ...Option) (*Manager, error) {

	host, port, _ := net.SplitHostPort(hostPort)

	options := &Options{
		v8Version:         "8.3",
		updateInterval:    time.Hour,
		timeout:           defaultTimeout,
		autoSetDefCluster: true,
	}

	for _, opt := range opts {
		opt(options)
	}

	manager := &Manager{
		Host:    host,
		Port:    port,
		options: options,
		log:     &nullLogger{},
		defClusterAuth: Auth{
			User: options.clusterAuth.User,
			Pwd:  options.clusterAuth.Password,
		},
		updateInterval: options.updateInterval,
		idxCluster:     make(map[string]ClusterInfo),
		idxServers:     make(map[string]ServerInfo),
	}

	err := manager.init()

	return manager, err

}

func (m *Manager) init() error {

	err := m.updateCluster()

	if err != nil {
		return err
	}

	m.detectDefCluster()

	m.process()

	return nil

}

func (m *Manager) detectDefCluster() {

	if m.options.autoSetDefCluster {

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

func clusterSigParams(cluster, user, pwd string) map[string]string {

	return map[string]string{
		"--cluster":      cluster,
		"--cluster-user": user,
		"--cluster-pwd":  pwd,
	}

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

func (m *Manager) do(command string, setParams ...map[string]string) *RawRespond {

	var args []string

	args = append(args, strings.Fields(command)...)

	params := mergeParams(setParams...)

	for key, value := range params {

		args = append(args, fmt.Sprintf("%s=%s", key, value))

	}

	raw, err := m.run(args...)

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

func (m *Manager) run(args ...string) ([]byte, error) {

	var (
		ctx     context.Context
		respond []byte
	)

	ctx = context.Background()

	if m.options.ctx != nil {
		ctx = m.options.ctx
	}

	if m.options.timeout > 0 {

		ctx, _ = context.WithTimeout(ctx, defaultTimeout)
	}

	rac := m.options.racPath // TODO Полечнеие по версии 1С

	args = append(args, m.ServerSig())

	cmd := exec.CommandContext(ctx, rac, args...)

	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)
	errch := make(chan error, 1)

	err := cmd.Start()
	if err != nil {
		return respond, fmt.Errorf("Произошла ошибка запуска:\n\terr:%v\n\tПараметры: %v\n\t", err.Error(), cmd.Args)
	}

	// запускаем в горутине т.к. наблюдалось что при выполнении RAC может происходить зависон, нам нужен таймаут
	go func() {
		errch <- cmd.Wait()
	}()

	select {
	case <-ctx.Done(): // timeout

		if ctx.Err() == context.DeadlineExceeded {
			m.log.Errorf("Выполнение команды прервано по таймауту\n\tПараметры: %v\n\t", cmd.Args)
		}

		return respond, ctx.Err()

	case err := <-errch:
		if err != nil {

			stderr := cmd.Stderr.(*bytes.Buffer).String()
			errText := fmt.Sprintf("Произошла ошибка запуска:\n\terr:%v\n\tПараметры: %v\n\t", err.Error(), cmd.Args)
			if stderr != "" {
				errText += fmt.Sprintf("StdErr:%v\n", stderr)
			}

			return respond, errors.New(errText)

		} else {

			in := cmd.Stdout.(*bytes.Buffer).Bytes()

			respond, err = decodeOutBytes(in)

			if err != nil {
				return respond, err
			}

			return respond, nil
		}
	}

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
