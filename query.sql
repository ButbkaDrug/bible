-- name: GetBookNames :many
SELECT * FROM books ORDER BY book_number;

-- name: GetChapters :many
SELECT * FROM verses
WHERE (book_number = ?)
AND (verse = 1)
ORDER BY book_number, chapter, verse;

-- name: GetVersesCollection :many
SELECT * FROM verses
WHERE (book_number = ?)
AND (chapter = ?)
AND (verse = ?)
ORDER BY book_number, chapter, verse;

-- name: GetVersesRange :many
SELECT * FROM verses
WHERE (book_number = ?)
AND (chapter = ?)
AND (verse BETWEEN ? and ?)
ORDER BY book_number, chapter, verse;

-- name: Search :many
SELECT
 *
FROM
 verses
WHERE
 text LIKE ?
ORDER BY
 book_number,
 chapter,
 verse;
