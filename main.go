package main

import (
	"backend/config"
	"backend/routes"
	"backend/verification"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Validator = &verification.TaskValidator{Validator: validator.New()}

	// Load configuration
	config.LoadEnv()

	// Setup routes
	routes.SetupRoutes(e)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Canal para recibir señales del sistema operativo
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Esperar a recibir una señal
	<-quit
	fmt.Println("Recibida señal de interrupción, iniciando apagado ordenado...")

	// Crear un contexto con tiempo de espera
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Apagar el servidor de forma ordenada
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Servidor apagado correctamente.")
}
