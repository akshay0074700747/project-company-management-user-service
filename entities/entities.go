package entities

type User struct {
	UserID string `gorm:"primaryKey"`
	Name   string
	Email  string `gorm:"unique;not null"`
	Phone  string `gorm:"unique;not null"`
}

type Roles struct {
	ID   uint   `gorm:"primaryKey"`
	Role string `gorm:"unique;not null"`
}

type Status struct {
	UserID    string `gorm:"foreignKey:UserID;references:users(user_id);constraint:OnDelete:CASCADE"`
	RoleID    uint   `gorm:"foreignKey:RoleID;references:roles(id);constraint:OnDelete:CASCADE"`
	Available bool
}

type RolesandProjects struct {
	UserID    string `gorm:"foreignKey:UserID;references:users(user_id);constraint:OnDelete:CASCADE"`
	RoleID    uint   `gorm:"foreignKey:RoleID;references:roles(id);constraint:OnDelete:CASCADE"`
	ProjectID string
}

type RolesandCompanies struct {
	UserID    string `gorm:"foreignKey:UserID;references:users(user_id);constraint:OnDelete:CASCADE"`
	RoleID    uint   `gorm:"foreignKey:RoleID;references:roles(id);constraint:OnDelete:CASCADE"`
	CompanyID string
}
