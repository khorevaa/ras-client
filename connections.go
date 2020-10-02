package rac

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
)

type ConnectionsCommandType string

func (c ConnectionsCommandType) Check(params map[string]string) error {

	var err error

	if val, ok := params["--cluster"]; !ok || len(val) == 0 {
		return errors.New("cluster must be identified")
	}

	switch c {
	case ConnectionsInfoCommand:

		if val, ok := params["--connection"]; !ok || len(val) == 0 {
			err = errors.New("connection uuid must be identified")
		}

	case ConnectionsDisconnectCommand:

		if val, ok := params["--connection"]; !ok || len(val) == 0 {
			err = errors.New("connection uuid must be identified")
		}

		if val, ok := params["--process"]; !ok || len(val) == 0 {
			err = errors.New("connection uuid must be identified")
		}

	}

	return err
}

func (c ConnectionsCommandType) Command() string {
	return string(c)
}

const (
	baseConnectionsCommand       ConnectionsCommandType = "connection"
	ConnectionsListCommand                              = baseConnectionsCommand + " list"
	ConnectionsInfoCommand                              = baseConnectionsCommand + " info"
	ConnectionsDisconnectCommand                        = baseConnectionsCommand + " disconnect"
)

type ConnectionFilterFunc func(info ConnectionInfo) bool

type ConnectionsList struct {
	Process    string
	Infobase   string
	Auth       Auth
	FilterFunc ConnectionFilterFunc
}

func (_ ConnectionsList) Command() DoCommand {
	return ConnectionsListCommand
}

func (i ConnectionsList) Values() map[string]string {

	user, pwd := i.Auth.Sig()

	return map[string]string{
		"--process":       i.Process,
		"--infobase":      i.Infobase,
		"--infobase-user": user,
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

	if i.FilterFunc != nil {
		list = i.filter(list)
	}

	res.parsedRespond = list

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

func (i *ConnectionsList) extractOptions(props []interface{}) {

	for _, prop := range props {

		switch opt := prop.(type) {
		case ConnectionFilterFunc:
			i.FilterFunc = opt
		default:
			continue
		}
	}

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

func (_ ConnectionsInfo) Command() DoCommand {
	return ConnectionsInfoCommand
}

type ConnectionsDisconnect struct {
	UUID     string
	Process  string
	Infobase string
	Auth     Auth
}

func (_ ConnectionsDisconnect) Command() DoCommand {
	return ConnectionsDisconnectCommand
}

func (i ConnectionsDisconnect) Values() map[string]string {
	user, pwd := i.Auth.Sig()

	return map[string]string{
		"--connection":    i.UUID,
		"--process":       i.Process,
		"--infobase":      i.Infobase,
		"--infobase-user": user,
		"--infobase-pwd":  pwd,
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

func (m *Manager) Connections(what interface{}, opts ...interface{}) (respond ConnectionsRespond, err error) {

	val, ok := what.(Valued)

	if !ok {
		return respond, ErrUnsupportedWhat
	}

	respond.RawRespond, err = m.Do(val, opts...)

	if err != nil {
		return respond, err
	}

	switch v := respond.parsedRespond.(type) {

	case ConnectionInfo:
		respond.Info = v
		respond.List = append(respond.List, v)
	case []ConnectionInfo:
		respond.List = v
	}

	return respond, nil

}

func (m *Manager) DisconnectConnection(what ConnectionSig, opts ...interface{}) (ConnectionsRespond, error) {

	if what == nil {
		return ConnectionsRespond{}, errors.New("connection sig is nil")
	}

	do := ConnectionsDisconnect{}
	do.UUID, do.Process, do.Infobase = what.ConnectionSig()

	return m.Connections(what, opts...)

}

func (m *Manager) ConnectionInfo(what ConnectionSig, opts ...interface{}) (ConnectionsRespond, error) {

	if what == nil {
		return ConnectionsRespond{}, errors.New("connection sig is nil")
	}

	do := ConnectionsInfo{}
	do.UUID, _, _ = what.ConnectionSig()

	return m.Connections(what, opts...)

}

func (m *Manager) ConnectionsList(what ConnectionsSig, opts ...interface{}) (ConnectionsRespond, error) {

	do := &ConnectionsList{}

	if what != nil {
		do.Process, do.Infobase, do.Auth = what.ConnectionsSig()
	}

	do.extractOptions(opts)

	return m.Connections(*do, opts...)

}

func (m *Manager) DisconnectAllConnection(what ConnectionSig, opts ...interface{}) error {

	doList := &ConnectionsList{}
	if what != nil {
		_, doList.Process, doList.Infobase = what.ConnectionSig()
	}

	doList.extractOptions(opts)

	respList, err := m.Connections(*doList, opts...)

	if err != nil {
		return err
	}

	var result *multierror.Error

	list := respList.List

	var errStack []error

	for _, conn := range list {

		do := ConnectionsDisconnect{}
		do.UUID, do.Process, _ = conn.ConnectionSig()

		_, err := m.Connections(what, opts...)

		if err != nil {
			errStack = append(errStack, fmt.Errorf("connection %s process %s err: %v", do.UUID, do.Process, err))
		}

	}

	_ = multierror.Append(result, errStack...)

	return result.ErrorOrNil()

}
