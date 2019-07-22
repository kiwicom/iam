package monitoring

import (
	"strings"
	"unicode"
)

func sanitize(value string) string {
	// Early return if the string is empty
	if value == "" {
		return value
	}

	// Must start with a letter, trim the string until we find a letter
	value = strings.TrimLeftFunc(value, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	// May contain alphanumerics, '_', '-', ':', ',', and '/'.
	// Other characters are converted to underscores.
	// We also replace ':', ',' because this is the value of tag : and , would split it
	value = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r):
			return unicode.ToLower(r)
		case unicode.IsNumber(r) || r == '-' || r == '_' || r == '/':
			return r
		default:
			return '_'
		}
	}, value)

	return value
}

// Tag will create from given key and value a tag for statsd metrics in a format key:value
// value will be sanitized by default:
// - value will be lowercase, starting with letter
// - value will contain only alphanumerics and underscores, commas and slashes
// sanitizing can be skipped by using UnsafeTag instead
func Tag(key, value string) string {
	return UnsafeTag(key, sanitize(value))
}

// UnsafeTag will create from given key and value a tag for statsd metrics in a format key:value.
// skips sanitization and as such should be used only with internal values that we know will be safe
func UnsafeTag(key, value string) string {
	var tag strings.Builder
	tag.Grow(1 + len(key) + len(value))
	tag.WriteString(key)
	tag.WriteRune(':')
	tag.WriteString(value)
	return tag.String()
}
