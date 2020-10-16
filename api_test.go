package ras_client

import (
	"context"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/v8platform/rac/serialize"
	"sync"
	"testing"
	"time"
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

	s.client = NewClient("app:1545")
	//err := s.client.CreateConnection()
	//s.r().NoError(err)

}

func (s *apiTestSuite) TestClustersList() {

	ctx := context.Background()

	clusters, err := s.client.GetClusters(ctx)

	s.r().NoError(err)
	pp.Println(clusters)

	//s.client.AuthenticateCluster(clusters[0].UUID, "admin", "admin")

	infobases, err := s.client.GetClusterInfobases(ctx, clusters[0].UUID)
	s.r().NoError(err)
	pp.Println(infobases)

	ib, _ := infobases.ByName("Тест_НИГ_03")
	s.r().NotNil(ib)
	pp.Println(ib)

}

func (s *apiTestSuite) TestGoClustersList() {

	n := 1

	wg := sync.WaitGroup{}
	i := 0
	for i < n {
		i++
		wg.Add(1)

		go func(i int) {
			defer func() {
				wg.Done()
			}()
			pp.Println("start", i)

			ctx := context.Background()
			timeout, _ := context.WithTimeout(ctx, 1000*time.Second)
			clusters, err := s.client.GetClusters(timeout)
			if err != nil {
				pp.Println("GetClusters", err)
				return
			}

			//s.r().NoErrorf(err, "go %d", num)
			//pp.Println(clusters)

			infobases, err := s.client.GetClusterInfobases(timeout, clusters[0].UUID)
			//s.r().NoErrorf(err, "go %d", num)
			//pp.Println(infobases)
			if err != nil {
				pp.Println("GetClusterInfobases", err)
				return
			}

			infobases.Each(func(info *serialize.InfobaseSummaryInfo) {
				wg.Add(1)
				go func() {
					defer wg.Done()

					conns, err := s.client.GetInfobaseConnections(ctx, clusters[0].UUID, info.UUID)
					if err != nil {
						pp.Println("GetInfobaseConnections", err)
						return
					}
					if len(conns) > 0 {
						pp.Println(info.Name, conns)
					}

				}()
			})

			//			pp.Println(ib)
			pp.Println("done", i)

		}(i)

	}

	wg.Wait()

	//err := s.client.pool.Close()
	//s.r().NoError(err)

}

//
//func (s *apiTestSuite) TestInfobaseBlockList() {
//
//	end, err := s.client.OpenEndpoint("9.0")
//	s.r().NoError(err)
//	//defer end.Close()
//
//	clusters, err := end.GetClusters()
//	s.r().NoError(err)
//
//	err = end.AuthenticateCluster(clusters[0].UUID, "", "")
//	s.r().NoError(err)
//
//	infobases, err := end.GetClusterInfobases(clusters[0].UUID)
//	s.r().NoError(err)
//
//	ib, _ := infobases.ByName("")
//	s.r().NotNil(ib)
//
//	err = end.AuthenticateInfobase(clusters[0].UUID, "", "")
//	s.r().NoError(err)
//
//	info, err := ib.FullInfo(end)
//	s.r().NoError(err)
//	//pp.Println(info)
//
//	_, err = info.GetSessions(end)
//	s.r().NoError(err)
//	_, err = info.GetConnections(end)
//	s.r().NoError(err)
//	_, err = info.GetLocks(end)
//	s.r().NoError(err)
//
//	pp.Println(info)
//
//	//blocker := info.Blocker(true)
//	//blocker.From(time.Now()).To(time.Now().Add(time.Minute)).Code("123")
//	//err = blocker.Block(end)
//	//s.r().NoError(err)
//	//
//	//info2, err := ib.FullInfo(end)
//	//s.r().NoError(err)
//	//
//	//pp.Println(info2)
//	//
//	//err = blocker.Unblock()
//	//s.r().NoError(err)
//	//pp.Println(info)
//}
//
//func (s *apiTestSuite) TestSessionList() {
//
//	end, err := s.client.OpenEndpoint("9.0")
//	s.r().NoError(err)
//	//defer end.Close()
//
//	clusters, err := end.GetClusters()
//	s.r().NoError(err)
//
//	cluster := clusters[0].UUID
//
//	err = end.AuthenticateCluster(cluster, "", "")
//	s.r().NoError(err)
//
//	sessions, err := end.GetClusterSessions(cluster)
//
//	s.r().NoError(err)
//	pp.Println(sessions)
//}
//
//func (s *apiTestSuite) TestLocksList() {
//
//	end, err := s.client.OpenEndpoint("9.0")
//	s.r().NoError(err)
//	//defer end.Close()
//
//	clusters, err := end.GetClusters()
//	s.r().NoError(err)
//
//	cluster := clusters[0].UUID
//
//	err = end.AuthenticateCluster(cluster, "", "")
//	s.r().NoError(err)
//
//	locks, err := end.GetClusterLocks(clusters[0].UUID)
//	s.r().NoError(err)
//
//	pp.Println(locks)
//
//}
//
//func (s *apiTestSuite) TestInfobasesList() {
//
//	end, err := s.client.OpenEndpoint("9.0")
//	s.r().NoError(err)
//	//defer end.Close()
//
//	clusters, err := end.GetClusters()
//	s.r().NoError(err)
//
//	cluster := clusters[0].UUID
//
//	err = end.AuthenticateCluster(cluster, "", "")
//	s.r().NoError(err)
//
//	list, err := end.GetClusterInfobases(clusters[0].UUID)
//	s.r().NoError(err)
//
//	pp.Println(list)
//
//}
