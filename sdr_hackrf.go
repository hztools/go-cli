// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2022
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

//go:build !sdr.nohackrf

package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"hz.tools/sdr"

	"hz.tools/sdr/hackrf"
)

func init() {
	addSdr(
		"hackrf",
		func(flags *pflag.FlagSet, prefix string) {},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			if err := hackrf.Init(); err != nil {
				return nil, err
			}
			dev, err := hackrf.Open()
			if err != nil {
				return nil, err
			}
			return dev, nil
		},
	)
}

// vim: foldmethod=marker
