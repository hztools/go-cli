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

//go:build !sdr.nouhd

package cli

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"hz.tools/sdr"
	"hz.tools/sdr/uhd"
)

func init() {
	addSdr(
		"uhd",
		func(flags *pflag.FlagSet, prefix string) {
			flags.Int(prefix+"uhd-rx-channel", 0, "rx channel to use")
			flags.IntSlice(prefix+"uhd-rx-channels", nil, "rx channels to use")
			flags.Int(prefix+"uhd-tx-channel", 0, "tx channel to use")
			flags.String(prefix+"uhd-sample-format", "i16", "[i8|i16|c64]")
			flags.String(prefix+"uhd-time-source", "", "clock source to use, check UHD docs for help")
			flags.Int(prefix+"uhd-buffer-length", 10, "Set the underlying buffer queue length")
			flags.String(prefix+"uhd-args", "", "underlying uhd arguments to pass to libuhd")
		},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			flags := c.Flags()
			rxChannels, err := flags.GetIntSlice(prefix + "uhd-rx-channels")
			if err != nil {
				return nil, err
			}
			rxChannel, err := flags.GetInt(prefix + "uhd-rx-channel")
			if err != nil {
				return nil, err
			}
			txChannel, err := flags.GetInt(prefix + "uhd-tx-channel")
			if err != nil {
				return nil, err
			}
			sampleFormatStr, err := flags.GetString(prefix + "uhd-sample-format")
			if err != nil {
				return nil, err
			}
			timeSource, err := flags.GetString(prefix + "uhd-time-source")
			if err != nil {
				return nil, err
			}
			bufLength, err := flags.GetInt(prefix + "uhd-buffer-length")
			if err != nil {
				return nil, err
			}

			uhdArgs, err := flags.GetString(prefix + "uhd-args")
			if err != nil {
				return nil, err
			}

			var sampleFormat sdr.SampleFormat
			switch sampleFormatStr {
			case "i8":
				sampleFormat = sdr.SampleFormatI8
			case "i16":
				sampleFormat = sdr.SampleFormatI16
			case "c64":
				sampleFormat = sdr.SampleFormatC64
			default:
				return nil, sdr.ErrSampleFormatUnknown
			}
			_ = bufLength
			uhd, err := uhd.Open(uhd.Options{
				Args:         uhdArgs,
				RxChannels:   rxChannels,
				RxChannel:    rxChannel,
				TxChannel:    txChannel,
				SampleFormat: sampleFormat,
			})
			if err != nil {
				return nil, err
			}
			if timeSource != "" {
				if err := uhd.SetTimeSource(timeSource); err != nil {
					sources, _ := uhd.GetTimeSources()
					log.Printf("Valid clock sources: %#v", sources)
					return nil, err
				}
			}
			return uhd, nil
		},
	)
}

// vim: foldmethod=marker
