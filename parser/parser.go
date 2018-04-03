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

package parser

import (
	"io"
	"strings"
)

type (
	// Parser contains the working information to parse the PDF file.
	Parser struct {
		ReadReadAtSeeker
		Length int64
		*Document
	}

	// ReadReadAtSeeker combines Reader, ReaderAt, and Seeker, and represents things that can be
	// parsed.
	ReadReadAtSeeker interface {
		io.Reader
		io.ReaderAt
		io.Seeker
	}
)

// NewParser creates a new PDF parser.
func NewParser(r ReadReadAtSeeker) (*Parser, error) {
	// Seek the to end of the PDF to find the length.
	length, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	return &Parser{
		ReadReadAtSeeker: r,
		Length:           length,
		Document:         NewDocument(),
	}, nil
}

// Parse parses the PDF file and builds a Document structure.
func (p *Parser) Parse() error {
	if err := p.parseHeader(); err != nil {
		return err
	}

	if err := p.parseStartXRef(); err != nil {
		return err
	}

	return nil
}

// maxInt64 calculates the maximum of 2 int64 numbers.
func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

// trimEOL removes the end-of-line character(s) CR and LF from the given string.
func trimEOL(s string) string {
	i := strings.IndexAny(s, "\r\n")
	if i >= 0 {
		return s[:i]
	}

	return s
}
