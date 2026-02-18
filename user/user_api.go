package user

import (
	"context"
	"log"

	"rjfield.com/backend/db"
	"rjfield.com/backend/generated/pb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	log.Printf("ListUsers() - received: %v", req)
	users, err := db.ListUsers(req.PageSize, req.PageToken)
	if err != nil {
		return nil, err
	}
	log.Printf("ListUsers() - returning: %v", users)
	return &pb.ListUsersReply{
		Users: users,
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	log.Printf("GetUser() - received: %v", req.GetUserId())
	u, err := db.ReadUser(req.UserId)
	if err != nil {
		return nil, err
	}
	log.Printf("GetUser() - returning: %v", u)
	return &pb.GetUserReply{
		User: u,
	}, nil
}

func (s *UserServer) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserReply, error) {
	log.Printf("AuthenticateUser() - received: %s / %s", req.GetUsername(), req.GetPassword())
	u, err := db.AuthenticateUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	log.Printf("AuthenticateUser() - returning: %s", u)
	return &pb.AuthenticateUserReply{
		Token:  "sample_token",
		UserId: u,
	}, nil
}
