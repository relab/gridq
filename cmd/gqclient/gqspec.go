package main

import (
	"sort"

	"github.com/relab/gridq/proto/gqrpc"
)

type GridQuorumSpec struct {
	rows, cols int
	printGrid  bool
	vgrid      *visualGrid
}

// ReadQF: All replicas from one row.
func (gqs *GridQuorumSpec) ReadQF(replies []*gqrpc.ReadResponse) (*gqrpc.ReadResponse, bool) {
	if len(replies) < gqs.rows {
		return nil, false
	}

	sort.Sort(ByRowCol(replies))

	qreplies := 1 // Counter for replies from the same row.
	row := replies[0].Row
	for i := 1; i < len(replies); i++ {
		if replies[i].Row != row {
			qreplies = 1
			row = replies[i].Row
			left := len(replies) - i - 1
			if qreplies+left < gqs.rows {
				// Not enough replies left.
				return nil, false
			}
			continue
		}
		qreplies++
		if qreplies == gqs.rows {
			gqs.vgrid.setRowQuorum(row)
			if gqs.printGrid {
				gqs.vgrid.print()
			}
			// Return the last reply. The replies forming a quorum
			// should be sorted using the timestamps, but we don't
			// want that logic to impact the benchmarks.
			return replies[i], true
		}
	}

	panic("an invariant was not handled")
}

// WriteQF: One replica from each row.
func (gqs *GridQuorumSpec) WriteQF(replies []*gqrpc.WriteResponse) (*gqrpc.WriteResponse, bool) {
	if len(replies) > 0 {
		return replies[0], true
	}
	panic("an invariant was not handled")
}
