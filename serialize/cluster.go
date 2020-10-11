package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
)

type ClusterInfo struct {
	UUID                       uuid.UUID `rac:"cluster"` // UUID cluster                    : 6d6958e1-a96c-4999-a995-698a0298161e
	Host                       string    // Host                          : Sport2
	Port                       int16     // Port                          : 1541
	Name                       string    // Name                          : "Новый кластер"
	ExpirationTimeout          int       // ExpirationTimeout expiration-timeout            : 0
	LifetimeLimit              int       // LifetimeLimit lifetime-limit                : 0
	MaxMemorySize              int       // MaxMemorySize max-memory-count               : 0
	MaxMemoryTimeLimit         int       // MaxMemoryTimeLimit max-memory-time-limit         : 0
	SecurityLevel              int       // SecurityLevel security-level                : 0
	SessionFaultToleranceLevel int       // SessionFaultToleranceLevel session-fault-tolerance-level : 0
	LoadBalancingMode          int       // LoadBalancingMode load-balancing-mode           : performance
	ErrorsCountThreshold       int       // ErrorsCountThreshold errors-count-threshold        : 0
	KillProblemProcesses       bool      // KillProblemProcesses kill-problem-processes        : 0
	KillByMemoryWithDump       bool      // KillByMemoryWithDump kill-by-memory-with-dump      : 0
	LifeTimeLimit              int
}

func (i *ClusterInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.IntPtr(&i.ExpirationTimeout, r) // expirationTimeout
	decoder.StringPtr(&i.Host, r)
	decoder.IntPtr(&i.LifeTimeLimit, r)
	decoder.ShortPtr(&i.Port, r)
	decoder.IntPtr(&i.MaxMemorySize, r)
	decoder.IntPtr(&i.MaxMemoryTimeLimit, r)
	decoder.StringPtr(&i.Name, r)
	decoder.IntPtr(&i.SecurityLevel, r)
	decoder.IntPtr(&i.SessionFaultToleranceLevel, r)
	decoder.IntPtr(&i.LoadBalancingMode, r)
	decoder.IntPtr(&i.ErrorsCountThreshold, r)
	decoder.BoolPtr(&i.KillProblemProcesses, r)

	if version > 8 {
		decoder.BoolPtr(&i.KillByMemoryWithDump, r)
	}

}
