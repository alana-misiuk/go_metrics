package repository

import "errors"

var ErrNotFound = errors.New("metric not found")

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.gauges[name] = value
	return nil
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	value, ok := s.gauges[name]
	if !ok {
		return 0, ErrNotFound
	}
	return value, nil
}

func (s *MemStorage) UpdateCounter(name string, value int64) error {
	s.counters[name] = s.counters[name] + value
	return nil
}

func (s *MemStorage) GetCounter(name string) (int64, error) {
	value, ok := s.counters[name]
	if !ok {
		return 0, ErrNotFound
	}
	return value, nil
}

func (s *MemStorage) GetAllGauges() map[string]float64 {
	return s.gauges
}

func (s *MemStorage) GetAllCounters() map[string]int64 {
	return s.counters
}
