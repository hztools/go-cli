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
	"runtime"
	"runtime/debug"
	"strings"

	hzdebug "hz.tools/sdr/debug"

	"github.com/spf13/cobra"
)

func version(cmd *cobra.Command, args []string) error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}
	fmt.Printf("%s %s\n", buildInfo.Main.Path, buildInfo.Main.Version)
	fmt.Printf("\n")

	for _, pkg := range buildInfo.Deps {
		if !strings.HasPrefix(pkg.Path, "hz.tools") {
			continue
		}
		fmt.Printf("  %s %s\n", pkg.Path, pkg.Version)
	}

	bi := hzdebug.ReadBuildInfo()

	fmt.Printf("\n")
	fmt.Printf("Understood IQ formats:\n")
	fmt.Printf("\n")
	for _, format := range bi.SampleFormats {
		fmt.Printf("    - %s\n", format.String())
	}

	fmt.Printf("\n")
	fmt.Printf("Compiled drivers:\n")
	fmt.Printf("\n")
	for _, radioDriver := range bi.RadioDrivers {
		fmt.Printf("    - %s\n", radioDriver)
	}

	fmt.Printf("\n")
	fmt.Printf("Compiled Binary:\n")
	fmt.Printf("\n")
	fmt.Printf("           SIMD: %t\n", bi.SIMD.Enabled)
	if bi.SIMD.Enabled {
		fmt.Printf("  SIMD Backends: %s\n", strings.Join(bi.SIMD.Backends, ", "))
	}
	fmt.Printf("      ByteOrder: %s\n", bi.HostEndianness)
	fmt.Printf("\n")
	fmt.Printf("  Compiled with: %s\n", runtime.Compiler)
	fmt.Printf("         GOARCH: %s\n", runtime.GOARCH)
	fmt.Printf("           GOOS: %s\n", runtime.GOOS)
	fmt.Printf("     Go Version: %s\n", runtime.Version())
	fmt.Printf("\n")

	return nil
}

// RegisterVersionSubcommand will register a standard Version command to the
// Cobra app.
func RegisterVersionSubcommand(rootCmd *cobra.Command) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "enumerate dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return version(cmd, args)
		},
	}

	rootCmd.AddCommand(versionCmd)
	return versionCmd
}

// vim: foldmethod=marker
