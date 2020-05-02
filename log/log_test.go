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

		want := `::debug::printf foo
::warning::warnf foo
::error::errorf foo
::debug file=foo.go,line=10,col=3::printf foo
::warning file=foo.go,line=10,col=3::warnf foo
::error file=foo.go,line=10,col=3::errorf foo
::debug file=foo.go::printf foo
::warning file=foo.go::warnf foo
::error file=foo.go::errorf foo
`

		assert.Equal(t, want, logThings())
	})

	t.Run("CI=false", func(t *testing.T) {
		goaction.CI = false

		want := `printf foo
warnf foo
errorf foo
file=foo.go,line=10,col=3: printf foo
file=foo.go,line=10,col=3: warnf foo
file=foo.go,line=10,col=3: errorf foo
file=foo.go: printf foo
file=foo.go: warnf foo
file=foo.go: errorf foo
`

		assert.Equal(t, want, logThings())
	})
}

func logThings() string {
	var b bytes.Buffer
	logger.SetOutput(&b)

	Printf("printf %s", "foo")
	Warnf("warnf %s", "foo")
	Errorf("errorf %s", "foo")

	p := token.Position{Filename: "foo.go", Line: 10, Column: 3}
	PrintfFile(p, "printf %s", "foo")
	WarnfFile(p, "warnf %s", "foo")
	ErrorfFile(p, "errorf %s", "foo")

	p = token.Position{Filename: "foo.go"}
	PrintfFile(p, "printf %s", "foo")
	WarnfFile(p, "warnf %s", "foo")
	ErrorfFile(p, "errorf %s", "foo")

	return b.String()
}
