package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsClient struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

// New inits the metrics and automatically registers the metrics with the default prometheus default register.
// When adding tags keep in mind tag cardinality, prometheus client keeps tags in memory, so you should not use
// any unbounded value in the tag value.
func New() MetricsClient {
	client := MetricsClient{
		counters: map[string]*prometheus.CounterVec{
			"db_write_counter": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "db_write_counter",
				Help: "A counter for write operations to the DB",
			}, []string{"acquiring_host"}),
			"db_read_counter": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "db_read_counter",
				Help: "A counter for read operations to the DB",
			}, []string{"status"}),
			"db_transacation_counter_request": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "db_transacation_counter_request",
				Help: "A counter for read transactions operations to the DB",
			}, []string{"status"}),
			"invocation_count": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "function_counter",
				Help: "A counter for invocations of a specific function",
			}, []string{"function"}),
			"invocation_result": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "invocation_result",
				Help: "A counter for the success/errors of invoking a function",
			}, []string{"function", "status"}),
			// TODO: Add your new counters here
		},
		histograms: map[string]*prometheus.HistogramVec{
			"sample-histogram": promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name: "sample_histogram",
				Help: "A sample on how to use histograms with prometheus",
				// this list represents the tags name for this metric
			}, []string{"event"}),
			"execution_time_ms": promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name: "execution_time_ms",
				Help: "Time to process",
				// this list represents the tags name for this metric
			}, []string{"function", "status"}),
		},
	}

	return client
}

func (m MetricsClient) Histogram(name string, value float64, tags []string) {
	if h, ok := m.histograms[name]; ok {
		h.WithLabelValues(tags...).Observe(value)

		return
	}
	// TODO: log if metric not found
}

func (m MetricsClient) Count(name string, value int64, tags []string) {
	if h, ok := m.counters[name]; ok {
		h.WithLabelValues(tags...).Add(float64(value))

		return
	}
	// TODO: log if metric not found
}
