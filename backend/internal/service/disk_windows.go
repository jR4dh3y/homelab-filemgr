//go:build windows

package service

import (
	"syscall"
	"unsafe"

	"github.com/homelab/filemanager/internal/config"
)

// getDiskUsage returns disk usage statistics for the given path (Windows)
func getDiskUsage(path string) (*diskUsage, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	var freeBytesAvailable, totalBytes, totalFreeBytes uint64

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	ret, _, err := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)

	if ret == 0 {
		return nil, err
	}

	used := totalBytes - totalFreeBytes
	usedPct := float64(0)
	if totalBytes > 0 {
		usedPct = float64(used) / float64(totalBytes) * config.PercentMultiplier
	}

	return &diskUsage{
		Total:      totalBytes,
		Free:       totalFreeBytes,
		Used:       used,
		UsedPct:    usedPct,
		Device:     path, // On Windows, the path is typically the drive letter
		FSType:     "",   // Could be extended to get filesystem type via GetVolumeInformation
		MountPoint: path,
	}, nil
}
