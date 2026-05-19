package main

import (
	"fmt"
	"io"
	"strings"
)

func joinLines(lines []string) string {
	return strings.Join(lines, "\n") + "\n"
}

func markdownEscapeCell(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", `\|`)
	return strings.TrimSpace(value)
}

func stripCodeCell(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`") && len(value) >= 2 {
		return strings.TrimSuffix(strings.TrimPrefix(value, "`"), "`")
	}
	return value
}

// writeText wraps plain text writes so every CLI output error is checked.
func writeText(writer io.Writer, text string) error {
	_, err := io.WriteString(writer, text)
	return err
}

// writeLine wraps line-oriented output so errcheck covers human-facing status
// messages as well as file writes.
func writeLine(writer io.Writer, values ...any) error {
	_, err := fmt.Fprintln(writer, values...)
	return err
}

// writeFormat wraps formatted output and keeps the command handlers free of
// unchecked fmt.Fprintf calls.
func writeFormat(writer io.Writer, format string, values ...any) error {
	_, err := fmt.Fprintf(writer, format, values...)
	return err
}
