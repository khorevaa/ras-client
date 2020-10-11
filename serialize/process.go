package serialize

import "time"

type ProcessInfo struct {
	UUID                string    `rac:"process"` // process              : 3ea9968d-159c-4b5f-9bdc-22b8ead96f74
	Host                string    //host                 : Sport1
	Port                string    //port                 : 1564
	Pid                 int       //pid                  : 5428
	Enable              bool      `rac:"is-enable"` //is-enable            : yes
	Running             bool      //running              : yes
	StartedAt           time.Time //started-at           : 2018-03-29T11:16:02
	Use                 string    //use                  : used
	AvailablePerfomance int       //available-perfomance : 100
	Capacity            int32     //capacity             : 1000
	Connections         int32     //connections          : 7
	MemorySize          int64     //memory-size          : 1518604
	MemoryExcessTime    int64     //memory-excess-time   : 0
	SelectionSize       int64     //selection-size       : 61341
	AvgBackCallTime     float64   //avg-back-call-time   : 0.000
	AvgCallTime         float64   //avg-call-time        : 0.483
	AvgDbCallTime       float64   //avg-db-call-time     : 0.124
	AvgLockCallTime     float64   //avg-lock-call-time   : 0.000
	AvgServerCallTime   float64   //avg-server-call-time : -0.265
	AvgThreads          float64   //avg-threads          : 0.281
	Reverse             bool      //reserve              : no
}

//
//
//final UUID processId = decoder.decodeUuid(buffer);
//final double avgBackCallTime = decoder.decodeDouble(buffer);
//final double avgCallTime = decoder.decodeDouble(buffer);
//final double avgDBCallTime = decoder.decodeDouble(buffer);
//final double avgLockCallTime = decoder.decodeDouble(buffer);
//final double avgServerCallTime = decoder.decodeDouble(buffer);
//final double avgThreads = decoder.decodeDouble(buffer);
//final int capacity = decoder.decodeInt(buffer);
//final int connections = decoder.decodeInt(buffer);
//final String hostName = decoder.decodeString(buffer);
//final boolean enable = decoder.decodeBoolean(buffer);
//final List<ILicenseInfo> licenses = parseLicenseInfos(buffer, decoder);
//final int mainPort = decoder.decodeUnsignedShort(buffer);
//final int memoryExcessTime = decoder.decodeInt(buffer);
//final int memorySize = decoder.decodeInt(buffer);
//final String pid = decoder.decodeString(buffer);
//final int running = decoder.decodeInt(buffer);
//final int selectionSize = decoder.decodeInt(buffer);
//final long startedAt = decoder.decodeLong(buffer);
//final int use = decoder.decodeInt(buffer);
//final int availablePerfomance = decoder.decodeInt(buffer);
//boolean reserve = false;
//if (version >= 9) {
//reserve = decoder.decodeBoolean(buffer);
//}
