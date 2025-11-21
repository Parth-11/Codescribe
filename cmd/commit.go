package cmd

import (
	"fmt"

	"github.com/Parth-11/Codescribe/internal/ai"
	"github.com/Parth-11/Codescribe/internal/git"
	"github.com/Parth-11/Codescribe/internal/prompt"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate AI commit messages and push",
	RunE: func(cmd *cobra.Command, args []string) error {
		diff, err := git.GetGitDiff()
		if err != nil {
			return err
		}
		if diff == "" {
			fmt.Println("No changes detected.")
			return nil
		}

		messages, err := ai.GenerateCommitMessages(diff)
		if err != nil {
			return err
		}

		choice, err := prompt.SelectCommitMessage(messages)
		if err != nil {
			return err
		}

		fmt.Println("Using commit:", choice)
		return git.CommitAndPush(choice)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
