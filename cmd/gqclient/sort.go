package main

import "github.com/relab/gridq/proto/gqrpc"

type ByRowCol []*gqrpc.ReadResponse

func (p ByRowCol) Len() int      { return len(p) }
func (p ByRowCol) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ByRowCol) Less(i, j int) bool {
	if p[i].Row < p[j].Row {
		return true
	} else if p[i].Row > p[j].Row {
		return false
	} else {
		return p[i].Col < p[j].Col
	}
}

type ByTimestamp []*gqrpc.ReadResponse

func (p ByTimestamp) Len() int           { return len(p) }
func (p ByTimestamp) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByTimestamp) Less(i, j int) bool { return p[i].State.Timestamp < p[j].State.Timestamp }

type ByColRow []*gqrpc.WriteResponse

func (p ByColRow) Len() int      { return len(p) }
func (p ByColRow) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ByColRow) Less(i, j int) bool {
	if p[i].Col < p[j].Col {
		return true
	} else if p[i].Col > p[j].Col {
		return false
	} else {
		return p[i].Row < p[j].Row
	}
}
