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
	"bufio"
	"io"
)

type reader struct {
	*bufio.Reader
}

// newReader creates and returns a new bufio.Reader wrapper.
func newReader(r io.Reader) reader {
	return reader{bufio.NewReader(r)}
}

// reader reads a line up to an end-of-line sequence (CR, LF, or CRLF) and returns the line,
// including the end-of-line character(s).
func (rdr reader) readLine() (string, error) {
	var runes []rune
	r, _, err := rdr.ReadRune()
	if err != nil {
		return "", err
	}
	for {
		if r == '\r' {
			runes = append(runes, r)
			r, _, err = rdr.ReadRune()
			if err != nil {
				return "", err
			}
			if r == '\n' {
				runes = append(runes, r)
				return string(runes), nil
			}

			if err = rdr.UnreadRune(); err != nil {
				return "", err
			}

			return string(runes), nil
		}

		if r == '\n' {
			runes = append(runes, r)
			return string(runes), nil
		}

		runes = append(runes, r)
		r, _, err = rdr.ReadRune()
		if err != nil {
			return "", err
		}
	}
}
