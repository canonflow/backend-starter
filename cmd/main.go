package main

import (
	"context"
	"fmt"

	"github.com/canonflow/backend-starter/internal/app"
	_ "github.com/canonflow/backend-starter/internal/app"
	"github.com/canonflow/backend-starter/internal/config"
	_ "github.com/canonflow/backend-starter/internal/config"
)

func main() {
	environment := config.Get[string](config.AppEnv)

	fmt.Println(environment)

	ctx := context.Background()
	ctx = app.NewContext(ctx)

	logger := app.LoggerFromContext(ctx)

	logger.Info("test")
	logger.Info("test2")
}
