// Copyright 2025 Jacob Philip. All rights reserved.
package inventory

import (
	"context"
	"testing"

	"github.com/jacobtrvl/inventory-management/pkg/observability"
	"github.com/jacobtrvl/inventory-management/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestInventory(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	db := store.NewMemDb()
	// Using real MetricsCollector for testing, since mock is not implemented.
	mc := observability.NewMetricsCollector()
	inventory := NewInventory(ctx, "products", db, mc)

	p := CreateRequest{
		ID:    "1",
		Name:  "Test Product",
		Price: 9.99,
		Stock: 100,
	}

	// Test Add
	id, _, err := inventory.Add(ctx, p)
	assert.NoError(err)
	assert.Equal(p.ID, id)
	// Adding same product again should fail
	_, _, err = inventory.Add(ctx, p)
	assert.Error(err)

	// Test Get
	product, _, err := inventory.Get(ctx, "1")
	assert.NoError(err)
	assert.NotNil(product)
	assert.Equal(p.Name, product.Name)

	_, _, err = inventory.Get(ctx, "2")
	assert.Error(err)

	allProducts, _, _, err := inventory.List(ctx, ListParams{})
	assert.NoError(err)
	assert.NotEmpty(allProducts)
	page := 1
	limit := 10

	products, _, _, err := inventory.List(ctx, ListParams{Page: &page, Limit: &limit})
	assert.NoError(err)
	assert.NotEmpty(products)
	assert.Equal(products, allProducts)

	// Test Update
	newName := "Updated Product"
	u := UpdateRequest{
		Name:  &newName,
	}

	_, err = inventory.Update(ctx, "1", u)
	assert.NoError(err)
	product, _, err = inventory.Get(ctx, "1")
	assert.NoError(err)
	assert.Equal(newName, product.Name)
	assert.Equal(p.Price, product.Price) // Unchanged
	assert.Equal(p.Stock, product.Stock) // Unchanged

	_, err = inventory.Update(ctx, "2", UpdateRequest{})
	assert.Error(err)

	// Test Delete
	_, err = inventory.Delete(ctx, "1")
	assert.NoError(err)

	_, err = inventory.Delete(ctx, "2")
	assert.Error(err)

	// Test List
	allProducts, _, _, err = inventory.List(ctx, ListParams{})
	assert.NoError(err)
	assert.Empty(allProducts)

	products, _, _, err = inventory.List(ctx, ListParams{Page: &page, Limit: &limit})
	assert.NoError(err)
	assert.Empty(products)
}
