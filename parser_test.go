package bible

import (
	"errors"
	"fmt"
	"testing"
)

func TestReadNumber(t *testing.T) {
	tests := []string{" 1 Jonh", "   2 Corinth.", "Matt.", "1Peter", ""}

	expt := []struct {
		number int
		s      string
	}{
		{number: 1, s: " Jonh"},
		{number: 2, s: " Corinth."},
		{number: 0, s: "Matt."},
		{number: 1, s: "Peter"},
		{number: 0, s: ""},
	}

	for i, test := range tests {
		num, s := readNumber(test)
		enum := expt[i].number
		es := expt[i].s

		if num != enum {
			t.Fatalf("Test[%d] failed with wrong number. Expected %d got %d", i, enum, num)
		}

		if len(s) != len(es) {
			t.Fatalf("Test[%d] failed. Ramaining strings are not the same. Expected len=%d got %d", i, len(s), len(es))
		}

		if s != es {
			t.Fatalf("Test[%d] failed. Ramaining strings are not the same. Expected %s got %s", i, s, es)
		}
	}
}

func TestReadString(t *testing.T) {
	tests := []string{" John 3:16", "    Corinth.2:1", "Matt.3", "Peter    5:1   ", "34", " 5:14"}

	expt := []struct {
		str      string
		reminder string
	}{
		{str: "John", reminder: " 3:16"},
		{str: "Corinth.", reminder: "2:1"},
		{str: "Matt.", reminder: "3"},
		{str: "Peter", reminder: "    5:1   "},
		{str: "", reminder: "34"},
		{str: "", reminder: "5:14"},
	}

	for i, test := range tests {
		str, reminder := readString(test)
		eStr := expt[i].str
		eReminder := expt[i].reminder

		if len(str) != len(eStr) {
			t.Fatalf("Test[%d] failed. Strings are not the same. Expected len=%d got %d", i, len(str), len(eStr))
		}

		if len(reminder) != len(eReminder) {
			t.Fatalf("Test[%d] failed. Ramaining strings are not the same. Expected len=%d got %d", i, len(reminder), len(eReminder))
		}

		if str != eStr {
			t.Fatalf("Test[%d] failed. Strings are not the same. Expected %s got %s", i, str, eStr)
		}

		if reminder != eReminder {
			t.Fatalf("Test[%d] failed. Ramaining strings are not the same. Expected %s got %s", i, reminder, eReminder)
		}
	}
}

func TestRequestType(t *testing.T) {
	tests := []string{
		"John 3:15-20",
		"1 John 3,4",
		"Acts 14-15",
		"Acts 14:38",
		"John 3:15-20, 14",
	}

	expectedREsults := []RequestType{
		RANGE,
		COLLECTION,
		RANGE,
		COLLECTION,
		MIXED,
	}

	for i, test := range tests {
		expectedType := expectedREsults[i]
		resultType := readRequestType(test)

		if resultType != expectedType {
			t.Fatalf("TEST[%d] expectedType %v got %v",
				i,
				resultType,
				expectedType,
			)
		}
	}
}

func TestIsName(t *testing.T) {
	tests := []string{
		"1 John 3",
		"John 3",
		"3",
		"20",
		" 15:20",
	}
	expectedResults := []bool{
		true,
		true,
		false,
		false,
		false,
	}

	for i, test := range tests {

		result := isName(test)
		expect := expectedResults[i]

		if result != expect {
			t.Fatalf("Test[%d] failed: %s expected to be %v", i, test, expect)

		}

	}
}

func TestParseGenericRequest(t *testing.T) {
	tests := []string{
		"John 3:14",
		"25",
		"3:15",
		"",
		"     .",
	}
	expectedResults := []referance{
		{
			book:    "John",
			chapter: 3,
			verse:   14,
		},
		{
			book:    "",
			chapter: 25,
			verse:   0,
		},
		{
			book:    "",
			chapter: 3,
			verse:   15,
		},
		{
			book:    "",
			chapter: 0,
			verse:   0,
		},
		{
			book:    "",
			chapter: 0,
			verse:   0,
		},
	}

	for i, test := range tests {
		result := parseGenericRequest(test)
		expect := expectedResults[i]

		if result.book != expect.book {
			t.Fatalf("TEST[%d failed: expected book=%s got %s]",
				i,
				expect.book,
				result.book,
			)
		}

		if result.chapter != expect.chapter {
			t.Fatalf("TEST[%d failed: expected chapter=%v got %v]",
				i,
				expect.chapter,
				result.chapter,
			)
		}

		if result.verse != expect.verse {
			t.Fatalf("TEST[%d failed: expected verse=%v got %v]",
				i,
				expect.verse,
				result.verse,
			)
		}
	}

}

func TestParseRangeRequest(t *testing.T) {
	tests := []string{
		"Acts 14-15",       // basic range with no verses specified
		"John 3:15-20",     // range within the same chapter
		"1 John 3:15-20",   // range in a book with a number
		"Acts 14:38-15:14", // range across chapters
		"Acts 15-14",       // malformed range
		"John A:B-C",       // non-numeric chapter/verse
		"1 John 1:1-3",     // book with multiple words
		"John 3:20-3:10",   // descending range (end before start)
	}

	expectedResults := []RangeRequest{
		{
			Start: referance{
				book:    "Acts",
				chapter: 14,
				verse:   0,
			},
			End: referance{
				book:    "Acts",
				chapter: 15,
				verse:   0,
			},
		},
		{
			Start: referance{
				book:    "John",
				chapter: 3,
				verse:   15,
			},
			End: referance{
				book:    "John",
				chapter: 3,
				verse:   20,
			},
		},
		{
			Start: referance{
				book:    "1 John",
				chapter: 3,
				verse:   15,
			},
			End: referance{
				book:    "1 John",
				chapter: 3,
				verse:   20,
			},
		},
		{
			Start: referance{
				book:    "Acts",
				chapter: 14,
				verse:   38,
			},
			End: referance{
				book:    "Acts",
				chapter: 15,
				verse:   14,
			},
		},
		{
			// Malformed range, should ideally fail or return an error
			Start: referance{
				book:    "",
				chapter: 0,
				verse:   0,
			},
			End: referance{
				book:    "",
				chapter: 0,
				verse:   0,
			},
		},
		{
			// Invalid chapter/verse (non-numeric), should fail or return a specific error
			Start: referance{
				book:    "John",
				chapter: 0,
				verse:   0,
			},
			End: referance{
				book:    "John",
				chapter: 0,
				verse:   0,
			},
		},
		{
			// Valid multi-word book name (1 John)
			Start: referance{
				book:    "1 John",
				chapter: 1,
				verse:   1,
			},
			End: referance{
				book:    "1 John",
				chapter: 1,
				verse:   3,
			},
		},
		{
			// Descending range, should handle it gracefully
			Start: referance{
				book:    "",
				chapter: 0,
				verse:   0,
			},
			End: referance{
				book:    "",
				chapter: 0,
				verse:   0,
			},
		},
	}

	for i, test := range tests {

		result, err := parseRangeRequest(test)
		expect := expectedResults[i]

		if (i == 4 || i == 7) && err == nil {
			t.Fatalf("TEST[%d] should fail: %s", i, err)
		}

		err = testRangeRequest(expect, result)

		if err != nil {
			t.Fatalf("TEST[%d] failed: %s", i, err)
		}
	}
}

func testRangeRequest(expect, result RangeRequest) error {

	if result.Start.book != expect.Start.book {
		return errors.New(fmt.Sprintf("start book do not match %s != %s",
			expect.Start.book,
			result.Start.book,
		))
	}
	if result.Start.chapter != expect.Start.chapter {
		return errors.New(fmt.Sprintf("start chapter do not match %v != %v",
			expect.Start.chapter,
			result.Start.chapter,
		))
	}
	if result.Start.verse != expect.Start.verse {
		return errors.New(fmt.Sprintf("start verse do not match %v != %v",
			expect.Start.verse,
			result.Start.verse,
		))
	}
	if result.End.book != expect.End.book {
		return errors.New(fmt.Sprintf("start book do not match %s != %s",
			expect.End.book,
			result.End.book,
		))
	}
	if result.End.chapter != expect.End.chapter {
		return errors.New(fmt.Sprintf("start chapter do not match %v != %v",
			expect.End.chapter,
			result.End.chapter,
		))
	}
	if result.End.verse != expect.End.verse {
		return errors.New(fmt.Sprintf("start verse do not match %v != %v",
			expect.End.verse,
			result.End.verse,
		))
	}
	return nil
}

func TestParseCollectionRequest(t *testing.T) {
	tests := []string{
		"John3:14",
		"John3:14,16,18",
		"1John 3,4",
		"1John 3, Matt.4,5",
	}

	expectedResults := []CollectionRequest{
		{
			Entries: []referance{
				{
					book:    "John",
					chapter: 3,
					verse:   14,
				},
			},
		},
		{
			Entries: []referance{
				{
					book:    "John",
					chapter: 3,
					verse:   14,
				},
				{
					book:    "John",
					chapter: 3,
					verse:   16,
				},
				{
					book:    "John",
					chapter: 3,
					verse:   18,
				},
			},
		},
		{
			Entries: []referance{
				{
					book:    "1 John",
					chapter: 3,
					verse:   0,
				},
				{
					book:    "1 John",
					chapter: 4,
					verse:   0,
				},
			},
		},
		{
			Entries: []referance{
				{
					book:    "1 John",
					chapter: 3,
					verse:   0,
				},
				{
					book:    "Matt.",
					chapter: 4,
					verse:   0,
				},
				{
					book:    "Matt.",
					chapter: 5,
					verse:   0,
				},
			},
		},
	}

	for i, test := range tests {
		result, _ := parseCollectionRequest(test)
		expect := expectedResults[i]

		if len(expect.Entries) != len(result.Entries) {
			t.Fatalf("TEST[%d] failed expected len %d got %d",
				i,
				len(expect.Entries),
				len(result.Entries),
			)
		}

		for j, r := range result.Entries {
			e := expect.Entries[j]

			err := testReferances(e, r)

			if err != nil {
				t.Fatalf("TEST[%d] CASE[%d] failed: %s", i, j, err)
			}
		}
	}
}

// e - expected result
// r - actual result
// will check if all the fields are the same. If not
// will return an error.
func testReferances(e, r referance) error {
	if r.book != e.book {
		return errors.New(fmt.Sprintf("expected book %v got %v",
			e.book,
			r.book,
		),
		)
	}
	if r.chapter != e.chapter {
		return errors.New(fmt.Sprintf("expected chapter %v got %v",
			e.chapter,
			r.chapter,
		),
		)
	}
	if r.verse != e.verse {
		return errors.New(fmt.Sprintf("expected verse %v got %v",
			e.verse,
			r.verse,
		),
		)
	}
	return nil
}

func TestParseMixedRequest(t *testing.T) {
	//FIX: need more test cases
	tests := []string{
		"Luke 14:5-15:8,10",
		"Mark 14:5, 8-10, 20",
	}

	expectedResults := []MixedRequest{
		{
			Entries: []Request{
				RangeRequest{
					Start: referance{
						book:    "Luke",
						chapter: 14,
						verse:   5,
					},
					End: referance{
						book:    "Luke",
						chapter: 15,
						verse:   8,
					},
				},
				CollectionRequest{
					Entries: []referance{
						{
							book:    "Luke",
							chapter: 15,
							verse:   10,
						},
					},
				},
			},
		},
		{
			//"Mark 14:5, 8-10, 20",
			Entries: []Request{
				CollectionRequest{
					Entries: []referance{
						{
							book:    "Mark",
							chapter: 14,
							verse:   5,
						},
					},
				},
				RangeRequest{
					Start: referance{
						book:    "Mark",
						chapter: 14,
						verse:   8,
					},
					End: referance{
						book:    "Mark",
						chapter: 14,
						verse:   10,
					},
				},
				CollectionRequest{
					Entries: []referance{
						{
							book:    "Mark",
							chapter: 14,
							verse:   20,
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		result, _ := parseMixedRequest(test)
		expect := expectedResults[i]

		if len(expect.Entries) != len(result.Entries) {
			t.Logf("%s => %#v\n", test, result)
			t.Fatalf("TEST[%d] failed. expected etries len %d got %d",
				i,
				len(expect.Entries),
				len(result.Entries),
			)
		}

		for j, r := range result.Entries {
			e := expect.Entries[j]

			if e.Type() != r.Type() {
				t.Fatalf("TEST[%d] CASE[%d] failed. expected request type %T got %T",
					i,
					j,
					e,
					r,
				)
			}

			switch r := r.(type) {
			case RangeRequest:
				e := e.(RangeRequest)
				err := testRangeRequest(e, r)

				if err != nil {
					t.Fatalf("TEST[%d] CASE[%d] failed: %s", i, j, err)
				}

			case CollectionRequest:
				//test collection request
				e := e.(CollectionRequest)
				for k, result := range r.Entries {
					err := testReferances(e.Entries[k], result)

					if err != nil {
						t.Fatalf("TEST[%d] CASE[%d] ENTRY[%d] failed: %s", i, j, k, err)
					}
				}

			}
		}
	}
}
