package main

import (
	"context"
	"fmt"
	"log"
	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/proto/query"
	_ "vitess.io/vitess/go/vt/vtctl/grpcvtctlclient"
	_ "vitess.io/vitess/go/vt/vtgate/grpcvtgateconn"
	"vitess.io/vitess/go/vt/vtgate/vtgateconn"
)

func main() {
	ctx := context.Background()
	conn, err := vtgateconn.Dial(ctx, "localhost:15991")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	session := conn.Session("commerce", &query.ExecuteOptions{})

	_, err = session.Execute(ctx, "INSERT INTO customer (customer_id, email) VALUES (1, 'foo') ON DUPLICATE KEY UPDATE email = 'foo'", nil)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; true; i++ {
		bindVariables := map[string]*query.BindVariable{
			"email": sqltypes.StringBindVariable(fmt.Sprintf("foo-%d", i)),
		}
		_, err = session.Execute(ctx, "UPDATE customer SET email = :email WHERE customer_id = 1", bindVariables)
		if err != nil {
			log.Fatal(err)
		}

	}
}
