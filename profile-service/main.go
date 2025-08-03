package main

import (
	"fmt"
	"github.com/jariwat/p_project/profile-service/helper"
	"log"
	"net/http"

	myMiddL "github.com/jariwat/p_project/profile-service/middleware"
	"github.com/jariwat/p_project/profile-service/service/profile"
	profile_repository "github.com/jariwat/p_project/profile-service/service/profile/repository"
	profile_usecase "github.com/jariwat/p_project/profile-service/service/profile/usecase"
	profile_handler "github.com/jariwat/p_project/profile-service/service/profile/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	APP_PORT    = helper.GetENV("APP_PORT", "3000")
	DB_HOST     = helper.GetENV("DB_HOST", "localhost")
	DB_NAME     = helper.GetENV("DB_NAME", "app_example")
	DB_USER     = helper.GetENV("DB_USER", "postgres")
	DB_PORT     = helper.GetENV("DB_PORT", "5432")
	DB_PASSWORD = helper.GetENV("DB_PASSWORD", "postgres")
)


func runMigrations() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

func gormDB() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", DB_HOST, DB_USER, DB_NAME, DB_PORT, DB_PASSWORD)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	runMigrations()

	return db
}

func main() {
	psqlClient := gormDB()

	g := gin.Default()

	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
		return
	})

	g.GET("/health-check", func(c *gin.Context) {
		resp := gin.H{
			"status": "ok",
		}
		c.JSON(http.StatusOK, resp)
		return
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
	g.Run(port)
	log.Println("Server running on port", APP_PORT)
}
