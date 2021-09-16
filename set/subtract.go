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
	"bufio"
	"io"
)

func Subtract(f1, f2 io.Reader) <-chan string {
	ch := make(chan string, 16)
	go func() {
		defer close(ch)
		scanner1 := bufio.NewScanner(f1)
		searcher := &rowSearcher{
			chRowsInBulk: readLinesInBulk(f2, 64),
		}
		var lastLine string
		var f2Exhausted bool
		for scanner1.Scan() {
			line := scanner1.Text()
			if len(line) == 0 {
				continue
			}
			if line == lastLine {
				continue
			}
			lastLine = line
			if f2Exhausted {
				ch <- line
				continue
			}
			for {
				found, inRange, exhausted := searcher.Search(line)
				if found {
					break
				}
				if inRange || exhausted {
					ch <- line
					f2Exhausted = exhausted
					break
				}
			}
		}
	}()
	return ch
}