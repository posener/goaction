package metadata

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestNew_flags(t *testing.T) {
	t.Parallel()

	content := `
// Package inputs tests parsing of input calls.
package inputs

import (
	"flag"
	"github.com/posener/goaction"
)

var (
	_ = flag.String("string", "", "string usage")
	_ = flag.String("string-default", "default", "string default usage")
	_ = flag.Int("int", 1, "int usage")
	_ = flag.Bool("bool-true", true, "bool true usage")
	_ = flag.Bool("bool-false", false, "bool false usage")

	s string
	i int
	b bool
)

func init() {
	flag.StringVar(&s, "string-var", "", "string var usage")
	flag.StringVar(&s, "string-var-default", "default", "string var default usage")
	flag.IntVar(&i, "int-var", 0, "int var usage")
	flag.BoolVar(&b, "bool-var-true", true, "bool var true usage")
	flag.BoolVar(&b, "bool-var-false", false, "bool var false usage")
}

func main() {
	_ = goaction.Getenv("env", "default", "usage of env", false)
}
`

	var want = Metadata{
		Name: "inputs",
		Desc: "Package inputs tests parsing of input calls.",
		Inputs: yaml.MapSlice{
			{"string", Input{tp: inputFlag, Desc: "string usage"}},
			{"string-default", Input{tp: inputFlag, Default: "default", Desc: "string default usage"}},
			{"int", Input{tp: inputFlag, Default: 1, Desc: "int usage"}},
			{"bool-true", Input{tp: inputFlag, Default: true, Desc: "bool true usage"}},
			{"bool-false", Input{tp: inputFlag, Default: false, Desc: "bool false usage"}},
			{"string-var", Input{tp: inputFlag, Desc: "string var usage"}},
			{"string-var-default", Input{tp: inputFlag, Default: "default", Desc: "string var default usage"}},
			{"int-var", Input{tp: inputFlag, Default: 0, Desc: "int var usage"}},
			{"bool-var-true", Input{tp: inputFlag, Default: true, Desc: "bool var true usage"}},
			{"bool-var-false", Input{tp: inputFlag, Default: false, Desc: "bool var false usage"}},

			{"env", Input{tp: inputEnv, Default: "default", Desc: "usage of env"}},
		},
		Runs: Runs{
			Using: "docker",
			Image: "Dockerfile",
			Args: []string{
				"\"-string=${{ inputs.string }}\"",
				"\"-string-default=${{ inputs.string-default }}\"",
				"\"-int=${{ inputs.int }}\"",
				"\"-bool-true=${{ inputs.bool-true }}\"",
				"\"-bool-false=${{ inputs.bool-false }}\"",
				"\"-string-var=${{ inputs.string-var }}\"",
				"\"-string-var-default=${{ inputs.string-var-default }}\"",
				"\"-int-var=${{ inputs.int-var }}\"",
				"\"-bool-var-true=${{ inputs.bool-var-true }}\"",
				"\"-bool-var-false=${{ inputs.bool-var-false }}\"",
			},
			Env: yaml.MapSlice{
				{"env", "\"${{ inputs.env }}\""},
			},
		},
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "flags.go", content, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	got, err := New(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, got, want)
}
