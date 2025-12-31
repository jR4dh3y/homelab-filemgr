//go:build !windows

package service

import "syscall"

// getDiskUsage returns disk usage statistics for the given path (Unix/Linux/macOS)
func getDiskUsage(path string) (*diskUsage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	usedPct := float64(0)
	if total > 0 {
		usedPct = float64(used) / float64(total) * 100
	}

	return &diskUsage{
		Total:   total,
		Free:    free,
		Used:    used,
		UsedPct: usedPct,
	}, nil
}
