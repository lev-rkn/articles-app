package transport

import (
	authv1 "auth-service/api/gen/proto"
	"auth-service/internal/lib/types"
	"auth-service/internal/service"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Регистрируем нашу реализацию обработчиков контракта proto в grpc сервер
func AuthRegister(gRPCServer *grpc.Server, auth service.AuthServiceInterface) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{authService: auth})
}

type serverAPI struct {
	// позволяет обеспечить обратную совместимость при изменении auth.proto
	// файла и позволит избежать ошибки, если забудем реализовать метод
	authv1.UnimplementedAuthServer
	authService service.AuthServiceInterface
}

var _ authv1.AuthServer = (*serverAPI)(nil)

func (s *serverAPI) Register(
	ctx context.Context,
	in *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrEmailRequired.Error())
	}
	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrPassRequired.Error())
	}

	uid, err := s.authService.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, types.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "failed to register user: "+err.Error())
	}

	return &authv1.RegisterResponse{UserId: uid}, nil
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	// TODO: автоматизировать валидацию полей в protobuf
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrEmailRequired.Error())
	}
	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrPassRequired.Error())
	}
	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, types.ErrAppIdRequired.Error())
	}
	if in.GetFingerprint() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrFingerprintRequired.Error())
	}

	tokenPair, err := s.authService.Login(
		ctx,
		in.GetEmail(),
		in.GetPassword(),
		in.GetAppId(),
		in.GetFingerprint(),
	)
	if err != nil {
		// TODO: нужно ошибки оборачивать в join, чтобы можно было отправить нужный код
		if errors.Is(err, types.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, types.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *serverAPI) Refresh(
	ctx context.Context,
	in *authv1.RefreshTokenRequest,
) (*authv1.RefreshTokenResponse, error) {
	if in.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrRefreshRequired.Error())
	}
	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, types.ErrAppIdRequired.Error())
	}
	if in.GetFingerprint() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrFingerprintRequired.Error())
	}

	tokenPair, err := s.authService.RefreshToken(
		ctx,
		in.GetRefreshToken(),
		in.GetFingerprint(),
	)
	if err != nil {
		if errors.Is(err, types.ErrUnidentifiedDevice) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, types.ErrRefreshTokenNotValid) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
