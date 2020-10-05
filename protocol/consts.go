package protocol

const (
	NAME                         = "swp"
	VERSION_1_0                  = 256
	MAGIC                        = 475223888
	DEFAULT_CONNECT_TIMEOUT      = 2000
	MESSAGE_ENCODED              = 1
	PROTOCOL_PARAMETERS_PREFIX   = "nipple.swp.protocol."
	CONNECTION_PARAMETERS_PREFIX = "nipple.swp.connection."
	SASL_METHODS                 = "sasl.methods"
	SECURE_REQUIRED              = "secure.required"
	KEEP_ALIVE_TIMEOUT           = "keep_alive.timeout"
	CONNECT_TIMEOUT              = "connect.timeout"
	ENDPO_TIMEOUT                = "endpo.timeout"
	ENDPO_ENCODING               = "endpo.encoding"
	SERVER_NAME                  = "server.name"
	SERVER_PORT                  = "server.port"
	SASL_USERNAME                = "sasl.username"
	SASL_PASSWORD                = "sasl.password"
)

func getConnectionParameterName(parameterName string) string {
	return "nipple.swp.connection." + parameterName
}

func getProtocolParameterName(parameterName string) string {
	return "nipple.swp.protocol." + parameterName
}

//public enum MessageType

type MessageType byte

const (
	NEGOTIATE         MessageType = 1
	CONNECT_ACK       MessageType = 2
	START_TLS         MessageType = 3
	DISCONNECT        MessageType = 4
	SASL_NEGOTIATE    MessageType = 5
	SASL_AUTH         MessageType = 6
	SASL_CHALLENGE    MessageType = 7
	SASL_SUCCESS      MessageType = 8
	SASL_FAILURE      MessageType = 9
	SASL_ABORT        MessageType = 10
	ENDPOINT_OPEN     MessageType = 11
	ENDPOINT_OPEN_ACK MessageType = 12
	ENDPOINT_CLOSE    MessageType = 13
	ENDPOINT_MESSAGE  MessageType = 14
	ENDPOINT_FAILURE  MessageType = 15
	KEEP_ALIVE        MessageType = 16
)

//public final class TypeCodes

const (
	BOOLEAN       = 1
	BYTE          = 2
	SHORT         = 3
	INT           = 4
	LONG          = 5
	FLOAT         = 6
	DOUBLE        = 7
	SIZE          = 8
	NULLABLE_SIZE = 9
	STRING        = 10
	UUID          = 11
	TYPE          = 12
	ENDPOINT_ID   = 13
)
