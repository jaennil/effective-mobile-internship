package main

import (
	"context"
	"log"
	"net"
	"time"

	"3/pb"
	"3/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


type userServiceServer struct {
	storage *service.UserStorage
    jwtManager *service.JWTManager
    pb.UnimplementedUserServiceServer
}

const (
    jwtSecret = "jwtSecret"
    tokenDuration = 24 * time.Hour
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

    jwtManager := service.NewJWTManager(jwtSecret, tokenDuration)

    interceptor := service.NewAuthInterceptor(jwtManager)

	grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(interceptor.Unary()),
    )

	storage := service.NewUserStorage()


    pb.RegisterUserServiceServer(grpcServer, &userServiceServer{storage: storage, jwtManager: jwtManager})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
        return nil, err
    }

    if err = s.storage.Save(user); err != nil {
        return nil, err
    }

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
