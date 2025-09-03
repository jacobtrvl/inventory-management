// Copyright 2025 Jacob Philip. All rights reserved.
package inventory

import (
	"time"

	"github.com/jacobtrvl/inventory-management/internal/store"
	"github.com/jacobtrvl/inventory-management/pkg/observability"
)

type Inventory struct {
	tableName string
	db        store.Store
	mc        *observability.MetricsCollector
}

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Pagination and filtering parameters for listing products.
// Currently only pagination is supported
type ListParams struct {
	Page  *int
	Limit *int
}

type ListMetadata struct {
	CurrentPage int  `json:"current_page"`
	NextPage    *int `json:"next_page,omitempty"`
}

type CreateRequest struct {
	ID    string  `json:"id,omitempty"`
	Name  string  `json:"name"`
	Price float64 `json:"price,omitempty"`
	Stock int     `json:"stock,omitempty"`
}

type UpdateRequest struct {
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
	Stock *int     `json:"stock,omitempty"`
}
