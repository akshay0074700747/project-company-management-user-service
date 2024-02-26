package db

import (
	"fmt"
	"log"

	"github.com/akshay0074700747/projectandCompany_management_user-service/config"
	"github.com/akshay0074700747/projectandCompany_management_user-service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.Config) *gorm.DB {

	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBhost, cfg.DBuser, cfg.DBname, cfg.DBport, cfg.DBpassword)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal("cannot connect to the db ", err)
	}

	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Roles{})
	db.AutoMigrate(&entities.RolesandCompanies{})
	db.AutoMigrate(&entities.Status{})
	db.AutoMigrate(&entities.RolesandProjects{})

	return db
}
