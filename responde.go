package rac

import "time"

type LockInfo struct {
	Connection  string    //connection : 00000000-0000-0000-0000-000000000000
	Session     string    //session    : 8b8a0817-4cb1-4e13-9a8f-472dde1a3b47
	Object      string    //object     : 00000000-0000-0000-0000-000000000000
	Locked      time.Time //locked     : 2020-10-01T08:30:00
	Description string    `rac:"descr"` //descr      : "БД(сеанс ,УППБоеваяБаза,разделяемая)"

}

type SessionInfo struct {
	UUID                          string    `rac:"session"`    // UUID session                          : 1fb5f037-99e8-4924-a99d-a9e687522d32
	ID                            int64     `rac:"session-id"` // ID session-id                       : 1
	Infobase                      string    // Infobase infobase               : aea71760-15b3-485a-9a35-506eb8a0b04a
	Connection                    string    // connection                      : 8adf4514-0379-4333-a153-0b2689edc415
	Process                       string    // process                         : 1af2e54f-d95a-4370-9b45-8277280cad23
	UserName                      string    // user-name                       : АКузнецов
	Host                          string    //host                             : Sport1
	AppId                         string    //app-id                           : Designer
	Locale                        string    //locale                           : ru_RU
	StartedAt                     time.Time //started-at                       : 2018-04-09T14:51:31
	LastActiveAt                  time.Time //last-active-at                   : 2018-05-14T11:12:33
	Hibernate                     bool      // hibernate                        : no
	PassiveSessionHibernateTime   int32     //passive-session-hibernate-time   : 1200
	HibernateDessionTerminateTime int32     //hibernate-session-terminate-time : 86400
	BlockedByDbms                 int64     //blocked-by-dbms                  : 0
	BlockedByLs                   int64     //blocked-by-ls                    : 0
	BytesAll                      int64     //bytes-all                        : 105972550
	BytesLast5min                 int64     `rac:"bytes-last-5min"` //bytes-last-5min                  : 0
	CallsAll                      int64     //calls-all                        : 119052
	CallsLast5min                 int64     `rac:"calls-last-5min"` //calls-last-5min                  : 0
	DbmsBytesAll                  int64     //dbms-bytes-all                   : 317824922
	DbmsBytesLast5min             int64     `rac:"dbms-bytes-last-5min"` //dbms-bytes-last-5min             : 0
	DbProcInfo                    string    //db-proc-info                     :
	DbProcTook                    int32     //db-proc-took                     : 0
	DbProcTookAt                  time.Time //db-proc-took-at                  :
	DurationAll                   int64     //duration-all                     : 66184
	DurationAllDbms               int64     //duration-all-dbms                : 43242
	DurationCurrent               int64     //duration-current                 : 0
	DurationCurrentDbms           int64     //duration-current-dbms            : 0
	DurationLast5Min              int64     `rac:"duration-last-5min"`      //duration-last-5min               : 0
	DurationLast5MinDbms          int64     `rac:"duration-last-5min-dbms"` //duration-last-5min-dbms          : 0
	MemoryCurrent                 int64     //memory-current                   : 0
	MemoryLast5min                int64     //memory-last-5min                 : 416379
	MemoryTotal                   int64     //memory-total                     : 23178863
	ReadCurrent                   int64     //read-current                     : 0
	ReadLast5min                  int64     //read-last-5min                   : 0
	ReadTotal                     int64     //read-total                       : 156162
	WriteCurrent                  int64     //write-current                    : 0
	WriteLast5min                 int64     ///write-last-5min                  : 0
	WriteTotal                    int64     //write-total                      : 1071457
	DurationCurrentService        int64     //duration-current-service         : 0
	DurationLast5minService       int64     //duration-last-5min-service       : 30
	DurationAllService            int64     //duration-all-service             : 515
	CurrentServiceName            string    //current-service-name             :
	CpuTimeCurrent                int64     //cpu-time-current                 : 0
	CpuTimeLast5min               int64     //cpu-time-last-5min               : 280
	CpuTimeTotal                  int64     //cpu-time-total                   : 6832
	DataSeparation                string    //data-separation                  : ''
	ClientIp                      string    //client-ip                        :

}

type ConnectionInfo struct {
	UUID        string    `rac:"connection"` //connection     : cd16cde9-6372-4664-ac61-b0ae5cb24478
	ID          int32     `rac:"conn-id"`    //conn-id        : 8714
	Host        string    //host           : srv-uk-term-09
	Process     string    //process        : 94232f94-be78-4acd-a11e-09911bd4f4ed
	Infobase    string    //infobase       : efa3672f-947a-4d84-bd58-b21997b83561
	Application string    //application    : "1CV8"
	ConnectedAt time.Time //connected-at   : 2020-10-01T07:29:57
	SessionID   int64     `rac:"session-number"` //session-number : 148542
	BlockedByLs int       //blocked-by-ls  : 0

}

func (i ConnectionInfo) Sig() (string, string, string) {
	return i.UUID, i.Process, i.Infobase
}

type ProcessInfo struct {
	UUID                string    `rac:"process"` // process              : 3ea9968d-159c-4b5f-9bdc-22b8ead96f74
	Host                string    //host                 : Sport1
	Port                string    //port                 : 1564
	Pid                 int       //pid                  : 5428
	Enable              bool      `rac:"is-enable"` //is-enable            : yes
	Running             bool      //running              : yes
	StartedAt           time.Time //started-at           : 2018-03-29T11:16:02
	Use                 string    //use                  : used
	AvailablePerfomance int       //available-perfomance : 100
	Capacity            int32     //capacity             : 1000
	Connections         int32     //connections          : 7
	MemorySize          int64     //memory-size          : 1518604
	MemoryExcessTime    int64     //memory-excess-time   : 0
	SelectionSize       int64     //selection-size       : 61341
	AvgBackCallTime     float64   //avg-back-call-time   : 0.000
	AvgCallTime         float64   //avg-call-time        : 0.483
	AvgDbCallTime       float64   //avg-db-call-time     : 0.124
	AvgLockCallTime     float64   //avg-lock-call-time   : 0.000
	AvgServerCallTime   float64   //avg-server-call-time : -0.265
	AvgThreads          float64   //avg-threads          : 0.281
	Reverse             bool      //reserve              : no
}

type ClusterInfo struct {
	UUID                       string `rac:"cluster"` // UUID cluster                    : 6d6958e1-a96c-4999-a995-698a0298161e
	Host                       string // Host                          : Sport2
	Port                       int    // Port                          : 1541
	Name                       string // Name                          : "Новый кластер"
	ExpirationTimeout          int    // ExpirationTimeout expiration-timeout            : 0
	LifetimeLimit              int    // LifetimeLimit lifetime-limit                : 0
	MaxMemorySize              int    // MaxMemorySize max-memory-size               : 0
	MaxMemoryTimeLimit         int    // MaxMemoryTimeLimit max-memory-time-limit         : 0
	SecurityLevel              int    // SecurityLevel security-level                : 0
	SessionFaultToleranceLevel int    // SessionFaultToleranceLevel session-fault-tolerance-level : 0
	LoadBalancingMode          string // LoadBalancingMode load-balancing-mode           : performance
	ErrorsCountThreshold       int    // ErrorsCountThreshold errors-count-threshold        : 0
	KillProblemProcesses       int    // KillProblemProcesses kill-problem-processes        : 0
	KillByMemoryWithDump       int    // KillByMemoryWithDump kill-by-memory-with-dump      : 0
}

func (i ClusterInfo) ClusterSig() (string, string, string) {
	return i.UUID, "", ""
}

type ServiceInfo struct {
	Name        string //name      : EventLogService
	MainOnly    bool   //main-only : 0
	Manager     string //manager   : ad2754ad-9415-4689-9559-74dc36b11592
	Description string `rac:"descr"` //descr     : "Сервис журналов регистрации"
}

type InfobaseInfo struct {
	UUID                                   string    `rac:"infobase"` //infobase : efa3672f-947a-4d84-bd58-b21997b83561
	Name                                   string    //name     : УППБоеваяБаза
	Description                            string    `rac:"descr"` //descr    : "УППБоеваяБаза"
	Dbms                                   string    //dbms                                       : MSSQLServer
	DbServer                               string    //db-server                                  : sql
	DbName                                 string    //db-name                                    : base
	DbUser                                 string    //db-user                                    : user
	SecurityLevel                          int       //security-level                             : 0
	LicenseDistribution                    string    //license-distribution                       : allow
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

}

type ServerInfo struct {
	UUID                                 string `rac:"server"` //server                                    : 82b8f05a-898e-48ec-9a5b-461bdf66b7d0
	AgentHost                            string //agent-host                                : app
	AgentPort                            int    //agent-port                                : 1540
	PortRange                            string //port-range                                : 1560:1591
	Name                                 string //name                                      : "Центральный сервер"
	Using                                string //using                                     : main
	DedicateManagers                     string //dedicate-managers                         : none
	InfobasesLimit                       int32  //infobases-limit                           : 8
	MemoryLimit                          int64  //memory-limit                              : 0
	ConnectionsLimit                     int32  //connections-limit                         : 128
	SafeWorkingProcessesMemoryLimit      int32  //safe-working-processes-memory-limit       : 0
	SafeCallMemoryLimit                  int32  //safe-call-memory-limit                    : 0
	ClusterPort                          int    //cluster-port                              : 1541
	CriticalTotalMemory                  int64  //critical-total-memory                     : 0
	TemporaryAllowedTotalMemory          int64  //temporary-allowed-total-memory            : 0
	TemporaryAllowedTotalMemoryTimeLimit int64  //temporary-allowed-total-memory-time-limit : 300

}

type ManagerInfo struct {
	UUID        string `rac:"manager"` //manager : 0e588a25-8354-4344-b935-53442312aa30
	PID         int    //pid     : 3388
	Using       string //using   : normal
	Host        string //host    : Sport1
	Port        int    //port    : 1541
	Description string `rac:"descr"` //descr   : "Главный менеджер кластера"

}

type LicenseInfo struct {
	Process           string //process            : 94232f94-be78-4acd-a11e-09911bd4f4ed
	Session           string //session            : e45c1c2b-b3ac-4fea-9f0c-0583ad65d117
	UserName          string //user-name          : User
	Host              string //host               : host
	AppId             string //app-id             : 1CV8
	FullName          string //full-name          :
	Series            string //series             : "ORG8A"
	IssuedByServer    bool   //issued-by-server   : yes
	LicenseType       string //license-type       : HASP
	Net               bool   //net                : yes
	MaxUsersAll       int32  //max-users-all      : 300
	MaxUsersCur       int32  //max-users-cur      : 300
	RmngrAddress      string //rmngr-address      : "app"
	RmngrPort         int    //rmngr-port         : 1541
	RmngrPid          int32  //rmngr-pid          : 2300
	ShortPresentation string //short-presentation : "Сервер, ORG8A Сет 300"
	FullPresentation  string //full-presentation  : "Сервер, 2300, app, 1541, ORG8A Сетевой 300"
}
