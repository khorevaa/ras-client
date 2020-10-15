package pool

import (
	"github.com/v8platform/rac/protocol/esig"
	"sort"
)

type IdleConns []*Conn

func (c IdleConns) Pop(sig esig.ESIG) *Endpoint {

	type finder struct {
		connIdx     int
		endpointIdx int
		order       int
		cap         int
	}

	var finders []finder
	var findConnIdx int
	var findEndpoint *Endpoint

	for idx, conn := range c {

		if len(conn.endpoints) == 0 {
			finders = append(finders, finder{idx, -1, 0, 0})
			continue
		}

		capEnd := len(conn.endpoints)

		for i, endpoint := range conn.endpoints {

			if esig.Equal(endpoint.sig, sig) {
				findEndpoint = endpoint
				findConnIdx = idx
				break
			}

			orderByte := 2

			if esig.HighEqual(endpoint.sig, sig) {
				orderByte = 1
			}

			finders = append(finders, finder{idx, i, orderByte, capEnd})

		}

		if findEndpoint != nil {
			break
		}

	}

	if findEndpoint != nil {
		c.remove(findConnIdx)
		return findEndpoint
	}

	if len(finders) == 0 {
		return nil
	}

	sort.Slice(finders, func(i, j int) bool {
		if finders[i].order < finders[j].order {
			return true
		}
		if finders[i].order > finders[j].order {
			return false
		}
		return finders[i].cap < finders[j].cap
	})

	f := finders[0]

	conn := c[f.connIdx]
	if f.endpointIdx == -1 {
		findEndpoint = &Endpoint{
			conn: conn,
		}
	} else {
		findEndpoint = conn.endpoints[f.endpointIdx]
	}

	c.remove(f.connIdx)
	return findEndpoint
}

func (c *IdleConns) remove(i int) {

	conns := *c
	conns[i] = conns[len(conns)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	*c = conns[:len(conns)-1]

}
