package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mattn/go-isatty"
	"github.com/nokute78/go-bit/pkg/bit"
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
	byte         uint64
	bit          uint64
	bitsize      uint64
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

func (cli *CLI) showBits(in string, b []byte, cnf *config) {
	if cnf.verbose {
		fmt.Fprintf(cli.OutStream, "%s (Byte:%d,Bit:%d,Size:%d): 0x%x\n", in, cnf.byte, cnf.bit, cnf.bitsize, b)
	} else {
		fmt.Fprintf(cli.OutStream, "0x%x\n", b)
	}
}

func (cli *CLI) readBits(in io.Reader, cnf *config) ([]byte, error) {
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return []byte{}, fmt.Errorf("ioutil.ReadAll: %s", err)
	}

	ret, err := bit.GetBits(buf, bit.Offset{Byte: cnf.byte, Bit: cnf.bit}, cnf.bitsize)
	if err != nil {
		return []byte{}, fmt.Errorf("GetBits: %s, len=%d size=%d", err, len(buf), cnf.bitsize)
	}
	return ret, nil
}

func (cli *CLI) readStdin(cnf *config) int {
	buf, err := cli.readBits(cli.InStream, cnf)
	if err != nil {
		fmt.Fprintf(cli.ErrStream, "readStdin :%s\n", err)
		return ExitCmdError
	}

	cli.showBits("(stdin)", buf, cnf)

	return ExitOK
}

func (cli *CLI) readFiles(files []string, cnf *config) (ret int) {
	ret = ExitCmdError
	for _, v := range files {
		f, err := os.Open(v)
		if err != nil {
			fmt.Fprintf(cli.ErrStream, "os.Open :%s\n", err)
			continue
		}
		defer f.Close()

		buf, err := cli.readBits(f, cnf)
		if err != nil {
			fmt.Fprintf(cli.ErrStream, "readFiles :%s\n", err)
			continue
		}
		cli.showBits(v, buf, cnf)
		ret = ExitOK
	}

	return ret
}

func (cli *CLI) checkOption(args []string) (*config, error) {
	config := &config{}

	cli.Flags = flag.NewFlagSet(filepath.Base(args[0]), flag.ExitOnError)

	cli.Flags.BoolVar(&config.verbose, "v", false, "verbose mode")
	cli.Flags.BoolVar(&config.showVersion, "V", false, "show version")
	cli.Flags.Uint64Var(&config.byte, "B", 0, "offset (in byte)")
	cli.Flags.Uint64Var(&config.bit, "b", 0, "offset (in bit)")
	cli.Flags.Uint64Var(&config.bitsize, "s", 0, "read size(in bit)")

	cli.Flags.Parse(args[1:])

	config.terminalMode = isatty.IsTerminal(cli.InStream.Fd())
	if cli.forceTerminal {
		// for testing
		config.terminalMode = true
	}

	if config.showVersion {
		return config, nil
	}

	if config.bitsize == 0 {
		return nil, fmt.Errorf("read size is 0")
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

	if cnf.terminalMode {
		ret = cli.readFiles(cli.Flags.Args(), cnf)
	} else {
		ret = cli.readStdin(cnf)
	}

	return ret
}

func main() {
	cli := &CLI{OutStream: os.Stdout, InStream: os.Stdin, ErrStream: os.Stderr}

	os.Exit(cli.Run(os.Args))
}
