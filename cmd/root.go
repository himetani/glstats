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
	"fmt"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	duration int
	dayFlag  bool
)

type DurationType int

const (
	DAY DurationType = iota
	MONTH
	YEAR
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "glstats",
	Short: "Analyze git log data",
	Long:  `Analyze git log data`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.glstats.yaml)")
	RootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 12, "Duration to analyze (default is 12) ")
	RootCmd.PersistentFlags().BoolVar(&dayFlag, "day", false, "Analyze by day (default is by month)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".glstats")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func GetTimesUntil(until time.Time, num int, durationType DurationType) []time.Time {
	var since time.Time
	var times []time.Time

	switch durationType {
	case MONTH:
		since = until.AddDate(0, -num, 0)
	case DAY:
		since = until.AddDate(0, 0, -num)
	default:
		since = until
	}

	for i := 0; ; i++ {
		var t time.Time
		switch durationType {
		case MONTH:
			t = since.AddDate(0, i, 0)
		case DAY:
			since = since.AddDate(0, 0, i)
		default:
		}

		if t.After(until) || t.Equal(until) {
			break
		}
		times = append(times, t)
	}

	return times
}
