package monitor

import (
	"context"
	"peekaping/src/modules/heartbeat"
	"time"
)

// UptimeResult holds the calculated uptime and average ping for a period.
type UptimeResult struct {
	Uptime  float64                `json:"uptime"` // 0-100
	AvgPing *float64               `json:"avgPing"`
	Points  []heartbeat.ChartPoint `json:"points"`
}

// UptimeGranularity defines the time bucket for aggregation.
type UptimeGranularity string

const (
	GranularityMinute UptimeGranularity = "minute"
	GranularityHour   UptimeGranularity = "hour"
	GranularityDay    UptimeGranularity = "day"
)

// UptimeCalculator provides methods to calculate uptime for a monitor.
type UptimeCalculator struct {
	HeartbeatService heartbeat.Service
}

func NewUptimeCalculator(heartbeatService heartbeat.Service) *UptimeCalculator {
	return &UptimeCalculator{HeartbeatService: heartbeatService}
}

// GetUptime calculates uptime for a monitor over a given period and granularity.
func (u *UptimeCalculator) GetUptime(ctx context.Context, monitorID string, since time.Time, until time.Time, granularity UptimeGranularity) (*UptimeResult, error) {
	// Use the heartbeat service to fetch chart points for the period
	points, err := u.fetchChartPoints(ctx, monitorID, since, until, granularity)
	if err != nil {
		return nil, err
	}

	var totalUp, totalDown int
	var totalPing float64
	var upCount int

	for _, pt := range points {
		totalUp += pt.Up
		totalDown += pt.Down
		if pt.Up > 0 {
			totalPing += pt.AvgPing * float64(pt.Up)
			upCount += pt.Up
		}
	}

	// Calculate uptime percentage
	total := totalUp + totalDown
	uptime := 0.0
	if total > 0 {
		uptime = float64(totalUp) / float64(total) * 100
	}

	// Calculate average ping
	var avgPing *float64
	if upCount > 0 {
		avg := totalPing / float64(upCount)
		avgPing = &avg
	}

	return &UptimeResult{
		Uptime:  uptime,
		AvgPing: avgPing,
		Points:  points,
	}, nil
}

// fetchChartPoints fetches heartbeat chart points for the given period and granularity.
func (u *UptimeCalculator) fetchChartPoints(ctx context.Context, monitorID string, since, until time.Time, granularity UptimeGranularity) ([]heartbeat.ChartPoint, error) {
	// For now, use the existing FindByMonitorIDAndTimeRange for minutely granularity.
	// For hourly/daily, you may want to add similar aggregation in the heartbeat repository.
	if granularity == GranularityMinute {
		points, err := u.HeartbeatService.FindByMonitorIDAndTimeRange(ctx, monitorID, since, until)
		if err != nil {
			return nil, err
		}
		return toChartPointSlice(points), nil
	}
	// For hour/day, you would need to implement aggregation in heartbeat repository.
	// For now, fallback to minutely and aggregate in Go.
	minutelyPoints, err := u.HeartbeatService.FindByMonitorIDAndTimeRange(ctx, monitorID, since, until)
	if err != nil {
		return nil, err
	}
	return aggregatePoints(minutelyPoints, granularity), nil
}

// toChartPointSlice converts []*ChartPoint to []ChartPoint
func toChartPointSlice(points []*heartbeat.ChartPoint) []heartbeat.ChartPoint {
	res := make([]heartbeat.ChartPoint, len(points))
	for i, p := range points {
		if p != nil {
			res[i] = *p
		}
	}
	return res
}

// aggregatePoints aggregates minutely points into hourly or daily buckets.
func aggregatePoints(points []*heartbeat.ChartPoint, granularity UptimeGranularity) []heartbeat.ChartPoint {
	bucketMap := make(map[int64]*heartbeat.ChartPoint)
	var bucketSize int64
	if granularity == GranularityHour {
		bucketSize = 3600 * 1000 // ms
	} else if granularity == GranularityDay {
		bucketSize = 86400 * 1000 // ms
	} else {
		return toChartPointSlice(points)
	}

	for _, pt := range points {
		if pt == nil {
			continue
		}
		bucket := pt.Timestamp - (pt.Timestamp % bucketSize)
		b, ok := bucketMap[bucket]
		if !ok {
			b = &heartbeat.ChartPoint{Timestamp: bucket, MinPing: pt.MinPing, MaxPing: pt.MaxPing}
			bucketMap[bucket] = b
		}
		b.Up += pt.Up
		b.Down += pt.Down
		if pt.MinPing < b.MinPing || b.MinPing == 0 {
			b.MinPing = pt.MinPing
		}
		if pt.MaxPing > b.MaxPing {
			b.MaxPing = pt.MaxPing
		}
		b.AvgPing += pt.AvgPing * float64(pt.Up)
	}
	// Finalize avgPing
	var result []heartbeat.ChartPoint
	for _, b := range bucketMap {
		if b.Up > 0 {
			b.AvgPing = b.AvgPing / float64(b.Up)
		}
		result = append(result, *b)
	}
	return result
}
