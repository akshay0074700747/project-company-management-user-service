package adapters

import "github.com/akshay0074700747/projectandCompany_management_user-service/entities"

type UserAdapterInterfaces interface {
	SignupUser(entities.User) (entities.User, error)
	AddRole(entities.Roles) (entities.Roles, error)
	GetRoles() ([]entities.Roles, error)
	IsExistingRole(string) (bool, error)
	SetStatus(entities.Status) error
	UpdateStatus(entities.Status) error
	IsUserStatusExist(entities.Status) (bool, error)
	GetIDbyEmail(string)(string,error)
	SearchUsers(uint)([]entities.SearchUsecase,error)
	GetUserDetails(string) (entities.User,error)
	GetRolebyID(uint)(string,error)
	InsertEmailandPassforTestPuropose(string,string)
	EditStatus(entities.Status)(error)
	UpdateUserDetails(entities.User)(error)
}
