package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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
// func (h *Handler) Login(c *gin.Context) {
// 	var user LoginUserReq
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	u, err := h.Service.Login(c.Request.Context(), &user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.SetCookie("cookie-session", u.accessToken, 21600, "/auth", os.Getenv("DOMAIN"), false, true)

// 	res := &LoginUserRes{
// 		ID: u.ID,
// 	}
// 	c.JSON(http.StatusOK, res)
// }

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("cookie-session", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout Successfully"})
}
