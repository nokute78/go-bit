package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func runHelper(cli *CLI, args []string, t *testing.T) {
	t.Helper()

	ret := cli.Run(args)

	if ret != ExitOK {
		t.Errorf("Return Code %d is not Success", ret)
	}
}

func TestReadFiles(t *testing.T) {
	tempfile, err := ioutil.TempFile("", "TestReadFiles")
	if err != nil {
		t.Fatalf("ioutil.TempFile error: %s", err)
	}
	defer os.Remove(tempfile.Name())

	n, err := tempfile.Write([]byte{0xff, 0x07})
	if err != nil {
		t.Fatalf("File.Write error: %s, n=%d", err, n)
	}

	retOutput := make([]byte, 16)
	out := bytes.NewBuffer(retOutput)

	cli := &CLI{OutStream: out, ErrStream: os.Stderr, forceTerminal: true}
	runHelper(cli, []string{"hoge", "-s", "11", tempfile.Name()}, t)

	if strings.Contains(string(retOutput), "0xff07") {
		t.Errorf("ReadFile Error. got %s want %s", string(retOutput), "0xff07")
	}
}

func TestShowVersion(t *testing.T) {
	retOutput := make([]byte, len(version)+8)
	out := bytes.NewBuffer(retOutput)

	cli := &CLI{OutStream: out, ErrStream: os.Stderr}
	runHelper(cli, []string{"hoge", "-V"}, t)

	if strings.Contains(string(retOutput), version) {
		t.Errorf("Version Error. got %s want %s", string(retOutput), version)
	}
}
