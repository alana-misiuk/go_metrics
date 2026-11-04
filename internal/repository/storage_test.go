package repository

import "testing"

func TestUpdateCounter_Accumulates(t *testing.T) {
	s := NewMemStorage()

	s.UpdateCounter("hits", 5)
	s.UpdateCounter("hits", 3)

	value, err := s.GetCounter("hits")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 8 {
		t.Errorf("expected 8, got %d", value)
	}
}

func TestUpdateGauge_Replaces(t *testing.T) {
	s := NewMemStorage()

	s.UpdateGauge("cpu", 1.0)
	s.UpdateGauge("cpu", 2.0)

	value, err := s.GetGauge("cpu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 2.0 {
		t.Errorf("expected 2.0, got %v", value)
	}
}

func TestGetGauge_NotFound(t *testing.T) {
	s := NewMemStorage()

	_, err := s.GetGauge("nonexistent")
	if err == nil {
		t.Error("expected error for missing metric, got nil")
	}
}
