package protocol

import (
	"bytes"
	"context"
	"github.com/k0kubun/pp"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/protocol/esig"
	"github.com/v8platform/rac/protocol/internal/pool"
	"github.com/v8platform/rac/protocol/messages"
	"github.com/v8platform/rac/protocol/types"
	"net"
	"strconv"
	"time"

	bus "github.com/asaskevich/EventBus"
	//"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

const readTimeout = time.Second * 15
const magic = 475223888
const protocolVersion = 256
const maxResponseChunkSize = 1460

var serviceVersions = []string{"3.0", "4.0", "5.0", "6.0", "7.0", "8.0", "9.0", "10.0"}

type Client struct {
	addr  string
	laddr net.Addr

	ctx          context.Context
	stopRoutines context.CancelFunc // остановить ping, read, и подобные горутины

	agentUser     string
	agentPassword string

	pool pool.EndpointPool

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
		PoolSize:           5,
		MinIdleConns:       1,
		MaxConnAge:         time.Hour,
		IdleTimeout:        10 * time.Minute,
		IdleCheckFrequency: 20 * time.Second,
		PoolTimeout:        10 * time.Minute,
	})

	//m.serviceVersion = serviceVersions[len(serviceVersions)-1]
	m.serviceVersion = "9.0"

	return m
}

func (c *Client) initConnection(ctx context.Context, conn *pool.Conn) error {

	negotiateMessage := NewNegotiateMessage(protocolVersion, c.codec.Version())

	err := c.sendRequestMessage(conn, negotiateMessage)

	if err != nil {
		return err
	}

	err = c.sendRequestMessage(conn, &ConnectMessage{params: map[string]interface{}{
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

		message, ok := err.(*messages.EndpointMessageFailure)
		if !ok {
			return nil, err
		}
		supportedVersion := detectSupportedVersion(message)
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
		format:    0, // defaultFormat,
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

	err := c.sendRequestMessage(conn, &OpenEndpointMessage{Version: c.serviceVersion})

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

		pp.Println(answer)
		panic("unknown answer type")
	}

}

func (c *Client) closeEndpoint(ctx context.Context, conn *pool.Conn, endpoint *pool.Endpoint) error {

	pp.Println("close endpoint", endpoint.ID())
	err := c.sendRequestMessage(conn, &CloseEndpointMessage{EndpointID: endpoint.ID()})

	//_, err = conn.GetPacket(ctx)

	if err != nil {
		return err
	}

	return nil
}
func (c *Client) sendRequestMessage(conn *pool.Conn, message types.RequestMessage) error {

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

		panic(pp.Sprintln(int(packet.Type), "packet", packet))

	default:

		panic(pp.Sprintln(int(packet.Type), "packet", packet))
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

func (c *Client) getEndpoint(ctx context.Context, sig esig.ESIG) (*pool.Endpoint, error) {

	return c.pool.Get(ctx, sig)

}

func (c *Client) putEndpoint(ctx context.Context, endpoint *pool.Endpoint) {

	c.pool.Put(ctx, endpoint)

}

func (c *Client) withEndpoint(ctx context.Context, sig esig.ESIG, fn func(context.Context, *pool.Endpoint) error) error {

	cn, err := c.getEndpoint(ctx, sig)
	if err != nil {
		return err
	}

	defer c.putEndpoint(ctx, cn)

	err = fn(ctx, cn)

	return err

}

func (c *Client) sendEndpointRequest(ctx context.Context, req types.EndpointRequestMessage) (interface{}, error) {

	var value interface{}

	err := c.withEndpoint(ctx, req.Sig(), func(ctx context.Context, p *pool.Endpoint) error {

		message, err := p.SendRequest(ctx, req)

		if err != nil {
			return err
		}

		value = message.Message

		return err
	})

	return value, err

}

func (c *Client) Disconnect() error {
	// stop all routines
	c.stopRoutines()

	//err := c.conn.Close()
	//if err != nil {
	//	return errors.Wrap(err, "closing TCP connection")
	//}

	// TODO: закрыть каналы

	// возвращаем в false, потому что мы теряем конфигурацию
	// сессии, и можем ее потерять во время отключения.

	return nil
}
