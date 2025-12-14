package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0-alpha"

func main() {
	rootCmd := &cobra.Command{
		Use:     "growth",
		Short:   "Git-native career development manager",
		Long:    `growth.md - Track your skills, goals, and learning paths in plain Markdown files powered by AI`,
		Version: version,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
