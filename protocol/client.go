package protocol

import (
	"bytes"
	"context"
	"github.com/k0kubun/pp"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/internal/pool"
	"github.com/v8platform/rac/protocol/types"
	"io"
	"net"
	"strconv"
	"time"

	bus "github.com/asaskevich/EventBus"
	//"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
)

const readTimeout = time.Second * 15
const magic = 475223888
const protocolVersion = 256
const maxResponseChunkSize = 1460

type Client struct {
	addr  string
	laddr net.Addr

	ctx          context.Context
	stopRoutines context.CancelFunc // остановить ping, read, и подобные горутины

	pool pool.EndpointPool

	// каналы, которые ожидают ответа rpc
	ackRespond chan ackRespond

	responses chan rawResponse
	endpoints map[int]*endpoint

	codec codec.Codec

	serviceVersion string

	// шина соо
	//бщений, используется для разных нотификаций, описанных в константах нотификации
	bus bus.Bus
}

func NewClient(addr string) *Client {

	m := new(Client)
	m.addr = addr
	m.codec = codec.NewCodec1_0()
	m.pool = pool.NewEndpointPool(&pool.Options{
		Dialer:             m.dialfunc,
		OpenEndpoint:       m.openEndpoint,
		CloseEndpoint:      m.closeEndpoint,
		InitConnection:     m.initConnection,
		PoolSize:           10,
		MinIdleConns:       2,
		MaxConnAge:         time.Hour,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 2 * time.Minute,
	})

	m.ackRespond = make(chan ackRespond)
	m.responses = make(chan rawResponse)
	m.endpoints = make(map[int]*endpoint)
	m.serviceVersion = serviceVersions[len(serviceVersions)-1]
	m.resetAck()

	return m
}

func (c *Client) initConnection(ctx context.Context, conn *pool.Conn) error {

	negotiateMessage := NewNegotiateMessage(protocolVersion, c.codec.Version())

	err := c.sendRequest(conn, negotiateMessage)

	if err != nil {
		return err
	}

	err = c.sendRequest(conn, &ConnectMessage{params: map[string]interface{}{
		"connect.timeout": int64(2000),
	}})

	packet, err := conn.GetPacket(ctx)

	if err != nil {
		return err
	}

	answer, err := c.tryParseMessage(packet)

	if err != nil {
		return err
	}

	if _, ok := answer.(*ConnectMessageAck); !ok {
		return errors.New("unknown ack")
	}

	return nil
}

func (c *Client) openEndpoint(ctx context.Context, conn *pool.Conn) (info pool.EndpointInfo, err error) {

	var ack *OpenEndpointMessageAck

	ack, err = c.tryOpenEndpoint(ctx, conn)
	if err != nil {

		supportedVersion := detectSupportedVersion(err)
		if len(supportedVersion) > 0 {
			return nil, errors.New(pp.Sprint("ras no supported service version", serviceVersions))
		}

		c.serviceVersion = supportedVersion
		ack, err = c.tryOpenEndpoint(ctx, conn)
	}

	if err != nil {
		return nil, err
	}

	endpointVersion, err := strconv.ParseFloat(ack.Version, 10)
	if err != nil {
		return nil, err
	}

	return endpointInfo{
		id:        ack.EndpointID,
		version:   int(endpointVersion),
		format:    defaultFormat,
		serviceID: ack.ServiceID,
		codec:     c.codec,
	}, nil
}

type endpointInfo struct {
	id        int
	version   int
	format    int16
	serviceID string
	codec     codec.Codec
}

func (e endpointInfo) ID() int {
	return e.id
}

func (e endpointInfo) Version() int {
	return e.version
}

func (e endpointInfo) Format() int16 {
	return e.format
}

func (e endpointInfo) ServiceID() string {
	return e.serviceID
}

func (e endpointInfo) Codec() codec.Codec {
	return e.codec
}

func (c *Client) tryOpenEndpoint(ctx context.Context, conn *pool.Conn) (*OpenEndpointMessageAck, error) {

	err := c.sendRequest(conn, &OpenEndpointMessage{Version: c.serviceVersion})

	packet, err := conn.GetPacket(ctx)

	if err != nil {
		return nil, err
	}

	answer, err := c.tryParseMessage(packet)

	if err != nil {
		return nil, err
	}

	switch t := answer.(type) {

	case *EndpointFailure:

		return nil, t

	case *OpenEndpointMessageAck:

		return t, nil

	default:
		panic("unknown answer type")
	}

}

func (c *Client) closeEndpoint(ctx context.Context, conn *pool.Conn, endpoint *pool.Endpoint) error {

	return nil
}
func (c *Client) sendRequest(conn *pool.Conn, message types.RequestMessage) error {

	body := bytes.NewBuffer([]byte{})
	message.Format(c.codec.Encoder(), body)
	packet := pool.NewPacket(byte(message.Type().Type()), body.Bytes())

	err := conn.SendPacket(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) tryParseMessage(packet *pool.Packet) (message types.ResponseMessage, err error) {
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

	case CONNECT_ACK.Type():

		decoder := c.codec.Decoder()

		message = &ConnectMessageAck{}
		message.Parse(decoder, packet)

	case KEEP_ALIVE.Type():
		// nothing
	case ENDPOINT_OPEN_ACK.Type():

		decoder := c.codec.Decoder()

		message = &OpenEndpointMessageAck{}
		message.Parse(decoder, packet)

	case ENDPOINT_FAILURE.Type():

		decoder := c.codec.Decoder()

		message = &EndpointFailure{}
		message.Parse(decoder, packet)

	case NULL_TYPE.Type():

		panic(NULL_TYPE.String())

	default:

		panic(packet)
	}

	return
}

func (c *Client) dialfunc(ctx context.Context) (net.Conn, error) {

	_, err := net.ResolveTCPAddr("tcp", c.addr)
	if err != nil {
		return nil, errors.Wrap(err, "resolving tcp")
	}

	var dialer net.Dialer

	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return nil, errors.Wrap(err, "dialing tcp")
	}

	return conn, nil

}

func (c *Client) CreateConnection() error {

	_, err := net.ResolveTCPAddr("tcp", c.addr)
	if err != nil {
		return errors.Wrap(err, "resolving tcp")
	}

	c.conn, err = net.Dial("tcp", c.addr)
	if err != nil {
		return errors.Wrap(err, "dialing tcp")
	}
	c.laddr = c.conn.LocalAddr()

	ctx, cancelfunc := context.WithCancel(context.Background())
	c.stopRoutines = cancelfunc
	c.ctx = ctx

	// start reading responses from the server
	c.startReadingResponses(ctx)

	// start keepalive pinging
	c.startRoutingResponses(ctx)

	_, err = c.SendRequest(NewNegotiateMessage(protocolVersion, c.codec.Version()))

	if err != nil {
		return err
	}

	_, err = c.SendRequest(&ConnectMessage{params: map[string]interface{}{
		"connect.timeout": int64(2000),
	}})

	return err
}

func (c *Client) SendRequest(req types.RequestMessage) (interface{}, error) {

	buf := NewBuffer()
	err := c.formatRequestMessage(req, buf)
	if err != nil {
		return nil, err
	}

	err = c.sendPacket(buf)

	if err != nil {
		return nil, err
	}

	if req.ResponseMessage() != nil {
		resp, err := c.waitAck(req)

		return resp, err
	}

	return nil, err

}

func (c *Client) Disconnect() error {
	// stop all routines
	c.stopRoutines()

	err := c.conn.Close()
	if err != nil {
		return errors.Wrap(err, "closing TCP connection")
	}

	// TODO: закрыть каналы

	// возвращаем в false, потому что мы теряем конфигурацию
	// сессии, и можем ее потерять во время отключения.

	return nil
}

type ackRespond struct {
	req  types.RequestMessage
	wait chan interface{}
}

func newAckRespond(req types.RequestMessage, wait chan interface{}) ackRespond {

	return ackRespond{
		req:  req,
		wait: wait,
	}

}

func (c *Client) resetAck() {
	c.ackRespond = make(chan ackRespond, 1)
}

// waitAck добавляет в список id сообщения, которому нужно подтверждение
// возвращает true, если ранее этого id не было
func (c *Client) waitAck(req types.RequestMessage) (interface{}, error) {

	wait := make(chan interface{})
	c.ackRespond <- newAckRespond(req, wait)

	select {
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	case resp := <-wait:

		switch typed := resp.(type) {

		case error:

			return nil, typed

		default:

			return typed, nil

		}
	}

}

func (c *Client) sendPacket(request *bytes.Buffer) error {

	_, err := request.WriteTo(c.conn)
	if err != nil {
		return errors.Wrap(err, "sending ack")
	}
	return nil
}

func (c *Client) formatRequestMessage(req types.RequestMessage, buf *bytes.Buffer) (err error) {

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

	e := c.codec.Encoder()

	switch req.Type() {

	case NEGOTIATE:

		req.Format(e, buf)

	default:

		body := NewBuffer()
		req.Format(e, body)

		e.Type(req.Type(), buf)
		e.Size(body.Len(), buf)

		_, err = body.WriteTo(buf)

	}

	return
}

func (c *Client) startReadingResponses(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:

				rawResp, err := c.readFromConn(ctx)
				if err != nil {
					pp.Errorf("error while reading from connection")
					continue
				}

				if rawResp.size == 0 {
					continue
				}

				c.responses <- rawResp

			}
		}
	}()
}

func (c *Client) readFromConn(ctx context.Context) (reso rawResponse, err error) {

	//c.conn.SetReadDeadline(readTimeout)

	header := make([]byte, 20)
	n, err := c.conn.Read(header)
	if n == 0 {
		return rawResponse{}, nil
	}

	dry.PanicIfErr(err)
	header = header[:n]

	buf := bytes.NewReader(header)
	dec := c.codec.Decoder()
	messageType := dec.Type(buf)
	size := dec.Size(buf)

	data := make([]byte, size)

	if size > maxResponseChunkSize {

		n = 0

		for size-n > 0 {

			buf := make([]byte, maxResponseChunkSize)
			offset := n
			nReaded, err := c.conn.Read(buf)
			dry.PanicIfErr(err)
			copy(data[offset:offset+nReaded], buf[:nReaded])

			n += nReaded

		}
	} else {
		reader := dry.NewCancelableReader(ctx, c.conn)

		n, err = reader.Read(data)

	}

	dry.PanicIfErr(err)
	dry.PanicIf(n != int(size), "expected read "+strconv.Itoa(int(size))+" bytes, got "+strconv.Itoa(n))

	resp := rawResponse{
		ConnectionMessageType(messageType),
		size,
		data,
	}

	return resp, nil // TODO Переделать на возврат сообщения
}

func (c *Client) startRoutingResponses(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case res := <-c.responses:

				switch res.t {

				case ENDPOINT_MESSAGE:

					c.receiveEndpointMessage(res)

				case ENDPOINT_FAILURE:

					panic(pp.Sprintln(string(res.payload)))

				case KEEP_ALIVE:

					pp.Println(KEEP_ALIVE)

				case NULL_TYPE:
					// Nothing to do
				default:

					c.receiveResponse(res)
				}

			}
		}
	}()
}

func (c *Client) receiveResponse(raw rawResponse) {

	//if len(c.ackRespond) == 0 {
	//	return
	//}

	ack := <-c.ackRespond

	resp := ack.req.ResponseMessage()
	d := c.codec.Decoder()
	r := bytes.NewReader(raw.payload)

	switch raw.Type().Type() {

	case resp.Type().Type():

		err := tryParse(resp, d, r)

		if err != nil {
			ack.wait <- err
			return
		}

		ack.wait <- resp

	default:
		return

	}
}

func tryParse(p types.ResponseMessage, decoder codec.Decoder, r io.Reader) (err error) {
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

	p.Parse(decoder, r)

	return
}

func (c *Client) registryNewEndpoint(ack *OpenEndpointMessageAck) types.Endpoint {

	end := newEndpoint(c, ack.EndpointID, ack.ServiceID, ack.Version)

	c.endpoints[end.Id] = end
	return end
}

func (c *Client) receiveEndpointMessage(res rawResponse) {

	d := c.codec.Decoder()
	r := bytes.NewReader(res.payload)
	endpointID := d.EndpointId(r)
	_ = d.Short(r) // Format уже записан в точке

	receiver, ok := c.endpoints[endpointID]

	if !ok {

		pp.Println("Не удалось определить точку получения сообщения", endpointID)
		return
	}

	message := EndpointMessage{
		endpoint:  receiver,
		bufReader: r,
	}

	receiver.PushMessage(message)

}

type rawResponse struct {
	t       types.Typed
	size    int
	payload []byte
}

func (r rawResponse) Type() types.Typed {
	return r.t
}

func (r rawResponse) Len() int {
	return r.size
}

func (r rawResponse) Size() int {
	return r.size
}

func (r rawResponse) Data() []byte {
	return r.payload
}

type nullRespondMessage struct {
}

func (_ nullRespondMessage) Parse(_ codec.Decoder, _ io.Reader) {

}

func (_ nullRespondMessage) Type() types.Typed {
	return NULL_TYPE
}
