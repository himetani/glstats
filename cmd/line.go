// Copyright © 2017 himetani
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/himetani/glstats/analyze"
	git "github.com/libgit2/git2go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var lineCmd = &cobra.Command{
	Use:   "line",
	Short: "Count lines of code by deploy",
	Long:  `Count lines of code by deploy`,
}

func init() {
	lineCmd.RunE = lineExec
	RootCmd.AddCommand(lineCmd)
}

func lineExec(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Argument are invlid")
	}

	repoPath := args[0]

	repo, _ := git.OpenRepository(repoPath)
	result, err := analyze.CountLine(repo, "deploy")
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Revision", "Tag", "Insertions", "Deletions"})
	for _, tc := range result {
		table.Append([]string{tc.Oid.String(), fmt.Sprint(strings.Join(tc.Tags, ",")), fmt.Sprint(tc.Ins), fmt.Sprint(tc.Del)})
	}
	table.Render()

	return nil
}
