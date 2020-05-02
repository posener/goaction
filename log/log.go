// Package log is an alternative package for standard library "log" package for logging in Github
// action environment. It behaves as expected both in CI mode and in non-CI mode.
//
// 	 import (
// 	-	"log"
// 	+	"github.com/posener/goaction/log"
// 	 )
package log

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/posener/goaction"
)

var logger *log.Logger

const (
	levelDebug level = "::debug%s::"
	levelWarn  level = "::warning%s::"
	levelError level = "::error%s::"
)

func init() {
	out := os.Stdout
	if !goaction.CI {
		out = os.Stderr
	}
	logger = log.New(out, "", 0)
}

type level string

func (l level) format(lc Loc) string {
	loc := lc.String()
	if !goaction.CI {
		if len(loc) > 0 {
			loc = loc + ": "
		}
		return loc
	}
	if len(loc) > 0 {
		loc = " " + loc
	}
	return fmt.Sprintf(string(l), loc)
}

// Printf logs a debug level message.
func Printf(format string, args ...interface{}) {
	PrintfFile(Loc{}, format, args...)
}

// Printf logs a debug level message with a file location.
func PrintfFile(f Loc, format string, args ...interface{}) {
	logger.Printf(levelDebug.format(f)+format, args...)
}

// Warnf logs a warning level message.
func Warnf(format string, args ...interface{}) {
	WarnfFile(Loc{}, format, args...)
}

// WarnfFile logs a warning level message with a file location.
func WarnfFile(f Loc, format string, args ...interface{}) {
	logger.Printf(levelWarn.format(f)+format, args...)
}

// Errorf logs an error level message.
func Errorf(format string, args ...interface{}) {
	ErrorfFile(Loc{}, format, args...)
}

// ErrorfFile logs an error level message with a file location.
func ErrorfFile(f Loc, format string, args ...interface{}) {
	logger.Printf(levelError.format(f)+format, args...)
}

// Fatalf logs an error level message, and fails the program.
func Fatalf(format string, args ...interface{}) {
	FatalfFile(Loc{}, format, args...)
}

// FatalfFile logs an error level message with a file location, and fails the program.
func FatalfFile(f Loc, format string, args ...interface{}) {
	logger.Fatalf(levelError.format(f)+format, args...)
}

// Fatal logs an error level message, and fails the program.
func Fatal(v ...interface{}) {
	FatalFile(Loc{}, v...)
}

// FatalFile logs an error level message with a file location, and fails the program.
func FatalFile(f Loc, v ...interface{}) {
	logger.Fatal(append([]interface{}{levelError.format(f)}, v...)...)
}

// Loc provides file location infromation for logging purposes.
type Loc struct {
	Path string // Path to file.
	Line int    // Line in file.
	Col  int    // Col is column in line.
}

func (f Loc) String() string {
	if f.Path == "" {
		return ""
	}
	parts := []string{"file=" + f.Path}
	if f.Line > 0 {
		parts = append(parts, fmt.Sprintf("line=%d", f.Line))
		if f.Col > 0 {
			parts = append(parts, fmt.Sprintf("col=%d", f.Col))
		}
	}
	return strings.Join(parts, ",")
}

// Mask a term in the logs (will appear as '*' instead.)
func Mask(term string) {
	if !goaction.CI {
		return
	}
	fmt.Println("::add-mask::" + term)
}
