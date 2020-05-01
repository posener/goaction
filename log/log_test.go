package log

import (
	"bytes"
	"testing"

	"github.com/posener/goaction"
	"github.com/stretchr/testify/assert"
)

func TestLogCI(t *testing.T) {
	if !goaction.CI {
		t.Skip("Only runs in CI")
	}

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
	got := logThings()

	assert.Equal(t, want, got)
}

func TestCLI(t *testing.T) {
	if goaction.CI {
		t.Skip("Only runs not in CI")
	}

	want := `printf foo
warnf foo
errorf foo
file=foo.go,line=10,col=3:printf foo
file=foo.go,line=10,col=3:warnf foo
file=foo.go,line=10,col=3:errorf foo
file=foo.go:printf foo
file=foo.go:warnf foo
file=foo.go:errorf foo
`
	got := logThings()

	assert.Equal(t, want, got)
}

func logThings() string {
	var b bytes.Buffer
	logger.SetOutput(&b)
	logger.SetFlags(0)

	Printf("printf %s", "foo")
	Warnf("warnf %s", "foo")
	Errorf("errorf %s", "foo")

	f := &FileLocate{Path: "foo.go", Line: 10, Col: 3}
	PrintfFile(f, "printf %s", "foo")
	WarnfFile(f, "warnf %s", "foo")
	ErrorfFile(f, "errorf %s", "foo")

	f = &FileLocate{Path: "foo.go"}
	PrintfFile(f, "printf %s", "foo")
	WarnfFile(f, "warnf %s", "foo")
	ErrorfFile(f, "errorf %s", "foo")

	return b.String()
}
