# goaction

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/posener/goaction)

Package goaction enables writing Github Actions in Go.

The idea is: write a standard Go script, one that works with `go run`, and use it as Github action.
The script's inputs - flags and environment variables, are set though the Github action API. This
project will generate all the required files for the script (This generation can be done
automattically with Github action integration). The library also exposes neat API to get workflow
information.

## Required Steps

- [x] Write a Go script.

- [x] Add `goaction` configuration in `.github/workflows/goaction.yml`.

- [x] Push the project to Github.

See simplest example for a Goaction script: [posener/goaction-example](https://github.com/posener/goaction-example).

## Writing a Goaction Script

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

## Goaction Configuration

In order to convert the repository to a Github action, goaction command line should run on the
**"main file"** (described above). This command can run manually (by [./cmd/goaction](./cmd/goaction)) but luckily
`goaction` also comes as a Github action :-)

Goaction Github action keeps the Github action file updated according to the main Go file
automatically. When a PR is made, goaction will post a review explaining what changes to expect.
When a new commit is pushed, Goaction makes sure that the Github action files are updated if needed.

Add the following content to `.github/workflows/goaction.yml`

```go
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
      uses: [./](./)
```

## Goaction Artifacts

[./action.yml](./action.yml): A "metadata" file for Github actions. If this file exists, the repository is
considered as Github action, and the file contains information that instructs how to invoke this
action. See [metadata syntax](https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions).
for more info.

[./Dockerfile](./Dockerfile): A file that contains instructions how to build a container, that is used for Github
actions. Github action uses this file in order to create a container image to the action. The
container can also be built and tested manually:

```go
$ docker build -t my-action .
$ docker run --rm my-action
```

## Annotations

Goaction parses Go script file and looks for annotations that extends the information that exists in
the function calls. Goaction annotations are a comments that start with `//goaction:` (no space
after slashes). They can only be set on a `var` definition. The following annotations are available:

* `//goaction:required` - sets an input definition to be "required".

* `//goaction:skip` - skips an input out output definition.

* `//goaction:description <description>` - add description for `os.Getenv`.

* `//goaction:default <value>` - add default value for `os.Getenv`.

## Using Goaction

A list of projects which are using Goaction (please send a PR if your project uses goaction and does
not appear her).

* [posener/goreadme](http://github.com/posener/goreadme)

## Sub Packages

* [actionutil](./actionutil): Package actionutil provides utility functions for Github actions.

* [log](./log): Package log is an alternative package for standard library "log" package for logging in Github action environment.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
