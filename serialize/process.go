package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

type ProcessInfoList []*ProcessInfo

func (l ProcessInfoList) Each(fn func(info *ProcessInfo)) {

	for _, info := range l {

		fn(info)

	}

}

func (l ProcessInfoList) Filter(fn func(info *ProcessInfo) bool) ProcessInfoList {

	return l.filter(fn, 0)

}

func (l ProcessInfoList) filter(fn func(info *ProcessInfo) bool, count int) (val ProcessInfoList) {

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

func (l *ProcessInfoList) Parse(decoder Decoder, version int, r io.Reader) {

	count := decoder.Size(r)
	var ls ProcessInfoList

	for i := 0; i < count; i++ {

		info := &ProcessInfo{}
		info.Parse(decoder, version, r)

		ls = append(ls, info)
	}

	*l = ls
}

type ProcessInfo struct {
	UUID                uuid.UUID `rac:"process"` // process              : 3ea9968d-159c-4b5f-9bdc-22b8ead96f74
	Host                string    //host                 : Sport1
	Port                int16     //port                 : 1564
	Pid                 string    //pid                  : 5428
	Enable              bool      `rac:"is-enable"` //is-enable            : yes
	Running             bool      //running              : yes
	StartedAt           time.Time //started-at           : 2018-03-29T11:16:02
	Use                 bool      //use                  : used
	AvailablePerfomance int       //available-perfomance : 100
	Capacity            int       //capacity             : 1000
	Connections         int       //connections          : 7
	MemorySize          int       //memory-size          : 1518604
	MemoryExcessTime    int       //memory-excess-time   : 0
	SelectionSize       int       //selection-size       : 61341
	AvgBackCallTime     float64   //avg-back-call-time   : 0.000
	AvgCallTime         float64   //avg-call-time        : 0.483
	AvgDbCallTime       float64   //avg-db-call-time     : 0.124
	AvgLockCallTime     float64   //avg-lock-call-time   : 0.000
	AvgServerCallTime   float64   //avg-server-call-time : -0.265
	AvgThreads          float64   //avg-threads          : 0.281
	Reverse             bool      //reserve              : no
	Licenses            LicenseInfoList

	ClusterID uuid.UUID
}

func (i *ProcessInfo) Parse(decoder Decoder, version int, r io.Reader) {

	decoder.UuidPtr(&i.UUID, r)

	decoder.DoublePtr(&i.AvgBackCallTime, r)
	decoder.DoublePtr(&i.AvgCallTime, r)
	decoder.DoublePtr(&i.AvgDbCallTime, r)
	decoder.DoublePtr(&i.AvgLockCallTime, r)
	decoder.DoublePtr(&i.AvgServerCallTime, r)
	decoder.DoublePtr(&i.AvgThreads, r)
	decoder.IntPtr(&i.Capacity, r)
	decoder.IntPtr(&i.Connections, r)
	decoder.StringPtr(&i.Host, r)
	decoder.BoolPtr(&i.Enable, r)
	decoder.StringPtr(&i.Host, r)

	licenseList := LicenseInfoList{}
	licenseList.Parse(decoder, version, r)
	i.Licenses = licenseList

	decoder.ShortPtr(&i.Port, r)
	decoder.IntPtr(&i.MemoryExcessTime, r)
	decoder.IntPtr(&i.MemorySize, r)

	decoder.StringPtr(&i.Pid, r)

	running := decoder.Int(r)
	if running == 1 {
		i.Running = true
	}

	decoder.IntPtr(&i.SelectionSize, r)
	decoder.TimePtr(&i.StartedAt, r)

	use := decoder.Int(r)
	if use == 1 {
		i.Use = true
	}

	decoder.IntPtr(&i.AvailablePerfomance, r)

	if version >= 9 {
		decoder.BoolPtr(&i.Reverse, r)
	}

	i.Licenses.Each(func(info *LicenseInfo) {
		info.ProcessID = i.UUID
	})

}
