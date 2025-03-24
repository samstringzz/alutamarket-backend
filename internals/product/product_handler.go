package product

import "context"

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateProduct(ctx context.Context, input *NewProduct) (*Product, error) {
	item, err := h.Service.CreateProduct(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) UpdateProduct(ctx context.Context, input *NewProduct) (*Product, error) {
	item, err := h.Service.UpdateProduct(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) CreateCategory(ctx context.Context, input *Category) (*Category, error) {
	item, err := h.Service.CreateCategory(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) CreateSubCategory(ctx context.Context, input SubCategory) (*Category, error) {
	item, err := h.Service.CreateSubCategory(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetProduct(ctx context.Context, productId, userId uint32) (*Product, error) {
	item, err := h.Service.GetProduct(ctx, productId, userId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetCategory(ctx context.Context, id uint32) (*Category, error) {
	item, err := h.Service.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetProducts(ctx context.Context, store string, categorySlug string, limit int, offset int) ([]*Product, int, error) {
    return h.Service.GetProducts(ctx, store, categorySlug, limit, offset)
}
func (h *Handler) GetCategories(ctx context.Context) ([]*Category, error) {
	item, err := h.Service.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetHandledProducts(ctx context.Context, user uint32, eventType string) ([]*HandledProduct, error) {
	item, err := h.Service.GetHandledProducts(ctx, user, eventType)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) RemoveHandledProduct(ctx context.Context, user uint32, eventType string) error {
	err := h.Service.RemoveHandledProduct(ctx, user, eventType)
	return err
}

func (h *Handler) AddHandledProduct(ctx context.Context, user, product uint32, eventType string) (*HandledProduct, error) {
	item, err := h.Service.AddHandledProduct(ctx, user, product, eventType)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetRecommendedProducts(ctx context.Context, query string) ([]*Product, error) {
	item, err := h.Service.GetRecommendedProducts(ctx, query)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Add this method if it's not already present
func (h *Handler) SearchProducts(ctx context.Context, query string) ([]*Product, error) {
	products, err := h.Service.SearchProducts(ctx, query)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (h *Handler) GetProductReviews(ctx context.Context, productId uint32, sellerId string) ([]*Review, error) {
	item, err := h.Service.GetProductReviews(ctx, productId, sellerId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// func (h *Handler) GetRecentlyViewedProducts(ctx context.Context, user uint32)([]*Product,error){
// 	item, err := h.Service.GetRecentlyViewedProducts(ctx, user)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return item, nil
// }

func (h *Handler) AddReview(ctx context.Context, input *Review) (*Review, error) {
	item, err := h.Service.AddReview(ctx, input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// func (h *Handler) AddRecentlyViewedProducts(ctx context.Context, userId,productId uint3,eventType string2)error{
// 	err := h.Service.AddHandledProduct(ctx, userId,productId,eventType)
// 	if err != nil {
// 		return  err
// 	}
// 	return nil
// }
