package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	var user model.UserLogin

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("email or password is empty"))
		return
	}

	var recordUser = model.User{
		Email:    user.Email,
		Password: user.Password,
	}

	token, err := u.userService.Login(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    *token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	// Cek apakah cookie session_token sudah ada
	if _, err := c.Request.Cookie("session_token"); err == nil {
		// Jika cookie sudah ada, ganti nilai cookie dengan tokenString yang baru
		http.SetCookie(c.Writer, cookie)
	} else {
		// Jika cookie belum ada, tambahkan cookie baru
		http.SetCookie(c.Writer, cookie)
	}

	// Mendapatkan user_id dari token JWT
	var claims model.Claims
	_, err = jwt.ParseWithClaims(*token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(model.JwtKey), nil
	})

	response := gin.H{
		"user_id": claims.Email,
		"message": "login success",
	}

	c.JSON(http.StatusOK, response)

	// TODO: answer here
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	userTaskCategory, err := u.userService.GetUserTaskCategory()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, userTaskCategory)
	// TODO: answer here
}
