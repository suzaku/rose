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
	"github.com/suzaku/rose/set"
	"os"

	"github.com/spf13/cobra"
)

var andCmd = &cobra.Command{
	Use:   "and [file1] [file2]",
	Short: "Outputs the intersection of file1 and file2",
	Long: `Outputs lines that appear in both file1 and file2. Both files must be sorted.
If they are not already sorted you can use the sort command. For example:

	rose and <(sort f1) <(sort f2)
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var files [2]*os.File
		for i := 0; i < 2; i++ {
			name := args[i]
			file, err := os.Open(name)
			if err != nil {
				return err
			}
			files[i] = file
		}
		for l := range set.Intersect(files[0], files[1]) {
			fmt.Println(l)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(andCmd)
}
