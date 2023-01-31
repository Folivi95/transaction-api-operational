package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
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
			"ingress_counter": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "ingress_counter",
				Help: "Counter for ingress topic",
				// this list represents the tags name for this metric
			}, []string{"topic"}),
			"egress_counter": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "egress_counter",
				Help: "Counter for egress topic",
				// this list represents the tags name for this metric
			}, []string{"status", "topic"}),
		},
		histograms: map[string]*prometheus.HistogramVec{
			"transaction_tokenization_time_ms": promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name: "transaction_tokenization_time_ms",
				Help: "Time to consume, process, and write back tokenized transaction",
				// this list represents the tags name for this metric
			}, []string{"status", "topic"}),
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

func NewPrometheusHook() func(entry zapcore.Entry) error {
	metric := promauto.NewCounter(prometheus.CounterOpts{
		Name: "error_log_entries",
		Help: "Metric for the number of errors in the logs",
	})

	return func(entry zapcore.Entry) error {
		if entry.Level == zapcore.ErrorLevel {
			metric.Inc()
		}
		return nil
	}
}
