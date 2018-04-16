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

package main

import (
	"os"
	"testing"
)

const labels_key = "FLYTE_LABELS"
const host_key = "FLYTE_API_URL"

func Test_getLabels(t *testing.T) {

	os.Setenv(labels_key, "one=1,two =   2 ,  three   =3")
	defer func() { os.Unsetenv(labels_key) }()

	labels, err := getLabels()
	if err != nil {
		t.Error("un-expected error")
	}

	if len(labels) != 3 {
		t.Errorf("expected 3 labels, got %d", len(labels))
	}

	if labels["one"] != "1" || labels["two"] != "2" || labels["three"] != "3" {
		t.Errorf("expected %v labels, got %v", map[string]string{"one": "1", "two": "2", "three": "3"}, labels)
	}
}

func Test_getLabelsNotSet(t *testing.T) {

	labels, err := getLabels()
	if err != nil {
		t.Error("un-expected error")
	}

	if len(labels) != 0 {
		t.Error("un-expected labels")
	}
}

func Test_getLabelsInvalidFormat(t *testing.T) {

	os.Setenv(labels_key, "one: 1, two: 2")
	defer func() { os.Unsetenv(labels_key) }()
	_, err := getLabels()
	if err == nil {
		t.Error("expected error")
	}
}

func Test_getHost(t *testing.T) {

	os.Setenv(host_key, "http://localhost")
	defer func() { os.Unsetenv(host_key) }()

	_, err := getHost()
	if err != nil {
		t.Error("un-expected error")
	}
}

func Test_getHostNotSet(t *testing.T) {

	_, err := getHost()
	if err == nil {
		t.Error("expected error")
	}
}

func Test_getHostInvalidValue(t *testing.T) {

	// missing protocol
	os.Setenv(host_key, ":invalid host")
	defer func() { os.Unsetenv(host_key) }()

	_, err := getHost()
	if err == nil {
		t.Error("expected error")
	}
}
