package product

import (
	"context"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(5) * time.Second,
	}
}

func (s *service) CreateCategory(c context.Context, req *Category) (*Category, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u := &Category{
		Name: req.Name,
	}
	r, err := s.Repository.CreateCategory(ctx, u)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) CreateSubCategory(c context.Context, req SubCategory) (*Category, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u := SubCategory{
		Name:   req.Name,
		CategoryID: req.CategoryID,
	}
	r, err := s.Repository.CreateSubCategory(ctx, u)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) GetCategories(c context.Context) ([]*Category, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) GetCategory(c context.Context, id uint32) (*Category, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) GetProduct(c context.Context, id uint32) (*Product, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) GetProducts(c context.Context) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// func (s *service) CreateProduct(c context.Context, req *Product) (*Product, error) {
// 	ctx, cancel := context.WithTimeout(c, s.timeout)
// 	defer cancel()

// 	u := &Product{
// 	Name: req.Name,
// 	Slug: req.Slug,
// 	Description: req.Description,
// 	Quantity: req.Quantity,
// 	Campus: req.Campus,
// 	Status: req.Status,
// 	Image: req.Image,
// 	Store: req.Store,
// 	Condition: req.Condition,
// 	Price: req.Price,
// 	Category: req.Category,
// 	SubCategory: req.SubCategory,
// 	Variant: req.Variant,
// 	}
// 	r, err := s.Repository.CreateProduct(ctx, u)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return r, nil
// }

// func (s *service) UpdateProduct(c context.Context, req *Product) (*Product, error) {
//     ctx, cancel := context.WithTimeout(c, s.timeout)
//     defer cancel()

//     // First, check if the product exists by its ID or another unique identifier
//     existingProduct, err := s.Repository.GetProduct(ctx, req.ID)
//     if err != nil {
//         return nil, err
//     }

//     // Update only the fields that are present in the req
//     if req.Name != "" {
//         existingProduct.Name = req.Name
//     }
//     if req.Slug != "" {
//         existingProduct.Slug = req.Slug
//     }
//     if req.Description != "" {
//         existingProduct.Description = req.Description
//     }
//     if req.Quantity != 0 {
//         existingProduct.Quantity = req.Quantity
//     }
//     if req.Campus != "" {
//         existingProduct.Campus = req.Campus
//     }
//     if req.Status != "" {
//         existingProduct.Status = req.Status
//     }
//     if req.Image != "" {
//         existingProduct.Image = req.Image
//     }
//     if req.Condition != "" {
//         existingProduct.Condition = req.Condition
//     }
//     if req.Price != "" {
//         existingProduct.Price = req.Price
//     }

//     if req.Variant != "" {
//         existingProduct.Variant = req.Variant
//     }

//     // Update the product in the repository
//     updatedProduct, err := s.Repository.UpdateProduct(ctx, existingProduct)
//     if err != nil {
//         return nil, err
//     }

//     return updatedProduct, nil
// }

// func (s *service) GetProducts(c context.Context) ([]*Product, error) {
// 	ctx, cancel := context.WithTimeout(c, s.timeout)
// 	defer cancel()

// 	r, err := s.Repository.GetSubCategory(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return r, nil
// }

