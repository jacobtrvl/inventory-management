// Copyright 2025 Jacob Philip. All rights reserved.
package observability

import (
	"log/slog"
	"sync"
	"time"
)

func NewMetricsCollector() *MetricsCollector {
	mc := &MetricsCollector{
		operations:    make(chan OperationMetric, 100),
		shutdown:      make(chan struct{}),
		statsSnapshot: sync.Map{},
	}
	mc.stats = make(map[OperationType]int64)
	go mc.run()
	return mc
}

func (mc *MetricsCollector) run() {
	ticker := time.NewTicker(SnapshotIntervalSeconds * time.Second)
	defer ticker.Stop()
	for {
		select {
		case op := <-mc.operations:
			key := op.Operation
			if !op.Success {
				key = OperationType(string(op.Operation) + "_fail")
			}
			// Locking is not required since this is already synchronized via channel & select
			mc.stats[key]++
		case <-ticker.C:
			// Take snapshot of current stats. Deep copy with locks (sync.Map).
			mc.TakeSnapshot()
		case <-mc.shutdown:
			slog.Info("MetricsCollector shutting down")
			return
		}
	}
}

func (mc *MetricsCollector) TakeSnapshot() {
	for k, v := range mc.stats {
		mc.statsSnapshot.Store(string(k), v)
	}
}

func (mc *MetricsCollector) GetStats() map[string]int64 {
	stats := make(map[string]int64)
	mc.statsSnapshot.Range(func(k, v any) bool {
		stats[k.(string)] = v.(int64)
		return true
	})
	return stats
}

func (mc *MetricsCollector) RecordOperation(op OperationType, success bool) {
	mc.operations <- OperationMetric{
		Operation: op,
		Success:   success,
	}
}

func (mc *MetricsCollector) Shutdown() {
	close(mc.shutdown)
}
