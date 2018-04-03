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
	"fmt"
	"io"
	"regexp"
	"strconv"
)

const maxPDFVersion = float32(1.7)

var headerRegex = regexp.MustCompile(`%PDF-(\d.\d)`)

// parseHeader parses the '%PDF-' file header. The location is expected to be at 0.
func (p *Parser) parseHeader() error {
	// Seek to location 0.
	if _, err := p.Seek(0, io.SeekStart); err != nil {
		return err
	}

	r := newReader(p.ReadReadAtSeeker)
	line, err := r.readLine()
	if err != nil {
		return err
	}

	m := headerRegex.FindStringSubmatch(line)
	if m == nil {
		return fmt.Errorf("cannot find PDF header")
	}

	v, err := strconv.ParseFloat(m[1], 32)
	if err != nil {
		return err
	}

	if float32(v) > maxPDFVersion {
		return fmt.Errorf("PDF version exceeds maximum %0.1f", maxPDFVersion)
	}

	p.Document.Version = float32(v)
	return nil
}
