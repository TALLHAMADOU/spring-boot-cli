package cmd

import (
	"fmt"
	"os"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// colorEnabled reports whether ANSI colors should be written to f. Colors are
// disabled when NO_COLOR is set (https://no-color.org) or the stream is not a
// terminal (e.g. piped or redirected output), avoiding garbled escape codes.
func colorEnabled(f *os.File) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// Evaluated once at startup so the TTY check is not repeated on every log call.
var (
	stdoutColor = colorEnabled(os.Stdout)
	stderrColor = colorEnabled(os.Stderr)
)

func colorize(enabled bool, color, s string) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}

func formatMsg(format string, a ...any) string {
	msg := fmt.Sprintf(format, a...)
	return strings.TrimSuffix(msg, "\n")
}

// Success prints a green success message to stdout.
func Success(format string, a ...any) {
	fmt.Fprintln(os.Stdout, colorize(stdoutColor, colorGreen, "✅ SUCCESS: "+formatMsg(format, a...)))
}

// Error prints a red error message to stderr.
func Error(format string, a ...any) {
	fmt.Fprintln(os.Stderr, colorize(stderrColor, colorRed, "❌ ERROR: "+formatMsg(format, a...)))
}

// Warning prints a yellow warning message to stdout.
func Warning(format string, a ...any) {
	fmt.Fprintln(os.Stdout, colorize(stdoutColor, colorYellow, "⚠️ WARNING: "+formatMsg(format, a...)))
}

// Info prints a blue info message to stdout.
func Info(format string, a ...any) {
	fmt.Fprintln(os.Stdout, colorize(stdoutColor, colorBlue, "ℹ️ INFO: "+formatMsg(format, a...)))
}
