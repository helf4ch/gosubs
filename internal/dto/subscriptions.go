package dto

import "github.com/google/uuid"

type Subscription struct {
	Id          uuid.UUID `json:"id"`
	ServiceName string    `json:"serviceName"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"userId"`
	StartDate   string    `json:"startDate"`
	EndDate     string    `json:"endDate,omitempty"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt,omitempty"`
}

// Create
type CreateSubscriptionRequest struct {
	ServiceName string    `json:"serviceName"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"userId"`
	StartDate   string    `json:"startDate"`
	EndDate     string    `json:"endDate,omitempty"`
}

type CreateSubscriptionResponse struct {
	Subscription
}

// Read
type ReadSubscriptionResponse struct {
	Subscription
}

// Update
type UpdateSubscriptionRequest struct {
	ServiceName string    `json:"serviceName"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"userId"`
	StartDate   string    `json:"startDate"`
	EndDate     string    `json:"endDate,omitempty"`
}

type UpdateSubscriptionResponse struct {
	Subscription
}

// Delete
type DeleteSubscriptionResponse struct {
	Count int `json:"count"`
}

// List
type ListSubscriptionResponse struct {
	List []Subscription `json:"list"`
}

// Total
type TotalSubscriptionResponse struct {
	Cost int `json:"cost"`
}
