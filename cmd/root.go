package cmd

import (
	"fmt"
	"os"
	"waifetch/internal"

	"github.com/spf13/cobra"
)

var (
	language      string
	listLanguages bool

	rootCmd = &cobra.Command{
		Use:   "waifetch",
		Short: "Display anime girls holding programming books in your terminal uWu",
		RunE: func(cmd *cobra.Command, args []string) error {
			fetcher := internal.NewFetcher()

			if listLanguages {
				languages, err := fetcher.GetLanguageFolders()
				if err != nil {
					return fmt.Errorf("failed to get languages: %w", err)
				}

				fmt.Println("Available programming languages:")
				for i, lang := range languages {
					if i > 0 && i%6 == 0 {
						fmt.Println()
					}
					fmt.Printf("%-20s", lang)
				}
				fmt.Println()
				return nil
			}

			return fetcher.FetchRandomImage(language)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&language, "lang", "l", "",
		"Specify programming language (e.g., C, Go, Rust)")

	rootCmd.PersistentFlags().BoolVarP(&listLanguages, "list", "", false,
		"List all available programming languages")

	rootCmd.Example = `  waifetch
  waifetch --list
  waifetch --lang Go
  waifetch -l Rust`
}
