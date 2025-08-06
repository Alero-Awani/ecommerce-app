package handlers

import (
	"Region-Simulator/internal/api/rest"
	"Region-Simulator/internal/dto"
	"Region-Simulator/internal/repository"
	"Region-Simulator/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type catalogHandler struct {
	svc service.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	// Create an instance of the catalog service and inject to the handler
	svc := service.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}

	handler := catalogHandler{
		svc: svc,
	}

	// Public Catalog Endpoints
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)
	app.Get("categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategoryById)

	// Private Catalog Endpoints
	selRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)
	// Categories
	selRoutes.Post("/categories", handler.CreateCategories)
	selRoutes.Patch("/categories/:id", handler.EditCategory)
	selRoutes.Delete("/categories/:id", handler.DeleteCategory)

	// Products
	selRoutes.Post("/products", handler.CreateProducts)
	selRoutes.Get("/products", handler.GetProducts)
	selRoutes.Get("/products/:id", handler.GetProduct)
	selRoutes.Put("/products/:id", handler.EditProducts)
	selRoutes.Patch("/products/:id", handler.UpdateStock) // update stock
	selRoutes.Delete("/products/:id", handler.DeleteProduct)
}
func (h *catalogHandler) GetCategories(ctx *fiber.Ctx) error {
	cats, err := h.svc.GetCategories()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "categories", cats)
}

func (h *catalogHandler) GetCategoryById(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	cat, err := h.svc.GetCategory(id)
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "category endpoint", cat)
}

func (h *catalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "create category request is not valid")
	}
	err = h.svc.CreateCategory(req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "category created successfully", nil)
}

func (h *catalogHandler) EditCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	req := dto.CreateCategoryRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "update category request is not valid")
	}
	updatedCategory, err := h.svc.EditCategory(id, req)

	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "edit category endpoint", updatedCategory)
}

func (h *catalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := h.svc.DeleteCategory(id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "delete category endpoint", nil)
}

func (h *catalogHandler) CreateProducts(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.CreateProductRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "create product request is not valid")
	}
	err = h.svc.CreateProduct(req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "product created successfully", nil)
}

func (h *catalogHandler) EditProducts(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, "edit product endpoint", nil)
}

func (h *catalogHandler) GetProducts(ctx *fiber.Ctx) error {
	products, err := h.svc.GetProducts()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "categories", products)
}

func (h *catalogHandler) GetProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "get product by ID", nil)
}

func (h *catalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "update stock endpoint", nil)
}

func (h *catalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "delete product endpoint", nil)

}
