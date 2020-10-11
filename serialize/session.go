package serialize

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type SessionInfo struct {
	UUID                          uuid.UUID `rac:"session"`    // UUID session                          : 1fb5f037-99e8-4924-a99d-a9e687522d32
	ID                            int64     `rac:"session-id"` // ID session-id                       : 1
	Infobase                      uuid.UUID // Infobase infobase               : aea71760-15b3-485a-9a35-506eb8a0b04a
	Connection                    uuid.UUID // connection                      : 8adf4514-0379-4333-a153-0b2689edc415
	Process                       uuid.UUID // process                         : 1af2e54f-d95a-4370-9b45-8277280cad23
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
	MemoryCurrent                 int64     //memory-current                   : 0
	MemoryLast5min                int64     //memory-last-5min                 : 416379
	MemoryTotal                   int64     //memory-total                     : 23178863
	ReadCurrent                   int64     //read-current                     : 0
	ReadLast5min                  int64     //read-last-5min                   : 0
	ReadTotal                     int64     //read-total                       : 156162
	WriteCurrent                  int64     //write-current                    : 0
	WriteLast5min                 int64     ///write-last-5min                  : 0
	WriteTotal                    int64     //write-total                      : 1071457
	DurationCurrentService        int64     //duration-current-service         : 0
	DurationLast5minService       int64     //duration-last-5min-service       : 30
	DurationAllService            int64     //duration-all-service             : 515
	CurrentServiceName            string    //current-service-name             :
	CpuTimeCurrent                int64     //cpu-time-current                 : 0
	CpuTimeLast5min               int64     //cpu-time-last-5min               : 280
	CpuTimeTotal                  int64     //cpu-time-total                   : 6832
	DataSeparation                string    //data-separation                  : ''
	ClientIp                      string    //client-ip                        :

}

//
//decoder.decodeUuid(buffer));
//builder.appId(decoder.decodeString(buffer)).
//blockedByDbms(decoder.decodeInt(buffer)).
//blockedByLs(decoder.decodeInt(buffer)).
//bytesAll(decoder.decodeLong(buffer)).
//bytesLast5Min(decoder.decodeLong(buffer)).
//callsAll(decoder.decodeInt(buffer)).
//callsLast5Min(decoder.decodeLong(buffer)).
//connectionId(decoder.decodeUuid(buffer)).
//dbmsBytesAll(decoder.decodeLong(buffer)).
//dbmsBytesLast5Min(decoder.decodeLong(buffer)).
//dbProcInfo(decoder.decodeString(buffer)).
//dbProcTook(decoder.decodeInt(buffer)).
//dbProcTookAt(dateFromTicks(decoder.decodeLong(buffer))).
//durationAll(decoder.decodeInt(buffer)).
//durationAllDbms(decoder.decodeInt(buffer)).
//durationCurrent(decoder.decodeInt(buffer)).
//durationCurrentDbms(decoder.decodeInt(buffer)).
//durationLast5Min(decoder.decodeLong(buffer)).
//durationLast5MinDbms(decoder.decodeLong(buffer)).
//host(decoder.decodeString(buffer)).
//infoBaseId(decoder.decodeUuid(buffer)).
//lastActiveAt(dateFromTicks(decoder.decodeLong(buffer))).
//hibernate(decoder.decodeBoolean(buffer)).
//passiveSessionHibernateTime(decoder.decodeInt(buffer)).
//hibernateSessionTerminateTime(decoder.decodeInt(buffer)).
//licenses(parseLicenseInfos(buffer, decoder)).
//locale(decoder.decodeString(buffer)).
//workingProcessId(decoder.decodeUuid(buffer)).
//sessionId(decoder.decodeInt(buffer)).
//startedAt(dateFromTicks(decoder.decodeLong(buffer))).
//userName(decoder.decodeString(buffer));
//if (version >= 4) {
//builder.memoryCurrent(decoder.decodeLong(buffer)).
//memoryLast5Min(decoder.decodeLong(buffer)).
//memoryTotal(decoder.decodeLong(buffer)).
//readBytesCurrent(decoder.decodeLong(buffer)).
//readBytesLast5Min(decoder.decodeLong(buffer)).
//readBytesTotal(decoder.decodeLong(buffer)).
//writeBytesCurrent(decoder.decodeLong(buffer)).
//writeBytesLast5Min(decoder.decodeLong(buffer)).
//writeBytesTotal(decoder.decodeLong(buffer));
//}
//if (version >= 5) {
//builder.durationCurrentService(decoder.decodeInt(buffer)).
//durationLast5MinService(decoder.decodeLong(buffer)).
//durationAllService(decoder.decodeInt(buffer)).
//currentServiceName(decoder.decodeString(buffer));
//}
//if (version >= 6) {
//builder.cpuTimeCurrent(decoder.decodeLong(buffer)).
//cpuTimeLast5Min(decoder.decodeLong(buffer)).
//cpuTimeAll(decoder.decodeLong(buffer));
//}
//if (version >= 7) {
//builder.dataSeparation(decoder.decodeString(buffer));
//}
//if (version >= 10) {
//builder.clientIPAddress(decoder.decodeString(buffer));
//}
