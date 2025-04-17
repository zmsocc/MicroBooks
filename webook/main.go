package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/practice/webook/internal/repository"
	"github.com/zmsocc/practice/webook/internal/repository/dao"
	"github.com/zmsocc/practice/webook/internal/service"
	"github.com/zmsocc/practice/webook/internal/web"
	"gorm.io/gorm"
)

func main() {
	db := dao.InitDB()
	server := web.InitWebServer()
	InitUser(server, db)
	server.Run(":8080")
}

func InitUser(server *gin.Engine, db *gorm.DB) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	u := web.NewUserHandler(us)
	u.RegisterRoutes(server)
}
