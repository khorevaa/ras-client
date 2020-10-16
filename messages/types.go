package messages

type EndpointMessageKind int

func (e EndpointMessageKind) Type() int {
	return int(e)
}

const (
	VOID_MESSAGE_KIND EndpointMessageKind = 0
	MESSAGE_KIND      EndpointMessageKind = 1
	EXCEPTION_KIND    EndpointMessageKind = 0xff
)
