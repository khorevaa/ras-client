package protocol

import (
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type apiTestSuite struct {
	suite.Suite
	client *Client
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}

func (s *apiTestSuite) r() *require.Assertions {
	return s.Require()
}

func (s *apiTestSuite) TearDownSuite() {

}
func (s *apiTestSuite) TearDownTest() {

}
func (s *apiTestSuite) SetupSuite() {

	s.client = NewClient("app4:1545")
	err := s.client.CreateConnection()
	s.r().NoError(err)

}

func (s *apiTestSuite) TestInfobaseBlockList() {

	end, err := s.client.OpenEndpoint("9.0")
	s.r().NoError(err)
	//defer end.Close()

	clusters, err := end.GetClusters()
	s.r().NoError(err)

	err = end.AuthenticateCluster(clusters[0].UUID, "", "")
	s.r().NoError(err)

	infobases, err := end.GetClusterInfobases(clusters[0].UUID)
	s.r().NoError(err)

	ib, _ := infobases.ByName("Тест_03")
	s.r().NotNil(ib)

	err = end.AuthenticateInfobase(clusters[0].UUID, "", "")
	s.r().NoError(err)

	info, err := ib.FullInfo(end)
	s.r().NoError(err)
	//pp.Println(info)

	_, err = info.GetSessions(end)
	s.r().NoError(err)
	_, err = info.GetConnections(end)
	s.r().NoError(err)
	_, err = info.GetLocks(end)
	s.r().NoError(err)

	pp.Println(info)

	//blocker := info.Blocker(true)
	//blocker.From(time.Now()).To(time.Now().Add(time.Minute)).Code("123")
	//err = blocker.Block(end)
	//s.r().NoError(err)
	//
	//info2, err := ib.FullInfo(end)
	//s.r().NoError(err)
	//
	//pp.Println(info2)
	//
	//err = blocker.Unblock()
	//s.r().NoError(err)
	//pp.Println(info)
}

func (s *apiTestSuite) TestSessionList() {

	end, err := s.client.OpenEndpoint("9.0")
	s.r().NoError(err)
	//defer end.Close()

	clusters, err := end.GetClusters()
	s.r().NoError(err)

	cluster := clusters[0].UUID

	err = end.AuthenticateCluster(cluster, "", "")
	s.r().NoError(err)

	sessions, err := end.GetClusterSessions(cluster)

	s.r().NoError(err)
	pp.Println(sessions)
}

func (s *apiTestSuite) TestLocksList() {

	end, err := s.client.OpenEndpoint("9.0")
	s.r().NoError(err)
	//defer end.Close()

	clusters, err := end.GetClusters()
	s.r().NoError(err)

	cluster := clusters[0].UUID

	err = end.AuthenticateCluster(cluster, "", "")
	s.r().NoError(err)

	locks, err := end.GetClusterLocks(clusters[0].UUID)
	s.r().NoError(err)

	pp.Println(locks)

}

func (s *apiTestSuite) TestInfobasesList() {

	end, err := s.client.OpenEndpoint("9.0")
	s.r().NoError(err)
	//defer end.Close()

	clusters, err := end.GetClusters()
	s.r().NoError(err)

	cluster := clusters[0].UUID

	err = end.AuthenticateCluster(cluster, "", "")
	s.r().NoError(err)

	list, err := end.GetClusterInfobases(clusters[0].UUID)
	s.r().NoError(err)

	pp.Println(list)

}