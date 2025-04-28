package pid

import (
	"github.com/shirou/gopsutil/v4/process"
	"strings"
)

func FindPIdsByName(processName string) ([]int32, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var pIds []int32
	for _, proc := range processes {
		name, _ := proc.Name()
		if strings.Contains(strings.ToLower(name), strings.ToLower(processName)) {
			pid := proc.Pid
			pIds = append(pIds, pid)
		}
	}
	return pIds, nil
}
