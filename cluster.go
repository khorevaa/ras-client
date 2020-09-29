package rac

type Cluster struct {

	// ID cluster                    : 6d6958e1-a96c-4999-a995-698a0298161e
	ID string `rac:"cluster"`
	// Host                          : Sport2
	Host string
	// Port                          : 1541
	Port int
	// Name                          : "Новый кластер"
	Name string
	// ExpirationTimeout expiration-timeout            : 0
	ExpirationTimeout int
	// LifetimeLimit lifetime-limit                : 0
	LifetimeLimit int
	// MaxMemorySize max-memory-size               : 0
	MaxMemorySize int
	// MaxMemoryTimeLimit max-memory-time-limit         : 0
	MaxMemoryTimeLimit int
	// SecurityLevel security-level                : 0
	SecurityLevel int
	// SessionFaultToleranceLevel session-fault-tolerance-level : 0
	SessionFaultToleranceLevel int
	// LoadBalancingMode load-balancing-mode           : performance
	LoadBalancingMode string
	// ErrorsCountThreshold errors-count-threshold        : 0
	ErrorsCountThreshold int
	// KillProblemProcesses kill-problem-processes        : 0
	KillProblemProcesses int
}
