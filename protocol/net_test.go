package protocol

import (
	"fmt"
	"github.com/k0kubun/pp"
	uuid "github.com/satori/go.uuid"
	"github.com/xelaj/go-dry"
	"testing"
)

func TestRASConn_CreateConnection(t *testing.T) {

	conn, err := NewRASConn("srv-uk-app22:1545")

	//conn, err := net.Dial("tcp", "srv-uk-app22:1545")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	err = conn.CreateConnection()
	dry.PanicIfErr(err)

	defer conn.Disconnect()

	resp, err := conn.OpenEndpoint("9.0")
	dry.PanicIfErr(err)

	pp.Println(resp)

	//err = conn.AuthenticateAgent("", "")
	//dry.PanicIfErr(err)

	resp2, err := conn.GetClusters()

	dry.PanicIfErr(err)

	pp.Println(resp2)

	id, _ := uuid.FromString(resp2[0].UUID)

	err = conn.AuthenticateCluster(id, "", "")

	dry.PanicIfErr(err)

	//r, err := conn.GetClusterManagers(id)
	//dry.PanicIfErr(err)
	//
	//pp.Println(r)

	r2, err := conn.GetClusterConnections(id)
	dry.PanicIfErr(err)

	pp.Println(r2)

}
