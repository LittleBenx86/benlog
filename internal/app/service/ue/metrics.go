package ue

import (
	"fmt"
	"os"

	"github.com/LittleBenx86/Benlog/internal/global/dependencies"

	"github.com/shirou/gopsutil/v3/process"
)

type MetricsService struct {
	*dependencies.Dependencies
}

func NewMetricsService(d *dependencies.Dependencies) *MetricsService {
	return &MetricsService{
		Dependencies: d,
	}
}

func (m *MetricsService) GetMetrics() (string, error) {
	return "empty metrics", nil
}

func (m *MetricsService) GetCpuMetrics() (string, error) {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f", cpuPercent), nil
}

func (m *MetricsService) GetMemMetrics() (string, error) {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}

	pmi, err := p.MemoryInfo()
	if err != nil {
		return "", err
	}

	realProcessMemMB := float64(pmi.RSS) / 1024 / 1024
	return fmt.Sprintf("%.2f", realProcessMemMB), nil
}
