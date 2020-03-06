package main

import (
	"net/http"
	"os"
	"time"
)

const targetHost = "server"

func main() {
	time.Sleep(time.Second)
	value := os.Getenv("HOSTNAME")
	resp, err := http.Get("http://" + targetHost + ":8080/scalechecker?value=" + value)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}
