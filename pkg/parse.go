package pkg

import (
	"encoding/csv"
	"strings"
)

func SplitString(input string, sep rune) []string {
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = sep
	fields, _ := reader.Read()
	return fields
}
