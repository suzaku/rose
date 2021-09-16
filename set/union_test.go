package set

import (
	"fmt"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestUnion(t *testing.T) {
	t.Run("Interleaving elements", func(t *testing.T) {
		x1 := "a\nc\ne\ng"
		x2 := "b\nd\nf\nh"
		ch := Union(strings.NewReader(x1), strings.NewReader(x2))
		lines := chToSlice(t, ch)
		assert.Equal(t, lines, []string{"a", "b", "c", "d", "e", "f", "g", "h"})
	})
	t.Run("Overlapping elements", func(t *testing.T) {
		x1 := "a\nc\ne\ng"
		x2 := "e\ng\nh\nj"
		ch := Union(strings.NewReader(x1), strings.NewReader(x2))
		lines := chToSlice(t, ch)
		assert.Equal(t, lines, []string{"a", "c", "e", "g", "h", "j"})
	})
	t.Run("Larger files", func(t *testing.T) {
		var sb1 strings.Builder
		for i := 0; i < 500; i++ {
			sb1.WriteString(fmt.Sprintf("%04d", i))
			sb1.WriteRune('\n')
		}
		x1 := sb1.String()

		var sb2 strings.Builder
		for i := 300; i < 600; i++ {
			sb2.WriteString(fmt.Sprintf("%04d", i))
			sb2.WriteRune('\n')
		}
		x2 := sb2.String()

		expect := make([]string, 0, 600)
		for i := 0; i < 600; i++ {
			expect = append(expect, fmt.Sprintf("%04d", i))
		}
		ch := Union(strings.NewReader(x1), strings.NewReader(x2))
		lines := chToSlice(t, ch)
		assert.Equal(t, lines, expect)
	})
}
