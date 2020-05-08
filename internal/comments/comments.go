package comments

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
)

var (
	docRequired = regexp.MustCompile("^//goaction:required$")
	docSkip     = regexp.MustCompile("^//goaction:skip$")
	docDefault  = regexp.MustCompile("^//goaction:default (.*)$")
	docDesc     = regexp.MustCompile("^//goaction:description (.*)$")
)

// Comments holds information from doc string.
type Comments struct {
	Required Bool
	Skip     Bool
	Default  String
	Desc     String
}

type Bool struct {
	token.Pos
	Value bool
}

type String struct {
	token.Pos
	Value string
}

// parseComment searches for a special doc is a comment group.
func (d *Comments) Parse(doc *ast.CommentGroup) {
	if doc == nil {
		return
	}
	for _, comment := range doc.List {
		txt := comment.Text
		pos := comment.Slash
		switch {
		case docRequired.MatchString(txt):
			d.Required = Bool{Value: true, Pos: pos}
		case docSkip.MatchString(txt):
			d.Skip = Bool{Value: true, Pos: pos}
		case docDefault.MatchString(txt):
			d.Default = String{Value: docDefault.FindStringSubmatch(txt)[1], Pos: pos}
		case docDesc.MatchString(txt):
			d.Desc = String{Value: strconv.Quote(docDesc.FindStringSubmatch(txt)[1]), Pos: pos}
		}
	}
}
