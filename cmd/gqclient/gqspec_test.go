package main

import (
	"strings"
	"testing"

	"github.com/relab/gridq/proto/gqrpc"
)

var gridReadQFTests = []struct {
	name    string
	replies []*gqrpc.ReadResponse
	rq      bool
}{
	{
		"nil input",
		nil,
		false,
	},
	{
		"len=0 input",
		[]*gqrpc.ReadResponse{},
		false,
	},
	{
		"no quorum (I)",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 0},
			{Row: 0, Col: 1},
		},
		false,
	},
	{
		"no quorum (II)",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 0},
			{Row: 1, Col: 0},
		},
		false,
	},
	{
		"no quorum (III)",
		[]*gqrpc.ReadResponse{
			{Row: 2, Col: 0},
			{Row: 1, Col: 0},
			{Row: 0, Col: 0},
		},
		false,
	},
	{
		"no quorum (IV)",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 0},
			{Row: 1, Col: 1},
			{Row: 0, Col: 1},
			{Row: 1, Col: 0},
		},
		false,
	},
	{
		"no quorum (V)",
		[]*gqrpc.ReadResponse{
			{Row: 2, Col: 2},
			{Row: 0, Col: 0},
			{Row: 1, Col: 1},
			{Row: 0, Col: 1},
			{Row: 2, Col: 0},
			{Row: 1, Col: 0},
		},
		false,
	},
	{
		"best-case quorum",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 0},
			{Row: 0, Col: 1},
			{Row: 0, Col: 2},
		},
		true,
	},
	{
		"col quorum",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 1},
			{Row: 1, Col: 1},
			{Row: 2, Col: 1},
		},
		false,
	},
	{
		"approx. worst-case quorum",
		[]*gqrpc.ReadResponse{
			{Row: 2, Col: 1},
			{Row: 3, Col: 2},
			{Row: 0, Col: 1},
			{Row: 2, Col: 2},
			{Row: 3, Col: 0},
			{Row: 0, Col: 0},
			{Row: 3, Col: 3},
		},
		true,
	},
}

func TestGridReadQF(t *testing.T) {
	rows, cols := 3, 3
	gqspec := &GridQuorumSpec{
		rows:      rows,
		cols:      cols,
		printGrid: false,
		vgrid:     newVisualGrid(rows, cols),
	}
	for _, test := range gridReadQFTests {
		t.Run(test.name, func(t *testing.T) {
			_, rquorum := gqspec.ReadQF(test.replies)
			if rquorum != test.rq {
				t.Errorf("got %t, want %t", rquorum, test.rq)
			}
		})
	}
}

func BenchmarkGridReadQF(b *testing.B) {
	rows, cols := 3, 3
	gqspec := &GridQuorumSpec{
		rows:      rows,
		cols:      cols,
		printGrid: false,
		vgrid:     newVisualGrid(rows, cols),
	}
	for _, test := range gridReadQFTests {
		if !strings.Contains(test.name, "case") {
			continue
		}
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gqspec.ReadQF(test.replies)
			}
		})
	}
}
