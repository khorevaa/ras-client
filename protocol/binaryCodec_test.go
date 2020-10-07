package protocol

import (
	"encoding/hex"

	"testing"
)

func TestNewDecoder(t *testing.T) {

	ddd, _ := hex.DecodeString("0e69")
	dec := NewDecoder(ddd)
	messageType := dec.decodeByte()
	t.Logf("message type: %v", messageType)
	size := dec.decodeByte()
	t.Logf("size: %v", size)

	testData := "010000010c014af13da5f0fc44cbae420f9fbbff76ef0000001e0c7372762d756b2d6170703232000000000605000000000000000021d09bd0bed0bad0b0d0bbd18cd0bdd18bd0b920d0bad0bbd0b0d181d182d0b5d180000000000000000000000000000000000000"
	bData, _ := hex.DecodeString(testData)

	decoder := NewDecoder(bData)

	messageFormat := decoder.decodeByte()
	t.Logf("message format: %v", messageFormat)
	//private static final byte VOID_MESSAGE = 0;
	//private static final byte MESSAGE = 1;
	//private static final byte EXCEPTION = -1;

	size2 := decoder.decodeShort()
	t.Logf("size: %v", size2)

	mType := decoder.decodeType()
	t.Logf("message type: %v", mType)
	//_ = decoder.decodeByte()
	//_ = decoder.decodeByte()

	EndpointId := decoder.decodeEndpointId()
	format := decoder.decodeByte()

	//size2 := decoder.decodeShort()
	//t.Logf("size: %v", size2)

	//b2 := decoder.decodeByte()
	//t.Logf("b2: %v", b2)

	t.Logf("endpoint: %v", EndpointId)
	t.Logf("format: %v", format)
	t.Logf("compression %v", format&0x1 != 0x0)

	info := &ClusterInfo{}
	info.UUID = decoder.decodeUUID().String()
	_ = decoder.decodeInt() // expirationTimeout
	info.Host = decoder.decodeString()
	info.ExpirationTimeout = int(decoder.decodeInt())
	info.Port = int(decoder.decodeUnsignedShort())
	info.MaxMemorySize = int(decoder.decodeInt())
	info.MaxMemoryTimeLimit = int(decoder.decodeInt())
	info.Name = decoder.decodeString()
	info.SecurityLevel = int(decoder.decodeInt())
	info.SessionFaultToleranceLevel = int(decoder.decodeInt())
	info.LoadBalancingMode = int(decoder.decodeInt()) // Не понтяно что
	info.ErrorsCountThreshold = int(decoder.decodeInt())
	info.KillProblemProcesses = decoder.decodeBoolean()
	info.KillByMemoryWithDump = decoder.decodeBoolean()

	t.Logf("info %v", info)
	//fmt.Sprintf()
	//public static IClusterInfo parseClusterInfo(final ChannelBuffer buffer, final IServiceWireDecoder decoder, final int version) throws ServiceWireCodecException {
	//	final UUID clusterId = decoder.decodeUuid(buffer);
	//	final ClusterInfo info = new ClusterInfo(clusterId);
	//	final int expirationTimeout = decoder.decodeInt(buffer);
	//	info.setExpirationTimeout(expirationTimeout);
	//	final String hostName = decoder.decodeString(buffer);
	//	info.setHostName(hostName);
	//	final int lifeTimeLimit = decoder.decodeInt(buffer);
	//	info.setLifeTimeLimit(lifeTimeLimit);
	//	final int mainPort = decoder.decodeUnsignedShort(buffer);
	//	info.setMainPort(mainPort);
	//	final int maxMemorySize = decoder.decodeInt(buffer);
	//	info.setMaxMemorySize(maxMemorySize);
	//	final int maxMemoryTimeLimit = decoder.decodeInt(buffer);
	//	info.setMaxMemoryTimeLimit(maxMemoryTimeLimit);
	//	final String name = decoder.decodeString(buffer);
	//	info.setName(name);
	//	final int securityLevel = decoder.decodeInt(buffer);
	//	info.setSecurityLevel(securityLevel);
	//	final int sessionFaultToleranceLevel = decoder.decodeInt(buffer);
	//	info.setSessionFaultToleranceLevel(sessionFaultToleranceLevel);
	//	final int loadBalancingMode = decoder.decodeInt(buffer);
	//	info.setLoadBalancingMode(loadBalancingMode);
	//	final int errorsCountThreshold = decoder.decodeInt(buffer);
	//	info.setClusterRecyclingErrorsCountThreshold(errorsCountThreshold);
	//	final boolean killProblemProcesses = decoder.decodeBoolean(buffer);
	//	info.setClusterRecyclingKillProblemProcesses(killProblemProcesses);
	//	boolean killByMemoryWithDump = false;
	//	if (version >= 8) {
	//	killByMemoryWithDump = decoder.decodeBoolean(buffer);
	//}
	//	info.setClusterRecyclingKillByMemoryWithDump(killByMemoryWithDump);
	//	return info;
	//}

}

func TestNewDecoder2(t *testing.T) {

	data2, _ := hex.DecodeString("0116010f636f6e6e6563742e74696d656f757404000007d0")

	decoder := NewDecoder(data2)
	mType := decoder.decodeNullableSize()
	t.Logf("count: %v", mType)
	//size := decoder.decodeSize()
	//t.Logf("size: %v", size)
	//
	count := decoder.decodeString()
	format := decoder.decodeShort()
	t.Logf("param: %v, value %v ", count, format)

}
