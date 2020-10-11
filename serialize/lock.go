package serialize

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type LockInfo struct {
	Connection  uuid.UUID //connection : 00000000-0000-0000-0000-000000000000
	Session     uuid.UUID //session    : 8b8a0817-4cb1-4e13-9a8f-472dde1a3b47
	Object      uuid.UUID //object     : 00000000-0000-0000-0000-000000000000
	Locked      time.Time //locked     : 2020-10-01T08:30:00
	Description string    `rac:"descr"` //descr      : "БД(сеанс ,УППБоеваяБаза,разделяемая)"

}
