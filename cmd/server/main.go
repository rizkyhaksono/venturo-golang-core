package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	_ "venturo-core/docs"
	"venturo-core/internal/server"
	"venturo-core/pkg/logger"
)

//	@title			Venturo Golang Core API
//	@version		1.0
//	@description	This is the API documentation for the Venturo Golang Core project.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Anis Fajar Fakhruddin
//	@contact.url	https://discord.com/users/858389159555497994
//	@contact.email	sina4science@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and a JWT.
func main() {
	// Get the app and the shared WaitGroup from our server setup
	app, wg := server.NewServer()
	logger.InitLogger()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server is starting", "port", 3000)
		if err := app.Listen(":3000"); err != nil {
			slog.Error("Server failed to start", "error", err)
		}
	}()

	// Block until a signal is received
	<-quit

	// Trigger the graceful shutdown, passing the shared WaitGroup
	server.GracefulShutdown(app, wg)
}
