/*
 * MIT LICENSE
 *
 * Copyright © 2018, G.Ralph Kuntz, MD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/grkuntzmd/pdf-to-go/parser"
)

const maxMemoryFile = 10000000

var (
	buildInfo  string
	buildStamp = "No BuildStamp provided"
	gitHash    = "No GitHash provided"
	version    = "No Version provided"
)

func main() {
	if buildInfo != "" {
		parts := strings.Split(buildInfo, "|")
		if len(parts) >= 3 {
			buildStamp = parts[0]
			gitHash = parts[1]
			version = parts[2]
		}
	}

	flag.CommandLine.Usage = usage
	flag.Parse()

	for _, n := range flag.Args() {
		f, err := os.Open(n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open %s for reading\n", n)
			continue
		}

		s, err := f.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get information for file %s\n", n)
			continue
		}

		// If the file is <= maxMemoryFile, just read it into a memory buffer, otherwise use the file.
		var r parser.ReadReadAtSeeker = f
		if s.Size() <= maxMemoryFile {
			b, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot read %s\n", n)
				continue
			}

			r = bytes.NewReader(b)
		}

		p, err := parser.NewParser(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating PDF parser for %s: %s\n", n, err)
			continue
		}

		err = p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing PDF for %s: %s\n", n, err)
			continue
		}
		fmt.Printf("Document: %+v\n", *p.Document)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s file ...\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "buildStamp: %s, gitHash: %s, version: %s\n", buildStamp, gitHash, version)
	flag.PrintDefaults()
}
