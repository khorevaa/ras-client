package rac

import "errors"

type ClustersCommandType string

func (c ClustersCommandType) Check(params map[string]string) error {

	var err error

	switch c {

	case ClustersInfoCommand:

		if val, ok := params["--cluster"]; !ok || len(val) == 0 {
			err = errors.New("cluster must be identified")
		}

	case ClustersRemoveCommand, ClustersUpdateCommand:

		if val, ok := params["--cluster"]; !ok || len(val) == 0 {
			err = errors.New("cluster must be identified")
		}

	}

	return err
}

func (c ClustersCommandType) Command() string {
	return string(c)
}

const (
	baseClustersCommand   ClustersCommandType = "cluster"
	ClustersListCommand                       = baseClustersCommand + " list"
	ClustersInfoCommand                       = baseClustersCommand + " info"
	ClustersInsertCommand                     = baseClustersCommand + " insert"
	ClustersRemoveCommand                     = baseClustersCommand + " remove"
	ClustersUpdateCommand                     = baseClustersCommand + " update"
)

type ClustersList struct {
}

func (_ ClustersList) Command() DoCommand {
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

func (_ ClustersInfo) Command() DoCommand {
	return ClustersInfoCommand
}

type ClustersUpdate struct {
	ClusterInfo
	Auth
}

func (_ ClustersUpdate) Command() DoCommand {
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

func (_ ClustersInsert) Command() DoCommand {
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

func (_ ClustersRemove) Command() DoCommand {
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

func (m *Manager) Clusters(what interface{}, opts ...interface{}) (respond ClustersRespond, err error) {

	val, ok := what.(Valued)

	if !ok {
		return respond, ErrUnsupportedWhat
	}

	respond.RawRespond, err = m.Do(val, opts...)

	if err != nil {
		return respond, err
	}

	switch v := respond.parsedRespond.(type) {

	case ClusterInfo:
		respond.Info = v
		respond.List = append(respond.List, v)
	case []ClusterInfo:
		respond.List = v
	}

	return respond, nil

}
