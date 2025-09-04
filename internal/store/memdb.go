// Copyright 2025 Jacob Philip. All rights reserved.
// Store package implements in memory databse operations
// Design:
// MemDb is a collection of tables
// Table can hold any type of key-value pairs
// We use two levels of locking
// - First level lock is for controlling access to the tables map (sync.Map)
// - Second level lock is for controlling access to the data map within each table
// This allows concurrent access to different tables, which improves performance
// Note: Concept of table is not necessary in given task, but introduced to
// demonstrate an extendable design.
// An alternative design is to use in-memory db packages like hashicorp/memdb

// For supporting efficient range queries and key lookups, I am following a design
// where each table holds a slice of data and a map of key to index in the slice.
// This allows O(1) lookups and efficient range queries.
// Trade-off is that delete operation is O(n) due to slice reallocation.
// In given task, delete operation can be assumed to be rare compared to get/add operations.
// A product table might be read heavy, with occasional writes and rare deletes.

package store

import (
	"fmt"
	"sync"
)

type MemDb struct {
	// tables holds the mapping of table names to their data.
	// Read-heavy data structure.
	// sync.Map is optimized for concurrent reads and infrequent writes.
	// But we are losing type safety by using sync.Map.
	// In production implementation, I may prefer a map with Mutex
	tables *sync.Map
}

type value struct {
	indexMap  map[any]int
	dataSlice []any
	mutex     sync.RWMutex
}

// NewMemDb initializes a new in-memory database.
func NewMemDb() *MemDb {
	t := &sync.Map{}
	return &MemDb{
		tables: t,
	}
}

// CreateTable creates a new table in the memdb if it does not already exist.
func (m *MemDb) CreateTable(name string) error {
	if m.tables == nil {
		return fmt.Errorf("tables map not initialized")
	}

	v := &value{
		indexMap:  make(map[any]int),
		dataSlice: []any{},
		mutex:     sync.RWMutex{},
	}

	m.tables.LoadOrStore(name, v)
	return nil
}

// DeleteTable deletes a table from the memdb.
func (m *MemDb) DeleteTable(name string) error {
	if m.tables == nil {
		return fmt.Errorf("tables map not initialized")
	}
	if _, exists := m.tables.Load(name); !exists {
		return fmt.Errorf("table %s does not exist", name)
	}
	m.tables.Delete(name)
	return nil
}

// Table map value cannot be updated, hence its safe to operate on data map
// without table lock. An expection is when table itself is deleted, but since operation
// that calls this function happens before table deletion, we are operating on a
// valid data at that point of time.
func (m *MemDb) getDataMap(table string) (*value, error) {
	t, exists := m.tables.Load(table)
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", table)
	}
	return t.(*value), nil
}

// Write inserts an item into the specified table in the memdb.
func (m *MemDb) Write(table string, key any, item any) error {
	v, err := m.getDataMap(table)
	if err != nil {
		return err
	}
	v.mutex.Lock()
	defer v.mutex.Unlock()
	if index, exists := v.indexMap[key]; exists {
		v.dataSlice[index] = item
		return nil
	}
	v.dataSlice = append(v.dataSlice, item)
	v.indexMap[key] = len(v.dataSlice) - 1
	return nil
}

// Read retrieves an item from the specified table and index in the memdb.
func (m *MemDb) Read(table string, id any) (any, error) {
	t, err := m.getDataMap(table)
	if err != nil {
		return nil, err
	}
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	index, exists := t.indexMap[id]
	if !exists {
		return nil, fmt.Errorf("item with id %s not found in table %s", id, table)
	}
	return t.dataSlice[index], nil
}

// ReadRange retrieves all itemms within the specified range [start, end) from the table.
// Returns slice of items, EOF status, and error (if any).
func (m *MemDb) ReadRange(table string, start, end int) ([]any, bool, error) {
	t, err := m.getDataMap(table)
	if err != nil {
		return nil, false, err
	}
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if start < 0 {
		start = 0
	}
	if end > len(t.dataSlice) {
		end = len(t.dataSlice)
	}
	if start > end {
		return nil, false, fmt.Errorf("invalid start or end index")
	}
	result := make([]any, end-start)
	copy(result, t.dataSlice[start:end])
	return result, end >= len(t.dataSlice), nil
}

// ReadAll retrieves all items from the specified table in the memdb.
func (m *MemDb) ReadAll(table string) ([]any, error) {
	t, err := m.getDataMap(table)
	if err != nil {
		return nil, err
	}
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	result := make([]any, len(t.dataSlice))
	copy(result, t.dataSlice)
	return result, nil
}

// Delete deletes an item from the specified table in the memdb.
// This is an O(n) operation due to slice reallocation and index remapping.
// We are assuming delete operations are rare compared to read/write operations.
func (m *MemDb) Delete(table string, id any) error {
	v, err := m.getDataMap(table)
	if err != nil {
		return err
	}
	v.mutex.Lock()
	defer v.mutex.Unlock()

	indexToDelete, exists := v.indexMap[id]
	if !exists {
		return fmt.Errorf("item with id %v not found in table %s", id, table)
	}

	// Remove the item from slice
	copy(v.dataSlice[indexToDelete:], v.dataSlice[indexToDelete+1:])
	v.dataSlice = v.dataSlice[:len(v.dataSlice)-1]

	// Remove from index map
	delete(v.indexMap, id)

	// Update all indices that were shifted down
	for key, index := range v.indexMap {
		if index > indexToDelete {
			v.indexMap[key] = index - 1
		}
	}

	return nil
}
