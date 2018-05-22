package metrics

type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
}

type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
}

type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}
