package protocol

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
	"testing"
)

func TestRASConn_CreateConnection(t *testing.T) {

	conn, err := NewRASConn("srv-uk-app10:1545")

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

}
