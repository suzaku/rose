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
	"github.com/suzaku/rose/set"
	"os"

	"github.com/spf13/cobra"
)

// subCmd represents the sub command
var subCmd = &cobra.Command{
	Use:   "sub [file1] [file2]",
	Short: "Outputs lines that appear in file1 but not in file2",
	Long: `Outputs lines that appear in file1 but not in file2.Both files must be sorted.
If they are not already sorted you can use the sort command. For example:

	rose sub <(sort f1) <(sort f2)`,
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
		set.Subtract(files[0], files[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(subCmd)
}
