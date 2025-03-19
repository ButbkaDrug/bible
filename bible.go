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
	Book       string
	Text       string
	BookNumber int
	Chapter    int
	Verse      int
}

func wrapBooks(books []repository.Book) []Verse {
	var result = make([]Verse, len(books))

	for i, b := range books {
		result[i].Book = b.LongName
		result[i].Text = b.ShortName

		result[i].BookNumber = int(b.BookNumber)
		result[i].Chapter = 0
		result[i].Verse = 0
	}

	return result
}

func wrapVerses(book string, verses []repository.Verse) []Verse {
	var result = make([]Verse, len(verses))

	for i, v := range verses {
		result[i].Book = book
		result[i].Text = v.Text

		result[i].BookNumber = int(v.BookNumber)
		result[i].Chapter = int(v.Chapter)
		result[i].Verse = int(v.Verse)

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

type Bible struct {
	ctx            context.Context
	db             *repository.Queries
	render         Renderer
	writer         io.Writer
	books          []repository.Book
	defBookNumbers []repository.Book

	query string
	env   string
}

func New(ctx context.Context, conn repository.DBTX, env string) *Bible {
	bible := &Bible{
		ctx:    ctx,
		db:     repository.New(conn),
		env:    env,
		writer: os.Stdout,
	}

	bible.init()

	return bible
}

func (app *Bible) init() *Bible {
	if app.ctx == nil {
		app.ctx = context.Background()
	}

	if app.render == nil {
		render := NewDefaultRender()
		if app.env == "" {
			render.Color()
		}
		app.render = render
	}

	books, err := app.GetBooks()
	if err != nil {
		log.Fatal("initialization failed: ", err)
	}

	defBooksNumbers := []repository.Book{
		{
			BookNumber: 10,

			ShortName: "Gen",
			LongName:  "Genesis",
			BookColor: "#ccccff",
		},
		{
			BookNumber: 20,
			ShortName:  "Exo",
			LongName:   "Exodus",
			BookColor:  "#ccccff",
		},
		{
			BookNumber: 30,
			ShortName:  "Lev",
			LongName:   "Leviticus",
			BookColor:  "#ccccff",
		},
		{
			BookNumber: 40,
			ShortName:  "Num",
			LongName:   "Numbers",
			BookColor:  "#ccccff",
		},
		{
			BookNumber: 50,
			ShortName:  "Deu",
			LongName:   "Deuteronomy",
			BookColor:  "#ccccff",
		},
		{
			BookNumber: 60,
			ShortName:  "Josh",
			LongName:   "Joshua",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 70,
			ShortName:  "Judg",
			LongName:   "Judges",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 80,
			ShortName:  "Ruth",
			LongName:   "Ruth",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 90,
			ShortName:  "1Sam",
			LongName:   "1 Samuel",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 100,
			ShortName:  "2Sam",
			LongName:   "2 Samuel",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 110,
			ShortName:  "1Kgs",
			LongName:   "1 Kings",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 120,
			ShortName:  "2Kgs",
			LongName:   "2 Kings",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 130,
			ShortName:  "1Chr",
			LongName:   "1 Chronicles",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 140,
			ShortName:  "2Chr",
			LongName:   "2 Chronicles",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 150,
			ShortName:  "Ezr",
			LongName:   "Ezra",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 160,
			ShortName:  "Neh",
			LongName:   "Nehemiah",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 165,
			ShortName:  "1Esd",
			LongName:   "1 Esdras",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 170,
			ShortName:  "Tob",
			LongName:   "Tobit",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 180,
			ShortName:  "Jdt",
			LongName:   "Judith",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 190,
			ShortName:  "Esth",
			LongName:   "Esther",
			BookColor:  "#ffcc99",
		},
		{
			BookNumber: 192,
			ShortName:  "EstGr",
			LongName:   "Greek Esther",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 220,
			ShortName:  "Job",
			LongName:   "Job",
			BookColor:  "#66ff99",
		},
		{
			BookNumber: 230,
			ShortName:  "Ps",
			LongName:   "Psalm",
			BookColor:  "#66ff99",
		},
		{
			BookNumber: 232,
			ShortName:  "Ps151",
			LongName:   "Psalm 151",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 240,
			ShortName:  "Prov",
			LongName:   "Proverbs",
			BookColor:  "#66ff99",
		},
		{
			BookNumber: 250,
			ShortName:  "Eccl",
			LongName:   "Ecclesiastes",
			BookColor:  "#66ff99",
		},
		{
			BookNumber: 260,
			ShortName:  "Song",
			LongName:   "Song of Solomon",
			BookColor:  "#66ff99",
		},
		{
			BookNumber: 270,
			ShortName:  "Wis",
			LongName:   "Wisdom",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 280,
			ShortName:  "Sir",
			LongName:   "Sirach",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 290,
			ShortName:  "Isa",
			LongName:   "Isaiah",
			BookColor:  "#ff9fb4",
		},
		{
			BookNumber: 300,
			ShortName:  "Jer",
			LongName:   "Jeremiah",
			BookColor:  "#ff9fb4",
		},
		{
			BookNumber: 305,
			ShortName:  "PrAz",
			LongName:   "Prayer of Azariah",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 310,
			ShortName:  "Lam",
			LongName:   "Lamentations",
			BookColor:  "#ff9fb4",
		},
		{
			BookNumber: 315,
			ShortName:  "EpJer",
			LongName:   "Letter of Jeremiah",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 320,
			ShortName:  "Bar",
			LongName:   "Baruch",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 323,
			ShortName:  "Sg3",
			LongName:   "Song of the Three Young Men",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 325,
			ShortName:  "Sus",
			LongName:   "Susanna",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 330,
			ShortName:  "Ezek",
			LongName:   "Ezekiel",
			BookColor:  "#ff9fb4",
		},
		{
			BookNumber: 340,
			ShortName:  "Dan",
			LongName:   "Daniel",
			BookColor:  "#ff9fb4",
		},
		{
			BookNumber: 345,
			ShortName:  "Bel",
			LongName:   "Bel and the Dragon",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 350,
			ShortName:  "Hos",
			LongName:   "Hosea",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 360,
			ShortName:  "Joel",
			LongName:   "Joel",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 370,
			ShortName:  "Am",
			LongName:   "Amos",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 380,
			ShortName:  "Oba",
			LongName:   "Obadiah",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 390,
			ShortName:  "Jona",
			LongName:   "Jonah",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 400,
			ShortName:  "Mic",
			LongName:   "Micah",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 410,
			ShortName:  "Nah",
			LongName:   "Nahum",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 420,
			ShortName:  "Hab",
			LongName:   "Habakkuk",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 430,
			ShortName:  "Zeph",
			LongName:   "Zephaniah",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 440,
			ShortName:  "Hag",
			LongName:   "Haggai",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 450,
			ShortName:  "Zech",
			LongName:   "Zechariah",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 460,
			ShortName:  "Mal",
			LongName:   "Malachi",
			BookColor:  "#ffff99",
		},
		{
			BookNumber: 462,
			ShortName:  "1Mac",
			LongName:   "1 Maccabees",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 464,
			ShortName:  "2Mac",
			LongName:   "2 Maccabees",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 466,
			ShortName:  "3Mac",
			LongName:   "3 Maccabees",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 467,
			ShortName:  "4Mac",
			LongName:   "4 Maccabees",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 468,
			ShortName:  "2Esd",
			LongName:   "2 Esdras",
			BookColor:  "#c0c0c0",
		},
		{
			BookNumber: 470,
			ShortName:  "Mat",
			LongName:   "Matthew",
			BookColor:  "#ff6600",
		},
		{
			BookNumber: 480,
			ShortName:  "Mar",
			LongName:   "Mark",
			BookColor:  "#ff6600",
		},
		{
			BookNumber: 490,
			ShortName:  "Luk",
			LongName:   "Luke",
			BookColor:  "#ff6600",
		},
		{
			BookNumber: 500,
			ShortName:  "John",
			LongName:   "John",
			BookColor:  "#ff6600",
		},
		{
			BookNumber: 510,
			ShortName:  "Acts",
			LongName:   "Acts",
			BookColor:  "#00ffff",
		},
		{
			BookNumber: 520,
			ShortName:  "Rom",
			LongName:   "Romans",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 530,
			ShortName:  "1Cor",
			LongName:   "1 Corinthians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 540,
			ShortName:  "2Cor",
			LongName:   "2 Corinthians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 550,
			ShortName:  "Gal",
			LongName:   "Galatians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 560,
			ShortName:  "Eph",
			LongName:   "Ephesians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 570,
			ShortName:  "Phil",
			LongName:   "Philippians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 580,
			ShortName:  "Col",
			LongName:   "Colossians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 590,
			ShortName:  "1Ths",
			LongName:   "1 Thessalonians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 600,
			ShortName:  "2Ths",
			LongName:   "2 Thessalonians",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 610,
			ShortName:  "1Tim",
			LongName:   "1 Timothy",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 620,
			ShortName:  "2Tim",
			LongName:   "2 Timothy",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 630,
			ShortName:  "Tit",
			LongName:   "Titus",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 640,
			ShortName:  "Phlm",
			LongName:   "Philemon",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 650,
			ShortName:  "Heb",
			LongName:   "Hebrews",
			BookColor:  "#ffff00",
		},
		{
			BookNumber: 660,
			ShortName:  "Jam",
			LongName:   "James",
			BookColor:  "#00ff00",
		},
		{
			BookNumber: 670,
			ShortName:  "1Pet",
			LongName:   "1 Peter",
			BookColor:  "#00ff00",
		},
		{
			BookNumber: 680,
			ShortName:  "2Pet",
			LongName:   "2 Peter",
			BookColor:  "#00ff00",
		},
		{
			BookNumber: 690,
			ShortName:  "1Jn",
			LongName:   "1 John",
			BookColor:  "#00ff00",
		},
		{
			BookNumber: 700,
			ShortName:  "2Jn",
			LongName:   "2 John",
			BookColor:  "#00ff00",
		},
		{
			BookNumber: 710,
			ShortName:  "3Jn",
			LongName:   "3 John",
			BookColor:  "#00ff00",
		},
		{
			BookNumber: 720,
			ShortName:  "Jud",
			LongName:   "Jude",
			BookColor:  "#00ff00",
		},
		{
			BookNumber: 730,
			ShortName:  "Rev",
			LongName:   "Revelation",
			BookColor:  "#ff7c80",
		},
		{
			BookNumber: 790,
			ShortName:  "PrMan",
			LongName:   "Prayer of Manasseh",
			BookColor:  "#c0c0c0",
		},
	}

	app.books = books
	app.defBookNumbers = defBooksNumbers

	return app
}

func (app *Bible) SetEnvironment(s string) *Bible {
	app.env = s
	return app

}

func (app *Bible) GetBooks() ([]repository.Book, error) {
	return app.db.GetBookNames(app.ctx)
}

func (app *Bible) GetChapters(n int) ([]Verse, error) {
	verses, err := app.db.GetChapters(app.ctx, float64(n))
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

func (app *Bible) Search(s string) ([]Verse, error) {
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

func (app *Bible) getBookName(num float64) string {
	for _, book := range app.books {
		if book.BookNumber == num {
			return book.LongName
		}
	}

	return fmt.Sprintf("undefined(%v)", num)
}

func (app *Bible) getBookNumber(s string) float64 {
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

func (app *Bible) GetVersesRange(r RangeRequest) ([]Verse, error) {
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

func (app *Bible) GetChapter(book int, chapter int) ([]Verse, error) {
	param := repository.GetVersesRangeParams{
		BookNumber: float64(book),
		Chapter:    float64(chapter),
		FromVerse:  0,
		ToVerse:    MAX_VERSE,
	}
	v, err := app.db.GetVersesRange(app.ctx, param)

	if err != nil {
		return []Verse{}, err
	}

	//no book name for now
	return wrapVerses("", v), nil
}

func (app *Bible) requestRange(name string, chapter, from, to float64) ([]repository.Verse, error) {
	bookNumber := app.getBookNumber(name)

	param := repository.GetVersesRangeParams{
		BookNumber: bookNumber,
		Chapter:    chapter,
		FromVerse:  from,
		ToVerse:    to,
	}
	return app.db.GetVersesRange(app.ctx, param)
}
func (app *Bible) GetVersesCollection(r CollectionRequest) ([]Verse, error) {
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

func (app *Bible) requestCollection(r Referance) ([]repository.Verse, error) {
	bookNumber := app.getBookNumber(r.Book())

	params := repository.GetVersesCollectionParams{
		BookNumber: bookNumber,
		Chapter:    r.Chapter(),
		Verse:      r.Verse(),
	}
	return app.db.GetVersesCollection(app.ctx, params)
}

func (app *Bible) SetRender(r Renderer) *Bible {
	app.render = r
	return app
}

func (app *Bible) SetQuery(s string) *Bible {
	app.query = s
	return app
}

func (app *Bible) SetContext(ctx context.Context) *Bible {
	app.ctx = ctx
	return app
}

func (app *Bible) SetDBConnection(conn repository.DBTX) *Bible {
	app.db = repository.New(conn)
	return app
}

func (app *Bible) SetWriter(w io.Writer) *Bible {
	app.writer = w
	return app
}

func (app *Bible) Execute() ([]Verse, error) {
	request, err := Parse(app.query)
	if err != nil {
		return []Verse{}, err
	}

	var verses []Verse
	switch r := request.(type) {
	case EmptyRequest:
		// I want to return list of books
		books, err := app.GetBooks()

		if err != nil {
			return []Verse{}, err
		}

		return wrapBooks(books), nil
	case ConcreteRequest:
		bookNumber := app.getBookNumber(r.ref.book)
		if bookNumber == 0 {
			break
		}
		return app.GetChapters(int(bookNumber))
	case RangeRequest:
		return app.GetVersesRange(r)
	case CollectionRequest:
		return app.GetVersesCollection(r)
	case MixedRequest:
		return []Verse{}, errors.New("MIXED REQUESTS ARE NOT IMPLEMENTED")
	}

	if len(verses) < 1 {
		verses, err = app.Search(app.query)

		if err != nil {
			return []Verse{}, err
		}

		app.render.SetHighlights(strings.Split(app.query, " "))

		return verses, nil
	}

	if len(verses) < 1 {
		return verses, errors.New("nothing found!")
	}

	return []Verse{}, err

}

func (app *Bible) Run() error {
	verses, err := app.Execute()
	if err != nil {
		return err
	}

	return app.render.Render(app.writer, verses)
}
