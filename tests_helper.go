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

func readConfig(t *testing.T, configPath string) (*Config, error) {
	b, err := ioutil.ReadFile(configPath)
	assert.NilError(t, err)
	c := Config{}
	err = yaml.Unmarshal(b, &c)
	assert.NilError(t, err)
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

func executeUp(t *testing.T, c *Config, configName string) {
	upOpts := verbWithOptions(c, c.Up)
	execCmd(t, c.Command, configName, upOpts)
}

func executeDown(t *testing.T, c *Config, configName string) {
	downOpts := verbWithOptions(c, c.Down)
	execCmd(t, c.Command, configName, downOpts)
}

func execCmd(t *testing.T, command string, configName string, opts []string) {
	cmd := icmd.Command(command, opts...)
	cmd.Dir = filepath.Join("tests", configName)
	icmd.RunCmd(cmd).Assert(t, icmd.Success)
}

func listDirs(t *testing.T, testDir string) []string {
	currDir, err := os.Getwd()
	assert.NilError(t, err)
	files, err := ioutil.ReadDir(filepath.Join(currDir, testDir))
	assert.NilError(t, err)
	var dirs []string
	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs
}

func listFiles(t *testing.T, dir string) []string {
	currDir, err := os.Getwd()
	assert.NilError(t, err)
	content, err := ioutil.ReadDir(filepath.Join(currDir, dir))
	assert.NilError(t, err)
	var configFiles []string
	for _, f := range content {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".yml") {
			configFiles = append(configFiles, f.Name())
		}
	}
	return configFiles
}

func checkDown(t *testing.T) {
	cli, err := client.NewEnvClient()
	assert.NilError(t, err)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	assert.NilError(t, err)
	assert.Assert(t, len(containers) == 0)
}

func getHttpBody(t *testing.T, address string) string {
	resp, err := http.Get(address)
	assert.NilError(t, err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NilError(t, err)
	return string(body)
}

type TestHelper struct {
	*testing.T
	testDir string
}

func (h TestHelper) testUpDown(fun func(t *testing.T)) {
	assert.Assert(h, fun != nil, "Test function cannot be `nil`")
	for _, f := range listFiles(h.T, commandsDir) {
		h.Run(f, func(t *testing.T) {
			c, err := readConfig(t, filepath.Join(commandsDir, f))
			assert.NilError(t, err)
			executeUp(t, c, h.testDir)
			fun(t)
			executeDown(t, c, h.testDir)
			checkDown(t)
		})
	}
}
