package bible

import (
	"fmt"
	"io"
	"regexp"
)

type defaultRender struct {
}

func (d defaultRender) sanitize(s string) string {
	re := regexp.MustCompile(`<f>.*?</f>`)
	return re.ReplaceAllString(s, "")

}

func (d defaultRender) Render(w io.Writer, verses []Verse) error {
	var title string
	for _, v := range verses {
		if title != v.Book {
			title = v.Book
			fmt.Fprintf(w, "%s\n", title)
		}

		text := d.sanitize(v.Text)

		fmt.Fprintf(w, "%v:%v %s\n", v.Chapter, v.Verse, text)
	}

	return nil
}

func NewDefaultRender() defaultRender {
	return defaultRender{}
}
