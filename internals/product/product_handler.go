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

func (h *Handler) UpdateProduct(ctx context.Context, input *Product) (*Product, error) {
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

func (h *Handler) GetProduct(ctx context.Context, id uint32) (*Product, error) {
	item, err := h.Service.GetProduct(ctx, id)
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

func (h *Handler) GetProducts(ctx context.Context, store string,limit int, offset int) ([]*Product, error) {
	item, err := h.Service.GetProducts(ctx, store,limit,offset)
	if err != nil {
		return nil, err
	}
	return item, nil
}
func (h *Handler) GetCategories(ctx context.Context) ([]*Category, error) {
	item, err := h.Service.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) GetWishListedProducts(ctx context.Context, user uint32) ([]*WishListedProduct, error) {
	item, err := h.Service.GetWishListedProducts(ctx, user)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *Handler) RemoveWishListedProduct(ctx context.Context, user uint32) error {
	err := h.Service.RemoveWishListedProduct(ctx, user)
	return err
}


func (h *Handler) AddWishListedProduct(ctx context.Context, user, product uint32) (*WishListedProduct, error) {
	item, err := h.Service.AddWishListedProduct(ctx, user, product)
	if err != nil {
		return nil, err
	}
	return item, nil
}
