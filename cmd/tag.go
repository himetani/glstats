package cmd

import (
	"fmt"
	"os"

	"github.com/himetani/glstats/git"
	"github.com/himetani/glstats/timeutil"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type TagFlags struct {
	since string
	until string
	//rname        string
	//input        string
	//output       string
	//substr       string
	//since        time.Time
	//until        time.Time
	repo string
}

type Tag struct{}

var tagFlags = &TagFlags{}

// countTagCmd represents the countTag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Count commits having tags by month",
	Long: `Count commits having tags by month

`,
	Run: func(cmd *cobra.Command, args []string) {
		tag := &Tag{}
		tag.execute()
	},
}

func init() {
	tagCmd.Flags().StringVarP(&tagFlags.repo, "repo", "r", "", "Specify the git path")
	tagCmd.Flags().StringVarP(&tagFlags.since, "since", "s", "", "Since date to be analyzed. Format is YYYY-MM-DD(default: 2014-01-01")
	tagCmd.Flags().StringVarP(&tagFlags.until, "until", "u", "", "Until date to be analyzed. Format is YYYY-MM-DD(default: now)")
	RootCmd.AddCommand(tagCmd)
}

func (c *Tag) execute() {
	times, err := timeutil.Divide(tagFlags.since, tagFlags.until, timeutil.MONTH)
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
	}

	analyzer := &git.TagAnalyzer{
		Path: tagFlags.repo,
	}

	tagCnts, err := analyzer.Analyze("deploy", times)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Month", "Count"})
	for _, tc := range tagCnts {
		table.Append([]string{tc.Time.Format("2006-01"), fmt.Sprint(tc.Cnt)})
	}
	table.Render()
}
