package transport

import (
	"auth-service/internal/lib/utils"
	"auth-service/internal/service"
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Создает экземпляр grpc сервера и регистрирует обработчики контрактов
func NewServer(ctx context.Context, service *service.Service) *grpc.Server {
	// Если мы встречаем панику во время работы нашего сервера, то мы ее обрабатываем.
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			utils.ErrorLog("Recovered from panic", err)
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	logger := logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		slog.Log(ctx, slog.Level(lvl), msg, fields...)
	})
	// создаем экземпляр grpc сервера с перехватчиками
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(logger),
	))

	AuthRegister(gRPCServer, service.Auth)

	return gRPCServer
}
