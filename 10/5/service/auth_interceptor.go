package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
    jwtManager *JWTManager
}

func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
    return &AuthInterceptor{jwtManager}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
            log.Println("--> unary interceptor: ", info.FullMethod)

            err := i.authorize(ctx, info.FullMethod)
            if err != nil {
                return nil, err
            }

            return handler(ctx, req)
    }
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) error {
    switch method {
    case "/jaennil.effective_mobile_internship_ten_three.UserService/CreateUser":
        return nil
    case "/jaennil.effective_mobile_internship_ten_three.UserService/AuthUser":
        return nil
    }

    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return status.Errorf(codes.Unauthenticated, "metadata is not provided")
    }

    values := md["authorization"]
    if len(values) == 0 {
        return status.Errorf(codes.Unauthenticated, "authentication token is not provided")
    }

    accessToken := values[0]
    _, err := i.jwtManager.Verify(accessToken)
    if err != nil {
        return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
    }

    return nil
}
