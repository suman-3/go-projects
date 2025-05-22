package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"githu.com/suman-3/go-projects/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	// init storage / database

	// init router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to my First Go Server"))
	})

	// init server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server started on port %s\n", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Server failed to start: %s", err)
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	server.Shutdown(ctx)

	slog.Info("Server stopped")

}
