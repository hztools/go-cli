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
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"runtime/pprof"
)

// Pprof will initialize the profiler. This will only do something if the
// RF_%s_PPROF envvar is set. The string value is the passed name,
// so if "RFUTILS" is passed, the full envvar name is:
//
// RF_RFUTILS_PPROF
//
// This root will than have .cpu and .memory appended to it, for
// the memory/heap profile and CPU profile.
func Pprof(name string) func() {
	pprofPath := os.Getenv(fmt.Sprintf("RF_%s_PPROF", name))
	if pprofPath == "" {
		return func() {}
	}

	cpuCloser, err := PprofCPU(fmt.Sprintf("%s.cpu", pprofPath))
	if err != nil {
		log.WithError(err).Warn("cli.Pprof: error profiling CPU")
		return func() {}
	}

	allCloser, err := PprofAll(pprofPath)
	if err != nil {
		log.WithError(err).Warn("cli.Pprof: error profiling")
		return func() {}
	}

	return func() {
		allCloser()
		cpuCloser()
	}
}

// PprofAll will write the Heap profile
func PprofAll(path string) (func(), error) {
	return func() {
		runtime.GC()
		for _, ptype := range []string{
			"goroutine",
			"heap",
			"allocs",
			"threadcreate",
			"block",
			"mutex",
		} {
			fPath := fmt.Sprintf("%s.%s", path, ptype)
			f, err := os.Create(fPath)
			if err != nil {
				log.Printf("cli.Pprof: Can't create file '%s': %s", fPath, err)
				continue
			}
			pprof.Lookup(ptype).WriteTo(f, 0)
			f.Close()
		}
	}, nil
}

// PprofCPU will create a pprof capture to a specified path.
func PprofCPU(path string) (func(), error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		return nil, err
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}, nil
}

// vim: foldmethod=marker
