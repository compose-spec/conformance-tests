package main

import (
	"path/filepath"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestSimpleUpDown(t *testing.T) {
	for _, f := range listFiles(t, "commands") {
		t.Run(f, func(t *testing.T) {
			c, err := readConfig(t, filepath.Join("commands", f))
			assert.NilError(t, err)
			for _, configName := range listDirs(t, "tests") {
				assert.NilError(t, err)
				executeUp(t, c, configName)
				time.Sleep(2 * time.Second) // FIXME Get a better way to do that
				executeDown(t, c, configName)
			}
		})
	}
}
