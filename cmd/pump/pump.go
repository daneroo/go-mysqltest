package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/daneroo/go-ted1k/ephemeral"
	"github.com/daneroo/go-ted1k/ipfs"
	"github.com/daneroo/go-ted1k/jsonl"
	"github.com/daneroo/go-ted1k/merge"
	"github.com/daneroo/go-ted1k/mysql"
	"github.com/daneroo/go-ted1k/postgres"
	"github.com/daneroo/go-ted1k/progress"
	"github.com/daneroo/go-ted1k/types"
	_ "github.com/go-sql-driver/mysql"
	shell "github.com/ipfs/go-ipfs-api"
)

const (
	myCredentials = "ted:secret@tcp(0.0.0.0:3306)/ted"
	// pgCredentials    = "postgres://postgres:secret@127.0.0.1:5432/ted"
	pgCredentials    = "postgres://postgres:secret@0.0.0.0:5432/ted"
	fmtRFC3339Millis = "2006-01-02T15:04:05.000Z07:00"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format(fmtRFC3339Millis) + " - " + string(bytes))
}

type entryWriter interface {
	Write(src <-chan []types.Entry) (int, error)
}
type entryReader interface {
	Read() <-chan []types.Entry
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	log.Printf("Starting TED1K pump\n") // TODO(daneroo): add version,buildDate

	tableNames := []string{"watt", "watt2"}
	db := mysql.Setup(tableNames, myCredentials)
	defer db.Close()
	conn := postgres.Setup(context.Background(), tableNames, pgCredentials)
	defer conn.Close(context.Background())
	sh := shell.NewShell("localhost:5001")

	// ephemeral
	if true {
		fmt.Println()
		doTest("ephemeral -> ephemeral", ephemeral.NewReader(), ephemeral.NewWriter())
		verify("ephemeral <-> ephemeral", ephemeral.NewReader(), ephemeral.NewReader())
	}

	// jsonl
	if true {
		fmt.Println()
		doTest("ephemeral -> jsonl", ephemeral.NewReader(), jsonl.NewWriter())
		doTest("jsonl -> ephemeral", jsonl.NewReader(), ephemeral.NewWriter())
		verify("ephemeral<->jsonl", ephemeral.NewReader(), jsonl.NewReader())
	}

	// ipfs
	if true {
		fmt.Println()
		iw := ipfs.NewWriter(sh)
		doTest("ephemeral -> ipfs", ephemeral.NewReader(), iw)
		dirCid := iw.Dw.Dir
		// dirCid := "QmYEZzGXRwzWArokCyEqpJnLrbp3F2WEUY46huWtu6TqL6"
		doTest("ipfs -> ephemeral", ipfs.NewReader(sh, dirCid), ephemeral.NewWriter())
		verify("ephemeral <-> ipfs", ephemeral.NewReader(), ipfs.NewReader(sh, dirCid))
	}

	// postgres
	if true {
		fmt.Println()
		doTest("ephemeral -> postgres", ephemeral.NewReader(), postgres.NewWriter(conn, tableNames[0]))
		doTest("postgres -> ephemeral", postgres.NewReader(conn, tableNames[0]), ephemeral.NewWriter())
		verify("ephemeral <-> postgres", ephemeral.NewReader(), postgres.NewReader(conn, tableNames[0]))
	}

	// mysql
	if true {
		fmt.Println()
		doTest("ephemeral -> mysql", ephemeral.NewReader(), mysql.NewWriter(db, tableNames[0]))
		doTest("mysql -> ephemeral", mysql.NewReader(db, tableNames[0]), ephemeral.NewWriter())
		verify("ephemeral <-> mysql", ephemeral.NewReader(), mysql.NewReader(db, tableNames[0]))
	}

	//  ** Mysql -> Flux
	// 116k/s (~200M entries , SSD, empty or full)
	// pipeToFlux(fromMysql(db))

}

func doTest(name string, r entryReader, w entryWriter) (int, error) {
	log.Printf("-=- %s\n", name)
	return w.Write(progress.Monitor(name, r.Read()))
}

func verify(name string, a, b entryReader) {
	log.Printf("-=- %s\n", name)
	vv := merge.Verify(a.Read(), progress.Monitor(name, b.Read()))
	log.Printf("Verified %s:\n", name)
	for _, v := range vv {
		log.Println(v)
	}
}

func gaps(in <-chan []types.Entry) {
	ephemeral.NewWriter().Write(progress.Gaps(in))
}
