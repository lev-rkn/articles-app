package main

import (
	"ads-service/internal/app"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := app.NewServer()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Listen server", "err", err.Error())
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// ждем сигнала о завершении программы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	// создем контекст с таймаутом в 5 секунд и передаем его в функцию выключения сервера.
	// через 5 секунд наш сервера закончит свою работу
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", "err", err.Error())
	}

	<-ctx.Done()
	slog.Info("timeout of 5 seconds, server exiting")
}
