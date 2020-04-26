package goaction

import (
	"os"
	"strconv"
	"strings"
)

// Github actions default environment variables.
// See https://help.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables#default-environment-variables
var (
	// CI is set to true when running under github action
	CI            = os.Getenv("CI") == "true"
	Home          = os.Getenv("HOME")
	Workflow      = os.Getenv("GITHUB_WORKFLOW")
	RunID         = os.Getenv("GITHUB_RUN_ID")
	RunNum        = os.Getenv("GITHUB_RUN_NUMBER")
	ActionID      = os.Getenv("GITHUB_ACTION")
	Actor         = os.Getenv("GITHUB_ACTOR")
	Repository    = os.Getenv("GITHUB_REPOSITORY")
	Workspace     = os.Getenv("GITHUB_WORKSPACE")
	SHA           = os.Getenv("GITHUB_SHA")
	Ref           = os.Getenv("GITHUB_REF")
	ForkedHeadRef = os.Getenv("GITHUB_HEAD_REF")
	ForkedBaseRef = os.Getenv("GITHUB_BASE_REF")

	eventName = os.Getenv("GITHUB_EVENT_NAME") // Use IsPush/IsPR instead.
	eventPath = os.Getenv("GITHUB_EVENT_PATH")

	repoParts = strings.Split(Repository, "/")
)

// Return environment variables from Github action. Providing a default value, usage string.
func Getenv(name string, value string, usage string) string {
	v := os.Getenv("INPUT_" + strings.ToUpper(name))
	if v == "" {
		return value
	}
	return v
}

// Owner returns the name of the owner of the Github repository.
func Owner() string {
	if len(repoParts) < 2 {
		return ""
	}
	return repoParts[0]
}

// Project returns the name of the project of the Github repository.
func Project() string {
	if len(repoParts) < 2 {
		return ""
	}
	return repoParts[1]
}

// IsPR returns true in push mode.
func IsPush() bool {
	return eventName == "push"
}

// IsPR returns true in pull request mode.
func IsPR() bool {
	return eventName == "pull_request"
}

// Branch returns the push branch for push flow or empty string for other flows.
func Branch() string {
	if IsPush() {
		return strings.Split(Ref, "/")[2]
	}
	return ""

}

// PrNum returns pull request number for PR flow or 0 in other flows.
func PrNum() int {
	if IsPR() {
		// Ref is in the form: "refs/pull/:prNumber/merge"
		// See https://help.github.com/en/actions/reference/events-that-trigger-workflows#pull-request-event-pull_request
		num, err := strconv.Atoi(strings.Split(Ref, "/")[2])
		if err != nil {
			panic(err) // Should not happen
		}
		return num
	}
	return 0
}

// IsForked return true if the action is running on a forked repository.
func IsForked() bool {
	return ForkedBaseRef != ""
}
