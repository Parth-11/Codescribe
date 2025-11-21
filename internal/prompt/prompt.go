package prompt

import (
	"github.com/manifoldco/promptui"
)

func SelectCommitMessage(messages []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select Commit Message",
		Items: messages,
	}

	_, result, err := prompt.Run()
	return result, err
}
