package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetCPU1() (float64, error) {
	c, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil {
		return 0, err
	}
	if len(c) > 0 {
		return c[0], nil
	}
	return 0, nil

}

func GetMEM1() (float64, error) {
	m, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return m.UsedPercent, nil

}
