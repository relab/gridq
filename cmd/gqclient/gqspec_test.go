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

const val = 42

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
			{Row: 0, Col: 0, State: &gqrpc.State{}},
			{Row: 0, Col: 1, State: &gqrpc.State{}},
		},
		false,
	},
	{
		"no quorum (II)",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 0, State: &gqrpc.State{}},
			{Row: 1, Col: 0, State: &gqrpc.State{}},
		},
		false,
	},
	{
		"no quorum (III)",
		[]*gqrpc.ReadResponse{
			{Row: 2, Col: 0, State: &gqrpc.State{}},
			{Row: 1, Col: 0, State: &gqrpc.State{}},
			{Row: 0, Col: 0, State: &gqrpc.State{}},
		},
		false,
	},
	{
		"no quorum (IV)",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 0, State: &gqrpc.State{}},
			{Row: 1, Col: 1, State: &gqrpc.State{}},
			{Row: 0, Col: 1, State: &gqrpc.State{}},
			{Row: 1, Col: 0, State: &gqrpc.State{}},
		},
		false,
	},
	{
		"no quorum (V)",
		[]*gqrpc.ReadResponse{
			{Row: 2, Col: 2, State: &gqrpc.State{}},
			{Row: 0, Col: 0, State: &gqrpc.State{}},
			{Row: 1, Col: 1, State: &gqrpc.State{}},
			{Row: 0, Col: 1, State: &gqrpc.State{}},
			{Row: 2, Col: 0, State: &gqrpc.State{}},
			{Row: 1, Col: 0, State: &gqrpc.State{}},
		},
		false,
	},
	{
		"col quorum",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 1, State: &gqrpc.State{}},
			{Row: 1, Col: 1, State: &gqrpc.State{}},
			{Row: 2, Col: 1, State: &gqrpc.State{}},
		},
		false,
	},
	{
		"best-case quorum",
		[]*gqrpc.ReadResponse{
			{Row: 0, Col: 0, State: &gqrpc.State{Timestamp: 2, Value: 9}},
			{Row: 0, Col: 1, State: &gqrpc.State{Timestamp: 3, Value: val}},
			{Row: 0, Col: 2, State: &gqrpc.State{Timestamp: 1, Value: 3}},
		},
		true,
	},
	{
		"approx. worst-case quorum",
		[]*gqrpc.ReadResponse{
			{Row: 1, Col: 0, State: &gqrpc.State{}},
			{Row: 2, Col: 1, State: &gqrpc.State{Timestamp: 2, Value: 9}},
			{Row: 0, Col: 1, State: &gqrpc.State{}},
			{Row: 1, Col: 1, State: &gqrpc.State{}},
			{Row: 2, Col: 0, State: &gqrpc.State{Timestamp: 3, Value: val}},
			{Row: 0, Col: 0, State: &gqrpc.State{}},
			{Row: 2, Col: 2, State: &gqrpc.State{Timestamp: 1, Value: 3}},
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
		"GQSort(3x3)",
		&GQSort{
			rows:      grows,
			cols:      gcols,
			printGrid: false,
			vgrid:     newVisualGrid(grows, gcols),
		},
	},
	{
		"GQMap(3x3)",
		&GQMap{
			rows:      grows,
			cols:      gcols,
			printGrid: false,
			vgrid:     newVisualGrid(grows, gcols),
		},
	},
	{
		"GQSliceOne(3x3)",
		&GQSliceOne{
			rows:      grows,
			cols:      gcols,
			printGrid: false,
			vgrid:     newVisualGrid(grows, gcols),
		},
	},
	{
		"GQSliceTwo(3x3)",
		&GQSliceTwo{
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
				replies := cloneReplies(test.replies)
				reply, rquorum := qspec.spec.ReadQF(replies)
				if rquorum != test.rq {
					t.Errorf("got %t, want %t", rquorum, test.rq)
				}
				if rquorum {
					if reply == nil || reply.State == nil {
						t.Fatalf("got nil as quorum value, want %d", val)
					}
					gotVal := reply.State.Value
					if gotVal != val {
						t.Errorf("got %d, want %d as quorum value", gotVal, val)
					}
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
				replies := cloneReplies(test.replies)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					qspec.spec.ReadQF(replies)
				}
			})
		}
	}
}

func cloneReplies(replies []*gqrpc.ReadResponse) []*gqrpc.ReadResponse {
	cloned := make([]*gqrpc.ReadResponse, len(replies))
	copy(cloned, replies)
	return cloned
}

func BenchmarkGridReadQFSuccessive(b *testing.B) {
	for _, qspec := range qspecs {
		for _, test := range gridReadQFTests {
			if !strings.Contains(test.name, "case") {
				continue
			}
			b.Run(qspec.name+"-"+test.name, func(b *testing.B) {
				replies := cloneReplies(test.replies)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for i := 0; i < len(replies); i++ {
						qspec.spec.ReadQF(replies[0 : i+1])
					}
				}
			})
		}
	}
}
