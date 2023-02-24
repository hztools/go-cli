// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2021
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// IsPty will check if a file is a CharDevice. This is useful when
// checking os.Stdout to see if it's going to a terminal.
func IsPty(f *os.File) bool {
	fi, _ := f.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// WarnInteractive will os.Exit(1) if:
//
//  - the passed FD is a CharDevice
//  - the --im-a-weirdo flag is false
//
// This will return an error if RegisterInteractive is not called
// on the cobra.Command.
//
func WarnInteractive(cmd *cobra.Command, f *os.File) error {
	isWeird, err := cmd.Flags().GetBool("im-a-weirdo")
	if err != nil {
		return err
	}

	if IsPty(f) && !isWeird {
		fmt.Printf(`
I'm cowardly refusing to write binary to your pretty
terminal. It makes unpleasant glyphs and the terminal
beeps too. No one likes that. Except weirdos.

If you're a weirdo, pass the --im-a-weirdo flag.

`)
		os.Exit(1)
	}
	return nil
}

// RegisterInteractiveFlags will register the --im-a-weirdo override for
// sending IQ data to an interactive PTY.
func RegisterInteractiveFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("im-a-weirdo", false, "set to true if you're weird")
}

// vim: foldmethod=marker
