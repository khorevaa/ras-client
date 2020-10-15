package pool

import (
	"context"
	"errors"
	"github.com/v8platform/rac/protocol/messages"
	"github.com/v8platform/rac/protocol/types"

	uuid "github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/esig"

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

type EndpointPooler interface {
	NewEndpoint(ctx context.Context) (*Endpoint, error)
	CloseEndpoint(endpoint *Endpoint) error

	Get(ctx context.Context, sig esig.ESIG) (*Endpoint, error)
	Put(context.Context, *Endpoint)
	Remove(context.Context, *Endpoint, error)

	Len() int
	IdleLen() int

	Close() error
}

var timers = sync.Pool{
	New: func() interface{} {
		t := time.NewTimer(time.Hour)
		t.Stop()
		return t
	},
}

var (
	ErrClosed         = errors.New("protocol: pool is closed")
	ErrUnknownMessage = errors.New("protocol: unknown message packet")
	ErrPoolTimeout    = errors.New("protocol: endpoint pool timeout")
)

func (p *EndpointPool) SetAuthHeader(uuid uuid.UUID, user, password string) {

	//p.authCacheMu.Lock()
	//p.authCache[uuid] = struct {
	//	user, password string
	//}{user, password}
	//p.authCacheMu.Unlock()

}

func (p *EndpointPool) openEndpoint(ctx context.Context, conn *Conn) (*Endpoint, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	if !conn.Inited {
		err := p.opt.InitConnection(ctx, conn)

		if err != nil {
			return nil, err
		}
	}

	openAck, err := p.opt.OpenEndpoint(ctx, conn)
	if err != nil {
		return nil, err
	}

	endpoint := NewEndpoint(openAck)
	endpoint.Inited = true
	conn.endpoints = append(conn.endpoints, endpoint)

	return endpoint, nil
}

type EndpointPool struct {
	opt *Options

	dialErrorsNum uint32 // atomic

	_closed uint32 // atomic

	lastDialErrorMu sync.RWMutex
	lastDialError   error

	queue chan struct{}

	connsMu   sync.Mutex
	conns     []*Conn
	idleConns IdleConns

	poolSize     int
	idleConnsLen int
}

func (p *EndpointPool) NewEndpoint(ctx context.Context) (*Endpoint, error) {

	if p.closed() {
		return nil, ErrClosed
	}

	err := p.waitTurn(ctx)
	if err != nil {
		return nil, err
	}

	for {
		p.connsMu.Lock()
		endpoint := p.popIdle(esig.ESIG{})
		p.connsMu.Unlock()

		if endpoint == nil {
			break
		}

		if p.isStaleConn(endpoint.conn) {
			_ = p.CloseConn(endpoint.conn)
			continue
		}

		if !endpoint.Inited {
			endpoint, err = p.openEndpoint(ctx, endpoint.conn)

			if err != nil {
				return nil, err
			}

		}

		return endpoint, nil
	}

	newcn, err := p.newConn(ctx, true)
	if err != nil {
		p.freeTurn()
		return nil, err
	}

	endpoint, err := p.openEndpoint(ctx, newcn)

	return endpoint, err

}

func (p *EndpointPool) CloseEndpoint(endpoint *Endpoint) error {
	panic("implement me")
}

var _ EndpointPooler = (*EndpointPool)(nil)

func NewEndpointPool(opt *Options) *EndpointPool {
	p := &EndpointPool{
		opt: opt,

		queue:     make(chan struct{}, opt.PoolSize),
		conns:     make([]*Conn, 0, opt.PoolSize),
		idleConns: make([]*Conn, 0, opt.PoolSize),
	}

	p.connsMu.Lock()
	p.checkMinIdleConns()
	p.connsMu.Unlock()

	if opt.IdleTimeout > 0 && opt.IdleCheckFrequency > 0 {
		go p.reaper(opt.IdleCheckFrequency)
	}

	return p
}

// Get returns existed connection from the pool or creates a new one.
func (p *EndpointPool) onRequest(ctx context.Context, endpoint *Endpoint, req types.EndpointRequestMessage) error {

	if esig.IsNul(req.Sig()) {
		return nil
	}

	if esig.Equal(endpoint.sig, req.Sig()) {
		return nil
	}

	if !esig.HighEqual(endpoint.sig, req.Sig()) {

		authMessage := endpoint.newEndpointMessage(messages.ClusterAuthenticateRequest{
			ClusterID: req.Sig().High(),
			User:      "",
			Password:  "",
		})

		_, err := endpoint.sendRequest(ctx, authMessage)

		return err
	}

	if uuid.Equal(req.Sig().Low(), uuid.Nil) {
		return nil
	}

	authMessage := endpoint.newEndpointMessage(messages.AuthenticateInfobaseRequest{
		ClusterID: req.Sig().High(),
		User:      "",
		Password:  "",
	})

	_, err := endpoint.sendRequest(ctx, authMessage)

	return err
}

// Get returns existed connection from the pool or creates a new one.
func (p *EndpointPool) Get(ctx context.Context, sig esig.ESIG) (*Endpoint, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	err := p.waitTurn(ctx)
	if err != nil {
		return nil, err
	}

	for {
		p.connsMu.Lock()
		endpoint := p.popIdle(sig)
		p.connsMu.Unlock()

		if endpoint == nil {
			break
		}

		if p.isStaleConn(endpoint.conn) {
			_ = p.CloseConn(endpoint.conn)
			continue
		}

		if !endpoint.Inited {
			endpoint, err = p.openEndpoint(ctx, endpoint.conn)

			if err != nil {
				return nil, err
			}

		}

		return endpoint, nil
	}

	newcn, err := p.newConn(ctx, true)
	if err != nil {
		p.freeTurn()
		return nil, err
	}

	endpoint, err := p.openEndpoint(ctx, newcn)

	return endpoint, err
}

func (p *EndpointPool) checkMinIdleConns() {
	if p.opt.MinIdleConns == 0 {
		return
	}
	for p.poolSize < p.opt.PoolSize && p.idleConnsLen < p.opt.MinIdleConns {
		p.poolSize++
		p.idleConnsLen++
		go func() {
			err := p.addIdleConn()
			if err != nil {
				p.connsMu.Lock()
				p.poolSize--
				p.idleConnsLen--
				p.connsMu.Unlock()
			}
		}()
	}
}

func (p *EndpointPool) addIdleConn() error {
	cn, err := p.dialConn(context.TODO(), true)
	if err != nil {
		return err
	}

	p.connsMu.Lock()
	p.conns = append(p.conns, cn)
	p.idleConns = append(p.idleConns, cn)
	p.connsMu.Unlock()
	return nil
}

func (p *EndpointPool) newConn(c context.Context, pooled bool) (*Conn, error) {
	cn, err := p.dialConn(c, pooled)
	if err != nil {
		return nil, err
	}

	p.connsMu.Lock()
	p.conns = append(p.conns, cn)
	if pooled {
		// If pool is full remove the cn on next Put.
		if p.poolSize >= p.opt.PoolSize {
			cn.pooled = false
		} else {
			p.poolSize++
		}
	}
	p.connsMu.Unlock()
	return cn, nil
}

func (p *EndpointPool) dialConn(c context.Context, pooled bool) (*Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.opt.PoolSize) {
		return nil, p.getLastDialError()
	}

	netConn, err := p.opt.Dialer(c)
	if err != nil {
		p.setLastDialError(err)
		if atomic.AddUint32(&p.dialErrorsNum, 1) == uint32(p.opt.PoolSize) {
			go p.tryDial()
		}
		return nil, err
	}

	cn := NewConn(netConn)
	cn.pooled = pooled
	return cn, nil
}

func (p *EndpointPool) tryDial() {
	for {
		if p.closed() {
			return
		}

		conn, err := p.opt.Dialer(context.TODO())
		if err != nil {
			p.setLastDialError(err)
			time.Sleep(time.Second)
			continue
		}

		atomic.StoreUint32(&p.dialErrorsNum, 0)
		_ = conn.Close()
		return
	}
}

func (p *EndpointPool) setLastDialError(err error) {
	p.lastDialErrorMu.Lock()
	p.lastDialError = err
	p.lastDialErrorMu.Unlock()
}

func (p *EndpointPool) getLastDialError() error {
	p.lastDialErrorMu.RLock()
	err := p.lastDialError
	p.lastDialErrorMu.RUnlock()
	return err
}

func (p *EndpointPool) getTurn() {
	p.queue <- struct{}{}
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

func (p *EndpointPool) freeTurn() {
	<-p.queue
}

func (p *EndpointPool) popIdle(sig esig.ESIG) *Endpoint {
	if len(p.idleConns) == 0 {
		return nil
	}

	endpoint := p.idleConns.Pop(sig)

	if endpoint == nil {
		return nil
	}

	p.idleConnsLen--
	p.checkMinIdleConns()
	return endpoint
}

func (p *EndpointPool) Put(ctx context.Context, cn *Endpoint) {
	if !cn.pooled {
		p.Remove(ctx, cn, nil)
		return
	}

	p.connsMu.Lock()
	p.idleConns = append(p.idleConns, cn.conn)
	p.idleConnsLen++
	p.connsMu.Unlock()
	p.freeTurn()
}

func (p *EndpointPool) Remove(ctx context.Context, cn *Endpoint, reason error) {
	p.removeConnWithLock(cn.conn)
	p.freeTurn()
	_ = p.closeConn(cn.conn)
}

func (p *EndpointPool) CloseConn(cn *Conn) error {
	p.removeConnWithLock(cn)
	return p.closeConn(cn)
}

func (p *EndpointPool) removeConnWithLock(cn *Conn) {
	p.connsMu.Lock()
	p.removeConn(cn)
	p.connsMu.Unlock()
}

func (p *EndpointPool) removeConn(cn *Conn) {
	for i, c := range p.conns {
		if c == cn {
			p.conns = append(p.conns[:i], p.conns[i+1:]...)
			if cn.pooled {
				p.poolSize--
				p.checkMinIdleConns()
			}
			return
		}
	}
}

func (p *EndpointPool) closeConn(cn *Conn) error {
	if p.opt.OnClose != nil {
		_ = p.opt.OnClose(cn)
	}
	return cn.Close()
}

// Len returns total number of connections.
func (p *EndpointPool) Len() int {
	p.connsMu.Lock()
	n := len(p.conns)
	p.connsMu.Unlock()
	return n
}

// IdleLen returns number of idle connections.
func (p *EndpointPool) IdleLen() int {
	p.connsMu.Lock()
	n := p.idleConnsLen
	p.connsMu.Unlock()
	return n
}

func (p *EndpointPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}

func (p *EndpointPool) Close() error {
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
		return ErrClosed
	}

	var firstErr error
	p.connsMu.Lock()
	for _, cn := range p.conns {
		if err := p.closeConn(cn); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.conns = nil
	p.poolSize = 0
	p.idleConns = nil
	p.idleConnsLen = 0
	p.connsMu.Unlock()

	return firstErr
}

func (p *EndpointPool) reaper(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for range ticker.C {
		if p.closed() {
			break
		}
		_, err := p.ReapStaleConns()
		if err != nil {
			continue
		}
	}
}

func (p *EndpointPool) ReapStaleConns() (int, error) {
	var n int
	for {
		p.getTurn()

		p.connsMu.Lock()
		cn := p.reapStaleConn()
		p.connsMu.Unlock()

		p.freeTurn()

		if cn != nil {
			_ = p.closeConn(cn)
			n++
		} else {
			break
		}
	}
	return n, nil
}

func (p *EndpointPool) reapStaleConn() *Conn {
	if len(p.idleConns) == 0 {
		return nil
	}

	cn := p.idleConns[0]
	if !p.isStaleConn(cn) {
		return nil
	}

	p.idleConns = append(p.idleConns[:0], p.idleConns[1:]...)
	p.idleConnsLen--
	p.removeConn(cn)

	return cn
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
