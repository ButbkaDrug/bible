package main

import (
	"context"
	"database/sql"
	"log"
	"os"
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

	app := bible.New(ctx, conn)
	if err := app.SetQuery("1 John 3:2").Execute(); err != nil {
		log.Fatal(err)
	}
}
