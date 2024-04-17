package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

type App struct {
	Listen  string      `json:"listen"`
	Backend BackendList `json:"backend"`
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.Backend.Check(r) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func (a *App) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(a); err != nil {
		return err
	}

	if a.Listen == "" {
		a.Listen = ":1064"
	}

	return nil
}

func start(path string) error {
	app := &App{}

	if err := app.load(path); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	server := http.Server{
		Addr:    app.Listen,
		Handler: app,
	}

	go func() {
		<-signalCh
		server.Close()
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func main() {
	configPath := "/etc/auth-proxy.conf"

	if len(os.Args) > 1 {
		configPath = os.Args[1]

		switch configPath {
		case "help", "-h", "--help":
			fmt.Fprintf(os.Stdout, "Usage: %s [config]\n", os.Args[0])
			os.Exit(0)

		case "version", "-v", "--version":
			fmt.Fprintf(os.Stdout, "Auth Proxy v0\n")
			os.Exit(0)
		}
	}

	if err := start(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
