package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"6/pb"
	"6/service"

	"github.com/KaranJagtiani/go-logstash"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


type userServiceServer struct {
    logstash *logstash_logger.Logstash
	storage *service.UserStorage
    jwtManager *service.JWTManager
    pb.UnimplementedUserServiceServer
}

func NewUserServiceServer(logstash *logstash_logger.Logstash, storage *service.UserStorage, jwtManager *service.JWTManager) *userServiceServer {
    return &userServiceServer{logstash: logstash, storage: storage, jwtManager: jwtManager}
}

const (
    jwtSecret = "jwtSecret"
    tokenDuration = 24 * time.Hour
)

func logstashError(message string, err error) map[string]interface{} {
    return map[string]interface{}{"message": message, "error": err}
}

func logstashMessage(message string) map[string]interface{} {
    return map[string]interface{}{"message": message}
}

func main() {
    logstash := logstash_logger.Init("localhost", 50000, "tcp", 5)

    jwtManager := service.NewJWTManager(jwtSecret, tokenDuration)
    authInterceptor := service.NewAuthInterceptor(jwtManager)

    prometheusMetrics := grpcprom.NewServerMetrics()
    reg := prometheus.NewRegistry()
    reg.MustRegister(prometheusMetrics)

	grpcServer := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            authInterceptor.Unary(),
            prometheusMetrics.UnaryServerInterceptor(),
        ),
    )

    prometheusMetrics.InitializeMetrics(grpcServer)

	storage := service.NewUserStorage()
    pb.RegisterUserServiceServer(grpcServer, NewUserServiceServer(logstash, storage, jwtManager))

    go func() {
        http.Handle("/metrics", promhttp.HandlerFor(
            reg,
            promhttp.HandlerOpts {
            },
        ))
        log.Fatal(http.ListenAndServe(":9091", nil))
    }()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
        logstash.Error(logstashError("failed to listen on tcp", err))
        os.Exit(1)
	}

	if err := grpcServer.Serve(listener); err != nil {
        logstash.Error(logstashError("failed to serve grpc server", err))
        os.Exit(1)
	}
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

    existingUser, err := s.storage.Find(req.GetUsername())
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot find user")
    }

    if existingUser != nil {
        return nil, status.Errorf(codes.AlreadyExists, "user already exists")
    }

    user, err := service.NewUser(req.GetUsername(), req.GetPassword())
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    if err = s.storage.Save(user); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    s.logstash.Debug(logstashMessage(fmt.Sprintf("users: %v", s.storage.Users)))

    log.Printf("users: %v", s.storage.Users)

	return &pb.CreateUserResponse{}, nil
}

func (s *userServiceServer) AuthUser(ctx context.Context, req *pb.AuthUserRequest) (*pb.AuthUserResponse, error) {
    user, err := s.storage.Find(req.GetUsername())
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
    }

    if user == nil || !user.IsPasswordCorrect(req.GetPassword()) {
        return nil, status.Errorf(codes.NotFound, "incorrect username/password")
    }

    token, err := s.jwtManager.GenerateJWT(user)
    if err != nil {
        log.Printf("cannot generate access token: %v", err)
        return nil, status.Errorf(codes.Internal, "cannot generate access token")
    }

    log.Printf("users: %+v", s.storage.Users)

    return &pb.AuthUserResponse{AccessToken: token}, nil
}

func (s *userServiceServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
    var usernames[]string
    for _, user := range s.storage.Users {
        usernames = append(usernames, user.Username)
    }

    return &pb.GetUsersResponse{Usernames: usernames}, nil
}

func (s *userServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    md, _ := metadata.FromIncomingContext(ctx)
    accessToken := md["authorization"][0]
    claims, err := s.jwtManager.Verify(accessToken)
    username := claims.Username

    user, err := s.storage.Find(username)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot find user")
    }

    if user == nil {
        return nil, status.Errorf(codes.NotFound, "user not found")
    }

    delete(s.storage.Users, user.Id)

    log.Printf("users: %+v", s.storage.Users)

    return &pb.DeleteUserResponse{}, nil
}

func (s *userServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    md, _ := metadata.FromIncomingContext(ctx)
    accessToken := md["authorization"][0]
    claims, err := s.jwtManager.Verify(accessToken)
    username := claims.Username
    user, err := s.storage.Find(username)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot find user")
    }

    if user == nil {
        return nil, status.Errorf(codes.NotFound, "cannot find user")
    }

    user.Username = req.GetUsername()
    if err = s.storage.Update(user.Clone()); err != nil {
        return nil, status.Errorf(codes.Internal, "cannot save user")
    }

    log.Printf("users: %+v", s.storage.Users)

    return &pb.UpdateUserResponse{}, nil
}
