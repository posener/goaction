package metadata

import (
	"fmt"
	"go/ast"
	"go/doc"
	"log"
	"strconv"

	"github.com/goccy/go-yaml"
)

const (
	inputFlag = "flag"
	inputEnv  = "env"
)

// Metadata represents the structure of Github actions metadata yaml file.
// See https://help.github.com/en/actions/building-actions/metadata-syntax-for-github-actions.
type Metadata struct {
	Name     string
	Desc     string        `yaml:"description,omitempty"`
	Inputs   yaml.MapSlice `yaml:",omitempty"` // map[string]Input
	Runs     Runs
	Branding struct {
		Icon  string `yaml:",omitempty"`
		Color string `yaml:",omitempty"`
	} `yaml:",omitempty"`

	err error
}

func (m *Metadata) AddInput(name string, in Input) {
	m.Inputs = append(m.Inputs, yaml.MapItem{name, in})
}

type Input struct {
	Default  interface{} `yaml:",omitempty"`
	Desc     string      `yaml:"description"`
	Required bool

	tp string
}

type Runs struct {
	Using string
	Image string
	Env   yaml.MapSlice `yaml:",omitempty"` // map[string]string
	Args  []string      `yaml:",omitempty"`
}

func New(f *ast.File) (Metadata, error) {
	m := Metadata{
		Name: f.Name.Name,
		Desc: strconv.Quote(doc.Synopsis(f.Doc.Text())),
		Runs: Runs{
			Using: "docker",
			Image: "Dockerfile",
		},
	}

	ast.Inspect(f, m.inspect)
	if m.err != nil {
		return m, m.err
	}

	var err error
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

func (m *Metadata) inspect(n ast.Node) bool {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return true
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	cnt := true

	defer func() {
		if r := recover(); r != nil {
			cnt = false
			err, ok := r.(error)
			if ok {
				m.err = err
			} else {
				panic(r)
			}
		}
	}()

	switch name(selector.X) {
	case "flag":
		return m.inspectFlagCall(selector, call)
	case "os":
		return m.inspectOSCall(selector, call)
	case "goaction":
		return m.inspectGoactionCall(selector, call)
	}
	return cnt
}

func (m *Metadata) inspectFlagCall(selector *ast.SelectorExpr, call *ast.CallExpr) bool {
	var in Input
	var inName string
	switch name(selector.Sel) {
	case "String":
		inName = unqoute(stringValue(call.Args[0]))
		in = stringFlag(call.Args[1], call.Args[2])
	case "StringVar":
		inName = unqoute(stringValue(call.Args[1]))
		in = stringFlag(call.Args[2], call.Args[3])
	case "Int":
		inName = unqoute(stringValue(call.Args[0]))
		in = intFlag(call.Args[1], call.Args[2])
	case "IntVar":
		inName = unqoute(stringValue(call.Args[1]))
		in = intFlag(call.Args[2], call.Args[3])
	case "Bool":
		inName = unqoute(stringValue(call.Args[0]))
		in = boolFlag(call.Args[1], call.Args[2])
	case "BoolVar":
		inName = unqoute(stringValue(call.Args[1]))
		in = boolFlag(call.Args[2], call.Args[3])
	default:
		return true
	}
	in.tp = inputFlag
	m.AddInput(inName, in)
	return true
}

func (m *Metadata) inspectOSCall(selector *ast.SelectorExpr, call *ast.CallExpr) bool {
	var in Input
	var inName string
	switch name(selector.Sel) {
	case "Getenv":
		// Github is passing all environment variables with "INPUT_" prefix. Therefore it is
		// required to use the goaction environment wrapper.
		key := stringValue(call.Args[0])
		log.Fatalf("Found `os.Getenv(%s)`, use `goaction.Getenv(%s)` instead", key, key)
	default:
		return true
	}
	in.tp = inputEnv
	m.AddInput(inName, in)
	return true
}

func (m *Metadata) inspectGoactionCall(selector *ast.SelectorExpr, call *ast.CallExpr) bool {
	var in Input
	var inName string
	switch name(selector.Sel) {
	case "Getenv":
		inName = unqoute(stringValue(call.Args[0]))
		in = stringFlag(call.Args[1], call.Args[2])
	default:
		return true
	}
	in.tp = inputEnv
	m.AddInput(inName, in)
	return true
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

func stringFlag(def ast.Expr, desc ast.Expr) Input {
	var in Input
	if v := unqoute(stringValue(def)); v != "" {
		in.Default = v
	}
	in.Desc = stringValue(desc)
	return in
}

func intFlag(def ast.Expr, desc ast.Expr) Input {
	var in Input
	if v := stringValue(def); v != "" {
		var err error
		in.Default, err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
	}
	in.Desc = stringValue(desc)
	return in
}

func boolFlag(def ast.Expr, desc ast.Expr) Input {
	var in Input
	if v := stringValue(def); v != "" {
		var err error
		in.Default, err = strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
	}
	in.Desc = stringValue(desc)
	return in
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
	if err == nil {
		return uq
	}
	return s
}

func stringValue(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.BasicLit:
		return x.Value
	case *ast.Ident:
		if x.Name == "true" || x.Name == "false" {
			return x.Name
		}
		panic(fmt.Errorf("unsupported identifier: %v", x.Name))
	default:
		panic(fmt.Errorf("unsupported expression: %T", e))
	}
}
