/*
Package goaction enables writing Github Actions in Go.

The idea is - write a Go script one that you could also run with `go run`, and use it as Github
action. This project will create all the required files generated specifically for your script. Your
script can also use this library in order to get workflow information.

Required steps

- [x] Create a script in Go.

- [x] Add `goaction` configuration in `.github/workflows/goaction.yml`.

- [x] Push the project to Github

- [x] Tell me about it, I'll link it below.

Writing Action Go Script

Create a Go project. Currently it must be using Go modules, and compilable with Go 1.14. The action
script is simply a main package located somewhere in this project.

One limitation is that this main package should contain a main file, and this file must contain all
the required inputs for this binary. This inputs are defined with the standard `flag` package for
command line arguments, or by `goaction.Getenv` for environment variables.

> These inputs define the API of the program. `goaction` automatically detect them and creates the
> `action.yml` file from them.

Additionally, goaction also provides a library that exposes all Github action envirnment in an
easy-to-use API. See the documentation for more information.

Code segments which should run only in Github action (called "CI mode"), and not when the main
package runs as a command line tool, should be protected by `if goaction.CI { ... }`.

Goaction Configuration

In order to convert the repository to a Github action, goaction command line should run on the main
Go file. This command can run manually (by ./cmd/goaction) but luckly `goaction` also comes as a
Github action :-)

Goaction Github action keeps the Github action file updated according to the main Go file
automatically. When a PR is made, goaction will post a review explaining what changes to expect.
When a new commit is pushed, goreadme makes sure that the Github action files are updated if needed.

Add the following content to `.github/workflows/goaction.yml`

	on:
	  pull_request:
	    branches: [master]
	  push:
	    branches: [master]
	jobs:
	  goaction:
	    runs-on: ubuntu-latest
	    steps:
	    - name: Check out repository
	      uses: actions/checkout@v2
	    - name: Update action files
	      uses: posener/goaction@v1
	      with:
	        path: <path to main file>.
	        # Optional: required only for commenting on PRs.
	        github-token: '${{ secrets.GITHUB_TOKEN }}'
	        # Other inputs... see ./action.yml

Goaction Artifacts

./action.yml: A "metadata" file for Github actions. If this file exists, the repository is
considered as Github action, and the file contains information that instructs how to invoke this
action. See (metadata syntax) https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions.
for more info.

./Dockerfile: A file that contains instructions how to build a container, that is used for Github
actions. Github action uses this file in order to create a container image to the action. The
container can also be built and tested manually:

	$ docker build -t my-action .
	$ docker run --rm my-action

Using Goaction

* (posener/goreadme) http://github.com/posener/goreadme

*/
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
	//
	// This variable can be used to protect code segments which should only run in Github action
	// mode and not in command line mode:
	//
	//  if goaction.CI {
	// 		// Code that should run only in Github action mode.
	// 	}
	CI = os.Getenv("CI") == "true"

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
	// In Github action mode, update the environment variable name according to match the Github
	// action modifications.
	if CI {
		name = "INPUT_" + strings.ToUpper(name)
	}
	v := os.Getenv(name)
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
