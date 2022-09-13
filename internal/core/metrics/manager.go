package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Report scheduler metrics manager
type Manager interface {

	// Register handled job
	RegisterJob(dur time.Duration, count uint64, labels []string)

	// Increase counter by status
	AddCount(count uint64, labels []string) error

	// Updates number of currently active workers
	UpdateWorkersNumber(diff float64)
}

type ManagerConfig struct {
	Namespace     string
	Subsystem     string
	CounterDesc   string
	HistogramDesc string
	Buckets       []float64
	SupportGauge  bool
	Labels        []string
}

type managerImpl struct {
	counter                *prometheus.CounterVec
	timeHistogram          prometheus.Histogram
	workersGauge           prometheus.Gauge
	workersGaugeRegistered bool
}

// Creates metrics manager
func NewManager(config *ManagerConfig) (Manager, error) {
	manager := new(managerImpl)

	err := manager.init(config)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *managerImpl) RegisterJob(duration time.Duration, count uint64, labelValues []string) {
	err := m.AddCount(count, labelValues)
	if err == nil {
		m.timeHistogram.Observe(duration.Seconds())
	}
}

func (m *managerImpl) AddCount(count uint64, labelValues []string) error {
	metric, err := m.counter.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		return err
	}

	metric.Add(float64(count))

	return nil
}

func (m *managerImpl) UpdateWorkersNumber(diff float64) {
	if !m.workersGaugeRegistered {
		return
	}

	m.workersGauge.Add(diff)
}

func (m *managerImpl) init(config *ManagerConfig) (err error) {
	err = m.registerJobCounter(config.Namespace, config.Subsystem, config.CounterDesc, config.Labels)
	if err != nil {
		return err
	}

	err = m.registerTimeHistogram(config.Namespace, config.Subsystem, config.Buckets, config.HistogramDesc)
	if err != nil {
		return err
	}

	if config.SupportGauge {
		err = m.enableGauge(config.Namespace, config.Subsystem)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *managerImpl) enableGauge(namespace string, subsystem string) error {
	if m.workersGaugeRegistered {
		return nil
	}

	err := m.registerWorkersGauge(namespace, subsystem)
	if err != nil {
		return err
	}
	m.workersGaugeRegistered = true

	return nil
}

func (m *managerImpl) registerJobCounter(namespace string, subsystem string, description string, labels []string) error {
	if len(labels) == 0 {
		labels = []string{"job_status"}
	}

	m.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "handled_total",
			Help:      description,
		},
		labels,
	)

	err := prometheus.Register(m.counter)
	if err != nil {
		are, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			m.counter = are.ExistingCollector.(*prometheus.CounterVec)
		} else {
			return err
		}
	}

	return nil
}

func (m *managerImpl) registerTimeHistogram(
	namespace string,
	subsystem string,
	buckets []float64,
	description string,
) error {
	m.timeHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "handling_seconds",
		Help:      description,
		Buckets:   buckets,
	})

	err := prometheus.Register(m.timeHistogram)
	if err != nil {
		are, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			m.timeHistogram = are.ExistingCollector.(prometheus.Histogram)
		} else {
			return err
		}
	}

	return nil
}

func (m *managerImpl) registerWorkersGauge(namespace string, subsystem string) error {
	m.workersGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "workers",
		Help:      "Number of currently active workers.",
	})

	err := prometheus.Register(m.workersGauge)
	if err != nil {
		are, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			m.workersGauge = are.ExistingCollector.(prometheus.Gauge)
		} else {
			return err
		}
	}

	return nil
}
