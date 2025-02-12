// TODO: I need to fix the problem when i pass verses to renderer
// book name is a number, and renderer have no idea what the string name is
// i need to be able to print book name - it is vital
package bible

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ButbkaDrug/bible/internal/repository"
)

const (
	MAX_VERSE float64 = 999
)

type Verse struct {
	Book    string
	Text    string
	Chapter int
	Verse   int
}

func wrapVerses(book string, verses []repository.Verse) []Verse {
	var result = make([]Verse, len(verses))

	for i, v := range verses {
		result[i].Book = book
		result[i].Chapter = int(v.Chapter)
		result[i].Verse = int(v.Verse)
		result[i].Text = v.Text

	}

	return result
}

type Referance interface {
	Book() string
	Chapter() float64
	Verse() float64
}

type Renderer interface {
	Render(io.Writer, []Verse) error
	SetHighlights([]string) Renderer
}

type app struct {
	ctx    context.Context
	db     *repository.Queries
	render Renderer
	writer io.Writer
	books  []repository.BooksAll

	query string
	env   string
}

func New(ctx context.Context, conn repository.DBTX) *app {
	return &app{
		ctx:    ctx,
		db:     repository.New(conn),
		writer: os.Stdout,
	}
}

func (app *app) init() *app {
	if app.ctx == nil {
		app.ctx = context.Background()
	}

	if app.render == nil {
		render := NewDefaultRender()
		if app.env == "" {
			render = render.Color()
		}

		app.render = render
	}

	books, err := app.getBookNames()
	if err != nil {
		log.Fatal("initialization failed: ", err)
	}

	app.books = books

	return app
}

func (app *app) SetEnvironment(s string) *app {
	app.env = s
	return app

}

func (app *app) getBookNames() ([]repository.BooksAll, error) {

	return app.db.GetBookNames(app.ctx)
}

func (app *app) Search(s string) ([]Verse, error) {
	var query string

	query = strings.Trim(s, " \n\r\t")
	query = strings.ReplaceAll(query, " ", "%")
	query = fmt.Sprintf("%%%s%%", query)

	verses, err := app.db.Search(app.ctx, query)

	if err != nil {
		return []Verse{}, err
	}

	var result []Verse

	for _, v := range verses {
		result = append(result, Verse{
			Book:    app.getBookName(v.BookNumber),
			Chapter: int(v.Chapter),
			Verse:   int(v.Verse),
			Text:    v.Text,
		})

	}
	return result, nil
}

func (app *app) getBookName(num float64) string {

	for _, book := range app.books {
		if book.BookNumber == num {
			return book.LongName
		}
	}

	return fmt.Sprintf("undefined(%v)", num)

}

func (app *app) getBookNumber(s string) float64 {

	for id, book := range app.books {
		if strings.ToLower(s) == strings.ToLower(book.LongName) {
			return float64(id)
		}

		if strings.ToLower(s) == strings.ToLower(book.ShortName) {
			return float64(id)
		}
	}

	return 0
}

func (app *app) GetVersesRange(r RangeRequest) ([]Verse, error) {
	if r.Simple() {
		resp, err := app.requestRange(r.Start.book, r.Start.chapter, r.Start.verse, r.End.verse)
		return wrapVerses(r.Start.book, resp), err
	}

	left, err := app.requestRange(r.Start.book, r.Start.chapter, r.Start.verse, MAX_VERSE)
	if err != nil {
		return []Verse{}, err
	}

	right, err := app.requestRange(r.End.book, r.End.chapter, 1, r.End.verse)
	if err != nil {
		return []Verse{}, err
	}

	return append(wrapVerses(r.Start.book, left), wrapVerses(r.End.book, right)...), nil
}

func (app *app) requestRange(name string, chapter, from, to float64) ([]repository.Verse, error) {
	bookNumber := app.getBookNumber(name)

	param := repository.GetVersesRangeParams{
		BookNumber: bookNumber,
		Chapter:    chapter,
		FromVerse:  from,
		ToVerse:    to,
	}
	return app.db.GetVersesRange(app.ctx, param)
}
func (app *app) GetVersesCollection(r CollectionRequest) ([]Verse, error) {
	var result []Verse

	for _, req := range r.Entries {
		resp, err := app.requestCollection(req)
		if err != nil {
			return []Verse{}, err
		}
		result = append(result, wrapVerses(req.book, resp)...)
	}

	return result, nil
}

func (app *app) requestCollection(r Referance) ([]repository.Verse, error) {
	bookNumber := app.getBookNumber(r.Book())

	params := repository.GetVersesCollectionParams{
		BookNumber: bookNumber,
		Chapter:    r.Chapter(),
		Verse:      r.Verse(),
	}
	return app.db.GetVersesCollection(app.ctx, params)
}

func (app *app) SetRender(r Renderer) *app {
	app.render = r
	return app
}

func (app *app) SetQuery(s string) *app {
	app.query = s
	return app
}

func (app *app) SetContext(ctx context.Context) *app {
	app.ctx = ctx
	return app
}

func (app *app) SetDBConnection(conn repository.DBTX) *app {
	app.db = repository.New(conn)
	return app
}
func (app *app) SetWriter(w io.Writer) *app {
	app.writer = w
	return app
}

func (app *app) Execute() error {
	app.init()
	//check that the query is set if not show help menu
	//parse query
	//=> detect if it's a search or referance
	//peform db lookup
	//render results
	request, err := Parse(app.query)
	if err != nil {
		return err
	}

	var verses []Verse
	switch r := request.(type) {
	case RangeRequest:
		verses, err = app.GetVersesRange(r)
		if err != nil {
			return err
		}
	case CollectionRequest:
		verses, err = app.GetVersesCollection(r)
		if err != nil {
			return err
		}
	case MixedRequest:
		return errors.New("MIXED REQUESTS ARE NOT IMPLEMENTED")
	}

	if len(verses) < 1 {
		verses, err = app.Search(app.query)

		if err != nil {
			return err
		}

		app.render.SetHighlights(strings.Split(app.query, " "))
	}

	if len(verses) < 1 {
		log.Fatal("noting was found! query: ", app.query)
	}

	app.render.Render(app.writer, verses)

	return nil
}
