package set

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSubtract(t *testing.T) {
	t.Run("Nothing left", func(t *testing.T) {
		testCases := []struct {
			a string
			b string
		}{
			{a: "", b: "aa\nb"},
			{a: "a1\na2\ng9", b: "a1\na2\na3\nb1\ng9"},
		}
		for _, c := range testCases {
			ch := Subtract(strings.NewReader(c.a), strings.NewReader(c.b))
			lines := chToSlice(t, ch)
			assert.Empty(t, lines)
		}
	})
	t.Run("Only in first file", func(t *testing.T) {
		testCases := []struct {
			a string
			b string
			expect []string
		}{
			{a: "a\nb\nc\nd", b: "c\nd\ne", expect: []string{"a", "b"}},
			{a: "x\ny\nz\n", b: "c\nd\ne\ny", expect: []string{"x", "z"}},
		}
		for _, tc := range testCases {
			ch := Subtract(strings.NewReader(tc.a), strings.NewReader(tc.b))
			lines := chToSlice(t, ch)
			assert.Equal(t, lines, tc.expect)
		}
	})
}
