[![Go Report Card](https://goreportcard.com/badge/github.com/jacobseth/csvpart)](https://goreportcard.com/report/github.com/jacobseth/csvpart)

# CSVPart

Golang CSV file splitter that allows you to partition by percentages

## Installation

As this tool is written in golang, you will need to
[download and install golang](https://golang.org/doc/install).

To build and install, use the `go get` tool:

    go get github.com/jacobseth/csvpart

## Basic usage

CSVpart will read a file from a provided filename, and partition it into a
number of files based on the provided percentages.

``` bash
NAME:
   CSVPart - Separate a CSV file into smaller ones based on percentage

USAGE:
   csvpart [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value     file to be partitioned
   --headers value  number of header lines to duplicate (default: 0)
   --whole          Assume provided percentage is a part of a whole, and fill the remainder (default: false)
   --help, -h       show help (default: false)
```

Resultant files are named based on the index of their respective percentages, in
the following format: `%d_some_data.csv`.

## Examples
### Split a CSV file 60% - 40%

``` bash
csvpart -headers 1 -file some_data.csv 60 40
```

### Split a CSV file 60% - ~40% using -whole flag

``` bash
csvpart -headers 1 -file some_data.csv 60
```

---

Pull/Feature requests welcome! :)
