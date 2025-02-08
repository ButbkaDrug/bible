// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package repository

import (
	"context"
)

const getBookNames = `-- name: GetBookNames :many
SELECT book_number, short_name, long_name, book_color, is_present FROM books_all ORDER BY book_number
`

func (q *Queries) GetBookNames(ctx context.Context) ([]BooksAll, error) {
	rows, err := q.db.QueryContext(ctx, getBookNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BooksAll
	for rows.Next() {
		var i BooksAll
		if err := rows.Scan(
			&i.BookNumber,
			&i.ShortName,
			&i.LongName,
			&i.BookColor,
			&i.IsPresent,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVersesCollection = `-- name: GetVersesCollection :many
SELECT book_number, chapter, verse, text FROM verses
WHERE (book_number = ?)
AND (chapter = ?)
AND (verse = ?)
ORDER BY book_number, chapter, verse
`

type GetVersesCollectionParams struct {
	BookNumber float64
	Chapter    float64
	Verse      float64
}

func (q *Queries) GetVersesCollection(ctx context.Context, arg GetVersesCollectionParams) ([]Verse, error) {
	rows, err := q.db.QueryContext(ctx, getVersesCollection, arg.BookNumber, arg.Chapter, arg.Verse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Verse
	for rows.Next() {
		var i Verse
		if err := rows.Scan(
			&i.BookNumber,
			&i.Chapter,
			&i.Verse,
			&i.Text,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVersesRange = `-- name: GetVersesRange :many
SELECT book_number, chapter, verse, text FROM verses
WHERE (book_number = ?)
AND (chapter = ?)
AND (verse BETWEEN ? and ?)
ORDER BY book_number, chapter, verse
`

type GetVersesRangeParams struct {
	BookNumber float64
	Chapter    float64
	FromVerse  float64
	ToVerse    float64
}

func (q *Queries) GetVersesRange(ctx context.Context, arg GetVersesRangeParams) ([]Verse, error) {
	rows, err := q.db.QueryContext(ctx, getVersesRange,
		arg.BookNumber,
		arg.Chapter,
		arg.FromVerse,
		arg.ToVerse,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Verse
	for rows.Next() {
		var i Verse
		if err := rows.Scan(
			&i.BookNumber,
			&i.Chapter,
			&i.Verse,
			&i.Text,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const search = `-- name: Search :many
SELECT
 book_number, chapter, verse, text
FROM
 verses
WHERE
 text LIKE ?
ORDER BY
 book_number,
 chapter,
 verse
`

func (q *Queries) Search(ctx context.Context, text string) ([]Verse, error) {
	rows, err := q.db.QueryContext(ctx, search, text)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Verse
	for rows.Next() {
		var i Verse
		if err := rows.Scan(
			&i.BookNumber,
			&i.Chapter,
			&i.Verse,
			&i.Text,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
