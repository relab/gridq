package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/relab/gridq/proto/gqrpc"
)

type register struct {
	sync.RWMutex
	row, col uint32
	state    gqrpc.State
}

func main() {
	port := flag.String("port", "8080", "port to listen on")
	id := flag.String("id", "", "id using the form 'row:col'")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *id == "" {
		fmt.Fprintf(os.Stderr, "no id given\n")
		flag.Usage()
		os.Exit(2)
	}

	row, col, err := parseRowCol(*id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing id: %v\n", err)
		flag.Usage()
		os.Exit(2)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatal(err)
	}

	register := &register{
		row: row,
		col: col,
	}

	grpcServer := grpc.NewServer()
	gqrpc.RegisterRegisterServer(grpcServer, register)
	log.Fatal(grpcServer.Serve(l))
}

func parseRowCol(id string) (uint32, uint32, error) {
	splitted := strings.Split(id, ":")
	if len(splitted) != 2 {
		return 0, 0, errors.New("id should have form 'row:col'")
	}
	row, err := strconv.Atoi(splitted[0])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing row from: %q", splitted[0])
	}
	col, err := strconv.Atoi(splitted[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing col from: %q", splitted[1])
	}
	return uint32(row), uint32(col), nil
}

func (r *register) Read(ctx context.Context, e *gqrpc.Empty) (*gqrpc.State, error) {
	r.RLock()
	state := r.state
	r.RUnlock()
	return &state, nil
}

func (r *register) Write(ctx context.Context, s *gqrpc.State) (*gqrpc.WriteResponse, error) {
	wr := &gqrpc.WriteResponse{}
	r.Lock()
	if s.Timestamp > r.state.Timestamp {
		r.state = *s
		wr.Written = true
	}
	r.Unlock()
	return wr, nil
}
