package main

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

const (
	defaultPingEntrypoint = "http://localhost:8080/ping"
	targetUrl             = "target:8080/ping"
	defaultTargetUrl      = defaultPingEntrypoint + "?address=" + targetUrl
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
		actual := getHttpBody(t, defaultTargetUrl)
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
		actual := getHttpBody(t, defaultTargetUrl)
		assert.Assert(t, actual == "{\"response\":\"Could not reach address: target:8080/ping\"}\n")
	})
}
