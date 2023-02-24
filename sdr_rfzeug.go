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

//go:build !sdr.norfzeug

package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"hz.tools/rfzeug/client"
	"hz.tools/sdr"
)

func init() {
	addSdr(
		"rfzeug",
		func(flags *pflag.FlagSet, prefix string) {
			flags.String(prefix+"rfzeug-uri", "http://localhost:3823", "rfzeug endpoint to talk to")
			flags.String(prefix+"rfzeug-host", "", "rfzeug host to resolve SRV records for")
			flags.String(prefix+"rfzeug-radio", "default", "radio to talk to")
			flags.String(prefix+"rfzeug-api-key", "", "API key to talk to the server with")
			flags.Int(prefix+"rfzeug-tos", 0, "IP TOS to set on connections")
		},
		func(cmd *cobra.Command, prefix string) (sdr.Sdr, error) {
			flags := cmd.Flags()
			host, err := flags.GetString(prefix + "rfzeug-host")
			if err != nil {
				return nil, err
			}
			endpoint, err := flags.GetString(prefix + "rfzeug-uri")
			if err != nil {
				return nil, err
			}
			tos, err := flags.GetInt(prefix + "rfzeug-tos")
			if err != nil {
				return nil, err
			}

			var resolveSrv bool
			if host != "" {
				// if --rfzeug-host is set, we're going to use that, not the
				// provided --rfzeug-uri and then look up the SRV record.
				endpoint = host
				resolveSrv = true
			}

			radio, err := flags.GetString(prefix + "rfzeug-radio")
			if err != nil {
				return nil, err
			}
			var auth client.Auth
			apiKey, err := flags.GetString(prefix + "rfzeug-api-key")
			if err != nil {
				return nil, err
			}
			if apiKey != "" {
				auth = client.BearerAuth{Token: apiKey}
			}
			c, err := client.New(client.Options{
				Auth:       auth,
				Endpoint:   endpoint,
				ResolveSrv: resolveSrv,
				TOS:        tos,
			})
			if err != nil {
				return nil, err
			}
			return c.Radio(radio).Sdr()
		},
	)
}

// vim: foldmethod=marker
