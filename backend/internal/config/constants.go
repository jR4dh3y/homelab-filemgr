package config

// Filesystem constants
const (
	// ProcMountsPath is the Linux path to the mounts file
	ProcMountsPath = "/proc/mounts"
)

// Disk usage constants
const (
	// PercentMultiplier is used to convert decimal to percentage
	PercentMultiplier = 100
)

// VirtualFilesystems contains filesystem types that should be filtered out during discovery
var VirtualFilesystems = map[string]bool{
	"sysfs": true, "proc": true, "devtmpfs": true, "devpts": true,
	"tmpfs": true, "securityfs": true, "cgroup": true, "cgroup2": true,
	"pstore": true, "debugfs": true, "tracefs": true, "configfs": true,
	"fusectl": true, "mqueue": true, "hugetlbfs": true, "binfmt_misc": true,
	"autofs": true, "overlay": true,
}
