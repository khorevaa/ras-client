package rac

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"testing"
)

type clusterTestSuite struct {
	suite.Suite
}

func TestClusterTestSuite(t *testing.T) {
	suite.Run(t, new(clusterTestSuite))
}

func (s *clusterTestSuite) r() *require.Assertions {
	return s.Require()
}

func (s *clusterTestSuite) AfterTest(suite, testName string) {

}
func (s *clusterTestSuite) BeforeTest(suite, testName string) {

}

func (s *clusterTestSuite) TearDownSuite() {

}
func (s *clusterTestSuite) TearDownTest() {

}
func (s *clusterTestSuite) SetupSuite() {

}

func (s *clusterTestSuite) TestClusterList() {

	m, _ := NewManager("srv-uk-app22:1545", WithNoUpdate(), WithPath("C:\\Program Files\\1cv8\\8.3.17.1549\\bin\\rac.exe"))
	needLen := len(m.idxCluster)
	resp, err := m.Clusters(ClustersList{})
	s.r().NoError(err)
	s.r().True(len(resp.List) == needLen, "len must be 1")

}
