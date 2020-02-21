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
	h.testUpDown(func(t *testing.T) {
		time.Sleep(time.Second) // FIXME Test when the up is complete here
	})
}

func TestSimpleNetwork(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_network"}
	h.testUpDown(func(t *testing.T) {
		actual := getHttpBody(t, pingUrl)
		assert.Assert(t, actual == "{\"response\":\"PONG FROM TARGET\"}\n")
	})
}

func TestSimpleNetworkFail(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_network"}
	h.testUpDown(func(t *testing.T) {
		actual := getHttpBody(t, "http://localhost:8080/ping?address=notatarget:8080/ping")
		assert.Assert(t, actual == "{\"response\":\"Could not reach address: notatarget:8080/ping\"}\n")
	})
}

func TestDifferentNetworks(t *testing.T) {
	h := TestHelper{T: t, testDir: "different_networks"}
	h.testUpDown(func(t *testing.T) {
		actual := getHttpBody(t, pingUrl)
		assert.Assert(t, actual == "{\"response\":\"Could not reach address: target:8080/ping\"}\n")
	})
}

func TestVolumeFile(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_volume"}
	h.testUpDown(func(t *testing.T) {
		actual := getHttpBody(t, volumeUrl+"test_volume.txt")
		assert.Assert(t, actual == "{\"response\":\"MYVOLUME\"}\n")
	})
}

func TestSecretFile(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_secretfile"}
	h.testUpDown(func(t *testing.T) {
		actual := getHttpBody(t, volumeUrl+"test_secret.txt")
		assert.Assert(t, actual == "{\"response\":\"MYSECRET\"}\n")
	})
}

// Ignored because docker-compose does not support that for now
func _TestConfigFile(t *testing.T) {
	h := TestHelper{T: t, testDir: "simple_configfile"}
	h.testUpDown(func(t *testing.T) {
		actual := getHttpBody(t, volumeUrl+"test_config.txt")
		assert.Assert(t, actual == "{\"response\":\"MYCONFIG\"}\n")
	})
}
