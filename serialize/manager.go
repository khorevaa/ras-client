package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
)

type ManagerInfo struct {
	UUID        uuid.UUID `rac:"manager"` //manager : 0e588a25-8354-4344-b935-53442312aa30
	PID         string    //pid     : 3388
	Using       string    //using   : normal
	Host        string    //host    : Sport1
	MainManager int
	Port        int16  //port    : 1541
	Description string `rac:"descr"` //descr   : "Главный менеджер кластера"

}

func (i *ManagerInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.Description, r)
	decoder.StringPtr(&i.Host, r)
	decoder.IntPtr(&i.MainManager, r)
	decoder.ShortPtr(&i.Port, r) // expirationTimeout
	decoder.StringPtr(&i.PID, r)

}
