package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestShowVersion(t *testing.T) {
	out, err := ioutil.TempFile("", "test_readbit")
	if err != nil {
		t.Fatalf("ioutil.TempFile error:%s", err)
	}
	defer os.Remove(out.Name())

	cli := &CLI{OutStream: out}
	args := []string{"hoge", "-V"}
	ret := cli.Run(args)

	if ret != ExitOK {
		t.Errorf("Return Code %d is not Success", ret)
	}

	retOutput := make([]byte, len(version)+8)

	n, err := out.Read(retOutput)
	if err != io.EOF {
		t.Errorf("File.Read error:%s, n=%d", err, n)
	}
	if strings.Contains(string(retOutput), version) {
		t.Errorf("Version Error. got %s want %s", string(retOutput), version)
	}
}
