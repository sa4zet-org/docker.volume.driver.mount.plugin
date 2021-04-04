package main

import (
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"os"
	"os/exec"
	"path"
	"sync"
)

type MountVolumeContext struct {
	fstype  string
	source  string
	target  string
	options string
	status  map[string]interface{}
}

type MountVolumeDriver struct {
	mu     sync.Mutex
	mounts map[string]*MountVolumeContext
}

func newMountVolumeDriver() *MountVolumeDriver {
	return &MountVolumeDriver{
		mounts: make(map[string]*MountVolumeContext),
	}
}

// Get the list of capabilities the driver supports.
// Supported scopes are global and local
// Scope allows cluster managers to handle the volume in different ways.
// For instance, a scope of global, signals to the cluster manager that it
// only needs to create the volume once instead of on each Docker host.
func (driver *MountVolumeDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "local"}}
}

// Instruct the plugin that the user wants to create a volume, given a user specified volume name.
// The plugin does not need to actually manifest the volume on the filesystem yet (until Mount is called).
// Opts is a map of driver specific options passed through from the user request.
func (driver *MountVolumeDriver) Create(createRequest *volume.CreateRequest) error {
	_, present := driver.mounts[createRequest.Name]
	if present {
		return fmt.Errorf(createRequest.Name + " already exists!")
	}
	target := path.Join("/mnt/docker.volume.driver.mounts", createRequest.Name)
	err := os.MkdirAll(target, 0755)
	if err != nil {
		return err
	}
	fstype, fstypePresent := createRequest.Options["fstype"]
	if !fstypePresent {
		return fmt.Errorf("fstype is required")
	}
	source, sourcePresent := createRequest.Options["source"]
	if !sourcePresent {
		return fmt.Errorf("source is required")
	}
	options, optionsPresent := createRequest.Options["options"]
	if !optionsPresent {
		return fmt.Errorf("options is required")
	}
	status := make(map[string]interface{})
	status["state"] = "created"
	driver.mu.Lock()
	driver.mounts[createRequest.Name] = &MountVolumeContext{
		fstype:  fstype,
		source:  source,
		target:  target,
		options: options,
		status:  status,
	}
	driver.mu.Unlock()
	return nil
}

func (driver *MountVolumeDriver) Get(getRequest *volume.GetRequest) (*volume.GetResponse, error) {
	driver.mu.Lock()
	mount, present := driver.mounts[getRequest.Name]
	driver.mu.Unlock()
	if !present {
		return nil, fmt.Errorf(getRequest.Name + " not found")
	}
	return &volume.GetResponse{
		Volume: &volume.Volume{
			Name:       getRequest.Name,
			Mountpoint: mount.target,
			Status:     mount.status,
		},
	}, nil
}

func (driver *MountVolumeDriver) List() (*volume.ListResponse, error) {
	var volumes []*volume.Volume
	driver.mu.Lock()
	for name, mount := range driver.mounts {
		volumes = append(volumes, &volume.Volume{
			Name:       name,
			Mountpoint: mount.target,
			Status:     mount.status,
		})
	}
	driver.mu.Unlock()
	return &volume.ListResponse{
		Volumes: volumes,
	}, nil
}

func (driver *MountVolumeDriver) Path(pathRequest *volume.PathRequest) (*volume.PathResponse, error) {
	driver.mu.Lock()
	mount, present := driver.mounts[pathRequest.Name]
	driver.mu.Unlock()
	if !present {
		return nil, fmt.Errorf(pathRequest.Name + " not found")
	}
	return &volume.PathResponse{
		Mountpoint: mount.target,
	}, nil
}

func (driver *MountVolumeDriver) Mount(mountRequest *volume.MountRequest) (*volume.MountResponse, error) {
	driver.mu.Lock()
	mount, present := driver.mounts[mountRequest.Name]
	if !present {
		driver.mu.Unlock()
		return nil, fmt.Errorf(mountRequest.Name + " not found")
	}
	if out, err := exec.Command(
		"mount",
		"-t",
		mount.fstype,
		"-o",
		mount.options,
		mount.source,
		mount.target,
	).CombinedOutput(); err != nil {
		mount.status["state"] = string(out)
		driver.mu.Unlock()
		return nil, err
	}
	mount.status["state"] = "mounted"
	driver.mu.Unlock()
	return &volume.MountResponse{
		Mountpoint: mount.target,
	}, nil
}

func (driver *MountVolumeDriver) Unmount(unmountRequest *volume.UnmountRequest) error {
	driver.mu.Lock()
	mount, present := driver.mounts[unmountRequest.Name]
	if !present {
		driver.mu.Unlock()
		return fmt.Errorf(unmountRequest.Name + " not found")
	}
	if out, err := exec.Command(
		"umount",
		"-f",
		mount.target,
	).CombinedOutput(); err != nil {
		mount.status["state"] = string(out)
		driver.mu.Unlock()
		return err
	}
	mount.status["state"] = "unmounted"
	driver.mu.Unlock()
	return nil
}

func (driver *MountVolumeDriver) Remove(removeRequest *volume.RemoveRequest) error {
	driver.mu.Lock()
	mount, present := driver.mounts[removeRequest.Name]
	if !present {
		driver.mu.Unlock()
		return fmt.Errorf(removeRequest.Name + " not found")
	}
	delete(driver.mounts, removeRequest.Name)
	err := os.Remove(mount.target)
	if err != nil {
		mount.status["state"] = "remove error"
		driver.mu.Unlock()
		return err
	}
	driver.mu.Unlock()
	return nil
}
