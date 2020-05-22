// Generates actionutil/githubapi.go
package main

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/posener/autogen"
	"golang.org/x/tools/go/loader"
)

//go:generate go run .

type fn struct {
	Field  *types.Var
	Method *types.Func
}

func (f fn) Params() []*types.Var {
	var vars []*types.Var
	params := f.Method.Type().(*types.Signature).Params()
	for i := 0; i < params.Len(); i++ {
		vars = append(vars, params.At(i))
	}
	return vars
}

func (f fn) Results() []*types.Var {
	var vars []*types.Var
	results := f.Method.Type().(*types.Signature).Results()
	for i := 0; i < results.Len(); i++ {
		vars = append(vars, results.At(i))
	}
	return vars
}

func (f fn) OtherParamsDefinition() string {
	others := f.Params()[3:]
	var params []string
	for _, other := range others {
		params = append(params, fmt.Sprintf("%s %s", other.Name(), cleanTypeName(other.Type())))
	}
	return strings.Join(params, ", ")
}

func (f fn) OtherParamsUse() string {
	others := f.Params()[3:]
	var params []string
	for _, other := range others {
		params = append(params, other.Name())
	}
	return strings.Join(params, ", ")
}

func (f fn) ResultsDefinition() string {
	var params []string
	for _, other := range f.Results() {
		params = append(params, fmt.Sprintf("%s %s", other.Name(), cleanTypeName(other.Type())))
	}
	ret := strings.Join(params, ", ")
	if len(params) > 1 {
		ret = "(" + ret + ")"
	}
	return ret
}

func (f fn) isValid() bool {
	params := f.Params()
	return len(params) >= 2 &&
		params[0].Name() == "ctx" &&
		params[1].Name() == "owner" &&
		params[2].Name() == "repo"
}

func cleanTypeName(tp types.Type) string {
	name := tp.String()
	if i := strings.LastIndex(name, "/"); i >= 0 {
		prefix := regexp.MustCompile(`^[\[\]\*]*`).FindString(name)
		name = prefix + name[i+1:]
	}
	return name
}

func main() {
	log.SetFlags(log.Lshortfile)

	// Load the github program.
	conf := loader.Config{AllowErrors: true}
	conf.Import("github.com/google/go-github/v31/github")
	stderr := os.Stderr
	os.Stderr, _ = os.Open(os.DevNull)
	program, err := conf.Load()
	os.Stderr.Close()
	os.Stderr = stderr
	if err != nil {
		log.Fatal(err)
	}

	// Get github package.
	var pkg *types.Package
	for pkg = range program.AllPackages {
		if pkg.Name() == "github" {
			break
		}
	}
	if pkg == nil {
		log.Fatal("Package github was not found.")
	}

	// Get `type Client struct`:
	client := pkg.Scope().Lookup("Client").Type().Underlying().(*types.Struct)

	var funcs []fn

	// Iterate Client fields and collect the services.
	for i := 0; i < client.NumFields(); i++ {
		field := client.Field(i)
		if !field.Exported() {
			continue
		}
		// The field is a pointer to a struct, get the pointer.
		fieldPointer, ok := field.Type().Underlying().Underlying().(*types.Pointer)
		if !ok {
			continue
		}
		// The pointer points on a struct that ends with the "Service" suffix.
		fieldType, ok := fieldPointer.Elem().(*types.Named)
		if !ok {
			continue
		}
		if !strings.HasSuffix(fieldType.Obj().Name(), "Service") {
			continue
		}

		// Iterate over the field type methods.
		for j := 0; j < fieldType.NumMethods(); j++ {
			f := fn{Field: field, Method: fieldType.Method(j)}
			if !f.isValid() {
				continue
			}
			// Skip Repositories.GetArchiveLink: https://github.com/google/go-github/issues/1533.
			if f.Field.Name() == "Repositories" && f.Method.Name() == "GetArchiveLink" {
				continue
			}
			funcs = append(funcs, f)
		}
	}

	err = autogen.Execute(
		funcs,
		autogen.Location(filepath.Join(autogen.ModulePath, "actionutil")))
	if err != nil {
		log.Fatal(err)
	}
}
