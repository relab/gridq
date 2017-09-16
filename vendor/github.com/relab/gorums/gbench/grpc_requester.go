package gbench

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	rpc "github.com/relab/gorums/dev"
	"github.com/tylertreat/bench"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// GrpcRequesterFactory implements RequesterFactory by creating a Requester which
// issues requests to a storage using the gRPC framework.
type GrpcRequesterFactory struct {
	Addrs             []string
	ReadQuorum        int
	WriteQuorum       int
	PayloadSize       int
	Timeout           time.Duration
	WriteRatioPercent int
	Concurrent        bool
}

// GetRequester returns a new Requester, called for each Benchmark connection.
func (r *GrpcRequesterFactory) GetRequester(uint64) bench.Requester {
	return &grpcRequester{
		addrs:       r.Addrs,
		readq:       r.ReadQuorum,
		writeq:      r.WriteQuorum,
		payloadSize: r.PayloadSize,
		timeout:     r.Timeout,
		writeRatio:  r.WriteRatioPercent,
		concurrent:  r.Concurrent,
		dialOpts: []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(time.Second),
		},
	}
}

type client struct {
	conn   *grpc.ClientConn
	client rpc.StorageClient
}

type grpcRequester struct {
	addrs       []string
	readq       int
	writeq      int
	payloadSize int
	timeout     time.Duration
	writeRatio  int
	concurrent  bool

	dialOpts []grpc.DialOption

	clients []*client

	state *rpc.State
	ctx   context.Context
}

func (gr *grpcRequester) Setup() error {
	gr.ctx = context.Background()
	gr.clients = make([]*client, len(gr.addrs))

	for i := 0; i < len(gr.clients); i++ {
		conn, err := grpc.Dial(gr.addrs[i], gr.dialOpts...)
		if err != nil {
			return fmt.Errorf("error connecting to %q: %v", gr.addrs[i], err)
		}
		gr.clients[i] = &client{
			conn:   conn,
			client: rpc.NewStorageClient(conn),
		}
	}

	// Set initial state.
	gr.state = &rpc.State{
		Value:     strings.Repeat("x", gr.payloadSize),
		Timestamp: time.Now().UnixNano(),
	}

	for i, c := range gr.clients {
		wreply, err := c.client.Write(gr.ctx, gr.state)
		if err != nil {
			return fmt.Errorf("%s: write rpc error: %v", gr.addrs[i], err)
		}
		if !wreply.New {
			return fmt.Errorf("%s: intital write reply was not marked as new", gr.addrs[i])
		}
	}

	return nil
}

func (gr *grpcRequester) Request() error {
	write := gr.doWrite()
	if gr.concurrent {
		return gr.concurrentReq(write)
	}
	return gr.singleReq(write)
}

func (gr *grpcRequester) singleReq(write bool) error {
	client := gr.clients[0].client
	if !write {
		_, err := client.Read(gr.ctx, &rpc.ReadRequest{})
		return err
	}
	gr.state.Timestamp = time.Now().UnixNano()
	_, err := client.Write(gr.ctx, gr.state)
	return err
}

func (gr *grpcRequester) concurrentReq(write bool) error {
	if write {
		_, err := gr.writeConcurrent()
		return err
	}
	_, err := gr.readConcurrent()
	return err
}

func (gr *grpcRequester) writeConcurrent() (*rpc.WriteResponse, error) {
	replies := make(chan *rpc.WriteResponse, len(gr.clients))
	for _, c := range gr.clients {
		go func(c *client) {
			gr.state.Timestamp = time.Now().UnixNano()
			rep, err := c.client.Write(gr.ctx, gr.state)
			if err != nil {
				panic("write error")
			}
			replies <- rep
		}(c)
	}

	count := 0
	for reply := range replies {
		count++
		if count >= gr.writeq {
			return reply, nil
		}
	}

	return nil, fmt.Errorf("write incomplete")
}

func (gr *grpcRequester) readConcurrent() (*rpc.State, error) {
	replies := make(chan *rpc.State, len(gr.clients))
	for _, c := range gr.clients {
		go func(c *client) {
			rep, err := c.client.Read(gr.ctx, &rpc.ReadRequest{})
			if err != nil {
				panic("read error")
			}
			replies <- rep
		}(c)
	}

	count := 0
	for reply := range replies {
		count++
		if count >= gr.readq {
			return reply, nil
		}
	}

	return nil, fmt.Errorf("read incomplete")
}

func (gr *grpcRequester) Teardown() error {
	for _, c := range gr.clients {
		_ = c.conn.Close()
		c.conn = nil
	}
	return nil
}

func (gr *grpcRequester) doWrite() bool {
	switch gr.writeRatio {
	case 0:
		return false
	case 100:
		return true
	default:
		x := rand.Intn(100)
		if x < gr.writeRatio {
			return true
		}
		return false
	}
}
