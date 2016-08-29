package main

import (
	"sort"

	"github.com/relab/gridq/proto/gqrpc"
)

type GQSort struct {
	rows, cols int
	printGrid  bool
	vgrid      *visualGrid
}

// ReadQF: All replicas from one row.
func (gqs *GQSort) ReadQF(replies []*gqrpc.ReadResponse) (*gqrpc.ReadResponse, bool) {
	if len(replies) < gqs.rows {
		return nil, false
	}

	sort.Sort(ByRowTimestamp(replies))

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
			if gqs.printGrid {
				gqs.vgrid.setRowQuorum(row)
				gqs.vgrid.print()
			}
			return replies[i-gqs.rows+1], true
		}
	}

	panic("an invariant was not handled")
}

// WriteQF: One replica from each row.
func (gqs *GQSort) WriteQF(replies []*gqrpc.WriteResponse) (*gqrpc.WriteResponse, bool) {
	if len(replies) < gqs.cols {
		return nil, false
	}

	sort.Sort(ByColRow(replies))

	qreplies := 1 // Counter for replies from the same row.
	col := replies[0].Col
	for i := 1; i < len(replies); i++ {
		if replies[i].Col != col {
			qreplies = 1
			col = replies[i].Col
			left := len(replies) - i - 1
			if qreplies+left < gqs.cols {
				// Not enough replies left.
				return nil, false
			}
			continue
		}
		qreplies++
		if qreplies == gqs.cols {
			if gqs.printGrid {
				gqs.vgrid.setColQuorum(col)
				gqs.vgrid.print()
			}
			// WriteResponses don't have timestamps.
			return replies[i], true
		}
	}

	panic("an invariant was not handled")
}

type GQMap struct {
	rows, cols int
	printGrid  bool
	vgrid      *visualGrid
}

// ReadQF: All replicas from one row.
//
// Note: It is not enough to just know that we have a quorum from a row, we also
// need to know what replies forms the quorum (both in practice and to be fair
// to GQSort above).
func (gqm *GQMap) ReadQF(replies []*gqrpc.ReadResponse) (*gqrpc.ReadResponse, bool) {
	if len(replies) < gqm.rows {
		return nil, false
	}

	rowReplies := make(map[uint32][]*gqrpc.ReadResponse)
	var row []*gqrpc.ReadResponse
	var found bool
	for _, reply := range replies {
		row, found = rowReplies[reply.Row]
		if !found {
			row = make([]*gqrpc.ReadResponse, 0, gqm.rows)
		}
		row = append(row, reply)
		if len(row) >= gqm.rows {
			if gqm.printGrid {
				gqm.vgrid.setRowQuorum(reply.Row)
				gqm.vgrid.print()
			}
			sort.Sort(ByTimestamp(row))
			return row[len(row)-1], true
		}
		rowReplies[reply.Row] = row
	}

	return nil, false
}

// WriteQF: One replica from each row.
func (gqm *GQMap) WriteQF(replies []*gqrpc.WriteResponse) (*gqrpc.WriteResponse, bool) {
	panic("not implemented, symmetric with read")
}

type GQSliceOne struct {
	rows, cols int
	printGrid  bool
	vgrid      *visualGrid
}

// ReadQF: All replicas from one row.
func (gqs *GQSliceOne) ReadQF(replies []*gqrpc.ReadResponse) (*gqrpc.ReadResponse, bool) {
	if len(replies) < gqs.rows {
		return nil, false
	}
	rowCount := make([]int, gqs.rows)
	repliesRM := make([]*gqrpc.ReadResponse, gqs.rows*gqs.cols) // row-major
	for _, reply := range replies {
		repliesRM[(int(reply.Row)*gqs.rows)+int(reply.Col)] = reply
		rowCount[reply.Row]++
		if rowCount[reply.Row] >= gqs.rows {
			start := int(reply.Row) * gqs.rows
			repliesRM = repliesRM[start : start+gqs.rows]
			sort.Sort(ByTimestamp(repliesRM))
			return repliesRM[len(repliesRM)-1], true
		}
	}

	return nil, false
}

// WriteQF: One replica from each row.
func (gqs *GQSliceOne) WriteQF(replies []*gqrpc.WriteResponse) (*gqrpc.WriteResponse, bool) {
	panic("not implemented, symmetric with read")
}

type GQSliceTwo struct {
	rows, cols int
	printGrid  bool
	vgrid      *visualGrid
}
type rowInfo struct {
	count int                 // Replies for this row.
	reply *gqrpc.ReadResponse // Reply with highest timestamp.
}

// ReadQF: All replicas from one row.
func (gqs *GQSliceTwo) ReadQF(replies []*gqrpc.ReadResponse) (*gqrpc.ReadResponse, bool) {
	if len(replies) < gqs.rows {
		return nil, false
	}
	rows := make([]*rowInfo, gqs.rows)
	for _, reply := range replies {
		ri := rows[int(reply.Row)]
		if rows[int(reply.Row)] == nil {
			ri = &rowInfo{}
			rows[int(reply.Row)] = ri
		}
		if ri.reply == nil {
			ri.reply = reply
		} else if reply.State.Timestamp > ri.reply.State.Timestamp {
			ri.reply = reply
		}
		ri.count++
		if ri.count >= gqs.rows {
			if gqs.printGrid {
				gqs.vgrid.setRowQuorum(reply.Row)
				gqs.vgrid.print()
			}
			return ri.reply, true
		}
	}

	return nil, false
}

// WriteQF: One replica from each row.
func (gqs *GQSliceTwo) WriteQF(replies []*gqrpc.WriteResponse) (*gqrpc.WriteResponse, bool) {
	panic("not implemented, symmetric with read")
}
