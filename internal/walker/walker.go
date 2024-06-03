package walker

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const curRoot = "."

type Copier struct {
	fsys    afero.Fs
	answes  map[string]interface{}
	funcMap map[string]any
}

func New(root string, fields map[string]interface{}) *Copier {
	c := &Copier{
		fsys:   afero.NewBasePathFs(afero.NewOsFs(), root),
		answes: fields,
	}

	c.funcMap = template.FuncMap{
		"ModuleName": func(username, projectname string) string {
			return "github.com/" + username + "/" + projectname
		},
	}

	fmt.Printf("Answers: %+v\n", fields)

	return c
}

func (cp *Copier) isNeedDelete(path string) (string, bool) {
	slog.Info("isNeedDelete", "income path", path)
	t, err := template.New("").Parse(path)
	if err != nil {
		slog.Error("check delete path", "message", err)
		return "", false
	}
	var b bytes.Buffer
	if err := t.Execute(&b, cp.answes); err != nil {
		slog.Error("renamed path", "message", err)
		return "", false
	}
	slog.Info("isNeedDelete", "formated_path", b.String())
	return b.String(), b.Len() == 0
}

func isExcludedPath(path string) bool {
	var excludedPath = []string{".git"}
	for i := range excludedPath {
		if strings.HasPrefix(path, excludedPath[i]) {
			return true
		}
	}
	return false
}

func (cp *Copier) Walk() error {
	return afero.Walk(cp.fsys, curRoot, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "walk files")
		}
		if path == curRoot || isExcludedPath(path) {
			return nil
		}

		fileName, needDelete := cp.isNeedDelete(path)
		fmt.Printf("Path: %s; Name: %s; need delete: %t\n", path, fileName, needDelete)
		if needDelete {
			return cp.fsys.RemoveAll(path)
		}

		if !info.IsDir() {
			var content []byte
			content, err = afero.ReadFile(cp.fsys, path)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			t, err := template.New("").Funcs(
				template.FuncMap(cp.funcMap),
			).Parse(string(content))
			if err != nil {
				return errors.Wrap(err, "parse content file")
			}
			var buf bytes.Buffer
			if err := t.Execute(&buf, cp.answes); err != nil {
				return errors.Wrap(err, "write content to buffer")
			}

			return afero.WriteFile(cp.fsys, fileName, buf.Bytes(), os.ModePerm)
		} else {
			exist, err := afero.DirExists(cp.fsys, fileName)
			if err != nil {
				return errors.Wrap(err, "check dir exist")
			}
			if !exist {
				err := cp.fsys.MkdirAll(fileName, os.ModePerm)
				return errors.Wrap(err, "create dir")
			}
		}

		return nil
	})
}
