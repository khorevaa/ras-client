package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

type ConnectionInfoList []ConnectionInfo

func (l ConnectionInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)

	for i := 0; i < count; i++ {

		info := &ConnectionInfo{}
		info.Parse(decoder, version, r)

		l = append(l, *info)
	}
}

type ConnectionInfo struct {
	UUID        uuid.UUID `rac:"connection"` //connection     : cd16cde9-6372-4664-ac61-b0ae5cb24478
	ID          int       `rac:"conn-id"`    //conn-id        : 8714
	Host        string    //host           : srv-uk-term-09
	Process     uuid.UUID //process        : 94232f94-be78-4acd-a11e-09911bd4f4ed
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
