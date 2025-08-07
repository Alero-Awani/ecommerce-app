package service

import (
	"ecommerce-app/config"
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"

	"github.com/pkg/errors"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {
	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	})
	return err
}

func (s CatalogService) EditCategory(id int, input dto.CreateCategoryRequest) (*domain.Category, error) {

	existCat, err := s.Repo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}
	if len(input.Name) > 0 {
		existCat.Name = input.Name
	}
	if len(input.ImageUrl) > 0 {
		existCat.ImageUrl = input.ImageUrl
	}
	if input.ParentId > 0 {
		existCat.ParentId = input.ParentId
	}
	if input.DisplayOrder > 0 {
		existCat.DisplayOrder = input.DisplayOrder
	}
	updatedCat, err := s.Repo.EditCategory(existCat)

	return updatedCat, err
}

func (s CatalogService) DeleteCategory(id int) error {
	err := s.Repo.DeleteCategory(id)
	if err != nil {
		return errors.New("category does not exist to delete")
	}
	return nil
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (s CatalogService) GetCategory(id int) (*domain.Category, error) {
	category, err := s.Repo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}
	return category, nil
}

// Products

func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		ImageUrl:    input.ImageUrl,
		Description: input.Description,
		CategoryId:  input.CategoryId,
		Price:       input.Price,
		Stock:       uint(input.Stock),
		UserId:      user.ID,
	})
	return err
}

func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, errors.New("could not fetch products")
	}
	return products, nil
}
