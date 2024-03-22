package injectdependency

import (
	"github.com/akshay0074700747/projectandCompany_management_user-service/config"
	"github.com/akshay0074700747/projectandCompany_management_user-service/db"
	"github.com/akshay0074700747/projectandCompany_management_user-service/internal/adapters"
	"github.com/akshay0074700747/projectandCompany_management_user-service/internal/services"
	"github.com/akshay0074700747/projectandCompany_management_user-service/internal/usecases"
)

func Initialize(cfg config.Config) *services.UserEngine {

	db := db.ConnectDB(cfg)
	adapter := adapters.NewUserAdapter(db)
	usecase := usecases.NewUserUsecases(adapter)
	server := services.NewUserServiceServer(usecase, "auth-service:50004")

	return services.NewUserEngine(server)
}
