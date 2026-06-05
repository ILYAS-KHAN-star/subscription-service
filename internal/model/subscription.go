package model

import "time"

type Subscription struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=1"`
	UserID      string  `json:"user_id" binding:"required,uuid"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type TotalCostRequest struct {
	UserID      string `form:"user_id" binding:"required,uuid"`
	ServiceName string `form:"service_name"`
	PeriodFrom  string `form:"period_from" binding:"required"`
	PeriodTo    string `form:"period_to" binding:"required"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Total         int64          `json:"total"`
	Page          int            `json:"page"`
	Limit         int            `json:"limit"`
	TotalPages    int            `json:"total_pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
