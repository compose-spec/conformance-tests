package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

const defaultPort = 8080

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

func main() {
	port := defaultPort
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
	e.Logger.Fatal(e.StartServer(s))
}
