package util

import (
	"github.com/mndrix/tap-go"
	rspec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/opencontainers/runtime-tools/cgroups"
)

// ValidateLinuxResourcesMemory validates linux.resources.memory.
func ValidateLinuxResourcesMemory(config *rspec.Spec, t *tap.T, state *rspec.State) error {
	cg, err := cgroups.FindCgroup()
	t.Ok((err == nil), "find memory cgroup")
	if err != nil {
		t.Diagnostic(err.Error())
		return nil
	}

	lm, err := cg.GetMemoryData(state.Pid, config.Linux.CgroupsPath)
	t.Ok((err == nil), "get memory cgroup data")
	if err != nil {
		t.Diagnostic(err.Error())
		return nil
	}

	memory := config.Linux.Resources.Memory
	checkOptionalValue(t, "memory limit", memory.Limit, lm.Limit)
	checkOptionalValue(t, "memory reservation", memory.Reservation, lm.Reservation)

	if _, ok := cg.(*cgroups.CgroupV2); ok {
		return nil
	}

	checkOptionalValue(t, "memory swap", memory.Swap, lm.Swap)
	checkOptionalValue(t, "memory kernel", memory.Kernel, lm.Kernel) //nolint:staticcheck // Ignore SA1019: memory.Kernel is deprecated
	checkOptionalValue(t, "memory kernelTCP", memory.KernelTCP, lm.KernelTCP)
	checkOptionalValue(t, "memory swappiness", memory.Swappiness, lm.Swappiness)
	checkOptionalValue(t, "memory oom", memory.DisableOOMKiller, lm.DisableOOMKiller)

	return nil
}
