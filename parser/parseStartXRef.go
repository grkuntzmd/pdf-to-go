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
	"bytes"
	"fmt"
	"strconv"
)

const trailerByteCount = 1024

var (
	eofMarker = []byte("%%EOF")
	startXRef = []byte("startxref")
)

func (p *Parser) parseStartXRef() error {
	// Seek to trailerByteCount bytes from the end of the PDF.
	skip := maxInt64(0, p.Length-trailerByteCount)
	b := make([]byte, trailerByteCount)
	if _, err := p.ReadAt(b, skip); err != nil {
		return err
	}

	i := bytes.LastIndex(b, eofMarker)
	if i < 0 {
		return fmt.Errorf("cannot find PDF end-of-file marker ('%s')", string(eofMarker))
	}

	// Found the %%EOF. Now look for the 'startxref'.
	i = bytes.LastIndex(b[:i], startXRef)
	if i < 0 {
		return fmt.Errorf("cannot find PDF startxref marker")
	}

	// Found 'startxref'. The next line contains the offset.
	r := newReader(bytes.NewReader(b[i:]))
	if _, err := r.readLine(); err != nil {
		return err
	}

	// The line will contain the offset of the cross reference table as a string.
	line, err := r.readLine()
	if err != nil {
		return err
	}
	p.StartXRefOffset, err = strconv.ParseInt(trimEOL(line), 10, 64)

	return err
}
