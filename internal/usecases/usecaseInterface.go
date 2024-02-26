package usecases

import "github.com/akshay0074700747/projectandCompany_management_user-service/entities"

type UserUsecaseInterfaces interface {
	SignupUser(entities.User,string) (entities.User, error)
	AddRole(entities.Roles) (entities.Roles, error)
	GetRoles() ([]entities.Roles, error)
	SetStatus(entities.Status) error
	GetIDbyEmail(string) (string,error)
	SearchUsers(uint) ([]entities.SearchUsecase, error)
	GetUserDetails(string) (entities.User, error)
	GetRolebyID(uint)(string,error)
	GetStreamofRoles([]uint)(map[uint32]string,error)
}
