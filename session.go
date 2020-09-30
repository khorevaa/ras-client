package rac

import "time"

type Session struct {
	UUID                          string    `rac:"session"`    // UUID session                          : 1fb5f037-99e8-4924-a99d-a9e687522d32
	ID                            int64     `rac:"session-id"` // ID session-id                       : 1
	Infobase                      string    // Infobase infobase               : aea71760-15b3-485a-9a35-506eb8a0b04a
	Connection                    string    // connection                      : 8adf4514-0379-4333-a153-0b2689edc415
	Process                       string    // process                         : 1af2e54f-d95a-4370-9b45-8277280cad23
	UserName                      string    // user-name                       : АКузнецов
	Host                          string    //host                             : Sport1
	AppId                         string    //app-id                           : Designer
	Locale                        string    //locale                           : ru_RU
	StartedAt                     time.Time //started-at                       : 2018-04-09T14:51:31
	LastActiveAt                  time.Time //last-active-at                   : 2018-05-14T11:12:33
	Hibernate                     bool      // hibernate                        : no
	PassiveSessionHibernateTime   int32     //passive-session-hibernate-time   : 1200
	HibernateDessionTerminateTime int32     //hibernate-session-terminate-time : 86400
	BlockedByDbms                 int64     //blocked-by-dbms                  : 0
	BlockedByLs                   int64     //blocked-by-ls                    : 0
	BytesAll                      int64     //bytes-all                        : 105972550
	BytesLast5min                 int64     `rac:"bytes-last-5min"` //bytes-last-5min                  : 0
	CallsAll                      int64     //calls-all                        : 119052
	CallsLast5min                 int64     `rac:"calls-last-5min"` //calls-last-5min                  : 0
	DbmsBytesAll                  int64     //dbms-bytes-all                   : 317824922
	DbmsBytesLast5min             int64     `rac:"dbms-bytes-last-5min"` //dbms-bytes-last-5min             : 0
	DbProcInfo                    string    //db-proc-info                     :
	DbProcTook                    int32     //db-proc-took                     : 0
	DbProcTookAt                  time.Time //db-proc-took-at                  :
	DurationAll                   int64     //duration-all                     : 66184
	DurationAllDbms               int64     //duration-all-dbms                : 43242
	DurationCurrent               int64     //duration-current                 : 0
	DurationCurrentDbms           int64     //duration-current-dbms            : 0
	DurationLast5Min              int64     `rac:"duration-last-5min"`      //duration-last-5min               : 0
	DurationLast5MinDbms          int64     `rac:"duration-last-5min-dbms"` //duration-last-5min-dbms          : 0

}
