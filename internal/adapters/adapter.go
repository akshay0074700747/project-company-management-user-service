package adapters

import (
	"github.com/akshay0074700747/projectandCompany_management_user-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_user-service/helpers"
	"gorm.io/gorm"
)

type UserAdapter struct {
	DB *gorm.DB
}

func NewUserAdapter(db *gorm.DB) *UserAdapter {
	return &UserAdapter{
		DB: db,
	}
}

func (user *UserAdapter) SignupUser(req entities.User) (entities.User, error) {

	query := "INSERT INTO users (user_id,name,email,phone) VALUES($1,$2,$3,$4) RETURNING user_id,name,email,phone"
	var res entities.User

	tx := user.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := user.DB.Raw(query, req.UserID, req.Name, req.Email, req.Phone).Scan(&res).Error; err != nil {
		helpers.PrintErr(err, "error on signup adapter")
		tx.Rollback()
		return entities.User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}

func (user *UserAdapter) AddRole(req entities.Roles) (entities.Roles, error) {

	query := "INSERT INTO roles (role) VALUES($1) RETURNING id,role"
	var res entities.Roles

	tx := user.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := user.DB.Raw(query, req.Role).Scan(&res).Error; err != nil {
		helpers.PrintErr(err, "error on add role adapter")
		tx.Rollback()
		return entities.Roles{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}

func (user *UserAdapter) GetRoles() ([]entities.Roles, error) {

	query := "SELECT * FROM roles"
	var res []entities.Roles

	if err := user.DB.Raw(query).Scan(&res).Error; err != nil {
		helpers.PrintErr(err, "error at getroles adapter")
		return []entities.Roles{}, err
	}

	return res, nil
}

func (user *UserAdapter) IsExistingRole(role string) (bool, error) {

	query := "SELECT * FROM roles WHERE role = $1"

	res := user.DB.Exec(query, role)
	if res.Error != nil {
		helpers.PrintErr(res.Error, "error at IsExistingRole")
		return true, res.Error
	}

	if res.RowsAffected != 0 {
		return true, nil
	}

	return false, nil
}

func (user *UserAdapter) SetStatus(req entities.Status) error {

	query := "INSERT INTO statuses (user_id,role_id,available) VALUES($1,$2,$3)"

	if err := user.DB.Exec(query, req.UserID, req.RoleID, req.Available).Error; err != nil {
		return err
	}

	return nil
}

func (user *UserAdapter) IsUserStatusExist(req entities.Status) (bool, error) {

	query := "SELECT * FROM statuses WHERE user_id = $1 AND role_id = $2"

	res := user.DB.Exec(query, req.UserID, req.RoleID)
	if res.Error != nil {
		return true, res.Error
	}

	if res.RowsAffected != 0 {
		return true, nil
	}

	return false, nil
}

func (user *UserAdapter) UpdateStatus(req entities.Status) error {

	query := "UPDATE statuses SET available = $1 WHERE user_id = $2 AND role_id = $3"

	if err := user.DB.Exec(query, req.Available, req.UserID, req.RoleID).Error; err != nil {
		return err
	}

	return nil
}

func (user *UserAdapter) GetIDbyEmail(email string) (string, error) {

	query := "SELECT user_id FROM users WHERE email = $1"
	var res string

	if err := user.DB.Raw(query, email).Scan(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (user *UserAdapter) SearchUsers(roleID uint) ([]entities.SearchUsecase, error) {

	query := "SELECT a.name,a.email FROM users a INNER JOIN statuses b ON a.user_id = b.user_id AND b.role_id = $1 AND b.available = true"
	var res []entities.SearchUsecase

	if err := user.DB.Raw(query, roleID).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (user *UserAdapter) GetUserDetails(userID string) (entities.User, error) {

	query := "SELECT * FROM users WHERE user_id = $1"
	var res entities.User

	if err := user.DB.Raw(query, userID).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}

	return res, nil
}

func (usr *UserAdapter) GetRolebyID(id uint) (string, error) {

	query := "SELECT role from roles where id = $1"
	var res string

	if err := usr.DB.Raw(query, id).Scan(&res).Error; err != nil {
		return "", err
	}

	return res, nil
}

func (usr *UserAdapter) InsertEmailandPassforTestPuropose(email, pass string) {

	query := "INSERT INTO tst_user_table (email,password) VALUES($1,$2)"
	if err := usr.DB.Exec(query, email, pass).Error; err != nil {
		helpers.PrintErr(err, "")
	}
}

func (usr *UserAdapter) EditStatus(req entities.Status) error {

	if err := usr.DB.Model(&entities.Status{}).Where("user_id = ?", req.UserID).Updates(req).Error; err != nil {
		return err
	}

	return nil
}

func (usr *UserAdapter) UpdateUserDetails(req entities.User)(error) {
	
	if err := usr.DB.Model(&entities.User{}).Where("user_id = ?", req.UserID).Updates(req).Error; err != nil {
		return err
	}

	return nil
}
