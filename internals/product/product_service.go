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

func (s *service) CreateProduct(c context.Context, req *NewProduct) (*Product, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) CreateSubCategory(c context.Context, req SubCategory) (*Category, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u := SubCategory{
		Name:       req.Name,
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

func (s *service) GetProducts(c context.Context, store string,limit int, offset int) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.GetProducts(ctx, store,limit,offset)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) UpdateProduct(c context.Context, req *Product) (*Product, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	// Update the product in the repository
	updatedProduct, err := s.Repository.UpdateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *service) AddWishListedProduct(ctx context.Context, userId, productId uint32) (*WishListedProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	addedWishlists, err := s.Repository.AddWishListedProduct(ctx, userId, productId)
	if err != nil {
		return nil, err
	}
	return addedWishlists, nil
}

func (s *service) GetWishListedProducts(ctx context.Context, userId uint32) ([]*WishListedProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	getWishlists, err := s.Repository.GetWishListedProducts(ctx, userId)
	if err != nil {
		return nil, err
	}
	return getWishlists, nil
}

func (s *service) RemoveWishListedProduct(ctx context.Context, userId uint32) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.Repository.RemoveWishListedProduct(ctx, userId)
	return err
}
func (s *service) DeleteProduct(ctx context.Context, id uint32) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.Repository.DeleteProduct(ctx, id)
	return err
}

func (s *service) GetRecommendedProducts(ctx context.Context, query string)([]*Product,error){
  ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	prd,err := s.Repository.GetRecommendedProducts(ctx, query)
	if err !=nil{
		return nil,err
	}
	return prd,nil
}
func (s *service) SearchProducts(ctx context.Context, query string)([]*Product,error){
  ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	prd,err := s.Repository.SearchProducts(ctx, query)
	if err !=nil{
		return nil,err
	}
	return prd,nil
}

