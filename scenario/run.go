package scenario

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/juju/errors"
)

var ctx = context.Background()

// Run writes about 1MB of data, and runs a query that returns about 5MB.
func Run(dg *dgo.Dgraph) error {
	rootUid, err := write(dg)
	if err != nil {
		return errors.Annotate(err, "cannot write")
	}
	fmt.Println("wrote root at:", rootUid)
	err = query(dg, rootUid)
	if err != nil {
		return errors.Annotate(err, "cannot query")
	}
	fmt.Println("queried successfully")
	return nil
}

// write returns the uid of a node with <link>s to 5 nodes, each of which has a
// <link> to another node, which has <content> weighing about 1MB.
func write(dg *dgo.Dgraph) (string, error) {
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	quads := []byte(fmt.Sprintf(`
        _:root <link> _:a .
        _:root <link> _:b .
        _:root <link> _:c .
        _:root <link> _:d .
        _:root <link> _:e .
        _:a <link> _:big .
        _:b <link> _:big .
        _:c <link> _:big .
        _:d <link> _:big .
        _:e <link> _:big .
        _:big <content> "%s" .
    `, strings.Repeat("x", 1000000)))
	resp, err := txn.Mutate(ctx, &api.Mutation{
		SetNquads: quads,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}
	return resp.Uids["root"], nil
}

// query looks for ->link->link->content starting at rootUid and prints how
// much JSON data was returned.
func query(dg *dgo.Dgraph, rootUid string) error {
	qStr := fmt.Sprintf(`{
  result(func: uid(%s)) {
    link {
      link {
        content
      }
    }
  }
}`, rootUid)

	txn := dg.NewReadOnlyTxn()
	resp, err := txn.Query(ctx, qStr)
	if err != nil {
		return err
	}
	fmt.Println("got bytes of json:", len(resp.Json))
	return nil
}
