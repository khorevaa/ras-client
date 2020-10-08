package protocol

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	GET_AGENT_ADMINS_REQUEST EndpointMessageType = iota
	GET_AGENT_ADMINS_RESPONSE
	GET_CLUSTER_ADMINS_REQUEST
	GET_CLUSTER_ADMINS_RESPONSE
	REG_AGENT_ADMIN_REQUEST
	REG_CLUSTER_ADMIN_REQUEST
	UNREG_AGENT_ADMIN_REQUEST
	UNREG_CLUSTER_ADMIN_REQUEST
	AUTHENTICATE_AGENT_REQUEST
	AUTHENTICATE_REQUEST
	ADD_AUTHENTICATION_REQUEST
	GET_CLUSTERS_REQUEST
	GET_CLUSTERS_RESPONSE
	GET_CLUSTER_INFO_REQUEST
	GET_CLUSTER_INFO_RESPONSE
	REG_CLUSTER_REQUEST
	REG_CLUSTER_RESPONSE
	UNREG_CLUSTER_REQUEST
	GET_CLUSTER_MANAGERS_REQUEST
	GET_CLUSTER_MANAGERS_RESPONSE
	GET_CLUSTER_MANAGER_INFO_REQUEST
	GET_CLUSTER_MANAGER_INFO_RESPONSE
	GET_WORKING_SERVERS_REQUEST
	GET_WORKING_SERVERS_RESPONSE
	GET_WORKING_SERVER_INFO_REQUEST
	GET_WORKING_SERVER_INFO_RESPONSE
	REG_WORKING_SERVER_REQUEST
	REG_WORKING_SERVER_RESPONSE
	UNREG_WORKING_SERVER_REQUEST
	GET_WORKING_PROCESSES_REQUEST
	GET_WORKING_PROCESSES_RESPONSE
	GET_WORKING_PROCESS_INFO_REQUEST
	GET_WORKING_PROCESS_INFO_RESPONSE
	GET_SERVER_WORKING_PROCESSES_REQUEST
	GET_SERVER_WORKING_PROCESSES_RESPONSE
	GET_CLUSTER_SERVICES_REQUEST
	GET_CLUSTER_SERVICES_RESPONSE
	CREATE_INFOBASE_REQUEST
	CREATE_INFOBASE_RESPONSE
	UPDATE_INFOBASE_SHORT_REQUEST
	UPDATE_INFOBASE_REQUEST
	DROP_INFOBASE_REQUEST
	GET_INFOBASES_SHORT_REQUEST
	GET_INFOBASES_SHORT_RESPONSE
	GET_INFOBASES_REQUEST
	GET_INFOBASES_RESPONSE
	GET_INFOBASE_SHORT_INFO_REQUEST
	GET_INFOBASE_SHORT_INFO_RESPONSE
	GET_INFOBASE_INFO_REQUEST
	GET_INFOBASE_INFO_RESPONSE
	GET_CONNECTIONS_SHORT_REQUEST
	GET_CONNECTIONS_SHORT_RESPONSE
	GET_INFOBASE_CONNECTIONS_SHORT_REQUEST
	GET_INFOBASE_CONNECTIONS_SHORT_RESPONSE
	GET_CONNECTION_INFO_SHORT_REQUEST
	GET_CONNECTION_INFO_SHORT_RESPONSE
	GET_INFOBASE_CONNECTIONS_REQUEST
	GET_INFOBASE_CONNECTIONS_RESPONSE
	DISCONNECT_REQUEST
	GET_SESSIONS_REQUEST
	GET_SESSIONS_RESPONSE
	GET_INFOBASE_SESSIONS_REQUEST
	GET_INFOBASE_SESSIONS_RESPONSE
	GET_SESSION_INFO_REQUEST
	GET_SESSION_INFO_RESPONSE
	TERMINATE_SESSION_REQUEST
	GET_LOCKS_REQUEST
	GET_LOCKS_RESPONSE
	GET_INFOBASE_LOCKS_REQUEST
	GET_INFOBASE_LOCKS_RESPONSE
	GET_CONNECTION_LOCKS_REQUEST
	GET_CONNECTION_LOCKS_RESPONSE
	GET_SESSION_LOCKS_REQUEST
	GET_SESSION_LOCKS_RESPONSE
	APPLY_ASSIGNMENT_RULES_REQUEST
	REG_ASSIGNMENT_RULE_REQUEST
	REG_ASSIGNMENT_RULE_RESPONSE
	UNREG_ASSIGNMENT_RULE_REQUEST
	GET_ASSIGNMENT_RULES_REQUEST
	GET_ASSIGNMENT_RULES_RESPONSE
	GET_ASSIGNMENT_RULE_INFO_REQUEST
	GET_ASSIGNMENT_RULE_INFO_RESPONSE
	GET_SECURITY_PROFILES_REQUEST
	GET_SECURITY_PROFILES_RESPONSE
	CREATE_SECURITY_PROFILE_REQUEST
	DROP_SECURITY_PROFILE_REQUEST
	GET_VIRTUAL_DIRECTORIES_REQUEST
	GET_VIRTUAL_DIRECTORIES_RESPONSE
	CREATE_VIRTUAL_DIRECTORY_REQUEST
	DROP_VIRTUAL_DIRECTORY_REQUEST
	GET_COM_CLASSES_REQUEST
	GET_COM_CLASSES_RESPONSE
	CREATE_COM_CLASS_REQUEST
	DROP_COM_CLASS_REQUEST
	GET_ALLOWED_ADDINS_REQUEST
	GET_ALLOWED_ADDINS_RESPONSE
	CREATE_ALLOWED_ADDIN_REQUEST
	DROP_ALLOWED_ADDIN_REQUEST
	GET_EXTERNAL_MODULES_REQUEST
	GET_EXTERNAL_MODULES_RESPONSE
	CREATE_EXTERNAL_MODULE_REQUEST
	DROP_EXTERNAL_MODULE_REQUEST
	GET_ALLOWED_APPLICATIONS_REQUEST
	GET_ALLOWED_APPLICATIONS_RESPONSE
	CREATE_ALLOWED_APPLICATION_REQUEST
	DROP_ALLOWED_APPLICATION_REQUEST
	GET_INTERNET_RESOURCES_REQUEST
	GET_INTERNET_RESOURCES_RESPONSE
	CREATE_INTERNET_RESOURCE_REQUEST
	DROP_INTERNET_RESOURCE_REQUEST
	INTERRUPT_SESSION_CURRENT_SERVER_CALL_REQUEST
	GET_RESOURCE_COUNTERS_REQUEST
	GET_RESOURCE_COUNTERS_RESPONSE
	GET_RESOURCE_COUNTER_INFO_REQUEST
	GET_RESOURCE_COUNTER_INFO_RESPONSE
	REG_RESOURCE_COUNTER_REQUEST
	UNREG_RESOURCE_COUNTER_REQUEST
	GET_RESOURCE_LIMITS_REQUEST
	GET_RESOURCE_LIMITS_RESPONSE
	GET_RESOURCE_LIMIT_INFO_REQUEST
	GET_RESOURCE_LIMIT_INFO_RESPONSE
	REG_RESOURCE_LIMIT_REQUEST
	UNREG_RESOURCE_LIMIT_REQUEST
	GET_COUNTER_VALUES_REQUEST
	GET_COUNTER_VALUES_RESPONSE
	CLEAR_COUNTER_VALUE_REQUEST
	GET_COUNTER_ACCUMULATED_VALUES_REQUEST
	GET_COUNTER_ACCUMULATED_VALUES_RESPONSE
	GET_AGENT_VERSION_REQUEST
	GET_AGENT_VERSION_RESPONSE
)

type EndpointMessageType int

func (t EndpointMessageType) Type() int {
	return int(t)
}

type GetClustersRequest struct{}

func (_ GetClustersRequest) Type() MessageType {
	return GET_CLUSTERS_REQUEST
}

func (_ GetClustersRequest) Format(_ *encoder) {

}

type GetClustersResponse struct {
	Clusters []*ClusterInfo
}

func (_ *GetClustersResponse) Type() MessageType {
	return GET_CLUSTERS_RESPONSE
}

func (r *GetClustersResponse) Parse(body []byte) error {

	decoder := NewDecoder(body)

	count := decoder.decodeSize()

	for i := 0; i < count; i++ {

		info := &ClusterInfo{}
		info.UUID = decoder.decodeUUID().String()
		_ = decoder.decodeInt() // expirationTimeout
		info.Host = decoder.decodeString()
		info.ExpirationTimeout = int(decoder.decodeInt())
		info.Port = int(decoder.decodeUnsignedShort())
		info.MaxMemorySize = int(decoder.decodeInt())
		info.MaxMemoryTimeLimit = int(decoder.decodeInt())
		info.Name = decoder.decodeString()
		info.SecurityLevel = int(decoder.decodeInt())
		info.SessionFaultToleranceLevel = int(decoder.decodeInt())
		info.LoadBalancingMode = int(decoder.decodeInt()) // Не понтяно что
		info.ErrorsCountThreshold = int(decoder.decodeInt())
		info.KillProblemProcesses = decoder.decodeBoolean()
		info.KillByMemoryWithDump = decoder.decodeBoolean()

		r.Clusters = append(r.Clusters, info)
	}

	return nil

}

type AuthenticateAgentRequest struct {
	user, password string
}

func (_ AuthenticateAgentRequest) Type() MessageType {
	return AUTHENTICATE_AGENT_REQUEST
}

func (r AuthenticateAgentRequest) Format(e *encoder) {

	e.encodeString(r.user)
	e.encodeString(r.password)

}

type GetClusterManagersRequest struct {
	ID uuid.UUID
}

func (_ GetClusterManagersRequest) Type() MessageType {
	return GET_CLUSTER_MANAGERS_REQUEST
}

func (r GetClusterManagersRequest) Format(e *encoder) {

	e.encodeUuid(r.ID)

}

type GetClusterManagersResponse struct {
	Managers []*ManagerInfo
}

func (_ *GetClusterManagersResponse) Type() MessageType {
	return GET_CLUSTER_MANAGERS_RESPONSE
}

func (r *GetClusterManagersResponse) Parse(body []byte) error {

	decoder := NewDecoder(body)

	count := decoder.decodeSize()

	for i := 0; i < count; i++ {

		info := &ManagerInfo{}
		info.UUID = decoder.decodeUUID().String()
		info.Description = decoder.decodeString()
		info.Host = decoder.decodeString()
		info.MainManager = int(decoder.decodeInt())    // expirationTimeout
		info.Port = int(decoder.decodeUnsignedShort()) // expirationTimeout
		info.PID = decoder.decodeString()

		r.Managers = append(r.Managers, info)
	}

	return nil

}

type ClusterAuthenticateRequest struct {
	ID             uuid.UUID
	user, password string
}

func (_ ClusterAuthenticateRequest) Type() MessageType {
	return AUTHENTICATE_REQUEST
}

func (r ClusterAuthenticateRequest) Format(e *encoder) {

	e.encodeUuid(r.ID)
	e.encodeString(r.user)
	e.encodeString(r.password)

}

type ManagerInfo struct {
	UUID        string `rac:"manager"` //manager : 0e588a25-8354-4344-b935-53442312aa30
	PID         string //pid     : 3388
	Using       string //using   : normal
	Host        string //host    : Sport1
	MainManager int
	Port        int    //port    : 1541
	Description string `rac:"descr"` //descr   : "Главный менеджер кластера"

}

type GetClusterServicesRequest struct {
	ID uuid.UUID
}

func (_ GetClusterServicesRequest) Type() MessageType {
	return GET_CLUSTER_SERVICES_REQUEST
}

func (r GetClusterServicesRequest) Format(e *encoder) {

	e.encodeUuid(r.ID)

}

type GetClusterServicesResponse struct {
	Services []*ServiceInfo
}

func (_ *GetClusterServicesResponse) Type() MessageType {
	return GET_CLUSTER_SERVICES_RESPONSE
}

func (r *GetClusterServicesResponse) Parse(body []byte) error {

	decoder := NewDecoder(body)

	count := decoder.decodeSize()

	for i := 0; i < count; i++ {

		info := &ServiceInfo{}
		info.Name = decoder.decodeString()
		info.Description = decoder.decodeString()
		info.MainOnly = decoder.decodeInt()

		idCount := decoder.decodeSize()

		for ii := 0; ii < idCount; ii++ {
			info.Manager = append(info.Manager, decoder.decodeUUID().String())
		}
		r.Services = append(r.Services, info)
	}

	return nil

}

type ServiceInfo struct {
	Name        string   //name      : EventLogService
	MainOnly    int32    //main-only : 0
	Manager     []string //manager   : ad2754ad-9415-4689-9559-74dc36b11592
	Description string   `rac:"descr"` //descr     : "Сервис журналов регистрации"
}

type GetInfobasesRequest struct {
	ID uuid.UUID
}

func (_ GetInfobasesRequest) Type() MessageType {
	return GET_INFOBASES_REQUEST
}

func (r GetInfobasesRequest) Format(e *encoder) {

	e.encodeUuid(r.ID)

}

type GetInfobasesResponse struct {
	Infobases []*InfobaseInfo
}

func (_ *GetInfobasesResponse) Type() MessageType {
	return GET_INFOBASES_RESPONSE
}

func (r *GetInfobasesResponse) Parse(body []byte) error {

	decoder := NewDecoder(body)

	count := decoder.decodeSize()

	for i := 0; i < count; i++ {

		info := &InfobaseInfo{}
		info.UUID = decoder.decodeUUID().String()
		info.DateOffset = decoder.decodeInt()
		info.Dbms = decoder.decodeString()
		info.DbName = decoder.decodeString()
		info.DbPwd = decoder.decodeString()
		info.DbServer = decoder.decodeString()
		info.DbUser = decoder.decodeString()
		info.DeniedFrom = dateFromTicks(decoder.decodeLong())
		info.DeniedMessage = decoder.decodeString()
		info.DeniedParameter = decoder.decodeString()
		info.DeniedTo = dateFromTicks(decoder.decodeLong())
		info.Description = decoder.decodeString()
		info.Locale = decoder.decodeString()
		info.Name = decoder.decodeString()
		info.PermissionCode = decoder.decodeString()
		info.ScheduledJobsDeny = decoder.decodeBoolean()
		info.SecurityLevel = int(decoder.decodeInt())
		info.SessionsDeny = decoder.decodeBoolean()
		info.LicenseDistribution = decoder.decodeInt()
		info.ExternalSessionManagerConnectionDtring = decoder.decodeString()
		info.ExternalSessionManagerRequired = decoder.decodeBoolean()
		info.SecurityProfileName = decoder.decodeString()
		info.SafeModeSecurityProfileName = decoder.decodeString()
		info.ReserveWorkingProcesses = decoder.decodeBoolean()

		r.Infobases = append(r.Infobases, info)
	}

	return nil

}

type InfobaseInfo struct {
	UUID                                   string    `rac:"infobase"` //infobase : efa3672f-947a-4d84-bd58-b21997b83561
	Name                                   string    //name     : УППБоеваяБаза
	Description                            string    `rac:"descr"` //descr    : "УППБоеваяБаза"
	Dbms                                   string    //dbms                                       : MSSQLServer
	DbServer                               string    //db-server                                  : sql
	DbName                                 string    //db-name                                    : base
	DbUser                                 string    //db-user                                    : user
	DbPwd                                  string    `rac:"="` //--db-pwd=<pwd>  пароль администратора базы данных
	SecurityLevel                          int       //security-level                             : 0
	LicenseDistribution                    int32     //license-distribution                       : allow
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
	DateOffset                             int32
	Locale                                 string
}

func dateFromTicks(ticks int64) time.Time {
	if ticks > 0 {

		timeT := (ticks - 621355968000000) / 10

		t := time.Unix(timeT, 0)
		return t

	}
	return time.Time{}
}

type GetInfobasesShortRequest struct {
	ID uuid.UUID
}

func (_ GetInfobasesShortRequest) Type() MessageType {
	return GET_INFOBASES_SHORT_REQUEST
}

func (r GetInfobasesShortRequest) Format(e *encoder) {

	e.encodeUuid(r.ID)

}

type GetInfobasesShortResponse struct {
	Infobases []*InfobaseInfo
}

func (_ *GetInfobasesShortResponse) Type() MessageType {
	return GET_INFOBASES_SHORT_RESPONSE
}

func (r *GetInfobasesShortResponse) Parse(body []byte) error {

	decoder := NewDecoder(body)

	count := decoder.decodeSize()

	for i := 0; i < count; i++ {

		info := &InfobaseInfo{}
		info.UUID = decoder.decodeUUID().String()
		info.Description = decoder.decodeString()
		info.Name = decoder.decodeString()

		r.Infobases = append(r.Infobases, info)
	}

	return nil

}

type GetConnectionsShortRequest struct {
	ID uuid.UUID
}

func (_ GetConnectionsShortRequest) Type() MessageType {
	return GET_CONNECTIONS_SHORT_REQUEST
}

func (r GetConnectionsShortRequest) Format(e *encoder) {

	e.encodeUuid(r.ID)

}

type GetConnectionsShortResponse struct {
	Connections []*ConnectionInfo
}

func (_ *GetConnectionsShortResponse) Type() MessageType {
	return GET_CONNECTIONS_SHORT_RESPONSE
}

func (r *GetConnectionsShortResponse) Parse(body []byte) error {

	decoder := NewDecoder(body)

	count := decoder.decodeSize()

	for i := 0; i < count; i++ {

		info := &ConnectionInfo{}
		info.UUID = decoder.decodeUUID().String()
		info.Application = decoder.decodeString()
		info.BlockedByLs = decoder.decodeInt()
		info.ID = int32(decoder.decodeInt())
		info.ConnectedAt = time.Unix(decoder.decodeLong(), 0)
		info.SessionID = decoder.decodeInt()
		info.Host = decoder.decodeString()
		info.Infobase = decoder.decodeUUID().String()

		info.Process = decoder.decodeUUID().String()
		//_ = decoder.decodeByte()

		r.Connections = append(r.Connections, info)
	}

	return nil

}

type ConnectionInfo struct {
	UUID        string    `rac:"connection"` //connection     : cd16cde9-6372-4664-ac61-b0ae5cb24478
	ID          int32     `rac:"conn-id"`    //conn-id        : 8714
	Host        string    //host           : srv-uk-term-09
	Process     string    //process        : 94232f94-be78-4acd-a11e-09911bd4f4ed
	Infobase    string    //infobase       : efa3672f-947a-4d84-bd58-b21997b83561
	Application string    //application    : "1CV8"
	ConnectedAt time.Time //connected-at   : 2020-10-01T07:29:57
	SessionID   int32     `rac:"session-number"` //session-number : 148542
	BlockedByLs int32     //blocked-by-ls  : 0

}
