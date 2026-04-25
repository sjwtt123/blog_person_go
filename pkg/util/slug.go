package util

import (
	"regexp"
	"strings"
)

var nonSlugChar = regexp.MustCompile(`[^a-z0-9\-]+`)
var multiDash = regexp.MustCompile(`-+`)

func Slugify(s string) string {
	v := strings.ToLower(strings.TrimSpace(s))
	v = strings.ReplaceAll(v, " ", "-")
	v = nonSlugChar.ReplaceAllString(v, "-")
	v = multiDash.ReplaceAllString(v, "-")
	v = strings.Trim(v, "-")
	return v
}
