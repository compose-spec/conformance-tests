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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

const (
	defaultHttpPort = 8080
	defaultUdpPort  = 10001
)

func getMapResponse(response string) map[string]string {
	return map[string]string{
		"response": response,
	}
}

type Ping struct {
	Address string `json:"address"`
}

func pingHandler(c echo.Context) error {
	p := new(Ping)
	if err := c.Bind(p); err != nil {
		c.Error(err)
		return err
	}
	if p.Address == "" {
		return c.JSON(
			http.StatusOK,
			getMapResponse("PONG FROM TARGET"),
		)
	}
	resp, err := http.Get(fmt.Sprintf("http://%s", p.Address))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			getMapResponse(fmt.Sprintf("Could not reach address: %s", p.Address)),
		)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			getMapResponse(fmt.Sprintf("Could not body from response: %s", err)),
		)
	}
	return c.String(http.StatusOK, string(body))
}

type GetFile struct {
	Filename string `json:"filename"`
}

func fileHandler(c echo.Context) error {
	g := new(GetFile)
	if err := c.Bind(g); err != nil {
		c.Error(err)
		return err
	}
	b, err := ioutil.ReadFile(filepath.Join("/volumes", g.Filename))
	if err != nil {
		c.Error(err)
		return err
	}
	return c.JSON(
		http.StatusOK,
		getMapResponse(string(b)),
	)
}

var udpValue string

func udpServer() {
	fmt.Println("Running UDP server...")
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", defaultUdpPort))
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	checkError(err, true)
	defer ServerConn.Close()

	var objmap map[string]string
	buf := make([]byte, 1024)
	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		checkError(err, true)
		err = json.Unmarshal(buf[:n], &objmap)
		checkError(err, true)
		udpValue = objmap["request"]
	}
}

func udpHandler(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		getMapResponse(udpValue),
	)
}

var scaleValues []string

type ScaleValue struct {
	Value string `json:"value"`
}

func scaleCheckHandler(c echo.Context) error {
	s := new(ScaleValue)
	if err := c.Bind(s); err != nil {
		c.Error(err)
		return err
	}
	if s.Value != "" && !containsString(scaleValues, s.Value) {
		scaleValues = append(scaleValues, s.Value)
	}
	return c.JSON(
		http.StatusOK,
		getMapResponse(fmt.Sprintf("%d", len(scaleValues))),
	)
}

func containsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func checkError(err error, exitOnError bool) {
	if err != nil {
		fmt.Println("Error: ", err)
		if exitOnError {
			os.Exit(0)
		}
	}
}

func main() {
	go udpServer()

	port := defaultHttpPort
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort != "" {
		port, _ = strconv.Atoi(httpPort)
	}
	e := echo.New()
	e.HideBanner = true
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	e.GET("/ping", pingHandler)
	e.GET("/volumefile", fileHandler)
	e.GET("/udp", udpHandler)
	e.GET("/scalechecker", scaleCheckHandler)
	e.Logger.Fatal(e.StartServer(s))
}
