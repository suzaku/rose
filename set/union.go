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
	"fmt"
	"os"
)

func Union(f1, f2 *os.File) {
	scanner1 := bufio.NewScanner(f1)
	chLine1 := make(chan string, 10)
	var tmp1 string
	go func() {
		for scanner1.Scan() {
			line := scanner1.Text()
			if len(line) == 0 {
				continue
			}
			chLine1 <- line
		}
		close(chLine1)
	}()
	scanner2 := bufio.NewScanner(f2)
	chLine2 := make(chan string, 10)
	var tmp2 string
	go func() {
		for scanner2.Scan() {
			line := scanner2.Text()
			if len(line) == 0 {
				continue
			}
			chLine2 <- line
		}
		close(chLine2)
	}()
	var lastLine string
outer:
	for {
		var l1, l2 string
		if len(tmp1) > 0 {
			l1, tmp1 = tmp1, ""
		} else {
			var ok bool
			select {
			case l1, ok = <-chLine1:
				if !ok {
					break outer
				}
			}
		}

		if len(tmp2) > 0 {
			l2, tmp2 = tmp2, ""
		} else {
			var ok bool
			select {
			case l2, ok = <-chLine2:
				if !ok {
					break outer
				}
			}
		}

		if l1 < l2 {
			if l1 != lastLine {
				fmt.Println(l1)
				lastLine = l1
			}
			tmp2 = l2
		} else {
			if l2 != lastLine {
				fmt.Println(l2)
				lastLine = l2
			}
			tmp1 = l1
		}
	}
	if len(tmp1) > 0 {
		fmt.Println(tmp1)
		lastLine = tmp1
	}
	if len(tmp2) > 0 {
		fmt.Println(tmp2)
		lastLine = tmp2
	}
	for l := range chLine1 {
		if l != lastLine {
			fmt.Println(l)
			lastLine = l
		}
	}
	for l := range chLine2 {
		if l != lastLine {
			fmt.Println(l)
			lastLine = l
		}
	}
}