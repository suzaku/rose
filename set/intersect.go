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

func Intersect(f1, f2 io.Reader, others ...io.Reader) <-chan string {
	ch := make(chan string, 16)
	const bulkSize = 64
	rowSets := make([]*rowSearcher, 0, len(others)+1)
	rowSets = append(rowSets, &rowSearcher{chRowsInBulk: readLinesInBulk(f2, bulkSize)})
	for _, f := range others {
		rowSets = append(rowSets, &rowSearcher{chRowsInBulk: readLinesInBulk(f, bulkSize)})
	}
	go func() {
		defer close(ch)
		var lastLine string
		for line := range readNonEmptyLines(f1) {
			if line == lastLine {
				continue
			}
			lastLine = line
			foundInAll := true
			for _, set := range rowSets {
				found, exhausted := set.Search(line)
				if exhausted {
					return
				}
				if !found {
					foundInAll = false
					break
				}
			}
			if foundInAll {
				ch <- line
			}
		}
	}()
	return ch
}
