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

//go:build !sdr.nopluto

package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"hz.tools/sdr"
	"hz.tools/sdr/pluto"
)

func init() {
	addSdr(
		"pluto",
		func(flags *pflag.FlagSet, prefix string) {
			flags.String(prefix+"pluto-uri", "ip:pluto.local", "plutosdr to connect to")
			flags.Bool(prefix+"pluto-loopback", false, "Set the PlutoSDR BIST Loopback (be sure gain is set low)")

			flags.Uint(prefix+"pluto-kbuf-rx", 0, "Set the number of kernel buffers for the RX channel")
			flags.Uint(prefix+"pluto-kbuf-tx", 0, "Set the number of kernel buffers for the TX channel")
		},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			flags := c.Flags()
			uri, err := flags.GetString(prefix + "pluto-uri")
			if err != nil {
				return nil, err
			}
			loopback, err := flags.GetBool(prefix + "pluto-loopback")
			if err != nil {
				return nil, err
			}
			kbufRx, err := flags.GetUint(prefix + "pluto-kbuf-rx")
			if err != nil {
				return nil, err
			}
			kbufTx, err := flags.GetUint(prefix + "pluto-kbuf-tx")
			if err != nil {
				return nil, err
			}
			p, err := pluto.OpenWithOptions(uri, pluto.Options{
				RxBufferLength:       1024 * 3,
				TxBufferLength:       1024 * 3,
				RxKernelBuffersCount: kbufRx,
				TxKernelBuffersCount: kbufTx,
			})
			if err != nil {
				return nil, err
			}
			if loopback {
				if err := p.SetLoopback(true); err != nil {
					return nil, err
				}
			}
			return p, nil
		},
	)
}

// vim: foldmethod=marker
