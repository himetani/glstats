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
	"time"

	"github.com/himetani/glstats/analyze"
	"github.com/himetani/glstats/timeutil"
	git "github.com/libgit2/git2go"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type TagFlags struct {
	since string
	until string
}

type Tag struct{}

var tagFlags = &TagFlags{}

// countTagCmd represents the countTag command
var tagCmd = &cobra.Command{
	Use:   "tag [repo]",
	Short: "Count commits having tags (default: by month)",
	Long:  `Count commits having tags (default: by month)`,
}

func init() {
	tagCmd.RunE = tagExec
	RootCmd.AddCommand(tagCmd)
}

func tagExec(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Argument are invlid")
	}

	repoPath := args[0]

	times := timeutil.GetTimesUntil(time.Now(), duration, timeutil.MONTH)

	repo, _ := git.OpenRepository(repoPath)

	tagCnts, err := analyze.CountTag(repo, "deploy", times)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Month", "Count"})
	for _, tc := range tagCnts {
		table.Append([]string{tc.Time.Format("2006-01"), fmt.Sprint(tc.Cnt)})
	}
	table.Render()

	return nil
}
