package metadata

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	code := `
// Package main tests parsing of input calls.
package main

import (
	"flag"
	"os"
	"github.com/posener/goaction"
)

var (
	_ = flag.String("string", "", "string usage")
	_ = flag.String("string-default", "default", "string default usage")
	_ = flag.Int("int", 1, "int usage")
	_ = flag.Bool("bool-true", true, "bool true usage")
	_ = flag.Bool("bool-false", false, "bool false usage")

	_ = os.Getenv("env")

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
	goaction.Output("out", "value", "output description")
}
`

	var want = Metadata{
		Name: "main",
		Desc: "\"Package main tests parsing of input calls.\"",
		Inputs: yaml.MapSlice{
			{"string", Input{tp: inputFlag, Desc: "\"string usage\""}},
			{"string-default", Input{tp: inputFlag, Default: "default", Desc: "\"string default usage\""}},
			{"int", Input{tp: inputFlag, Default: 1, Desc: "\"int usage\""}},
			{"bool-true", Input{tp: inputFlag, Default: true, Desc: "\"bool true usage\""}},
			{"bool-false", Input{tp: inputFlag, Default: false, Desc: "\"bool false usage\""}},
			{"env", Input{tp: inputEnv}},
			{"string-var", Input{tp: inputFlag, Desc: "\"string var usage\""}},
			{"string-var-default", Input{tp: inputFlag, Default: "default", Desc: "\"string var default usage\""}},
			{"int-var", Input{tp: inputFlag, Default: 0, Desc: "\"int var usage\""}},
			{"bool-var-true", Input{tp: inputFlag, Default: true, Desc: "\"bool var true usage\""}},
			{"bool-var-false", Input{tp: inputFlag, Default: false, Desc: "\"bool var false usage\""}},
		},
		Outputs: yaml.MapSlice{
			{"out", Output{Desc: "\"output description\""}},
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

	got, err := parse(code)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

// Tests cases of goaction:required comment.
func TestNewRequired(t *testing.T) {
	t.Parallel()

	code := `
package main

import (
	"flag"
	"os"
	"github.com/posener/goaction"
)

var (
	// Test two following definitions. The required should apply only to the first.
	//goaction:required
	_ = flag.String("simple1", "", "simple1")
	_ = flag.String("simple2", "", "simple2")

	// Test multiple definitions.
	//goaction:required
	_, _ = flag.String("multi1", "", "multi1"), flag.String("multi2", "", "multi2")
)

// Test var definition.
//goaction:required
var _ = flag.String("var", "", "var")

// Test var block.
//goaction:required
var (
	_ = flag.String("block1", "", "block1")
	_ = flag.String("block2", "", "block2")
)

var (
	// Test environment variable required and description.
	//goaction:required
	_ = os.Getenv("env")
)
`

	var wantInputs = yaml.MapSlice{
		{"simple1", Input{tp: inputFlag, Desc: "\"simple1\"", Required: true}},
		{"simple2", Input{tp: inputFlag, Desc: "\"simple2\""}},
		{"multi1", Input{tp: inputFlag, Desc: "\"multi1\"", Required: true}},
		{"multi2", Input{tp: inputFlag, Desc: "\"multi2\"", Required: true}},
		{"var", Input{tp: inputFlag, Desc: "\"var\"", Required: true}},
		{"block1", Input{tp: inputFlag, Desc: "\"block1\"", Required: true}},
		{"block2", Input{tp: inputFlag, Desc: "\"block2\"", Required: true}},
		{"env", Input{tp: inputEnv, Required: true}},
	}

	got, err := parse(code)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, wantInputs, got.Inputs)
}

func TestNewDefaultDesc(t *testing.T) {
	t.Parallel()

	code := `
package main

import (
	"flag"
	"os"
	"github.com/posener/goaction"
)

// Test environment variable required and description.
//goaction:default default
//goaction:description input from environment variable
var	_ = os.Getenv("env")
`

	var wantInputs = yaml.MapSlice{
		{"env", Input{tp: inputEnv, Default: "default", Desc: "\"input from environment variable\""}},
	}

	got, err := parse(code)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, wantInputs, got.Inputs)
}

// Tests cases of goaction:skip comment.
func TestNewSkip(t *testing.T) {
	t.Parallel()

	code := `
package main

import (
	"flag"
	"github.com/posener/goaction"
)

var (
	// Test two following definitions. The skip should apply only to the first.
	//goaction:skip
	_ = flag.String("simple1", "", "simple1")
	_ = flag.String("simple2", "", "simple2")

	// Test multiple definitions.
	//goaction:skip
	_, _ = flag.String("multi1", "", "multi1"), flag.String("multi2", "", "multi2")
)

// Test var definition.
//goaction:skip
var _ = flag.String("var", "", "var")

// Test var block.
//goaction:skip
var (
	_ = flag.String("block1", "", "block1")
	_ = flag.String("block2", "", "block2")
)

var (
	// Test environment variable required and description.
	//goaction:skip
	_ = os.Getenv("env")
)
`

	var wantInputs = yaml.MapSlice{
		{"simple2", Input{tp: inputFlag, Desc: "\"simple2\""}},
	}

	got, err := parse(code)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, wantInputs, got.Inputs)
}

func TestNewInvalidAnnotations(t *testing.T) {
	t.Parallel()

	codes := []string{
		`
package main
import "flag"

//goaction:description description
var _ = flag.String("simple1", "", "simple1")
`,
		`
package main
import "flag"

//goaction:default default
var _ = flag.String("simple1", "", "simple1")
`,
	}

	for _, code := range codes {
		t.Run(code, func(t *testing.T) {
			_, err := parse(strings.TrimSpace(code))
			assert.Error(t, err)
		})
	}
}

func TestMarshal(t *testing.T) {
	m := Metadata{
		Name: "name",
		Desc: "description",
		Inputs: yaml.MapSlice{
			{"in2", Input{tp: "tp2", Default: 1, Desc: "description 2"}},
			{"in1", Input{tp: "tp1", Default: "string", Desc: "description 1"}},
		},
		Runs: Runs{
			Using: "using",
			Image: "image",
			Args:  []string{"arg1", "arg2"},
			Env: yaml.MapSlice{
				{"key2", "value2"},
				{"key1", "value1"},
			},
		},
	}

	m.Branding.Color = "color"
	m.Branding.Icon = "icon"

	want := `name: name
description: description
inputs:
  in2:
    default: 1
    description: description 2
    required: false
  in1:
    default: string
    description: description 1
    required: false
runs:
  using: using
  image: image
  env:
    key2: value2
    key1: value1
  args:
  - arg1
  - arg2
branding:
  icon: icon
  color: color
`
	got, err := yaml.Marshal(m)
	require.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func parse(code string) (Metadata, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", code, parser.ParseComments)
	if err != nil {
		return Metadata{}, err
	}
	pkg := &ast.Package{
		Name:  "main",
		Files: map[string]*ast.File{"main.go": f},
	}
	return New(pkg)
}
