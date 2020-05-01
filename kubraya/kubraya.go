package kubraya

import (
	"strings"
)

// KubrayaSeparator is the default separator between kubraya words
var KubrayaSeparator = "_"

// IsKubraya tells if a word is a kubraya
func IsKubraya(str string) bool {
	return strings.Contains(str, KubrayaSeparator)
}

// SplitKubraya splits Kubraya in chunks
func SplitKubraya(kubraya string) []string {
	return strings.Split(kubraya, KubrayaSeparator)
}
