package actionutil

import (
	"github.com/posener/script"
)

func GitConfig(name, email string) error {
	err := script.Exec("git", "config", "user.name", name).ToStdout()
	if err != nil {
		return err
	}
	return script.Exec("git", "config", "user.email", email).ToStdout()
}
