package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
)

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
