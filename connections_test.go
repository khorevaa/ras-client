package rac

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type connectionsTestSuite struct {
	suite.Suite
}

func TestConnectionsTestSuite(t *testing.T) {
	suite.Run(t, new(connectionsTestSuite))
}

func (s *connectionsTestSuite) r() *require.Assertions {
	return s.Require()
}

func (s *connectionsTestSuite) AfterTest(suite, testName string) {

}
func (s *connectionsTestSuite) BeforeTest(suite, testName string) {

}

func (s *connectionsTestSuite) TearDownSuite() {

}
func (s *connectionsTestSuite) TearDownTest() {

}
func (s *connectionsTestSuite) SetupSuite() {

}

func (s *connectionsTestSuite) TestConnectionsList() {

	m, _ := NewManager("srv-uk-app22:1545", WithNoUpdate(), WithPath("C:\\Program Files\\1cv8\\8.3.17.1549\\bin\\rac.exe"))

	resp, err := m.ConnectionList(nil)
	s.r().NoError(err)
	s.r().True(len(resp.List) > 0, "len must be 1")

}
