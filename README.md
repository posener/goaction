# goaction

[![Build Status](https://travis-ci.org/posener/goaction.svg?branch=master)](https://travis-ci.org/posener/goaction)
[![codecov](https://codecov.io/gh/posener/goaction/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/goaction)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/posener/goaction)

```go
Package goaction enables writing Github Actions in Go.
```

The idea is - write a Go script, and use it as Github action. This script can be seemlesly run using
both `go run` command line, or by Github action! This project will create all the required files
(`Dockerfile`, `action.yml`) generated specifically for your script. Your script can also use this
library in order to get workflow information.

The things that you are required to do:

1. In a github repository, have a main Go file that contains your script logic.

2. Add `goaction` configuration in `.github/workflows/goaction.yml`.

1. How to write the main Go file

The main Go file, is a single file, and can be located anywhere in the repository. It does **not**
need to contain all the action logic. It **does** need to contain all the inputs definitions.

Input definitions, as explained, should only be defined in the main Go file can be set using `flag`
package to get inputs from command line flags, or by `goreadme.Getenv` to get inputs from
environment variables.

> The inputs define the API of the program. `goaction` automatically detect these calls and creates
> the `action.yml` file from them.

Goaction also provides a library that exposes all Github action envirnment which easy-to-use
variabls / functions. See the documentation for more information.

Code segments which we want to run only in Github action, and not when invoking the script with as
a command line tool, should be protected by `if goaction.CI { ... }`.

2. `goaction` configuration

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
```

## Goaction artifacts

action.yml: Github uses this file to infer that this repository is a github action, and how to run
it.

Dockerfile: Github action uses this file in order to create a container image to the action. It can
also be built and run manually:

```go
$ docker build -t my-action .
$ docker run --rm my-action
```

## Sub Packages

* [actionutil](./actionutil)

* [cmd/goaction](./cmd/goaction): Creates action files for Go code

* [metadata](./metadata)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
