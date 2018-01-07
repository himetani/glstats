glstats
====

`glstats` is command to show git log stats data.

[![Build Status](https://travis-ci.org/himetani/glstats.svg?branch=master)](https://travis-ci.org/himetani/glstats)

# Usage
```bash
Usage:
  glstats [command]

Available Commands:
  help        Help about any command
  tag         Show stats by tag

Flags:
      --config string   config file (default is $HOME/.glstats.yaml)
      --day             Analyze by day (default is by month)
  -d, --duration int    Duration to analyze (default 12)
  -h, --help            help for glstats
```

# Example
```bash
% ./glstats tag ./glstats-sample-submodule
### Count tag by month
+---------+-------+
|  MONTH  | COUNT |
+---------+-------+
| 2016-01 |     0 |
| 2016-02 |     0 |
| 2016-03 |     0 |
| 2016-04 |     0 |
| 2016-05 |     0 |
| 2016-06 |     0 |
| 2016-07 |     0 |
| 2016-08 |     0 |
| 2016-09 |     0 |
| 2016-10 |     0 |
| 2016-11 |     0 |
| 2016-12 |     0 |
| 2017-01 |     0 |
| 2017-02 |     0 |
| 2017-03 |     0 |
| 2017-04 |     0 |
| 2017-05 |     0 |
| 2017-06 |     0 |
| 2017-07 |     0 |
| 2017-08 |     0 |
| 2017-09 |     0 |
| 2017-10 |     0 |
| 2017-11 |     3 |
| 2017-12 |     2 |
| 2018-01 |     0 |
+---------+-------+
### summary of commit statistics by tag
+------------------------------------------+-----------------+-----------+------------+-----------+
|                 REVISION                 |       TAG       | COMMITNUM | INSERTIONS | DELETIONS |
+------------------------------------------+-----------------+-----------+------------+-----------+
| 65d04b818726554c198866e5f9fbef65d6064a46 | deploy/20180106 |         1 |          0 |         3 |
| 0c254ca47f0924a4d0874c88499315a0987d8d3f | deploy/20180102 |         1 |          0 |         0 |
| e7cec8f34445ef794e54c0d7f6bacef97d99bf5a | deploy/20171124 |         2 |          7 |         0 |
| 264f0767fb0cb4f34eb49d63022a443cefb75783 | deploy/20171123 |         1 |          2 |         0 |
+------------------------------------------+-----------------+-----------+------------+-----------+
```


