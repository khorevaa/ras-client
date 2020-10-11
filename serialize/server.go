package serialize

type ServerInfo struct {
	UUID                                 string `rac:"server"` //server                                    : 82b8f05a-898e-48ec-9a5b-461bdf66b7d0
	AgentHost                            string //agent-host                                : app
	AgentPort                            int    //agent-port                                : 1540
	PortRange                            string //port-range                                : 1560:1591
	Name                                 string //name                                      : "Центральный сервер"
	Using                                string //using                                     : main
	DedicateManagers                     string //dedicate-managers                         : none
	InfobasesLimit                       int32  //infobases-limit                           : 8
	MemoryLimit                          int64  //memory-limit                              : 0
	ConnectionsLimit                     int32  //connections-limit                         : 128
	SafeWorkingProcessesMemoryLimit      int32  //safe-working-processes-memory-limit       : 0
	SafeCallMemoryLimit                  int32  //safe-call-memory-limit                    : 0
	ClusterPort                          int    //cluster-port                              : 1541
	CriticalTotalMemory                  int64  //critical-total-memory                     : 0
	TemporaryAllowedTotalMemory          int64  //temporary-allowed-total-memory            : 0
	TemporaryAllowedTotalMemoryTimeLimit int64  //temporary-allowed-total-memory-time-limit : 300

}

//
//final UUID workingServerId = decoder.decodeUuid(buffer);
//final String hostName = decoder.decodeString(buffer);
//final int mainPort = decoder.decodeUnsignedShort(buffer);
//final String name = decoder.decodeString(buffer);
//final boolean mainServer = decoder.decodeBoolean(buffer);
//final long safeWorkingProcessesMemoryLimit = decoder.decodeLong(buffer);
//final long safeCallMemoryLimit = decoder.decodeLong(buffer);
//final int infoBasesPerWorkingProcessLimit = decoder.decodeInt(buffer);
//final long workingProcessMemoryLimit = decoder.decodeLong(buffer);
//final int connectionsPerWorkingProcessLimit = decoder.decodeInt(buffer);
//final int clusterMainPort = decoder.decodeUnsignedShort(buffer);
//final boolean dedicatedManagers = decoder.decodeBoolean(buffer);
//final List<IPortRangeInfo> portRanges = parsePortRangeInfos(buffer, decoder);
//long criticalProcessesTotalMemory = 0L;
//long temporaryAllowedProcessesTotalMemory = 0L;
//long temporaryAllowedProcessesTotalMemoryTimeLimit = 0L;
//if (version >= 8) {
//criticalProcessesTotalMemory = decoder.decodeLong(buffer);
//temporaryAllowedProcessesTotalMemory = decoder.decodeLong(buffer);
//temporaryAllowedProcessesTotalMemoryTimeLimit = decoder.decodeLong(buffer);
//}
