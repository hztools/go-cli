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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"hz.tools/rf"
	"hz.tools/sdr"
)

// RegisterSDRFlagsWithPrefix will register the flags but with a string prefix,
// such as "rx-" or "tx-"
func RegisterSDRFlagsWithPrefix(c *cobra.Command, prefix string) {
	flags := pflag.NewFlagSet("", pflag.ExitOnError)

	flags.String(prefix+"sdr", "rtl", fmt.Sprintf("[%s]", strings.Join(allSdrNames(allSdrConstructors), "|")))

	flags.String(prefix+"gains", "", "NAME=1.0,NAME2=2.5")
	flags.String(prefix+"agc", "", "[on|manual]")

	flags.String(prefix+"frequency", "", "frequency to set the SDR to")
	flags.Uint(prefix+"sample-rate", 2.5e6, "samples per second")

	for _, addSdrFlags := range allSdrFlags {
		addSdrFlags(flags, prefix)
	}

	EnvRegister("RF_", flags)

	c.Flags().AddFlagSet(flags)
}

// RegisterSDRFlags will set the SDR related flags on the cobra.Command's pflag
// object.
func RegisterSDRFlags(c *cobra.Command) {
	RegisterSDRFlagsWithPrefix(c, "")
}

// CreateGainMap will parse the string format GainMap and return it as a
// map of string to float values.
func CreateGainMap(gains string) (map[string]float32, error) {
	if gains == "" {
		return nil, nil
	}

	gainsMap := map[string]float32{}

	for _, gainSetting := range strings.Split(gains, ",") {
		gainKV := strings.Split(gainSetting, "=")
		if len(gainKV) != 2 {
			return nil, fmt.Errorf("Can't parse gain %s", gainSetting)
		}

		gainValue, err := strconv.ParseFloat(gainKV[1], 32)
		if err != nil {
			return nil, err
		}

		gainsMap[gainKV[0]] = float32(gainValue)
	}

	return gainsMap, nil
}

func createGainMap(c *cobra.Command, prefix string) (map[string]float32, error) {
	gains, err := c.Flags().GetString(prefix + "gains")
	if err != nil {
		return nil, err
	}
	return CreateGainMap(gains)
}

// LoadSDR will return an sdr.Sdr defined by the configured CLI flags,
// or an error.
func LoadSDR(c *cobra.Command) (sdr.Sdr, rf.Hz, uint, error) {
	return LoadSDRWithPrefix(c, "")
}

// LoadSDRWithPrefix will return an sdr.Sdr define by the configured CLI flags,
// as well as the provided prefix prepended to the CLI flags.
func LoadSDRWithPrefix(c *cobra.Command, prefix string) (sdr.Sdr, rf.Hz, uint, error) {
	dev, err := loadSDRWithPrefix(c, prefix)
	if err != nil {
		return nil, rf.Hz(0), 0, err
	}

	agc, err := c.Flags().GetString(prefix + "agc")
	if err != nil {
		return nil, rf.Hz(0), 0, err
	}

	switch agc {
	case "manual":
		if err := dev.SetAutomaticGain(false); err != nil {
			return nil, rf.Hz(0), 0, err
		}
	case "on":
		if err := dev.SetAutomaticGain(true); err != nil {
			return nil, rf.Hz(0), 0, err
		}
	case "":
		break
	default:
		return nil, rf.Hz(0), 0, fmt.Errorf("unknown gain mode")
	}

	gainsMap, err := createGainMap(c, prefix)
	if err != nil {
		return nil, rf.Hz(0), 0, err
	}

	if gainsMap != nil {
		if err := sdr.SetGainStages(dev, gainsMap); err != nil {
			return nil, rf.Hz(0), 0, err
		}
	}
	flags := c.Flags()

	sps, err := flags.GetUint(prefix + "sample-rate")
	if err != nil {
		return nil, rf.Hz(0), 0, err
	}

	if err := dev.SetSampleRate(sps); err != nil {
		return nil, rf.Hz(0), 0, err
	}

	rsps, err := dev.GetSampleRate()
	if err == nil {
		sps = rsps
	}

	var frequency rf.Hz
	freqString, err := flags.GetString(prefix + "frequency")
	if err != nil {
		return nil, rf.Hz(0), 0, err
	}

	if freqString != "" {
		frequency, err = rf.ParseHz(freqString)
		if err != nil {
			return nil, rf.Hz(0), 0, err
		}
		if err := dev.SetCenterFrequency(frequency); err != nil {
			return nil, rf.Hz(0), 0, err
		}

		rFrequency, err := dev.GetCenterFrequency()
		if err == nil {
			frequency = rFrequency
		}

		log.WithFields(log.Fields{
			"frequency":      frequency,
			"frequency.band": frequency.ITUBandName(),
		}).Info("Center Frequency set")
	}

	return dev, frequency, sps, nil
}

// sdrConstructor is used internally to register different SDR backends
// into loadSDRWithPrefix without having a massive switch statement
// when invoked by LoadSDR (aka LoadSDRWithPrefix)
type sdrConstructor func(*cobra.Command, string) (sdr.Sdr, error)

// sdrFlagSet is used internally to register CLI flag arguments when
// invoked by RegisterSDRFlags (aka RegisterSDRFlagsWithPrefix)
type sdrFlagSet func(*pflag.FlagSet, string)

var (
	allSdrConstructors = map[string]sdrConstructor{}
	allSdrFlags        = map[string]func(*pflag.FlagSet, string){}
)

func allSdrNames(allSdrs map[string]sdrConstructor) []string {
	ret := []string{}
	for name := range allSdrs {
		ret = append(ret, name)
	}
	return ret
}

func addSdr(name string, flags func(*pflag.FlagSet, string), c sdrConstructor) {
	allSdrFlags[name] = flags
	allSdrConstructors[name] = c
}

// loadSDRWithPrefix will return an sdr.Sdr defined by the configured CLI flags,
// or an error.
func loadSDRWithPrefix(c *cobra.Command, prefix string) (sdr.Sdr, error) {
	flags := c.Flags()

	sdrType, err := flags.GetString(prefix + "sdr")
	if err != nil {
		return nil, err
	}

	sdrConstructor, ok := allSdrConstructors[sdrType]
	if !ok {
		return nil, fmt.Errorf("rfutil: unknown sdr type: %s", sdrType)
	}
	return sdrConstructor(c, prefix)
}

// vim: foldmethod=marker
