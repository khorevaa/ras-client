package rac

const (
	ClustersCommand       = "cluster"
	ClustersListCommand   = ClustersCommand + " list"
	ClustersInfoCommand   = ClustersCommand + " info"
	ClustersInsertCommand = ClustersCommand + " insert"
	ClustersRemoveCommand = ClustersCommand + " remove"
	ClustersUpdateCommand = ClustersCommand + " update"
)

type ClustersList struct {
}

func (_ ClustersList) Command() string {
	return ClustersListCommand
}

func (i ClustersList) Values() map[string]string {

	return map[string]string{}

}

func (i ClustersList) Parse(res *RawRespond) error {

	var list []ClusterInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &list)
	res.Error = err
	res.parsedRespond = list
	return err

}

type ClustersInfo struct {
	UUID string
}

func (i ClustersInfo) Values() map[string]string {

	return map[string]string{
		"--cluster": i.UUID,
	}

}

func (i ClustersInfo) Parse(res *RawRespond) error {

	var info ClusterInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ ClustersInfo) Command() string {
	return ClustersInfoCommand
}

type ClustersUpdate struct {
	ClusterInfo
	Auth
}

func (_ ClustersUpdate) Command() string {
	return ClustersUpdateCommand
}

func (i ClustersUpdate) Values() map[string]string {

	return map[string]string{
		"--cluster":    i.UUID,
		"--agent-user": i.User,
		"--agent-pwd":  i.Pwd,
	}

}

func (i ClustersUpdate) Parse(res *RawRespond) error {

	if !res.Status {
		return res.Error
	}

	return nil
}

type ClustersInsert struct {
	ClusterInfo
	Auth
}

func (_ ClustersInsert) Command() string {
	return ClustersInsertCommand
}

func (i ClustersInsert) Values() map[string]string {

	return map[string]string{
		"--cluster":    i.UUID,
		"--agent-user": i.User,
		"--agent-pwd":  i.Pwd,
	}

}

func (i ClustersInsert) Parse(res *RawRespond) error {

	if !res.Status {
		return res.Error
	}

	return nil
}

type ClustersRemove struct {
	ClusterInfo
	Auth
}

func (_ ClustersRemove) Command() string {
	return ClustersRemoveCommand
}

func (i ClustersRemove) Values() map[string]string {

	return map[string]string{
		"--cluster":      i.UUID,
		"--cluster-user": i.User,
		"--cluster-pwd":  i.Pwd,
	}

}

func (i ClustersRemove) Parse(res *RawRespond) error {

	if !res.Status {
		return res.Error
	}

	return nil
}

type ClustersRespond struct {
	*RawRespond
	List []ClusterInfo
	Info ClusterInfo
}

func (m *Manager) Clusters(what interface{}, opts ...interface{}) (ClustersRespond, error) {

	var (
		method  string
		params  = make(map[string]string)
		parser  RespondParser
		respond ClustersRespond
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

	raw := m.do(method, params, doOptions.Values())

	if raw.Error != nil {
		return respond, raw.Error
	}

	err := parser.Parse(raw)

	if err != nil {
		return respond, ErrUnsupportedWhat
	}

	respond.RawRespond = raw

	switch v := raw.parsedRespond.(type) {

	case ClusterInfo:
		respond.Info = v
		respond.List = append(respond.List, v)
	case []ClusterInfo:
		respond.List = v
	}

	return respond, nil

}
