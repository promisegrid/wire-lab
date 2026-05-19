package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func writeText(writer io.Writer, text string) error {
	_, err := io.WriteString(writer, text)
	return err
}

func writeLine(writer io.Writer, text string) error {
	_, err := fmt.Fprintln(writer, text)
	return err
}

func writeFormat(writer io.Writer, format string, values ...any) error {
	_, err := fmt.Fprintf(writer, format, values...)
	return err
}

func ensureParent(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0o755)
}
