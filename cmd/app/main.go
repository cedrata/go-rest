package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cedrata/go-rest/pkg/middleware"
)

type testHandler struct{}

func (testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello"))
	time.Sleep(7 * time.Second)
}

func main() {
	log.Println("server starting")

	mux := http.NewServeMux()
	mux.Handle("/test", &testHandler{})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server started")

	<-done
	log.Println("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}

	middleware.LogMiddleware()
}
