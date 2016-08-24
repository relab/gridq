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
