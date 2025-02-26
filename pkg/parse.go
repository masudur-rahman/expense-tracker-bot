package pkg

import (
	"encoding/csv"
	"strings"
)

func SplitString(input string, sep rune) []string {
	input = strings.ReplaceAll(input, "”", "\"")
	input = strings.ReplaceAll(input, "“", "\"")
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = sep
	fields, _ := reader.Read()
	return fields
}
