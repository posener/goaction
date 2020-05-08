package main

import (
	"flag"
	"os"

	"github.com/posener/goaction/log"
)

var (
	message = flag.String("message", "", "Input from first job.")
	arg     = flag.Int("arg", 0, "Int value.")

	fail = false
)

func main() {
	flag.Parse()

	// Tests that input from first job passed correctly.

	if got, want := *message, "message"; got != want {
		errorf(`flag.String("message") = %s, want %s`, got, want)
	}

	if got, want := *arg, 42; got != want {
		errorf(`flag.Int("arg") = %d, want %d`, got, want)
	}

	// Test that the environment variable is not set with standard name:
	if got, want := os.Getenv("env"), "env"; got != want {
		errorf(`os.Getenv("env") = %s, want %s`, got, want)
	}

	if fail {
		log.Fatalf("Failed...")
	}
}

func errorf(f string, args ...interface{}) {
	log.Printf(f, args...)
	fail = true
}
