package handler

import (
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

// GetAll godoc
// @Summary Get all products
// @Description Retrieve all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {array} domain.Product
// @Failure 500
// @Failure 400
// @Router /api/product_service/products [get]
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

// GetByID godoc
// @Summary Get product by ID
// @Description Retrieve a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400
// @Failure 500
// @Router /api/product_service/products/{id} [get]
func (ph *productHandler) GetByID(c *gin.Context) {
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

// Create godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags Products
// @Accept json
// @Produce json
// @Param product body domain.ProductRequest true "Product Request"
// @Success 201 {object} domain.Product
// @Failure 400
// @Failure 500
// @Router /api/product_service/products [post]
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

// Update godoc
// @Summary Update an existing product
// @Description Update a product with the provided details
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body domain.ProductRequest true "Product Request"
// @Success 200 {object} domain.Product
// @Failure 400
// @Failure 500
// @Router /api/product_service/products/{id} [put]
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

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/product_service/products/{id} [delete]
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
