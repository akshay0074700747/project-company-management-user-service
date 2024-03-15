package services

import (
	"context"
	"errors"
	"io"

	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/authpb"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"github.com/akshay0074700747/projectandCompany_management_user-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_user-service/helpers"
	"github.com/akshay0074700747/projectandCompany_management_user-service/internal/usecases"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (user *UserServiceServer) SignupUser(ctx context.Context, req *userpb.SignupUserRequest) (*userpb.UserResponce, error) {

	res, err := user.Usecase.SignupUser(entities.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Mobile,
	}, req.Password)
	if err != nil {
		helpers.PrintErr(err, "error at signup user usecase")
		return nil, err
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
		return errors.New("cannot handle user request now , please try again later")
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

func (user *UserServiceServer) AddRoles(ctx context.Context, req *userpb.AddRoleReq) (*emptypb.Empty, error) {

	if _, err := user.Usecase.AddRole(entities.Roles{
		Role: req.Role,
	}); err != nil {
		helpers.PrintErr(err, "error at AddRole usecase")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (user *UserServiceServer) SearchforMembers(req *userpb.SearchReq, stream userpb.UserService_SearchforMembersServer) error {

	res, err := user.Usecase.SearchUsers(uint(req.RoleID))
	if err != nil {
		helpers.PrintErr(err, "error occured at SearchUsers usecase")
		return errors.New("cannot search now... please try agin later")
	}

	for _, v := range res {
		if err := stream.Send(&userpb.SearchRes{
			Name:  v.Name,
			Email: v.Email,
		}); err != nil {
			helpers.PrintErr(err, "error at sending stream")
			return errors.New("cannot serch now please try again later")
		}
	}

	return nil
}

func (user *UserServiceServer) GetUserDetails(mctx context.Context, req *userpb.GetUserDetailsReq) (*userpb.GetUserDetailsRes, error) {

	var err error
	res, err := user.Usecase.GetUserDetails(req.UserID)
	if err != nil {
		helpers.PrintErr(err, "error at GetUserDetails usecase")
		return nil, err
	}

	result := &userpb.GetUserDetailsRes{
		UserID: req.UserID,
		Email:  res.Email,
		Mobile: res.Phone,
		Name:   res.Name,
	}

	if req.RoleID != 0 {
		result.Role, err = user.Usecase.GetRolebyID(uint(req.RoleID))
		if err != nil {
			helpers.PrintErr(err, "error at GetRolebyID usecase")
			return nil, err
		}
	}

	return result, err
}

func (user *UserServiceServer) GetRole(ctx context.Context, req *userpb.GetRoleReq) (*userpb.GetRoleRes, error) {

	role, err := user.Usecase.GetRolebyID(uint(req.ID))
	if err != nil {
		helpers.PrintErr(err, "error at GetRolebyID usecase")
		return nil, err
	}

	return &userpb.GetRoleRes{
		Role: role,
	}, nil
}

func (user *UserServiceServer) GetStreamofRoles(stream userpb.UserService_GetStreamofRolesServer) error {
	var roleIDs []uint
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			helpers.PrintErr(err, "error at recieving from stream")
			return err
		}
		roleIDs = append(roleIDs, uint(req.RoleID))
	}

	res, err := user.Usecase.GetStreamofRoles(roleIDs)
	if err != nil {
		helpers.PrintErr(err, "error at GetStreamofRoles usecase")
		return err
	}

	if err := stream.SendAndClose(&userpb.GetStreamofRolesRes{
		RoleIDsWithNames: res,
	}); err != nil {
		helpers.PrintErr(err, "error at sending to stream")
		return err
	}

	return nil
}

func (user *UserServiceServer) GetStreamofUserDetails(stream userpb.UserService_GetStreamofUserDetailsServer) error {
	var role string
	var userr entities.User
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			helpers.PrintErr(err, "error at reci to stream")
			return err
		}
		if req.RoleID != 0 {
			role, err = user.Usecase.GetRolebyID(uint(req.RoleID))
			if err != nil {
				helpers.PrintErr(err, "error at GetRolebyID usecase")
				return err
			}
		}
		userr, err = user.Usecase.GetUserDetails(req.UserID)
		if err != nil {
			helpers.PrintErr(err, "error at GetUserDetails usecase")
			return err
		}

		if err = stream.Send(&userpb.GetUserDetailsRes{
			UserID: userr.UserID,
			Mobile: userr.Phone,
			Email:  userr.Email,
			Name:   userr.Name,
			Role:   role,
		}); err != nil {
			helpers.PrintErr(err, "error at sending to stream")
			return err
		}
	}

	return nil
}

func (user *UserServiceServer) EditStatus(ctx context.Context, req *userpb.EditStatusReq) (*emptypb.Empty, error) {

	if err := user.Usecase.EditStatus(entities.Status{
		RoleID:    uint(req.RoleID),
		Available: req.IsAvailable,
		UserID:    req.UserID,
	}); err != nil {
		helpers.PrintErr(err, "error happened at EditStatus usease")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (user *UserServiceServer) UpdateUserDetails(ctx context.Context, req *userpb.UpdateUserDetailsReq) (*emptypb.Empty, error) {

	if err := user.Usecase.UpdateUserDetails(entities.User{
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Mobile,
		UserID: req.UserID,
	}); err != nil {
		helpers.PrintErr(err, "errro happened at UpdateUserDetails adapter")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

