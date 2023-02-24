// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2020
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
	"strings"

	"github.com/spf13/pflag"
)

// Turn a flag into an environment variable.
func createEnvName(prefix string, flag *pflag.Flag) string {
	return fmt.Sprintf("%s%s", prefix, strings.Replace(strings.ToUpper(flag.Name), "-", "_", -1))
}

// EnvRegister will set the default values for all flags in the FlagSet to values
// taken from the environment.
func EnvRegister(prefix string, flagSet *pflag.FlagSet) {
	EnvRegisterWithOverrides(prefix, flagSet, map[string]string{})
}

// EnvRegisterWithOverrides will set the default values for all flags in the
// FlagSet to values taken from the environment.
func EnvRegisterWithOverrides(prefix string, flagSet *pflag.FlagSet, overrides map[string]string) {
	flagSet.VisitAll(func(flag *pflag.Flag) {
		envName, ok := overrides[flag.Name]
		if !ok {
			envName = createEnvName(prefix, flag)
		}
		flag.Usage = fmt.Sprintf("%s (${%s})", flag.Usage, envName)
		value := os.Getenv(envName)
		if value == "" {
			return
		}
		flag.Value.Set(value)
	})
}

// vim: foldmethod=marker
