package serialize

import (
	"github.com/hashicorp/go-multierror"
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

type ConnectionCloser interface {
	DisconnectConnection(cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) error
}

type ConnectionSig interface {
	Sig() (cluster uuid.UUID, process uuid.UUID, connection uuid.UUID)
}

type ConnectionInfoList []ConnectionInfo

func (l ConnectionInfoList) ByID(id uuid.UUID) (ConnectionInfo, bool) {

	if id == uuid.Nil {
		return ConnectionInfo{}, false
	}

	fn := func(info ConnectionInfo) bool {
		return uuid.Equal(info.UUID, id)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return ConnectionInfo{}, false

}

func (l ConnectionInfoList) ByProcess(id uuid.UUID) (ConnectionInfo, bool) {

	if id == uuid.Nil {
		return ConnectionInfo{}, false
	}

	fn := func(info ConnectionInfo) bool {
		return uuid.Equal(info.Process, id)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return ConnectionInfo{}, false

}

func (l ConnectionInfoList) ByInfobase(id uuid.UUID) (ConnectionInfo, bool) {

	if id == uuid.Nil {
		return ConnectionInfo{}, false
	}

	fn := func(info ConnectionInfo) bool {
		return uuid.Equal(info.Infobase, id)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return ConnectionInfo{}, false

}

func (l ConnectionInfoList) Find(fn func(info ConnectionInfo) bool) (ConnectionInfo, bool) {

	val := l.filter(fn, 1)

	if len(val) == 0 {
		return ConnectionInfo{}, false
	}

	return val[0], true

}

func (l ConnectionInfoList) Filter(fn func(info ConnectionInfo) bool) ConnectionInfoList {

	return l.filter(fn, 0)

}

func (l ConnectionInfoList) Each(fn func(info ConnectionInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l ConnectionInfoList) Disconnect(closer ConnectionCloser) (n int, err error) {

	var mErr *multierror.Error

	l.Each(func(info ConnectionInfo) {

		errDisconnect := closer.DisconnectConnection(info.Sig())

		if errDisconnect != nil {
			multierror.Append(mErr, errDisconnect)
			return
		}

		n += 1

	})

	return

}

func (l ConnectionInfoList) filter(fn func(info ConnectionInfo) bool, count int) (val ConnectionInfoList) {

	n := 0

	for _, info := range l {

		if n == count {
			break
		}

		result := fn(info)

		if result {
			n += 1
			val = append(val, info)
		}

	}

	return

}

func (l *ConnectionInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls ConnectionInfoList

	for i := 0; i < count; i++ {

		info := &ConnectionInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, *info)
	}

	*l = ls
}

type ConnectionInfo struct {
	UUID        uuid.UUID `rac:"connection"` //connection     : cd16cde9-6372-4664-ac61-b0ae5cb24478
	ID          int       `rac:"conn-id"`    //conn-id        : 8714
	Host        string    //host           : srv-uk-term-09
	Process     uuid.UUID //process        : 94232f94-be78-4acd-a11e-09911bd4f4ed
	Cluster     uuid.UUID `rac:"-"` //cluster        : 94232f94-be78-4acd-a11e-09911bd4f4ed
	Infobase    uuid.UUID //infobase       : efa3672f-947a-4d84-bd58-b21997b83561
	Application string    //application    : "1CV8"
	ConnectedAt time.Time //connected-at   : 2020-10-01T07:29:57
	SessionID   int       `rac:"session-number"` //session-number : 148542
	BlockedByLs int       //blocked-by-ls  : 0
}

func (i *ConnectionInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Application, r)
	decoder.IntPtr(&i.BlockedByLs, r)
	decoder.TimePtr(&i.ConnectedAt, r)
	decoder.IntPtr(&i.ID, r)
	decoder.StringPtr(&i.Host, r)
	decoder.UuidPtr(&i.Infobase, r)
	decoder.UuidPtr(&i.Process, r)
	decoder.IntPtr(&i.SessionID, r)

}

func (info ConnectionInfo) Sig() (cluster uuid.UUID, process uuid.UUID, connection uuid.UUID) {
	return info.Cluster, info.Process, info.UUID
}
