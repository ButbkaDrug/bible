package bible

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type defaultRender struct {
}

type line struct {
	s string
}

func (l *line) RemoveFootnoteTage() *line {

	re := regexp.MustCompile(`<f>.*?</f>`)
	l.s = re.ReplaceAllString(l.s, "")

	return l
}

func (l *line) ConvertPageBrakes() *line {
	l.s = strings.ReplaceAll(l.s, "<pb/>", "\n")
	return l
}

func (l *line) Build() string {
	return l.s
}

func (d defaultRender) Render(w io.Writer, verses []Verse) error {
	var title string
	for _, v := range verses {
		if title != v.Book {
			title = v.Book
			fmt.Fprintf(w, "%s\n", title)
		}

		line := &line{s: v.Text}

		text := line.
			RemoveFootnoteTage().
			ConvertPageBrakes().
			Build()

		fmt.Fprintf(w, "%v:%v %s\n", v.Chapter, v.Verse, text)
	}

	return nil
}

func NewDefaultRender() defaultRender {
	return defaultRender{}
}
