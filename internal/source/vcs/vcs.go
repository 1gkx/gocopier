package vcs

import (
	"context"
	"errors"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var (
	gitPostfix = ".git"

	replacements = map[*regexp.Regexp]string{
		regexp.MustCompile(`^gh:?(.*)\.git$`): `https://github.com/$1.git`,
		regexp.MustCompile(`^gh:?(.*)$`):      `https://github.com/$1.git`,
		regexp.MustCompile(`^gl:?(.*)\.git$`): `https://gitlab.com/$1.git`,
		regexp.MustCompile(`^gl:?(.*)$`):      `https://gitlab.com/$1.git`,
	}
)

func startWith(source string, prefs ...string) bool {
	for i := range prefs {
		if strings.HasPrefix(source, prefs[i]) {
			return true
		}
	}
	return false
}

func endWith(source, end string) bool {
	return strings.HasSuffix(source, end)
}

type vcs struct {
	source, destination string
}

func New(sourceDst string) (*vcs, error) {
	for pattern, replacement := range replacements {
		sourceDst = pattern.ReplaceAllString(sourceDst, replacement)
	}

	if !startWith(sourceDst, "https://", "http://") {
		return nil, errors.New("vsc source must be started as https:// or http://")
	}

	if !endWith(sourceDst, gitPostfix) {
		sourceDst += gitPostfix
	}

	return &vcs{source: sourceDst}, nil
}

func (s *vcs) GetConfigFile() string {
	return path.Join(s.destination, "copier.yaml")
}

func (s *vcs) CopyTo(ctx context.Context, destination string) error {
	var err error

	cmd := exec.CommandContext(ctx,
		"git", "clone", "--progress", "--verbose",
		s.source,
		destination,
	)

	err = cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}
