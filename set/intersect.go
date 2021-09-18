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

func Intersect(f1, f2 io.Reader) <-chan string {
	ch := make(chan string, 16)
	searcher := &rowSearcher{
		chRowsInBulk: readLinesInBulk(f2, 64),
	}
	go func() {
		defer close(ch)
		var lastLine string
		for line := range readNonEmptyLines(f1) {
			if line == lastLine {
				continue
			}
			lastLine = line
			for {
				found, exhausted := searcher.Search(line)
				if found {
					ch <- line
				}
				if exhausted {
					return
				}
				break
			}
		}
	}()
	return ch
}
