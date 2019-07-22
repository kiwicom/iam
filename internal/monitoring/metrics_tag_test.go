package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		value, expected, msg string
	}{
		{"", "", "should handle empty string"},
		{"r,test", "r_test", "should replace commas"},
		{"r:test", "r_test", "should replace colons"},
		{"Rtest", "rtest", "should convert to lowercase"},
		{"1test", "test", "should start with letter"},
		{"test-/_", "test-/_", "should keep underscores, minuses and slashes"},
		{"r,test:test", "r_test_test", "should replace multiple occurrences"},
		{"r::test_tes🤖t,test", "r__test_tes_t_test", "should sanitize"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, sanitize(test.value), test.msg)
	}
}

func TestTag(t *testing.T) {
	tag := Tag("key", "1not:really,nice_tag🤖right")
	assert.Equal(t, "key:not_really_nice_tag_right", tag, "Tag should call sanitize by default")
}

func TestUnsafeTag(t *testing.T) {
	tag := UnsafeTag("key", "1not:a,nice-tag")
	assert.Equal(t, "key:1not:a,nice-tag", tag, "UnsafeTag should not call sanitize by default")
}
