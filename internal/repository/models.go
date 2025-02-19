// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

type Book struct {
	BookNumber float64
	ShortName  string
	LongName   string
	BookColor  string
}

type BooksAll struct {
	BookNumber float64
	ShortName  string
	LongName   string
	BookColor  string
	IsPresent  bool
}

type Info struct {
	Name  string
	Value string
}

type Introduction struct {
	BookNumber   float64
	Introduction string
}

type Story struct {
	BookNumber     float64
	Chapter        float64
	Verse          float64
	OrderIfSeveral float64
	Title          string
}

type Verse struct {
	BookNumber float64
	Chapter    float64
	Verse      float64
	Text       string
}
