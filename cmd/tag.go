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
	"time"

	"github.com/himetani/glstats/repo"
	git "github.com/libgit2/git2go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// countTagCmd represents the countTag command
var tagCmd = &cobra.Command{
	Use:   "tag [repoPath] [tagSubStr]",
	Short: "Count commits having tags (default: by month)",
	Long:  `Count commits having tags (default: by month)`,
}

var (
	all            bool
	countTag       bool
	countCommit    bool
	countInsAndDel bool
)

func init() {
	tagCmd.RunE = tagExec
	tagCmd.Flags().BoolVar(&all, "all", false, "Show all stats(--count-tag, --count-commit, --count-ins-and-del)")
	tagCmd.Flags().BoolVar(&countTag, "count-tag", false, "Show the summary of tag number counted by month ")
	tagCmd.Flags().BoolVar(&countCommit, "count-commit", false, "Show the summary of commit number counted by tag")
	tagCmd.Flags().BoolVar(&countInsAndDel, "count-ins-and-del", false, "Show the summary of insertions and deletions line number counted by tag")
	RootCmd.AddCommand(tagCmd)
}

func tagExec(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Argument are invlid")
	}

	repoPath := args[0]
	tagSubStr := args[1]

	times := GetTimesUntil(time.Now(), duration, MONTH)

	r, _ := git.OpenRepository(repoPath)

	if !(countTag || countCommit || countInsAndDel) {
		all = true
	}

	if all || countTag {
		tagCnts, err := repo.CountTagBy(r, tagSubStr, times)
		if err != nil {
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Month", "Count"})
		for _, tc := range tagCnts {
			table.Append([]string{tc.Time.Format("2006-01"), fmt.Sprint(tc.Cnt)})
		}

		fmt.Println("### Cound tag by month")
		table.Render()
	}

	if all || countCommit {
		taggedCommits, err := repo.CountCommit(r, tagSubStr)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Tag", "Commit Count"})
		for _, tc := range taggedCommits {
			table.Append([]string{fmt.Sprint(tc.Tags), fmt.Sprint(tc.Cnt)})
		}
		fmt.Println("### Count commit number counted by tag")
		table.Render()
	}

	if all || countInsAndDel {
		result, err := repo.CountLine(r, tagSubStr)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Revision", "Tag", "Insertions", "Deletions"})
		for _, tc := range result {
			table.Append([]string{tc.Oid.String(), fmt.Sprint(strings.Join(tc.Tags, ",")), fmt.Sprint(tc.Ins), fmt.Sprint(tc.Del)})
		}
		fmt.Println("### Count insertions and deletions line number counted by tag")
		table.Render()

	}

	return nil
}
