package set

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
	"time"
)

func TestIntersect(t *testing.T) {
	t.Run("No common elements", func(t *testing.T) {
		f1 := strings.NewReader("1\n3\n5\n7\n9")
		f2 := strings.NewReader("2\n4\n6\n8")
		ch := Intersect(f1, f2)
		lines := chToSlice(t, ch)
		assert.Equal(t, lines, []string(nil))
	})
	t.Run("Output common elements", func(t *testing.T) {
		x1 := "a\nb\nc\nd\ne\nf"
		x2 := "b\nd\ne\ng\nh\ni"
		x3 := "e\ni"
		testCases := []struct {
			a string
			b string
			expect []string
		}{
			{a: x1, b: x2, expect: []string{"b", "d", "e"}},
		    {a: x2, b: x1, expect: []string{"b", "d", "e"}},
			{a: x1, b: x3, expect: []string{"e"}},
			{a: x3, b: x2, expect: []string{"e", "i"}},
		}
		for _, tc := range testCases {
			ch := Intersect(strings.NewReader(tc.a), strings.NewReader(tc.b))
			lines := chToSlice(t, ch)
			assert.Equal(t, lines, tc.expect)
		}
	})
}

func chToSlice(t *testing.T, ch <-chan string) []string {
	var result []string
	for {
		timer := time.NewTimer(100 * time.Millisecond)
		select {
		case l, ok := <-ch:
			fmt.Println(l)
			if !ok {
				return result
			}
			result = append(result, l)
		 case <-timer.C:
		 	t.Fatal("Timeout waiting for channel")
		}
		timer.Stop()
	}
}