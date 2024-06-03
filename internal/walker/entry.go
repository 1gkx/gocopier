package walker

import (
	"bytes"
	"log/slog"
	"regexp"
	"text/template"
)

var regTemplate = regexp.MustCompile("{{.*}}")

type Entry struct {
	originEntry   string
	formatedEntry string
	isTemplate    bool
}

func NewEntry(et string) Entry {
	tmpl := regTemplate.MatchString(et)
	return Entry{
		originEntry:   et,
		formatedEntry: et,
		isTemplate:    tmpl,
	}
}

func (e *Entry) Execute(data map[string]any) error {
	t, err := template.New("").Parse(e.originEntry)
	if err != nil {
		slog.Error("check delete path", "message", err)
		return err
	}
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		slog.Error("renamed path", "message", err)
		return err
	}
	e.formatedEntry = b.String()

	return nil
}
