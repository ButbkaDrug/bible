package bible

import (
	"context"
	"fmt"
	"strings"

	"github.com/ButbkaDrug/bible/internal/repository"
)

type Referance interface {
	Book() float64
	Chapter() float64
	Verses() []float64
	//should return ether "collection" or "range"
	Type() string
}

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

func (app *app) GetVersesRange(r Referance) ([]repository.Verse, error) {

	param := repository.GetVersesRangeParams{
		BookNumber: r.Book(),
		Chapter:    r.Chapter(),
		FromVerse:  r.Verses()[0],
		ToVerse:    r.Verses()[1],
	}
	return app.db.GetVersesRange(app.ctx, param)
}

func (app *app) GetVersesCollection(r Referance) ([]repository.Verse, error) {
	params := repository.GetVersesCollectionParams{
		BookNumber: r.Book(),
		Chapter:    r.Chapter(),
		Number:     r.Verses(),
	}
	return app.db.GetVersesCollection(app.ctx, params)
}
