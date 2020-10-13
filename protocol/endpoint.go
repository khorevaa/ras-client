package protocol

import (
	"bytes"
	"context"
	"errors"
	"github.com/k0kubun/pp"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/messages"
	"github.com/v8platform/rac/protocol/types"
	"io"
	"strconv"
)

const defaultFormat = 256

type endpoint struct {
	conn      *Client
	Id        int
	serviceID string
	version   string
	format    int16

	codec codec.Codec

	messages chan EndpointMessage
	wait     chan chan EndpointMessage

	opened bool
	ctx    context.Context
}

func (e *endpoint) Close() {
	panic("implement me")
}

func newEndpoint(conn *Client, id int, serviceID string, version string) *endpoint {

	end := &endpoint{
		conn:      conn,
		Id:        id,
		serviceID: serviceID,
		version:   version,
		format:    defaultFormat,
		codec:     conn.codec,
		messages:  make(chan EndpointMessage),
		wait:      make(chan chan EndpointMessage),
		opened:    true,
		ctx:       conn.ctx,
	}

	end.processMessages(conn.ctx)

	return end
}

func (e *endpoint) Version() int {

	v, err := strconv.ParseFloat(e.version, 10)
	if err != nil {
		panic(err)
	}

	return int(v)
}

func (e *endpoint) newMessage(req types.EndpointRequestMessage) *EndpointMessage {

	endpointMessage := &EndpointMessage{
		kind:     req.Kind(),
		endpoint: e,
		req:      req,
	}

	return endpointMessage

}

func (e *endpoint) SendMessage(req types.EndpointRequestMessage) (interface{}, error) {

	if !e.opened {
		return nil, errors.New("endpoint is not opened")
	}

	m := e.newMessage(req)

	_, err := e.conn.SendRequest(m)

	if err != nil {
		return nil, err
	}

	if req.Kind() == messages.VOID_MESSAGE_KIND {
		return nil, err
	}

	return e.waitAck(req)

}

func (e *endpoint) waitAck(req types.EndpointRequestMessage) (interface{}, error) {

	wait := make(chan EndpointMessage)
	e.wait <- wait

	select {
	case <-e.ctx.Done():

		return nil, e.ctx.Err()

	case m := <-wait:

		r := m.GetMessage()

		resp := req.ResponseMessage()

		err := e.tryParse(resp.Type(), resp, r)

		return resp, err

	}

}

func (e *endpoint) tryParse(t types.Typed, p codec.BinaryParser, r io.Reader) (err error) {
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

	decoder := e.codec.Decoder()

	kind := messages.EndpointMessageKind(decoder.Type(r))

	switch kind {

	case messages.VOID_MESSAGE_KIND:
		return
	case messages.EXCEPTION_KIND:

		fail := &EndpointMessageFailure{EndpointID: e.Id}
		fail.Parse(decoder, r)
		return fail

	case messages.MESSAGE_KIND:

		respondType := decoder.Type(r)

		if t.Type() != respondType {
			//pp.Println("decoded type", respondType)
			return
		}

		p.Parse(decoder, e.Version(), r)
	}

	return
}

func (e *endpoint) Format(c codec.Encoder, w io.Writer) error {

	c.EndpointId(e.Id, w)
	c.Short(e.format, w)

	return nil
}

func (e *endpoint) PushMessage(message EndpointMessage) {
	e.messages <- message
}

func (e *endpoint) processMessages(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				// TODO Сделать закрытие всех ожидащих каналов
				return
			case message := <-e.messages:

				//if len(e.wait) == 0 {
				//	pp.Println(message)
				//	return
				//}

				wait := <-e.wait
				wait <- message
			}
		}
	}()
}

type EndpointMessageFailure struct {
	ServiceID  string
	Message    string
	EndpointID int
}

func (m *EndpointMessageFailure) Parse(decoder codec.Decoder, r io.Reader) {

	decoder.StringPtr(&m.ServiceID, r)
	decoder.StringPtr(&m.Message, r)

	//pp.Println("EndpointMessageFailure", m)

}

func (m *EndpointMessageFailure) String() string {
	return pp.Sprintln(m)
}

func (m *EndpointMessageFailure) Type() types.Typed {
	return messages.EXCEPTION_KIND
}

func (m *EndpointMessageFailure) Error() string {
	return pp.Sprintf("endpoint: %s service: %s msg: %s", m.EndpointID, m.ServiceID, m.Message)
}

type EndpointMessage struct {
	endpoint *endpoint
	kind     types.Typed

	bufReader *bytes.Reader
	req       types.EndpointRequestMessage
}

func (m *EndpointMessage) ResponseMessage() types.ResponseMessage {
	return nil // Для данного типа мы ничего не возвращает, отдельно обрабатывается
}

func (m *EndpointMessage) Type() types.Typed {
	return ENDPOINT_MESSAGE
}

func (m *EndpointMessage) GetMessage() *bytes.Reader {
	return m.bufReader
}

func (m *EndpointMessage) Parse(_ codec.Decoder, _ io.Reader) {}

func (m *EndpointMessage) String() string {
	return pp.Sprintln(m)
}

func (m *EndpointMessage) Format(encoder codec.Encoder, w io.Writer) {

	encoder.EndpointId(m.endpoint.Id, w)
	//encoder.Short(m.endpoint.format, w)
	encoder.Short(0, w)
	encoder.Type(m.kind, w)
	encoder.Type(m.req.Type(), w) // МАГИЯ без этого байта требует авторизации на центральном кластере

	m.req.Format(encoder, m.endpoint.Version(), w) // запись тебя сообщения

}
