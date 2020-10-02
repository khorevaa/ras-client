package rac

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type InfobasesCommandType string

func (c InfobasesCommandType) Check(params map[string]string) error {

	errs := &multierror.Error{}

	if val, ok := params["--cluster"]; !ok || len(val) == 0 {

		_ = multierror.Append(errs, errors.New("cluster must be identified"))
	}

	switch c {
	case InfobasesCreateCommand:

		if val, ok := params["--name"]; !ok || len(val) == 0 {
			_ = multierror.Append(errs, errors.New("name must be identified"))
		}

		if val, ok := params["--dbms"]; !ok || len(val) == 0 {
			_ = multierror.Append(errs, errors.New("dbms must be identified"))
		}

		if val, ok := params["--db-server"]; !ok || len(val) == 0 {
			_ = multierror.Append(errs, errors.New("server must be identified"))
		}

		if val, ok := params["--db-name"]; !ok || len(val) == 0 {
			_ = multierror.Append(errs, errors.New("db name must be identified"))
		}

		if val, ok := params["--ocale"]; !ok || len(val) == 0 {
			_ = multierror.Append(errs, errors.New("locale must be identified"))
		}

	case InfobasesSummaryInfoCommand, InfobasesSummaryUpdateCommand,
		InfobasesDropCommand, InfobasesUpdateCommand, InfobasesInfoCommand:

		if val, ok := params["--infobase"]; !ok || len(val) == 0 {
			_ = multierror.Append(errs, errors.New("infobase must be identified"))
		}
	}

	return errs.ErrorOrNil()
}

func (c InfobasesCommandType) Command() string {
	return string(c)
}

const (
	baseInfobasesCommand          InfobasesCommandType = "infobase"
	infobasesSummaryCommand                            = baseInfobasesCommand + " summary"
	InfobasesListCommand                               = infobasesSummaryCommand + " list"
	InfobasesSummaryInfoCommand                        = infobasesSummaryCommand + " info"
	InfobasesSummaryUpdateCommand                      = infobasesSummaryCommand + " update"
	InfobasesInfoCommand                               = baseInfobasesCommand + " info"
	InfobasesCreateCommand                             = baseInfobasesCommand + " create"
	InfobasesDropCommand                               = baseInfobasesCommand + " drop"
	InfobasesUpdateCommand                             = baseInfobasesCommand + " update"
)

type InfobasesList struct{}

func (_ InfobasesList) Command() DoCommand {
	return InfobasesListCommand
}

func (i InfobasesList) Values() map[string]string {

	return map[string]string{}

}

func (i InfobasesList) Parse(res *RawRespond) error {

	var list []InfobaseSummaryInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &list)
	res.Error = err
	res.parsedRespond = list

	return err

}

type InfobasesSummaryInfo struct {
	UUID string
}

func (i InfobasesSummaryInfo) Values() map[string]string {

	return map[string]string{
		"--infobase": i.UUID,
	}

}

func (i InfobasesSummaryInfo) Parse(res *RawRespond) error {

	var info InfobaseSummaryInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ InfobasesSummaryInfo) Command() DoCommand {
	return InfobasesSummaryInfoCommand
}

type InfobaseFullInfo struct {
	UUID string
	Auth Auth
}

func (i InfobaseFullInfo) Values() map[string]string {

	return map[string]string{
		"--infobase":      i.UUID,
		"--infobase-user": i.Auth.User,
		"--infobase-pwd":  i.Auth.Pwd,
	}

}

func (i InfobaseFullInfo) Parse(res *RawRespond) error {

	var info InfobaseInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ InfobaseFullInfo) Command() DoCommand {
	return InfobasesInfoCommand
}

type InfobaseDrop struct {
	UUID          string
	Auth          Auth
	DropDatabase  bool
	ClearDatabase bool
}

func (i InfobaseDrop) Values() map[string]string {

	val := map[string]string{
		"--infobase":      i.UUID,
		"--infobase-user": i.Auth.User,
		"--infobase-pwd":  i.Auth.Pwd,
	}

	switch {
	case i.DropDatabase:
		val["--drop-database"] = strconv.FormatBool(i.DropDatabase)

	case i.ClearDatabase:
		val["--clear-database"] = strconv.FormatBool(i.ClearDatabase)

	}

	return val
}

func (i InfobaseDrop) Parse(res *RawRespond) error {

	var info interface{}

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ InfobaseDrop) Command() DoCommand {
	return InfobasesDropCommand
}

const (
	LicenseDistributionAllow = "allow"
	LicenseDistributionDeni  = "deni"
	DbmsMSSQLServer          = "MSSQLServer"
	DbmsPostgreSQL           = "PostgreSQL"
	DbmsIBMDB2               = "IBMDB2"
	DbmsOracleDatabase       = "OracleDatabase"
)

type YasNoBoolType int

const (
	NullBool YasNoBoolType = iota
	NoBool
	YesBool
)

func (t YasNoBoolType) Bool() bool {

	switch t {

	case NullBool:
		return false
	case YesBool:
		return true
	case NoBool:
		return false
	}
	return false
}

func (t YasNoBoolType) FromBool(v bool) YasNoBoolType {

	b := NoBool
	if v {
		b = YesBool
	}

	return b
}

func (t YasNoBoolType) String() string {

	switch t {

	case NullBool:
		return ""
	case YesBool:
		return "yes"
	case NoBool:
		return "no"
	}

	return ""

}

type InfobaseUpdate struct {
	UUID     string `rac:"infobase"` //infobase : efa3672f-947a-4d84-bd58-b21997b83561
	Dbms     string //dbms                                       : MSSQLServer
	DbServer string //db-server                                  : sql
	DbName   string //db-name                                    : base
	DbUser   string //db-user                                    : user
	DbPwd    string //db-pwd    пароль администратора базы данных
	//SecurityLevel                          int       //security-level                             : 0
	LicenseDistribution                    string        //license-distribution                       : allow
	ScheduledJobsDeny                      YasNoBoolType //scheduled-jobs-deny                        : off
	SessionsDeny                           YasNoBoolType //sessions-deny                              : off
	DeniedFrom                             time.Time     //denied-from                                :
	DeniedMessage                          string        //denied-message                             : "Выполняется обновление базы"
	DeniedParameter                        string        //denied-parameter                           :
	DeniedTo                               time.Time     //denied-to                                  :
	PermissionCode                         string        //permission-code                            : "123"
	ExternalSessionManagerConnectionDtring string        //external-session-manager-connection-string :
	ExternalSessionManagerRequired         YasNoBoolType //external-session-manager-required          : no
	SecurityProfileName                    string        //security-profile-name                      :
	SafeModeSecurityProfileName            string        //safe-mode-security-profile-name            :
	ReserveWorkingProcesses                YasNoBoolType //reserve-working-processes                  : no

	Auth Auth
}

func (i InfobaseUpdate) Values() map[string]string {

	val := map[string]string{
		"--infobase-user": i.Auth.User,
		"--infobase-pwd":  i.Auth.Pwd,
	}

	rv := reflect.ValueOf(&i)
	ri := reflect.Indirect(rv)

	rt := reflect.TypeOf(i)

	for i := 0; i < ri.NumField(); i++ {

		fieldName := NameMapping(rt.Field(i).Name)

		tag := rt.Field(i).Tag.Get(TagNamespace)

		tags := strings.Split(tag, ",")

		if len(tags) > 0 && len(tags[0]) > 0 {
			fieldName = tags[0]
		}

		if tags[0] == "-" {
			continue
		}

		paramName := "--" + fieldName

		value := ri.Field(i).Interface()

		switch v := value.(type) {
		case YasNoBoolType:
			if v == NullBool {
				continue
			}

			val[paramName] = v.String()

		case int:
			if v == -1 {
				continue
			}
			val[paramName] = fmt.Sprintf("%d", v)

		case time.Time:
			if v.IsZero() {
				continue
			}
			val[paramName] = v.Format("2006-01-02T15:04:05")
		case string:
			if len(v) == 0 {
				continue
			}
			val[paramName] = v
		}

	}

	return val
}

func (i InfobaseUpdate) Parse(res *RawRespond) error {

	var info InfobaseInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ InfobaseUpdate) Command() DoCommand {
	return InfobasesUpdateCommand
}

type InfobaseUpdateDescription struct {
	UUID        string
	Description string
	Auth        Auth
}

func (i InfobaseUpdateDescription) Values() map[string]string {

	val := map[string]string{

		"--infobase":      i.UUID,
		"--infobase-user": i.Auth.User,
		"--infobase-pwd":  i.Auth.Pwd,
		"--descr":         i.Description,
	}

	return val
}

func (i InfobaseUpdateDescription) Parse(res *RawRespond) error {

	var info InfobaseSummaryInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ InfobaseUpdateDescription) Command() DoCommand {
	return InfobasesSummaryUpdateCommand
}

type InfobaseCreate struct {
	Name                string //name     : УППБоеваяБаза
	CreateDatabase      bool
	Description         string `rac:"descr"` //descr    : "УППБоеваяБаза"
	Dbms                string //dbms                                       : MSSQLServer
	DbServer            string //db-server                                  : sql
	DbName              string //db-name                                    : base
	Locale              string
	DbUser              string //db-user                                    : user
	DbPwd               string //db-pwd    пароль администратора базы данных
	DataOffset          int
	SecurityLevel       int           //security-level                             : 0
	ScheduledJobsDeny   YasNoBoolType //scheduled-jobs-deny                        : off
	LicenseDistribution string        //license-distribution                       : allow

}

func (i InfobaseCreate) Values() map[string]string {

	val := map[string]string{}

	rv := reflect.ValueOf(&i)
	ri := reflect.Indirect(rv)

	rt := reflect.TypeOf(i)

	for i := 0; i < ri.NumField(); i++ {

		fieldName := NameMapping(rt.Field(i).Name)

		tag := rt.Field(i).Tag.Get(TagNamespace)

		tags := strings.Split(tag, ",")

		if len(tags) > 0 && len(tags[0]) > 0 {
			fieldName = tags[0]
		}

		if tags[0] == "-" {
			continue
		}

		paramName := "--" + fieldName

		value := ri.Field(i).Interface()

		switch v := value.(type) {
		case YasNoBoolType:
			if v == NullBool {
				continue
			}

			val[paramName] = v.String()

		case int:
			if v == 0 {
				continue
			}
			val[paramName] = fmt.Sprintf("%d", v)

		case bool:
			if !v {
				continue
			}
			val[paramName] = strconv.FormatBool(v)
		case string:
			if len(v) == 0 {
				continue
			}
			val[paramName] = v
		}

	}

	return val
}

func (i InfobaseCreate) Parse(res *RawRespond) error {

	var info InfobaseInfo

	if !res.Status {
		return res.Error
	}

	err := Unmarshal(res.raw, &info)
	res.Error = err
	res.parsedRespond = info
	return err
}

func (_ InfobaseCreate) Command() DoCommand {
	return InfobasesCreateCommand
}

type InfobasesRespond struct {
	*RawRespond
	List        []InfobaseSummaryInfo
	Info        InfobaseInfo
	SummaryInfo InfobaseSummaryInfo
}

func (m *Manager) Infobases(what interface{}, opts ...interface{}) (respond InfobasesRespond, err error) {

	val, ok := what.(Valued)

	if !ok {
		return respond, ErrUnsupportedWhat
	}

	respond.RawRespond, err = m.Do(val, opts...)

	if err != nil {
		return respond, err
	}

	switch v := respond.parsedRespond.(type) {

	case InfobaseInfo:
		respond.Info = v
		respond.List = append(respond.List, v.Summary())
	case []InfobaseSummaryInfo:
		respond.List = v
	case InfobaseSummaryInfo:
		respond.SummaryInfo = v
	}

	return respond, nil

}

func (m *Manager) InfobasesList(opts ...interface{}) (list []InfobaseSummaryInfo, err error) {

	do := InfobasesList{}
	respond, err := m.Infobases(do, opts...)

	if err != nil {
		return
	}

	return respond.List, nil

}

func (m *Manager) InfobaseByName(name string, opts ...interface{}) (respond InfobaseSummaryInfo, err error) {

	list, err := m.InfobasesList(opts...)

	if err != nil {
		return
	}

	for _, info := range list {
		if strings.EqualFold(info.Name, name) {
			respond = info
			break
		}
	}

	return respond, err

}

func (m *Manager) DropInfobase(do InfobaseDrop, opts ...interface{}) (err error) {

	_, err = m.Infobases(do, opts...)

	return
}

func (m *Manager) InfobaseInfo(do InfobaseFullInfo, opts ...interface{}) (respond InfobaseInfo, err error) {

	resp, err := m.Infobases(do, opts...)
	if err != nil {
		return
	}

	respond = resp.Info

	return
}

func (m *Manager) UpdateInfobase(changes interface{}, opts ...interface{}) (respond InfobaseInfo, err error) {

	var do InfobaseUpdate

	switch u := changes.(type) {

	case InfobaseUpdate:
		do = u
	case *InfobaseUpdate:
		do = *u
	case *InfobaseInfo, InfobaseInfo:

		var update InfobaseInfo

		sig := u.(InfobaseSig)

		getFullInfo := InfobaseFullInfo{
			UUID: sig.InfobaseSig(),
		}

		ibInfo, err := m.InfobaseInfo(getFullInfo, opts...)
		if err != nil {
			return respond, err
		}

		switch t := u.(type) {
		case InfobaseInfo:
			update = t
		case *InfobaseInfo:
			update = *t
		}

		do = ibInfo.UpdateChanges(update)

	}

	res, err := m.Infobases(do, opts...)
	if err != nil {
		return
	}

	respond = res.Info
	return respond, err
}

func (m *Manager) UpdateDescription(do InfobaseUpdateDescription, opts ...interface{}) (respond InfobaseSummaryInfo, err error) {

	res, err := m.Infobases(do, opts...)
	if err != nil {
		return
	}

	respond = res.SummaryInfo
	return
}

func (m *Manager) CreateInfobase(do InfobaseCreate, opts ...interface{}) (respond InfobaseInfo, err error) {

	res, err := m.Infobases(do, opts...)
	if err != nil {
		return
	}

	respond = res.Info
	return

}
