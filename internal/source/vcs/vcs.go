package vcs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

func formatedLocalPath(s string) (string, error) {
	absPath, err := filepath.Abs(s)
	if err != nil {
		return "", err
	}
	entity, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}
	if !entity.IsDir() {
		return "", fmt.Errorf("%s is not path", s)
	}
	return absPath, nil
}

func formatedDestinationPath(s string) (string, error) {
	absPath, err := formatedLocalPath(s)
	if err != nil {
		return "", err
	}
	dirEntries, err := os.ReadDir(absPath)
	if err != nil {
		slog.Error("could not read directory", "message", err)
		return "", err
	}
	if len(dirEntries) > 0 {
		return "", fmt.Errorf(
			"directory '%s' already exists and is not empty. Please choose a different name",
			absPath,
		)
	}
	return absPath, nil
}

func (s *vcs) CopyTo(ctx context.Context, destination string) error {
	var err error

	// validate destination
	s.destination, err = formatedDestinationPath(destination)
	if err != nil {
		return err
	}

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
