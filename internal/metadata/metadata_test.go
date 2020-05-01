package metadata

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	code := `
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

	_ = goaction.Getenv("env", "default", "env usage")

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
	goaction.Output("out", "value", "out description")
}
`

	var want = Metadata{
		Name: "inputs",
		Desc: "\"Package inputs tests parsing of input calls.\"",
		Inputs: yaml.MapSlice{
			{"string", Input{tp: inputFlag, Default: "", Desc: "\"string usage\""}},
			{"string-default", Input{tp: inputFlag, Default: "default", Desc: "\"string default usage\""}},
			{"int", Input{tp: inputFlag, Default: 1, Desc: "\"int usage\""}},
			{"bool-true", Input{tp: inputFlag, Default: true, Desc: "\"bool true usage\""}},
			{"bool-false", Input{tp: inputFlag, Default: false, Desc: "\"bool false usage\""}},
			{"env", Input{tp: inputEnv, Default: "default", Desc: "\"env usage\""}},
			{"string-var", Input{tp: inputFlag, Default: "", Desc: "\"string var usage\""}},
			{"string-var-default", Input{tp: inputFlag, Default: "default", Desc: "\"string var default usage\""}},
			{"int-var", Input{tp: inputFlag, Default: 0, Desc: "\"int var usage\""}},
			{"bool-var-true", Input{tp: inputFlag, Default: true, Desc: "\"bool var true usage\""}},
			{"bool-var-false", Input{tp: inputFlag, Default: false, Desc: "\"bool var false usage\""}},
		},
		Outputs: yaml.MapSlice{
			{"out", Output{Desc: "\"out description\""}},
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
	assert.Equal(t, got, want)
}

func TestNewDocStr(t *testing.T) {
	t.Parallel()

	code := `
package dockstr

import (
	"flag"
	"github.com/posener/goaction"
)

var (
	// Test simple case
	//goaction:required
	_ = flag.String("simple", "", "simple")

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
	_ = goaction.Getenv("env", "", "env")
)
`

	var wantInputs = yaml.MapSlice{
		{"simple", Input{tp: inputFlag, Default: "", Desc: "\"simple\"", Required: true}},
		{"multi1", Input{tp: inputFlag, Default: "", Desc: "\"multi1\"", Required: true}},
		{"multi2", Input{tp: inputFlag, Default: "", Desc: "\"multi2\"", Required: true}},
		{"var", Input{tp: inputFlag, Default: "", Desc: "\"var\"", Required: true}},
		{"block1", Input{tp: inputFlag, Default: "", Desc: "\"block1\"", Required: true}},
		{"block2", Input{tp: inputFlag, Default: "", Desc: "\"block2\"", Required: true}},
		{"env", Input{tp: inputEnv, Default: "", Desc: "\"env\"", Required: true}},
	}

	got, err := parse(code)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, got.Inputs, wantInputs)
}

func TestOsGetenvFailure(t *testing.T) {
	t.Parallel()

	code := `
package required
import "os"
var env = os.Getenv("env")
`

	_, err := parse(code)
	assert.Error(t, err)
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
	return New(f)
}
