package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

type InfobaseSummaryList []InfobaseSummaryInfo

func (l InfobaseSummaryList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &InfobaseSummaryInfo{}
		info.Parse(decoder, version, r)

		l = append(l, *info)
	}
}

type InfobaseSummaryInfo struct {
	UUID        uuid.UUID `rac:"infobase"` //infobase : efa3672f-947a-4d84-bd58-b21997b83561
	Name        string    //name     : УППБоеваяБаза
	Description string    `rac:"descr"` //descr    : "УППБоеваяБаза"

}

func (i InfobaseSummaryInfo) InfobaseSig() string {
	return i.UUID.String()
}

func (i *InfobaseSummaryInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Description, r)
	decoder.StringPtr(&i.Name, r)

}

type InfobaseInfo struct {
	UUID                                   uuid.UUID `rac:"infobase"` //infobase : efa3672f-947a-4d84-bd58-b21997b83561
	Name                                   string    //name     : УППБоеваяБаза
	Description                            string    `rac:"descr"` //descr    : "УППБоеваяБаза"
	Dbms                                   string    //dbms                                       : MSSQLServer
	DbServer                               string    //db-server                                  : sql
	DbName                                 string    //db-name                                    : base
	DbUser                                 string    //db-user                                    : user
	DbPwd                                  string    `rac:"="` //--db-pwd=<pwd>  пароль администратора базы данных
	SecurityLevel                          int       //security-level                             : 0
	LicenseDistribution                    int       //license-distribution                       : allow
	ScheduledJobsDeny                      bool      //scheduled-jobs-deny                        : off
	SessionsDeny                           bool      //sessions-deny                              : off
	DeniedFrom                             time.Time //denied-from                                :
	DeniedMessage                          string    //denied-message                             : "Выполняется обновление базы"
	DeniedParameter                        string    //denied-parameter                           :
	DeniedTo                               time.Time //denied-to                                  :
	PermissionCode                         string    //permission-code                            : "123"
	ExternalSessionManagerConnectionDtring string    //external-session-manager-connection-string :
	ExternalSessionManagerRequired         bool      //external-session-manager-required          : no
	SecurityProfileName                    string    //security-profile-name                      :
	SafeModeSecurityProfileName            string    //safe-mode-security-profile-name            :
	ReserveWorkingProcesses                bool      //reserve-working-processes                  : no
	DateOffset                             int
	Locale                                 string
}

func (i *InfobaseInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.IntPtr(&i.DateOffset, r)
	decoder.StringPtr(&i.Dbms, r)
	decoder.StringPtr(&i.DbName, r)
	decoder.StringPtr(&i.DbPwd, r)
	decoder.StringPtr(&i.DbServer, r)
	decoder.StringPtr(&i.DbUser, r)
	decoder.TimePtr(&i.DeniedFrom, r)
	decoder.StringPtr(&i.DeniedMessage, r)
	decoder.StringPtr(&i.DeniedParameter, r)
	decoder.TimePtr(&i.DeniedTo, r)
	decoder.StringPtr(&i.Description, r)
	decoder.StringPtr(&i.Locale, r)
	decoder.StringPtr(&i.Name, r)
	decoder.StringPtr(&i.PermissionCode, r)
	decoder.BoolPtr(&i.ScheduledJobsDeny, r)
	decoder.IntPtr(&i.SecurityLevel, r)
	decoder.BoolPtr(&i.SessionsDeny, r)
	decoder.IntPtr(&i.LicenseDistribution, r)
	decoder.StringPtr(&i.ExternalSessionManagerConnectionDtring, r)
	decoder.BoolPtr(&i.ExternalSessionManagerRequired, r)
	decoder.StringPtr(&i.SecurityProfileName, r)
	decoder.StringPtr(&i.SafeModeSecurityProfileName, r)
	if version > 9 {
		decoder.BoolPtr(&i.ReserveWorkingProcesses, r)
	}

}
func (i InfobaseInfo) Format(encoder Encoder, version int, w io.Writer) {

	encoder.Uuid(i.UUID, w)
	encoder.Int(i.DateOffset, w)
	encoder.String(i.Dbms, w)
	encoder.String(i.DbName, w)
	encoder.String(i.DbPwd, w)
	encoder.String(i.DbServer, w)
	encoder.String(i.DbUser, w)
	encoder.Time(i.DeniedFrom, w)
	encoder.String(i.DeniedMessage, w)
	encoder.String(i.DeniedParameter, w)
	encoder.Time(i.DeniedTo, w)
	encoder.String(i.Description, w)
	encoder.String(i.Locale, w)
	encoder.String(i.Name, w)
	encoder.String(i.PermissionCode, w)
	encoder.Bool(i.ScheduledJobsDeny, w)
	encoder.Int(i.SecurityLevel, w)
	encoder.Bool(i.SessionsDeny, w)
	encoder.Int(i.LicenseDistribution, w)
	encoder.String(i.ExternalSessionManagerConnectionDtring, w)
	encoder.Bool(i.ExternalSessionManagerRequired, w)
	encoder.String(i.SecurityProfileName, w)
	encoder.String(i.SafeModeSecurityProfileName, w)
	if version > 9 {
		encoder.Bool(i.ReserveWorkingProcesses, w)
	}

}

func (i InfobaseInfo) Summary() InfobaseSummaryInfo {
	return InfobaseSummaryInfo{
		UUID:        i.UUID,
		Name:        i.Name,
		Description: i.Description,
	}

}

func (i InfobaseInfo) InfobaseSig() string {
	return i.UUID.String()
}
