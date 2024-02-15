package usecases

import "github.com/akshay0074700747/projectandCompany_management_user-service/entities"

type UserUsecaseInterfaces interface {
	SignupUser(entities.User) (entities.User, error)
	AddRole(entities.Roles) (entities.Roles, error)
	GetRoles() ([]entities.Roles, error)
	SetStatus(entities.Status) error
	GetIDbyEmail(string) (string,error)
}
