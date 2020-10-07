package protocol

import (
	"context"
	"net"
	"strconv"
	"time"

	bus "github.com/asaskevich/EventBus"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
)

const readTimeout = time.Second * 15
const magic = 475223888

type RASConn struct {
	addr         net.Addr
	laddr        net.Addr
	conn         net.Conn
	packetConn   net.PacketConn
	stopRoutines context.CancelFunc // остановить ping, read, и подобные горутины

	// ключ авторизации. изменять можно только через setAuthKey
	authKey []byte

	// хеш ключа авторизации. изменять можно только через setAuthKey
	authKeyHash []byte

	// каналы, которые ожидают ответа rpc. ответ записывается в канал и удаляется
	ackRespond   chan map[MessageType]RespondMessage
	ackWait      chan RespondMessage
	ackWaitError chan error

	Endpoint       *Endpoint
	EndpointClosed bool

	// шина соо
	//бщений, используется для разных нотификаций, описанных в константах нотификации
	bus bus.Bus
}

func NewRASConn(addr string) (*RASConn, error) {
	m := new(RASConn)
	m.addr, _ = net.ResolveTCPAddr("tcp", addr)
	m.ackRespond = make(chan map[MessageType]RespondMessage)
	m.resetAck()
	return m, nil
}

func (m *RASConn) CreateConnection() error {
	// connect
	tcpAddr, err := net.ResolveTCPAddr("tcp", m.addr.String())
	if err != nil {
		return errors.Wrap(err, "resolving tcp")
	}

	m.conn, err = net.Dial("tcp", tcpAddr.String())
	if err != nil {
		return errors.Wrap(err, "dialing tcp")
	}
	m.laddr = m.conn.LocalAddr()

	ctx, cancelfunc := context.WithCancel(context.Background())
	m.stopRoutines = cancelfunc

	// start reading responses from the server
	m.startReadingResponses(ctx)

	// start keepalive pinging
	m.startPinging(ctx)

	err = m.VoidRequest(NewNegotiateMessage(256, 256))

	if err != nil {
		return err
	}
	ack := &ConnectMessageAck{}

	_, err = m.SendRequest(&ConnectMessage{params: map[string]interface{}{
		"connect.timeout": int64(2000),
	}}, ack)

	pp.Println(ack)
	return err
}

func (m *RASConn) resetAck() {
	m.ackWait = make(chan RespondMessage, 1)
	m.ackWaitError = make(chan error, 1)
}

// waitAck добавляет в список id сообщения, которому нужно подтверждение
// возвращает true, если ранее этого id не было
func (m *RASConn) waitAck(resp []RespondMessage) (RespondMessage, error) {

	if resp == nil || len(resp) == 0 {
		return &nullRespondMessage{}, nil
	}

	ackRespond := make(map[MessageType]RespondMessage)

	for _, message := range resp {
		ackRespond[message.Type()] = message
	}

	m.ackRespond <- ackRespond

	select {

	case err := <-m.ackWaitError:

		if err != nil {
			return nil, err
		}

	case r := <-m.ackWait:

		pp.Println(r)
		return r, nil
	}

	return nil, nil
}

func (m *RASConn) sendPacket(request []byte, resp []RespondMessage) (RespondMessage, error) {

	pp.Println("writing message", request)
	_, err := m.conn.Write(request)
	if err != nil {
		return nil, errors.Wrap(err, "sending request")
	}

	return m.waitAck(resp)
}

func (m *RASConn) VoidRequest(req RequestMessage) (err error) {

	requestData := formatRequestMessage(req)
	_, err = m.sendPacket(requestData, nil)

	return err
}

func (m *RASConn) SendRequest(req RequestMessage, resp ...RespondMessage) (RespondMessage, error) {

	requestData := formatRequestMessage(req)

	r, err := m.sendPacket(requestData, resp)

	return r, err
}

func (m *RASConn) SendEndpointRequest(req RequestMessage, resp ...RespondMessage) (RespondMessage, error) {

	if m.Endpoint == nil || m.EndpointClosed {
		return nullRespondMessage{}, errors.New("endpoint is in active")
	}

	endpointMessage := &EndpointMessage{}
	endpointMessage.addResponse(resp...)

	waitResp := []RespondMessage{endpointMessage, &EndpointFeature{}}

	body := m.formatEndpointRequestMessage(req, VOID_MESSAGE_KIND)

	requestData := formatMessageType(ENDPOINT_MESSAGE, body)

	r, err := m.sendPacket(requestData, waitResp)

	return r, err
}

func formatRequestMessage(req RequestMessage) []byte {

	switch req.Type() {

	case NEGOTIATE:
		enc := NewEncoder()
		req.Format(enc)
		return enc.Bytes()

	default:

		enc := NewEncoder()
		req.Format(enc)
		body := enc.Bytes()
		return formatMessageType(req.Type(), body)

	}
}

func (m *RASConn) formatEndpointRequestMessage(req RequestMessage, kind EndpointMessageKind) []byte {

	enc := NewEncoder()
	req.Format(enc)
	body := enc.Bytes()

	enc = NewEncoder()

	enc.encodeEndpointId(m.Endpoint.Id)
	enc.encodeShort(m.Endpoint.Fornat)
	enc.encodeByte(byte(kind))

	enc.encodeType(req.Type()) // МАГИЯ без этого байта требует авторизации на центральном кластере

	header := enc.Bytes()

	buf := make([]byte, len(header)+len(body))
	copy(buf, header)
	copy(buf[len(header):], body)

	return buf

}

func formatMessageType(mType MessageType, body []byte) []byte {

	enc := NewEncoder()
	enc.encodeType(mType)
	enc.encodeSize(len(body))
	header := enc.Bytes()

	buf := make([]byte, len(header)+len(body))
	copy(buf, header)
	copy(buf[len(header):], body)
	return buf
}

func (m *RASConn) Disconnect() error {
	// stop all routines
	m.stopRoutines()

	err := m.conn.Close()
	if err != nil {
		return errors.Wrap(err, "closing TCP connection")
	}

	// TODO: закрыть каналы

	// возвращаем в false, потому что мы теряем конфигурацию
	// сессии, и можем ее потерять во время отключения.

	return nil
}

// startPinging пингует сервер что все хорошо, клиент в сети
// нужно просто запустить
func (m *RASConn) startPinging(ctx context.Context) {
	ticker := time.Tick(60 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				//_, err := m.Ping(0xCADACADA)
				//dry.PanicIfErr(err)
			}
		}
	}()
}

func (m *RASConn) startReadingResponses(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ack := <-m.ackRespond:

				rawResp, err := m.readFromConn(ctx)

				if err != nil {
					m.ackWaitError <- err
					return
				}

				parser, ok := ack[rawResp.Type()]
				if !ok {
					m.ackWaitError <- errors.New("got unsupporsed type")
					return
				}

				err = parser.Parse(rawResp.Data())
				if err != nil {
					m.ackWaitError <- err
					return
				}

				m.ackWait <- parser

			}
		}
	}()
}

func (m *RASConn) readFromConn(ctx context.Context) (reso rawRespond, err error) {

	header := make([]byte, 256)
	n, err := m.conn.Read(header)
	dry.PanicIfErr(err)
	header = header[:n]

	dec := NewDecoder(header)
	messageType := dec.decodeType()

	pp.Println("messageType", messageType)
	size := dec.decodeSize()
	pp.Println("count", size)

	data := make([]byte, size)
	reader := dry.NewCancelableReader(ctx, m.conn)
	n, err = reader.Read(data)
	dry.PanicIfErr(err)
	dry.PanicIf(n != int(size), "expected read "+strconv.Itoa(int(size))+" bytes, got "+strconv.Itoa(n))

	resp := rawRespond{
		ConnectionMessageType(messageType),
		size,
		data,
	}

	return resp, nil // TODO Переделать на возврат сообщения
}

type rawRespond struct {
	t       MessageType
	size    int
	payload []byte
}

func (r rawRespond) Type() MessageType {
	return r.t
}

func (r rawRespond) Len() int {
	return r.size
}

func (r rawRespond) Size() int {
	return r.size
}

func (r rawRespond) Data() []byte {
	return r.payload
}

type RequestMessage interface {
	Type() MessageType
	Format(enc *encoder)
}

type RespondMessage interface {
	Parse(body []byte) error
	Type() MessageType
}

type nullRespondMessage struct {
}

func (_ nullRespondMessage) Parse(_ []byte) error {
	return nil
}

func (_ nullRespondMessage) Type() MessageType {
	return NULL_TYPE
}

type Endpoint struct {
	Id        int
	ServiceID string
	Version   string
	Fornat    int
}
