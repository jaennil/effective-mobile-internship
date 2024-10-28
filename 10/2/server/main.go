package main

import (
	userpb "2/user"
	"context"
	"log"
	"net"
    "fmt"

	"google.golang.org/grpc"
)

type UserStorage struct {
	users  map[int32]*userpb.User
	nextID int32
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users:  make(map[int32]*userpb.User),
		nextID: 1,
	}
}

type userServiceServer struct {
	storage *UserStorage
	userpb.UnimplementedUserServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	storage := NewUserStorage()

	userpb.RegisterUserServiceServer(grpcServer, &userServiceServer{storage: storage})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user := req.GetUser()
    user.Id = s.storage.nextID
    s.storage.nextID++
	s.storage.users[user.Id] = user

	return &userpb.CreateUserResponse{User: user}, nil
}

func (s *userServiceServer) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
    var users []*userpb.User
    for _, user := range s.storage.users {
        users = append(users, user)
    }
    return &userpb.GetUsersResponse{Users: users}, nil
}

func (s *userServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
    id := req.GetId()
    if _, exists := s.storage.users[id]; !exists {
        return nil, fmt.Errorf("user with id `%v` not found", id)
    }

    delete(s.storage.users, id)

    return &userpb.DeleteUserResponse{}, nil
}

func (s *userServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
    user := req.GetUser()
    
    if _, exists := s.storage.users[user.Id]; !exists {
        return nil, fmt.Errorf("user with id `%v` not found", user.Id)
    }

    s.storage.users[user.Id] = user

    return &userpb.UpdateUserResponse{User: user}, nil
}
