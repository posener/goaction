/*
Package goaction enables writing Github Actions in Go.

The idea is: write a standard Go script, one that works with `go run`, and use it as Github action.
The script's inputs - flags and environment variables, are set though the Github action API. This
project will generate all the required files for the script (This generation can be done
automattically with Github action integration). The library also exposes neat API to get workflow
information.

Required Steps

- [x] Write a Go script.

- [x] Add `goaction` configuration in `.github/workflows/goaction.yml`.

- [x] Push the project to Github.

See simplest example for a Goaction script: (posener/goaction-example) https://github.com/posener/goaction-example.

Writing a Goaction Script

Write Github Action by writing Go code! Just start a Go module with a main package, and execute it
as a Github action using Goaction, or from the command line using `go run`.

A go executable can get inputs from the command line flag and from environment variables. Github
actions should have a `action.yml` file that defines this API. Goaction bridges the gap by parsing
the Go code and creating this file automatically for you.

The main package inputs should be defined with the standard `flag` package for command line
arguments, or by `os.Getenv` for environment variables. These inputs define the API of the program
and `goaction` automatically detect them and creates the `action.yml` file from them.

Additionally, goaction also provides a library that exposes all Github action environment in an
easy-to-use API. See the documentation for more information.

Code segments which should run only in Github action (called "CI mode"), and not when the main
package runs as a command line tool, should be protected by a `if goaction.CI { ... }` block.

Goaction Configuration

In order to convert the repository to a Github action, goaction command line should run on the
**"main file"** (described above). This command can run manually (by ./cmd/goaction) but luckily
`goaction` also comes as a Github action :-)

Goaction Github action keeps the Github action file updated according to the main Go file
automatically. When a PR is made, goaction will post a review explaining what changes to expect.
When a new commit is pushed, Goaction makes sure that the Github action files are updated if needed.

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
	        # Optional: required only for commenting on PRs.
	        github-token: '${{ secrets.GITHUB_TOKEN }}'
	    # Optional: now that the script is a Github action, it is possible to run it in the
	    # workflow.
	    - name: Example
	      uses: ./

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

Annotations

Goaction parses Go script file and looks for annotations that extends the information that exists in
the function calls. Goaction annotations are a comments that start with `//goaction:` (no space
after slashes). They can only be set on a `var` definition. The following annotations are available:

* `//goaction:required` - sets an input definition to be "required".

* `//goaction:skip` - skips an input out output definition.

* `//goaction:description <description>` - add description for `os.Getenv`.

* `//goaction:default <value>` - add default value for `os.Getenv`.

Using Goaction

A list of projects which are using Goaction (please send a PR if your project uses goaction and does
not appear her).

* (posener/goreadme) http://github.com/posener/goreadme

*/
package goaction

import (
	"fmt"
	"log"
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
	// The path to the GitHub home directory used to store user data. For example, /github/home.
	Home = os.Getenv("HOME")
	// The name of the workflow
	Workflow = os.Getenv("GITHUB_WORKFLOW")
	// 	A unique number for each run within a repository. This number does not change if you re-run
	// the workflow run.
	RunID = os.Getenv("GITHUB_RUN_ID")
	// 	A unique number for each run of a particular workflow in a repository. This number begins at
	// 1 for the workflow's first run, and increments with each new run. This number does not change
	// if you re-run the workflow run.
	RunNum = os.Getenv("GITHUB_RUN_NUMBER")
	// The unique identifier (id) of the action.
	ActionID = os.Getenv("GITHUB_ACTION")
	// The name of the person or app that initiated the workflow. For example, octocat.
	Actor = os.Getenv("GITHUB_ACTOR")
	// The owner and repository name. For example, octocat/Hello-World.
	Repository = os.Getenv("GITHUB_REPOSITORY")
	// The name of the webhook event that triggered the workflow.
	Event = EventType(os.Getenv("GITHUB_EVENT_NAME"))
	// 	The GitHub workspace directory path. The workspace directory contains a subdirectory with a
	// copy of your repository if your workflow uses the actions/checkout action. If you don't use
	// the actions/checkout action, the directory will be empty. For example,
	// /home/runner/work/my-repo-name/my-repo-name.
	Workspace = os.Getenv("GITHUB_WORKSPACE")
	// The commit SHA that triggered the workflow. For example,
	// ffac537e6cbbf934b08745a378932722df287a53.
	SHA = os.Getenv("GITHUB_SHA")
	// The branch or tag ref that triggered the workflow. For example, refs/heads/feature-branch-1.
	// If neither a branch or tag is available for the event type, the variable will not exist.
	Ref = os.Getenv("GITHUB_REF")
	//	Only set for forked repositories. The branch of the head repository.
	ForkedHeadRef = os.Getenv("GITHUB_HEAD_REF")
	// Only set for forked repositories. The branch of the base repository.
	ForkedBaseRef = os.Getenv("GITHUB_BASE_REF")

	eventPath = os.Getenv("GITHUB_EVENT_PATH")

	repoParts = strings.Split(Repository, "/")
)

func init() {
	if CI {
		// Set the default logging to stdout since Github actions treats stderr as error level logs.
		log.SetOutput(os.Stdout)
	}
}

// Setenv sets an environment variable that will only be visible for all following Github actions in
// the current workflow, but not in the current action.
// See https://help.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-environment-variable.
func Setenv(name string, value string) {
	if !CI {
		return
	}
	// Store in the given environment variable name such that programs that expect this environment
	// variable (not through goaction) can get it.
	fmt.Printf("::set-env name=%s::%s\n", name, value)
}

// Export sets an environment variable that will also be visible for all following Github actions in
// the current workflow.
func Export(name string, value string) error {
	err := os.Setenv(name, value)
	if err != nil {
		return err
	}
	Setenv(name, value)
	return nil
}

// Output sets Github action output.
// See https://help.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-output-parameter.
func Output(name string, value string, desc string) {
	if !CI {
		return
	}
	fmt.Printf("::set-output name=%s::%s\n", name, value)
}

// AddPath prepends a directory to the system PATH variable for all subsequent actions in the
// current job. The currently running action cannot access the new path variable.
// See https://help.github.com/en/actions/reference/workflow-commands-for-github-actions#adding-a-system-path.
func AddPath(path string) {
	if !CI {
		return
	}
	fmt.Printf("::add-path::%s\n", path)
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

// Branch returns the push branch for push flow or empty string for other flows.
func Branch() string {
	if Event == EventPush {
		return strings.Split(Ref, "/")[2]
	}
	return ""

}

// PrNum returns pull request number for PR flow or -1 in other flows.
func PrNum() int {
	if Event == EventPullRequest {
		// Ref is in the form: "refs/pull/:prNumber/merge"
		// See https://help.github.com/en/actions/reference/events-that-trigger-workflows#pull-request-event-pull_request
		num, err := strconv.Atoi(strings.Split(Ref, "/")[2])
		if err != nil {
			panic(err) // Should not happen
		}
		return num
	}
	return -1
}

// IsForked return true if the action is running on a forked repository.
func IsForked() bool {
	return ForkedBaseRef != ""
}
