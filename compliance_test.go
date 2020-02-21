package main

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

const (
	pingEntrypoint = "http://localhost:8080/ping"
	pingTargetUrl  = "target:8080/ping"
	pingUrl        = pingEntrypoint + "?address=" + pingTargetUrl

	volumefileEntrypoint = "http://localhost:8080/volumefile"
	volumeUrl            = volumefileEntrypoint + "?filename="
)

func TestSimpleLifecycle(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_lifecycle"}
	h.TestUpDown(func(t *testing.T) {
		time.Sleep(time.Second)
	})
}

func TestSimpleNetwork(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_network"}
	h.TestUpDown(func(t *testing.T) {
		actual := h.getHttpBody(pingUrl)
		assert.Assert(t, actual == "{\"response\":\"PONG FROM TARGET\"}\n")
	})
}

func TestSimpleNetworkFail(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_network"}
	h.TestUpDown(func(t *testing.T) {
		actual := h.getHttpBody("http://localhost:8080/ping?address=notatarget:8080/ping")
		assert.Assert(t, actual == "{\"response\":\"Could not reach address: notatarget:8080/ping\"}\n")
	})
}

func TestDifferentNetworks(t *testing.T) {
	h := TestHelper{T: t, testDir: "different_networks"}
	h.TestUpDown(func(t *testing.T) {
		actual := h.getHttpBody(pingUrl)
		assert.Assert(t, actual == "{\"response\":\"Could not reach address: target:8080/ping\"}\n")
	})
}

func TestVolumeFile(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_volume"}
	h.TestUpDown(func(t *testing.T) {
		actual := h.getHttpBody(volumeUrl + "test_volume.txt")
		assert.Assert(t, actual == "{\"response\":\"MYVOLUME\"}\n")
	})
}

func TestSecretFile(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_secretfile"}
	h.TestUpDown(func(t *testing.T) {
		actual := h.getHttpBody(volumeUrl + "test_secret.txt")
		assert.Assert(t, actual == "{\"response\":\"MYSECRET\"}\n")
	})
}

func TestConfigFile(t *testing.T) {
	h := TestHelper{
		T:            t,
		testDir:      "simple_configfile",
		skipCommands: []string{"docker-composeV1"},
	}
	h.TestUpDown(func(t *testing.T) {
		actual := h.getHttpBody(volumeUrl + "test_config.txt")
		assert.Assert(t, actual == "{\"response\":\"MYCONFIG\"}\n")
	})
}
