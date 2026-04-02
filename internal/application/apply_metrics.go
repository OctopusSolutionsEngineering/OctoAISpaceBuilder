package application

import (
	"sync"
	"time"
)

type metricsBucket struct {
	minuteEpoch int64
	success     int64
	failed      int64
}

type terraformApplyCallMetrics struct {
	mu             sync.Mutex
	buckets        [60]metricsBucket
	allTimeSuccess int64
	allTimeFailed  int64
}

type terraformApplyCallMetricsSnapshot struct {
	SuccessLast1Minute   int64
	FailedLast1Minute    int64
	SuccessLast15Minutes int64
	FailedLast15Minutes  int64
	SuccessLast1Hour     int64
	FailedLast1Hour      int64
	SuccessAllTime       int64
	FailedAllTime        int64
}

var terraformApplyMetrics = &terraformApplyCallMetrics{}
var terraformPlanMetrics = &terraformApplyCallMetrics{}

func (m *terraformApplyCallMetrics) Record(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	minuteEpoch := time.Now().Unix() / 60
	bucket := &m.buckets[minuteEpoch%60]

	if bucket.minuteEpoch != minuteEpoch {
		bucket.minuteEpoch = minuteEpoch
		bucket.success = 0
		bucket.failed = 0
	}

	if success {
		bucket.success++
		m.allTimeSuccess++
		return
	}

	bucket.failed++
	m.allTimeFailed++
}

func (m *terraformApplyCallMetrics) Snapshot(now time.Time) terraformApplyCallMetricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	currentMinute := now.Unix() / 60

	snapshot := terraformApplyCallMetricsSnapshot{
		SuccessAllTime: m.allTimeSuccess,
		FailedAllTime:  m.allTimeFailed,
	}

	for _, bucket := range m.buckets {
		if bucket.minuteEpoch == 0 {
			continue
		}

		age := currentMinute - bucket.minuteEpoch
		if age < 0 || age >= 60 {
			continue
		}

		snapshot.SuccessLast1Hour += bucket.success
		snapshot.FailedLast1Hour += bucket.failed

		if age < 15 {
			snapshot.SuccessLast15Minutes += bucket.success
			snapshot.FailedLast15Minutes += bucket.failed
		}

		if age < 1 {
			snapshot.SuccessLast1Minute += bucket.success
			snapshot.FailedLast1Minute += bucket.failed
		}
	}

	return snapshot
}
