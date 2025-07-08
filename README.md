# E-Commerce Microservices Platform

A comprehensive e-commerce platform built with Go using microservices architecture. This project demonstrates modern backend development practices including API Gateway, service isolation, caching, and comprehensive API documentation.

## 🏗️ Architecture

This project follows a microservices architecture pattern with the following components:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Client/Web    │    │   API Gateway    │    │  Swagger UI     │
│   Application   │◄──►│    (Port 8080)   │◄──►│  Documentation  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                    ┌───────────┼───────────┐
                    │           │           │
            ┌───────▼────┐ ┌────▼─────┐ ┌───▼──────┐
            │ Product    │ │ Customer │ │ Cart     │
            │ Service    │ │ Service  │ │ Service  │
            │(Port 8082) │ │(Port 8081)│ │(Port 8083)│
            └────────────┘ └──────────┘ └──────────┘
                    │           │           │
                    └───────────┼───────────┘
                                │
                    ┌───────────▼───────────┐
                    │     MongoDB + Redis   │
                    │   (Database & Cache)  │
                    └───────────────────────┘
```

## 🚀 Features

### Core Services
- **API Gateway**: Central entry point for all client requests with routing, CORS, and request/response logging
- **Product Service**: Manage product catalog with caching support
- **Customer Service**: Handle customer registration, authentication, and profile management  
- **Cart Service**: Shopping cart functionality with Redis caching
- **Order Service** : On developing 
- **Payment Service** : On developing
- **Shipment Service** : On developing
- **Notification Service** : On developing

### Technical Features
- **Microservices Architecture**: Loosely coupled services with clear separation of concerns
- **RESTful APIs**: Well-designed REST endpoints following HTTP standards
- **API Documentation**: Comprehensive Swagger/OpenAPI documentation with examples
- **Caching Layer**: Redis integration for improved performance
- **Database**: MongoDB for data persistence
- **Logging**: Structured logging with Logrus
- **CORS Support**: Cross-origin resource sharing configuration

## 📋 Prerequisites

- Go 1.21 or higher
- MongoDB 6.0+
- Redis 7.0+

## 🛠️ Installation & Setup

### Option 1: Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd intern-project-v3
   ```

2. **Set up environment variables**
   
   Create `.env` files for each service or use the default values:
   
   ```bash
   # API Gateway
   PORT=8080
   PRODUCT_SERVICE_HOST=http://localhost:8082
   CUSTOMER_SERVICE_HOST=http://localhost:8081
   CART_SERVICE_HOST=http://localhost:8083
   MONGODB_URI=mongodb://localhost:27017
   DB_NAME=auth_db
   
   # Product Service
   PORT=8082
   MONGODB_URI=mongodb://localhost:27017
   DB_NAME=product_db
   REDIS_ADDR=localhost:6379
   
   # Customer Service  
   PORT=8081
   MONGODB_URI=mongodb://localhost:27017
   DB_NAME=customer_db
   
   # Cart Service
   PORT=8083
   MONGODB_URI=mongodb://localhost:27017
   DB_NAME=cart_db
   REDIS_ADDR=localhost:6379
   ```

3. **Install dependencies**
   ```bash
   # For each service
   cd api_gateway && go mod tidy
   cd ../product_service && go mod tidy
   cd ../customer_service && go mod tidy
   cd ../cart_service && go mod tidy
   ```

4. **Start MongoDB and Redis**
   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:6
   docker run -d -p 6379:6379 --name redis redis:7
   ```

5. **Run services**
   ```bash
   # Terminal 1 - API Gateway
   cd api_gateway && go run cmd/main.go
   
   # Terminal 2 - Product Service
   cd product_service && go run cmd/main.go
   
   # Terminal 3 - Customer Service
   cd customer_service && go run cmd/main.go
   
   # Terminal 4 - Cart Service
   cd cart_service && go run cmd/main.go
   ```

## 📖 API Documentation

### Quick Access
- **Swagger UI**: http://localhost:8080/swagger-ui/ (Interactive API documentation)
- **Health Check**: http://localhost:8080/health

### Authentication
All protected endpoints require JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### API Endpoints Reference

#### 🔐 Authentication & User Management

##### Register Customer
```http
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "0123456789",
  "password": "securepassword123"
}
```
**Response:** Returns customer data and authentication token.

##### Customer Login
```http
GET /api/users
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```
**Response:** Returns JWT token and user information.

#### 📦 Product Management

##### Get All Products
```http
GET /api/product_service/products
```
**Query Parameters:**
- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `category` (optional): Filter by product category
- `search` (optional): Search products by name or description

**Response:** Array of product objects with pagination metadata.

##### Get Product by ID
```http
GET /api/product_service/products/{id}
```
**Path Parameters:**
- `id` (required): MongoDB ObjectID of the product

**Response:** Single product object or 404 if not found.

##### Create Product
```http
POST /api/product_service/products
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "iPhone 15 Pro",
  "description": "Latest Apple smartphone with advanced features",
  "price": 999.99,
  "category": "Electronics",
  "stock": 100,
  "images": ["url1", "url2"],
  "specifications": {
    "brand": "Apple",
    "model": "iPhone 15 Pro",
    "color": "Space Black"
  }
}
```
**Authorization:** Admin role required.

##### Update Product
```http
PUT /api/product_service/products/{id}
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Updated Product Name",
  "price": 899.99,
  "stock": 150
}
```
**Note:** Partial updates supported - only include fields you want to change.

##### Delete Product
```http
DELETE /api/product_service/products/{id}
Authorization: Bearer <admin-token>
```
**Authorization:** Admin role required.

#### 👥 Customer Management

##### Get All Customers
```http
GET /api/customer_service/customers
Authorization: Bearer <admin-token>
```
**Authorization:** Admin role required.
**Query Parameters:**
- `page`, `limit`: Pagination options
- `search`: Search by name or email

##### Get Customer by ID
```http
GET /api/customer_service/customers/{id}
Authorization: Bearer <token>
```
**Authorization:** Users can only access their own data unless admin.

##### Update Customer Profile
```http
PUT /api/customer_service/customers/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Name",
  "phone": "0987654321"
}
```
**Note:** Email cannot be changed after registration.

#### 🛒 Cart Management

##### Get User's Cart
```http
GET /api/cart_service/cart
Authorization: Bearer <token>
```
**Response:** Returns cart with items, quantities, and total price.

##### Add Item to Cart
```http
POST /api/cart_service/cart/items
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": "507f1f77bcf86cd799439011",
  "quantity": 2
}
```
**Business Rules:**
- If item already exists, quantity will be updated
- Stock availability is checked before adding
- Maximum quantity per item: 10

##### Update Cart Item Quantity
```http
PUT /api/cart_service/cart/items/{product_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "quantity": 5
}
```
**Note:** Set quantity to 0 to remove item from cart.

##### Remove Item from Cart
```http
DELETE /api/cart_service/cart/items/{product_id}
Authorization: Bearer <token>
```

##### Clear Entire Cart
```http
DELETE /api/cart_service/cart
Authorization: Bearer <token>
```

### Response Format

#### Success Response
```json
{
  "status": "success",
  "data": {
    // Response data here
  },
  "message": "Operation completed successfully"
}
```

#### Error Response
```json
{
  "status": "error",
  "error": "Error description",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

### HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

### Data Types

#### Product Object
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "category": "Category Name",
  "stock": 50,
  "images": ["url1", "url2"],
  "specifications": {},
  "created_at": "2025-07-08T10:00:00Z",
  "updated_at": "2025-07-08T10:00:00Z"
}
```

#### Customer Object
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "0123456789",
  "created_at": "2025-07-08T10:00:00Z",
  "updated_at": "2025-07-08T10:00:00Z"
}
```

#### Cart Object
```json
{
  "id": "507f1f77bcf86cd799439011",
  "customer_id": "507f1f77bcf86cd799439011",
  "items": [
    {
      "product_id": "507f1f77bcf86cd799439011",
      "product_name": "Product Name",
      "price": 99.99,
      "quantity": 2,
      "subtotal": 199.98
    }
  ],
  "total_items": 2,
  "total_price": 199.98,
  "updated_at": "2025-07-08T10:00:00Z"
}
```

## 🧪 Testing

### Using Swagger UI
1. Navigate to http://localhost:8080/swagger-ui/
2. Select the service you want to test from the dropdown
3. Try the API endpoints with the provided examples

### Using cURL

```bash
# Get all products
curl -X GET http://localhost:8080/api/product_service/products

# Create a new customer
curl -X POST http://localhost:8080/api/customer_service/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com", 
    "phone": "0123456789"
  }'

# Add item to cart
curl -X POST http://localhost:8080/api/cart_service/cart/items \
  -H "Content-Type: application/json" \
  -H "X-User-ID: customer_id_here" \
  -d '{
    "product_id": "product_id_here",
    "quantity": 2
  }'
```

## 🏗️ Project Structure

```
intern-project-v3/
├── api_gateway/
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── app/
│   │   ├── handler/
│   │   ├── repository/
│   │   └── usecase/
│   ├── middleware/
│   ├── config/
│   └── static/swagger-ui/
├── product_service/
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── app/
│   │   ├── domain/
│   │   ├── handler/
│   │   ├── repository/
│   │   └── usecase/
│   ├── config/
│   ├── docs/
│   └── proto/
├── customer_service/
│   └── (similar structure)
├── cart_service/
│   └── (similar structure)
├── docker-compose.yml
└── README.md
```

## 🔧 Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: MongoDB
- **Cache**: Redis
- **Documentation**: Swagger/OpenAPI 3.0
- **Logging**: Logrus
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker
- **Communication**: HTTP/REST, gRPC (for internal service communication)


## 🔄 Development Workflow

1. **Add new features** to individual services
2. **Update API documentation** using Swagger annotations
3. **Test APIs** using Swagger UI or cURL
4. **Regenerate docs** with `swag init -g cmd/main.go -o docs`
5. **Update README** if needed
