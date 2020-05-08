// Package metadata loads main go file to a datastructure that describes Github action metadata.
package metadata

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"strconv"

	"github.com/goccy/go-yaml"
)

const (
	inputFlag = "flag"
	inputEnv  = "env"

	docRequired = "//goaction:required"
	docSkip     = "//goaction:skip"
)

type ErrParse struct {
	Pos token.Pos
	error
}

// Metadata represents the structure of Github actions metadata yaml file.
// See https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions.
type Metadata struct {
	Name    string
	Desc    string        `yaml:"description,omitempty"`
	Inputs  yaml.MapSlice `yaml:",omitempty"` // map[string]Input
	Outputs yaml.MapSlice `yaml:",omitempty"` // map[string]Output
	Runs    Runs
	// Branding of Github action.
	// See https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions#branding
	Branding struct {
		Icon  string `yaml:",omitempty"`
		Color string `yaml:",omitempty"`
	} `yaml:",omitempty"`
}

// Input for a Github action.
// See https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions#inputs.
type Input struct {
	Default  interface{} `yaml:",omitempty"`
	Desc     string      `yaml:"description,omitempty"`
	Required bool

	tp string
}

// Output for Github action.
// See https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions#outputs.
type Output struct {
	Desc string `yaml:"description"`
}

// Runs section for "Docker" Github action.
// See https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions#runs-for-docker-actions.
type Runs struct {
	Using string // Alwasy "docker"
	Image string
	Env   yaml.MapSlice `yaml:",omitempty"` // map[string]string
	Args  []string      `yaml:",omitempty"`
}

func New(pkg *ast.Package) (Metadata, error) {
	// pkgDoc := doc.New(pkg, "", doc.AllDecls)
	m := Metadata{
		Name: pkg.Name,
		// Desc: strconv.Quote(doc.Synopsis(pkgDoc.Doc)),
		Runs: Runs{
			Using: "docker",
			Image: "Dockerfile",
		},
	}

	var err error
	ast.Inspect(pkg, func(n ast.Node) bool {
		defer func() {
			e := recover()
			if e == nil {
				return
			}
			var ok bool
			err, ok = e.(ErrParse)
			if !ok {
				panic(e)
			}
		}()
		return m.inspect(n, docStr{})
	})
	if err != nil {
		return m, err
	}
	m.Runs.Args, err = calcArgs(m.Inputs)
	if err != nil {
		return m, err
	}
	m.Runs.Env, err = calcEnv(m.Inputs)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (m *Metadata) AddInput(name string, in Input) {
	m.Inputs = append(m.Inputs, yaml.MapItem{Key: name, Value: in})
}

func (m *Metadata) AddOutput(name string, out Output) {
	m.Outputs = append(m.Outputs, yaml.MapItem{Key: name, Value: out})
}

// Inspect might panic with `parseError` when parsing failed.
func (m *Metadata) inspect(n ast.Node, d docStr) bool {
	switch v := n.(type) {
	case *ast.File:
		if v.Doc != nil {
			m.Desc = strconv.Quote(doc.Synopsis(v.Doc.Text()))
		}
		return true
	case *ast.GenDecl:
		// Decleration definition, catches "var ( ... )" segments.
		m.inspectDecl(v, d)
		return false
	case *ast.ValueSpec:
		// Value definition, catches "v := package.Func(...)"" calls."
		m.inspectValue(v, d)
		return false // Covered all inspections, no need to inspect down this node.
	case *ast.CallExpr:
		m.inspectCall(v, d)
		return true // Continue inspecting, maybe there is another call in this call.
	}
	return true
}

func (m *Metadata) inspectDecl(decl *ast.GenDecl, d docStr) {
	// Decleration can be IMPORT, CONST, TYPE, VAR. We are only interested in VAR.
	if decl.Tok != token.VAR {
		return
	}
	d.parse(decl.Doc)
	if d.skip {
		return
	}
	for _, spec := range decl.Specs {
		m.inspect(spec, d)
	}
}

func (m *Metadata) inspectValue(value *ast.ValueSpec, d docStr) {
	d.parse(value.Doc)
	if d.skip {
		return
	}
	for _, v := range value.Values {
		call, ok := v.(*ast.CallExpr)
		if !ok {
			continue
		}
		m.inspectCall(call, d)
	}
}

func (m *Metadata) inspectCall(call *ast.CallExpr, d docStr) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// Full call name, of the form: 'package.Function'.
	fullName := name(selector.X) + "." + name(selector.Sel)

	switch fullName {
	default:
		return
	case "flag.String":
		m.AddInput(
			unqoute(stringValue(call.Args[0])),
			Input{
				Default:  omitEmpty(unqoute(stringValue(call.Args[1]))),
				Desc:     stringValue(call.Args[2]),
				Required: d.required,
				tp:       inputFlag,
			})
	case "flag.StringVar":
		m.AddInput(
			unqoute(stringValue(call.Args[1])),
			Input{
				Default:  omitEmpty(unqoute(stringValue(call.Args[2]))),
				Desc:     stringValue(call.Args[3]),
				Required: d.required,
				tp:       inputFlag,
			})
	case "flag.Int":
		m.AddInput(
			unqoute(stringValue(call.Args[0])),
			Input{
				Default:  intValue(call.Args[1]),
				Desc:     stringValue(call.Args[2]),
				Required: d.required,
				tp:       inputFlag,
			})
	case "flag.IntVar":
		m.AddInput(
			unqoute(stringValue(call.Args[1])),
			Input{
				Default:  intValue(call.Args[2]),
				Desc:     stringValue(call.Args[3]),
				Required: d.required,
				tp:       inputFlag,
			})
	case "flag.Bool":
		m.AddInput(
			unqoute(stringValue(call.Args[0])),
			Input{
				Default:  boolValue(call.Args[1]),
				Desc:     stringValue(call.Args[2]),
				Required: d.required,
				tp:       inputFlag,
			})
	case "flag.BoolVar":
		m.AddInput(
			unqoute(stringValue(call.Args[1])),
			Input{
				Default:  boolValue(call.Args[2]),
				Desc:     stringValue(call.Args[3]),
				Required: d.required,
				tp:       inputFlag,
			})
	case "os.Getenv":
		m.AddInput(
			unqoute(stringValue(call.Args[0])),
			Input{
				Required: d.required,
				tp:       inputEnv,
			})
	case "goaction.Getenv":
		m.AddInput(
			unqoute(stringValue(call.Args[0])),
			Input{
				Default:  omitEmpty(unqoute(stringValue(call.Args[1]))),
				Desc:     stringValue(call.Args[2]),
				Required: d.required,
				tp:       inputEnv,
			})
	case "goaction.Output":
		m.AddOutput(
			unqoute(stringValue(call.Args[0])),
			Output{
				Desc: stringValue(call.Args[2]),
			})
	}
}

func calcArgs(inputs yaml.MapSlice /* map[string]Input */) ([]string, error) {
	var args []string
	for _, mapItem := range inputs {
		name := mapItem.Key.(string)
		input := mapItem.Value.(Input)
		if input.tp != inputFlag {
			continue
		}
		args = append(args, fmt.Sprintf("\"-%s=${{ inputs.%s }}\"", name, name))
	}
	return args, nil
}

func calcEnv(inputs yaml.MapSlice /* map[string]Input */) (yaml.MapSlice /* map[string]string */, error) {
	var envs yaml.MapSlice
	for _, mapItem := range inputs {
		name := mapItem.Key.(string)
		input := mapItem.Value.(Input)
		if input.tp != inputEnv {
			continue
		}
		envs = append(envs, yaml.MapItem{name, fmt.Sprintf("\"${{ inputs.%s }}\"", name)})
	}
	return envs, nil
}

func stringValue(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.BasicLit:
		return x.Value
	case *ast.Ident:
		if x.Name == "true" || x.Name == "false" {
			return x.Name
		}
		panic(ErrParse{error: fmt.Errorf("unsupported identifier: %v", x.Name), Pos: e.Pos()})
	default:
		panic(ErrParse{error: fmt.Errorf("unsupported expression: %T", e), Pos: e.Pos()})
	}
}

func intValue(e ast.Expr) int {
	v, err := strconv.Atoi(stringValue(e))
	if err != nil {
		panic(ErrParse{error: err, Pos: e.Pos()})
	}
	return v
}

func boolValue(e ast.Expr) bool {
	v, err := strconv.ParseBool(stringValue(e))
	if err != nil {
		panic(ErrParse{error: err, Pos: e.Pos()})
	}
	return v
}

// doc holds information from doc string.
type docStr struct {
	required bool
	skip     bool
}

// parseComment searches for a special doc is a comment group.
func (d *docStr) parse(doc *ast.CommentGroup) {
	if doc == nil {
		return
	}
	for _, comment := range doc.List {
		switch comment.Text {
		case docRequired:
			d.required = true
		case docSkip:
			d.skip = true
		}
	}
}

func name(e ast.Expr) string {
	id, ok := e.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

func unqoute(s string) string {
	uq, err := strconv.Unquote(s)
	if err != nil {
		return s
	}
	return uq
}

func omitEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
