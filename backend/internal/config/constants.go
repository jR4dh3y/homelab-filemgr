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
	"autofs": true, "overlay": true, "efivarfs": true, "nsfs": true,
	"ramfs": true, "rpc_pipefs": true, "nfsd": true, "squashfs": true,
}

// ExcludedMountPointPrefixes contains mount point path prefixes that should be filtered out
// These are typically Docker bind mounts, system paths, or container-specific mounts
// Also includes /host_root equivalents for when running in Docker
var ExcludedMountPointPrefixes = []string{
	// Container-internal mounts
	"/etc/",                   // Docker bind mounts like /etc/hosts, /etc/resolv.conf, /etc/hostname
	"/app/",                   // Application config bind mounts
	"/proc/",                  // Process filesystem
	"/sys/",                   // System filesystem
	"/dev/",                   // Device filesystem
	"/run/",                   // Runtime data
	"/var/lib/docker",         // Docker internal paths
	"/snap/",                  // Snap packages
	"/boot/",                  // Boot partition
	// Host paths via /host_root bind mount - exclude duplicates and system paths
	"/host_root/etc/",         // Host's etc directory
	"/host_root/proc/",        // Host's proc
	"/host_root/sys/",         // Host's sys
	"/host_root/dev/",         // Host's dev
	"/host_root/run/",         // Host's run
	"/host_root/var/lib/docker", // Host's docker
	"/host_root/snap/",        // Host's snap packages
	"/host_root/boot/",        // Host's boot partition
	// These are duplicates - already accessible via direct mounts
	"/host_root/media/",       // Already accessible via /media/devmon mount
	"/host_root/home/",        // Already accessible via /home/user mount
}

// ExcludedMountPointSuffixes contains mount point path suffixes that should be filtered out
var ExcludedMountPointSuffixes = []string{
	"/efivars", // EFI variables
	"/efi",     // EFI partition (small boot partition)
}

// AllowedMountPointPrefixes contains mount point prefixes that should be included
// if they pass the exclusion checks above
var AllowedMountPointPrefixes = []string{
	"/",                  // Root filesystem (exact match handled specially)
	"/home",              // User home directories (container)
	"/media/",            // Removable media (container)
	"/mnt/",              // Manual mounts (container)
	"/host_root",         // Host root (exact match - the main host filesystem)
	"/host_root/mnt/",    // Host manual mounts (rclone, etc. that aren't in /media)
}

