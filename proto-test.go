package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	binlogdatapb "vitess.io/vitess/go/vt/proto/binlogdata"
)

// ✅ s48 - 1 bad record (533) -> bad_record_table=_bf_checkpoint, ddl_modified_table=_bf_checkpoint, time_updated=2023-04-06 17:03:15 alter table _bf_checkpoint add column source_id bigint(20) default 48 after id, add UNIQUE KEY tmp_source_id (source_id)

// ✅ s21 - 1 bad record (1107) -> bad_record_table=_bf_verify_audit, ddl_modified_table=_bf_verify_audit, time_updated=2023-04-04 18:59:38 alter table _bf_verify_audit add INDEX resolutions (source_id, is_resolved, failure_count, updated_at)

// ✅ s78 - 2 bad records (786, 802) -> [786] bad_record_table=_bf_verify_audit, ddl_modified_table=_bf_verify_audit, time_updated=2023-04-04 18:53:16 alter table _bf_verify_audit add column source_id bigint(20) default null after id
// 															     -> [802] bad_record_table=_bf_checkpoint, ddl_modified_table=_bf_checkpoint, time_upated=2023-04-06 16:43:10 alter table _bf_checkpoint add column source_id bigint(20) default 78 after id, add UNIQUE KEY tmp_source_id (source_id)
// s88 (706, 766, 771, 775)
const schemaVersionIdStart = 0
var getSchemaVersions = fmt.Sprintf("select id, schemax from _vt.schema_version where id > %v order by id asc", schemaVersionIdStart)
const mysqlUser = "sys.vt_dba.1"
const mysqlHost = "127.0.0.1"
const mysqlPort = "3355"
const mysqlDb = "_vt"

func main() {
	// Open a connection to the database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home dir %v", err)
		return
	}
	password, err := ioutil.ReadFile(fmt.Sprintf("%s/Desktop/dontleakme.txt", homeDir))
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, strings.TrimSpace(string(password)), mysqlHost, mysqlPort, mysqlDb)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getSchemaVersions)
	if err != nil {
		fmt.Println("Failed to prepare query:", err)
		return
	}
	defer stmt.Close()

	// Execute the query to retrieve the protobuf message
	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("Error reading schema_tracking table %v, will operate with the latest available schema", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var data []byte
		err := rows.Scan(&id, &data)
		if err != nil {
			fmt.Println("Failed to retrieve data:", err)
			return
		} else {
			fmt.Printf("Retrieved id %d with %d bytes of data\n", id, len(data))
		}

		// Parse the protobuf message
		sch := &binlogdatapb.MinimalSchema{}
		if err := sch.UnmarshalVT(data); err != nil {
			fmt.Printf("Error on id %d: %v", id, err)
			return
		}
	}
	fmt.Println("✅ - All schema versions unmarshaled successfully")
}
