package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/icmd"
)

const (
	commandsDir = "commands"
)

type Config struct {
	Name       string `yaml:"name"`
	Command    string `yaml:"command"`
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
}

func (h TestHelper) TestUpDown(fun func(t *testing.T)) {
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
			fun(t)
			h.executeDown(c)
			h.checkDown()
		})
	}
}

func (h TestHelper) readConfig(configPath string) (*Config, error) {
	b, err := ioutil.ReadFile(configPath)
	assert.NilError(h.T, err)
	c := Config{}
	err = yaml.Unmarshal(b, &c)
	assert.NilError(h.T, err)
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
	assert.NilError(h.T, err)
	files, err := ioutil.ReadDir(filepath.Join(currDir, testDir))
	assert.NilError(h.T, err)
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
	assert.NilError(h.T, err)
	content, err := ioutil.ReadDir(filepath.Join(currDir, dir))
	assert.NilError(h.T, err)
	var configFiles []string
	for _, f := range content {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".yml") {
			configFiles = append(configFiles, f.Name())
		}
	}
	return configFiles
}

func (h TestHelper) checkDown() {
	cli, err := client.NewEnvClient()
	assert.NilError(h.T, err)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	assert.NilError(h.T, err)
	assert.Assert(h.T, len(containers) == 0)
}

func (h TestHelper) getHttpBody(address string) string {
	resp, err := http.Get(address)
	assert.NilError(h.T, err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(h.T, err)
	return string(body)
}
