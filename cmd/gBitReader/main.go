/*
   Copyright 2019 Takahiro Yamashita
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mattn/go-isatty"
)

const version string = "0.0.1"

// Exit status
const (
	ExitOK int = iota
	ExitArgError
	ExitCmdError
)

type config struct {
	showVersion  bool
	verbose      bool
	terminalMode bool
}

// CLI has In/Out/Err streams.
// Flags is option.
type CLI struct {
	OutStream     io.Writer
	InStream      *os.File
	ErrStream     io.Writer
	Flags         *flag.FlagSet
	forceTerminal bool
}

func (cli *CLI) checkOption(args []string) (*config, error) {
	config := &config{}

	cli.Flags = flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)

	cli.Flags.BoolVar(&config.showVersion, "V", false, "show version")

	cli.Flags.Parse(args[1:])

	config.terminalMode = isatty.IsTerminal(cli.InStream.Fd())
	if cli.forceTerminal {
		// for testing
		config.terminalMode = true
	}

	if config.showVersion {
		return config, nil
	}

	if config.terminalMode && cli.Flags.NArg() == 0 {
		return nil, fmt.Errorf("no files")
	}

	return config, nil
}

// Run executes real main function.
func (cli *CLI) Run(args []string) (ret int) {
	cnf, err := cli.checkOption(args)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "Error:%s\n", err)
		return ExitArgError
	}

	if cnf.showVersion {
		fmt.Fprintf(cli.OutStream, "Ver: %s\n", version)
		return ExitOK
	}

	return ret
}

func main() {
	cli := &CLI{OutStream: os.Stdout, InStream: os.Stdin, ErrStream: os.Stderr}

	os.Exit(cli.Run(os.Args))
}
