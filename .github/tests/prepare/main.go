package main

import (
	"os"

	"github.com/posener/goaction"
	"github.com/posener/goaction/log"
)

var fail = false

func main() {
	// Write an output for the second action.
	goaction.Output("out", "message", "output of first action")
	// Set an environment variable for the second action.
	err := goaction.Setenv("set", "set")
	if err != nil {
		log.Fatal(err)
	}
	// Set an environment variable for the second action using export.
	err = goaction.Export("export", "export")
	if err != nil {
		log.Fatal(err)
	}

	// Setenv environment variable should not be able to be accessed by this action.
	if got, want := os.Getenv("set"), ""; got != want {
		errorf(`os.Getenv("set") = %s, want %s`, got, want)
	}

	// Exported environment variable should be able to be accessed by this action.
	if got, want := os.Getenv("export"), "export"; got != want {
		errorf(`os.Getenv("export") = %s, want %s`, got, want)
	}

	if fail {
		log.Fatalf("Failed...")
	}
}

func errorf(f string, args ...interface{}) {
	log.Printf(f, args...)
	fail = true
}
