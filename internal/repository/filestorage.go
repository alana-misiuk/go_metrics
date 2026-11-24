package repository

import (
	"encoding/json"
	models "go_metrics/internal/model"
	"os"
	"path/filepath"
)

type FileStorage struct {
	path string
	mem  *MemStorage
}

func NewFileStorage(path string, mem *MemStorage, restore bool) (*FileStorage, error) {
	fs := &FileStorage{path: path, mem: mem}
	if restore {
		if err := fs.restore(); err != nil {
			return nil, err
		}
	}
	return fs, nil
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
	tmp, err := os.CreateTemp(filepath.Dir(f.path), ".metrics-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, f.path)
}

func (f *FileStorage) restore() error {
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
