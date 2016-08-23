package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/relab/gridq/proto/gqrpc"
)

const localAddrs = ":8080,:8081,:8082,:8083"

func main() {
	saddrs := flag.String("addrs", localAddrs, "server addresses separated by ','")
	srows := flag.Int("rows", 0, "number of rows")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	rows := *srows
	addrs := strings.Split(*saddrs, ",")
	if len(addrs) == 0 {
		dief("no server addresses provided")
	}
	if rows == 0 {
		dief("rows must be > 0")
	}
	if len(addrs)%rows != 0 {
		dief("%d addresse(s) and %d row(s) do not provide a complete grid", len(addrs), rows)
	}
	cols := len(addrs) / rows

	log.Println("#addrs:", len(addrs), "rows:", rows, "cols:", cols)

	mgr, err := gqrpc.NewManager(
		addrs,
		gqrpc.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithTimeout(5*time.Second),
		),
	)
	if err != nil {
		dief("error creating manager: %v", err)
	}

	ids := mgr.NodeIDs()

	conf, err := mgr.NewConfiguration(ids, nil, time.Second)
	if err != nil {
		dief("error creating config: %v", err)
	}

	_ = conf
}

func dief(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprint(os.Stderr, "\n")
	flag.Usage()
	os.Exit(2)
}
