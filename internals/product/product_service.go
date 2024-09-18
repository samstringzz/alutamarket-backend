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

func (s *service) GetProduct(c context.Context, productId, userId uint32) (*Product, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.GetProduct(ctx, productId, userId)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) GetProducts(c context.Context, store string, limit int, offset int) ([]*Product, int, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, count, err := s.Repository.GetProducts(ctx, store, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return r, count, nil
}

func (s *service) UpdateProduct(c context.Context, req *NewProduct) (*Product, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	// Update the product in the repository
	updatedProduct, err := s.Repository.UpdateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *service) AddHandledProduct(ctx context.Context, userId, productId uint32, eventType string) (*HandledProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	addedWishlists, err := s.Repository.AddHandledProduct(ctx, userId, productId, eventType)
	if err != nil {
		return nil, err
	}
	return addedWishlists, nil
}

func (s *service) GetHandledProducts(ctx context.Context, userId uint32, eventType string) ([]*HandledProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	getWishlists, err := s.Repository.GetHandledProducts(ctx, userId, eventType)
	if err != nil {
		return nil, err
	}
	return getWishlists, nil
}

func (s *service) RemoveHandledProduct(ctx context.Context, userId uint32, eventType string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.Repository.RemoveHandledProduct(ctx, userId, eventType)
	return err
}
func (s *service) DeleteProduct(ctx context.Context, id uint32) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.Repository.DeleteProduct(ctx, id)
	return err
}

func (s *service) GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	prd, err := s.Repository.GetRecommendedProducts(ctx, query)
	if err != nil {
		return nil, err
	}
	return prd, nil
}
func (s *service) SearchProducts(ctx context.Context, query string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	prd, err := s.Repository.SearchProducts(ctx, query)
	if err != nil {
		return nil, err
	}
	return prd, nil
}

// func (s *service) AddRecentlyViewedProducts(ctx context.Context, userId,productId uint32)error{
//   ctx, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()
// 	err := s.Repository.AddRecentlyViewedProducts(ctx, userId,productId)
// 	if err !=nil{
// 		return err
// 	}
// 	return nil
// }

// func (s *service) GetRecentlyViewedProducts(ctx context.Context, userId uint32)([]*Product,error){
//   ctx, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()
// 	prd,err := s.Repository.GetRecentlyViewedProducts(ctx, userId)
// 	if err !=nil{
// 		return nil,err
// 	}
// 	return prd,nil
// }

func (s *service) AddReview(ctx context.Context, input *Review) (*Review, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	review, err := s.Repository.AddReview(ctx, input)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (s *service) GetProductReviews(ctx context.Context, productId uint32, sellerId string) ([]*Review, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	reviews, err := s.Repository.GetProductReviews(ctx, productId, sellerId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
