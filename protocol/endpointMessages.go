package protocol

import (
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
	"io"
	"time"
)

const (
	NULL_ENDPOINT_RESPONSE   EndpointMessageType = -1
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

type GetClustersRequest struct {
	response *GetClustersResponse
}

func (_ *GetClustersRequest) Format(encoder codec.Encoder, version int, w io.Writer) {}

func (_ *GetClustersRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (r *GetClustersRequest) ResponseMessage() types.EndpointResponseMessage {

	if r.response == nil {
		r.response = &GetClustersResponse{}
	}

	return r.response
}

func (r *GetClustersRequest) Response() *GetClustersResponse {
	return r.response
}

func (_ GetClustersRequest) Type() types.Typed {
	return GET_CLUSTERS_REQUEST
}

type GetClustersResponse struct {
	Clusters []*ClusterInfo
}

func (res *GetClustersResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &ClusterInfo{}
		decoder.UuidPtr(&info.UUID, r)
		decoder.IntPtr(&info.ExpirationTimeout, r) // expirationTimeout
		decoder.StringPtr(&info.Host, r)
		decoder.IntPtr(&info.LifeTimeLimit, r)
		decoder.ShortPtr(&info.Port, r)
		decoder.IntPtr(&info.MaxMemorySize, r)
		decoder.IntPtr(&info.MaxMemoryTimeLimit, r)
		decoder.StringPtr(&info.Name, r)
		decoder.IntPtr(&info.SecurityLevel, r)
		decoder.IntPtr(&info.SessionFaultToleranceLevel, r)
		decoder.IntPtr(&info.LoadBalancingMode, r)
		decoder.IntPtr(&info.ErrorsCountThreshold, r)
		decoder.BoolPtr(&info.KillProblemProcesses, r)

		if version > 0 {
			decoder.BoolPtr(&info.KillByMemoryWithDump, r)
		}

		res.Clusters = append(res.Clusters, info)
	}

}

func (_ *GetClustersResponse) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ *GetClustersResponse) Type() types.Typed {
	return GET_CLUSTERS_RESPONSE
}

type AuthenticateAgentRequest struct {
	user, password string
}

func (_ AuthenticateAgentRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ AuthenticateAgentRequest) ResponseMessage() types.EndpointResponseMessage {
	return nullEndpointResponse()
}

func (_ AuthenticateAgentRequest) Type() types.Typed {
	return AUTHENTICATE_AGENT_REQUEST
}

func (r AuthenticateAgentRequest) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.String(r.user, w)
	encoder.String(r.password, w)

}

type GetClusterManagersRequest struct {
	ID uuid.UUID

	response *GetClusterManagersResponse
}

func (_ *GetClusterManagersRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (r *GetClusterManagersRequest) ResponseMessage() types.EndpointResponseMessage {

	if r.response == nil {
		r.response = &GetClusterManagersResponse{}
	}

	return r.response
}

func (_ *GetClusterManagersRequest) Type() types.Typed {
	return GET_CLUSTER_MANAGERS_REQUEST
}

func (r *GetClusterManagersRequest) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.Uuid(r.ID, w)

}

func (r *GetClusterManagersRequest) Response() *GetClusterManagersResponse {
	return r.response
}

type GetClusterManagersResponse struct {
	Managers []*ManagerInfo
}

func (res *GetClusterManagersResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &ManagerInfo{}
		decoder.UuidPtr(&info.UUID, r)
		decoder.StringPtr(&info.Description, r)
		decoder.StringPtr(&info.Host, r)
		decoder.IntPtr(&info.MainManager, r)
		decoder.ShortPtr(&info.Port, r) // expirationTimeout
		decoder.StringPtr(&info.PID, r)

		res.Managers = append(res.Managers, info)
	}
}

func (_ *GetClusterManagersResponse) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ *GetClusterManagersResponse) Type() types.Typed {
	return GET_CLUSTER_MANAGERS_RESPONSE
}

type ClusterAuthenticateRequest struct {
	ID             uuid.UUID
	user, password string
}

func (r ClusterAuthenticateRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
	encoder.String(r.user, w)
	encoder.String(r.password, w)
}

func (_ ClusterAuthenticateRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ ClusterAuthenticateRequest) ResponseMessage() types.EndpointResponseMessage {
	return nullEndpointResponse()
}

func (_ ClusterAuthenticateRequest) Type() types.Typed {
	return AUTHENTICATE_REQUEST
}

type ManagerInfo struct {
	UUID        uuid.UUID `rac:"manager"` //manager : 0e588a25-8354-4344-b935-53442312aa30
	PID         string    //pid     : 3388
	Using       string    //using   : normal
	Host        string    //host    : Sport1
	MainManager int
	Port        int16  //port    : 1541
	Description string `rac:"descr"` //descr   : "Главный менеджер кластера"

}

type GetClusterServicesRequest struct {
	ID       uuid.UUID
	response *GetClusterServicesResponse
}

func (r *GetClusterServicesRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
}

func (_ *GetClusterServicesRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (r *GetClusterServicesRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetClusterServicesResponse{}
	}

	return r.response
}

func (r *GetClusterServicesRequest) Response() *GetClusterServicesResponse {

	return r.response
}

func (_ *GetClusterServicesRequest) Type() types.Typed {
	return GET_CLUSTER_SERVICES_REQUEST
}

type GetClusterServicesResponse struct {
	Services []*ServiceInfo
}

func (res *GetClusterServicesResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &ServiceInfo{}
		decoder.StringPtr(&info.Name, r)
		decoder.StringPtr(&info.Description, r)
		decoder.IntPtr(&info.MainOnly, r)

		idCount := decoder.Size(r)

		for ii := 0; ii < idCount; ii++ {
			info.Manager = append(info.Manager, decoder.Uuid(r).String())
		}
		res.Services = append(res.Services, info)
	}

}

func (_ *GetClusterServicesResponse) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ *GetClusterServicesResponse) Type() types.Typed {
	return GET_CLUSTER_SERVICES_RESPONSE
}

type ServiceInfo struct {
	Name        string   //name      : EventLogService
	MainOnly    int      //main-only : 0
	Manager     []string //manager   : ad2754ad-9415-4689-9559-74dc36b11592
	Description string   `rac:"descr"` //descr     : "Сервис журналов регистрации"
}

type GetInfobasesRequest struct {
	ID uuid.UUID
}

func (_ GetInfobasesRequest) Type() types.Typed {
	return GET_INFOBASES_REQUEST
}

func (r GetInfobasesRequest) Format(e *encoder) {

	e.encodeUuid(r.ID)

}

type GetInfobasesResponse struct {
	Infobases []*InfobaseInfo
}

func (_ *GetInfobasesResponse) Type() types.Typed {
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
	SafeModeSecurityProfileName            string    //safe-Mode-security-profile-name            :
	ReserveWorkingProcesses                bool      //reserve-working-processes                  : no
	DateOffset                             int
	Locale                                 string
}

func (i *InfobaseInfo) Parse(decoder codec.Decoder, version int, r io.Reader) {

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

type GetConnectionsShortRequest struct {
	ID       uuid.UUID
	response *GetConnectionsShortResponse
}

func (_ *GetConnectionsShortRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (r *GetConnectionsShortRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetConnectionsShortResponse{}
	}

	return r.response
}

func (_ *GetConnectionsShortRequest) Type() types.Typed {
	return GET_CONNECTIONS_SHORT_REQUEST
}

func (r *GetConnectionsShortRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
}

func (r *GetConnectionsShortRequest) Response() *GetConnectionsShortResponse {
	return r.response
}

type GetConnectionsShortResponse struct {
	Connections []*ConnectionInfo
}

func (_ *GetConnectionsShortResponse) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ *GetConnectionsShortResponse) Type() types.Typed {
	return GET_CONNECTIONS_SHORT_RESPONSE
}

func (res *GetConnectionsShortResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &ConnectionInfo{}
		info.Parse(decoder, version, r)

		res.Connections = append(res.Connections, info)
	}

}

type ConnectionInfo struct {
	UUID        uuid.UUID `rac:"connection"` //connection     : cd16cde9-6372-4664-ac61-b0ae5cb24478
	ID          int       `rac:"conn-id"`    //conn-id        : 8714
	Host        string    //host           : srv-uk-term-09
	Process     uuid.UUID //process        : 94232f94-be78-4acd-a11e-09911bd4f4ed
	Infobase    uuid.UUID //infobase       : efa3672f-947a-4d84-bd58-b21997b83561
	Application string    //application    : "1CV8"
	ConnectedAt time.Time //connected-at   : 2020-10-01T07:29:57
	SessionID   int       `rac:"session-number"` //session-number : 148542
	BlockedByLs int       //blocked-by-ls  : 0
}

func (i *ConnectionInfo) Parse(decoder codec.Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Application, r)
	decoder.IntPtr(&i.BlockedByLs, r)
	decoder.TimePtr(&i.ConnectedAt, r)
	decoder.IntPtr(&i.ID, r)
	decoder.StringPtr(&i.Host, r)
	decoder.UuidPtr(&i.Infobase, r)
	decoder.UuidPtr(&i.Process, r)
	decoder.IntPtr(&i.SessionID, r)

}

type CreateInfobaseRequest struct {
	ID       uuid.UUID
	Infobase *serialize.InfobaseInfo
	response *CreateInfobaseResponse
	Mode     int // 1 - создавать базу на сервере, 0 - не создавать

}

func (_ *CreateInfobaseRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (r *CreateInfobaseRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &CreateInfobaseResponse{}
	}

	return r.response
}

func (_ *CreateInfobaseRequest) Type() types.Typed {
	return CREATE_INFOBASE_REQUEST
}

func (r *CreateInfobaseRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)

	r.Infobase.Format(encoder, version, w)

	encoder.Int(r.Mode, w)
}

func (r *CreateInfobaseRequest) Response() *CreateInfobaseResponse {
	return r.response
}

type CreateInfobaseResponse struct {
	ID uuid.UUID
}

func (_ *CreateInfobaseResponse) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ *CreateInfobaseResponse) Type() types.Typed {
	return CREATE_INFOBASE_RESPONSE
}

func (res *CreateInfobaseResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&res.ID, r)

}

type GetInfobaseInfoRequest struct {
	ID         uuid.UUID
	InfobaseId uuid.UUID
	response   *GetInfobaseInfoResponse
}

func (_ *GetInfobaseInfoRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (r *GetInfobaseInfoRequest) ResponseMessage() types.EndpointResponseMessage {
	if r.response == nil {
		r.response = &GetInfobaseInfoResponse{}
	}

	return r.response
}

func (_ *GetInfobaseInfoRequest) Type() types.Typed {
	return GET_INFOBASE_INFO_REQUEST
}

func (r *GetInfobaseInfoRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
	encoder.Uuid(r.InfobaseId, w)
}

func (r *GetInfobaseInfoRequest) Response() *GetInfobaseInfoResponse {
	return r.response
}

type GetInfobaseInfoResponse struct {
	infobase *serialize.InfobaseInfo
}

func (_ *GetInfobaseInfoResponse) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (_ *GetInfobaseInfoResponse) Type() types.Typed {
	return GET_INFOBASE_INFO_RESPONSE
}

func (res *GetInfobaseInfoResponse) Parse(decoder codec.Decoder, version int, r io.Reader) {

	info := &serialize.InfobaseInfo{}
	info.Parse(decoder, version, r)
	res.infobase = info

}

type DropInfobaseRequest struct {
	ID         uuid.UUID
	InfobaseId uuid.UUID
	Mode       int
}

func (_ *DropInfobaseRequest) Kind() EndpointMessageKind {
	return MESSAGE_KIND
}

func (r *DropInfobaseRequest) ResponseMessage() types.EndpointResponseMessage {

	return nullEndpointResponse()
}

func (_ *DropInfobaseRequest) Type() types.Typed {
	return DROP_INFOBASE_REQUEST
}

func (r *DropInfobaseRequest) Format(encoder codec.Encoder, version int, w io.Writer) {
	encoder.Uuid(r.ID, w)
	encoder.Uuid(r.InfobaseId, w)
	encoder.Int(r.Mode, w)
}

//GET_INFOBASE_INFO_REQUEST
