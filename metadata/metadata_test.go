package metadata

import (
	"go/parser"
	"go/token"
	"testing"

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
		Inputs: map[string]Input{
			"string":             Input{tp: inputFlag, Desc: "string usage"},
			"string-default":     Input{tp: inputFlag, Default: "default", Desc: "string default usage"},
			"int":                Input{tp: inputFlag, Default: 1, Desc: "int usage"},
			"bool-true":          Input{tp: inputFlag, Default: true, Desc: "bool true usage"},
			"bool-false":         Input{tp: inputFlag, Default: false, Desc: "bool false usage"},
			"string-var":         Input{tp: inputFlag, Desc: "string var usage"},
			"string-var-default": Input{tp: inputFlag, Default: "default", Desc: "string var default usage"},
			"int-var":            Input{tp: inputFlag, Default: 0, Desc: "int var usage"},
			"bool-var-true":      Input{tp: inputFlag, Default: true, Desc: "bool var true usage"},
			"bool-var-false":     Input{tp: inputFlag, Default: false, Desc: "bool var false usage"},

			"env": Input{tp: inputEnv, Default: "default", Desc: "usage of env"},
		},
		Runs: Runs{
			Using: "docker",
			Image: "Dockerfile",
			Args: []string{
				"\"-bool-false=${{ inputs.bool-false }}\"",
				"\"-bool-true=${{ inputs.bool-true }}\"",
				"\"-bool-var-false=${{ inputs.bool-var-false }}\"",
				"\"-bool-var-true=${{ inputs.bool-var-true }}\"",
				"\"-int=${{ inputs.int }}\"",
				"\"-int-var=${{ inputs.int-var }}\"",
				"\"-string=${{ inputs.string }}\"",
				"\"-string-default=${{ inputs.string-default }}\"",
				"\"-string-var=${{ inputs.string-var }}\"",
				"\"-string-var-default=${{ inputs.string-var-default }}\"",
			},
			Env: map[string]string{
				"env": "\"${{ inputs.env }}\"",
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
