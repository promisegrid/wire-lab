package main

import "strings"

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
