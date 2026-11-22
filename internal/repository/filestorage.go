package repository

import (
	"encoding/json"
	models "go_metrics/internal/model"
	"os"
)

type FileStorage struct {
	path string
	mem  *MemStorage
}

func NewFileStorage(path string, mem *MemStorage) *FileStorage {
	return &FileStorage{path: path, mem: mem}
}

func (f *FileStorage) Save() error {
	gauges := f.mem.GetAllGauges()
	counters := f.mem.GetAllCounters()

	metrics := make([]models.Metrics, 0, len(gauges)+len(counters))
	for name, value := range gauges {
		v := value
		metrics = append(metrics, models.Metrics{ID: name, MType: models.Gauge, Value: &v})
	}
	for name, delta := range counters {
		d := delta
		metrics = append(metrics, models.Metrics{ID: name, MType: models.Counter, Delta: &d})
	}

	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0644)
}

func (f *FileStorage) Restore() error {
	data, err := os.ReadFile(f.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var metrics []models.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	for _, m := range metrics {
		switch m.MType {
		case models.Gauge:
			if m.Value != nil {
				f.mem.UpdateGauge(m.ID, *m.Value)
			}
		case models.Counter:
			if m.Delta != nil {
				f.mem.UpdateCounter(m.ID, *m.Delta)
			}
		}
	}
	return nil
}
