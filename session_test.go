package rac

import (
	"log"
	"testing"
)

func TestUnmarshalSession(t *testing.T) {

	data := `session                          : 1fb5f037-99e8-4924-a99d-a9e687522d32
session-id                       : 1
infobase                         : aea71760-15b3-485a-9a35-506eb8a0b04a
connection                       : 8adf4514-0379-4333-a153-0b2689edc415
process                          : 1af2e54f-d95a-4370-9b45-8277280cad23
user-name                        : АКузнецов
host                             : Sport1
app-id                           : Designer
locale                           : ru_RU
started-at                       : 2018-04-09T14:51:31
last-active-at                   : 2018-05-14T11:12:33
hibernate                        : no
passive-session-hibernate-time   : 1200
hibernate-session-terminate-time : 86400
blocked-by-dbms                  : 0
blocked-by-ls                    : 0
bytes-all                        : 105972550
bytes-last-5min                  : 0
calls-all                        : 119052
calls-last-5min                  : 0
dbms-bytes-all                   : 317824922
dbms-bytes-last-5min             : 0
db-proc-info                     : 
db-proc-took                     : 0
db-proc-took-at                  : 
duration-all                     : 66184
duration-all-dbms                : 43242
duration-current                 : 0
duration-current-dbms            : 0
duration-last-5min               : 0
duration-last-5min-dbms          : 0
`

	var test []Session

	_ = Unmarshal([]byte(data), &test)

	log.Println(test)
	var test2 Session

	_ = Unmarshal([]byte(data), &test2)

	log.Println(test2)

}
