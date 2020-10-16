package messages

import (
	"github.com/k0kubun/pp"
	"github.com/v8platform/rac/protocol/codec"
	"github.com/v8platform/rac/types"
	"io"
)

type EndpointMessageFailure struct {
	ServiceID  string
	Message    string
	EndpointID int
}

func (m *EndpointMessageFailure) Parse(decoder codec.Decoder, r io.Reader) {

	decoder.StringPtr(&m.ServiceID, r)
	decoder.StringPtr(&m.Message, r)

}

func (m *EndpointMessageFailure) String() string {
	return pp.Sprintln(m)
}

func (m *EndpointMessageFailure) Type() types.Typed {
	return EXCEPTION_KIND
}

func (m *EndpointMessageFailure) Error() string {
	return pp.Sprintf("endpoint: %s service: %s msg: %s", m.EndpointID, m.ServiceID, m.Message)
}
