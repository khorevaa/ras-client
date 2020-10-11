package protocol

import (
	"github.com/k0kubun/pp"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/types"
	"io"
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

func (r *ConnectMessageAck) Type() types.Typed {
	return CONNECT_ACK
}

func (r *ConnectMessageAck) Parse(codec codec.Decoder, w io.Reader) {

}

type ConnectMessage struct {
	params   map[string]interface{}
	response *ConnectMessageAck
}

func (m *ConnectMessage) ResponseMessage() types.ResponseMessage {

	if m.response == nil {
		m.response = &ConnectMessageAck{}
	}

	return m.response
}

func (m *ConnectMessage) String() string {
	return ""
}

func (m *ConnectMessage) Type() types.Typed {
	return CONNECT
}

func (m ConnectMessage) Format(c codec.Encoder, w io.Writer) {

	size := len(m.params)
	if size == 0 {
		c.Null(w)
		return
	}

	c.NullableSize(size, w)

	for key, value := range m.params {

		c.String(key, w)
		c.TypedValue(value, w)

	}

}

type NegotiateMessage struct {
	magic           int
	ProtocolVersion int16
	CodecVersion    int16
}

func (n NegotiateMessage) ResponseMessage() types.ResponseMessage {
	return &nullRespondMessage{}
}

func (n NegotiateMessage) Type() types.Typed {
	return NEGOTIATE
}

func (n NegotiateMessage) Format(c codec.Encoder, w io.Writer) {

	c.Int(n.magic, w)
	c.Short(n.ProtocolVersion, w)
	c.Short(n.ProtocolVersion, w)

}

func NewNegotiateMessage(protocol, codec int16) NegotiateMessage {
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
	ack      *OpenEndpointMessageAck
}

func (m *OpenEndpointMessage) String() string {
	return pp.Sprintln(m)
}

func (m *OpenEndpointMessage) Type() types.Typed {
	return ENDPOINT_OPEN
}
func (m *OpenEndpointMessage) ResponseMessage() types.ResponseMessage {

	if m.ack == nil {
		m.ack = &OpenEndpointMessageAck{}
	}

	return m.ack
}

func (m OpenEndpointMessage) Format(c codec.Encoder, w io.Writer) {

	c.String(endpointPrefix, w)
	c.String(m.Version, w)
	size := len(m.params)
	if size == 0 {
		c.Null(w)
		return
	}

	c.NullableSize(size, w)

	for key, value := range m.params {

		c.String(key, w)
		c.TypedValue(value, w)

	}

}

type OpenEndpointMessageAck struct {
	ServiceID  string
	Version    string
	EndpointID int

	params map[string]interface{}
}

func (m *OpenEndpointMessageAck) Parse(c codec.Decoder, r io.Reader) {

	c.StringPtr(&m.ServiceID, r)
	c.StringPtr(&m.Version, r)

	m.EndpointID = c.EndpointId(r)

	// TODO params

}

func (m *OpenEndpointMessageAck) String() string {
	return pp.Sprintln(m)
}

func (m *OpenEndpointMessageAck) Type() types.Typed {
	return ENDPOINT_OPEN_ACK
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

func (m *EndpointFailure) Parse(c codec.Decoder, r io.Reader) {

	c.StringPtr(&m.ServiceID, r)
	c.StringPtr(&m.Version, r)

	m.EndpointID = c.EndpointId(r)

	classError := c.String(r)

	pp.Printf(classError)

	errMessage := c.String(r)
	errSize := c.Size(r)

	pp.Printf(errMessage, errSize)

	if errSize > 0 {

		panic("TODO ")

	}

	causeService := c.String(r)
	causeMessage := c.String(r)

	m.err = &causeError{
		service: causeService,
		msg:     causeMessage,
	}
}

func (m *EndpointFailure) String() string {
	return m.err.Error()
}

func (m *EndpointFailure) Type() types.Typed {
	return ENDPOINT_FAILURE
}

func (m *EndpointFailure) Error() string {

	return m.err.Error()
}

//
//type EndpointMessageFailure struct {
//	serviceID  string
//	Message    string
//	EndpointID int
//}
//
//func (m *EndpointMessageFailure) Parse(body []byte) error {
//
//	dec := NewDecoder(body)
//	m.serviceID = dec.decodeString()
//	m.Message = dec.decodeString()
//
//	respBody, err := ioutil.ReadAll(dec) ///Читаем то что осталось
//
//	if err != nil {
//		return err
//	}
//
//	pp.Println("EndpointMessageFailure", respBody)
//
//	return nil
//}
//
//func (m *EndpointMessageFailure) String() string {
//	return m.Message
//}
//
//func (m *EndpointMessageFailure) Type() Typed {
//	return EXCEPTION_KIND
//}
//
//func (m *EndpointMessageFailure) Error() string {
//	return pp.Sprintf("endpoint: %s service: %s msg: %s", m.EndpointID, m.serviceID, m.Message)
//}

//type EndpointMessage struct {
//	raw []byte
//
//	endpointID int
//	format     int
//
//	kind EndpointMessageKind
//
//	respondType  Typed
//	waitResponse ResponseMessage
//	err          *EndpointMessageFailure
//}
//
//func (m *EndpointMessage) WaitResponse(r ResponseMessage) {
//
//	m.waitResponse = r
//
//}
//
//func (m *EndpointMessage) Parse(body []byte) error {
//
//	decoder := NewDecoder(body)
//	m.raw = body
//	m.endpointID = decoder.decodeEndpointId()
//	m.format = int(decoder.decodeShort())
//	m.kind = EndpointMessageKind(decoder.decodeByte())
//
//	switch m.kind {
//
//	case VOID_MESSAGE_KIND:
//		return nil
//	case EXCEPTION_KIND:
//
//		respBody, err := ioutil.ReadAll(decoder) ///Читаем то что осталось
//
//		if err != nil {
//			return err
//		}
//
//		m.err = &EndpointMessageFailure{EndpointID: m.endpointID}
//		err = m.err.Parse(respBody)
//
//		if err != nil {
//			return err
//		}
//
//	case MESSAGE_KIND:
//
//		m.respondType = EndpointMessageType(decoder.decodeUnsignedByte())
//
//		respBody, err := ioutil.ReadAll(decoder) ///Читаем то что осталось
//
//		if err != nil {
//			return err
//		}
//
//		if m.respondType != m.waitResponse.Type() {
//			return errors.New("не совпадает ожидаем и тип полученныго ответа")
//		}
//
//		err = m.waitResponse.Parse(respBody)
//		if err != nil {
//			return err
//		}
//
//	default:
//		dry.PanicIf(true, "неизвестный тип сообщения ответа")
//	}
//
//	return nil
//}
//
//func (m *EndpointMessage) String() string {
//	return ""
//}
//
//func (m *EndpointMessage) Type() Typed {
//	return ENDPOINT_MESSAGE
//}
//
//func (m EndpointMessage) Format(enc *encoder) {
//
//}
