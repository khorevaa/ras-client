package pool

import (
	"bytes"
	"context"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol"
	"github.com/v8platform/rac/protocol/codec"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Что есть
// 1. Клиент подключения к сервесу.
// 2. На клиенте есть несколько endpoint
//   - требует авторизации на кластере
//   - требует авториазции базы ( 1 раз для каждой базы,
//     если меняешь базы или используется другиой пароль доступа
//     нужна переавторизация
// 3. Деление endpoint по:
//   - общая инфорация (сообщения в целом на кластер)
//   - информация по базам (сообщения с ключем infobase
//

// Надо организовать пул клиентов
// пул endpoint для переключения авторизации по базам
// Клиент читает сообщения для endpoint и отправляем его в ответ.

/*

Алгоритм работы с подключением

1. Создается клиент подключения (далее клиент)
2. Создается 1 проверочное соединение.
2.1. Выполняются начальные команды:
	- NewNegotiateMessage(protocolVersion, c.codec.Version()))
	- &ConnectMessage{params: map[string]interface{}{
		"connect.timeout": int64(2000) Ъ}
2.2. Соединение переходит в ожидание, при отсутствии ошибок. При ошибке клиент не создается.
3. Основной цикл работы
   Запрос на данные - Ответ пользователю
3.1. Открывает точку обмена
3.2. Авторизация на кластере -> Возможна ошибка прав/авторизации -> возврат ошибки
3.3. Если необходимо авторизация на информациооной базе -> Возможна ошибка прав/авторизации -> возврат ошибки
3.4. Выполнение запроса -> Возможна ошибка парсинга -> возврат ошибки
3.5. Ожидание ответа. Для запросов (VIOD_MESSAGE) не ошидания ответа, переход сразу к пункту 3.8
3.6. Разбор ответа. -> возможна ошибка запроса
3.7. Отправка ответа пользователю
3.8. Перевод точки обмена в ожидание. По двум критериям
	- запрос был только на данные кластера (переиспользование для аналогичных запросов)
	- была авторизация по ИБ. (переиспользование для запросов по данной базе)
	  по истечении, н-минут переход в исползование по другим базам, с повторной авторизацией

4. Цикла работы точки обмена
4.1. Открытие
4.2. Отправка собщения
	 Установка блокировки на соединение -> Запись ланных в соединение
4.2. Чтение данных из соединения -> Получение сообщения
     Снятие блокировки на соедиенение
4.3. Ожидание для срока жизни / повторение цикла с пункта 4.2.
4.4. Закрытие точки обмена
4.5. Завершение при закрытии соединения

5. Работа с открытым соединением
5.1. Блокировка использования другими точками обмена
5.2. Запись данных
5.3. Ожидание ответа -> чтение данных. При открытой точке обмена всегда приходит ответ на посланный запрос
	 Даже если он не требует явного ответа, например Авторизация на кластере или в информационной базе
5.4. Разблокировка по таймауту или при получении ответа


*/

type RequestNeedInfobaseAuth interface {
	NeedInfobaseAuth()
}

var timers = sync.Pool{
	New: func() interface{} {
		t := time.NewTimer(time.Hour)
		t.Stop()
		return t
	},
}

var (
	ErrClosed      = errors.New("protocol: pool is closed")
	ErrPoolTimeout = errors.New("protocol: endpoint pool timeout")
)

type EndpointInfo interface {
	ID() int
	Version() int
	Format() string
	ServiceID() string
	Codec() codec.Codec
}

type Endpoint struct {
	id        int
	version   int
	format    string
	serviceID string
	codec     codec.Codec

	createdAt time.Time
	usedAt    uint32 // atomic
	pooled    bool
	Inited    bool

	clusterAuth  uuid.UUID
	infobaseAuth uuid.UUID

	onRequest func(e *Endpoint, ctx context.Context, conn net.Conn, m protocol.EndpointMessage) (*protocol.EndpointMessage, error)
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

func (e *Endpoint) sendRequest(ctx context.Context, conn net.Conn, m protocol.EndpointMessage) (*protocol.EndpointMessage, error) {

	body := bytes.NewBuffer([]byte{})

	m.Format(e.codec.Encoder(), e.version, body)

	packet := NewPacket(byte(m.Type().Type()), body.Bytes())

	err := packet.Write(conn)

	//conn.Unlock()
}

func (e *Endpoint) SendRequestCtx(ctx context.Context, conn net.Conn, m protocol.EndpointMessage) (*protocol.EndpointMessage, error) {

	//conn.Lock()

	if e.onRequest != nil {

		_, err := e.onRequest(e, ctx, conn, m)

		if err != nil {
			return nil, err
		}

	}

	//conn.Unlock()
}

type Options struct {
	Opener func(context.Context) (EndpointInfo, error)
	Closer func(context.Context, EndpointInfo) error

	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

type EndpointPool struct {
	opt *Options

	dialErrorsNum uint32 // atomic

	lastDialErrorMu sync.RWMutex
	lastDialError   error

	openfunc  func(ctx context.Context) error
	OnRequest func(ctx context.Context, req RequestNeedInfobaseAuth) error

	authCache   map[uuid.UUID]struct{ user, password string }
	authCacheMu sync.Mutex

	_closed uint32 // atomic

	queue chan struct{}

	endpointsMu      sync.Mutex
	endpoints        map[string][]*Endpoint // 1. UUID - Кластера  2 UUID - информационой базы
	idleEndpoints    map[string][]*Endpoint
	idleEndpointsLen int
	poolSize         int
}

func (p *EndpointPool) SetAuthHeader(uuid uuid.UUID, user, password string) {

	p.authCacheMu.Lock()
	p.authCache[uuid] = struct {
		user, password string
	}{user, password}
	p.authCacheMu.Unlock()

}

func (p *EndpointPool) Get(ctx context.Context, cluster, infobase uuid.UUID) (*Endpoint, error) {

	if p.closed() {
		return nil, ErrClosed
	}

	err := p.waitTurn(ctx)
	if err != nil {
		return nil, err
	}

	hash := cluster.String() + "_" + infobase.String()

	for {
		p.endpointsMu.Lock()
		cn := p.popIdle(hash)
		p.endpointsMu.Unlock()

		if cn == nil {
			break
		}

		if p.isStaleConn(cn) {
			_ = p.CloseConn(cn)
			continue
		}

		//atomic.AddUint32(&p.stats.Hits, 1)
		return cn, nil
	}

	newcn, err := p.newEndpoint(ctx, true)
	if err != nil {
		p.freeTurn()
		return nil, err
	}

	return newcn, nil

}

func (p *EndpointPool) newEndpoint(c context.Context, pooled bool) (*Endpoint, error) {
	cn, err := p.openEndpoint(c, pooled)
	if err != nil {
		return nil, err
	}

	p.endpointsMu.Lock()
	p.endpoints[clearHash] = append(p.endpoints[clearHash], cn)
	if pooled {
		// If pool is full remove the cn on next Put.
		if p.poolSize >= p.opt.PoolSize {
			cn.pooled = false
		} else {
			p.poolSize++
		}
	}
	p.endpointsMu.Unlock()
	return cn, nil
}

func (p *EndpointPool) isStaleConn(cn *Conn) bool {
	if p.opt.IdleTimeout == 0 && p.opt.MaxConnAge == 0 {
		return false
	}

	now := time.Now()
	if p.opt.IdleTimeout > 0 && now.Sub(cn.UsedAt()) >= p.opt.IdleTimeout {
		return true
	}
	if p.opt.MaxConnAge > 0 && now.Sub(cn.createdAt) >= p.opt.MaxConnAge {
		return true
	}

	return false
}

func (p *EndpointPool) openEndpoint(c context.Context, pooled bool) (*Endpoint, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.opt.PoolSize) {
		return nil, p.getLastOpenError()
	}

	endpoint, err := p.opt.Opener(c)
	if err != nil {
		p.setLastOpenError(err)
		if atomic.AddUint32(&p.dialErrorsNum, 1) == uint32(p.opt.PoolSize) {
			go p.tryOpen()
		}
		return nil, err
	}

	cn := NewEndpoint(endpoint)
	cn.pooled = pooled
	return cn, nil
}

func NewEndpoint(endpoint EndpointInfo) *Endpoint {

	return &Endpoint{
		ID:      endpoint.ID(),
		Version: endpoint.Version(),
		Format:  endpoint.Format(),
	}
}

func (p *EndpointPool) tryOpen() {
	for {
		if p.closed() {
			return
		}

		endpoint, err := p.opt.Opener(context.TODO())
		if err != nil {
			p.setLastOpenError(err)
			time.Sleep(time.Second)
			continue
		}

		atomic.StoreUint32(&p.dialErrorsNum, 0)
		_ = p.opt.Closer(context.TODO(), endpoint)
		return
	}
}

func (p *EndpointPool) popIdle(hash string) *Endpoint {

	if p.idleEndpointsLen == 0 {
		return nil
	}

	idleendpoints, ok := p.idleEndpoints[hash]
	if !ok {
		return p.popIdle(clearHash)
	}

	if len(idleendpoints) == 0 {
		return nil
	}

	idx := len(idleendpoints) - 1
	cn := idleendpoints[idx]

	p.idleEndpoints[hash] = idleendpoints[:idx]
	p.idleEndpointsLen--
	p.checkMinIdleEndpoints()
	return cn
}

var clearHash = uuid.UUID{}.String() + "_" + uuid.UUID{}.String()

func (p *EndpointPool) checkMinIdleEndpoints() {
	if p.opt.MinIdleConns == 0 {
		return
	}
	for p.poolSize < p.opt.PoolSize && p.idleEndpointsLen < p.opt.MinIdleConns {
		p.poolSize++
		p.idleEndpointsLen++
		go func() {

			err := p.addIdleEndpoint()
			if err != nil {
				p.endpointsMu.Lock()
				p.poolSize--
				p.idleEndpointsLen--
				p.endpointsMu.Unlock()
			}
		}()
	}
}

func (p *EndpointPool) setLastOpenError(err error) {
	p.lastDialErrorMu.Lock()
	p.lastDialError = err
	p.lastDialErrorMu.Unlock()
}

func (p *EndpointPool) getLastOpenError() error {
	p.lastDialErrorMu.RLock()
	err := p.lastDialError
	p.lastDialErrorMu.RUnlock()
	return err
}

func (p *EndpointPool) addIdleEndpoint() error {

	newEndpoint := &Endpoint{}

	p.endpointsMu.Lock()
	p.endpoints[clearHash] = append(p.endpoints[clearHash], newEndpoint)
	p.idleEndpoints[clearHash] = append(p.endpoints[clearHash], newEndpoint)
	p.endpointsMu.Unlock()

	return nil
}

func (p *EndpointPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}

func (p *EndpointPool) waitTurn(c context.Context) error {
	select {
	case <-c.Done():
		return c.Err()
	default:
	}

	select {
	case p.queue <- struct{}{}:
		return nil
	default:
	}

	timer := timers.Get().(*time.Timer)
	timer.Reset(p.opt.PoolTimeout)

	select {
	case <-c.Done():
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return c.Err()
	case p.queue <- struct{}{}:
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return nil
	case <-timer.C:
		timers.Put(timer)
		//atomic.AddUint32(&p.stats.Timeouts, 1)
		return ErrPoolTimeout
	}
}
func (p *EndpointPool) getTurn() {
	p.queue <- struct{}{}
}

func (p *EndpointPool) freeTurn() {
	<-p.queue
}

func (p *EndpointPool) CloseConn(cn *Endpoint) error {
	p.removeConnWithLock(cn)
	return p.closeConn(cn)
}

func (p *EndpointPool) removeConnWithLock(cn *Endpoint) {
	p.endpointsMu.Lock()
	p.removeConn(cn)
	p.endpointsMu.Unlock()
}

func (p *EndpointPool) removeConn(cn *Endpoint) {

	for _, endpoints := range p.endpoints {

		for i, c := range endpoints {
			if c == cn {
				endpoints = append(endpoints[:i], endpoints[i+1:]...)
				if cn.pooled {
					p.poolSize--
					p.checkMinIdleEndpoints()
				}
				return
			}
		}

	}

}

func (p *EndpointPool) closeConn(cn *Endpoint) error {

	return p.opt.Closer(context.Background(), cn)

}

type exchange struct {
	ID int16

	authCluster  uuid.UUID
	authInfobase uuid.UUID
}
