// Copyright 2025 Jacob Philip. All rights reserved.
package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
	db := NewMemDb()

	err := db.CreateTable("testTable")
	assert.NoError(t, err, "expected no error on creating table")

	err = db.CreateTable("testTable")
	assert.NoError(t, err, "expected no error on creating existing table")
}

func TestInsertAndGet(t *testing.T) {
	db := NewMemDb()
	err := db.CreateTable("testTable")
	assert.NoError(t, err, "expected no error on creating table")

	err = db.Write("testTable", "key1", "value1")
	assert.NoError(t, err, "expected no error on inserting value")

	err = db.Write("testTable", "key2", "value2")
	assert.NoError(t, err, "expected no error on inserting value")

	val, err := db.Read("testTable", "key1")
	assert.NoError(t, err, "expected no error on getting value")
	assert.Equal(t, "value1", val, "expected value to match inserted value")

	val, err = db.Read("testTable", "key3")
	assert.Error(t, err, "expected error on getting non-existent key")
	assert.Nil(t, val, "expected nil value for non-existent key")

	err = db.Delete("testTable", "key1")
	assert.NoError(t, err, "expected no error on deleting existing key")

	val, err = db.Read("testTable", "key1")
	assert.Error(t, err, "expected error on getting deleted key")
	assert.Nil(t, val, "expected nil value for deleted key")

	items, err := db.ReadAll("testTable")
	assert.NoError(t, err, "expected no error on getting all items")
	assert.Equal(t, 1, len(items), "expected one item in the table")
	assert.Equal(t, "value2", items[0], "expected remaining item to match inserted value")
}
