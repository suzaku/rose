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
	"container/heap"
	"fmt"
	"io"
)

type rowSource struct {
	ch  <-chan string
	top string
}

func (rs *rowSource) Top() string {
	if len(rs.top) == 0 {
		var ok bool
		select {
		case rs.top, ok = <-rs.ch:
			if !ok {
				return ""
			}
		}
	}
	return rs.top
}

func (rs *rowSource) Next() bool {
	rs.top = ""
	var ok bool
	select {
	case rs.top, ok = <-rs.ch:
		if !ok {
			return false
		}
	}
	return true
}

func (rs *rowSource) String() string {
	return fmt.Sprintf("<rowSource: top=%s>", rs.top)
}

type rowSources []*rowSource

func (rs *rowSources) Push(x interface{}) {
	*rs = append(*rs, x.(*rowSource))
}

func (rs *rowSources) Pop() interface{} {
	old := *rs
	n := len(old)
	x := old[n-1]
	*rs = old[0 : n-1]
	return x
}

func (rs *rowSources) Len() int {
	return len(*rs)
}

func (rs *rowSources) Less(i, j int) bool {
	return (*rs)[i].Top() < (*rs)[j].Top()
}

func (rs *rowSources) Swap(i, j int) {
	(*rs)[i], (*rs)[j] = (*rs)[j], (*rs)[i]
}

func Union(files ...io.Reader) <-chan string {
	resultCh := make(chan string, 16)

	sources := make(rowSources, len(files))
	for i, f := range files {
		sources[i] = &rowSource{
			ch: readNonEmptyLines(f),
		}
	}

	heap.Init(&sources)

	go func() {
		var lastRow string
		for len(sources) > 0 {
			popped := heap.Pop(&sources)
			src := popped.(*rowSource)
			row := src.Top()
			if row != lastRow {
				resultCh <- row
				lastRow = row
			}
			if src.Next() {
				heap.Push(&sources, src)
			}
		}
		close(resultCh)
	}()

	return resultCh
}
