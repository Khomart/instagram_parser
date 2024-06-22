package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/khomart/instagram_recipe_parser/internal/config"
	"github.com/khomart/instagram_recipe_parser/internal/downloader"
	"github.com/khomart/instagram_recipe_parser/internal/parser"
	"github.com/khomart/instagram_recipe_parser/internal/server"
)

func main() {
	slog.Info("Starting application")
	c, err := config.NewConfig()
	if err != nil {
		slog.Error("Error initializing config", "error", err)
		os.Exit(1)
	}
	d := downloader.NewDownloader()
	p := parser.NewParser(c)
	rh := server.NewRequestHandler(d, p)
	r, err := server.NewRouter(rh)
	if err != nil {
		slog.Error(fmt.Sprintf("Error setup router: %s", err.Error()))
		os.Exit(1)
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	slog.Info("Starting the HTTP server", "address", addr)
	err = r.Serve(addr)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting the server: %s", err.Error()))
		os.Exit(1)

	}
}
