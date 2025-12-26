// Package cmd contains the CLI commands for bury-it.
package cmd

import (
	"fmt"
	"os"

	"github.com/deanhigh/bury-it/internal/archive"
	"github.com/spf13/cobra"
)

// Version is set at build time.
var Version = "dev"

var (
	sourceFlag      string
	graveyardFlag   string
	nameFlag        string
	dropHistoryFlag bool
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

  # Full GitHub URL with custom name
  bury-it -s https://github.com/deanhigh/experiment -g /path/to/graveyard --name my-old-experiment`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no flags provided, show help (FR-5.1)
		if sourceFlag == "" && graveyardFlag == "" {
			_ = cmd.Help()
			return
		}

		// Validate required flags (FR-5.3)
		if sourceFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: --source is required")
			fmt.Fprintln(os.Stderr, "")
			_ = cmd.Help()
			os.Exit(1)
		}

		if graveyardFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: --graveyard is required")
			fmt.Fprintln(os.Stderr, "")
			_ = cmd.Help()
			os.Exit(1)
		}

		// Execute archive
		result, err := archive.Archive(archive.Options{
			Source:      sourceFlag,
			Graveyard:   graveyardFlag,
			Name:        nameFlag,
			DropHistory: dropHistoryFlag,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Success message (FR-5.4)
		fmt.Println("")
		fmt.Printf("Successfully buried %s!\n", result.ProjectName)
		fmt.Println("")
		fmt.Println("Next steps:")
		fmt.Printf("  1. Review the changes in %s\n", result.ProjectPath)
		fmt.Println("  2. Commit the graveyard repository:")
		fmt.Printf("     cd %s && git commit -m \"Bury %s\"\n", graveyardFlag, result.ProjectName)
		fmt.Println("  3. Archive or delete the original repository")
	},
}

func init() {
	rootCmd.Flags().StringVarP(&sourceFlag, "source", "s", "", "source repository (GitHub URL, owner/repo, or local path)")
	rootCmd.Flags().StringVarP(&graveyardFlag, "graveyard", "g", "", "local path to the graveyard repository")
	rootCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "override the project name in the graveyard")
	rootCmd.Flags().BoolVar(&dropHistoryFlag, "drop-history", false, "archive only the latest state, discard git history")

	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("bury-it version {{.Version}}\n")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
