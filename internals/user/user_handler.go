package user

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)
type Handler struct{
	Service
}

func NewHandler(s Service) *Handler{
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateUser (c *gin.Context){
var u CreateUserReq
if err:= c.ShouldBindJSON(&u); err!=nil{
	c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
	return
}
res,err:=h.Service.CreateUser(c.Request.Context(),&u)
if err !=nil{
	c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
	return
}
c.JSON(http.StatusOK,res)
}

func (h *Handler) Login (c *gin.Context){
	var user LoginUserReq
	if err:= c.ShouldBindJSON(&user); err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	u,err:= h.Service.Login(c.Request.Context(),&user)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.SetCookie("cookie-session", u.accessToken, 21600, "/auth", os.Getenv("DOMAIN"), false, true)

	res:= &LoginUserRes{
		ID: u.ID,
	}
	c.JSON(http.StatusOK,res)
}
func (h *Handler) Logout(c *gin.Context){
	c.SetCookie("cookie-session","",-1,"","",false,true)
	c.JSON(http.StatusOK,gin.H{"message":"Logout Successfully"})
}