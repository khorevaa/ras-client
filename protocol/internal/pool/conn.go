package pool

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const maxResponseChunkSize = 1460

// IOError is the data type for errors occurring in case of failure.
type IOError struct {
	Connection net.Conn
	Error      error
}

type Conn struct {
	connMU *sync.Mutex

	_locked uint32
	netConn net.Conn
	onError func(err IOError)

	endpoints []*Endpoint

	createdAt time.Time
	usedAt    uint32 // atomic
	pooled    bool
	Inited    bool
}

func NewConn(netConn net.Conn) *Conn {

	cn := &Conn{
		createdAt: time.Now(),
	}
	cn.SetNetConn(netConn)
	cn.SetUsedAt(time.Now())

	return cn
}

func (cn *Conn) SendPacket(packet *Packet) error {

	err := packet.Write(cn.netConn)
	cn.SetUsedAt(time.Now())
	return err
}

func (cn *Conn) GetPacket(ctx context.Context) (packet *Packet, err error) {

	return cn.readContext(ctx)
}

func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadUint32(&cn.usedAt)
	return time.Unix(int64(unix), 0)
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&cn.usedAt, uint32(tm.Unix()))
}

func (cn *Conn) RemoteAddr() net.Addr {
	return cn.netConn.RemoteAddr()
}

func (cn *Conn) SetNetConn(netConn net.Conn) {
	cn.netConn = netConn
}

func (cn *Conn) Close() error {
	return cn.netConn.Close()
}

func (conn *Conn) lock() {

	conn.connMU.Lock()
	atomic.StoreUint32(&conn._locked, 1)
}

func (conn *Conn) unlock() {

	atomic.StoreUint32(&conn._locked, 0)
	conn.connMU.Unlock()

}

func (conn *Conn) Locked() bool {
	return atomic.LoadUint32(&conn._locked) == 1
}

func (conn *Conn) readContext(ctx context.Context) (*Packet, error) {

	recvDone := make(chan *Packet)
	errChan := make(chan error)

	go conn.readPacket(recvDone, errChan)

	// setup the cancellation to abort reads in process
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
			// Close() can be used if this isn't necessarily a TCP connection
		case err := <-errChan:
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				go conn.readPacket(recvDone, errChan)
				continue
			}
			return nil, err
		case packet := <-recvDone:
			return packet, nil
		}
	}

}

func (conn *Conn) readPacket(recvDone chan *Packet, errChan chan error) {

	conn.lock()
	defer conn.unlock()

	err := conn.netConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		errChan <- err
	}

	typeBuffer := make([]byte, 1)

	_, err = conn.netConn.Read(typeBuffer)

	if err != nil {
		conn.onError(IOError{conn.netConn, err})
		errChan <- err
	}

	size, err := decodeSize(conn.netConn)

	if err != nil {
		conn.onError(IOError{conn.netConn, err})
		errChan <- err
	}

	data := make([]byte, size)
	readLength := 0
	n := 0

	for readLength < len(data) {
		n, err = conn.netConn.Read(data[readLength:])
		readLength += n

		if err != nil {
			conn.onError(IOError{conn.netConn, err})
			errChan <- err
		}
	}

	recvDone <- NewPacket(typeBuffer[0], data)

}

func decodeSize(r io.Reader) (int, error) {
	ff := 0xFFFFFF80
	b1, err := readByte(r)

	if err != nil {
		return 0, err
	}
	cur := int(b1 & 0xFF)
	size := cur & 0x7F
	for shift := 7; (cur & ff) != 0x0; {

		b1, err = readByte(r)

		if err != nil {
			return 0, err
		}

		cur = int(b1 & 0xFF)
		size += (cur & 0x7F) << shift
		shift += 7

		return size, nil
	}

	return size, nil
}

func readByte(r io.Reader) (byte, error) {

	byteBuffer := make([]byte, 1)
	_, err := r.Read(byteBuffer)

	return byteBuffer[0], err
}
