package cgroups

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	rspec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/opencontainers/runtime-tools/specerror"
)

// CgroupV2 used for cgroupv2 validation
type CgroupV2 struct {
	MountPath string
}

// GetBlockIOData gets cgroup blockio data
func (cg *CgroupV2) GetBlockIOData(pid int, cgPath string) (*rspec.LinuxBlockIO, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

// GetCPUData gets cgroup cpus data
func (cg *CgroupV2) GetCPUData(pid int, cgPath string) (*rspec.LinuxCPU, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

// GetDevicesData gets cgroup devices data
func (cg *CgroupV2) GetDevicesData(pid int, cgPath string) ([]rspec.LinuxDeviceCgroup, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

// GetHugepageLimitData gets cgroup hugetlb data
func (cg *CgroupV2) GetHugepageLimitData(pid int, cgPath string) ([]rspec.LinuxHugepageLimit, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

// GetNetworkData gets cgroup network data
func (cg *CgroupV2) GetNetworkData(pid int, cgPath string) (*rspec.LinuxNetwork, error) {
	return nil, fmt.Errorf("unimplemented yet")
}

// GetPidsData gets cgroup pid ints data
func (cg *CgroupV2) GetPidsData(pid int, cgPath string) (*rspec.LinuxPids, error) {
	contents, err := cg.readControllerFile(pid, cgPath, "pids.max")
	if err != nil {
		return nil, err
	}

	res, err := parseCgroup2Int(contents)
	if err != nil {
		return nil, err
	}

	return &rspec.LinuxPids{Limit: &res}, nil
}

func (cg *CgroupV2) GetMemoryData(pid int, cgPath string) (*rspec.LinuxMemory, error) {
	lm := &rspec.LinuxMemory{}

	contents, err := cg.readControllerFile(pid, cgPath, "memory.max")
	if err != nil {
		return nil, err
	}
	limit, err := parseCgroup2Int(contents)
	if err != nil {
		return nil, err
	}
	lm.Limit = &limit

	contents, err = cg.readControllerFile(pid, cgPath, "memory.low")
	if err != nil {
		return nil, err
	}
	reservation, err := parseCgroup2Int(contents)
	if err != nil {
		return nil, err
	}
	lm.Reservation = &reservation

	return lm, nil
}

func (cg *CgroupV2) readControllerFile(pid int, cgPath, fileName string) ([]byte, error) {
	dirPath, err := cg.resolveCgroupPath(pid, cgPath)
	if err != nil {
		return nil, err
	}

	contents, err := os.ReadFile(filepath.Join(dirPath, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, specerror.NewError(specerror.CgroupsPathAttach, fmt.Errorf("The runtime MUST consistently attach to the same place in the cgroups hierarchy given the same value of `cgroupsPath`"), rspec.Version)
		}
		return nil, err
	}

	return contents, nil
}

func (cg *CgroupV2) resolveCgroupPath(pid int, cgPath string) (string, error) {
	if filepath.IsAbs(cgPath) {
		path := filepath.Join(cg.MountPath, cgPath)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return "", specerror.NewError(specerror.CgroupsAbsPathRelToMount, fmt.Errorf("In the case of an absolute path, the runtime MUST take the path to be relative to the cgroups mount point"), rspec.Version)
			}
			return "", err
		}
		return path, nil
	}

	subPath, err := getUnifiedPath(pid)
	if err != nil {
		return "", err
	}
	if !hasCgroupPathSuffix(subPath, cgPath) {
		return "", fmt.Errorf("cgroup subsystem %s is not mounted as expected", "unified")
	}

	return filepath.Join(cg.MountPath, subPath), nil
}

func parseCgroup2Int(contents []byte) (int64, error) {
	trimmed := strings.TrimSpace(string(contents))
	if trimmed == "max" {
		return -1, nil
	}

	return strconv.ParseInt(trimmed, 10, 64)
}

func hasCgroupPathSuffix(fullPath, cgPath string) bool {
	if cgPath == "" {
		return true
	}

	fullPath = filepath.Clean(fullPath)
	cgPath = filepath.Clean(cgPath)

	if fullPath == cgPath {
		return true
	}

	return strings.HasSuffix(fullPath, string(filepath.Separator)+cgPath)
}
