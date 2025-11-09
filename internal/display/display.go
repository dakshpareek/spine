package display

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	green      = color.New(color.FgGreen)
	yellow     = color.New(color.FgYellow)
	blue       = color.New(color.FgCyan)
	red        = color.New(color.FgRed)
	bold       = color.New(color.Bold)
	symbolOK   = green.Sprint("✓")
	symbolWarn = yellow.Sprint("⚠")
	symbolInfo = blue.Sprint("⏳")
	symbolErr  = red.Sprint("✗")
)

// Success formats a success message with a green check mark.
func Success(format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", symbolOK, fmt.Sprintf(format, args...))
}

// Warning formats a warning message with a yellow indicator.
func Warning(format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", symbolWarn, fmt.Sprintf(format, args...))
}

// Info formats an informational message with a blue indicator.
func Info(format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", symbolInfo, fmt.Sprintf(format, args...))
}

// Error formats an error message with a red X indicator.
func Error(format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s", symbolErr, fmt.Sprintf(format, args...))
}

// Bold returns the text wrapped in bold formatting.
func Bold(format string, args ...interface{}) string {
	return bold.Sprintf(format, args...)
}
