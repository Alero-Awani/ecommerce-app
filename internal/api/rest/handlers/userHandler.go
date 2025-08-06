package handlers

import (
	"Region-Simulator/internal/api/rest"
	"Region-Simulator/internal/dto"
	"Region-Simulator/internal/repository"
	"Region-Simulator/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type userHandler struct {
	svc service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App
	// Create an instance of user service & inject to handler
	svc := service.UserService{
		Repo:   repository.NewUserRepository(rh.DB),
		CRepo:  repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}
	handler := userHandler{
		svc: svc,
	}

	pubRoutes := app.Group("/users")

	// Public endpoints
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)

	pvtRoutes := pubRoutes.Group("/", rh.Auth.Authorize)
	// Private endpoints
	pvtRoutes.Get("/verify", handler.GetVerificationCode)
	pvtRoutes.Post("/verify", handler.Verify)
	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Get("/profile", handler.GetProfile)
	pvtRoutes.Patch("/profile", handler.UpdateProfile)

	pvtRoutes.Post("/cart", handler.AddToCart)
	pvtRoutes.Get("/cart", handler.GetCart)

	pvtRoutes.Post("order", handler.CreateOrder)
	pvtRoutes.Get("order", handler.GetOrders)
	pvtRoutes.Get("/order/:id", handler.GetOrder)

	pvtRoutes.Post("/become-seller", handler.BecomeSeller)
}

func (h *userHandler) Register(ctx *fiber.Ctx) error {
	// to create user
	user := dto.UserSignUp{}
	err := ctx.BodyParser(&user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	token, err := h.svc.Signup(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error on signup",
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Register",
		"token":   token,
	})
}

func (h *userHandler) Login(ctx *fiber.Ctx) error {
	loginInput := dto.UserLogin{}
	err := ctx.BodyParser(&loginInput)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	token, err := h.svc.Login(loginInput.Email, loginInput.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please provide the correct login information",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login",
		"token":   token,
	})
}

func (h *userHandler) Verify(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	// request to accept the verification code
	var req dto.VerificationCodeInput

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide a valid input",
		})
	}

	err := h.svc.VerifyCode(user.ID, req.Code)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Verified Successfully",
	})
}

func (h *userHandler) GetVerificationCode(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	// Create the verification Code and update the user profile in DB
	err := h.svc.GetVerificationCode(user)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to get verification code",
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "the user now has a verification code",
	})
}

func (h *userHandler) CreateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.ProfileInput{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	log.Printf("User: %v", user)
	err := h.svc.CreateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to create profile",
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "profile created successfully",
	})
}

func (h *userHandler) GetProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	log.Println(user)

	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to get profile",
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "get profile",
		"profile": profile,
	})
}

func (h *userHandler) UpdateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	err := h.svc.UpdateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to update profile",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "profile updated successfully",
	})
}

func (h *userHandler) AddToCart(ctx *fiber.Ctx) error {
	req := dto.CreateCartRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please  provide a valid product and qty",
		})
	}
	user := h.svc.Auth.GetCurrentUser(ctx)
	log.Println(user)

	// call user service and perform create cart
	cartItems, err := h.svc.CreateCart(req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "cart created successfully", cartItems)
}

func (h *userHandler) GetCart(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	cart, err := h.svc.FindCart(user.ID)
	if err != nil {
		return rest.InternalError(ctx, errors.New("cart does not exist"))
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "get cart",
		"cart":    cart,
	})
}

func (h *userHandler) CreateOrder(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	orderRef, err := h.svc.CreateOrder(user)
	
	if err != nil {
		log.Println("Error creating order:", err)
		return rest.InternalError(ctx, errors.New("unable to create order"))
	}


	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "order created successfully",
		"orderRef": orderRef,
	})
}

func (h *userHandler) GetOrders(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	orders, err := h.svc.GetOrders(user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "get orders",
		"order": orders,
	})
}

func (h *userHandler) GetOrder(ctx *fiber.Ctx) error {
	orderId, _ := strconv.Atoi(ctx.Params("id"))
	user := h.svc.Auth.GetCurrentUser(ctx)

	order, err := h.svc.GetOrderById(uint(orderId), user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "get order",
		"order": order,
	})
}

func (h *userHandler) BecomeSeller(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.SellerInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	token, err := h.svc.BecomeSeller(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Failed to become seller",
		})
	}

	// Convert User to Seller

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "become seller",
		"token":   token,
	})
}
