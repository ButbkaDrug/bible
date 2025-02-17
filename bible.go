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
	ctx            context.Context
	db             *repository.Queries
	render         Renderer
	writer         io.Writer
	books          []repository.BooksAll
	defBookNumbers []repository.BooksAll

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

	defBooksNumbers := []repository.BooksAll{
		{
			BookNumber: 10,

			ShortName: "Gen",
			LongName:  "Genesis",
			BookColor: "#ccccff",
			IsPresent: true,
		},
		{
			BookNumber: 20,
			ShortName:  "Exo",
			LongName:   "Exodus",
			BookColor:  "#ccccff",
			IsPresent:  true,
		},
		{
			BookNumber: 30,
			ShortName:  "Lev",
			LongName:   "Leviticus",
			BookColor:  "#ccccff",
			IsPresent:  true,
		},
		{
			BookNumber: 40,
			ShortName:  "Num",
			LongName:   "Numbers",
			BookColor:  "#ccccff",
			IsPresent:  true,
		},
		{
			BookNumber: 50,
			ShortName:  "Deu",
			LongName:   "Deuteronomy",
			BookColor:  "#ccccff",
			IsPresent:  true,
		},
		{
			BookNumber: 60,
			ShortName:  "Josh",
			LongName:   "Joshua",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 70,
			ShortName:  "Judg",
			LongName:   "Judges",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 80,
			ShortName:  "Ruth",
			LongName:   "Ruth",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 90,
			ShortName:  "1Sam",
			LongName:   "1 Samuel",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 100,
			ShortName:  "2Sam",
			LongName:   "2 Samuel",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 110,
			ShortName:  "1Kgs",
			LongName:   "1 Kings",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 120,
			ShortName:  "2Kgs",
			LongName:   "2 Kings",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 130,
			ShortName:  "1Chr",
			LongName:   "1 Chronicles",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 140,
			ShortName:  "2Chr",
			LongName:   "2 Chronicles",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 150,
			ShortName:  "Ezr",
			LongName:   "Ezra",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 160,
			ShortName:  "Neh",
			LongName:   "Nehemiah",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 165,
			ShortName:  "1Esd",
			LongName:   "1 Esdras",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 170,
			ShortName:  "Tob",
			LongName:   "Tobit",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 180,
			ShortName:  "Jdt",
			LongName:   "Judith",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 190,
			ShortName:  "Esth",
			LongName:   "Esther",
			BookColor:  "#ffcc99",
			IsPresent:  true,
		},
		{
			BookNumber: 192,
			ShortName:  "EstGr",
			LongName:   "Greek Esther",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 220,
			ShortName:  "Job",
			LongName:   "Job",
			BookColor:  "#66ff99",
			IsPresent:  true,
		},
		{
			BookNumber: 230,
			ShortName:  "Ps",
			LongName:   "Psalm",
			BookColor:  "#66ff99",
			IsPresent:  true,
		},
		{
			BookNumber: 232,
			ShortName:  "Ps151",
			LongName:   "Psalm 151",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 240,
			ShortName:  "Prov",
			LongName:   "Proverbs",
			BookColor:  "#66ff99",
			IsPresent:  true,
		},
		{
			BookNumber: 250,
			ShortName:  "Eccl",
			LongName:   "Ecclesiastes",
			BookColor:  "#66ff99",
			IsPresent:  true,
		},
		{
			BookNumber: 260,
			ShortName:  "Song",
			LongName:   "Song of Solomon",
			BookColor:  "#66ff99",
			IsPresent:  true,
		},
		{
			BookNumber: 270,
			ShortName:  "Wis",
			LongName:   "Wisdom",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 280,
			ShortName:  "Sir",
			LongName:   "Sirach",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 290,
			ShortName:  "Isa",
			LongName:   "Isaiah",
			BookColor:  "#ff9fb4",
			IsPresent:  true,
		},
		{
			BookNumber: 300,
			ShortName:  "Jer",
			LongName:   "Jeremiah",
			BookColor:  "#ff9fb4",
			IsPresent:  true,
		},
		{
			BookNumber: 305,
			ShortName:  "PrAz",
			LongName:   "Prayer of Azariah",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 310,
			ShortName:  "Lam",
			LongName:   "Lamentations",
			BookColor:  "#ff9fb4",
			IsPresent:  true,
		},
		{
			BookNumber: 315,
			ShortName:  "EpJer",
			LongName:   "Letter of Jeremiah",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 320,
			ShortName:  "Bar",
			LongName:   "Baruch",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 323,
			ShortName:  "Sg3",
			LongName:   "Song of the Three Young Men",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 325,
			ShortName:  "Sus",
			LongName:   "Susanna",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 330,
			ShortName:  "Ezek",
			LongName:   "Ezekiel",
			BookColor:  "#ff9fb4",
			IsPresent:  true,
		},
		{
			BookNumber: 340,
			ShortName:  "Dan",
			LongName:   "Daniel",
			BookColor:  "#ff9fb4",
			IsPresent:  true,
		},
		{
			BookNumber: 345,
			ShortName:  "Bel",
			LongName:   "Bel and the Dragon",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 350,
			ShortName:  "Hos",
			LongName:   "Hosea",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 360,
			ShortName:  "Joel",
			LongName:   "Joel",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 370,
			ShortName:  "Am",
			LongName:   "Amos",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 380,
			ShortName:  "Oba",
			LongName:   "Obadiah",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 390,
			ShortName:  "Jona",
			LongName:   "Jonah",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 400,
			ShortName:  "Mic",
			LongName:   "Micah",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 410,
			ShortName:  "Nah",
			LongName:   "Nahum",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 420,
			ShortName:  "Hab",
			LongName:   "Habakkuk",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 430,
			ShortName:  "Zeph",
			LongName:   "Zephaniah",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 440,
			ShortName:  "Hag",
			LongName:   "Haggai",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 450,
			ShortName:  "Zech",
			LongName:   "Zechariah",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 460,
			ShortName:  "Mal",
			LongName:   "Malachi",
			BookColor:  "#ffff99",
			IsPresent:  true,
		},
		{
			BookNumber: 462,
			ShortName:  "1Mac",
			LongName:   "1 Maccabees",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 464,
			ShortName:  "2Mac",
			LongName:   "2 Maccabees",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 466,
			ShortName:  "3Mac",
			LongName:   "3 Maccabees",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 467,
			ShortName:  "4Mac",
			LongName:   "4 Maccabees",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 468,
			ShortName:  "2Esd",
			LongName:   "2 Esdras",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
		{
			BookNumber: 470,
			ShortName:  "Mat",
			LongName:   "Matthew",
			BookColor:  "#ff6600",
			IsPresent:  true,
		},
		{
			BookNumber: 480,
			ShortName:  "Mar",
			LongName:   "Mark",
			BookColor:  "#ff6600",
			IsPresent:  true,
		},
		{
			BookNumber: 490,
			ShortName:  "Luk",
			LongName:   "Luke",
			BookColor:  "#ff6600",
			IsPresent:  true,
		},
		{
			BookNumber: 500,
			ShortName:  "John",
			LongName:   "John",
			BookColor:  "#ff6600",
			IsPresent:  true,
		},
		{
			BookNumber: 510,
			ShortName:  "Acts",
			LongName:   "Acts",
			BookColor:  "#00ffff",
			IsPresent:  true,
		},
		{
			BookNumber: 520,
			ShortName:  "Rom",
			LongName:   "Romans",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 530,
			ShortName:  "1Cor",
			LongName:   "1 Corinthians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 540,
			ShortName:  "2Cor",
			LongName:   "2 Corinthians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 550,
			ShortName:  "Gal",
			LongName:   "Galatians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 560,
			ShortName:  "Eph",
			LongName:   "Ephesians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 570,
			ShortName:  "Phil",
			LongName:   "Philippians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 580,
			ShortName:  "Col",
			LongName:   "Colossians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 590,
			ShortName:  "1Ths",
			LongName:   "1 Thessalonians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 600,
			ShortName:  "2Ths",
			LongName:   "2 Thessalonians",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 610,
			ShortName:  "1Tim",
			LongName:   "1 Timothy",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 620,
			ShortName:  "2Tim",
			LongName:   "2 Timothy",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 630,
			ShortName:  "Tit",
			LongName:   "Titus",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 640,
			ShortName:  "Phlm",
			LongName:   "Philemon",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 650,
			ShortName:  "Heb",
			LongName:   "Hebrews",
			BookColor:  "#ffff00",
			IsPresent:  true,
		},
		{
			BookNumber: 660,
			ShortName:  "Jam",
			LongName:   "James",
			BookColor:  "#00ff00",
			IsPresent:  true,
		},
		{
			BookNumber: 670,
			ShortName:  "1Pet",
			LongName:   "1 Peter",
			BookColor:  "#00ff00",
			IsPresent:  true,
		},
		{
			BookNumber: 680,
			ShortName:  "2Pet",
			LongName:   "2 Peter",
			BookColor:  "#00ff00",
			IsPresent:  true,
		},
		{
			BookNumber: 690,
			ShortName:  "1Jn",
			LongName:   "1 John",
			BookColor:  "#00ff00",
			IsPresent:  true,
		},
		{
			BookNumber: 700,
			ShortName:  "2Jn",
			LongName:   "2 John",
			BookColor:  "#00ff00",
			IsPresent:  true,
		},
		{
			BookNumber: 710,
			ShortName:  "3Jn",
			LongName:   "3 John",
			BookColor:  "#00ff00",
			IsPresent:  true,
		},
		{
			BookNumber: 720,
			ShortName:  "Jud",
			LongName:   "Jude",
			BookColor:  "#00ff00",
			IsPresent:  true,
		},
		{
			BookNumber: 730,
			ShortName:  "Rev",
			LongName:   "Revelation",
			BookColor:  "#ff7c80",
			IsPresent:  true,
		},
		{
			BookNumber: 790,
			ShortName:  "PrMan",
			LongName:   "Prayer of Manasseh",
			BookColor:  "#c0c0c0",
			IsPresent:  false,
		},
	}

	app.books = books
	app.defBookNumbers = defBooksNumbers

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
	for _, book := range app.books {
		if strings.ToLower(s) == strings.ToLower(book.LongName) {
			return float64(book.BookNumber)
		}

		if strings.ToLower(s) == strings.ToLower(book.ShortName) {
			return float64(book.BookNumber)
		}
	}

	for _, book := range app.defBookNumbers {
		if strings.ToLower(s) == strings.ToLower(book.LongName) {
			return float64(book.BookNumber)
		}

		if strings.ToLower(s) == strings.ToLower(book.ShortName) {
			return float64(book.BookNumber)
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
