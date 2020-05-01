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

var (
	levelDebug level = "::debug%s::"
	levelWarn  level = "::warning%s::"
	levelError level = "::error%s::"
)

func init() {
	out := os.Stdout
	if !goaction.CI {
		levelDebug = ""
		levelWarn = ""
		levelError = ""
		out = os.Stderr
	}
	logger = log.New(out, "", log.Lshortfile|log.Ltime)
}

type level string

func (p level) format(f *FileLocate) string {
	fStr := f.String()
	if len(p) == 0 {
		if len(fStr) == 0 {
			return ""
		}
		return fStr + ":"
	}
	if len(fStr) > 0 {
		fStr = " " + fStr
	}
	return fmt.Sprintf(string(p), fStr)
}

// Printf logs a debug level message.
func Printf(format string, args ...interface{}) {
	PrintfFile(nil, format, args...)
}

// Printf logs a debug level message with a file location.
func PrintfFile(f *FileLocate, format string, args ...interface{}) {
	logger.Printf(levelDebug.format(f)+format, args...)
}

// Warnf logs a warning level message.
func Warnf(format string, args ...interface{}) {
	WarnfFile(nil, format, args...)
}

// WarnfFile logs a warning level message with a file location.
func WarnfFile(f *FileLocate, format string, args ...interface{}) {
	logger.Printf(levelWarn.format(f)+format, args...)
}

// Errorf logs an error level message.
func Errorf(format string, args ...interface{}) {
	ErrorfFile(nil, format, args...)
}

// ErrorfFile logs an error level message with a file location.
func ErrorfFile(f *FileLocate, format string, args ...interface{}) {
	logger.Printf(levelError.format(f)+format, args...)
}

// Fatalf logs an error level message, and fails the program.
func Fatalf(format string, args ...interface{}) {
	FatalfFile(nil, format, args...)
}

// FatalfFile logs an error level message with a file location, and fails the program.
func FatalfFile(f *FileLocate, format string, args ...interface{}) {
	logger.Fatalf(levelError.format(f)+format, args...)
}

// Fatal logs an error level message, and fails the program.
func Fatal(v ...interface{}) {
	FatalFile(nil, v...)
}

// FatalFile logs an error level message with a file location, and fails the program.
func FatalFile(f *FileLocate, v ...interface{}) {
	logger.Fatal(append([]interface{}{levelError.format(f)}, v...)...)
}

// FileLocate provides file infromation for logging purposes.
type FileLocate struct {
	Path string
	Line int
	Col  int
}

func (f *FileLocate) String() string {
	if f == nil {
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
