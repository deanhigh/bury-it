// Package cmd contains the CLI commands for bury-it.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time.
var Version = "dev"

var (
	source      string
	graveyard   string
	dropHistory bool
)

var rootCmd = &cobra.Command{
	Use:   "bury-it",
	Short: "Sunset experimental projects while preserving their history",
	Long: `bury-it is a CLI tool to sunset experimental projects by archiving them
into a local "graveyard" repository while optionally preserving their full git history.

It supports both remote GitHub repositories and local git repositories as sources.`,
	Example: `  # Bury a GitHub repository
  bury-it --source deanhigh/old-project --graveyard ~/graveyard

  # Bury a local repository without preserving history
  bury-it --source ./my-experiment --graveyard ~/graveyard --drop-history

  # Full GitHub URL
  bury-it -s https://github.com/deanhigh/experiment -g /path/to/graveyard`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no flags provided, show help
		if source == "" && graveyard == "" {
			cmd.Help()
			return
		}

		// Validate required flags
		if source == "" {
			fmt.Fprintln(os.Stderr, "Error: --source is required")
			fmt.Fprintln(os.Stderr, "")
			cmd.Help()
			os.Exit(1)
		}

		if graveyard == "" {
			fmt.Fprintln(os.Stderr, "Error: --graveyard is required")
			fmt.Fprintln(os.Stderr, "")
			cmd.Help()
			os.Exit(1)
		}

		// TODO: Implement archive logic
		fmt.Printf("Source: %s\n", source)
		fmt.Printf("Graveyard: %s\n", graveyard)
		fmt.Printf("Drop history: %v\n", dropHistory)
		fmt.Println("")
		fmt.Println("Archive logic not yet implemented.")
	},
}

func init() {
	rootCmd.Flags().StringVarP(&source, "source", "s", "", "source repository (GitHub URL, owner/repo, or local path)")
	rootCmd.Flags().StringVarP(&graveyard, "graveyard", "g", "", "local path to the graveyard repository")
	rootCmd.Flags().BoolVar(&dropHistory, "drop-history", false, "archive only the latest state, discard git history")

	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("bury-it version {{.Version}}\n")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
