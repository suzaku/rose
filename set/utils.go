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
	current      []string
}

func (rs *rowSearcher) Search(row string) (found bool, exhausted bool) {
	if len(rs.current) == 0 {
		var ok bool
		if rs.current, ok = <-rs.chRowsInBulk; !ok {
			exhausted = true
			return
		}
	}
	i := sort.SearchStrings(rs.current, row)
	if i >= len(rs.current) {
		rs.current = nil
		return rs.Search(row)
	}
	if rs.current[i] == row {
		found = true
	}
	return
}

func readNonEmptyLines(r io.Reader) <-chan string {
	scanner := bufio.NewScanner(r)
	chLine := make(chan string, 10)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 {
				continue
			}
			chLine <- line
		}
		close(chLine)
	}()
	return chLine
}
