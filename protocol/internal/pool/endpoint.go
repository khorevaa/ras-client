package pool

import (
	"bytes"
	"context"
	"errors"
	"github.com/k0kubun/pp"
	"github.com/v8platform/rac/protocol"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/esig"
	"github.com/v8platform/rac/protocol/messages"
	"github.com/v8platform/rac/protocol/types"
	"github.com/v8platform/rac/serialize"
	"io"
	"sync/atomic"
	"time"
)

func NewEndpoint(endpoint EndpointInfo) *Endpoint {

	return &Endpoint{
		id:        endpoint.ID(),
		version:   endpoint.Version(),
		format:    endpoint.Format(),
		serviceID: endpoint.ServiceID(),
	}
}

type EndpointInfo interface {
	ID() int
	Version() int
	Format() int16
	ServiceID() string
	Codec() codec.Codec
}

type Endpoint struct {
	id        int
	version   int
	format    int16
	serviceID string
	codec     codec.Codec

	conn      *Conn
	createdAt time.Time
	usedAt    uint32 // atomic
	pooled    bool
	Inited    bool

	sig       esig.ESIG
	onRequest func(ctx context.Context, endpoint *Endpoint, req types.EndpointRequestMessage) error
}

func (e *Endpoint) Sig() esig.ESIG {
	return e.sig
}

func (e *Endpoint) SetSig(sig esig.ESIG) {
	e.sig = sig
}

func (cn *Endpoint) UsedAt() time.Time {
	unix := atomic.LoadUint32(&cn.usedAt)
	return time.Unix(int64(unix), 0)
}

func (cn *Endpoint) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&cn.usedAt, uint32(tm.Unix()))
}

func (e *Endpoint) ID() int {
	panic("implement me")
}

func (e *Endpoint) Version() int {
	panic("implement me")
}

func (e *Endpoint) Format() string {
	panic("implement me")
}

func (e *Endpoint) ServiceID() string {
	panic("implement me")
}

func (e *Endpoint) Codec() codec.Codec {
	panic("implement me")
}

type UnknownMessageError struct {
	Type     byte
	Data     []byte
	Endpoint *Endpoint
	err      error
}

func (m *UnknownMessageError) Error() string {

	return m.err.Error()

}

func (e *Endpoint) sendRequest(ctx context.Context, message *EndpointMessage) (*EndpointMessage, error) {

	body := bytes.NewBuffer([]byte{})

	message.Format(e.codec.Encoder(), e.version, body)

	packet := NewPacket(byte(message.Type.Type()), body.Bytes())

	err := e.conn.SendPacket(packet)
	if err != nil {
		return nil, err
	}

	answer, err := e.conn.GetPacket(ctx)

	if err != nil {
		return nil, err
	}

	return e.tryParseMessage(answer)

}

func (e *Endpoint) sendVoidRequest(ctx context.Context, conn *Conn, m protocol.EndpointMessage) error {

	body := bytes.NewBuffer([]byte{})

	m.Format(e.codec.Encoder(), e.version, body)

	packet := NewPacket(byte(m.Type().Type()), body.Bytes())

	err := conn.SendPacket(packet)
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) tryParseMessage(packet *Packet) (message *EndpointMessage, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch val := e.(type) {

			case string:

				err = errors.New(val)

			case error:
				err = val
			default:
				panic(e)
			}
		}
	}()

	switch int(packet.Type) {

	case protocol.ENDPOINT_MESSAGE.Type():

		decoder := e.codec.Decoder()

		endpointID := decoder.EndpointId(packet)
		format := decoder.Short(packet)

		message = &EndpointMessage{
			EndpointID:     endpointID,
			EndpointFormat: format,
		}

		message.Parse(decoder, e.version, packet)

	case protocol.ENDPOINT_FAILURE.Type():

		panic(pp.Sprintln(string(packet.Data))) // TODO Гдето есть парсер

	default:

		return nil, &UnknownMessageError{
			packet.Type,
			packet.Data,
			e,
			ErrUnknownMessage}
	}

	return
}

func (e *Endpoint) tryFormatMessage(message *EndpointMessage, writer io.Writer) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch val := e.(type) {

			case string:

				err = errors.New(val)

			case error:
				err = val
			default:
				panic(e)
			}
		}
	}()

	encoder := e.codec.Encoder()
	message.Format(encoder, e.version, writer)

	return
}

func (m *EndpointMessage) Parse(decoder codec.Decoder, version int, reader io.Reader) {

	kind := messages.EndpointMessageKind(decoder.Type(reader))
	m.Kind = kind

	switch kind {

	case messages.VOID_MESSAGE_KIND:
		return
	case messages.EXCEPTION_KIND:

		fail := &protocol.EndpointMessageFailure{EndpointID: m.EndpointID}
		fail.Parse(decoder, reader)
		m.Message = fail

	case messages.MESSAGE_KIND:

		respondType := decoder.Type(reader)
		pp.Println(respondType)

		var parser codec.BinaryParser
		// TODO Сделать получение ответа по типу
		parser.Parse(decoder, version, reader)

		m.Message = parser
	}

	return

}

func (m *EndpointMessage) Format(encoder codec.Encoder, version int, w io.Writer) {

	encoder.EndpointId(m.EndpointID, w)
	encoder.Short(m.EndpointFormat, w)
	encoder.Type(m.Kind, w)
	encoder.Type(m.Type, w) // МАГИЯ без этого байта требует авторизации на центральном кластере

	formatter := m.Message.(codec.BinaryWriter)
	formatter.Format(encoder, version, w) // запись тебя сообщения

}

type EndpointMessage struct {
	EndpointID     int
	EndpointFormat int16
	Kind           messages.EndpointMessageKind

	Message interface{}
	Type    serialize.Typed
}

func (e *Endpoint) SendRequest(ctx context.Context, req types.EndpointRequestMessage) (*EndpointMessage, error) {

	if e.onRequest != nil {

		err := e.onRequest(ctx, e, req)

		if err != nil {
			return nil, err
		}

	}

	message := e.newEndpointMessage(req)
	answer, err := e.sendRequest(ctx, message)

	return answer, err

}

func (e *Endpoint) newEndpointMessage(req types.EndpointRequestMessage) *EndpointMessage {

	message := &EndpointMessage{
		EndpointID:     e.id,
		EndpointFormat: e.format,
		Kind:           messages.EndpointMessageKind(req.Type().Type()),
		Message:        req,
	}

	return message

}
