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

type MessageType interface {
	Type() int
}

type ConnectionMessageType int

const (
	NEGOTIATE         ConnectionMessageType = 0
	CONNECT           ConnectionMessageType = 1
	CONNECT_ACK       ConnectionMessageType = 2
	START_TLS         ConnectionMessageType = 3
	DISCONNECT        ConnectionMessageType = 4
	SASL_NEGOTIATE    ConnectionMessageType = 5
	SASL_AUTH         ConnectionMessageType = 6
	SASL_CHALLENGE    ConnectionMessageType = 7
	SASL_SUCCESS      ConnectionMessageType = 8
	SASL_FAILURE      ConnectionMessageType = 9
	SASL_ABORT        ConnectionMessageType = 10
	ENDPOINT_OPEN     ConnectionMessageType = 11
	ENDPOINT_OPEN_ACK ConnectionMessageType = 12
	ENDPOINT_CLOSE    ConnectionMessageType = 13
	ENDPOINT_MESSAGE  ConnectionMessageType = 14
	ENDPOINT_FAILURE  ConnectionMessageType = 15
	KEEP_ALIVE        ConnectionMessageType = 16

	NULL_TYPE ConnectionMessageType = 127
)

func (m ConnectionMessageType) String() string {

	switch m {

	case CONNECT:
		return "CONNECT"
	case ENDPOINT_FAILURE:
		return "ENDPOINT_FAILURE"
	case ENDPOINT_MESSAGE:
		return "ENDPOINT_MESSAGE"
	case CONNECT_ACK:
		return "CONNECT_ACK"
	case ENDPOINT_CLOSE:
		return "ENDPOINT_CLOSE"
	default:
		return "неизвестный тим ответа"
	}

}

func (m ConnectionMessageType) Type() int {
	return int(m)
}

//public final class TypeCodes

type ServiceWireType int

const (
	BOOLEAN       ServiceWireType = 1
	BYTE                          = 2
	SHORT                         = 3
	INT                           = 4
	LONG                          = 5
	FLOAT                         = 6
	DOUBLE                        = 7
	SIZE                          = 8
	NULLABLE_SIZE                 = 9
	STRING                        = 10
	UUID                          = 11
	TYPE                          = 12
	ENDPOINT_ID                   = 13
)

func (t ServiceWireType) raw() byte {
	return byte(t)
}
func (t ServiceWireType) Type() int {
	return int(t)
}
