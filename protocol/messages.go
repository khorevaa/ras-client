package protocol

import (
	"github.com/k0kubun/pp"
	"io/ioutil"
)

type EndpointMessageKind int

const (
	EXCEPTION_KIND    EndpointMessageKind = -1
	VOID_MESSAGE_KIND EndpointMessageKind = iota
	MESSAGE_KIND
)

type ConnectMessageAck struct {
	data []byte
}

func (r *ConnectMessageAck) Type() MessageType {
	return CONNECT_ACK
}

func (r *ConnectMessageAck) Parse(body []byte) error {

	r.data = body

	return nil

}

type ConnectMessage struct {
	params map[string]interface{}
}

func (m *ConnectMessage) String() string {
	return ""
}

func (m *ConnectMessage) Type() MessageType {
	return CONNECT
}

func (m ConnectMessage) Format(enc *encoder) {

	size := len(m.params)
	if size == 0 {
		enc.encodeNull()
		return
	}

	enc.encodeNullableSize(size)

	for key, value := range m.params {

		enc.encodeString(key)
		enc.encodeTypedValue(value)

	}

}

type NegotiateMessage struct {
	magic           int
	ProtocolVersion int
	CodecVersion    int
}

func (n NegotiateMessage) Type() MessageType {
	return NEGOTIATE
}

func (n NegotiateMessage) Format(enc *encoder) {

	enc.encodeInt(n.magic)
	enc.encodeShort(n.ProtocolVersion)
	enc.encodeShort(n.ProtocolVersion)

}

func NewNegotiateMessage(protocol, codec int) NegotiateMessage {
	return NegotiateMessage{
		magic:           magic,
		ProtocolVersion: protocol,
		CodecVersion:    codec,
	}
}

const endpointPrefix = "v8.service.Admin.Cluster"
const endpointParamPrefix = "endpoint."

type OpenEndpointMessage struct {
	Encoding string
	Version  string
	params   map[string]interface{}
}

func (m *OpenEndpointMessage) String() string {
	return ""
}

func (m *OpenEndpointMessage) Type() MessageType {
	return ENDPOINT_OPEN
}

func (m OpenEndpointMessage) Format(enc *encoder) {

	enc.encodeString(endpointPrefix)
	enc.encodeString(m.Version)
	size := len(m.params)
	if size == 0 {
		enc.encodeNull()
		return
	}

	enc.encodeNullableSize(size)

	for key, value := range m.params {

		enc.encodeString(key)
		enc.encodeTypedValue(value)

	}

}

type OpenEndpointMessageAck struct {
	ServiceID  string
	Version    string
	EndpointID int

	params map[string]interface{}
}

func (m *OpenEndpointMessageAck) Parse(body []byte) error {
	//panic("implement me")

	dec := NewDecoder(body)
	m.ServiceID = dec.decodeString()
	m.Version = dec.decodeString()
	m.EndpointID = dec.decodeEndpointId()

	// TODO params

	return nil
}

func (m *OpenEndpointMessageAck) String() string {
	return ""
}

func (m *OpenEndpointMessageAck) Type() MessageType {
	return ENDPOINT_OPEN_ACK
}

func (m OpenEndpointMessageAck) Format(enc *encoder) {

	size := len(m.params)
	if size == 0 {
		enc.encodeNull()
		return
	}

	enc.encodeNullableSize(size)

	for key, value := range m.params {

		enc.encodeString(key)
		enc.encodeTypedValue(value)

	}

}

type EndpointFeature struct {
	ServiceID  string
	Version    string
	EndpointID string
	trace      string
	Error      error
}

type causeError struct {
	service string
	msg     string
}

func (e *causeError) Error() string {

	return pp.Sprintf("service: %s err: %s", e.service, e.msg)

}

func (m *EndpointFeature) Parse(body []byte) error {

	dec := NewDecoder(body)
	m.ServiceID = dec.decodeString()
	m.Version = dec.decodeString()

	m.EndpointID = dec.decodeString()
	classError := dec.decodeString()

	pp.Printf(classError)

	errMessage := dec.decodeString()
	errSize := dec.decodeSize()

	pp.Printf(errMessage, errSize)

	if errSize > 0 {

		panic("TODO ")

	}

	causeService := dec.decodeString()
	causeMessage := dec.decodeString()

	m.Error = &causeError{
		service: causeService,
		msg:     causeMessage,
	}

	return nil
}

func (m *EndpointFeature) String() string {
	return m.Error.Error()
}

func (m *EndpointFeature) Type() MessageType {
	return ENDPOINT_FAILURE
}

func (m EndpointFeature) Format(enc *encoder) {

}

type EndpointMessage struct {
	raw []byte

	endpointID int
	format     int

	kind EndpointMessageKind

	respondType MessageType
	Respond     map[MessageType]RespondMessage
}

func (m *EndpointMessage) addResponse(r ...RespondMessage) {

	if m.Respond == nil {
		m.Respond = make(map[MessageType]RespondMessage)
	}
	for _, message := range r {
		m.Respond[message.Type()] = message
	}

}

func (m *EndpointMessage) Parse(body []byte) error {

	decoder := NewDecoder(body)
	m.raw = body

	m.endpointID = decoder.decodeEndpointId()
	m.format = int(decoder.decodeShort())

	m.kind = EndpointMessageKind(decoder.decodeByte())
	m.respondType = EndpointMessageType(decoder.decodeUnsignedByte())

	respBody, err := ioutil.ReadAll(decoder) ///Читаем то что осталось

	if err != nil {
		return err
	}

	typedFormat, ok := m.Respond[m.respondType]

	if ok {

		err = typedFormat.Parse(respBody)
		if err != nil {
			return err
		}

	}

	return nil
}

func (m *EndpointMessage) String() string {
	return ""
}

func (m *EndpointMessage) Type() MessageType {
	return ENDPOINT_MESSAGE
}

func (m EndpointMessage) Format(enc *encoder) {

}
