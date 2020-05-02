package log

import (
	"bytes"
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

	f := Loc{Path: "foo.go", Line: 10, Col: 3}
	PrintfFile(f, "printf %s", "foo")
	WarnfFile(f, "warnf %s", "foo")
	ErrorfFile(f, "errorf %s", "foo")

	f = Loc{Path: "foo.go"}
	PrintfFile(f, "printf %s", "foo")
	WarnfFile(f, "warnf %s", "foo")
	ErrorfFile(f, "errorf %s", "foo")

	return b.String()
}
