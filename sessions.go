package rac

import (
	"errors"
	"strconv"
)

type SessionCommandType string

func (c SessionCommandType) Check(params map[string]string) error {

	var err error

	if val, ok := params["--cluster"]; !ok || len(val) == 0 {
		return errors.New("cluster must be identified")
	}

	switch c {

	case SessionsInfoCommand, SessionsTerminateCommand:

		if val, ok := params["--session"]; !ok || len(val) == 0 {
			err = errors.New("session must be identified")
		}

	}

	return err
}

func (c SessionCommandType) Command() string {
	return string(c)
}

const (
	baseSessionsCommand      SessionCommandType = "session"
	SessionsListCommand                         = baseSessionsCommand + " list"
	SessionsInfoCommand                         = baseSessionsCommand + " info"
	SessionsTerminateCommand                    = baseSessionsCommand + " terminate"
)

type SessionFilterFunc func(info SessionsInfo) bool

type SessionsList struct {
	Infobase   string
	Licenses   bool
	FilterFunc SessionFilterFunc
}

func (_ SessionsList) Command() DoCommand {
	return SessionsListCommand
}

func (i SessionsList) Values() map[string]string {

	return map[string]string{
		"--infobase": i.Infobase,
		"--licenses": strconv.FormatBool(i.Licenses),
	}
}

func (i SessionsList) Parse(res *RawRespond) error {

	var err error
	var list interface{}

	list = []SessionsInfo{}

	if i.Licenses {
		list = []LicenseInfo{}
	}

	if !res.Status {
		return res.Error
	}

	err = Unmarshal(res.raw, &list)

	if l, ok := list.([]SessionsInfo); ok {
		list = i.filter(l)
	}

	res.Error = err
	res.parsedRespond = list
	return err

}

func (i SessionsList) filter(list []SessionsInfo) []SessionsInfo {

	if i.FilterFunc == nil {
		return list
	}

	var filtered []SessionsInfo

	for _, info := range list {

		if i.FilterFunc(info) {
			filtered = append(filtered, info)
		}

	}

	return filtered

}

func (i SessionsList) extractOptions(props []interface{}) {
	for _, prop := range props {

		switch opt := prop.(type) {
		case SessionFilterFunc:
			i.FilterFunc = opt
		default:
			continue
		}
	}
}

type SessionsInfo struct {
	UUID     string
	Licenses bool
}

func (i SessionsInfo) Values() map[string]string {

	return map[string]string{
		"--session":  i.UUID,
		"--licenses": strconv.FormatBool(i.Licenses),
	}

}

func (_ SessionsInfo) Command() DoCommand {
	return SessionsInfoCommand
}

func (i SessionsInfo) Parse(res *RawRespond) error {

	var err error
	var info interface{}

	info = SessionsInfo{}

	if i.Licenses {
		info = LicenseInfo{}
	}

	if !res.Status {
		return res.Error
	}

	err = Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err

}

type SessionsTerminate struct {
	UUID string
}

func (_ SessionsTerminate) Command() DoCommand {
	return SessionsTerminateCommand
}
func (i SessionsTerminate) Values() map[string]string {

	return map[string]string{
		"--session": i.UUID,
	}

}

func (i SessionsTerminate) Parse(res *RawRespond) error {

	var err error
	var info interface{}

	if !res.Status {
		return res.Error
	}

	err = Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err

}

type SessionsRespond struct {
	*RawRespond
	List []SessionsInfo
	Info SessionsInfo
}

type LicensesRespond struct {
	*RawRespond
	List []LicenseInfo
	Info LicenseInfo
}

func (m *Manager) Sessions(what interface{}, opts ...interface{}) (respond SessionsRespond, err error) {

	val, ok := what.(Valued)

	if !ok {
		return respond, ErrUnsupportedWhat
	}

	respond.RawRespond, err = m.Do(val, opts...)

	if err != nil {
		return respond, err
	}

	switch v := respond.parsedRespond.(type) {

	case SessionsInfo:
		respond.Info = v
		respond.List = append(respond.List, v)
	case []SessionsInfo:
		respond.List = v
	}

	return respond, nil

}

func (m *Manager) Licenses(what interface{}, opts ...interface{}) (respond LicensesRespond, err error) {
	val, ok := what.(Valued)

	if !ok {
		return respond, ErrUnsupportedWhat
	}

	respond.RawRespond, err = m.Do(val, opts...)

	if err != nil {
		return respond, err
	}

	switch v := respond.parsedRespond.(type) {
	case LicenseInfo:
		respond.Info = v
		respond.List = append(respond.List, v)
	case []LicenseInfo:
		respond.List = v
	}

	return respond, nil

}

type SessionsSig interface {
	SessionsSig() (uuid string)
}

func (m *Manager) SessionsList(what InfobaseSig, opts ...interface{}) (SessionsRespond, error) {

	do := &SessionsList{}

	if what != nil {
		do.Infobase = what.InfobaseSig()
	}

	do.extractOptions(opts)

	return m.Sessions(*do, opts...)

}

func (m *Manager) LicensesList(what InfobaseSig, opts ...interface{}) (LicensesRespond, error) {

	do := &SessionsList{
		Licenses: true,
	}

	if what != nil {
		do.Infobase = what.InfobaseSig()
	}

	do.extractOptions(opts)

	return m.Licenses(*do, opts...)

}

func (m *Manager) SessionInfo(what SessionsSig, opts ...interface{}) (SessionsRespond, error) {

	do := SessionsInfo{}

	if what != nil {
		do.UUID = what.SessionsSig()
	}

	return m.Sessions(do, opts...)

}

func (m *Manager) LicenseInfo(what SessionsSig, opts ...interface{}) (LicensesRespond, error) {

	do := SessionsInfo{
		Licenses: true,
	}

	if what != nil {
		do.UUID = what.SessionsSig()
	}

	return m.Licenses(do, opts...)

}

func (m *Manager) TerminateSession(what SessionsSig, opts ...interface{}) (SessionsRespond, error) {

	do := SessionsTerminate{}

	if what != nil {
		do.UUID = what.SessionsSig()
	}

	return m.Sessions(do, opts...)

}
