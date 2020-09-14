package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	var headers int
	var filename string
	var whole bool

	app := &cli.App{
		Name:  "CSVPart",
		Usage: "Separate a CSV file into smaller ones based on percentage",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "filename",
				Usage:       "name of file to be partitioned",
				Destination: &filename,
				Required:    true,
			},
			&cli.IntFlag{
				Name:        "headers",
				Usage:       "number of header lines to duplicate",
				Destination: &headers,
				Required:    true,
			},
			&cli.BoolFlag{
				Name:        "whole",
				Usage:       "Assume provided percentage is a part of a whole, and fill the remainder",
				Destination: &whole,
				Value:       false,
				DefaultText: "false",
			},
		},
		Action: func(c *cli.Context) error {
			srcFile, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			lines, err := lineCount(srcFile)
			if err != nil {
				log.Fatal(err)
			}

			percs, err := parsePercs(c.Args().Slice())
			if err != nil {
				log.Fatal(err)
			}

			lcs, err := linesFromPercs(lines, percs, headers, whole)
			if err != nil {
				log.Fatal(err)
			}

			srcBuff := bufio.NewReader(srcFile)
			hBytes := make([]byte, 0)
			srcFile.Seek(0, 0)
			for i := 0; i < headers; i++ {
				h, err := srcBuff.ReadBytes('\n')
				if err != nil {
					log.Fatal(err)
				}
				hBytes = append(hBytes, h...)
			}

			newFiles := make([]string, 0)
			for idx, lc := range lcs {
				fname := fmt.Sprintf("%d_%s.csv", idx, strings.Split(filename, ".")[0])
				f, err := os.Create(fname)
				if err != nil {
					removeFiles(newFiles)
					log.Fatal(err)
				}
				newFiles = append(newFiles, fname)

				f.Write(hBytes)
				for i := 0; i < lc; i++ {
					l, err := srcBuff.ReadBytes('\n')
					if err != nil {
						removeFiles(newFiles)
						log.Fatal(err)
					}
					f.Write(l)
				}
			}

			fmt.Println("Done 👍")
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func removeFiles(fileNames []string) (err error) {
	for _, name := range fileNames {
		err = os.Remove(name)
	}
	return
}

// Extract line counts for the provided percentages
func linesFromPercs(lines int, percs []float32, headers int, whole bool) ([]int, error) {
	realLc := lines - headers
	if realLc < 0 {
		return nil, errors.New("Headers cannot be larger than file's line count")
	}

	lcs := make([]int, 0)
	lSum := 0

	for _, p := range percs {
		lc := int(math.Round(float64((p / 100) * float32(realLc))))
		lSum = lSum + lc

		lcs = append(lcs, lc)
	}

	// If lSum == realLc, the provided percentages account for all the lines already
	if whole && lSum != realLc {
		rem := realLc - lSum
		if lSum > realLc {
			return nil, errors.New("Line sum larger than line count. Attempt without -whole flag")
		}
		lcs = append(lcs, rem)
	}

	return lcs, nil
}

// based on: https://stackoverflow.com/a/24563853
func lineCount(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// Parse percentage arg strings, also check sanity
func parsePercs(pStrings []string) ([]float32, error) {
	pFloats := make([]float32, 0)
	sum := float32(0.0)
	for _, ps := range pStrings {
		pf, err := strconv.ParseFloat(ps, 32)
		if err != nil {
			return nil, err
		}

		sum = sum + float32(pf)
		if sum > 100 {
			return nil, errors.New("Error: Supplied percentages sum larger than 100\n")
		}

		pFloats = append(pFloats, float32(pf))
	}
	return pFloats, nil
}
