package cmd

import (
	"fmt"
	"os"

	"github.com/Parth-11/Codescribe/internal/ai"
	"github.com/Parth-11/Codescribe/internal/git"
	"github.com/Parth-11/Codescribe/internal/prompt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codescribe",
	Short: "AI-powered conventional commit generator using Groq",
	RunE: func(cmd *cobra.Command, args []string) error {

		diff, err := git.GetGitDiff()
		if err != nil {
			return err
		}
		if diff == "" {
			fmt.Println("No changes detected.")
			return nil
		}

		fmt.Println("Changes detected. Generating commit messages...")

		messages, err := ai.GenerateCommitMessages(diff)
		if err != nil {
			return err
		}

		choice, err := prompt.SelectCommitMessage(messages)
		if err != nil {
			return err
		}

		fmt.Println("Using commit message:")
		fmt.Println(choice)

		if err := git.CommitAndPush(choice); err != nil {
			return err
		}

		fmt.Println("Commit pushed successfully!")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
