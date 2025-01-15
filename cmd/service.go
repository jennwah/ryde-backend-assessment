package main

import (
	"context"
	"errors"
	"fmt"

	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jennwah/ryde-backend-engineer/internal/api"
	"github.com/jennwah/ryde-backend-engineer/internal/config"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed loading application config: %w", err))
	}

	postgres, err := postgresql.New(cfg.Postgres)
	if err != nil {
		panic(fmt.Errorf("failed initializing connection with database: %w", err))
	}

	router := gin.Default()
	router.Use(gin.Recovery())

	apiController := api.New(postgres.DB, logger, cfg)
	apiController.RegisterHandlers(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = postgres.Close(); err != nil {
		logger.Error("database connection clean up err", slog.Any("error", err))
	}

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
