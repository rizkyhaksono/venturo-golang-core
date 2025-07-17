package http

import (
	"venturo-core/internal/service"
	"venturo-core/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new report handler.
func NewReportHandler(s *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: s}
}

// GetInventoryReport handles the GET /api/v1/reports/inventory request.
// @Summary      Get Inventory Report
// @Description  Generate a comprehensive inventory report with on-hand quantities and transaction history
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "Bearer JWT token"
// @Param        item_id query string false "Filter by specific item ID"
// @Param        outlet_id query string false "Filter by specific outlet ID"
// @Success      200      {object}  response.ApiResponse{data=[]service.InventoryReportItem} "Inventory report generated successfully"
// @Failure      400      {object}  response.ApiResponse "Bad Request"
// @Failure      401      {object}  response.ApiResponse "Unauthorized"
// @Failure      500      {object}  response.ApiResponse "Internal Server Error"
// @Router       /reports/inventory [get]
func (h *ReportHandler) GetInventoryReport(c *fiber.Ctx) error {
	// Parse query parameters
	var input service.InventoryReportInput

	// Parse item_id if provided
	if itemIDStr := c.Query("item_id"); itemIDStr != "" {
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err)
		}
		input.ItemID = &itemID
	}

	// Parse outlet_id if provided
	if outletIDStr := c.Query("outlet_id"); outletIDStr != "" {
		outletID, err := uuid.Parse(outletIDStr)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err)
		}
		input.OutletID = &outletID
	}

	// Generate the report
	report, err := h.reportService.GenerateInventoryReport(c.Context(), input)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.StatusOK, report)
}
