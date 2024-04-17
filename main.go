package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type App struct {
	Listen  string      `json:"listen"`
	Backend BackendList `json:"backend"`
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := a.Backend.Do(r.Context(), r); err != nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (a *App) load() error {
	configPath := os.Args[1]
	if configPath == "" {
		configPath = "/etc/auth-proxy.conf"
	}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	config := &App{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return err
	}

	if config.Listen == "" {
		config.Listen = ":1064"
	}

	return nil
}

func start() error {
	app := &App{}

	if err := app.load(); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	return http.ListenAndServe(app.Listen, app)
}

func main() {
	if err := start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
