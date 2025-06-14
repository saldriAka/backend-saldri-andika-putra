package api

import (
	"context"
	"saldri/backend-saldri-andika-putra/domain"
	"saldri/backend-saldri-andika-putra/dto"
	"saldri/backend-saldri-andika-putra/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type usersApi struct {
	usersService domain.UsersService
	store        *session.Store
}

func NewUsersApi(app *fiber.App, usersService domain.UsersService, auth fiber.Handler, store *session.Store) {
	u := &usersApi{usersService: usersService, store: store}

	app.Post("/auth/register", u.Register)
	app.Post("/auth/login", u.Login)

	users := app.Group("/api/users", auth)
	users.Get("/:id", u.Profile)
	users.Get("/", u.List)
}

func (u *usersApi) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError("invalid body"))
	}
	if fails := util.Validate(req); len(fails) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseErrorData("validation error", fails))
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	if err := u.usersService.Register(ctx, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess("user registered"))
}

func (u *usersApi) Login(c *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError("invalid body"))
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	res, err := u.usersService.Login(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.CreateResponseError("login failed: " + err.Error()))
	}

	sess, err := u.store.Get(c)
	if err == nil {
		sess.Set("user_id", res.ID)
		sess.Set("role", res.Role)
		sess.Set("name", res.Name)
		sess.Set("email", res.Email)
		sess.Save()
	}

	return c.JSON(dto.CreateResponseSuccess(res))
}

func (u *usersApi) Profile(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	userID := c.Params("id")

	data, err := u.usersService.GetProfile(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.CreateResponseError(err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess(data))
}

func (u *usersApi) List(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	page, limit, _ := util.SafePaginationParams(c)
	users, total, err := u.usersService.List(ctx, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return dto.PaginateAndRespond(c, users, int(total))
}
