package serialize

import (
	"context"
	"github.com/hashicorp/go-multierror"
	uuid "github.com/satori/go.uuid"
	"io"
	"sync"
	"time"
)

type SessionCloser interface {
	TerminateSession(ctx context.Context, cluster uuid.UUID, session uuid.UUID, msg string) error
}

type SessionnSig interface {
	Sig() (cluster uuid.UUID, session uuid.UUID)
}

type SessionInfoList []*SessionInfo

func (l SessionInfoList) ByID(id uuid.UUID) (*SessionInfo, bool) {

	if id == uuid.Nil {
		return nil, false
	}

	fn := func(info *SessionInfo) bool {
		return uuid.Equal(info.UUID, id)
	}

	val := l.filter(fn, 1)

	if len(val) > 0 {
		return val[0], true
	}

	return nil, false

}

func (l SessionInfoList) ByProcess(id uuid.UUID) (SessionInfoList, bool) {

	if id == uuid.Nil {
		return SessionInfoList{}, false
	}

	fn := func(info *SessionInfo) bool {
		return uuid.Equal(info.ProcessID, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l SessionInfoList) ByInfobase(id uuid.UUID) (SessionInfoList, bool) {

	if id == uuid.Nil {
		return SessionInfoList{}, false
	}

	fn := func(info *SessionInfo) bool {
		return uuid.Equal(info.InfobaseID, id)
	}

	val := l.filter(fn, 0)

	return val, true

}

func (l SessionInfoList) Find(fn func(info *SessionInfo) bool) (SessionInfoList, bool) {

	val := l.filter(fn, 0)

	return val, true

}

func (l SessionInfoList) First(fn func(info *SessionInfo) bool) (*SessionInfo, bool) {

	val := l.filter(fn, 1)

	if len(val) == 0 {
		return nil, false
	}

	return val[0], true

}

func (l SessionInfoList) Filter(fn func(info *SessionInfo) bool) SessionInfoList {

	return l.filter(fn, 0)

}

func (l SessionInfoList) Each(fn func(info *SessionInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l SessionInfoList) TerminateSessions(ctx context.Context, closer SessionCloser, msg string) error {

	var mErr *multierror.Error
	var muErr sync.Mutex
	var wg sync.WaitGroup
	l.Each(func(info *SessionInfo) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			errDisconnect := closer.TerminateSession(ctx, info.ClusterID, info.UUID, msg)

			if errDisconnect != nil {
				muErr.Lock()
				_ = multierror.Append(mErr, errDisconnect)
				muErr.Unlock()
			}

		}()

	})
	wg.Wait()

	return mErr.ErrorOrNil()
}

func (l SessionInfoList) filter(fn func(info *SessionInfo) bool, count int) (val SessionInfoList) {

	n := 0

	for _, info := range l {

		if n == count {
			break
		}

		result := fn(info)

		if result {
			n += 1
			val = append(val, info)
		}

	}

	return

}

func (l *SessionInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls SessionInfoList

	for i := 0; i < count; i++ {

		info := &SessionInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type SessionInfo struct {
	UUID                          uuid.UUID `rac:"session"`    // UUID session                          : 1fb5f037-99e8-4924-a99d-a9e687522d32
	ID                            int       `rac:"session-id"` // ID session-id                       : 1
	InfobaseID                    uuid.UUID // InfobaseID infobase               : aea71760-15b3-485a-9a35-506eb8a0b04a
	ConnectionID                  uuid.UUID // connection                      : 8adf4514-0379-4333-a153-0b2689edc415
	ProcessID                     uuid.UUID // process                         : 1af2e54f-d95a-4370-9b45-8277280cad23
	UserName                      string    // user-name                       : АКузнецов
	Host                          string    //host                             : Sport1
	AppId                         string    //app-id                           : Designer
	Locale                        string    //locale                           : ru_RU
	StartedAt                     time.Time //started-at                       : 2018-04-09T14:51:31
	LastActiveAt                  time.Time //last-active-at                   : 2018-05-14T11:12:33
	Hibernate                     bool      // hibernate                        : no
	PassiveSessionHibernateTime   int       //passive-session-hibernate-time   : 1200
	HibernateDessionTerminateTime int       //hibernate-session-terminate-time : 86400
	BlockedByDbms                 int       //blocked-by-dbms                  : 0
	BlockedByLs                   int       //blocked-by-ls                    : 0
	BytesAll                      int64     //bytes-all                        : 105972550
	BytesLast5min                 int64     `rac:"bytes-last-5min"` //bytes-last-5min                  : 0
	CallsAll                      int       //calls-all                        : 119052
	CallsLast5min                 int64     `rac:"calls-last-5min"` //calls-last-5min                  : 0
	DbmsBytesAll                  int64     //dbms-bytes-all                   : 317824922
	DbmsBytesLast5min             int64     `rac:"dbms-bytes-last-5min"` //dbms-bytes-last-5min             : 0
	DbProcInfo                    string    //db-proc-info                     :
	DbProcTook                    int       //db-proc-took                     : 0
	DbProcTookAt                  time.Time //db-proc-took-at                  :
	DurationAll                   int       //duration-all                     : 66184
	DurationAllDbms               int       //duration-all-dbms                : 43242
	DurationCurrent               int       //duration-current                 : 0
	DurationCurrentDbms           int       //duration-current-dbms            : 0
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
	DurationCurrentService        int       //duration-current-service         : 0
	DurationLast5minService       int64     //duration-last-5min-service       : 30
	DurationAllService            int       //duration-all-service             : 515
	CurrentServiceName            string    //current-service-name             :
	CpuTimeCurrent                int64     //cpu-time-current                 : 0
	CpuTimeLast5min               int64     //cpu-time-last-5min               : 280
	CpuTimeTotal                  int64     //cpu-time-total                   : 6832
	DataSeparation                string    //data-separation                  : ''
	ClientIPAddress               string    //client-ip                        :

	Licenses  *LicenseInfoList
	ClusterID uuid.UUID
}

func (i SessionInfo) Sig() (cluster uuid.UUID, session uuid.UUID) {
	return i.ClusterID, i.UUID
}

func (i *SessionInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)
	decoder.StringPtr(&i.AppId, r)
	decoder.IntPtr(&i.BlockedByDbms, r)
	decoder.IntPtr(&i.BlockedByLs, r)
	decoder.Int64Ptr(&i.BytesAll, r)
	decoder.Int64Ptr(&i.BytesLast5min, r)
	decoder.IntPtr(&i.CallsAll, r)
	decoder.Int64Ptr(&i.CallsLast5min, r)
	decoder.UuidPtr(&i.ConnectionID, r)
	decoder.Int64Ptr(&i.DbmsBytesAll, r)
	decoder.Int64Ptr(&i.DbmsBytesLast5min, r)
	decoder.StringPtr(&i.DbProcInfo, r)
	decoder.IntPtr(&i.DbProcTook, r)
	decoder.TimePtr(&i.DbProcTookAt, r)
	decoder.IntPtr(&i.DurationAll, r)
	decoder.IntPtr(&i.DurationAllDbms, r)
	decoder.IntPtr(&i.DurationCurrent, r)
	decoder.IntPtr(&i.DurationCurrentDbms, r)
	decoder.Int64Ptr(&i.DurationLast5Min, r)
	decoder.Int64Ptr(&i.DurationLast5MinDbms, r)

	decoder.StringPtr(&i.Host, r)
	decoder.UuidPtr(&i.InfobaseID, r)
	decoder.TimePtr(&i.LastActiveAt, r)
	decoder.BoolPtr(&i.Hibernate, r)
	decoder.IntPtr(&i.PassiveSessionHibernateTime, r)
	decoder.IntPtr(&i.HibernateDessionTerminateTime, r)

	licenseList := &LicenseInfoList{}
	licenseList.Parse(decoder, version, r)
	i.Licenses = licenseList

	decoder.StringPtr(&i.Locale, r)
	decoder.UuidPtr(&i.ProcessID, r)
	decoder.IntPtr(&i.ID, r)
	decoder.TimePtr(&i.StartedAt, r)
	decoder.StringPtr(&i.UserName, r)

	if version >= 4 {
		decoder.Int64Ptr(&i.MemoryCurrent, r)
		decoder.Int64Ptr(&i.MemoryLast5min, r)
		decoder.Int64Ptr(&i.MemoryTotal, r)
		decoder.Int64Ptr(&i.ReadCurrent, r)
		decoder.Int64Ptr(&i.ReadLast5min, r)
		decoder.Int64Ptr(&i.ReadTotal, r)
		decoder.Int64Ptr(&i.WriteCurrent, r)
		decoder.Int64Ptr(&i.WriteLast5min, r)
		decoder.Int64Ptr(&i.WriteTotal, r)
	}

	if version >= 5 {
		decoder.IntPtr(&i.DurationCurrentService, r)
		decoder.Int64Ptr(&i.DurationLast5minService, r)
		decoder.IntPtr(&i.DurationAllService, r)
		decoder.StringPtr(&i.CurrentServiceName, r)
	}

	if version >= 6 {
		decoder.Int64Ptr(&i.CpuTimeCurrent, r)
		decoder.Int64Ptr(&i.CpuTimeLast5min, r)
		decoder.Int64Ptr(&i.CpuTimeTotal, r)
	}

	if version >= 7 {
		decoder.StringPtr(&i.DataSeparation, r)
	}

	if version >= 10 {
		decoder.StringPtr(&i.ClientIPAddress, r)
	}

	i.Licenses.Each(func(info *LicenseInfo) {
		info.SessionID = i.UUID
		info.ProcessID = i.ProcessID
	})

}
