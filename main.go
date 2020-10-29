package main

import (
	"fmt"
	"os"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/fwereade/dgraph-messages/grpcutil"
	"github.com/fwereade/dgraph-messages/scenario"
)

func main() {
	conn, err := grpcutil.Connect(readArgs())
	check(err)
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	err = scenario.Run(dg)
	check(err)
}

func readArgs() (spec grpcutil.Spec) {
	switch args := os.Args[1:]; len(args) {
	case 3:
		spec.Token = args[2]
		fallthrough
	case 2:
		spec.Addr = args[1]
		switch args[0] {
		case "huge":
			spec.Huge = true
		case "normal":
		default:
			panic(doc())
		}
	default:
		panic(doc())
	}
	return
}

func doc() string {
	return fmt.Sprintf("%s [huge|normal] <addr> [<token-for-slash>]", os.Args[0])
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
