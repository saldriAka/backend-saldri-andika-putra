package main

import (
	"fmt"
	"saldri/backend-saldri-andika-putra/dto"
	"saldri/backend-saldri-andika-putra/internal/api"
	"saldri/backend-saldri-andika-putra/internal/config"
	"saldri/backend-saldri-andika-putra/internal/connection"
	"saldri/backend-saldri-andika-putra/internal/repository"
	"saldri/backend-saldri-andika-putra/internal/service"

	"github.com/go-playground/validator/v10"
	jwtMid "github.com/gofiber/contrib/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	cnf := config.Get()
	store := session.New()
	dbConnection := connection.GetDatabase(cnf.Database)
	validate := validator.New()

	app := fiber.New()

	jwtMiddleware := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ContextKey: "user",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				dto.CreateResponseError("endpoint perlu token, silakan login"),
			)
		},
	})

	productRepo := repository.NewProductRepository(dbConnection)
	txRepo := repository.NewTransactionRepository(dbConnection)
	usersRepo := repository.NewUsersRepository(dbConnection)

	productService := service.NewProductService(productRepo)
	txService := service.NewTransactionService(txRepo, productRepo)

	usersService := service.NewUsersService(usersRepo, cnf)

	api.NewProductApi(app, productService, jwtMiddleware, store, validate)
	api.NewTransactionApi(app, txService, usersService, jwtMiddleware, store, validate)
	api.NewUsersApi(app, usersService, jwtMiddleware, store)

	fmt.Println("Listening on", cnf.Server.Host+":"+cnf.Server.Port)
	_ = app.Listen(":" + cnf.Server.Port)
}
