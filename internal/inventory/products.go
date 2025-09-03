// Copyright 2025 Jacob Philip. All rights reserved.
// Package inventory provides inventory management functionalities.
// Currently a sample Product struct is defined.
// Error handling and status codes handling are not ideal.
// It is not a good practice to return HTTP status codes from internal packages.
// Currently we are doing it for simplicity. This needs improvement.
package inventory

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jacobtrvl/inventory-management/internal/store"
	"github.com/jacobtrvl/inventory-management/pkg/observability"
)

func NewInventory(ctx context.Context, table string, db *store.MemDb, mc *observability.MetricsCollector) *Inventory {
	db.CreateTable(table)
	return &Inventory{
		tableName: table,
		db:        db,
		mc:        mc,
	}
}

func (p *Inventory) Add(ctx context.Context, req CreateRequest) (string, int, error) {
	_, err := p.db.Read(p.tableName, req.ID)
	if err == nil {
		p.mc.RecordOperation(observability.OpInsert, false)
		return "", http.StatusConflict, fmt.Errorf("product with ID %s already exists", req.ID)
	}
	product := Product{
		ID:    req.ID,
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	if product.ID == "" {
		product.ID = generateID()
	}
	currentTime := time.Now()
	product.CreatedAt = currentTime
	product.UpdatedAt = currentTime
	// Validate ID length. In a real system, more validations would be needed.
	// Some validations can be done during request binding in the API layer
	if len(product.ID) > 255 {
		p.mc.RecordOperation(observability.OpInsert, false)
		return "", http.StatusBadRequest, fmt.Errorf("product ID too long: maximum 255 characters, got %d", len(product.ID))
	}
	err = p.db.Write(p.tableName, product.ID, product)
	if err != nil {
		p.mc.RecordOperation(observability.OpInsert, false)
		slog.ErrorContext(ctx, "Failed to add product", "id", product.ID, "error", err)
		return "", http.StatusInternalServerError, err
	}
	p.mc.RecordOperation(observability.OpInsert, true)
	slog.DebugContext(ctx, "Product added", "id", product.ID)
	return product.ID, http.StatusCreated, nil
}

func (p *Inventory) Get(ctx context.Context, id string) (Product, int, error) {
	item, err := p.db.Read(p.tableName, id)
	if err != nil {
		p.mc.RecordOperation(observability.OpGet, false)
		return Product{}, http.StatusNotFound, err
	}
	p.mc.RecordOperation(observability.OpGet, true)
	slog.DebugContext(ctx, "Product retrieved", "id", id)
	return item.(Product), http.StatusOK, nil
}

func (p *Inventory) Update(ctx context.Context, id string, req UpdateRequest) (int, error) {
	product, err := p.db.Read(p.tableName, id); 
	if err != nil {
		p.mc.RecordOperation(observability.OpUpdate, false)
		return http.StatusNotFound, err
	}
	pd := product.(Product)
	if req.Name != nil {
		pd.Name = *req.Name
	}
	if req.Price != nil {
		pd.Price = *req.Price
	}
	if req.Stock != nil {
		pd.Stock = *req.Stock
	}
	pd.UpdatedAt = time.Now()
	if err := p.db.Write(p.tableName, id, pd); err != nil {
		p.mc.RecordOperation(observability.OpUpdate, false)
		slog.ErrorContext(ctx, "Failed to update product", "id", id, "error", err)
		return http.StatusInternalServerError, err
	}

	p.mc.RecordOperation(observability.OpUpdate, true)
	slog.DebugContext(ctx, "Product updated", "id", id)
	return http.StatusOK, nil
}

func (p *Inventory) Delete(ctx context.Context, id string) (int, error) {
	_, err := p.db.Read(p.tableName, id)
	if err != nil {
		p.mc.RecordOperation(observability.OpDelete, false)
		return http.StatusNotFound, err
	}
	if err := p.db.Delete(p.tableName, id); err != nil {
		p.mc.RecordOperation(observability.OpDelete, false)
		slog.ErrorContext(ctx, "Failed to delete product", "id", id, "error", err)
		return http.StatusInternalServerError, err
	}
	p.mc.RecordOperation(observability.OpDelete, true)
	slog.DebugContext(ctx, "Product deleted", "id", id)
	return http.StatusOK, nil
}

// List returns a list of products based on the provided ListParams.
// Returns Products slice, EOF status, and error (if any).
func (p *Inventory) List(ctx context.Context, params ListParams) ([]Product, *ListMetadata, int, error) {
	// Filtering is not supported due to in-memory DB limitations.
	// Without an underlying DB with query capabilities filtering can be error-prone and inefficient.
	if params.Page == nil && params.Limit == nil {
		list, err := p.GetAllItems(ctx)
		if err != nil {
			p.mc.RecordOperation(observability.OpList, false)
			return nil, nil, http.StatusInternalServerError, err
		}
		p.mc.RecordOperation(observability.OpList, true)
		return list, nil, http.StatusOK, nil
	}
	if params.Page == nil || params.Limit == nil {
		p.mc.RecordOperation(observability.OpList, false)
		return nil, nil, http.StatusBadRequest, fmt.Errorf("both page and limit must be provided")
	}
	page := *params.Page
	limit := *params.Limit

	list, eof, err := p.NoFilter(ctx, (page-1)*limit, page*limit)
	if err != nil {
		p.mc.RecordOperation(observability.OpList, false)
		return nil, nil, http.StatusInternalServerError, err
	}
	var nextPage *int
	if !eof {
		np := page + 1
		nextPage = &np
	}
	p.mc.RecordOperation(observability.OpList, true)
	return list, &ListMetadata{
		CurrentPage: page,
		NextPage:    nextPage,
	}, http.StatusOK,nil
}

func (p *Inventory) NoFilter(ctx context.Context, start, end int) ([]Product, bool, error) {
	var unfiltered []Product
	items, eof, err := p.db.ReadRange(p.tableName, start, end)
	if err != nil {
		p.mc.RecordOperation(observability.OpList, false)
		slog.ErrorContext(ctx, "Failed to retrieve products", "error", err)
		return nil, false, err
	}
	for _, item := range items {
		unfiltered = append(unfiltered, item.(Product))
	}

	return unfiltered, eof, nil
}

func (p *Inventory) GetAllItems(ctx context.Context) ([]Product, error) {
	items, err := p.db.ReadAll(p.tableName)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to retrieve all products", "error", err)
		return nil, err
	}
	products := make([]Product, 0, len(items))
	for _, item := range items {
		products = append(products, item.(Product))
	}
	return products, nil
}

func (p *Inventory) GetStats() map[string]int64 {
	return p.mc.GetStats()
}

func generateID() string {
	return uuid.New().String()
}
