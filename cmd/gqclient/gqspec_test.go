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

const grows, gcols = 3, 3

var qspecs = []struct {
	name string
	spec gqrpc.QuorumSpec
}{
	{
		"GQSort",
		&GQSort{
			rows:      grows,
			cols:      gcols,
			printGrid: false,
			vgrid:     newVisualGrid(grows, gcols),
		},
	},
	{
		"GQMap",
		&GQMap{
			rows:      grows,
			cols:      gcols,
			printGrid: false,
			vgrid:     newVisualGrid(grows, gcols),
		},
	},
	{
		"GQSlice",
		&GQSlice{
			rows:      grows,
			cols:      gcols,
			printGrid: false,
			vgrid:     newVisualGrid(grows, gcols),
		},
	},
}

func TestGridReadQF(t *testing.T) {
	for _, qspec := range qspecs {
		for _, test := range gridReadQFTests {
			t.Run(qspec.name+"-"+test.name, func(t *testing.T) {
				_, rquorum := qspec.spec.ReadQF(test.replies)
				if rquorum != test.rq {
					t.Errorf("got %t, want %t", rquorum, test.rq)
				}
			})
		}
	}
}

func BenchmarkGridReadQF(b *testing.B) {
	for _, qspec := range qspecs {
		for _, test := range gridReadQFTests {
			if !strings.Contains(test.name, "case") {
				continue
			}
			b.Run(qspec.name+"-"+test.name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					qspec.spec.ReadQF(test.replies)
				}
			})
		}
	}
}
