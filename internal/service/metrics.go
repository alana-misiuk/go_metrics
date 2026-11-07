package service

import (
	"errors"
	"fmt"
	models "go_metrics/internal/model"
	"go_metrics/internal/repository"
	"strconv"
)

var (
	ErrNameRequired = errors.New("metric name is required")
	ErrUnknownType  = errors.New("unknown metric type")
	ErrInvalidValue = errors.New("invalid metric value")
	ErrNotFound     = errors.New("metric not found")
)

type Storage interface {
	UpdateGauge(name string, value float64) error
	GetGauge(name string) (float64, error)
	GetAllGauges() map[string]float64

	UpdateCounter(name string, value int64) error
	GetCounter(name string) (int64, error)
	GetAllCounters() map[string]int64
}

type MetricsService struct {
	storage Storage
}

func NewMetricsService(storage Storage) *MetricsService {
	return &MetricsService{storage: storage}
}

func (s *MetricsService) Update(metricType, name, rawValue string) error {
	if name == "" {
		return ErrNameRequired
	}

	switch metricType {
	case models.Gauge:
		value, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			return fmt.Errorf("%w: %q is not a float", ErrInvalidValue, rawValue)
		}
		if err := s.storage.UpdateGauge(name, value); err != nil {
			return fmt.Errorf("update gauge %q: %w", name, err)
		}
	case models.Counter:
		value, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %q is not an integer", ErrInvalidValue, rawValue)
		}
		if err := s.storage.UpdateCounter(name, value); err != nil {
			return fmt.Errorf("update counter %q: %w", name, err)
		}
	default:
		return fmt.Errorf("%w: %q", ErrUnknownType, metricType)
	}
	return nil
}

func (s *MetricsService) Get(metricType, name string) (models.Metrics, error) {
	m := models.Metrics{ID: name, MType: metricType}

	switch metricType {
	case models.Gauge:
		value, err := s.storage.GetGauge(name)
		if err != nil {
			return models.Metrics{}, notFound(err)
		}
		m.Value = &value
	case models.Counter:
		value, err := s.storage.GetCounter(name)
		if err != nil {
			return models.Metrics{}, notFound(err)
		}
		m.Delta = &value
	default:
		return models.Metrics{}, fmt.Errorf("%w: %q", ErrUnknownType, metricType)
	}
	return m, nil
}

func (s *MetricsService) GetAll() (gauges map[string]float64, counters map[string]int64) {
	return s.storage.GetAllGauges(), s.storage.GetAllCounters()
}

func notFound(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
