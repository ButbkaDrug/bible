package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/ButbkaDrug/bible/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	conn, err := sql.Open("sqlite", os.Getenv("PATH_TO_BIBLE"))
	if err != nil {
		log.Fatal("database connection error: ", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	db := repository.New(conn)
	books, err := db.GetBookNames(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", books)

}
