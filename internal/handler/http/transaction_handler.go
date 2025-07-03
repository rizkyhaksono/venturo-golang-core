package http

import (
	"errors"
	"venturo-core/internal/model"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: s}
}

// generateInvoiceCode generates a unique invoice code.
type CreateTransactionPayload struct {
	Items []struct {
		ProductID   uuid.UUID `json:"product_id" validate:"required"`
		ProductName string    `json:"product_name" validate:"required"`
		Category    uint8     `json:"category" validate:"required,min=1,max=3"`
		Qty         int8      `json:"qty" validate:"required,min=1"`
		Price       int32     `json:"price" validate:"required,min=0"`
	} `json:"items" validate:"required,min=1"`
	Note string `json:"note"`
}

// CreateTransaction is the handler for creating a new transaction.
// @Summary      Create a new transaction
// @Description  Creates a transaction with multiple detail items for the authenticated user.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  body      CreateTransactionPayload  true  "Transaction Payload"
// @Success      201      {object}  response.ApiResponse{data=model.Transaction} "Successfully created transaction"
// @Failure      400      {object}  response.ApiResponse "Bad Request"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(uuid.UUID)

	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, errors.New("unauthorized"))
	}

	payload := new(CreateTransactionPayload)
	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	// Map payload to service input
	serviceInput := service.CreateTransactionInput{
		UserID: userID,
		Note:   payload.Note,
	}

	for _, item := range payload.Items {
		serviceInput.Items = append(serviceInput.Items, struct {
			ProductID   uuid.UUID
			ProductName string
			Category    model.ProductCategory
			Qty         int8
			Price       int32
		}{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Category:    model.ProductCategory(item.Category),
			Qty:         item.Qty,
			Price:       item.Price,
		})
	}

	transaction, err := h.transactionService.CreateTransaction(c.Context(), serviceInput)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, transaction)
}
