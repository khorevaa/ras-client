package rac

import (
	"context"
	"time"
)

type Option func(o *Options)

type Options struct {
	v8Version         string
	racPath           string
	updateInterval    time.Duration
	timeout           time.Duration
	ctx               context.Context
	autoSetDefCluster bool
	clusterAuth       Auth
}

type DoOption func(o *DoOptions)

type DoOptions struct {
	Cluster string
	Process string

	InfobaseAuth *Auth
	ClusterAuth  *Auth
}

func (o *DoOptions) copy() *DoOptions {

	newO := *o

	return &newO

}

func (o *DoOptions) Values() map[string]string {

	values := make(map[string]string)

	if o.InfobaseAuth != nil {

		values["--infobase-user"] = o.InfobaseAuth.User
		values["--infobase-pwd"] = o.InfobaseAuth.Pwd

	}

	if o.ClusterAuth != nil {

		values["--cluster-user"] = o.ClusterAuth.User
		values["--cluster-pwd"] = o.ClusterAuth.Pwd

	}

	if len(o.Cluster) > 0 {
		values["--cluster"] = o.Cluster
	}

	if len(o.Process) > 0 {
		values["--process"] = o.Process
	}

	return values

}

func WithInfobaseAuth(user, pwd string) DoOption {

	return func(o *DoOptions) {
		if len(user) == 0 {
			return
		}

		o.InfobaseAuth = &Auth{
			User: user,
			Pwd:  pwd,
		}
	}
}

func WithProcess(process string) DoOption {

	return func(o *DoOptions) {
		if len(process) == 0 {
			return
		}

		o.Process = process
	}
}

func WithAuth(user, pwd string) DoOption {

	return func(o *DoOptions) {
		if len(user) == 0 {
			return
		}

		o.ClusterAuth = &Auth{
			User: user,
			Pwd:  pwd,
		}
	}
}

func WithContext(ctx context.Context) Option {

	return func(o *Options) {
		if ctx == nil {
			return
		}

		o.ctx = ctx
	}

}

func WithClusterAuth(user, pwd string) Option {

	return func(o *Options) {
		if len(user) == 0 {
			return
		}

		o.clusterAuth = Auth{
			User: user,
			Pwd:  pwd,
		}
	}

}

func WithTimeout(timeout time.Duration) Option {

	return func(o *Options) {
		if timeout == 0 {
			return
		}

		o.timeout = timeout
	}

}

func WithNoUpdate() Option {

	return func(o *Options) {
		o.updateInterval = 0
	}

}

func WithUpdate(duration time.Duration) Option {

	return func(o *Options) {
		o.updateInterval = duration
	}
}

func WithV8Version(v8version string) Option {

	return func(o *Options) {
		if len(v8version) == 0 {
			return
		}

		o.v8Version = v8version
	}

}

func WithPath(path string) Option {

	return func(o *Options) {
		if len(path) == 0 {
			return
		}

		o.racPath = path
	}

}
