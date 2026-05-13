package translator

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const subsystem = "translator"

const (
	translationCountMetricName    = "translation_count"
	translationDurationMetricName = "translation_duration_seconds"
	translationFailuresMetricName = "translation_failures_count"
)

const (
	translationCountMetricDesc    = "Number of times Translate was called"
	translationDurationMetricDesc = "Time spend translating resources"
	translationFailuresMetricDesc = "Number of resource translation failures"
)

var (
	registerOnce     sync.Once
	translationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      translationCountMetricName,
			Help:      translationCountMetricDesc,
			Subsystem: subsystem,
		},
		[]string{"result"},
	)
	translationDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      translationDurationMetricName,
			Help:      translationDurationMetricDesc,
			Subsystem: subsystem,
			Buckets:   []float64{0.001, 0.01, 0.1, 1, 5, 10},
		},
		[]string{"result"},
	)
	translationFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      translationFailuresMetricName,
			Help:      translationFailuresMetricDesc,
			Subsystem: subsystem,
		},
	)
)

type TranslatorMetricsProvider struct{}

func (TranslatorMetricsProvider) NewTranslationFailureMetric() prometheus.Counter {
	return translationFailures
}

func (TranslatorMetricsProvider) NewTranslationDurationMetric(result string) prometheus.Observer {
	return translationDurationSeconds.With(prometheus.Labels{"result": result})
}

func (TranslatorMetricsProvider) NewTranslationCountMetric(result string) prometheus.Counter {
	return translationCount.With(prometheus.Labels{"result": result})
}

func RegisterMetrics(registry prometheus.Registerer) {
	registerOnce.Do(func() {
		registry.MustRegister(translationCount, translationDurationSeconds, translationFailures)
	})
}
