package bible

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type RequestType int

const (
	RANGE RequestType = iota
	COLLECTION
	MIXED
)

type Request interface {
	Type() RequestType
}

type RangeRequest struct {
	Start referance
	End   referance
}

func (r RangeRequest) Type() RequestType {
	return RANGE
}

type CollectionRequest struct {
	Entries []referance
}

func (r CollectionRequest) Type() RequestType {
	return COLLECTION
}

type MixedRequest struct {
	Entries []Request
}

func (r MixedRequest) Type() RequestType {
	return MIXED
}

type referance struct {
	book    string
	chapter float64
	verse   float64
}

func (r referance) Book() string     { return r.book }
func (r referance) Chapter() float64 { return r.chapter }
func (r referance) Verses() float64  { return r.verse }

func Parse(s string) (Request, error) {

	switch readRequestType(s) {
	case COLLECTION:
		return parseCollectionRequest(s)
	case RANGE:
		return parseRangeRequest(s)
	case MIXED:
		return parseMixedRequest(s)
	}

	return nil, nil
}

// parses range expression will split on dash and parse
// two sides of the request separately
func parseRangeRequest(s string) (RangeRequest, error) {
	var r RangeRequest

	left, right, _ := strings.Cut(s, "-")

	start := parseGenericRequest(left)
	end := parseGenericRequest(right)

	if end.book == "" {
		end.book = start.book
	}

	if start.verse > 0 && end.verse == 0 {
		end.verse = end.chapter
		end.chapter = start.chapter
	}

	if end.chapter < start.chapter {
		msg := "END chapter cannot be smaller then START"
		err := fmt.Sprintf("ERROR: %s. start: %#v end: %#v", msg, start, end)
		return r, errors.New(err)
	}

	if start.chapter == end.chapter && start.verse > end.verse {
		msg := "END verse cannot be smaller then START verse within the same chapter"
		err := fmt.Sprintf("ERROR: %s. start: %#v end: %#v", msg, start, end)
		return r, errors.New(err)
	}

	r.Start = start
	r.End = end

	return r, nil
}

func parseGenericRequest(s string) referance {
	//typical case is (NUMBER?)NAME CHAPTER:VERSE
	// but NAME -- we will ignore this case for now
	// or CHPATER:VERSE
	// or CHAPTER?VERSE
	var r referance

	if isName(s) {
		r.book, s = parseName(s)
	}

	left, right, found := strings.Cut(s, ":")

	leftNum, _ := readNumber(left)

	r.chapter = float64(leftNum)

	if !found {
		return r
	}

	rightNum, _ := readNumber(right)

	r.verse = float64(rightNum)

	return r
}

// will check if the beginning of the string contains
// book name in the form of ?number string
func isName(s string) bool {
	_, s = readNumber(s)
	name, s := readString(s)

	if name == "" {
		return false
	}

	return true
}

func parseName(s string) (string, string) {
	num, s := readNumber(s)
	name, s := readString(s)

	if num > 0 {
		name = fmt.Sprintf("%d %s", num, name)
	}

	return name, s
}

func peek(s string) byte {
	if len(s) < 1 {
		return 0
	}
	return s[0]
}

// collection request can contain single verse, collection of verses,
// collection of ranges and etc
func parseCollectionRequest(s string) (CollectionRequest, error) {
	//John 3:14,18,20

	parts := strings.Split(s, ",")

	var refs = make([]referance, len(parts))

	for i, entry := range parts {
		refs[i] = parseGenericRequest(entry)
		if i == 0 {
			continue
		}

		prev := refs[i-1]

		if refs[i].book == "" {
			refs[i].book = prev.book
		}

		if (refs[i].chapter != 0 && refs[i].verse == 0) &&
			(prev.chapter != 0 && prev.verse != 0) {
			refs[i].verse = refs[i].chapter
			refs[i].chapter = prev.chapter

		}
	}

	return CollectionRequest{
		Entries: refs,
	}, nil
}

func parseMixedRequest(s string) (MixedRequest, error) { return MixedRequest{}, nil }

func readRequestType(s string) RequestType {
	commas := strings.Count(s, ",")
	dashes := strings.Count(s, "-")

	if commas > 0 && dashes > 0 {
		return MIXED
	}

	if dashes > 0 {
		return RANGE
	}

	return COLLECTION
}

func readString(s string) (string, string) {
	var str string

	s = skipWhitespace(s)

	for i, r := range s {
		if !unicode.IsLetter(r) && r != rune('.') {
			str = s[:i]
			s = s[i:]
			break
		}

	}

	return str, s

}

func skipWhitespace(s string) string {

	for i, char := range s {
		if char != rune(' ') {
			s = s[i:]
			break
		}
	}

	return s

}

// will keep reading string until first non digit char
// will return nubmer and rest or the string. if number conversion fails
// will return -1
func readNumber(s string) (int, string) {
	var numRunes []rune
	var number int

	s = skipWhitespace(s)

	for i, r := range s {

		if !unicode.IsNumber(r) {
			s = s[i:]
			break
		}

		numRunes = append(numRunes, r)

	}

	if len(numRunes) < 1 {
		return number, s
	}

	number, err := strconv.Atoi(string(numRunes))

	if err != nil {
		number = -1
	}

	return number, s
}
