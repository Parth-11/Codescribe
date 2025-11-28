package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetGitDiff returns both unstaged and staged diffs
func GetGitDiff() (string, error) {

	//Stage everything BEFORE checking diff
	exec.Command("git", "add", ".").Run()

	//Fetch staged changes
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()

	diff := out.String()

	//No changes? return empty
	if strings.TrimSpace(diff) == "" {
		return "", nil
	}

	return diff, nil
}

// CommitAndPush performs commit + push
func CommitAndPush(message string) error {

	// Make sure everything is staged again (safe)
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Create the commit
	commitCmd := exec.Command("git", "commit", "-m", message)

	// Suppress output noise but capture errors
	err := commitCmd.Run()
	if err != nil {
		// If there is nothing to commit → avoid failure
		if strings.Contains(err.Error(), "nothing to commit") {
			fmt.Println("Nothing new to commit.")
			return nil
		}
		return fmt.Errorf("git commit failed: %w", err)
	}

	// Push commit to remote
	if err := exec.Command("git", "push").Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}
