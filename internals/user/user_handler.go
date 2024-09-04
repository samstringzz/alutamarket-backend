package user

import (
	"context"

	"github.com/Chrisentech/aluta-market-api/internals/store"
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

func (h *Handler) VerifyOTP(ctx context.Context, input *VerifyOTPReq) (*User, error) {
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

func (h *Handler) ToggleStoreFollowStatus(ctx context.Context, userId, storeId uint32) error {
	err := h.Service.ToggleStoreFollowStatus(ctx, userId, storeId)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) UpdateUser(ctx context.Context, user *User) (*User, error) {
	usr, err := h.Service.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
func (h *Handler) CreateStore(ctx context.Context, user *store.Store) (*store.Store, error) {
	usr, err := h.Service.CreateStore(ctx, user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (h *Handler) CreateDVAAccount(ctx context.Context, req *DVADetails) (string, error) {
	resp, err := h.Service.CreateDVAAccount(ctx, req)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (h *Handler) GetMyDVA(ctx context.Context, email string) (*Account, error) {
	resp, err := h.Service.GetMyDVA(ctx, email)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
