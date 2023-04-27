package main

import (
	"context"
	"log"
	"time"
	"vitess.io/vitess/go/vt/proto/query"
	_ "vitess.io/vitess/go/vt/vtctl/grpcvtctlclient"
	_ "vitess.io/vitess/go/vt/vtgate/grpcvtgateconn"
	"vitess.io/vitess/go/vt/vtgate/vtgateconn"
)

// Give the tablet some time to process each DDL. We're not trying to race DDLs against each other.
const ddlSleep = 3 * time.Second

func main() {
	ctx := context.Background()
	conn, err := vtgateconn.Dial(ctx, "localhost:15991")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	session := conn.Session("commerce", &query.ExecuteOptions{})

	for {
		_, err = session.Execute(ctx, "ALTER TABLE customer ADD foo bigint", nil)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(ddlSleep)

		_, err = session.Execute(ctx, "ALTER TABLE customer DROP foo", nil)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(ddlSleep)
	}
}
