package user

import (
	"context"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}


func (h *Handler) CreateUser(ctx context.Context, input *CreateUserReq) (*CreateUserRes, error) {
	user, err := h.Service.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (h *Handler) VerifyOTP(ctx context.Context, input *User) (*User, error) {
	user, err := h.Service.VerifyOTP(ctx, input)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (h *Handler) Login(ctx context.Context, input *LoginUserReq) (*LoginUserRes, error) {
	res, err := h.Service.Login(ctx, input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *Handler) GetUsers(ctx context.Context) ([]*User, error) {
	user, err := h.Service.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (h *Handler) ToggleStoreFollowStatus(ctx context.Context,userId, storeId uint32) error {
	 err := h.Service.ToggleStoreFollowStatus(ctx,userId, storeId )
	if err != nil {
		return err
	}
	return  nil
}
