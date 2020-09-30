package rac

import "time"

type Session struct {
	UUID         string    `rac:"session"`    // UUID session                          : 1fb5f037-99e8-4924-a99d-a9e687522d32
	ID           int       `rac:"session-id"` // ID session-id                       : 1
	Infobase     string    // Infobase infobase                         : aea71760-15b3-485a-9a35-506eb8a0b04a
	Connection   string    // connection                      : 8adf4514-0379-4333-a153-0b2689edc415
	Process      string    // process                          : 1af2e54f-d95a-4370-9b45-8277280cad23
	UserName     string    // user-name                        : АКузнецов
	Host         string    //host                             : Sport1
	AppId        string    //app-id                           : Designer
	Locale       string    //locale                           : ru_RU
	StartedAt    time.Time //started-at                       : 2018-04-09T14:51:31
	LastActiveAt time.Time //last-active-at                   : 2018-05-14T11:12:33
	//hibernate                        : no
	//passive-session-hibernate-time   : 1200
	//hibernate-session-terminate-time : 86400
	//blocked-by-dbms                  : 0
	//blocked-by-ls                    : 0
	//bytes-all                        : 105972550
	//bytes-last-5min                  : 0
	//calls-all                        : 119052
	//calls-last-5min                  : 0
	//dbms-bytes-all                   : 317824922
	//dbms-bytes-last-5min             : 0
	//db-proc-info                     :
	//db-proc-took                     : 0
	//db-proc-took-at                  :
	//duration-all                     : 66184
	//duration-all-dbms                : 43242
	//duration-current                 : 0
	//duration-current-dbms            : 0
	//duration-last-5min               : 0
	//duration-last-5min-dbms          : 0

}
