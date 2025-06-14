package api

import (
	"fmt"
	"saldri/backend-saldri-andika-putra/domain"
	"saldri/backend-saldri-andika-putra/dto"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type productApi struct {
	service domain.ProductService
	store   *session.Store
}

func NewProductApi(app *fiber.App, service domain.ProductService, auth fiber.Handler, store *session.Store) {
	handler := &productApi{service: service, store: store}

	products := app.Group("/api/products", auth)

	products.Get("/", handler.List)
	products.Get("/:id", handler.Detail)
	products.Post("/", handler.Create)
	products.Put("/:id", IsMerchant(store), handler.Update)
	products.Delete("/:id", IsMerchant(store), handler.Delete)
}

// ✅ List semua produk
func (h *productApi) List(c *fiber.Ctx) error {
	products, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal mengambil data produk"))
	}
	return c.JSON(dto.CreateResponseSuccess(products))
}

// ✅ Ambil detail produk
func (h *productApi) Detail(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError("ID produk tidak valid"))
	}

	product, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.CreateResponseError("produk tidak ditemukan"))
	}
	return c.JSON(dto.CreateResponseSuccess(product))
}

// ✅ Buat produk baru (khusus merchant)
func (h *productApi) Create(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError("Gagal mengambil session"))
	}

	roleVal := sess.Get("role")
	role, ok := roleVal.(string)
	if !ok || role != "merchant" {
		return c.Status(fiber.StatusForbidden).
			JSON(dto.CreateResponseError("Hanya merchant yang dapat membuat produk"))
	}

	userIDVal := sess.Get("user_id")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(dto.CreateResponseError("Pengguna tidak valid"))
	}

	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError("Format data tidak valid"))
	}

	req.MerchantID = userID

	if err := h.service.Create(c.Context(), req); err != nil {
		fmt.Println("Error saat membuat produk:", err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError("Gagal membuat produk: " + err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess("Produk berhasil dibuat"))
}

// ✅ Update produk (khusus merchant pemilik produk)
func (h *productApi) Update(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	merchantID := sess.Get("user_id").(string)

	idParam := c.Params("id")
	productID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError("ID tidak valid"))
	}

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError("format tidak valid"))
	}

	err = h.service.Update(c.Context(), productID, req, merchantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError(err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess("produk berhasil diupdate"))
}

// ✅ Hapus produk (khusus merchant pemilik produk)
func (h *productApi) Delete(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	merchantID := sess.Get("user_id").(string)

	idParam := c.Params("id")
	productID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError("ID tidak valid"))
	}

	err = h.service.Delete(c.Context(), productID, merchantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.CreateResponseError(err.Error()))
	}

	return c.JSON(dto.CreateResponseSuccess("produk berhasil dihapus"))
}

// ✅ Middleware: hanya merchant yang bisa akses
func IsMerchant(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.CreateResponseError("gagal ambil sesi"))
		}

		role := sess.Get("role")
		if role == nil || role != "merchant" {
			return c.Status(fiber.StatusForbidden).JSON(dto.CreateResponseError("akses ditolak, bukan merchant"))
		}

		return c.Next()
	}
}
