package main

import (
	"context"
	"database/sql"
	"fmt"
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

	verses, err := app.Search("love your neighbor")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", verses)

	chapter, err := app.GetVersesRange(470, 19, 1, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("chapter: %#v\n", chapter)

}
