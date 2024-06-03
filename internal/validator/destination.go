package validator

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// doesDirectoryExistAndIsNotEmpty checks if the directory exists and is not empty
func doesDirectoryExistAndIsNotEmpty(name string) bool {
	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			slog.Error("could not read directory", "message", err)
			return false
		}
		if len(dirEntries) > 0 {
			return true
		}
	}
	return false
}

func ValidateAndFormattedDestinationPath() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		// validate path
		// check exists
		// check right
		if doesDirectoryExistAndIsNotEmpty(args[1]) {
			return fmt.Errorf(
				"directory '%s' already exists and is not empty. Please choose a different name",
				args[1],
			)
		}

		return nil
	}
}
