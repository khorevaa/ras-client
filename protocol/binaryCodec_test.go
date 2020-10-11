package protocol

import (
	"encoding/hex"
	"github.com/k0kubun/pp"
	"io/ioutil"

	"testing"
)

func TestNewDecoder(t *testing.T) {

	ddd, _ := hex.DecodeString("0e69")
	dec := NewDecoder(ddd)
	messageType := dec.decodeByte()
	t.Logf("message type: %v", messageType)
	size := dec.decodeByte()
	t.Logf("count: %v", size)

	testData := "010000010c014af13da5f0fc44cbae420f9fbbff76ef0000001e0c7372762d756b2d6170703232000000000605000000000000000021d09bd0bed0bad0b0d0bbd18cd0bdd18bd0b920d0bad0bbd0b0d181d182d0b5d180000000000000000000000000000000000000"
	bData, _ := hex.DecodeString(testData)

	decoder := NewDecoder(bData)

	EndpointId := decoder.decodeEndpointId()
	format := decoder.decodeShort()

	messageKind := decoder.decodeByte()
	//t.Logf("message kind: %d", EndpointMessageKind(messageKind))
	//private static final byte VOID_MESSAGE = 0;
	//private static final byte MESSAGE = 1;
	//private static final byte EXCEPTION = -1;
	mType := decoder.decodeUnsignedByte()
	chunkCount := decoder.decodeSize()

	//t.Logf("message type: %v", mType)
	//_ = decoder.decodeByte()
	//_ = decoder.decodeByte()
	//t.Logf("chunk count: %v", size2)

	pp.Println("message kind:", int(messageKind), "endpoint", EndpointId, "chunk count", int(chunkCount),
		"format", int(format), "respond type", int(mType))
	//size2 := decoder.decodeShort()
	//t.Logf("count: %v", size2)

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
	//count := decoder.decodeSize()
	//t.Logf("count: %v", count)
	//
	count := decoder.decodeString()
	format := decoder.decodeShort()
	t.Logf("param: %v, value %v ", count, format)

}

func TestNewDecoder3(t *testing.T) {

	data2, _ := hex.DecodeString(connectionsData)

	decoder := NewDecoder(data2)
	//m.raw = body
	endpointID := decoder.decodeEndpointId()
	format := int(decoder.decodeShort())
	kind := EndpointMessageKind(decoder.decodeByte())
	respondType := EndpointMessageType(decoder.decodeUnsignedByte())

	pp.Println(endpointID, format, kind, respondType)

	respBody, _ := ioutil.ReadAll(decoder) ///Читаем то что осталось

	parser := &GetConnectionsShortResponse{}
	parser.Parse(respBody)

}

//
//protected Object decode(final ChannelHandlerContext ctx, final Channel channel, final ChannelBuffer buffer, final MessageDecoderState state) throws Exception {
//if (buffer.readableBytes() == 0) {
//return null;
//}
//switch (state) {
//case READ_MESSAGE_TYPE: {
//this.messageType = this.wireFormat.decodeType(buffer, (IServiceWireDecoder)this.codec);
//this.checkpoint((Enum)MessageDecoderState.READ_MESSAGE_LENGTH);
//}
//case READ_MESSAGE_LENGTH: {
//this.length = this.codec.decodeSize(buffer);
//this.checkpoint((Enum)MessageDecoderState.READ_MESSAGE_BODY);
//}
//case READ_MESSAGE_BODY: {
//if (buffer.readableBytes() < this.length) {
//return null;
//}
//final ChannelBuffer body = buffer.readSlice((int)this.length);
//this.checkpoint((Enum)MessageDecoderState.READ_MESSAGE_TYPE);
//if (MessageDecoderHandler.LOGGER.isLoggable(Level.FINE)) {
//MessageDecoderHandler.LOGGER.log(Level.FINE, "Trying to parse message " + this.messageType.toString());
//}
//final IServiceWireMessage msg = this.parseMessage(this.messageType, body);
//if (MessageDecoderHandler.LOGGER.isLoggable(Level.FINE)) {
//MessageDecoderHandler.LOGGER.log(Level.FINE, "Parsed message " + msg.toString());
//}
//return msg;
//}
//default: {
//throw new err("Shouldn't reach here.");
//}
//}
//}
//
//private IServiceWireMessage parseMessage(final Type type, final ChannelBuffer buffer) throws Exception {
//if (type == Type.ENDPOINT_MESSAGE) {
//final EndpointId endpointId = this.codec.decodeEndpointId(buffer);
//final short format = this.codec.decodeShort(buffer);
//return (IServiceWireMessage)new EndpointMessage(endpointId, format, buffer.slice());
//}
//return this.wireFormat.parseMessage(buffer, (IServiceWireDecoder)this.codec, type);
