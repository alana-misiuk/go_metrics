package repository

import "errors"

type MemStorage struct {
	gauges map[string] float64
	counters map[string] int64
} 

type Storage interface {
	UpdateGauge(name string, value float64)
	GetGauge(name string) (float64, error)

	UpdateCounter(name string, value int64)
	GetCounter(name string) (int64, error)
} 

func NewMemStorage() *MemStorage {
	return &MemStorage{
        gauges:   make(map[string]float64),
        counters: make(map[string]int64),
    }
}

func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.gauges[name] = value
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	value, ok := s.gauges[name]
	if !ok {
		return 0, errors.New("metric not found")
	}
	return value, nil
}

func (s *MemStorage) UpdateCounter(name string, value int64) {
	s.counters[name] = s.counters[name] + value
}

func (s *MemStorage) GetCounter(name string) (int64, error) {
	value, ok := s.counters[name]
	if !ok {
		return 0, errors.New("metric not found")
	}
	return value, nil
}