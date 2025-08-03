package main

import (
	"fmt"
	"jariwat/p_project/profile-service/helper"
	"log"
	"net/http"

	myMiddL "jariwat/p_project/profile-service/middleware"
	"jariwat/p_project/profile-service/service/profile"
	profile_repository "jariwat/p_project/profile-service/service/profile/repository"
	profile_usecase "jariwat/p_project/profile-service/service/profile/usecase"
	profile_handler "jariwat/p_project/profile-service/service/profile/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	APP_PORT    = helper.GetENV("APP_PORT", "3000")
	DB_HOST     = helper.GetENV("DB_HOST", "localhost")
	DB_NAME     = helper.GetENV("DB_NAME", "app_example")
	DB_USER     = helper.GetENV("DB_USER", "postgres")
	DB_PORT     = helper.GetENV("DB_PORT", "5432")
	DB_PASSWORD = helper.GetENV("DB_PASSWORD", "postgres")
)

func gormDB() gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", DB_HOST, DB_USER, DB_NAME, DB_PORT, DB_PASSWORD)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

func main() {
	psqlClient := gormDB()

	g := gin.Default()

	g.GET("/", func(c *gin.Context) {
		return c.String(http.StatusOK, "Hello, World!")
	})

	g.GET("/health-check", func(c *gin.Context) {
		resp := gin.H{
			"status": "ok",
		}
		return c.JSON(http.StatusOK, resp)
	})

	// init openapi middleware here
	mw, err := myMiddL.CreateOpenapiMiddleware(profile.GetSwagger)
	if err != nil {
		panic(err)
	}
	g.Use(mw)

	/* repository */
	profileRepo := profile_repository.NewPsqlProfileRepository(psqlClient)

	/* usecase */
	profileUsecase := profile_usecase.NewProfileUsecase(profileRepo)

	/* handler */
	profileHandler := profile_handler.NewProfileHandler(profileUsecase)

	/* inject route */
	profile.RegisterHandlers(g, profileHandler)

	/* serve */
	port := fmt.Sprintf(":%s", APP_PORT)
	g.Logger.Fatal(g.Start(port))
}
