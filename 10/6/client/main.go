package main

import (
	"context"
	"log"
	"time"

	userpb "3/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

    user1 := createUser(client, &userpb.User{Username: "nikita2004", Age: 20})
    user2 := createUser(client, &userpb.User{Username: "nikita2005", Age: 19})
    getUsers(client)
    deleteUser(client, user1.Id)
    getUsers(client)
    updatedUser2 := user2
    updatedUser2.Username = "newusername"
    updateUser(client, updatedUser2)
    getUsers(client)
}

func createUser(client userpb.UserServiceClient, user *userpb.User) *userpb.User {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &userpb.CreateUserRequest{User: user}
	res, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	log.Printf("user created: %v\n", res.GetUser())
    return res.GetUser()
}

func getUsers(client userpb.UserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &userpb.GetUsersRequest{}
	res, err := client.GetUsers(ctx, req)
	if err != nil {
		log.Fatalf("failed to get users: %v", err)
	}

	log.Printf("got users: %v\n", res.GetUsers())
}

func deleteUser(client userpb.UserServiceClient, id int32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

    req := &userpb.DeleteUserRequest{Id: id}
	_, err := client.DeleteUser(ctx, req)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	log.Printf("deleted user: %v\n", id)
}

func updateUser(client userpb.UserServiceClient, user *userpb.User) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

    req := &userpb.UpdateUserRequest{User: user}
	_, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("updated user: %v\n", user)
}
