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

//go:build !sdr.noairspyhf

package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"hz.tools/sdr"

	"hz.tools/sdr/airspyhf"
)

func init() {
	addSdr(
		"airspyhf",
		func(flags *pflag.FlagSet, prefix string) {
			flags.Uint64(prefix+"airspy-serial", 0, "device serial to use")
			flags.Bool(prefix+"airspy-dsp", false, "Enable or disable Airspy DSP")
		},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			flags := c.Flags()
			serial, err := flags.GetUint64(prefix + "airspy-serial")
			if err != nil {
				return nil, err
			}

			setDsp, err := flags.GetBool(prefix + "airspy-dsp")
			if err != nil {
				return nil, err
			}
			var dev *airspyhf.Sdr
			if serial == 0 {
				dev, err = airspyhf.Open()
			} else {
				dev, err = airspyhf.OpenBySerial(serial)
			}
			if err != nil {
				return nil, err
			}

			if err := dev.SetDSP(setDsp); err != nil {
				return nil, err
			}

			return dev, nil
		},
	)
}

// vim: foldmethod=marker
