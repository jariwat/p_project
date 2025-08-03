package main

import (
	"fmt"
	"jariwat/p_project/profile-service/helper"
	"log"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

var (
	APP_PORT     = helper.GetENV("APP_PORT", "3000")
	DB_HOST      = helper.GetENV("DB_HOST", "localhost")
	DB_NAME      = helper.GetENV("DB_NAME", "app_example")
	DB_USER      = helper.GetENV("DB_USER", "postgres")
	DB_PORT      = helper.GetENV("DB_PORT", "5432")
	DB_PASSWORD  = helper.GetENV("DB_PASSWORD", "postgres")
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

	e := echo.New()
	e.HTTPErrorHandler = helperMiddl.SentryCapture(e)
	helperRoute.RegisterVersion(e)

	e.Use(echoMiddL.RequestLoggerWithConfig(echoMiddL.RequestLoggerConfig{
		LogError: true,
		LogValuesFunc: func(c echo.Context, values echoMiddL.RequestLoggerValues) error {
			if values.Error != nil {
				c.Logger().Error(values.Error)
			}
			return nil
		},
	}))
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/health-check", func(c echo.Context) error {
		resp := echo.Map{
			"status": "ok",
		}
		return c.JSON(http.StatusOK, resp)
	})

	/* repository */
	userRepo := user_repository.NewPsqlUserRepository(psqlClient)

	/* usecase */
	userUsecase := user_usecase.NewUserUsecase(userRepo)

	// init openapi middleware here
	mw, err := myMiddL.CreateOpenapiMiddleware(
		jwks,
		userUsecase,
		user.GetSwagger,
		company.GetSwagger,
	)
	if err != nil {
		panic(err)
	}
	e.Use(mw)

	/* handler */
	userHandler := user_handler.NewUserHandler(userUsecase)

	/* inject route */
	user.RegisterHandlers(e, userHandler)

	/* serve */
	port := fmt.Sprintf(":%s", APP_PORT)
	e.Logger.Fatal(e.Start(port))
}
