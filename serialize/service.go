package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
)

type ServiceInfo struct {
	Name        string      //name      : EventLogService
	MainOnly    int         //main-only : 0
	Manager     []uuid.UUID //manager   : ad2754ad-9415-4689-9559-74dc36b11592
	Description string      `rac:"descr"` //descr     : "Сервис журналов регистрации"
	ClusterID   uuid.UUID
}

func (i *ServiceInfo) Parse(decoder Decoder, _ int, r io.Reader) {

	decoder.StringPtr(&i.Name, r)
	decoder.StringPtr(&i.Description, r)
	decoder.IntPtr(&i.MainOnly, r)

	idCount := decoder.Size(r)

	for ii := 0; ii < idCount; ii++ {
		i.Manager = append(i.Manager, decoder.Uuid(r))
	}
}
