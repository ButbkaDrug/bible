package bible

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type lineBuilder struct {
	highlightStyle string
	quoteTagStyle  string
	JesusTagStyle  string
	terminator     string

	supVerses          bool
	withChapterNumbers bool
	withVerseNumbers   bool
	chapter            int
	verse              int
	pre                string
	s                  string
	highlights         []string
}

func NewLineBuilder(v Verse) *lineBuilder {
	return &lineBuilder{
		chapter:        v.Chapter,
		verse:          v.Verse,
		s:              v.Text,
		highlights:     []string{},
		highlightStyle: "\033[1;033m",
		quoteTagStyle:  "\033[1;34m",
		JesusTagStyle:  "\033[1;31m",
		terminator:     "\033[0m",
	}

}

func NewLineBuilderWithHighlights(v Verse, hl []string) *lineBuilder {
	return &lineBuilder{
		chapter:        v.Chapter,
		verse:          v.Verse,
		s:              v.Text,
		highlights:     hl,
		highlightStyle: "\033[1;033m",
		quoteTagStyle:  "\033[1;34m",
		JesusTagStyle:  "\033[1;31m",
		terminator:     "\033[0m",
	}
}

//<pb/> - Paragraph Break
//<f></f> - Footnote
//<t></t> - quoTe?
//<J></J> - Jesus?

// \033[<style>;<foreground_color>;<background_color>mYour text\033[0m
// \033[ – The escape character to start the color code.
// <style> – Optional. Style options (e.g., bold, underline).
// <foreground_color> – The color of the text.
// <background_color> – The color of the background.
// m – Indicates the end of the color code.
// \033[0m – Resets the formatting back to normal.
func (l *lineBuilder) BoldQuotes() *lineBuilder {

	l.s = strings.ReplaceAll(l.s, "<t>", l.quoteTagStyle)
	l.s = strings.ReplaceAll(l.s, "</t>", l.terminator)
	return l

}

func (l *lineBuilder) WithChapterNumber() *lineBuilder {
	l.withChapterNumbers = true
	return l
}

func (l *lineBuilder) WithVerseNumbers() *lineBuilder {
	l.withVerseNumbers = true
	return l
}

func (l *lineBuilder) SuperscriptVerses() *lineBuilder {
	l.supVerses = true
	return l
}

// Function to convert an integer to superscript Unicode characters
// I was lazy and ask a friend to do this one for me
func toSuperscript(num int) string {
	// Mapping from regular digits to superscript Unicode characters
	superscriptMap := map[rune]rune{
		'0': '⁰', '1': '¹', '2': '²', '3': '³', '4': '⁴',
		'5': '⁵', '6': '⁶', '7': '⁷', '8': '⁸', '9': '⁹',
	}

	// Convert the number to a string
	numStr := strconv.Itoa(num)

	// Create a slice to hold the superscript characters
	var superscriptStr []rune

	// Convert each character of the number string to its superscript equivalent
	for _, char := range numStr {
		if superscript, ok := superscriptMap[char]; ok {
			superscriptStr = append(superscriptStr, superscript)
		} else {
			// If the character is not a digit (shouldn't happen here), keep it as is
			superscriptStr = append(superscriptStr, char)
		}
	}

	// Return the superscript string
	return string(superscriptStr)
}

func (l *lineBuilder) RemoveFootnoteTage() *lineBuilder {

	re := regexp.MustCompile(`<f>.*?</f>`)
	l.s = re.ReplaceAllString(l.s, "")

	return l
}

func (l *lineBuilder) ConvertPageBrakes() *lineBuilder {
	count := strings.Count(l.s, "<pb/>")
	for range count {
		idx := strings.Index(l.s, "<pb/>")
		if idx == 0 {
			l.pre = "\n"
			l.s = strings.Replace(l.s, "<pb/>", "", 1)
		}
		l.s = strings.Replace(l.s, "<pb/>", "\n", 1)

	}

	return l
}

func (l *lineBuilder) RemoveQuoteTags() *lineBuilder {
	l.s = strings.ReplaceAll(l.s, "<t>", "")
	l.s = strings.ReplaceAll(l.s, "</t>", "")
	return l
}

func (l *lineBuilder) RemoveJesusTags() *lineBuilder {
	l.s = strings.ReplaceAll(l.s, "<J>", "")
	l.s = strings.ReplaceAll(l.s, "</J>", "")
	return l
}

func (l *lineBuilder) ColorJesusTags() *lineBuilder {
	l.s = strings.ReplaceAll(l.s, "<J>", l.JesusTagStyle)
	l.s = strings.ReplaceAll(l.s, "</J>", l.terminator)

	return l
}

func (l *lineBuilder) buildVerse() string {
	if !l.withVerseNumbers {
		return ""
	}

	if l.supVerses {
		return toSuperscript(l.verse)
	}

	return fmt.Sprintf("%d", l.verse)
}

func (l *lineBuilder) buildChpater() string {
	if !l.withChapterNumbers {
		return ""
	}
	return fmt.Sprintf("%d", l.chapter)
}

func (l *lineBuilder) Build() string {
	verse := l.buildVerse()
	chapter := l.buildChpater()
	sep := " "

	if l.withChapterNumbers {
		sep = ":"
	}
	return fmt.Sprintf("%s%v%s%v%s", l.pre, chapter, sep, verse, l.s)
}

func (l *lineBuilder) Highlight() *lineBuilder {
	for _, s := range l.highlights {
		l.s = strings.ReplaceAll(
			l.s,
			s,
			fmt.Sprintf("%s%s%s", l.highlightStyle, s, l.terminator),
		)
	}

	return l
}

// will go over verses and build the referance for the chapter
// if verses are consecative will build a-b type of referance string
// if verses are not consecative will build a,b type of referance
func versesInChapter(chapter int, verses []Verse) string {
	var result string

	nums := extractVerseNumbers(chapter, verses)

	if len(nums) < 1 {
		return fmt.Sprintf("%d", chapter) // or it shuld be an empty string?
	}

	if len(nums) == 1 {
		return fmt.Sprintf("%d:%d", chapter, nums[0])
	}

	for i := range nums {
		if i == 0 {
			result = fmt.Sprintf("%d:%d", chapter, nums[i])
			continue
		}
		consec := findLastConsecotive(nums[i:])

		if consec == nums[i] {
			result = fmt.Sprintf("%s,%d", result, consec)
		} else {
			return fmt.Sprintf("%s-%d", result, consec)
		}

	}

	return result
}

func findLastConsecotive(nums []int) int {
	if len(nums) < 1 {
		return 0
	}

	for i, n := range nums {
		if i == 0 {
			continue
		}

		if n-nums[i-1] > 1 {
			return nums[i-1]
		}
	}

	return nums[len(nums)-1]

}

func extractVerseNumbers(chapter int, verses []Verse) []int {
	var ints []int

	for _, v := range verses {
		if chapter != v.Chapter {
			break
		}

		ints = append(ints, v.Verse)
	}

	return ints
}

type lineDirector struct{}

func NewLineDirector() *lineDirector {
	return &lineDirector{}
}

func (l *lineDirector) CreatePlainLine(b *lineBuilder) string {
	return b.RemoveFootnoteTage().
		ConvertPageBrakes().
		RemoveQuoteTags().
		RemoveJesusTags().
		WithVerseNumbers().
		SuperscriptVerses().
		Build()
}
func (l *lineDirector) CreateColoredLine(b *lineBuilder) string {
	return b.Highlight().
		RemoveJesusTags().
		RemoveQuoteTags().
		RemoveFootnoteTage().
		ConvertPageBrakes().
		WithVerseNumbers().
		SuperscriptVerses().
		Build()
}

type defaultRender struct {
	hl    []string
	color bool
}

func NewDefaultRender() *defaultRender {
	return &defaultRender{}
}

func (d *defaultRender) SetHighlights(ss []string) Renderer {
	d.hl = ss
	return d
}

func (d *defaultRender) Color() *defaultRender {
	d.color = true
	return d
}

// FIX: bad design
func (d *defaultRender) Render(w io.Writer, verses []Verse) error {
	var title string
	var chapter int
	var ignoreNextNewLine bool
	var director = NewLineDirector()
	for i, v := range verses {
		if title != v.Book {
			title = v.Book
		}

		if chapter != v.Chapter {
			chapter = v.Chapter
			ignoreNextNewLine = true
			pref := "\n\n"
			if i == 0 {
				pref = ""
			}
			fmt.Fprintf(w, "%s%s %s\n", pref, title, versesInChapter(chapter, verses[i:]))
		}

		builder := NewLineBuilderWithHighlights(v, d.hl)

		var text string
		if d.color {
			text = director.CreateColoredLine(builder)
		} else {

			text = director.CreatePlainLine(builder)
		}

		//every time chapter is new. not only the fist time
		if ignoreNextNewLine {
			text = strings.Replace(text, "\n", "", 1)
			ignoreNextNewLine = false
		}

		fmt.Fprintf(w, "%s", text)
	}

	return nil
}
