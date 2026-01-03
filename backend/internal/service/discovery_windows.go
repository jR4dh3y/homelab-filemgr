//go:build windows

package service

import (
	"github.com/homelab/filemanager/internal/model"
	"github.com/homelab/filemanager/internal/pkg/filesystem"
)

// DiscoverMountPoints on Windows returns mount points as-is (no auto-discovery support yet)
func DiscoverMountPoints(fs filesystem.FS, mountPoints []model.MountPoint) []model.MountPoint {
	return mountPoints
}
