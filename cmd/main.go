package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ButbkaDrug/bible"
	_ "modernc.org/sqlite"
)

func main() {

	HOME, err := os.UserHomeDir()
	BIBLE_DIR := filepath.Join(HOME, ".config", "bible-cli")
	TRANSLATION := "ESV"
	EXT := "SQLite3"

	if err != nil {
		log.Fatal("failed to find home directory!")
	}

	if env_path := os.Getenv("BIBLECLI"); env_path != "" {
		BIBLE_DIR = env_path
	}

	if env_translation := os.Getenv("TRANSLATION"); env_translation != "" {
		TRANSLATION = env_translation
	}

	DATABASE := filepath.Join(BIBLE_DIR, fmt.Sprintf("%s.%s", TRANSLATION, EXT))

	conn, err := sql.Open("sqlite", DATABASE)
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
