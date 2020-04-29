package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/posener/script"
)

const (
	tmplGlob   = "internal/genevents/*.go.gotmpl"
	tmplSuffix = ".gotmpl"
)

var (
	events = []string{
		"check_run",
		"check_suite",
		"create",
		"delete",
		"deployment",
		"fork",
		"gollum",
		"issue_comment",
		"issues",
		"label",
		"milestone",
		"page_build",
		"project",
		"project_card",
		"public",
		"pull_request",
		"pull_request_review",
		"pull_request_review_comment",
		"push",
		"registry_package",
		"release",
		"status",
		"watch",
		"schedule",
		"repository_dispatch",
	}

	// Events for which not to generate an info function.
	skipInfo = map[string]bool{
		"schedule":         true,
		"registry_package": true,
	}
)

var tmpl = template.Must(template.New("template").
	Funcs(template.FuncMap{
		"camel":        camel,
		"pretty":       pretty,
		"funcName":     funcName,
		"retVal":       retVal,
		"skipInfoFunc": skipInfoFunc,
	}).
	ParseGlob(tmplGlob))

func main() {
	for _, t := range tmpl.Templates() {
		out := strings.TrimSuffix(filepath.Base(t.Name()), tmplSuffix)
		f, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = t.Execute(f, events)
		if err != nil {
			panic(err)
		}

		// Format the file.
		log.Printf("Writing %s", out)
		err = script.ExecHandleStderr(os.Stderr, "goimports", "-w", out).ToStdout()
		if err != nil {
			panic(err)
		}
	}
}

func camel(name string) string {
	parts := strings.Split(name, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func pretty(name string) string {
	return strings.ReplaceAll(name, "_", " ")
}

func funcName(name string) string {
	return "Get" + camel(name)
}

func retVal(name string) string {
	return "github." + camel(name) + "Event"
}

func skipInfoFunc(event string) bool {
	return skipInfo[event]
}
