package usecases

import (
	"errors"

	"github.com/akshay0074700747/projectandCompany_management_user-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_user-service/helpers"
	"github.com/akshay0074700747/projectandCompany_management_user-service/internal/adapters"
)

type UserUsecases struct {
	Adapter adapters.UserAdapterInterfaces
}

func NewUserUsecases(adapter adapters.UserAdapterInterfaces) *UserUsecases {
	return &UserUsecases{
		Adapter: adapter,
	}
}

func (user *UserUsecases) SignupUser(req entities.User) (entities.User, error) {

	if req.Name == "" {
		return entities.User{}, errors.New("the name field cannot be empty")
	}

	isValid, err := helpers.IsValidEmail(req.Email)
	if err != nil {
		helpers.PrintErr(err, "error occured while validating emial")
		return entities.User{}, errors.New("cannot verify email right now!")
	}

	if !isValid {
		return entities.User{}, errors.New("the given email format is not valid!")
	}

	if !helpers.IsValidPhoneNumber(req.Phone) {
		return entities.User{}, errors.New("the given phone number format is not valid!")
	}

	req.UserID = helpers.GenUuid()

	res, err := user.Adapter.SignupUser(req)
	if err != nil {
		helpers.PrintErr(err, "error happened at signupuser adapter")
		return entities.User{}, err
	}

	return res, nil
}

func (user *UserUsecases) AddRole(req entities.Roles) (entities.Roles, error) {

	isExisting, err := user.Adapter.IsExistingRole(req.Role)
	if err != nil {
		return entities.Roles{}, err
	}

	if isExisting {
		return entities.Roles{}, errors.New("role already exists...")
	}

	res, err := user.AddRole(req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (user *UserUsecases) GetRoles() ([]entities.Roles, error) {

	res, err := user.Adapter.GetRoles()
	if err != nil {
		return []entities.Roles{}, err
	}

	return res, nil
}

func (user *UserUsecases) SetStatus(req entities.Status) error {

	isExisting, err := user.Adapter.IsUserStatusExist(req)
	if err != nil {
		helpers.PrintErr(err, "error at IsUserStatusExist adpter")
		return err
	}

	if isExisting {
		if err := user.Adapter.UpdateStatus(req); err != nil {
			helpers.PrintErr(err, "erorr at UpdateStatus adapter")
			return err
		}
	} else {
		if err := user.Adapter.SetStatus(req); err != nil {
			helpers.PrintErr(err, "erorr at SetStatus adapter")
			return err
		}
	}

	return nil
}

func (user *UserUsecases) GetIDbyEmail(email string) (string, error) {

	res, err := user.Adapter.GetIDbyEmail(email)
	if err != nil {
		helpers.PrintErr(err, "error at GetIDbyEmail adapter")
		return res, err
	}

	return res, nil
}
