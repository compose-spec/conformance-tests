/*
   Copyright 2020 The Compose Specification Authors.

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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/icmd"
)

const (
	commandsDir           = "commands"
	baseSpecReference     = "https://github.com/compose-spec/compose-spec/blob/master/spec.md"
	defaultHealthCheckURL = "http://127.0.0.1:8080/ping"
)

type Config struct {
	Name       string `yaml:"name"`
	Command    string `yaml:"command"`
	PsCommand  string `yaml:"ps_command"`
	GlobalOpts []Opt  `yaml:"global_opts,omitempty"`
	Up         Verb   `yaml:"up,omitempty"`
	Down       Verb   `yaml:"down,omitempty"`
}

type Verb struct {
	Name string `yaml:"name"`
	Opts []Opt  `yaml:"opts,omitempty"`
}

type Opt struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}

type TestHelper struct {
	*testing.T
	testDir      string
	skipCommands []string
	specRef      string
}

type Response struct {
	Response string `yaml:"response"`
}

func (h TestHelper) TestUpDown(fun func()) {
	assert.Assert(h, fun != nil, "Test function cannot be `nil`")
	for _, f := range h.listFiles(commandsDir) {
		h.Run(f, func(t *testing.T) {
			c, err := h.readConfig(filepath.Join(commandsDir, f))
			assert.NilError(t, err)
			for _, v := range h.skipCommands {
				if v == c.Name {
					t.SkipNow()
				}
			}
			h.executeUp(c)
			h.waitHTTPReady(defaultHealthCheckURL, 5*time.Second)
			fun()
			h.executeDown(c)
			h.checkCleanUp(c)
		})
	}
}

func (h TestHelper) waitHTTPReady(url string, timeout time.Duration) {
	limit := time.Now().Add(timeout)
	for limit.After(time.Now()) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
	}
}

func (h TestHelper) Check(expected, actual string) {
	assert.Check(h.T, expected == actual, h.assertSpecReferenceMessage(expected, actual))
}

func (h TestHelper) NilError(e error) {
	assert.Check(h.T, e == nil, h.assertSpecReferenceMessage("", fmt.Sprintf("%q", e)))
}

func (h TestHelper) assertSpecReferenceMessage(expected, actual string) string {
	return fmt.Sprintf("\n- expected: %q\n+ actual: %q\n%s", expected, actual, h.specReferenceMessage())
}

func (h TestHelper) specReferenceMessage() string {
	return "Please refer to: " + h.getSpecReference()
}

func (h TestHelper) getSpecReference() string {
	if h.specRef != "" {
		return baseSpecReference + "#" + h.specRef
	}
	return baseSpecReference
}

func (h TestHelper) readConfig(configPath string) (*Config, error) {
	b, err := ioutil.ReadFile(configPath)
	h.NilError(err)
	c := Config{}
	err = yaml.Unmarshal(b, &c)
	h.NilError(err)
	return &c, nil
}

func verbWithOptions(c *Config, v Verb) []string {
	var gOpts []string
	for _, o := range c.GlobalOpts {
		gOpts = append(gOpts, o.Name)
		if o.Value != "" {
			gOpts = append(gOpts, o.Value)
		}
	}
	vOpts := append(gOpts, v.Name)
	for _, o := range v.Opts {
		vOpts = append(vOpts, o.Name)
		if o.Value != "" {
			vOpts = append(vOpts, o.Value)
		}
	}
	return vOpts
}

func (h TestHelper) executeUp(c *Config) {
	upOpts := verbWithOptions(c, c.Up)
	h.execCmd(c, upOpts)
}

func (h TestHelper) executeDown(c *Config) {
	downOpts := verbWithOptions(c, c.Down)
	h.execCmd(c, downOpts)
}

func (h TestHelper) execCmd(c *Config, opts []string) {
	cmd := icmd.Command(c.Command, opts...)
	cmd.Dir = filepath.Join("tests", h.testDir)
	icmd.RunCmd(cmd).Assert(h.T, icmd.Success)
}

func (h TestHelper) listDirs(testDir string) []string {
	currDir, err := os.Getwd()
	h.NilError(err)
	files, err := ioutil.ReadDir(filepath.Join(currDir, testDir))
	h.NilError(err)
	var dirs []string
	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs
}

func (h TestHelper) listFiles(dir string) []string {
	currDir, err := os.Getwd()
	h.NilError(err)
	content, err := ioutil.ReadDir(filepath.Join(currDir, dir))
	h.NilError(err)
	var configFiles []string
	for _, f := range content {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".yml") {
			configFiles = append(configFiles, f.Name())
		}
	}
	return configFiles
}

func (h TestHelper) checkCleanUp(c *Config) {
	command := strings.Split(c.PsCommand, " ")
	cmd := icmd.Command(command[0], command[1:]...)
	ret := icmd.RunCmd(cmd).Assert(h.T, icmd.Success)
	out := strings.Trim(ret.Stdout(), "\n")
	nLines := len(strings.Split(out, "\n")) - 1
	assert.Check(
		h.T,
		0 == nLines,
		"Problem checking containers' state. "+
			"There shouldn't be any containers before or after a test.")
}

func (h TestHelper) getHTTPBody(address string) string {
	resp, err := http.Get(address)
	h.NilError(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	h.NilError(err)
	return string(body)
}

func jsonResponse(content string) string {
	return fmt.Sprintf("{\"response\":\"%s\"}\n", content)
}
