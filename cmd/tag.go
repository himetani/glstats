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

	"github.com/himetani/glstats/stats"
	git "github.com/libgit2/git2go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// countTagCmd represents the countTag command
var tagCmd = &cobra.Command{
	Use:   "tag [repoPath]",
	Short: "Show stats analyzed by tag",
	Long:  `Show stats analyzed by tag`,
}

var (
	all        bool
	count      bool
	commitStat bool
	tagSubstr  string
)

func init() {
	tagCmd.RunE = tagExec
	tagCmd.Flags().BoolVar(&all, "all", false, "Show all stats(--count-tag, --count-commit, --count-ins-and-del)")
	tagCmd.Flags().BoolVar(&count, "count", false, "Show the summary of tag number counted by month ")
	tagCmd.Flags().BoolVar(&commitStat, "commit-stat", false, "Show the summary of commit statistics by tag")
	tagCmd.Flags().StringVarP(&tagSubstr, "tag-substr", "s", "", "tag substring to analyze")
	RootCmd.AddCommand(tagCmd)
}

func tagExec(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Argument are invlid")
	}

	repoPath := args[0]

	times := GetTimesUntil(time.Now(), duration, MONTH)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return err
	}

	if !(count || commitStat) {
		all = true
	}

	if all || count {
		tagCnts, err := stats.CountTagBy(repo, tagSubstr, times)
		if err != nil {
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Month", "Count"})
		for _, tc := range tagCnts {
			table.Append([]string{tc.Time.Format("2006-01"), fmt.Sprint(tc.Cnt)})
		}

		fmt.Println("### Count tag by month")
		table.Render()
	}

	if !(all || commitStat) {
		return nil
	}

	taggedCommitMap, err := stats.GetTaggedCommitMap(repo, tagSubstr)
	if err != nil {
		return err
	}
	stats, err := stats.GetStats(repo, taggedCommitMap)
	if err != nil {
		return err
	}

	if all || commitStat {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Revision", "Tag", "CommitNum", "Insertions", "Deletions"})
		for _, s := range stats {
			table.Append([]string{s.Revision, fmt.Sprint(strings.Join(s.Tags, ",")), fmt.Sprint(s.Cnt), fmt.Sprint(s.Ins), fmt.Sprint(s.Del)})
		}
		fmt.Println("### summary of commit statistics summaried by tag")
		table.Render()

	}

	return nil
}
