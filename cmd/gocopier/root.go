package gocopier

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/1gkx/gocopier/internal/configurator"
	"github.com/1gkx/gocopier/internal/source"
	"github.com/1gkx/gocopier/internal/walker"
	"github.com/1gkx/gocopier/internal/wizard"
	"github.com/spf13/cobra"
)

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

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "gocopier",
		Short: "short description gocopier app",
		Long:  "long description gocopier app",
		Args:  cobra.MatchAll(cobra.ExactArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			// read config file

			src, err := source.New(args[0])
			if err != nil {
				slog.Error("parse source path", "message", err)
				cobra.CheckErr(err)
			}

			destination, err := formatedDestinationPath(args[1])
			if err != nil {
				slog.Error("parse destination path", "message", err)
				cobra.CheckErr(err)
			}

			// dowload or copy files from source to destination
			if err = src.CopyTo(cmd.Context(), destination); err != nil {
				slog.Error("Copy files", "message", err)
				cobra.CheckErr(err)
			}

			// read file with questions
			q := configurator.Read(src.GetConfigFile())
			// create survey
			answers, err := wizard.New(q)
			if err != nil {
				slog.Error("wizad", "message", err)
				cobra.CheckErr(err)
			}
			// run generator
			// run templator
			if err := walker.New(args[1], answers).Walk(); err != nil {
				slog.Error("walker", "message", err)
				cobra.CheckErr(err)
			}
		},
	}

	// rootCmd.Flags().String("config", "", "конфигурационный файл")

	if err := rootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}
