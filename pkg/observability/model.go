// Copyright 2025 Jacob Philip. All rights reserved.
package observability

import "sync"

type OperationType string

const (
	SnapshotIntervalSeconds = 30
)

const (
	OpInsert OperationType = "insert"
	OpUpdate OperationType = "update"
	OpDelete OperationType = "delete"
	OpGet    OperationType = "get"
	OpList   OperationType = "list"
)

type OperationMetric struct {
	Operation OperationType
	Success   bool
}
type MetricsCollector struct {
	operations chan OperationMetric
	shutdown   chan struct{}
	stats      map[OperationType]int64
	// we keep a copy of statsSnapshot for reporting without locks
	// this is updated periodically from statsCurrent
	// to avoid locking during metric reporting
	statsSnapshot sync.Map
}
