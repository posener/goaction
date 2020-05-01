package main

import (
	"flag"
	"log"
	"os"
)

var fail = false

func main() {
	flag.Parse()

	if got, want := os.Getenv("set"), "set"; got != want {
		errorf(`os.Getenv("set") = %s, want %s`, got, want)
	}

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
