package handler

import (
	"customer_service/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	customerUsecase domain.CustomerUsecase
}

func NewCustomerHandler(customerUsecase domain.CustomerUsecase) *customerHandler {
	return &customerHandler{
		customerUsecase: customerUsecase,
	}
}

func (ch *customerHandler) GetAll(c *gin.Context) {
	customers, err := ch.customerUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No customers found"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (ch *customerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	customer, err := ch.customerUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer",
			"details": err.Error()})
		return
	}

	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (ch *customerHandler) Create(c *gin.Context) {
	var customerReq domain.CustomerRequest
	if err := c.ShouldBindJSON(&customerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	customer, err := ch.customerUsecase.Create(c.Request.Context(), &customerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create customer",
			"details": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Customer created successfully",
		"customer": customer})
}

func (ch *customerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var customerReq domain.CustomerRequest
	if err := c.ShouldBindJSON(&customerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if customerReq.Name == "" && customerReq.Email == "" && customerReq.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
		return
	}

	customer, err := ch.customerUsecase.Update(c.Request.Context(), id, &customerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update customer",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Customer updated successfully",
		"customer": customer})
}

func (ch *customerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	customer, err := ch.customerUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete customer",
			"details": err.Error()})
		return
	}

	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Customer deleted successfully",
		"customer": customer})
}
