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
	"go/token"
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

var (
	// File position formattings:
	formatFile string
	formatLine string
	formatCol  string
	formatJoin string
)

func init() {
	out := os.Stdout
	if !goaction.CI {
		out = os.Stderr
	}
	logger = log.New(out, "", 0)
	initFormats()
}

// initFormats initializes format strings. Exists for testing purposes.
func initFormats() {
	if goaction.CI {
		formatFile = "file=%s"
		formatLine = "line=%d"
		formatCol = "col=%d"
		formatJoin = ","
	} else {
		formatFile = "%s"
		formatLine = "+%d"
		formatCol = ":%d"
		formatJoin = ""
	}
}

type level string

func (l level) format(p token.Position) string {
	pos := posString(p)
	if !goaction.CI {
		if len(pos) > 0 {
			pos = pos + ": "
		}
		return pos
	}
	if len(pos) > 0 {
		pos = " " + pos
	}
	return fmt.Sprintf(string(l), pos)
}

// Printf logs an info level message.
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// Debugf logs a debug level message. To view these logs, set secret ACTIONS_STEP_DEBUG=true at
// https://github.com/<repo>/settings/secrets/new.
func Debugf(format string, args ...interface{}) {
	DebugfFile(token.Position{}, format, args...)
}

// DebugfFile logs a debug level message with a file location. To view these logs, set secret
// variable ACTIONS_STEP_DEBUG=true at https://github.com/<repo>/settings/secrets/new.
func DebugfFile(p token.Position, format string, args ...interface{}) {
	logger.Printf(levelDebug.format(p)+format, args...)
}

// Warnf logs a warning level message.
func Warnf(format string, args ...interface{}) {
	WarnfFile(token.Position{}, format, args...)
}

// WarnfFile logs a warning level message with a file location.
func WarnfFile(p token.Position, format string, args ...interface{}) {
	logger.Printf(levelWarn.format(p)+format, args...)
}

// Errorf logs an error level message.
func Errorf(format string, args ...interface{}) {
	ErrorfFile(token.Position{}, format, args...)
}

// ErrorfFile logs an error level message with a file location.
func ErrorfFile(p token.Position, format string, args ...interface{}) {
	logger.Printf(levelError.format(p)+format, args...)
}

// Fatalf logs an error level message, and fails the program.
func Fatalf(format string, args ...interface{}) {
	FatalfFile(token.Position{}, format, args...)
}

// FatalfFile logs an error level message with a file location, and fails the program.
func FatalfFile(p token.Position, format string, args ...interface{}) {
	logger.Fatalf(levelError.format(p)+format, args...)
}

// Fatal logs an error level message, and fails the program.
func Fatal(v ...interface{}) {
	FatalFile(token.Position{}, v...)
}

// FatalFile logs an error level message with a file location, and fails the program.
func FatalFile(p token.Position, v ...interface{}) {
	logger.Fatal(append([]interface{}{levelError.format(p)}, v...)...)
}

func posString(p token.Position) string {
	if p.Filename == "" {
		return ""
	}

	parts := []string{fmt.Sprintf(formatFile, p.Filename)}
	if p.Line > 0 {
		parts = append(parts, fmt.Sprintf(formatLine, p.Line))
		if p.Column > 0 {
			parts = append(parts, fmt.Sprintf(formatCol, p.Column))
		}
	}
	return strings.Join(parts, formatJoin)
}

// Mask a term in the logs (will appear as '*' instead.)
func Mask(term string) {
	if !goaction.CI {
		return
	}
	fmt.Println("::add-mask::" + term)
}
