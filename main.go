// zsh-unicode-fix reads UTF8 text from os.Stdin, transforms it using the
// following rules, then writes the results to os.Stdout.
//
//	1. Lines where len(s) == utf8.RuneCountInString(s) are unchanged,
//	   otherwise each rune of a line is examined in turn.
//	2. Runes of only 1 byte are left unchanged.
//	3. Runes of 2 or 3 bytes are reformatted using `$'\u%04X'`.
//	4. Runes longer than 3 bytes are reformatted using `$'\U%X'`
//

// Copyright © 2024 Timothy E. Peoples
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
}

func run() error {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		fmt.Println(encodeString(s.Text()))
	}

	return s.Err()
}

func encodeString(s string) string {
	if len(s) == utf8.RuneCountInString(s) {
		// s contains no UTF8 runes; leave it alone.
		return s
	}

	var sb strings.Builder

	for _, r := range []rune(s) {
		var fmtstr string
		switch utf8.RuneLen(r) {
		case 1:
			sb.WriteRune(r)
			continue

		case 2, 3:
			fmtstr = `$'\u%04X'`

		default:
			fmtstr = `$'\U%X'`
		}

		fmt.Fprintf(&sb, fmtstr, r)
	}

	return sb.String()
}
