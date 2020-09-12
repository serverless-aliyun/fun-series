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

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	svc := &service{domain: getenv("DOMAIN", "www.rrys2020.com")}
	ctrl := &controller{svc}

	rg := r.Group("/series")
	rg.GET("", ctrl.search)
	rg.GET("/:seriesId", ctrl.detail)
	rg.GET("/:seriesId/episodes", ctrl.episodes)

	start(&http.Server{
		Addr:    fmt.Sprintf(":%s", getenv("FC_SERVER_PORT", "9000")),
		Handler: r,
	})
}

func start(srv *http.Server) {
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err)
		}
	}()

	log.Printf("Start Server @ %s", srv.Addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown:%s", err)
	}
	<-ctx.Done()
	log.Print("Server exiting")
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
