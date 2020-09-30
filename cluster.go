package rac

type Cluster struct {
	ID                         string `rac:"cluster"` // ID cluster                    : 6d6958e1-a96c-4999-a995-698a0298161e
	Host                       string // Host                          : Sport2
	Port                       int    // Port                          : 1541
	Name                       string // Name                          : "Новый кластер"
	ExpirationTimeout          int    // ExpirationTimeout expiration-timeout            : 0
	LifetimeLimit              int    // LifetimeLimit lifetime-limit                : 0
	MaxMemorySize              int    // MaxMemorySize max-memory-size               : 0
	MaxMemoryTimeLimit         int    // MaxMemoryTimeLimit max-memory-time-limit         : 0
	SecurityLevel              int    // SecurityLevel security-level                : 0
	SessionFaultToleranceLevel int    // SessionFaultToleranceLevel session-fault-tolerance-level : 0
	LoadBalancingMode          string // LoadBalancingMode load-balancing-mode           : performance
	ErrorsCountThreshold       int    // ErrorsCountThreshold errors-count-threshold        : 0
	KillProblemProcesses       int    // KillProblemProcesses kill-problem-processes        : 0

}
