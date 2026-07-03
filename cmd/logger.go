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

func formatMsg(format string, a ...interface{}) string {
	msg := fmt.Sprintf(format, a...)
	return strings.TrimSuffix(msg, "\n")
}

// Success prints a green success message
func Success(format string, a ...interface{}) {
	fmt.Printf("%s✅ SUCCESS: %s%s\n", colorGreen, formatMsg(format, a...), colorReset)
}

// Error prints a red error message to stderr
func Error(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s❌ ERROR: %s%s\n", colorRed, formatMsg(format, a...), colorReset)
}

// Warning prints a yellow warning message
func Warning(format string, a ...interface{}) {
	fmt.Printf("%s⚠️ WARNING: %s%s\n", colorYellow, formatMsg(format, a...), colorReset)
}

// Info prints a blue info message
func Info(format string, a ...interface{}) {
	fmt.Printf("%sℹ️ INFO: %s%s\n", colorBlue, formatMsg(format, a...), colorReset)
}
