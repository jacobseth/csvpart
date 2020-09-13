package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

func main() {
	var headers int
	var filename string
	var whole bool

	app := &cli.App{
		Name: "CSVPart",
		Usage: "Separate a CSV file into smaller ones based on percentage",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "file",
				Usage: "file to be partitioned",
				Destination: &filename,
				Required: true,
			},
				&cli.IntFlag{
					Name: "headers",
					Usage: "number of header lines to duplicate",
					Destination: &headers,
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "whole",
					Usage: "Assume provided percentage is a part of a whole, and fill the remainder",
					Destination: &whole,
					Value: false,
					DefaultText: "false",
				},
			},
		Action: func(c *cli.Context) error {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			lines, err := lineCount(file)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", lines)

			percs, err := parsePercs(c.Args().Slice())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", percs)

			lcAndPercs, err := linesToPercs(lines, percs, headers, whole)
			if err != nil {
				log.Fatal(err)
			}
			for _, t := range lcAndPercs {
				fmt.Printf("lines:%d, perc:%f\n", t.lines, t.perc)
			}

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Line count and percentage
type lcAndPerc struct {
	lines int
	perc float32
}

// Extract line counts for the provided percentages
func linesToPercs(lines int, percs []float32, headers int, whole bool) ([]lcAndPerc, error) {
	realLc := lines - headers
	if realLc < 0 {
		return nil, errors.New("Headers cannot be larger than file's line count")
	}

	lcAndPercs := make([]lcAndPerc, 0)
	lSum := 0

	for _, p := range percs {
		lc :=  int(math.Round(float64((p / 100) * float32(realLc))))
		lSum = lSum + lc

		lcAndPercs = append(lcAndPercs, lcAndPerc{lc, p})
	}

	// If lSum == realLc, the provided percentages account for all the lines already
	if whole && lSum != realLc {
		rem := realLc - lSum
		if lSum > realLc {
			return nil, errors.New("Line sum larger than line count. Attemp without -whole flag")
		}

		remPerc := 100 / (float32(realLc) / float32(rem))
		lcAndPercs = append(lcAndPercs, lcAndPerc{rem, float32(remPerc)})
	}

	return lcAndPercs, nil
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
