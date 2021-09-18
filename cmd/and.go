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

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/suzaku/rose/set"

	"github.com/spf13/cobra"
)

var andCmd = &cobra.Command{
	Use:   "and [file1] [file2] [file3 ...]",
	Short: "Outputs the intersection of two or more files.",
	Long: `Outputs lines that appear in all the specified files. Both files must be sorted.
If they are not already sorted you can use the sort command. For example:

	rose and <(sort f1) <(sort f2)
`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		files := make([]io.Reader, len(args))
		for i := 0; i < len(args); i++ {
			name := args[i]
			file, err := os.Open(name)
			if err != nil {
				return err
			}
			files[i] = file
		}
		for l := range set.Intersect(files[0], files[1], files[2:]...) {
			fmt.Println(l)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(andCmd)
}
