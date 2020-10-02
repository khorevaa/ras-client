package rac

import (
	"errors"
	"time"
)

type ManagerOptions struct {
	V8Version       string
	RacPath         string
	UpdateInterval  time.Duration
	Timeout         time.Duration
	TryTimeoutCount int
	DetectCluster   bool
	ClusterAuth     struct {
		User string
		Pwd  string
	}
	Logger Logger
	Runner Runner
}

type Option func(o *Options)

type Options struct {
	Cluster  string
	Process  string
	Infobase string

	InfobaseAuth *Auth
	ClusterAuth  *Auth
}

func (o *Options) Options(opts ...interface{}) (err error) {

	defer func() {
		if e := recover(); e != nil {
			switch errT := e.(type) {
			case error:
				err = errT
			case string:
				err = errors.New(errT)
			default:
				panic(e)
			}
		}
	}()

	for _, opt := range opts {

		optFn, ok := opt.(Option)

		if ok {
			optFn(o)
		}

	}

	return
}

func (o *Options) copy() *Options {

	newO := *o

	return &newO

}

func (o *Options) Values() map[string]string {

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

func WithInfobaseAuth(user, pwd string) Option {

	return func(o *Options) {
		if len(user) == 0 {
			return
		}

		o.InfobaseAuth = &Auth{
			User: user,
			Pwd:  pwd,
		}
	}
}

func WithInfobase(infobase interface{}) Option {

	return func(o *Options) {

		switch i := infobase.(type) {

		case InfobaseSig:
			o.Infobase = i.InfobaseSig()

		case string:
			if len(i) == 0 {
				return
			}

			o.Infobase = i
		default:
			panic(errors.New("unsupported infobase type"))
		}

	}
}

func WithProcess(process string) Option {

	return func(o *Options) {
		if len(process) == 0 {
			return
		}

		o.Process = process
	}
}

func WithAuth(user, pwd string) Option {

	return func(o *Options) {
		if len(user) == 0 {
			return
		}

		o.ClusterAuth = &Auth{
			User: user,
			Pwd:  pwd,
		}
	}
}
