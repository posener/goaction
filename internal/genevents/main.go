// Generates events.go and events_test.go files.
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/posener/script"
)

const tmplGlob = "internal/genevents/*.go.gotmpl"

type event struct {
	Name             string
	SkipEventGetFunc bool
}

func (e event) CamelCase() string {
	parts := strings.Split(e.Name, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func (e event) Pretty() string {
	return strings.ReplaceAll(e.Name, "_", " ")
}

func (e event) EventGetFuncName() string {
	return "Get" + e.CamelCase()
}

func (e event) GithubReturnValue() string {
	return "github." + e.CamelCase() + "Event"
}

var events = []event{
	{Name: "check_run"},
	{Name: "check_suite"},
	{Name: "create"},
	{Name: "delete"},
	{Name: "deployment"},
	{Name: "fork"},
	{Name: "gollum"},
	{Name: "issue_comment"},
	{Name: "issues"},
	{Name: "label"},
	{Name: "milestone"},
	{Name: "page_build"},
	{Name: "project"},
	{Name: "project_card"},
	{Name: "public"},
	{Name: "pull_request"},
	{Name: "pull_request_review"},
	{Name: "pull_request_review_comment"},
	{Name: "push"},
	{Name: "registry_package", SkipEventGetFunc: true},
	{Name: "release"},
	{Name: "status"},
	{Name: "watch"},
	{Name: "schedule", SkipEventGetFunc: true},
	{Name: "repository_dispatch"},
}

var tmpl = template.Must(template.ParseGlob(tmplGlob))

func main() {
	for _, t := range tmpl.Templates() {
		out := outFileName(t.Name())
		f, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		log.Printf("Writing %s", out)
		err = t.Execute(f, events)
		if err != nil {
			panic(err)
		}

		// Format the file.
		err = script.ExecHandleStderr(os.Stderr, "goimports", "-w", out).ToStdout()
		if err != nil {
			panic(err)
		}
	}
}

func outFileName(templateName string) string {
	name := filepath.Base(templateName)
	// Remove .gotmpl suffix.
	return name[:strings.LastIndex(name, ".")]
}
