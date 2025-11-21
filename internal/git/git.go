package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// capture git diff
func GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

// commit + push
func CommitAndPush(message string) error {

	// stage all changes
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// commit
	if err := exec.Command("git", "commit", "-m", message).Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	// push
	if err := exec.Command("git", "push").Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}
