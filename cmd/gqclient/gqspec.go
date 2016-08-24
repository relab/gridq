package main

import (
	"sort"
	"sync"

	"github.com/relab/gridq/proto/gqrpc"
)

type GridQuorumSpec struct {
	rows, cols int
	printGrid  bool

	mu    sync.Mutex // Protects vgrid below if requests are issued concurrently.
	vgrid visualGrid
}

// ReadQF: All replicas from one row.
func (gqs *GridQuorumSpec) ReadQF(replies []*gqrpc.ReadResponse) (*gqrpc.ReadResponse, bool) {
	if len(replies) < gqs.rows {
		return nil, false
	}

	sort.Sort(ByRowCol(replies))

	qreplies := 0 // Count replies from the same row.
	row := replies[0].Row
	for i := 1; i < len(replies); i++ {
		if replies[i].Row != row {
			qreplies = 1
			row = replies[0].Row
			if len(replies) == FOO { // TODO: Cont.
				return nil, false // Not enough replies left.
			}
			continue
		}
		qreplies++
		if qreplies == gqs.rows {
			gqs.vgrid.setRowQuorum(row)
			if gqs.printGrid {
				gqs.vgrid.print()
			}
			return nil, true
		}
	}

	return nil, false
}

// WriteQF: One replica from each row.
func (gqs *GridQuorumSpec) WriteQF(replies []*gqrpc.WriteResponse) (*gqrpc.WriteResponse, bool) {
	if len(replies) > 0 {
		return replies[0], true
	}
	return nil, false
}
