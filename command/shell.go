/*
Copyright (C) 2018 Expedia Group.

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

package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
	"os/exec"
	"syscall"
)

var cmdExec func(string) (string, string, int)

func init() {
	cmdExec = func(script string) (string, string, int) {
		var outbuf, errbuf bytes.Buffer
		var exitCode int

		cmd := exec.Command("bash", "-c", script)
		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf

		err := cmd.Run()
		if err != nil {
			// if failure get exit code
			exitError, _ := err.(*exec.ExitError)
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()

		} else {
			// if success get exit code (will be 0)
			ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		}
		return outbuf.String(), errbuf.String(), exitCode
	}
}

var ShellCommand = flyte.Command{
	Name:         "Shell",
	OutputEvents: []flyte.EventDef{outputEventDef},
	Handler:      shellHandler,
}

func shellHandler(input json.RawMessage) flyte.Event {
	var script string
	if err := json.Unmarshal(input, &script); err != nil {
		log.Print(err)
		return newOutputEvent("", err.Error(), 1)
	}

	if script == "" {
		err := errors.New("No script supplied to handler")
		log.Print(err)
		return newOutputEvent("", err.Error(), 1)
	}

	stdout, stderr, exitCode := cmdExec(script)
	if stderr != "" {
		err := fmt.Errorf("script: %q -> %s; error: -> %v; exitCode: -> %v", script, stdout, stderr, exitCode)
		log.Print(err)
		return newOutputEvent(stdout, stderr, exitCode)
	}

	log.Printf("script: %q -> %s", script, stdout)
	return newOutputEvent(stdout, stderr, exitCode)
}

var outputEventDef = flyte.EventDef{
	Name: "Output",
}

func newOutputEvent(output string, err string, exitCode int) flyte.Event {
	return flyte.Event{
		EventDef: outputEventDef,
		Payload: shellOutputPayload{
			Stdout:   output,
			Stderr:   err,
			ExitCode: exitCode,
		},
	}
}

type shellOutputPayload struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}
