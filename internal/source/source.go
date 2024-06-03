package source

import (
	"context"

	"github.com/1gkx/gocopier/internal/source/vcs"
)

// var availableVcsSchemes = []string{
// 	"gh:", "gl:", "https://", "git@", "git://", "git+",
// }

type Source interface {
	CopyTo(ctx context.Context, destination string) error
	GetConfigFile() string
}

func New(sourceDst string) (Source, error) {
	return vcs.New(sourceDst)
}
