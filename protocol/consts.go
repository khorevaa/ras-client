package protocol

type ConnectionMessageType int

const (
	NEGOTIATE         ConnectionMessageType = 0
	CONNECT           ConnectionMessageType = 1
	CONNECT_ACK       ConnectionMessageType = 2
	START_TLS         ConnectionMessageType = 3 // Deprecated: Нереализовано в апи
	DISCONNECT        ConnectionMessageType = 4
	SASL_NEGOTIATE    ConnectionMessageType = 5  // Deprecated: Нереализовано в апи
	SASL_AUTH         ConnectionMessageType = 6  // Deprecated: Нереализовано в апи
	SASL_CHALLENGE    ConnectionMessageType = 7  // Deprecated: Нереализовано в апи
	SASL_SUCCESS      ConnectionMessageType = 8  // Deprecated: Нереализовано в апи
	SASL_FAILURE      ConnectionMessageType = 9  // Deprecated: Нереализовано в апи
	SASL_ABORT        ConnectionMessageType = 10 // Deprecated: Нереализовано в апи
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
