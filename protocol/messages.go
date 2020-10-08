package protocol

import (
	"errors"
	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
	"io/ioutil"
)

type EndpointMessageKind int

func (e EndpointMessageKind) Type() int {
	return int(e)
}

const (
	VOID_MESSAGE_KIND EndpointMessageKind = 0
	MESSAGE_KIND      EndpointMessageKind = 1
	EXCEPTION_KIND    EndpointMessageKind = 0xff
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

type EndpointFailure struct {
	ServiceID  string
	Version    string
	EndpointID int
	trace      string
	err        error
}

type causeError struct {
	service string
	msg     string
}

func (e *causeError) Error() string {

	return pp.Sprintf("service: %s err: %s", e.service, e.msg)

}

func (m *EndpointFailure) Parse(body []byte) error {

	dec := NewDecoder(body)
	m.ServiceID = dec.decodeString()
	m.Version = dec.decodeString()

	m.EndpointID = dec.decodeEndpointId()
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

	m.err = &causeError{
		service: causeService,
		msg:     causeMessage,
	}

	return nil
}

func (m *EndpointFailure) String() string {
	return m.err.Error()
}

func (m *EndpointFailure) Type() MessageType {
	return ENDPOINT_FAILURE
}

func (m *EndpointFailure) Error() string {

	return m.err.Error()
}

type EndpointMessageFailure struct {
	ServiceID  string
	Message    string
	EndpointID int
}

func (m *EndpointMessageFailure) Parse(body []byte) error {

	dec := NewDecoder(body)
	m.ServiceID = dec.decodeString()
	m.Message = dec.decodeString()

	respBody, err := ioutil.ReadAll(dec) ///Читаем то что осталось

	if err != nil {
		return err
	}

	pp.Println("EndpointMessageFailure", respBody)

	return nil
}

func (m *EndpointMessageFailure) String() string {
	return m.Message
}

func (m *EndpointMessageFailure) Type() MessageType {
	return EXCEPTION_KIND
}

func (m *EndpointMessageFailure) Error() string {
	return pp.Sprintf("endpoint: %s service: %s msg: %s", m.EndpointID, m.ServiceID, m.Message)
}

type EndpointMessage struct {
	raw []byte

	endpointID int
	format     int

	kind EndpointMessageKind

	respondType  MessageType
	waitResponse RespondMessage
	err          *EndpointMessageFailure
}

func (m *EndpointMessage) WaitResponse(r RespondMessage) {

	m.waitResponse = r

}

func (m *EndpointMessage) Parse(body []byte) error {

	decoder := NewDecoder(body)
	m.raw = body
	m.endpointID = decoder.decodeEndpointId()
	m.format = int(decoder.decodeShort())
	m.kind = EndpointMessageKind(decoder.decodeByte())

	switch m.kind {

	case VOID_MESSAGE_KIND:
		return nil
	case EXCEPTION_KIND:

		respBody, err := ioutil.ReadAll(decoder) ///Читаем то что осталось

		if err != nil {
			return err
		}

		m.err = &EndpointMessageFailure{EndpointID: m.endpointID}
		err = m.err.Parse(respBody)

		if err != nil {
			return err
		}

	case MESSAGE_KIND:

		m.respondType = EndpointMessageType(decoder.decodeUnsignedByte())

		respBody, err := ioutil.ReadAll(decoder) ///Читаем то что осталось

		if err != nil {
			return err
		}

		if m.respondType != m.waitResponse.Type() {
			return errors.New("не совпадает ожидаем и тип полученныго ответа")
		}

		err = m.waitResponse.Parse(respBody)
		if err != nil {
			return err
		}

	default:
		dry.PanicIf(true, "неизвестный тип сообщения ответа")
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
