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
	"bytes"
	"os"
	"strings"
	"testing"
)

func runHelper(cli *CLI, args []string, t *testing.T) {
	t.Helper()

	ret := cli.Run(args)

	if ret != ExitOK {
		t.Errorf("Return Code %d is not ExitOK", ret)
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
