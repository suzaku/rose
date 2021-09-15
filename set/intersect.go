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
	"sort"
)

func Intersect(f1, f2 io.Reader) <-chan string {
	ch := make(chan string, 16)
	scanner1 := bufio.NewScanner(f1)
	searcher := &rowSearcher{
		chRowsInBulk: readLinesInBulk(f2, 64),
	}
	go func() {
		defer close(ch)
		var lastLine string
		for scanner1.Scan() {
			line := scanner1.Text()
			if len(line) == 0 {
				continue
			}
			if line == lastLine {
				continue
			}
			lastLine = line
			for {
				found, inRange, exhausted := searcher.Search(line)
				if found {
					ch <- line
					break
				}
				if inRange {
					break
				}
				if exhausted {
					return
				}
			}
		}
	}()
	return ch
}

func readLinesInBulk(reader io.Reader, bulkSize int) <-chan []string {
	scanner2 := bufio.NewScanner(reader)
	chReadLines2 := make(chan []string, 2)
	go func() {
		readLines2 := make([]string, 0, bulkSize)
		for scanner2.Scan() {
			line := scanner2.Text()
			if len(line) == 0 {
				continue
			}
			readLines2 = append(readLines2, line)
			if len(readLines2) >= bulkSize {
				chReadLines2 <- readLines2
				readLines2 = make([]string, 0, bulkSize)
			}
		}
		if len(readLines2) > 0 {
			chReadLines2 <- readLines2
		}
		close(chReadLines2)
	}()
	return chReadLines2
}

type rowSearcher struct {
	chRowsInBulk <-chan []string
	current []string
}

func (rs *rowSearcher) Search(row string) (found bool, inRange bool, exhausted bool) {
	if len(rs.current) == 0 {
		var ok bool
		if rs.current, ok = <-rs.chRowsInBulk; !ok {
		  exhausted = true
		}
	}
	i := sort.SearchStrings(rs.current, row)
	if i < len(rs.current) {
		inRange = true
		if rs.current[i] == row {
			found = true
		}
	} else {
		rs.current = nil
	}
	return
}
