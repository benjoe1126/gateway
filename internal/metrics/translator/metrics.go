package translator

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const subsystem = "translator"

const (
	translationCountMetricName    = "envoy_gateway_translation_count"
	translationDurationMetricName = "envoy_gateway_translation_duration"
	translatedResourcesMetricName = "envoy_gateway_translated_resources"
	translationFailuresMetricName = "envoy_gateway_translation_failures"
)

const (
	translationCountMetricDesc    = "Number of times Translate was called"
	translationDurationMetricDesc = "Time spend translating resources"
	translatedResourcesMetricDesc = "Number of resources translated"
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
	translatedResources = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      translatedResourcesMetricName,
			Help:      translatedResourcesMetricDesc,
			Subsystem: subsystem,
		},
		[]string{"resource", "status"},
	)
	translationFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: translationFailuresMetricName,
			Help: translationFailuresMetricDesc,
		},
	)
)

type TranslatorMetricsProvider struct{}

func (TranslatorMetricsProvider) NewTranslationFailureMetric() prometheus.Counter {
	return translationFailures
}

func (TranslatorMetricsProvider) NewTranslatedResourcesMetric(resource, status string) prometheus.Counter {
	return translatedResources.With(prometheus.Labels{"resource": resource, "status": status})
}

func (TranslatorMetricsProvider) NewTranslationDurationMetric(result string) prometheus.Observer {
	return translationDurationSeconds.With(prometheus.Labels{"result": result})
}

func (TranslatorMetricsProvider) NewTranslationCountMetric(result string) prometheus.Counter {
	return translationCount.With(prometheus.Labels{"result": result})
}

func RegisterMetrics(registry prometheus.Registerer) {
	registerOnce.Do(func() {
		registry.MustRegister(translationCount, translationDurationSeconds, translatedResources, translationFailures)
	})
}
