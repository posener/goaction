# goaction

[![Build Status](https://travis-ci.org/posener/goaction.svg?branch=master)](https://travis-ci.org/posener/goaction)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/posener/goaction)

Package goaction enables writing Github Actions in Go.

The idea is - write a Go script one that you could also run with `go run`, and use it as Github
action. This project will create all the required files generated specifically for your script. Your
script can also use this library in order to get workflow information.

## Required steps

- [x] Create a script in Go.

- [x] Add `goaction` configuration in `.github/workflows/goaction.yml`.

- [x] Push the project to Github

- [x] Tell me about it, I'll link it below.

## Writing Action Go Script

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

## Goaction Configuration

In order to convert the repository to a Github action, goaction command line should run on the main
Go file. This command can run manually (by [./cmd/goaction](./cmd/goaction)) but luckly `goaction` also comes as a
Github action :-)

Goaction Github action keeps the Github action file updated according to the main Go file
automatically. When a PR is made, goaction will post a review explaining what changes to expect.
When a new commit is pushed, goreadme makes sure that the Github action files are updated if needed.

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
      path: <path to main file>.
      # Optional: required only for commenting on PRs.
	  github-token: '${{ secrets.GITHUB_TOKEN }}'
	  # Other falgs... see [./action.yml](./action.yml)
```

## Goaction Artifacts

[./action.yml](./action.yml): A "metadata" file for Github actions. If this file exists, the repository is
considered as Github action, and the file contains information that instructs how to invoke this
action. See [https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions](https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions).
for more info.

[./Dockerfile](./Dockerfile): A file that contains instructions how to build a container, that is used for Github
actions. Github action uses this file in order to create a container image to the action. The
container can also be built and tested manually:

```go
$ docker build -t my-action .
$ docker run --rm my-action
```

## Using Goaction

* [posener/goreadme]([http://github.com/posener/goreadme](http://github.com/posener/goreadme))

## Sub Packages

* [actionutil](./actionutil): Package actionutil provides utility functions for Github actions.

* [cmd/goaction](./cmd/goaction): Creates action files for Go code

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
