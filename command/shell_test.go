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
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkingAsExpected(t *testing.T) {
	cmdExec = func(string) (string, string, int) {
		return "hello\nworld", "", 0
	}

	input := toJson("echo hello; echo world", t)
	actualEvent := shellHandler(input)
	expectedEvent := newOutputEvent("hello\nworld", "", 0)

	assert.Equal(t, expectedEvent, actualEvent)
}

func TestRunScriptNoScript(t *testing.T) {
	cmdExec = func(string) (string, string, int) {
		t.Fatal("Command should not be attempted to be executed if empty")
		return "", "", -1
	}

	input := toJson("", t)
	actualEvent := shellHandler(input)
	expectedEvent := newOutputEvent("", "No script supplied to handler", 1)

	assert.Equal(t, expectedEvent, actualEvent)
}

func TestRunScriptExecError(t *testing.T) {
	cmdExec = func(string) (string, string, int) {
		return "some output", "Exec Error", 1
	}

	input := toJson("foo", t)
	actualEvent := shellHandler(input)
	expectedEvent := newOutputEvent("some output", "Exec Error", 1)

	assert.Equal(t, expectedEvent, actualEvent)
}

func TestNoInput(t *testing.T) {
	cmdExec = func(string) (string, string, int) {
		t.Fatal("If no arguments specified handler should not be called")
		return "", "", -1
	}

	actualEvent := shellHandler(nil)
	expectedEvent := newOutputEvent("", "unexpected end of JSON input", 1)

	assert.Equal(t, expectedEvent, actualEvent)
}

func toJson(i interface{}, t *testing.T) []byte {

	b, err := json.Marshal(i)
	if err != nil {
		t.Errorf("error marshalling: %v", err)
	}
	return b
}
