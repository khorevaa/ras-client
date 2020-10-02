package rac

import (
	hex2 "encoding/hex"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	str, _ := hex2.DecodeString(hex)

	println(str)

}
