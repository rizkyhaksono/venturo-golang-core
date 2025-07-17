package http

import (
	"errors"
	"venturo-core/internal/service"
	"venturo-core/pkg/response"
	"venturo-core/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

// NewInventoryHandler creates a new inventory handler.
func NewInventoryHandler(s *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: s}
}

// StockInPayload defines the expected JSON payload for stock-in operation.
type StockInPayload struct {
	ItemID   uuid.UUID `json:"item_id" validate:"required"`
	OutletID uuid.UUID `json:"outlet_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required,min=1"`
}

// StockIn handles the POST /api/v1/inventory/stock-in request.
// @Summary      Stock In
// @Description  Add stock to inventory by creating a positive inventory ledger entry
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer JWT token"
// @Param        payload body StockInPayload true "Stock in data"
// @Success      201      {object}  response.ApiResponse{data=model.InventoryLedger} "Stock in successful"
// @Failure      400      {object}  response.ApiResponse "Bad Request"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Failure      500      {object}  response.ApiResponse "Internal Server Error"
// @Router       /inventory/stock-in [post]
func (h *InventoryHandler) StockIn(c *fiber.Ctx) error {
	// Parse and validate the request payload
	payload := new(StockInPayload)
	if err := c.BodyParser(payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, errors.New("cannot parse JSON"))
	}

	if errs := validator.ValidateStruct(payload); errs != nil {
		return response.ValidationError(c, errs)
	}

	// Map payload to service input
	serviceInput := service.StockInInput{
		ItemID:   payload.ItemID,
		OutletID: payload.OutletID,
		Quantity: payload.Quantity,
	}

	// Call the service
	ledger, err := h.inventoryService.StockIn(c.Context(), serviceInput)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusCreated, ledger)
}
