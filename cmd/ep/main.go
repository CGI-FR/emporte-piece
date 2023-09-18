// Copyright (C) 2023 CGI France
//
// This file is part of emporte-piece.
//
// Emporte-piece is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Emporte-piece is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with emporte-piece.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/cgi-fr/emporte-piece/internal/infra"
	"github.com/cgi-fr/emporte-piece/pkg/filetree"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	name      string // Provisioned by ldflags.
	version   string // Provisioned by ldflags.
	commit    string // Provisioned by ldflags.
	buildDate string // Provisioned by ldflags.
	builtBy   string // Provisioned by ldflags.

	verbosity string
	jsonlog   bool
	debug     bool
	colormode string

	outputDir string
	format    string
)

func main() {
	cobra.OnInitialize(initLog)

	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:     fmt.Sprintf("%v path/to/template/dir < context.yaml", name),
		Short:   "Emporte-pièce",
		Long:    `Emporte-pièce will help you bootstrap your project.`,
		Version: fmt.Sprintf(`%v (commit=%v date=%v by=%v)`, version, commit, buildDate, builtBy),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Info().
				Str("verbosity", verbosity).
				Bool("log-json", jsonlog).
				Bool("debug", debug).
				Str("color", colormode).
				Msg("start")
		},
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args[0], outputDir, format); err != nil {
				log.Fatal().Err(err).Msg("end")
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			log.Info().Int("return", 0).Msg("end")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "info",
		"set level of log verbosity : none (0), error (1), warn (2), info (3), debug (4), trace (5)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "add debug information to logs (very slow)")
	rootCmd.PersistentFlags().BoolVar(&jsonlog, "log-json", false, "output logs in JSON format")
	rootCmd.PersistentFlags().StringVar(&colormode, "color", "auto", "use colors in log outputs : yes, no or auto")

	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", ".", "output directory")
	rootCmd.PersistentFlags().
		StringVarP(&format, "format", "f", "yaml", "format of context data : yaml, json or jsonl (default=yaml)")

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("error when executing command")
		os.Exit(1)
	}
}

func run(_ *cobra.Command, templateDir, outputDir, format string) error {
	var contextReader infra.ContextReader

	switch strings.ToLower(format) {
	case "yaml", "yml":
		contextReader = infra.NewContextReaderYAML(os.Stdin)
	case "json":
		contextReader = infra.NewContextReaderJSON(os.Stdin)
	case "jsonl":
		contextReader = infra.NewContextReaderJSONL(os.Stdin)
	}

	for contextReader.HasNext() {
		context, err := contextReader.Next()
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		driver := filetree.NewDriver(infra.FileSystem{})

		err = driver.Develop(templateDir, outputDir, context)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func initLog() {
	color := false

	switch strings.ToLower(colormode) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS != "windows" {
			color = true
		}
	case "yes", "true", "1", "on", "enable":
		color = true
	}

	if jsonlog {
		log.Logger = zerolog.New(os.Stderr)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !color}) //nolint:exhaustruct
	}

	if debug {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	setVerbosity()
}

func setVerbosity() {
	switch verbosity {
	case "trace", "5":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug", "4":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info", "3":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "2":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error", "1":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
