package protocol

import (
	"github.com/k0kubun/pp"
	uuid "github.com/satori/go.uuid"
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
	Managers []*ClusterInfo
}

func (_ *GetClusterManagersResponse) Type() MessageType {
	return GET_CLUSTER_MANAGERS_RESPONSE
}

func (r *GetClusterManagersResponse) Parse(body []byte) error {

	pp.Println(body)

	//decoder := NewDecoder(body)
	////mType := decoder.decodeType()
	////pp.Println("message type: %s", mType)
	////_ = decoder.decodeByte()
	////_ = decoder.decodeByte()
	//
	////t.Logf("endpoint: %v", EndpointId)
	////t.Logf("format: %v", format)
	////t.Logf("compression %v", format&0x1 != 0x0)
	//
	//info := &ClusterInfo{}
	//info.UUID = decoder.decodeUUID().String()
	//_ = decoder.decodeInt() // expirationTimeout
	//info.Host = decoder.decodeString()
	//info.ExpirationTimeout = int(decoder.decodeInt())
	//info.Port = int(decoder.decodeUnsignedShort())
	//info.MaxMemorySize = int(decoder.decodeInt())
	//info.MaxMemoryTimeLimit = int(decoder.decodeInt())
	//info.Name = decoder.decodeString()
	//info.SecurityLevel = int(decoder.decodeInt())
	//info.SessionFaultToleranceLevel = int(decoder.decodeInt())
	//info.LoadBalancingMode = int(decoder.decodeInt()) // Не понтяно что
	//info.ErrorsCountThreshold = int(decoder.decodeInt())
	//info.KillProblemProcesses = decoder.decodeBoolean()
	//info.KillByMemoryWithDump = decoder.decodeBoolean()
	//
	//r.Clusters = append(r.Clusters, info)

	return nil

}
