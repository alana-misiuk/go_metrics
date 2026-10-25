package agent

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
)

type AgentStats struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	PollCount     int64
	RandomValue   float64
}

func CollectMetrics(stats *AgentStats) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	stats.Alloc = float64(mem.Alloc)
	stats.BuckHashSys = float64(mem.BuckHashSys)
	stats.Frees = float64(mem.Frees)
	stats.GCCPUFraction = mem.GCCPUFraction
	stats.GCSys = float64(mem.GCSys)
	stats.HeapAlloc = float64(mem.HeapAlloc)
	stats.HeapIdle = float64(mem.HeapIdle)
	stats.HeapInuse = float64(mem.HeapInuse)
	stats.HeapObjects = float64(mem.HeapObjects)
	stats.HeapReleased = float64(mem.HeapReleased)
	stats.HeapSys = float64(mem.HeapSys)
	stats.LastGC = float64(mem.LastGC)
	stats.Lookups = float64(mem.Lookups)
	stats.MCacheInuse = float64(mem.MCacheInuse)
	stats.MCacheSys = float64(mem.MCacheSys)
	stats.MSpanInuse = float64(mem.MSpanInuse)
	stats.MSpanSys = float64(mem.MSpanSys)
	stats.Mallocs = float64(mem.Mallocs)
	stats.NextGC = float64(mem.NextGC)
	stats.NumForcedGC = float64(mem.NumForcedGC)
	stats.NumGC = float64(mem.NumGC)
	stats.OtherSys = float64(mem.OtherSys)
	stats.PauseTotalNs = float64(mem.PauseTotalNs)
	stats.StackInuse = float64(mem.StackInuse)
	stats.StackSys = float64(mem.StackSys)
	stats.Sys = float64(mem.Sys)
	stats.TotalAlloc = float64(mem.TotalAlloc)
	stats.PollCount++
	stats.RandomValue = rand.Float64()
}

func SendMetrics(stats *AgentStats) {
	const base = "http://localhost:8080/update"

	gauges := map[string]float64{
		"Alloc":         stats.Alloc,
		"BuckHashSys":   stats.BuckHashSys,
		"Frees":         stats.Frees,
		"GCCPUFraction": stats.GCCPUFraction,
		"GCSys":         stats.GCSys,
		"HeapAlloc":     stats.HeapAlloc,
		"HeapIdle":      stats.HeapIdle,
		"HeapInuse":     stats.HeapInuse,
		"HeapObjects":   stats.HeapObjects,
		"HeapReleased":  stats.HeapReleased,
		"HeapSys":       stats.HeapSys,
		"LastGC":        stats.LastGC,
		"Lookups":       stats.Lookups,
		"MCacheInuse":   stats.MCacheInuse,
		"MCacheSys":     stats.MCacheSys,
		"MSpanInuse":    stats.MSpanInuse,
		"MSpanSys":      stats.MSpanSys,
		"Mallocs":       stats.Mallocs,
		"NextGC":        stats.NextGC,
		"NumForcedGC":   stats.NumForcedGC,
		"NumGC":         stats.NumGC,
		"OtherSys":      stats.OtherSys,
		"PauseTotalNs":  stats.PauseTotalNs,
		"StackInuse":    stats.StackInuse,
		"StackSys":      stats.StackSys,
		"Sys":           stats.Sys,
		"TotalAlloc":    stats.TotalAlloc,
		"RandomValue":   stats.RandomValue,
	}

	for name, value := range gauges {
		url := fmt.Sprintf("%s/gauge/%s/%v", base, name, value)
		http.Post(url, "text/plain", nil)
	}

	url := fmt.Sprintf("%s/counter/PollCount/%d", base, stats.PollCount)
	http.Post(url, "text/plain", nil)
}