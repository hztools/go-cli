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

//go:build !sdr.nortl

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"hz.tools/fftw"
	"hz.tools/sdr"
	"hz.tools/sdr/rtl"
	"hz.tools/sdr/rtl/kerberos"
	"hz.tools/sdr/rtltcp"
)

func init() {

	addSdr(
		"rtl",
		func(flags *pflag.FlagSet, prefix string) {
			flags.String(prefix+"rtl-serial", "", "serial number to use")
			flags.Uint(prefix+"rtl-device-index", 0, "device index to use")
			flags.Bool(prefix+"rtl-bias-t", false, "Set bias-T state")
		},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			flags := c.Flags()
			serial, err := flags.GetString(prefix + "rtl-serial")
			if err != nil {
				return nil, err
			}
			deviceIndex, err := flags.GetUint(prefix + "rtl-device-index")
			if err != nil {
				return nil, err
			}
			if serial != "" {
				if deviceIndex != 0 {
					return nil, fmt.Errorf("rfutil: can't set both serial and index")
				}
				deviceIndex, err = rtl.DeviceIndexBySerial(serial)
				if err != nil {
					return nil, err
				}
			}
			rtlCount := rtl.DeviceCount()
			if rtlCount == 0 {
				return nil, fmt.Errorf("rfutil: no rtl devices found")
			}

			if rtlCount < deviceIndex {
				return nil, fmt.Errorf("rfutil: index isn't valid")
			}
			dev, err := rtl.New(deviceIndex, 0)
			if err != nil {
				return nil, err
			}

			biasT, err := flags.GetBool(prefix + "rtl-bias-t")
			if err != nil {
				return nil, err
			}
			if err := dev.SetBiasT(biasT); err != nil {
				return nil, err
			}
			return dev, nil
		},
	)

	addSdr(
		"rtltcp",
		func(flags *pflag.FlagSet, prefix string) {
			flags.String(prefix+"rtltcp-host", "localhost", "fqdn to connect to")
			flags.Uint(prefix+"rtltcp-port", 1234, "remote port to use")
		},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			flags := c.Flags()
			fqdn, err := flags.GetString(prefix + "rtltcp-host")
			if err != nil {
				return nil, err
			}
			port, err := flags.GetUint(prefix + "rtltcp-port")
			if err != nil {
				return nil, err
			}
			addr := fmt.Sprintf("%s:%d", fqdn, port)
			return rtltcp.Dial("tcp", addr)
		},
	)

	addSdr(
		"kerberos-coherent",
		func(flags *pflag.FlagSet, prefix string) {},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			return kerberos.NewCoherent(fftw.Plan, 0, 1, 2, 3, 0)
		},
	)

	addSdr(
		"kerberos-offset",
		func(flags *pflag.FlagSet, prefix string) {},
		func(c *cobra.Command, prefix string) (sdr.Sdr, error) {
			return kerberos.NewOffset(fftw.Plan, 0, 1, 2, 3, 0)
		},
	)
}

// vim: foldmethod=marker
