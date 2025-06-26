package handler

import (
	"fmt"
	"net/http"
	"product_service/internal/domain"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(productUsecase domain.ProductUsecase) *productHandler {
	return &productHandler{
		productUsecase: productUsecase,
	}
}

func (ph *productHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	products, err := ph.productUsecase.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No products found"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ph *productHandler) GetByID(c *gin.Context) {
	fmt.Println("Handler is called")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	ctx := c.Request.Context()
	product, err := ph.productUsecase.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product", "details": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ph *productHandler) Create(c *gin.Context) {
	var productReq domain.ProductRequest
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx := c.Request.Context()
	product, err := ph.productUsecase.Create(ctx, &productReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (ph *productHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var productReq domain.ProductRequest
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if productReq.Name == "" && productReq.Price <= 0 && productReq.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
		return
	}

	ctx := c.Request.Context()
	product, err := ph.productUsecase.Update(ctx, id, &productReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

func (ph *productHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	ctx := c.Request.Context()
	product, err := ph.productUsecase.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product", "details": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
		"product": product,
	})
}
