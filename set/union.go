/*
Copyright © 2021 Xuecong Liao <satorulogic@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package set

import (
	"io"
)

func Union(f1, f2 io.Reader) <-chan string {
	ch := make(chan string, 16)
	chLine1 := readNonEmptyLines(f1)
	chLine2 := readNonEmptyLines(f2)

	go func() {
		defer close(ch)
		var tmp1, tmp2 string
		var lastLine string
	outer:
		for {
			var l1, l2 string
			if len(tmp1) > 0 {
				l1, tmp1 = tmp1, ""
			} else {
				var ok bool
				if l1, ok = <-chLine1; !ok {
					break outer
				}
			}

			if len(tmp2) > 0 {
				l2, tmp2 = tmp2, ""
			} else {
				var ok bool
				if l2, ok = <-chLine2; !ok {
					break outer
				}
			}

			if l1 < l2 {
				if l1 != lastLine {
					ch <- l1
					lastLine = l1
				}
				tmp2 = l2
			} else {
				if l2 != lastLine {
					ch <- l2
					lastLine = l2
				}
				tmp1 = l1
			}
		}
		if len(tmp1) > 0 {
			ch <- tmp1
			lastLine = tmp1
		}
		if len(tmp2) > 0 {
			ch <- tmp2
			lastLine = tmp2
		}
		for l := range chLine1 {
			if l != lastLine {
				ch <- l
				lastLine = l
			}
		}
		for l := range chLine2 {
			if l != lastLine {
				ch <- l
				lastLine = l
			}
		}
	}()
	return ch
}
