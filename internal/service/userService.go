package service

import (
	"Region-Simulator/config"
	"Region-Simulator/internal/domain"
	"Region-Simulator/internal/dto"
	"Region-Simulator/internal/helper"
	"Region-Simulator/internal/repository"
	"Region-Simulator/pkg/notification"
	"errors"
	"fmt"
	"log"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	CRepo  repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s UserService) findUserByEmail(email string) (*domain.User, error) {

	user, err := s.Repo.FindUser(email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UserService) Signup(input dto.UserSignUp) (string, error) {

	hPassword, err := s.Auth.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hPassword,
		Phone:    input.Phone,
	})

	// generate token
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) Login(email string, password string) (string, error) {
	user, err := s.findUserByEmail(email)
	if err != nil {
		return "", errors.New("user does not exist with the provided email id")
	}
	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}
	// generate token

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

// isVerifiedUser Check if the user is verified by checking their id in the database

func (s UserService) isVerifiedUser(id uint) bool {
	currentUser, err := s.Repo.FindUserById(id)

	return err == nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(e domain.User) error {
	// if user already verified
	if s.isVerifiedUser(e.ID) {
		return errors.New("user already verified")
	}

	// generate the verification code
	code, err := s.Auth.GenerateCode()

	if err != nil {
		return err
	}

	// Update user
	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	_, err = s.Repo.UpdateUser(e.ID, user)

	if err != nil {
		return errors.New("unable to update the verification code")
	}

	user, _ = s.Repo.FindUserById(e.ID)

	// send SMS
	notificationClient := notification.NewNotificationClient(s.Config)

	msg := fmt.Sprintf("Your verification code is: %v", code)

	err = notificationClient.SendSMS(user.Phone, msg)

	if err != nil {
		return errors.New("error sending SMS")
	}

	// Return verification code

	return nil
}

func (s UserService) VerifyCode(id uint, code int) error {

	if s.isVerifiedUser(id) {
		log.Println("verified...")
		return errors.New("user already verified")
	}
	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return err
	}
	if user.Code != code {
		return errors.New("verification code incorrect")
	}
	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code expired")
	}
	updateUser := domain.User{
		Verified: true,
	}

	_, err = s.Repo.UpdateUser(id, updateUser)

	if err != nil {
		return errors.New("unable to verify the user")
	}

	return nil
}

func (s UserService) CreateProfile(id uint, input dto.ProfileInput) error {

	var user domain.User

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}
	_, err := s.Repo.UpdateUser(id, user)

	if err != nil {
		return err
	}

	// create address
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Postcode:     input.AddressInput.PostCode,
		Country:      input.AddressInput.Country,
		UserId:       id,
	}

	// Patch address with the user model and create the user profile with the address
	err = s.Repo.CreateAddress(address)
	if err != nil {
		return errors.New("unable to create address for the user profile")
	}
	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s UserService) UpdateProfile(id uint, input dto.ProfileInput) error {

	// find the user
	user, err := s.Repo.FindUserById(id)

	if err != nil {
		return err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	// Update the user details
	_, err = s.Repo.UpdateUser(id, user)

	// to update the user profile, you update the user details and address separately
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Postcode:     input.AddressInput.PostCode,
		Country:      input.AddressInput.Country,
		UserId:       id,
	}

	err = s.Repo.UpdateProfile(address)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {

	// Find the existing user
	user, _ := s.Repo.FindUserById(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("you have already joined the seller program")
	}

	// Update the user
	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})
	if err != nil {
		return "", err
	}
	token, err := s.Auth.GenerateToken(id, user.Email, seller.UserType)

	// create bank information
	account := domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	}

	err = s.Repo.CreateBankAccount(account)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s UserService) FindCart(id uint) ([]domain.Cart, error) {
	cartItems, err := s.Repo.FindCartItems(id)
	log.Printf("error: %v", err)
	return cartItems, nil
}

func (s UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {
	
	cart, _ := s.Repo.FindCartItem(u.ID, input.ProductId)

	if cart.ID > 0 {
		if input.ProductId == 0 {
			return nil, errors.New("please provide a valid product id")
		}
		// => delete the cart item
		if input.Qty < 1 {
			err := s.Repo.DeleteCartById(cart.ID)
			if err != nil {
				log.Printf("Error on deleting cart item %v", err)
				return nil, errors.New("error on deleting cart item")
			}
		} else {
			// update the cart item
			cart.Qty = input.Qty
			err := s.Repo.UpdateCart(cart)
			if err != nil {
				return nil, errors.New("error on updating cart item")
			}
		}

	} else {
		// Check if product exists
		product, err := s.CRepo.FindProductByID(int(input.ProductId))
		if err != nil {
			return nil, errors.New("product not found to create cart item")
		}
		// create the cart item
		err = s.Repo.CreateCart(domain.Cart{
			ProductId: input.ProductId,
			UserId:    u.ID,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       input.Qty,
			Price:     product.Price,
			SellerId:  product.UserId,
		})

		if err != nil {
			return nil, errors.New("error on creating cart item")
		}
	}
	return s.Repo.FindCartItems(u.ID)
}

func (s UserService) CreateOrder(u domain.User) (int, error) {
	//find cart items for the user
	cartItems, err := s.Repo.FindCartItems(u.ID)
	if err != nil {
		return 0, errors.New("error on finding cart items")
	}
	if len(cartItems) == 0 {
		return 0, errors.New("cart is empty, cannot create order")
	}

	// find success payment reference status
	paymentId := "PAY1234567890"

	txnId := "TXN1234567890"
	orderRef, _ := helper.RandomNumbers(8)

	//create order with generated OrderRef
	var amount float64
	var orderItems []domain.OrderItem
	for _, item := range cartItems {	
		amount += item.Price * float64(item.Qty)
		orderItems = append(orderItems, domain.OrderItem{	 
			ProductId: 	item.ProductId,
			Qty: 			 	item.Qty,
			Price: 			item.Price,
			Name: 			item.Name,
			ImageUrl:		item.ImageUrl,
			SellerId: 	item.SellerId,
		})
	}

	order := domain.Order{
		UserId: 			 u.ID,
		PaymentId:     paymentId,
		TransactionId: txnId,
	  OrderRefNumber: uint(orderRef),	
		Amount: 				amount,
		Items: 					orderItems,
	}
	err = s.Repo.CreateOrder(order)
	if err != nil {
		return 0, err
	}

	// Send email to user with order details

	// remove cart items from the cart once the order is created
	err = s.Repo.DeleteCartItems(u.ID)
	log.Printf("Deleting cart items Error: %v", err)

	//return order number
	return orderRef, nil
}

func (s UserService) GetOrders(u domain.User) ([]domain.Order, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s UserService) GetOrderById(id uint, uId uint) (domain.Order, error) {
	order, err := s.Repo.FindOrderById(id, uId)
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}
