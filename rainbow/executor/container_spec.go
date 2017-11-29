package executor

import (
	"context"
	ns "github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"path/filepath"
)

func containerSpec(ctx context.Context, id string) (*specs.Spec, error) {
	ns, err := ns.NamespaceRequired(ctx)
	if err != nil {
		return nil, err
	}

	spec := &specs.Spec{
		Version: specs.Version,
		Root: &specs.Root{
			Path: "rootfs",
		},
		Process: containerProcessSpec(),
		Mounts:  containerMountsSpec(),
		Linux:   containerLinuxSpec(ns, id),
	}
	return spec, nil
}

func containerProcessSpec() *specs.Process {
	return &specs.Process{
		Env:             []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		Cwd:             "/",
		NoNewPrivileges: true,
		User: specs.User{
			UID: 0,
			GID: 0,
		},
		Capabilities: &specs.LinuxCapabilities{
			Bounding:    containerProcessCapabilitiesSpec(),
			Permitted:   containerProcessCapabilitiesSpec(),
			Inheritable: containerProcessCapabilitiesSpec(),
			Effective:   containerProcessCapabilitiesSpec(),
		},
		Rlimits: []specs.POSIXRlimit{
			{
				Type: "RLIMIT_NOFILE",
				Hard: uint64(1024),
				Soft: uint64(1024),
			},
		},
	}
}

func containerProcessCapabilitiesSpec() []string {
	return []string{
		"CAP_CHOWN",
		"CAP_DAC_OVERRIDE",
		"CAP_FSETID",
		"CAP_FOWNER",
		"CAP_MKNOD",
		"CAP_NET_RAW",
		"CAP_SETGID",
		"CAP_SETUID",
		"CAP_SETFCAP",
		"CAP_SETPCAP",
		"CAP_NET_BIND_SERVICE",
		"CAP_SYS_CHROOT",
		"CAP_KILL",
		"CAP_AUDIT_WRITE",
	}
}

func containerMountsSpec() []specs.Mount {
	return []specs.Mount{
		{
			Destination: "/proc",
			Type:        "proc",
			Source:      "proc",
		},
		{
			Destination: "/dev",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options:     []string{"nosuid", "strictatime", "mode=755", "size=65536k"},
		},
		{
			Destination: "/dev/pts",
			Type:        "devpts",
			Source:      "devpts",
			Options:     []string{"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"},
		},
		{
			Destination: "/dev/shm",
			Type:        "tmpfs",
			Source:      "shm",
			Options:     []string{"nosuid", "noexec", "nodev", "mode=1777", "size=65536k"},
		},
		{
			Destination: "/dev/mqueue",
			Type:        "mqueue",
			Source:      "mqueue",
			Options:     []string{"nosuid", "noexec", "nodev"},
		},
		{
			Destination: "/sys",
			Type:        "sysfs",
			Source:      "sysfs",
			Options:     []string{"nosuid", "noexec", "nodev", "ro"},
		},
		{
			Destination: "/run",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options:     []string{"nosuid", "strictatime", "mode=755", "size=65536k"},
		},
	}
}

func containerLinuxSpec(ns string, id string) *specs.Linux {
	memoryLimit := int64(536870912)
	return &specs.Linux{
		MaskedPaths: []string{
			"/proc/kcore",
			"/proc/latency_stats",
			"/proc/timer_list",
			"/proc/timer_stats",
			"/proc/sched_debug",
			"/sys/firmware",
			"/proc/scsi",
		},
		ReadonlyPaths: []string{
			"/proc/asound",
			"/proc/bus",
			"/proc/fs",
			"/proc/irq",
			"/proc/sys",
			"/proc/sysrq-trigger",
		},
		CgroupsPath: filepath.Join("/", ns, id),
		Resources: &specs.LinuxResources{
			Memory: &specs.LinuxMemory{
				Limit: &memoryLimit,
			},
			Devices: []specs.LinuxDeviceCgroup{
				{
					Allow:  false,
					Access: "rwm",
				},
			},
		},
		Namespaces: []specs.LinuxNamespace{
			{
				Type: specs.PIDNamespace,
			},
			{
				Type: specs.IPCNamespace,
			},
			{
				Type: specs.UTSNamespace,
			},
			{
				Type: specs.MountNamespace,
			},
			{
				Type: specs.NetworkNamespace,
			},
		},
	}
}
