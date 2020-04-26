package actionutil

import (
	"fmt"

	"github.com/posener/goaction"
	"github.com/posener/script"
)

func GitConfig(name, email string) error {
	err := script.Exec("git", "config", "user.name", name).ToStdout()
	if err != nil {
		return err
	}
	return script.Exec("git", "config", "user.email", email).ToStdout()
}

// Returns diff of changes in a given file. The file will be staged after running this command.
func GitDiff(path string) (string, error) {
	// Add files to git, in case it does not exists
	err := script.Exec("git", "add", path).ToStdout()
	if err != nil {
		return "", fmt.Errorf("git add for %s: %s", path, err)
	}
	return script.Exec("git", "diff", "--staged", "--no-color", path).Tail(-5).ToString()
}

// Commits and pushes a list of files
func GitCommitPush(paths []string, message string) error {
	err := script.Exec("git", "reset").ToStdout()
	if err != nil {
		return fmt.Errorf("git reset: %s", err)
	}
	err = script.Exec("git", append([]string{"add"}, paths...)...).ToStdout()
	if err != nil {
		return fmt.Errorf("git add: %s", err)
	}
	err = script.Exec("git", "commit", "-m", message).ToStdout()
	if err != nil {
		return fmt.Errorf("git commit: %s", err)
	}
	return script.Exec("git", "push", "origin", "HEAD:"+goaction.Branch()).ToStdout()
}
