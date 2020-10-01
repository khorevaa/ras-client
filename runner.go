package rac

import (
	"errors"
)

type DoCommand interface {
	Check(params map[string]string) error
	Command() string
}

type CommonCommandType string

func (_ CommonCommandType) Check(_ map[string]string) error {
	return nil
}
func (c CommonCommandType) Command() string {
	return string(c)
}

var ErrUnsupportedWhat = errors.New("unsupported what argument")

type ConnectionSig interface {
	ConnectionSig() (uuid string, process string, infobase string)
}

type ConnectionsSig interface {
	ConnectionsSig() (process string, infobase string, auth Auth)
}
type InfobaseSig interface {
	InfobaseSig() (uuid string)
}

type InfobaseAuth interface {
	InfobaseAuth() (usr string, pwd string)
}

type ClusterSig interface {
	ClusterSig() (uuid string)
}

type ClusterAuth interface {
	ClusterAuth() (user string, pwd string)
}

type AuthSig interface {
	Auth() (user string, pwd string)
}

type Valued interface {
	Values() map[string]string
	Command() DoCommand
	RespondParser
}

type RespondParser interface {
	Parse(raw *RawRespond) error
}

type Auth struct {
	User string
	Pwd  string
}

func (a Auth) Sig() (usr string, pwd string) {
	return a.User, a.Pwd
}

func clusterSigParams(sig ClusterSig) map[string]string {

	if sig == nil {
		return make(map[string]string)
	}

	return map[string]string{
		"--cluster": sig.ClusterSig(),
	}
}

func infobaseSigParams(sig InfobaseSig) map[string]string {

	if sig == nil {
		return make(map[string]string)
	}

	return map[string]string{
		"--infobase": sig.InfobaseSig(),
	}
}

func clusterAuthParams(auth ClusterAuth) map[string]string {

	if auth == nil {
		return make(map[string]string)
	}

	user, pwd := auth.ClusterAuth()

	return map[string]string{
		"--cluster-user": user,
		"--cluster-pwd":  pwd,
	}
}

func infobaseAuthParams(auth InfobaseAuth) map[string]string {

	if auth == nil {
		return make(map[string]string)
	}

	user, pwd := auth.InfobaseAuth()

	return map[string]string{
		"--infobase-user": user,
		"--infobase-pwd":  pwd,
	}
}

type RawRespond struct {
	Status        bool
	raw           []byte
	parsedRespond interface{}
	Error         error
}

func newRawRespond(data []byte, err error) *RawRespond {

	res := &RawRespond{
		raw:    data,
		Error:  err,
		Status: true,
	}

	if err != nil {
		res.Status = false
	}

	return res
}

func extractOptions(how []interface{}) *DoOptions {

	var opts DoOptions

	for _, prop := range how {

		switch opt := prop.(type) {
		case *DoOptions:
			opts = *opt.copy()
		case DoOptions:
			opts = *opt.copy()
		case DoOption:
			opt(&opts)

		default:
			panic("unsupported doOption")
		}
	}

	return &opts
}

func extractParams(props []interface{}) (params map[string]string) {

	for _, prop := range props {

		clusterSig := prop.(ClusterSig)
		params = mergeParams(params, clusterSigParams(clusterSig))
		infobaseSig := prop.(InfobaseSig)
		params = mergeParams(params, infobaseSigParams(infobaseSig))
		clusterAuth := prop.(ClusterAuth)
		params = mergeParams(params, clusterAuthParams(clusterAuth))
		infobaseAuth := prop.(InfobaseAuth)
		params = mergeParams(params, infobaseAuthParams(infobaseAuth))

	}

	return params
}

func (m *Manager) embedParams() (params map[string]string) {

	params = mergeParams(params, clusterSigParams(m))
	params = mergeParams(params, clusterAuthParams(m))

	return params

}
