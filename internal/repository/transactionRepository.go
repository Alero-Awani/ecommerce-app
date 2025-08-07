package repository

import (
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

type transactionStorage struct {
	db *gorm.DB
}


func (t *transactionStorage) CreatePayment(payment *domain.Payment) error {
 err := t.db.Create(payment).Error
 if err != nil {
	return err
 }
 return nil
}

func (t *transactionStorage) FindOrders(uId uint) ([]domain.OrderItem, error) {
	var orders []domain.OrderItem
	err := t.db.Model(&domain.OrderItem{}).
	Where("user_id = ?", uId).
	Preload("Product").
	Preload("Seller").
	Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}		

func (t *transactionStorage) FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error) {
	var order dto.SellerOrderDetails
	err := t.db.Model(&domain.OrderItem{}).
	Where("user_id = ? AND id = ?", uId, id).
	Preload("Product").
	Preload("Seller").
	First(&order).Error
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}
	return order, nil
}	

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
 return &transactionStorage{db: db}
}