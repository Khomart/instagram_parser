package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/khomart/instagram_parser/internal/downloader"
	"github.com/khomart/instagram_parser/internal/parser"
	"github.com/khomart/instagram_parser/internal/server"
)

func main() {
	slog.Info("Starting application")
	d := downloader.NewDownloader()
	p := parser.NewParser()
	rh := server.NewRequestHandler(d, p)
	r, err := server.NewRouter(rh)
	if err != nil {
		slog.Error(fmt.Sprintf("Error setup router: %s", err.Error()))
		os.Exit(1)
	}
	addr := fmt.Sprintf("%s:%s", "localhost", "8080")
	slog.Info("Starting the HTTP server", "listen_address", addr)
	err = r.Serve(addr)
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting the server: %s", err.Error()))
		os.Exit(1)

	}
}
