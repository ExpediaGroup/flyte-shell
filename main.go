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
	"errors"
	"fmt"
	"github.com/ExpediaGroup/flyte-shell/command"
	"github.com/HotelsDotCom/flyte-client/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {

	host, err := getHost()
	if err != nil {
		log.Fatal(err)
	}

	labels, err := getLabels()
	if err != nil {
		log.Fatal(err)
	}

	helpUrl, err := url.Parse("https://github.com/ExpediaGroup/flyte-shell/blob/master/README.md")
	if err != nil {
		log.Fatal(err)
	}

	packDef := flyte.PackDef{
		Name:     "Shell",
		Labels:   labels,
		HelpURL:  helpUrl,
		Commands: []flyte.Command{command.ShellCommand},
	}

	p := flyte.NewPack(packDef, client.NewClient(host, 10*time.Second))
	p.Start()
	// Sleeps forever
	select {}
}

func getHost() (*url.URL, error) {

	host := os.Getenv("FLYTE_API_URL")
	if host == "" {
		return nil, errors.New("FLYTE_API_URL env. variable is not set")
	}
	return url.Parse(host)
}

func getLabels() (map[string]string, error) {

	labelsString := os.Getenv("FLYTE_LABELS")
	labels := make(map[string]string)

	if labelsString == "" {
		return labels, nil
	}

	// labels format: 'key=value,key=value'
	for _, label := range strings.Split(labelsString, ",") {
		items := strings.SplitN(label, "=", 2)
		if len(items) != 2 {
			return nil, fmt.Errorf("invalid format of LABELS env. variable %q", labelsString)
		}
		labels[strings.TrimSpace(items[0])] = strings.TrimSpace(items[1])
	}
	return labels, nil
}
