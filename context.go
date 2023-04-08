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
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func printStack() {
	os.Stderr.Write(stack())
}

func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// Context will return a context.Context tied to CLI application. This will
// use a --timeout flag, if set, and intercept C-c to cancel the context.
func Context(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		log.WithError(err).Warn("Internal Error: RegisterContextFlags was not called on the cobra.Command. Timeouts ignored.")
	}
	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	c := make(chan os.Signal)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	once := sync.Once{}
	go func() {
		for {
			sig := <-c

			once.Do(func() {
				time.Sleep(time.Second * 3)
				log.Warn("Something hung our exit for 3 seconds. Dumping stack")
				printStack()
				os.Exit(1)
			})

			switch sig {
			case syscall.SIGTERM, syscall.SIGINT:
				log.Info("C-c hit, requesting shutdown")
				cancel()
			}
		}
	}()

	return ctx, cancel
}

// RegisterContextFlags will register the --timeout flag for the context.
func RegisterContextFlags(flags *pflag.FlagSet) {
	flags.Duration("timeout", time.Duration(0), "time to wait before requesting exit")
}

// vim: foldmethod=marker
