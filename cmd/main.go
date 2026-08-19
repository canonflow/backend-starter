package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/canonflow/backend-starter/internal/app"
	_ "github.com/canonflow/backend-starter/internal/app"
	"github.com/canonflow/backend-starter/internal/config"
	_ "github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/core"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func main() {
	fiberApp := app.NewFiber()
	db := app.NewDatabase(config.Get[string](config.DBDriver), core.GetLogger())

	appHost := config.Get[string](config.AppHost)
	appPort := config.Get[string](config.AppPort)
	app.Bootstrap(&app.BootstrapConfig{
		DB:  db,
		App: fiberApp,
	})

	go func() {
		if err := fiberApp.Listen(appHost+":"+appPort, fiber.ListenConfig{
			EnablePrefork: true,
		}); err != nil {
			core.GetLogger().Info("Server Error", zap.String("error", err.Error()))
		}
	}()

	// Block until an interrupt/terminate signal is received
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	core.GetLogger().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fiberApp.ShutdownWithContext(ctx); err != nil {
		core.GetLogger().Fatal("forced shutdown: %v", zap.String("error", err.Error()))
	}

	core.GetLogger().Info("Server exited gracefully")
}
