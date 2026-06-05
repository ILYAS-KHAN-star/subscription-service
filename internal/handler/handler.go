package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"subscription-service/internal/model"
	"subscription-service/internal/service"
)

type Handler struct {
	service *service.Service
	logger  *zap.Logger
}

func NewHandler(service *service.Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// CreateSubscription godoc
// @Summary      Create subscription
// @Description  Создает новую подписку пользователя
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request  body      model.CreateSubscriptionRequest  true  "Subscription data"
// @Success      201      {object}  model.Subscription
// @Failure      400      {object}  model.ErrorResponse
// @Failure      500      {object}  model.ErrorResponse
// @Router       /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	sub, err := h.service.CreateSubscription(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// GetSubscription godoc
// @Summary      Get subscription by ID
// @Description  Возвращает подписку по её ID
// @Tags         subscriptions
// @Produce      json
// @Param        id    path      int  true  "Subscription ID"
// @Success      200   {object}  model.Subscription
// @Failure      400   {object}  model.ErrorResponse
// @Failure      404   {object}  model.ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	sub, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// UpdateSubscription godoc
// @Summary      Update subscription
// @Description  Обновляет данные существующей подписки
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id       path      int                             true  "Subscription ID"
// @Param        request  body      model.UpdateSubscriptionRequest  true  "Update data"
// @Success      200      {object}  model.Subscription
// @Failure      400      {object}  model.ErrorResponse
// @Failure      404      {object}  model.ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	var req model.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	sub, err := h.service.UpdateSubscription(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary      Delete subscription
// @Description  Удаляет подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id  path  int  true  "Subscription ID"
// @Success      204  "No Content"
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.service.DeleteSubscription(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSubscriptions godoc
// @Summary      List subscriptions
// @Description  Возвращает список подписок с возможностью фильтрации
// @Tags         subscriptions
// @Produce      json
// @Param        user_id     query  string  false  "Filter by user ID"
// @Param        service_name query  string  false  "Filter by service name"
// @Param        page         query  int     false  "Page number"          default(1)
// @Param        limit        query  int     false  "Items per page"       default(10)
// @Success      200         {object}  model.ListSubscriptionsResponse
// @Failure      500         {object}  model.ErrorResponse
// @Router       /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	subs, total, currentPage, err := h.service.ListSubscriptions(c.Request.Context(), userID, serviceName, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, model.ListSubscriptionsResponse{
		Subscriptions: subs,
		Total:         total,
		Page:          currentPage,
		Limit:         limit,
		TotalPages:    totalPages,
	})
}

// GetTotalCost godoc
// @Summary      Calculate total cost
// @Description  Calculate total cost of active subscriptions for a user within a date range
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query  string  true   "User ID (UUID)"
// @Param        service_name  query  string  false  "Filter by service name"
// @Param        period_from   query  string  true   "Start date (MM-YYYY)"
// @Param        period_to     query  string  true   "End date (MM-YYYY)"
// @Success      200          {object}  model.TotalCostResponse
// @Failure      400          {object}  model.ErrorResponse
// @Failure      500          {object}  model.ErrorResponse
// @Router       /total-cost [get]
func (h *Handler) GetTotalCost(c *gin.Context) {
	var req model.TotalCostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	total, err := h.service.GetTotalCost(c.Request.Context(), req.UserID, req.ServiceName, req.PeriodFrom, req.PeriodTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.TotalCostResponse{TotalCost: total})
}
