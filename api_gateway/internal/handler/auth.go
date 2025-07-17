package handler

import (
	"api_gateway/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *authHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} domain.LoginResponse
// @Failure 500
// @Failure 400
// @Router /api/auth/login [get]
func (ah *authHandler) Login(c *gin.Context) {
	var loginReq domain.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	ctx := c.Request.Context()
	loginRes, err := ah.authUsecase.Login(ctx, &loginReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": loginRes.AccessToken, "user": loginRes.User})
}

// Register godoc
// @Summary Register
// @Description Register
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "User email"
// @Failure 500
// @Failure 400
// @Router /api/auth/register [post]
func (ah *authHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var regiserReq domain.RegisterRequest
	if err := c.ShouldBindJSON(&regiserReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	newUser, err := ah.authUsecase.Register(ctx, &regiserReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "register successfully",
		"email":   newUser.Email})
}

// RegisterAdmin godoc
// @Summary Register admin account
// @Description Register admin account
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "User email"
// @Failure 500
// @Failure 400
// @Router /api/auth/register [post]
func (ah *authHandler) RegisterAdmin(c *gin.Context) {
	ctx := c.Request.Context()
	var regiserReq domain.RegisterRequest
	if err := c.ShouldBindJSON(&regiserReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	newUser, err := ah.authUsecase.RegisterAdmin(ctx, &regiserReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "register successfully admin account",
		"email":   newUser.Email})
}
