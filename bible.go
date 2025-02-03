package bible

import (
	"context"
	"fmt"
	"strings"

	"github.com/ButbkaDrug/bible/internal/repository"
)

type app struct {
	ctx context.Context
	db  *repository.Queries
}

func New(ctx context.Context, conn repository.DBTX) *app {

	return &app{
		ctx: ctx,
		db:  repository.New(conn),
	}
}

func (app *app) Search(s string) ([]repository.Verse, error) {
	var query string

	query = strings.Trim(s, " \n\r\t")
	query = strings.ReplaceAll(query, " ", "%")
	query = fmt.Sprintf("%%%s%%", query)

	return app.db.Search(app.ctx, query)
}

func (app *app) getBookNumber(s string) (float64, error) {
	books, err := app.db.GetBookNames(app.ctx)
	if err != nil {
		return 0, err
	}

	for _, book := range books {
		if s == book.LongName {
			return book.BookNumber, nil
		}

		if s == book.ShortName {
			return book.BookNumber, nil
		}
	}

	return 0, nil
}

func (app *app) GetVersesRange(book, chapter, from, until float64) ([]repository.Verse, error) {

	param := repository.GetVersesRangeParams{
		BookNumber: book,
		Chapter:    chapter,
		FromVerse:  from,
		ToVerse:    until,
	}
	return app.db.GetVersesRange(app.ctx, param)
}
