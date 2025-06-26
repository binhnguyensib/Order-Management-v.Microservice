package handler

import (
	"cart_service/internal/domain"

	"github.com/gin-gonic/gin"
)

type cartHandler struct {
	cartUsecase domain.CartUsecase
}

func NewCartHandler(cartUsecase domain.CartUsecase) *cartHandler {
	return &cartHandler{
		cartUsecase: cartUsecase,
	}
}

func (ch *cartHandler) AddToCart(c *gin.Context) {
	ctx := c.Request.Context()
	customerID := c.Param("id")
	var cartItem domain.CartItemRequest
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	cart, err := ch.cartUsecase.AddToCart(ctx, customerID, &cartItem)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cart)
}

func (ch *cartHandler) GetCartByCustomerId(c *gin.Context) {
	ctx := c.Request.Context()
	customerID := c.Param("id")
	cart, err := ch.cartUsecase.GetCartByCustomerId(ctx, customerID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cart)
}

func (ch *cartHandler) UpdateCartItem(c *gin.Context) {
	ctx := c.Request.Context()
	customerID := c.Param("id")
	var cartItem domain.CartItemRequest
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	cart, err := ch.cartUsecase.UpdateCartItem(ctx, customerID, &cartItem)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cart)
}

func (ch *cartHandler) RemoveCartItem(c *gin.Context) {
	ctx := c.Request.Context()
	customerID := c.Param("id")
	productID := c.Param("product_id")
	cart, err := ch.cartUsecase.RemoveCartItem(ctx, customerID, productID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cart)
}

func (ch *cartHandler) ClearCart(c *gin.Context) {
	ctx := c.Request.Context()
	customerID := c.Param("id")
	err := ch.cartUsecase.ClearCart(ctx, customerID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Cart cleared successfully"})
}
