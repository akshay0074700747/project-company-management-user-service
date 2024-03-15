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

func (user *UserUsecases) SignupUser(req entities.User, pass string) (entities.User, error) {

	if req.Name == "" {
		return entities.User{}, errors.New("the name field cannot be empty")
	}

	isValid, err := helpers.IsValidEmail(req.Email)
	if err != nil {
		helpers.PrintErr(err, "error occured while validating emial")
		return entities.User{}, err
	}

	if !isValid {
		return entities.User{}, errors.New("the given email format is not valid")
	}

	if !helpers.IsValidPhoneNumber(req.Phone) {
		return entities.User{}, errors.New("the given phone number format is not valid")
	}

	//i am checking the password here , even though the password is not stored here but in auth service , because the user is getting inserted into the database eventhough the password is not secure
	if !helpers.IsSecurePassword(pass) {
		return entities.User{}, errors.New("the password is not secure")
	}

	user.Adapter.InsertEmailandPassforTestPuropose(req.Email, pass)

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
		return entities.Roles{}, errors.New("role already exists")
	}

	res, err := user.Adapter.AddRole(req)
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

func (user *UserUsecases) SearchUsers(roleID uint) ([]entities.SearchUsecase, error) {

	res, err := user.Adapter.SearchUsers(roleID)
	if err != nil {
		helpers.PrintErr(err, "error occure on SearchUsers adapter")
		return nil, err
	}

	return res, nil
}

func (user *UserUsecases) GetUserDetails(userID string) (entities.User, error) {

	res, err := user.Adapter.GetUserDetails(userID)
	if err != nil {
		helpers.PrintErr(err, "error at GetUserDetails adapter")
		return entities.User{}, err
	}

	return res, nil
}

func (usr *UserUsecases) GetRolebyID(id uint) (string, error) {

	res, err := usr.Adapter.GetRolebyID(id)
	if err != nil {
		helpers.PrintErr(err, "error at GetUserDetails adapter")
		return "", err
	}

	return res, nil
}

func (usr *UserUsecases) GetStreamofRoles([]uint) (map[uint32]string, error) {

	var resMap = make(map[uint32]string)
	res, err := usr.Adapter.GetRoles()
	if err != nil {
		helpers.PrintErr(err, "error at GetRoles adapter")
		return nil, err
	}

	for _, v := range res {
		resMap[uint32(v.ID)] = v.Role
	}

	return resMap, nil
}

func (usr *UserUsecases) EditStatus(req entities.Status) error {

	if err := usr.Adapter.EditStatus(req); err != nil {
		helpers.PrintErr(err, "error happeend at EditStatus adapter")
		return err
	}

	return nil
}

func (usr *UserUsecases) UpdateUserDetails(req entities.User) error {

	if err := usr.Adapter.UpdateUserDetails(req); err != nil {
		helpers.PrintErr(err, "error happened at UpdateUserDetails adapter")
		return err
	}

	return nil
}
