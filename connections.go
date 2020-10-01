package rac

const (
	ConnectionsCommand           = "connection"
	ConnectionsListCommand       = ConnectionsCommand + " list"
	ConnectionsInfoCommand       = ConnectionsCommand + " info"
	ConnectionsDisconnectCommand = ConnectionsCommand + " disconnect"
)

type InfobaseSig interface {
	Sig() (uuid string, auth AuthSig)
}

type AuthSig interface {
	Sig() (usr string, pwd string)
}

type ConnectionSig interface {
	Sig() (string, string, string)
}

type ConnectionSigFilter interface {
	Filter() (process string, infobase InfobaseSig, filterFunc func(info ConnectionInfo) bool)
}

type ConnectionsList struct {
	Process    string
	Infobase   InfobaseSig
	FilterFunc func(info ConnectionInfo) bool
}

func (_ ConnectionsList) Command() string {
	return ConnectionsListCommand
}

func (i ConnectionsList) Values() map[string]string {

	var (
		ib, usr, pwd string
		auth         AuthSig
	)

	if i.Infobase != nil {
		ib, auth = i.Infobase.Sig()
		if auth != nil {
			usr, pwd = auth.Sig()
		}
	}

	return map[string]string{
		"--process":       i.Process,
		"--infobase":      ib,
		"--infobase-user": usr,
		"--infobase-pwd":  pwd,
	}

}

func (i ConnectionsList) Parse(res *RawRespond) error {

	var list []ConnectionInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &list)
	res.Error = err
	res.parsedRespond = list

	if i.FilterFunc != nil {

	}

	return err

}

func (i ConnectionsList) filter(list []ConnectionInfo) []ConnectionInfo {

	if i.FilterFunc == nil {
		return list
	}

	var filtered []ConnectionInfo

	for _, info := range list {

		if i.FilterFunc(info) {
			filtered = append(filtered, info)
		}

	}

	return filtered

}

type ConnectionsInfo struct {
	UUID string
}

func (i ConnectionsInfo) Values() map[string]string {

	return map[string]string{
		"--connection": i.UUID,
	}

}

func (i ConnectionsInfo) Parse(res *RawRespond) error {

	var info ConnectionInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ ConnectionsInfo) Command() string {
	return ConnectionsInfoCommand
}

type ConnectionsDisconnect struct {
	UUID     string
	Process  string
	Infobase string
	Auth
}

func (_ ConnectionsDisconnect) Command() string {
	return ConnectionsDisconnectCommand
}

func (i ConnectionsDisconnect) Values() map[string]string {

	return map[string]string{
		"--connection":    i.UUID,
		"--process":       i.Process,
		"--infobase":      i.Infobase,
		"--infobase-user": i.User,
		"--infobase-pwd":  i.Pwd,
	}

}

func (i ConnectionsDisconnect) Parse(res *RawRespond) error {

	if !res.Status {
		return res.Error
	}

	return nil
}

type ConnectionsRespond struct {
	*RawRespond
	List []ConnectionInfo
	Info ConnectionInfo
}

func (m *Manager) Connections(cluster Clusterable, what interface{}, opts ...interface{}) (ConnectionsRespond, error) {

	var (
		method  string
		params  = make(map[string]string)
		parser  RespondParser
		respond ConnectionsRespond
	)

	switch v := what.(type) {
	case valued:

		method = v.Command()
		params = v.Values()
		parser = v.(RespondParser)

	default:
		return respond, ErrUnsupportedWhat
	}

	doOptions := extractOptions(opts)

	raw := m.do(method, params, clusterSigParams(cluster.ClusterSig()), doOptions.Values())

	if raw.Error != nil {
		return respond, raw.Error
	}

	err := parser.Parse(raw)

	if err != nil {
		return respond, ErrUnsupportedWhat
	}

	respond.RawRespond = raw

	switch v := raw.parsedRespond.(type) {

	case ConnectionInfo:
		respond.Info = v
		respond.List = append(respond.List, v)
	case []ConnectionInfo:
		respond.List = v
	}

	return respond, nil

}

func (m *Manager) DisconnectConnection(cluster Clusterable, what ConnectionSig, opts ...interface{}) (ConnectionsRespond, error) {

	do := ConnectionsDisconnect{}
	do.UUID, do.Process, do.Infobase = what.Sig()

	return m.Connections(cluster, what, opts)

}

func (m *Manager) ConnectionInfo(cluster Clusterable, what ConnectionSig, opts ...interface{}) (ConnectionsRespond, error) {

	do := ConnectionsInfo{}
	do.UUID, _, _ = what.Sig()

	return m.Connections(cluster, what, opts)

}

func (m *Manager) ConnectionList(cluster Clusterable, filter ConnectionSigFilter, opts ...interface{}) (ConnectionsRespond, error) {

	do := ConnectionsList{}
	do.Process, do.Infobase, do.FilterFunc = filter.Filter()

	return m.Connections(cluster, do, opts)

}
