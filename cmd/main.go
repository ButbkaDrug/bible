package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ButbkaDrug/bible"
	_ "modernc.org/sqlite"
)

func main() {

	conn, err := sql.Open("sqlite", os.Getenv("PATH_TO_BIBLE"))
	if err != nil {
		log.Fatalf("database connection error: %s", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var env = os.Getenv("BIBLE_ENV")

	var query string
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	}

	app := bible.New(ctx, conn).
		SetQuery(query).
		SetEnvironment(env)

	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
