package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

func main() {
	var headers int
	var filename string

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
			fmt.Printf("%v", lines)

			percs, err := parsePercs(c.Args().Slice())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v", percs)

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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
	pFloats := make([]float32, 0, 0)
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
