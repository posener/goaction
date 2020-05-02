package log

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/posener/goaction"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	old := goaction.CI
	defer func() { goaction.CI = old }()

	t.Run("CI=true", func(t *testing.T) {
		goaction.CI = true
		initFormats()

		want := `::debug::debugf foo
printf foo
::warning::warnf foo
::error::errorf foo
::debug file=foo.go,line=10,col=3::debugf foo
::warning file=foo.go,line=10,col=3::warnf foo
::error file=foo.go,line=10,col=3::errorf foo
::debug file=foo.go::debugf foo
::warning file=foo.go::warnf foo
::error file=foo.go::errorf foo
`

		assert.Equal(t, want, logThings())
	})

	t.Run("CI=false", func(t *testing.T) {
		goaction.CI = false
		initFormats()

		want := `debugf foo
printf foo
warnf foo
errorf foo
foo.go+10:3: debugf foo
foo.go+10:3: warnf foo
foo.go+10:3: errorf foo
foo.go: debugf foo
foo.go: warnf foo
foo.go: errorf foo
`

		assert.Equal(t, want, logThings())
	})
}

func logThings() string {
	var b bytes.Buffer
	logger.SetOutput(&b)

	Debugf("debugf %s", "foo")
	Printf("printf %s", "foo")
	Warnf("warnf %s", "foo")
	Errorf("errorf %s", "foo")

	p := token.Position{Filename: "foo.go", Line: 10, Column: 3}
	DebugfFile(p, "debugf %s", "foo")
	WarnfFile(p, "warnf %s", "foo")
	ErrorfFile(p, "errorf %s", "foo")

	p = token.Position{Filename: "foo.go"}
	DebugfFile(p, "debugf %s", "foo")
	WarnfFile(p, "warnf %s", "foo")
	ErrorfFile(p, "errorf %s", "foo")

	return b.String()
}
