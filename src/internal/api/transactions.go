package api

import (
	"fmt"
	"saldri/backend-saldri-andika-putra/domain"
	"saldri/backend-saldri-andika-putra/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type transactionApi struct {
	service     domain.TransactionService
	userService domain.UsersService
	store       *session.Store
}

func NewTransactionApi(app *fiber.App, service domain.TransactionService, userService domain.UsersService, jwtMid fiber.Handler, store *session.Store) {
	h := &transactionApi{
		service:     service,
		userService: userService,
		store:       store,
	}

	group := app.Group("/api/transactions", jwtMid)

	group.Get("/customers", h.GetCustomersByMerchant)
	group.Post("/", h.Create)
	group.Get("/", h.GetMyTransactions)
}

func (h *transactionApi) Create(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal mendapatkan sesi"))
	}

	userIDRaw := sess.Get("user_id")
	if userIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.CreateResponseError("user belum login"))
	}
	userID := userIDRaw.(string)

	var req dto.CreateTransactionRequest

	fmt.Println("BODY RAW:", string(c.Body()))

	if err := c.BodyParser(&req); err != nil {
		fmt.Println("UNMARSHAL ERROR:", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError("format data tidak valid"))
	}

	fmt.Printf("Parsed Request: %+v\n", req)

	if err := h.service.Create(c.Context(), req, userID); err != nil {
		fmt.Println("ERROR TX:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal membuat transaksi: " + err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess("transaksi berhasil dibuat"))
}

func (h *transactionApi) GetMyTransactions(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal mendapatkan sesi"))
	}

	userIDRaw := sess.Get("user_id")
	if userIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.CreateResponseError("user belum login"))
	}
	userID := userIDRaw.(string)

	txList, err := h.service.GetUserTransactions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal mengambil transaksi: " + err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess(txList))
}

func (h *transactionApi) GetCustomersByMerchant(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal mendapatkan sesi"))
	}

	userIDRaw := sess.Get("user_id")
	if userIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.CreateResponseError("user belum login"))
	}
	merchantID := userIDRaw.(string)

	// Opsional: Validasi role "merchant" jika dibutuhkan
	user, err := h.userService.GetProfile(c.Context(), merchantID)
	if err != nil || user.Role != "merchant" {
		return c.Status(fiber.StatusForbidden).JSON(dto.CreateResponseError("akses ditolak"))
	}

	customers, err := h.service.GetCustomersByMerchantID(c.Context(), merchantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal mengambil data customer: " + err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess(customers))
}
