package rac

import (
	"log"
	"testing"
)

func TestUnmarshalProcess(t *testing.T) {

	data := `process              : 3ea9968d-159c-4b5f-9bdc-22b8ead96f74
host                 : Sport1
port                 : 1564
pid                  : 5428
is-enable            : yes
running              : yes
started-at           : 2018-03-29T11:16:02
use                  : used
available-perfomance : 100
capacity             : 1000
connections          : 7
memory-size          : 1518604
memory-excess-time   : 0
selection-size       : 61341
avg-back-call-time   : 0.000
avg-call-time        : 0.483
avg-db-call-time     : 0.124
avg-lock-call-time   : 0.000
avg-server-call-time : -0.265
avg-threads          : 0.281
`

	var test []ProcessInfo

	_ = Unmarshal([]byte(data), &test)

	log.Println(test)
	var test2 ProcessInfo

	_ = Unmarshal([]byte(data), &test2)

	log.Println(test2)

}
