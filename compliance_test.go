package main

import (
	"fmt"
	"testing"
	"time"
)

const (
	pingEntrypoint = "http://localhost:8080/ping"
	pingTargetUrl  = "target:8080/ping"
	pingUrl        = pingEntrypoint + "?address=" + pingTargetUrl

	volumefileEntrypoint = "http://localhost:8080/volumefile"
	volumeUrl            = volumefileEntrypoint + "?filename="
)

func TestSimpleLifecycle(t *testing.T) {
	h := TestHelper{
		T:       t,
		testDir: "simple_lifecycle",
	}
	h.TestUpDown(func() {
		time.Sleep(time.Second)
	})
}

func TestSimpleNetwork(t *testing.T) {
	h := TestHelper{
		T:       t,
		testDir: "simple_network",
		specRef: "Networks-top-level-element",
	}
	h.TestUpDown(func() {
		actual := h.getHttpBody(pingUrl)
		expected := jsonResponse("PONG FROM TARGET")
		h.Check(expected, actual)
	})
}

func TestSimpleNetworkFail(t *testing.T) {
	h := TestHelper{
		T:       t,
		testDir: "simple_network",
		specRef: "Networks-top-level-element",
	}
	h.TestUpDown(func() {
		actual := h.getHttpBody("http://localhost:8080/ping?address=notatarget:8080/ping")
		expected := jsonResponse("Could not reach address: notatarget:8080/ping")
		h.Check(expected, actual)
	})
}

func TestDifferentNetworks(t *testing.T) {
	h := TestHelper{
		T:       t,
		testDir: "different_networks",
		specRef: "Networks-top-level-element",
	}
	h.TestUpDown(func() {
		actual := h.getHttpBody(pingUrl)
		expected := jsonResponse("Could not reach address: target:8080/ping")
		h.Check(expected, actual)
	})
}

func TestVolumeFile(t *testing.T) {
	h := TestHelper{
		T:       t,
		testDir: "simple_volume",
		specRef: "volumes-top-level-element",
	}
	h.TestUpDown(func() {
		actual := h.getHttpBody(volumeUrl + "test_volume.txt")
		expected := jsonResponse("MYVOLUME")
		h.Check(expected, actual)

	})
}

func TestSecretFile(t *testing.T) {
	h := TestHelper{
		T:       t,
		testDir: "simple_secretfile",
		specRef: "secrets-top-level-element",
	}
	h.TestUpDown(func() {
		actual := h.getHttpBody(volumeUrl + "test_secret.txt")
		expected := jsonResponse("MYSECRET")
		h.Check(expected, actual)
	})
}

func TestConfigFile(t *testing.T) {
	h := TestHelper{
		T:            t,
		testDir:      "simple_configfile",
		skipCommands: []string{"docker-composeV1"},
		specRef:      "configs-top-level-element",
	}
	h.TestUpDown(func() {
		actual := h.getHttpBody(volumeUrl + "test_config.txt")
		expected := jsonResponse("MYCONFIG")
		h.Check(expected, actual)
	})
}

func jsonResponse(content string) string {
	return fmt.Sprintf("{\"response\":\"%s\"}\n", content)
}
