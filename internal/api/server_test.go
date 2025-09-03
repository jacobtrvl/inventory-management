// Copyright 2025 Jacob Philip. All rights reserved.
package api

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jacobtrvl/inventory-management/internal/inventory"
	"github.com/jacobtrvl/inventory-management/internal/store"
	"github.com/jacobtrvl/inventory-management/pkg/observability"
)

func TestAddAndGet(t *testing.T) {
	i := inventory.NewInventory(context.Background(),
		"products", store.NewMemDb(), observability.NewMetricsCollector())
	SetupRouter(context.Background(), i)
	body := `{"id":"1","name":"Test Product","price":9.99,"stock":100}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	w := httptest.NewRecorder()
	r := SetupRouter(context.Background(), i)
	r.e.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("Expected status 201, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/products/1", nil)
	w = httptest.NewRecorder()
	r.e.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	respBody := strings.TrimSpace(w.Body.String())
	if !strings.Contains(respBody, `"id":"1"`) {
		t.Fatalf("Expected response to contain product id 1, got %s", respBody)
	}


	req = httptest.NewRequest("GET", "/products", nil)
	w = httptest.NewRecorder()
	r.e.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	respBody = strings.TrimSpace(w.Body.String())
	if !strings.Contains(respBody, `"id":"1"`) {
		t.Fatalf("Expected response to contain product id 1, got %s", respBody)
	}

}
