package rac

import (
	"bytes"
	hex2 "encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type infobasesTestSuite struct {
	suite.Suite
}

func TestInfobasesTestSuite(t *testing.T) {
	suite.Run(t, new(infobasesTestSuite))
}

func (s *infobasesTestSuite) r() *require.Assertions {
	return s.Require()
}

func (s *infobasesTestSuite) AfterTest(suite, testName string) {

}
func (s *infobasesTestSuite) BeforeTest(suite, testName string) {

}

func (s *infobasesTestSuite) TearDownSuite() {

}
func (s *infobasesTestSuite) TearDownTest() {

}
func (s *infobasesTestSuite) SetupSuite() {

}

func (s *infobasesTestSuite) TestInfobasesList() {

	m, _ := NewManager("srv-uk-app22:1545")

	resp, err := m.InfobasesList()
	s.r().NoError(err)
	s.r().True(len(resp) > 0, "len must be 1")

}

func (s *infobasesTestSuite) TestInfobaseUpdate() {

	//m, _ := NewManager("srv-uk-app22:1545")
	//
	//resp, err := m.InfobasesList()
	//s.r().NoError(err)
	//s.r().True(len(resp) > 0, "len must be 1")

	update := InfobaseCreate{
		Name: "test",
	}

	val := update.Values()
	s.r().True(len(val) > 0, "must be more 0")

	hex := "1c53575001000100"
	_, _ = hex2.DecodeString(hex)

	val2 := 123009009090
	fmt.Println(fmt.Sprintf("%b", val2))
	val3 := []byte(strconv.Itoa(val2))
	fmt.Println(fmt.Sprintf("%v", val3))
	encodeedSize := encodeSize(val2)
	fmt.Println(fmt.Sprintf("%v", encodeedSize))
	decodedSize := decodeSize(encodeedSize)
	s.r().Equal(val2, decodedSize)

	fmt.Println(fmt.Sprintf("%b", decodedSize))
	val3 = []byte(strconv.Itoa(decodedSize))
	fmt.Println(fmt.Sprintf("%v", val3))

}
func decodeSize(b []byte) int {
	buf := bytes.NewBuffer(b)
	ff := 0xFFFFFF80
	b1, _ := buf.ReadByte()
	cur := int(b1 & 0xFF)
	size := cur & 0x7F
	for shift := 7; (cur & ff) != 0x0; {

		b1, _ = buf.ReadByte()
		cur = int(b1 & 0xFF)
		size += (cur & 0x7F) << shift
		shift += 7
	}

	return size

}

func encodeSize(val int) []byte {

	var buf []byte
	var b1 int

	b := bytes.NewBuffer(buf)

	msb := val >> 7
	if msb != 0 {
		b1 = -128
	} else {
		b1 = 0
	}

	b.WriteByte(byte(b1 | (val & 0x7F)))

	for val = msb; val > 0; val = msb {

		msb >>= 7
		if msb != 0 {
			b1 = -128
		} else {
			b1 = 0
		}

		b.WriteByte(byte(b1 | (val & 0x7F)))

	}

	return b.Bytes()

}

//
//public void encodeSize(final ChannelBuffer buffer, int val) {
//int n = 0;
//int msb = val >>> 7;
//BinaryCodecV1_0.temp[n++] = (byte)(((msb != 0) ? -128 : 0) | (val & 0x7F));
//for (val = msb; val > 0; val = msb) {
//msb >>>= 7;
//BinaryCodecV1_0.temp[n++] = (byte)(((msb != 0) ? -128 : 0) | (val & 0x7F));
//}
//buffer.writeBytes(BinaryCodecV1_0.temp, 0, n);
//}

//
//public int decodeSize(final ChannelBuffer buffer) throws ServiceWireCodecException {
//int cur = buffer.readByte() & 0xFF;
//int size = cur & 0x7F;
//for (int shift = 7; (cur & 0xFFFFFF80) != 0x0; cur = (buffer.readByte() & 0xFF), size += (cur & 0x7F) << shift, shift += 7) {}
//return size;
