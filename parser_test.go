package rac

import (
	"log"
	"testing"
)

func TestUnmarshal(t *testing.T) {

	var example = `cluster                       : 6d6958e1-a96c-4999-a995-698a0298161e
host                          : Sport2
port                          : 1541
name                          : "Новый кластер"
expiration-timeout            : 0
lifetime-limit                : 0
max-memory-size               : 0
max-memory-time-limit         : 0
security-level                : 0
session-fault-tolerance-level : 0
load-balancing-mode           : performance
errors-count-threshold        : 0
kill-problem-processes        : 0
`

	var test []ClusterInfo

	_ = Unmarshal([]byte(example), &test)

	var test2 ClusterInfo

	_ = Unmarshal([]byte(example), &test2)

	log.Println(test)

}
