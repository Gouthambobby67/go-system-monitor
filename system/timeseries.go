package system

import "time"

// TimeSeriesPoint represents a single datapoint in a time series
type TimeSeriesPoint struct {
    Timestamp time.Time
    Value     float64
}

// TimeSeries is a simple series of TimeSeriesPoints used by the UI for sparklines
type TimeSeries struct {
    Points []TimeSeriesPoint
}
