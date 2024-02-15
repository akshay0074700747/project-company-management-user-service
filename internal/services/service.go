package services

import (
	"context"
	"errors"

	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/authpb"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"github.com/akshay0074700747/projectandCompany_management_user-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_user-service/helpers"
	"github.com/akshay0074700747/projectandCompany_management_user-service/internal/usecases"
	"github.com/golang/protobuf/ptypes/empty"
)

type UserServiceServer struct {
	Usecase  usecases.UserUsecaseInterfaces
	AuthConn authpb.AuthServiceClient
	userpb.UnimplementedUserServiceServer
}

func NewUserServiceServer(usecase usecases.UserUsecaseInterfaces, authAddr string) *UserServiceServer {
	authConn, _ := helpers.DialGrpc(authAddr)

	return &UserServiceServer{
		Usecase:  usecase,
		AuthConn: authpb.NewAuthServiceClient(authConn),
	}
}

func (user *UserServiceServer) Signupuser(ctx context.Context, req *userpb.SignupUserRequest) (*userpb.UserResponce, error) {

	res, err := user.Usecase.SignupUser(entities.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Mobile,
	})
	if err != nil {
		helpers.PrintErr(err, "error at signup user usecase")
		return nil, errors.New("sorry for the inconveniance...Cannot onboard new users now. Please Try again later")
	}

	if _, err = user.AuthConn.InsertUser(ctx, &authpb.InsertUserReq{
		UserID:   res.UserID,
		Email:    res.Email,
		Password: req.Password,
	}); err != nil {
		helpers.PrintErr(err, "error at insertuser rpc call at signup user")
		return nil, err
	}

	return &userpb.UserResponce{
		Id:     res.UserID,
		Name:   res.Name,
		Email:  req.Email,
		Mobile: res.Phone,
	}, nil
}

func (user *UserServiceServer) GetRoles(emp *empty.Empty, stream userpb.UserService_GetRolesServer) error {

	res, err := user.Usecase.GetRoles()
	if err != nil {
		helpers.PrintErr(err, "error at getroles usecase")
		return errors.New("cannot handle user request now , please try again later...")
	}

	for _, v := range res {
		if err := stream.Send(&userpb.Role{
			ID:   uint32(v.ID),
			Role: v.Role,
		}); err != nil {
			helpers.PrintErr(err, "error at sending stream in getroles")
			return err
		}
	}

	return nil
}

func (user *UserServiceServer) SetStatus(ctx context.Context, req *userpb.StatusReq) (*empty.Empty, error) {

	if err := user.Usecase.SetStatus(entities.Status{
		UserID:    req.UserID,
		RoleID:    uint(req.RoleID),
		Available: req.IsAvailable,
	}); err != nil {
		helpers.PrintErr(err, "error at setstatus usecase")
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}

func (user *UserServiceServer) GetByEmail(ctx context.Context, req *userpb.GetByEmailReq) (*userpb.GetByEmailRes, error) {

	res, err := user.Usecase.GetIDbyEmail(req.Email)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetIDbyEmail usecase")
		return nil, err
	}

	return &userpb.GetByEmailRes{
		UserID: res,
	}, nil
}
