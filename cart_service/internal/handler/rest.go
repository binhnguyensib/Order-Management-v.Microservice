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

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a product to the customer's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param cartItem body domain.CartItemRequest true "Cart Item Request"
// @Success 200 {object} domain.Cart
// @Failure 400
// @Failure 500
// @Router /api/user/cart_service/carts/{id} [post]
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

// GetCartByCustomerId godoc
// @Summary Get cart by customer ID
// @Description Retrieve the cart for a specific customer
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} domain.Cart
// @Failure 500
// @Failure 400
// @Router /api/user/cart_service/carts/{id} [get]
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

// AddToCart godoc
// @Summary Update cart item
// @Description Update an existing item in customer's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param cartItem body domain.CartItemRequest true "Cart Item Request"
// @Success 200 {object} domain.Cart
// @Failure 400
// @Failure 500
// @Router /api/user/cart_service/carts/{id} [put]
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

// RemoveCartItem godoc
// @Summary Remove item from cart
// @Description Remove a product from the customer's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} domain.Cart
// @Failure 400
// @Failure 500
// @Router /api/user/cart_service/carts/{id}/item/{product_id} [delete]
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

// ClearCart godoc
// @Summary Clear cart
// @Description Clear all items from the customer's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200
// @Failure 500
// @Router /api/user/cart_service/carts/{id} [delete]
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
