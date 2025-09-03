// Copyright 2025 Jacob Philip. All rights reserved.
package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jacobtrvl/inventory-management/internal/inventory"
)

type Router struct {
	e *gin.Engine
	i *inventory.Inventory
}

type ResponseFormat struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
	Meta  any    `json:"meta,omitempty"`
}

func SetupRouter(ctx context.Context, i *inventory.Inventory) Router {
	router := gin.Default()
	r := Router{
		e: router,
		i: i,
	}

	router.GET("/products", r.ListProducts)
	router.GET("/products/:id", r.getProducts)
	router.POST("/products", r.addProduct)
	router.PUT("/products/:id", r.updateProduct)
	router.DELETE("/products/:id", r.deleteProduct)
	// Basic metrics endpoint returning JSON format
	// In production system, should be replaced with Prometheus Instrumentation
	router.GET("/metrics", r.metricsHandler)

	return r
}

func (r Router) Run(addr ...string) error {
	return r.e.Run(addr...)
}

func (r Router) Handler(addr ...string) http.Handler {
	return r.e.Handler()
}

func (r Router) metricsHandler(c *gin.Context) {
	i := r.i
	// In a large system, its important to pass right status codes and error messages
	// This involves status code or its mapping being send from internal packages
	// We are currently not doing that for simplicity
	if i == nil {
		c.JSON(http.StatusInternalServerError, ResponseFormat{Error: "inventory not initialized"})
		return
	}

	stats := i.GetStats()
	if len(stats) == 0 {
		c.JSON(http.StatusOK, ResponseFormat{Data: "No metrics available yet"})
		return
	}

	c.JSON(http.StatusOK, ResponseFormat{Data: stats})
}

func (r Router) getProducts(c *gin.Context) {
	i := r.i
	id := c.Param("id")
	product, status, err := i.Get(c.Request.Context(), id)
	if err != nil {
		errMessage := "Product not found: " + err.Error()
		c.JSON(status, ResponseFormat{Error: errMessage})
		return
	}
	c.JSON(http.StatusOK, ResponseFormat{Data: product})
}

func (r Router) ListProducts(c *gin.Context) {
	i := r.i
	var limit, page *int

	if ls := c.Query("limit"); ls != "" {
		l, err := strconv.Atoi(ls)
		if err != nil || l <= 0 {
			c.JSON(http.StatusBadRequest, ResponseFormat{Error: "Invalid limit parameter"})
			return
		}
		limit = &l
	}
	if ps := c.Query("page"); ps != "" {
		pg, err := strconv.Atoi(ps)
		if err != nil || pg <= 0 {
			c.JSON(http.StatusBadRequest, ResponseFormat{Error: "Invalid page parameter"})
			return
		}
		page = &pg
	}

	params := inventory.ListParams{
		Page:  page,
		Limit: limit,
	}

	products, meta, status, err := i.List(c.Request.Context(), params)
	if err != nil {
		errMessage := "Failed to retrieve products: " + err.Error()
		c.JSON(status, ResponseFormat{Error: errMessage})
		return
	}
	if limit == nil && page == nil {
		c.JSON(http.StatusOK, ResponseFormat{Data: products})
		return
	}


	c.JSON(http.StatusOK, ResponseFormat{
		Data: products,
		Meta: meta,
	})

}

func (r Router) addProduct(c *gin.Context) {
	i := r.i
	var productReq inventory.CreateRequest
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, ResponseFormat{Error: "Invalid input"})
		return
	}

	productID, status, err := i.Add(c.Request.Context(), productReq)
	if err != nil {
		c.JSON(status, ResponseFormat{Error: "Failed to add product"})
		return
	}
	c.JSON(http.StatusCreated, ResponseFormat{
		Data: map[string]string{
			"message":    fmt.Sprintf("Added product %s", productID),
			"product_id": productID,
		},
	})
}

func (r Router) updateProduct(c *gin.Context) {
	i := r.i
	id := c.Param("id")
	var updatedProduct inventory.UpdateRequest
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, ResponseFormat{Error: "Invalid input"})
		return
	}
	if status, err := i.Update(c.Request.Context(), id, updatedProduct); err != nil {
		errMessage := "Failed to update product: " + err.Error()
		c.JSON(status, ResponseFormat{Error: errMessage})
		return
	}
	c.JSON(http.StatusOK, ResponseFormat{
		Data: map[string]string{
			"message": "Product updated",
			"id":      id,
		},
	})
}

func (r Router) deleteProduct(c *gin.Context) {
	i := r.i
	id := c.Param("id")
	if status, err := i.Delete(c.Request.Context(), id); err != nil {
		errMessage := "Failed to delete product: " + err.Error()
		c.JSON(status, ResponseFormat{Error: errMessage})
		return
	}
	c.JSON(http.StatusOK, ResponseFormat{
		Data: map[string]string{
			"message": "Product deleted",
			"id":      id,
		},
	})
}
