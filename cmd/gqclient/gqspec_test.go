package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/relab/gridq/proto/gqrpc"
)

func TestMain(m *testing.M) {
	silentLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	grpclog.SetLogger(silentLogger)
	grpc.EnableTracing = false
	res := m.Run()
	os.Exit(res)
}

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
		"col quorum",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 1},
			{Row: 1, Col: 1},
			{Row: 2, Col: 1},
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
		"approx. worst-case quorum",
		[]*gqrpc.ReadResponse{
			{Row: 1, Col: 0},
			{Row: 2, Col: 1},
			{Row: 0, Col: 1},
			{Row: 1, Col: 1},
			{Row: 2, Col: 0},
			{Row: 0, Col: 0},
			{Row: 2, Col: 2},
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

var gridWriteQFTests = []struct {
	name    string
	replies []*gqrpc.WriteResponse
	rq      bool
}{
	{
		"nil input",
		nil,
		false,
	},
	{
		"len=0 input",
		[]*gqrpc.WriteResponse{},
		false,
	},
	{
		"no quorum (I)",
		[]*gqrpc.WriteResponse{
			{Row: 0, Col: 0},
			{Row: 1, Col: 0},
		},
		false,
	},
	{
		"no quorum (II)",
		[]*gqrpc.WriteResponse{
			{Row: 0, Col: 0},
			{Row: 0, Col: 1},
		},
		false,
	},
	{
		"no quorum (III)",
		[]*gqrpc.WriteResponse{
			{Row: 0, Col: 0},
			{Row: 0, Col: 1},
			{Row: 0, Col: 2},
		},
		false,
	},
	{
		"no quorum (IV)",
		[]*gqrpc.WriteResponse{
			{Row: 1, Col: 1},
			{Row: 0, Col: 0},
			{Row: 1, Col: 0},
			{Row: 0, Col: 1},
		},
		false,
	},
	{
		"no quorum (V)",
		[]*gqrpc.WriteResponse{
			{Row: 2, Col: 2},
			{Row: 0, Col: 0},
			{Row: 1, Col: 1},
			{Row: 1, Col: 0},
			{Row: 0, Col: 2},
			{Row: 0, Col: 1},
		},
		false,
	},
	{
		"row quorum",
		[]*gqrpc.WriteResponse{
			{Row: 0, Col: 0},
			{Row: 0, Col: 1},
			{Row: 0, Col: 2},
		},
		false,
	},

	{
		"best-case quorum",
		[]*gqrpc.WriteResponse{
			{Row: 0, Col: 0},
			{Row: 1, Col: 0},
			{Row: 2, Col: 0},
		},
		true,
	},
	{
		"approx. worst-case quorum",
		[]*gqrpc.WriteResponse{
			{Row: 0, Col: 1},
			{Row: 1, Col: 2},
			{Row: 1, Col: 0},
			{Row: 1, Col: 1},
			{Row: 0, Col: 2},
			{Row: 0, Col: 0},
			{Row: 2, Col: 2},
		},
		true,
	},
}

func TestGridWriteQF(t *testing.T) {
	rows, cols := 3, 3
	gqspec := &GridQuorumSpec{
		rows:      rows,
		cols:      cols,
		printGrid: false,
		vgrid:     newVisualGrid(rows, cols),
	}
	for _, test := range gridWriteQFTests {
		t.Run(test.name, func(t *testing.T) {
			_, rquorum := gqspec.WriteQF(test.replies)
			if rquorum != test.rq {
				t.Errorf("got %t, want %t", rquorum, test.rq)
			}
		})
	}
}

func BenchmarkGridWriteQF(b *testing.B) {
	rows, cols := 3, 3
	gqspec := &GridQuorumSpec{
		rows:      rows,
		cols:      cols,
		printGrid: false,
		vgrid:     newVisualGrid(rows, cols),
	}
	for _, test := range gridWriteQFTests {
		if !strings.Contains(test.name, "case") {
			continue
		}
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gqspec.WriteQF(test.replies)
			}
		})
	}
}
