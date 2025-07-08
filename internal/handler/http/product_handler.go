package http

import (
	"errors"
	"strconv"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

// CreateProduct handles the multipart/form-data request to create a product.
// @Summary      Create a new product
// @Description  Creates a new product with the provided details.
// @Tags         Products
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        name   formData  string  true  "Product Name"
// @Param        price  formData  int     true  "Product Price"
// @Param        stock  formData  int     true  "Product Stock"
// @Param        image  formData  file    false "Product Image"
// @Success      201    {object}  response.ApiResponse{data=model.Product} "Successfully created product"
// @Failure      400    {object}  response.ApiResponse "Bad Request"
// @Failure      401    {object}  response.ApiResponse "Unauthorized"
// @Failure      500    {object}  response.ApiResponse "Internal Server Error"
// @Router       /products [post]
// CreateProduct godoc
// CreateProduct handles the creation of a new product.
// It expects a multipart/form-data request with fields for name, price, stock, and an optional image file.
// It validates the price and stock fields, and returns an error if they are not in the correct format.
// If successful, it returns the created product with a 201 status code.
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid price format"))
	}
	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("invalid stock format"))
	}

	input := service.CreateProductInput{
		Name:  c.FormValue("name"),
		Price: int32(price),
		Stock: int16(stock),
	}

	file, err := c.FormFile("image")
	if err == nil {
		input.Image = file
	}

	product, err := h.productService.CreateProduct(c.Context(), input)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, product)
}
